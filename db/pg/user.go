package pg

import (
	"crypto/md5"
	"encoding/hex"
	"sccsmsserver/pub"
	"time"

	"go.uber.org/zap"
)

// User 定义用户对象
type User struct {
	UserID     int32  `db:"id" json:"id"`
	UserCode   string `db:"usercode" json:"code"`
	Username   string `db:"username" json:"name"`
	Password   string `db:"password" json:"password"`
	Mobile     string `db:"mobile" json:"mobile"`
	Email      string `db:"email" json:"email"`
	IsOperator int16  `db:"isoperator" json:"isoperator"`
	// OperatingPost OperatingPost `db:"op_id" json:"operatingpost"` //工作岗位
	Avatar File `db:"file_id" json:"avatar"`
	// Dept          SimpDept      `db:"dept_id" json:"department"`
	Description string      `db:"description" json:"description"`
	Gender      int16       `db:"gender" json:"gender"`
	Locked      int16       `db:"locked" json:"locked"`
	Status      int16       `db:"status" json:"status"`
	SystemFlag  int16       `db:"systemflag" json:"systemflag"`
	MenuList    SystemMenus `json:"menulist"`
	Roles       []Role      `json:"roles"`
	Person      Person      `json:"person"`
	CreateDate  time.Time   `db:"create_time" json:"createdate"`
	CreateUser  Person      `db:"createuserid" json:"createuser"`
	ModifyDate  time.Time   `db:"modify_time" json:"modifydate"`
	ModifyUser  Person      `db:"modifyuserid" json:"modifyuser"`
	Dr          int16       `db:"dr" json:"dr"` //删除标志
	Ts          time.Time   `db:"ts" json:"ts"`
}

// Initialize user table
func initSysUser() (isFinish bool, err error) {
	// Step 1: Check if a record exists for the system defaut user 'admin' in the sysuser table.
	sqlStr := "select count(id) as rownum from sysuser where id=10000 and dr=0"
	hasRecord, isFinish, err := genericCheckRecord("sysuser", sqlStr)
	// Step 2: Exit if the record exists or an error occurs
	if hasRecord || !isFinish || err != nil {
		return
	}
	// Step 3: Insert a record for the system default user 'admin' into the sysuser table.
	sqlStr = `insert into sysuser(id,username,password,createtime,description,
		systemflag,usercode,createuserid) 
		values(10000,'admin',$1,now(),'System default',
		1,'admin',10000)`
	_, err = db.Exec(sqlStr, EncryptPassword(pub.DefaultPassword))
	if err != nil {
		isFinish = false
		zap.L().Error("initSysUser db.Exec failed:", zap.Error(err))
		return isFinish, err
	}
	return
}

// encryptPassword 加密密码
func EncryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(pub.Secret))

	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
