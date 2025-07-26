package pg

import (
	"time"

	"go.uber.org/zap"
)

// Role Struct. A role is a collection of users with the same attributes
type Role struct {
	ID          int32     `db:"id" json:"ID"`
	Name        string    `db:"name" json:"name"  binding:"required"`
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
	// Step 1: Check if a record exists for the system default role 'systemadmin' in the sysrole table.
	sqlStr := "select count(id) as rownum from sysrole where id=10000 and dr=0"
	hasRecord, isFinish, err := genericCheckRecord("sysrole", sqlStr)
	// Step 2: Exit if the record exists or an error occurs
	if hasRecord || !isFinish || err != nil {
		return
	}
	// Step 3: Insert a record for the system default role 'systemadmin' into the sysrole table.
	sqlStr = `insert into sysrole(id,name,description,systemflag,alluserflag) 
	values(10000,'systemadmin','System default',1,0)`
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysrole insert systemadmin db.Exec failed:", zap.Error(err))
		return isFinish, err
	}

	// Step 4: Check if a record exists for the system default role 'public' in the sysrole table
	sqlStr = "select count(id) as rownum from sysrole where id=10001 and dr=0"
	hasRecord, isFinish, err = genericCheckRecord("sysrole public", sqlStr)
	// Step 5: Exit if the record exists or an error occurs
	if hasRecord || !isFinish || err != nil {
		return
	}
	// Step 6: Insert a record for the system default role 'public' into the sysrole table.
	sqlStr = `insert into sysrole(id,name,description,systemflag,alluserflag) 
	values(10001,'public','system default',1,1)`
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysrole Insert public db.Exec failed: ", zap.Error(err))
		return isFinish, err
	}
	return
}

// Initialize the sysuserrole table
func initSysUserRole() (isFinish bool, err error) {
	// Step 1: Check if a record exists for the sysuserrole table.
	sqlStr := "select count(id) as rownum from sysuserrole where userid=10000"
	// Step 2: Exit if the record exists or an error occurs,
	hasRecord, isFinish, err := genericCheckRecord("sysuserrole", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	// Step 3: Insert the 'admin' user and 'systemadmin' role mapping into the sysuserrole table
	sqlStr = "insert into sysuserrole(userid,roleid,ts) values(10000,10000,now())"
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysUserRole inserting the admin and sysadmin mapping failed:", zap.Error(err))
		return isFinish, err
	}
	// Step 4: Insert the 'admin' user and 'public' role mapping into the sysuserrole table.
	sqlStr = "insert into sysuserrole(userid,roleid,ts) values(10000,10001,now())"
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysuserRole inserting the admin and public mapping failed:", zap.Error(err))
		return isFinish, err
	}
	return
}

// Initialize the sysrolemenu table
func initSysRoleMenu() (isFinish bool, err error) {
	// Step 1: Check if a record exists for the sysrolemenu table
	sqlStr := "select count(id) as rownum from sysrolemenu where roleid=10000"
	hasRecord, isFinish, err := genericCheckRecord("sysrolemenu", sqlStr)
	// Step 2: Exit if the record exists or an error occurs.
	if hasRecord || !isFinish || err != nil {
		return
	}
	// Setp 3: Insert all menus for the systemadmin role into the sysrolemenu table.
	// This means the systemadmin role will have all menu permissions.
	sqlStr = `insert into sysrolemenu(roleid,menuid,selected,indeterminate) 
	values(10000,$1,true,false)`
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysRoleMenu while preparing to write systemadmin role data failed:", zap.Error(err))
		return isFinish, err
	}
	for _, menu := range SysFunctionList {
		_, err = stmt.Exec(menu.ID)
		if err != nil {
			isFinish = false
			zap.L().Error("initSysRoleMenu Failed to write systemadmin role menu "+string(menu.Title)+" to the sysrolemenu table:", zap.Error(err))
			return isFinish, err
		}
	}
	stmt.Close()
	// Step 3: Insert menus for the public role into the sysrolemenu table.
	sqlStr = `insert into sysrolemenu(roleid,menuid,selected,indeterminate) 
	values(10001,$1,true,false)`
	stmt, err = db.Prepare(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysRoleMenu while preparing to write public role data failed:", zap.Error(err))
		return isFinish, err
	}
	for _, menu := range PublicFunctionList {
		_, err = stmt.Exec(menu.ID)
		if err != nil {
			isFinish = false
			zap.L().Error("initSysRoleMenu Failed to write public role menu "+string(menu.Title)+"to the sysrolemenu table:", zap.Error(err))
			return isFinish, err
		}
	}
	stmt.Close()
	return
}
