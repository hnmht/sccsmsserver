package pg

import (
	"time"

	"go.uber.org/zap"
)

// Role Struct. A role is a collection of users with the same attributes
type Role struct {
	ID          int32     `db:"id" json:"ID"`
	Name        string    `db:"rolename" json:"name"  binding:"required"`
	Description string    `db:"description" json:"description" `
	SystemFlag  int16     `db:"systemflag" json:"systemFlag" `
	AllUserFlag int16     `db:"alluserflag" json:"allUserFlag"`
	Member      []Person  `json:"member"`
	CreateTime  time.Time `db:"createtime" json:"createTime"`
	Creator     Person    `db:"creatorid" json:"creator"`
	ModifyTime  time.Time `db:"modifytime" json:"modifyTime"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	Dr          int16     `db:"dr" json:"dr"`
	Ts          time.Time `json:"ts"`
}

// Initialize the sysrole table
func initSysrole() (isFinish bool, err error) {
	// 检查role表中是否已经存在预置数据systemadmin
	sqlStr := "select count(id) as rownum from sysrole where id=10000"
	hasRecord, isFinish, err := genericCheckRecord("sysrole", sqlStr)
	if hasRecord || !isFinish || err != nil {
		return
	}
	//插入sysAdmin角色
	sqlStr = "insert into sysrole(id,rolename,description,systemflag,alluserflag) values(10000,'systemadmin','系统预置角色',1,0)"
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("检查角色表中插入预置数据sysadmin出现错误", zap.Error(err))
		return isFinish, err
	}

	//2.3.2 检查role表中是否已经存在预置数据public
	sqlStr = "select count(id) as rownum from sysrole where id=10001"
	hasRecord, isFinish, err = genericCheckRecord("sysrole public", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//插入public角色
	sqlStr = "insert into sysrole(id,rolename,description,systemflag,alluserflag) values(10001,'public','系统预置角色',1,1)"
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("检查角色表中插入预置数据public出现错误", zap.Error(err))
		return isFinish, err
	}
	return
}
