package pg

import (
	"database/sql"
	"sccsmsserver/i18n"
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
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Creator     Person    `db:"creatorid" json:"creator"`
	ModifyDate  time.Time `db:"modifytime" json:"modifyDate"`
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

// Get role list
func GetRoles() (roles []Role, resStatus i18n.ResKey, err error) {
	roles = make([]Role, 0)
	// Retrieve from sysrole table
	sqlStr := `select a.id, a.name,a.description,a.systemflag,a.alluserflag,
	a.createtime,a.creatorid,a.modifytime,a.modifierid,a.dr,a.ts 
	from sysrole a
	where a.dr = 0
	order by a.name`
	rows, err := db.Query(sqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetRoles db.Query failed:", zap.Error(err))
		return
	}
	defer rows.Close()
	// Extract data from database query results
	for rows.Next() {
		var role Role
		err = rows.Scan(&role.ID, &role.Name, &role.Description, &role.SystemFlag, &role.AllUserFlag,
			&role.CreateDate, &role.Creator.ID, &role.ModifyDate, &role.Modifier.ID, &role.Dr, &role.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetRoles rows.next failed", zap.Error(err))
			return
		}
		// Get creator details
		if role.Creator.ID > 0 {
			resStatus, err = role.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get modifier details
		if role.Modifier.ID > 0 {
			resStatus, err = role.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get role members
		resStatus, err = role.GetMembers()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		// Add to roles slice
		roles = append(roles, role)
	}
	return
}

// Get Role members
func (role *Role) GetMembers() (resStatus i18n.ResKey, err error) {
	role.Member = make([]Person, 0)
	resStatus = i18n.StatusOK
	// Retrieve data from database
	sqlStr := `select a.id,a.code,a.name,a.fileid,a.deptid,
	COALESCE((select b.code from department b where b.id = a.deptid),'') as deptcode,
	COALESCE((select b.name from department b where b.id = a.deptid),'') as deptname,
	COALESCE(a.description,'') as description,
	COALESCE(a.mobile,'') as mobile,
	COALESCE(a.email,'') as email,
	a.gender,a.systemflag,a.status,a.createtime,a.ts,a.dr from sysuser a
	where a.dr=0  and a.id in (select c.userid from sysuserrole c where roleid = $1)`
	rows, err := db.Query(sqlStr, role.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		resStatus = i18n.StatusInternalError
		zap.L().Error("role.GetMembers db.Qurey failed", zap.Error(err))
		return
	}
	defer rows.Close()
	// Extract data from database query results
	for rows.Next() {
		var p Person
		err = rows.Scan(&p.ID, &p.Code, &p.Name, &p.Avatar.ID, &p.DeptID,
			&p.DeptCode, &p.DeptName, &p.Description, &p.Mobile, &p.Email,
			&p.Gender, &p.SystemFlag, &p.Status, &p.CreateDate, &p.Ts, &p.Dr)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetUserByRoleID rows.Next() failed", zap.Error(err))
			return
		}
		// Get Avatar details
		if p.Avatar.ID > 0 {
			resStatus, err = p.Avatar.GetFileInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Add to member slice
		role.Member = append(role.Member, p)
	}
	return
}
