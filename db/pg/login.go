package pg

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sccsmsserver/cache"
	"sccsmsserver/i18n"
	"sccsmsserver/pkg/jwt"
	"sccsmsserver/pkg/mysf"
	"sccsmsserver/pkg/security"
	"sccsmsserver/pub"
	"sccsmsserver/setting"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// User login parameters
type ParamLogin struct {
	UserCode   string `json:"userCode" binding:"required"`
	Password   string `json:"password" binding:"required"`
	ClientIP   string `json:"clientIp"`
	ClientType string `json:"clientType"`
	UserAgent  string `json:"userAgent"`
}

// User login failure record struct
type UserLoginFault struct {
	UserID    int32     `db:"user_id" json:"userID"`
	UserCode  string    `db:"usercode" json:"userCode"`
	ClientIp  string    `db:"clientip" json:"clientIp"`
	UserAgent string    `db:"useragent" json:"useragent"`
	Type      int16     `db:"type" json:"type"` //0 default 1 invalid: password 2 User does not exist
	Ts        time.Time `db:"ts" json:"ts"`
}

// IP address lockout struct
type IpLock struct {
	ClientIp  string    `json:"clientIp"`
	StartTime time.Time `json:"startTime"`
}

// User login
func Login(p *ParamLogin) (resStatus i18n.ResKey, token string, err error) {
	token = ""
	// RSA decrypt the password field
	op, err := base64.StdEncoding.DecodeString(p.Password)
	if err != nil {
		resStatus = i18n.CodeInternalError
		zap.L().Error("Login base64.StdEncoding.DecodeString failed", zap.Error(err))
		return
	}
	oriPassword, err := security.ScRsa.Decrypt(op)
	if err != nil {
		resStatus = i18n.CodeInternalError
		zap.L().Error("Login security.ScRsa.Decrypt(op) failed", zap.Error(err))
		return
	}
	p.Password = string(oriPassword)

	user := &User{
		Code:     p.UserCode,
		Password: p.Password,
	}
	// Record the password used during user login
	oPassword := user.Password

	// Check if the user exists.
	sqlStr := "select id,password,status,locked from sysuser where code = $1 and dr=0 and isoperator=1 limit 1"
	err = db.QueryRow(sqlStr, user.Code).Scan(&user.ID, &user.Password, &user.Status, &user.Locked)
	if err != nil && err != sql.ErrNoRows {
		resStatus = i18n.CodeInternalError
		zap.L().Error("Login query user information failed:", zap.Error(err))
		return
	}
	if err == sql.ErrNoRows {
		resStatus = i18n.StatusUserNotExist
		// Process when the user doesn't exist
		var ulf UserLoginFault
		ulf.UserID = 0
		ulf.UserCode = p.UserCode
		ulf.ClientIp = p.ClientIP
		ulf.UserAgent = p.UserAgent
		ulf.Type = 2
		ulf.process()
		return
	}

	// Check if the user is disabled.
	if user.Status != 0 {
		resStatus = i18n.StatusUserDisabled
		return
	}
	// Check if the user is locked
	if user.Locked != 0 {
		resStatus = i18n.StatusUserLocked
		return
	}

	// Check if the password matchs.
	password := encryptPassword(oPassword)
	result := strings.EqualFold(user.Password, password)

	if !result {
		resStatus = i18n.StatusInvalidPassword
		// Process if the passwords don't match
		var ulfp UserLoginFault
		ulfp.UserID = user.ID
		ulfp.UserCode = user.Code
		ulfp.ClientIp = p.ClientIP
		ulfp.UserAgent = p.UserAgent
		ulfp.Type = 1
		ulfp.process()
		return
	}

	// Generate Token ID
	tokenID := strconv.FormatInt(mysf.GenID(), 10)
	// Generate Json Web Token (JWT)
	token, expireTime, err := jwt.GenToken(user.ID, user.Code, tokenID)
	if err != nil {
		resStatus = i18n.CodeInternalError
		zap.L().Error("Login GenToken failed", zap.Error(err))
		return
	}

	// Get Person information
	person := Person{
		ID: user.ID,
	}
	resStatus, err = person.GetPersonInfoByID()
	if err != nil {
		return
	}

	// Add the user to online users.
	ou := OnlineUser{
		User:       person,
		TokenID:    tokenID,
		ClientType: p.ClientType,
		FromIp:     p.ClientIP,
		ExpireTime: expireTime,
	}

	resStatus, err = ou.Add()
	if err != nil {
		return
	}
	zap.L().Info("login Success", zap.String("user:", p.UserCode))

	resStatus = i18n.StatusOK
	return
}

// Add a new user login request failure record.
func (ulf *UserLoginFault) add() (err error) {
	sqlStr := `insert into sysloginfault(userid,usercode,clientip,useragent,type) 
	values($1,$2,$3,$4,$5)`
	_, err = db.Exec(sqlStr, &ulf.UserID, &ulf.UserCode, &ulf.ClientIp, &ulf.UserAgent, &ulf.Type)
	if err != nil {
		zap.L().Error("UserLoginFault.Add db.Exec failed", zap.Error(err))
		return
	}
	return
}

// Handle user login failure
func (ulf *UserLoginFault) process() (err error) {
	// Write login failure record.
	err = ulf.add()
	if err != nil {
		zap.L().Error("UserLoginFault Process AddLog failed", zap.Error(err))
		return
	}
	// Check login failure type,
	// if value is 1, indicating an incorrect password, the handle the password error.
	if ulf.Type == 1 {
		ulf.handlingInvalidPassword()
		return
	}
	// A value of 2 indicates the user doesn't exist.
	if ulf.Type == 2 {
		ulf.handlingUserNotExist()
		return
	}

	return
}

// User login password error handling
func (ulf *UserLoginFault) handlingInvalidPassword() (err error) {
	// Query user password error count within 30 minutes
	var pwdFaultNum int32
	sqlStr := `select count(id) as faultnum from sysloginfault 
	where ts > (current_timestamp - interval '30 minutes') 
	and type = 1 
	and user_id = $1  `
	err = db.QueryRow(sqlStr, &ulf.UserID).Scan(&pwdFaultNum)
	if err != nil {
		zap.L().Error("UserLoginFault TreatmentInvalidPassword QueryRow failed", zap.Error(err))
		return
	}
	// If the invalid password count exceeds the threshold, lock the user.
	if pwdFaultNum > setting.Conf.UserLockTh {
		lockUser(ulf.UserID)
	}
	return
}

// Lock the user
func lockUser(userID int32) (err error) {
	sqlStr := `update sysuser set locked=1,ts=current_timestamp where id=$1 and dr=0`
	_, err = db.Exec(sqlStr, &userID)
	if err != nil {
		zap.L().Error("LockUser db.exec failed", zap.Error(err))
		return
	}
	return
}

// User not exist error handling
func (ulf *UserLoginFault) handlingUserNotExist() (err error) {
	var userNotExistNum int32
	sqlStr := `select count(id) as falultnum from sysloginfault 
	where ts > (current_timestamp - interval '30 minutes') 
	and type=2 and clientip=$1`
	err = db.QueryRow(sqlStr, ulf.ClientIp).Scan(&userNotExistNum)
	if err != nil {
		zap.L().Error("UserLoginFault TreatmentUserNotExist", zap.Error(err))
		return
	}
	if userNotExistNum > setting.Conf.IpLockTh {
		// Write the IP lockout record to cache.
		key := fmt.Sprintf("%s%s%s", pub.IPBlack, ":", ulf.ClientIp)
		l := IpLock{ulf.ClientIp, time.Now()}
		jsonL, _ := json.Marshal(l)
		err = cache.SetOther(key, jsonL)
	}
	return
}
