package pg

import (
	"crypto/md5"
	"encoding/hex"
	"sccsmsserver/pub"
	"time"

	"go.uber.org/zap"
)

// User Master Data
type User struct {
	ID          int32       `db:"id" json:"id"`
	Code        string      `db:"code" json:"code"`
	Name        string      `db:"name" json:"name"`
	Password    string      `db:"password" json:"password"`
	Mobile      string      `db:"mobile" json:"mobile"`
	Email       string      `db:"email" json:"email"`
	IsOperator  int16       `db:"isoperator" json:"isOperator"`
	Position    Position    `db:"positionid" json:"position"`
	Avatar      File        `db:"fileid" json:"avatar"`
	Dept        SimpDept    `db:"deptid" json:"department"`
	Description string      `db:"description" json:"description"`
	Gender      int16       `db:"gender" json:"gender"`
	Locked      int16       `db:"locked" json:"locked"`
	Status      int16       `db:"status" json:"status"`
	SystemFlag  int16       `db:"systemflag" json:"systemFlag"`
	MenuList    SystemMenus `json:"menuList"`
	Roles       []Role      `json:"roles"`
	Person      Person      `json:"person"`
	CreateDate  time.Time   `db:"createtime" json:"createDate"`
	Creator     Person      `db:"creatorid" json:"creator"`
	ModifyDate  time.Time   `db:"modifytime" json:"modifyDate"`
	Modifier    Person      `db:"modifierid" json:"modifier"`
	Dr          int16       `db:"dr" json:"dr"`
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
	sqlStr = `insert into sysuser(id,name,password,createtime,description,
		systemflag,code,creatorid) 
		values(10000,'admin',$1,now(),'System default',
		1,'admin',10000)`
	_, err = db.Exec(sqlStr, encryptPassword(pub.DefaultPassword))
	if err != nil {
		isFinish = false
		zap.L().Error("initSysUser db.Exec failed:", zap.Error(err))
		return isFinish, err
	}
	return
}

// encryptPassword 加密密码
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(pub.Secret))

	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
