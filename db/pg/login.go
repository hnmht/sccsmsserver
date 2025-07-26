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

// ParamLogin 用户登录请求参数
type ParamLogin struct {
	UserCode   string `json:"usercode" binding:"required"`
	Password   string `json:"password" binding:"required"`
	ClientIP   string `json:"clientip"`
	ClientType string `json:"clientType"`
	UserAgent  string `json:"useragent"`
}

// UserLoginFault 用户登录请求失败记录
type UserLoginFault struct {
	UserID    int32     `db:"user_id" json:"id"`
	UserCode  string    `db:"usercode" json:"code"`
	ClientIp  string    `db:"clientip" json:"clientip"`
	UserAgent string    `db:"useragent" json:"useragent"`
	Type      int16     `db:"type" json:"type"` //0 默认 1 密码错误 2 用户名不存在
	Ts        time.Time `db:"ts" json:"ts"`
}

// IPLock IP锁定结构体
type IpLock struct {
	ClientIp  string    `json:"clientip"`
	StartTime time.Time `json:"starttime"`
}

// Login 用户登录
func Login(p *ParamLogin) (resStatus i18n.ResKey, token string, err error) {
	token = ""
	//解密rsaPassword
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
	oPassword := user.Password //用户登录时的密码

	//检查用户是否存在
	sqlStr := "select id,password,status,locked from sysuser where usercode = $1 and dr=0 and isoperator=1 limit 1"
	err = db.QueryRow(sqlStr, user.Code).Scan(&user.ID, &user.Password, &user.Status, &user.Locked)
	if err != nil && err != sql.ErrNoRows {
		resStatus = i18n.CodeInternalError
		zap.L().Error("Login query user information failed:", zap.Error(err))
		return
	}
	if err == sql.ErrNoRows {
		resStatus = i18n.StatusUserNotExist
		//用户不存在处理
		var ulf UserLoginFault
		ulf.UserID = 0
		ulf.UserCode = p.UserCode
		ulf.ClientIp = p.ClientIP
		ulf.UserAgent = p.UserAgent
		ulf.Type = 2
		ulf.Process()
		return
	}

	//检查用户数量是否超过最大授权数
	var userNumber int32
	sqlStr = `select count(id) as usernumber from sysuser where dr=0 and status=0 and systemflag=0 and isoperator=1`
	err = db.QueryRow(sqlStr).Scan(&userNumber)
	if err != nil {
		resStatus = i18n.CodeInternalError
		zap.L().Error("Login check userNumber db.QueryRow(sqlStr) failed", zap.Error(err))
		return
	}

	//检查用户是否停用
	if user.Status != 0 {
		resStatus = i18n.StatusUserDisabled
		return
	}
	//检查用户是否锁定
	if user.Locked != 0 {
		resStatus = i18n.StatusUserLocked
		return
	}

	//检查密码是否匹配
	password := EncryptPassword(oPassword)
	result := strings.EqualFold(user.Password, password)

	if !result {
		resStatus = i18n.StatusInvalidPassword
		//用户输入密码错误处理
		var ulfp UserLoginFault
		ulfp.UserID = user.ID
		ulfp.UserCode = user.Code
		ulfp.ClientIp = p.ClientIP
		ulfp.UserAgent = p.UserAgent
		ulfp.Type = 1
		ulfp.Process()
		return
	}

	//生成tokenID
	tokenID := strconv.FormatInt(mysf.GenID(), 10)
	//生成JWT
	token, expireTime, err := jwt.GenToken(user.ID, user.Code, tokenID)
	if err != nil {
		resStatus = i18n.CodeInternalError
		zap.L().Error("Login GenToken failed", zap.Error(err))
		return
	}

	person := Person{
		ID: user.ID,
	}

	//获取人员信息
	/* resStatus, err = person.GetPersonInfoByID()
	if err != nil {
		return
	} */

	//增加在线用户
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
	return
}

// UserLoginFault AddLog 增加用户登录请求失败记录
func (ulf *UserLoginFault) AddLog() (err error) {
	sqlStr := `insert into sysloginfault(user_id,usercode,clientip,useragent,type) values($1,$2,$3,$4,$5)`
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		zap.L().Error("UserLoginFault AddLog perpare failed", zap.Error(err))
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(&ulf.UserID, &ulf.UserCode, &ulf.ClientIp, &ulf.UserAgent, &ulf.Type)
	if err != nil {
		zap.L().Error("UserLoginFault AddLog exec failed", zap.Error(err))
		return
	}
	return
}

// UserLoginFalut Process 用户登录失败处理
func (ulf *UserLoginFault) Process() (err error) {
	//1 写入日志
	err = ulf.AddLog()
	if err != nil {
		zap.L().Error("UserLoginFault Process AddLog failed", zap.Error(err))
		return
	}
	//2 检查登录失败类型，如果为1 密码错误 则进行密码错误处理
	if ulf.Type == 1 {
		ulf.TreatmentInvalidPassword()
		return
	}
	//3 如果30分钟内同一ip用户名不存在
	if ulf.Type == 2 {
		ulf.TreatmentUserNotExist()
		return
	}

	return
}

// UserLoginFault TreatmentInvalidPassword
func (ulf *UserLoginFault) TreatmentInvalidPassword() (err error) {
	//查询30分钟内用户密码错误次数
	var pwdFaultNum int32
	sqlStr := `select count(id) as faultnum from sysloginfault where ts > (current_timestamp - interval '30 minutes') and type = 1 and  user_id = $1  `
	err = db.QueryRow(sqlStr, &ulf.UserID).Scan(&pwdFaultNum)
	if err != nil {
		zap.L().Error("UserLoginFault TreatmentInvalidPassword QueryRow failed", zap.Error(err))
		return
	}
	//如果密码错误次数大于阈值,则锁定用户
	if pwdFaultNum > setting.Conf.UserLockTh {
		LockUser(ulf.UserID)
	}
	return
}

// 锁定用户
func LockUser(userID int32) (err error) {
	sqlStr := `update sysuser set locked = 1 ,ts=current_timestamp where id=$1 and dr=0`
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		zap.L().Error("LockUser db.prepare failed", zap.Error(err))
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(&userID)
	if err != nil {
		zap.L().Error("LockUser stmt.exec failed", zap.Error(err))
		return
	}
	return
}

// UserLoginFault TreatmentUserNotExist 用户不存在处理
func (ulf *UserLoginFault) TreatmentUserNotExist() (err error) {
	var userNotExistNum int32
	sqlStr := `select count(id) as falultnum from sysloginfault where ts > (current_timestamp - interval '30 minutes') and type=2 and clientip=$1`
	err = db.QueryRow(sqlStr, ulf.ClientIp).Scan(&userNotExistNum)
	if err != nil {
		zap.L().Error("UserLoginFault TreatmentUserNotExist", zap.Error(err))
		return
	}
	if userNotExistNum > setting.Conf.IpLockTh {
		//将ip锁定记录写入缓存
		key := fmt.Sprintf("%s%s%s", pub.IPBlack, ":", ulf.ClientIp)
		l := IpLock{ulf.ClientIp, time.Now()}
		jsonL, _ := json.Marshal(l)
		err = cache.SetOther(key, jsonL)
		// blacklist.AddBlackListItem(ulf.ClientIp)
	}
	return
}
