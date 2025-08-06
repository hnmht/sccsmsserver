package pg

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"sccsmsserver/cache"
	"sccsmsserver/i18n"
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

// User password change struct
type ParamChangePwd struct {
	UserID        int32  `json:"id" binding:"required"`
	UserCode      string `json:"code"`
	UserName      string `json:"name"`
	Password      string `json:"password" binding:"required"`
	NewPassword   string `json:"newpwd" binding:"required"`
	ConfirmNewPwd string `json:"confirmnewpwd" binding:"required"`
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

// Encrypt Password
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write(pub.Md5Secret)
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// Get User Information by ID
func (user *User) GetUserInfoByID() (resStatus i18n.ResKey, err error) {
	// Query user information from the database.
	sqlStr := `select code,name,mobile,email,fileid,
		isoperator,positionid,deptid,COALESCE(description,''),gender,
		locked,systemflag,createtime,creatorid,modifytime,
		modifierid,ts,dr 
		from sysuser where id = $1`
	err = db.QueryRow(sqlStr, user.ID).Scan(&user.Code, &user.Name, &user.Mobile, &user.Email, &user.Avatar.ID,
		&user.IsOperator, &user.Position.ID, &user.Dept.ID, &user.Description, &user.Gender,
		&user.Locked, &user.SystemFlag, &user.CreateDate, &user.Creator.ID, &user.ModifyDate,
		&user.Modifier.ID, &user.Ts, &user.Dr)
	if err != nil && err != sql.ErrNoRows {
		resStatus = i18n.StatusInternalError
		zap.L().Error("dap.GetUserInfoByID failed", zap.Error(err))
		return
	}
	if err == sql.ErrNoRows {
		resStatus = i18n.StatusUserNotExist
		return
	}

	// Get user menu list.
	resStatus, err = user.GetUserMenusByID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Get user's assigned roles
	resStatus, err = user.GetUserRolesByID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Get user avatar.
	if user.Avatar.ID > 0 {
		resStatus, err = user.Avatar.GetFileInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Person infomation corresponding to the user.
	user.Person.ID = user.ID
	if user.Person.ID > 0 {
		resStatus, err = user.Person.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Postion detail.
	if user.Position.ID > 0 {
		resStatus, err = user.Position.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Department detail.
	if user.Dept.ID > 0 {
		resStatus, err = user.Dept.GetSimpDeptInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	return i18n.StatusOK, nil
}

// Get User Menu list by User ID
func (u *User) GetUserMenusByID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Query menu list from database
	sqlStr := `select id,fatherid,title,path,icon,
	component 
	from sysmenu 
	where id in (select menuid from sysrolemenu where roleid in (select roleid from sysuserrole where userid=$1)) 
	order by id`
	rows, err := db.Query(sqlStr, u.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("User.GetUserMenusByID db.Query failed", zap.Error(err))
		return
	}
	defer rows.Close()

	var menus SystemMenus
	for rows.Next() {
		var menu SystemMenu
		err = rows.Scan(&menu.ID, &menu.FatherID, &menu.Title, &menu.Path, &menu.Icon,
			&menu.Component)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("User.GetUserMenusByID rows.Next() failed", zap.Error(err))
			return
		}
		menus = append(menus, menu)
	}
	u.MenuList = menus
	return
}

// Get user assigned roles
func (u *User) GetUserRolesByID() (resStatus i18n.ResKey, err error) {
	// Get User Assigned roles from database
	sqlStr := `select id,name,description,systemflag,alluserflag 
	from sysrole 
	where id in (select roleid from sysuserrole where userid = $1)`
	rows, err := db.Query(sqlStr, u.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("User.GetUserRoleByID db.Query faield", zap.Error(err))
		return
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		err = rows.Scan(&role.ID, &role.Name, &role.Description, &role.SystemFlag, &role.AllUserFlag)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("User.GetUserRoleByID rows.next() failed", zap.Error(err))
			return
		}
		roles = append(roles, role)
	}
	u.Roles = roles
	return i18n.StatusOK, nil
}

// Check if the user name exists.
func (user *User) CheckUserNameExist() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var count int32
	sqlStr := "select count(id) from sysuser where  dr=0 and name = $1 "
	err = db.QueryRow(sqlStr, user.Name).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("User.CheckUserNameExist db.QueryRow failed", zap.Error(err))
		return
	}
	if count > 0 {
		resStatus = i18n.StatusUserNameExist
	}
	return
}

// Check if the user code exists.
func (user *User) CheckUserCodeExist() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var count int32
	sqlStr := "select count(id) from sysuser where dr=0 and code = $1 and id <> $2"
	err = db.QueryRow(sqlStr, user.Code, user.ID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("User.CheckUserCodeExist db.QueryRow failed", zap.Error(err))
		return
	}
	if count > 0 {
		resStatus = i18n.StatusUserCodeExist
	}
	return
}

// Add user
func (user *User) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Encrypt the password field useing MD5
	user.Password = encryptPassword(user.Password)
	// Check if the user code exists.
	resStatus, err = user.CheckUserCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Begin a database transaction.
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("User.Add db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Insert a record to sysuser table.
	sqlStr1 := `insert into 
	sysuser(code,name,password,mobile,email,
		isoperator,positionid,fileid,deptid,description,
		gender,status,locked,creatorid) 
		values($1,$2,$3,$4,$5,
		$6,$7,$8,$9,$10,
		$11,$12,$13,$14) returning id`

	err = db.QueryRow(sqlStr1,
		user.Code, user.Name, user.Password, user.Mobile, user.Email,
		user.IsOperator, user.Position.ID, user.Avatar.ID, user.Dept.ID, user.Description,
		user.Gender, user.Status, user.Locked, user.Creator.ID).Scan(&user.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("User.Add db.QueryRow failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Pre-processing for Insert records into the sysuserrole table.
	sqlStr2 := "insert into sysuserrole(userid,roleid) values($1,$2)"
	stmt2, err := tx.Prepare(sqlStr2)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("User.Add tx.Prepare(sqlStr2) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer stmt2.Close()

	for _, role := range user.Roles {
		_, err = stmt2.Exec(user.ID, role.ID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("User.Add stmt2.Exec failed", zap.Error(err))
			tx.Rollback()
			return
		}
	}
	return
}

// Edit user
func (user *User) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	//Encrypt the password field using MD5
	if user.Password != "" {
		user.Password = encryptPassword(user.Password)
	}
	// Check if the user code exists.
	resStatus, err = user.CheckUserCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Begin a database transaction.
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("User.Edit db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Update the record in the sysuser table.
	var uSql string
	if user.Password != "" {
		uSql = `update sysuser set usercode=$1,username=$2,mobile=$3, email=$4, file_id=$5,
		isoperator=$6,op_id=$7,dept_id=$8,description=$9, gender=$10,
		status=$11,locked=$12,modify_time=now(), modifyuserid = $13,ts=current_timestamp,password=$14 
		where id=$15 and ts=$16 and dr=0`
	} else {
		uSql = `update sysuser set usercode=$1,username=$2,mobile=$3, email=$4, file_id=$5,
		isoperator=$6,op_id=$7,dept_id=$8,description=$9, gender=$10,
		status=$11,locked=$12,modify_time=now(), modifyuserid=$13,ts=current_timestamp 
		where id=$14 and ts=$15 and dr=0`
	}

	var res sql.Result
	if user.Password != "" {
		res, err = db.Exec(uSql, user.Code, user.Name, user.Mobile, user.Email, user.Avatar.ID,
			user.IsOperator, user.Position.ID, user.Dept.ID, user.Description, user.Gender,
			user.Status, user.Locked, user.Modifier.ID, user.Password,
			user.ID, user.Ts)
	} else {
		res, err = db.Exec(uSql, user.Code, user.Name, user.Mobile, user.Email, user.Avatar.ID,
			user.IsOperator, user.Position.ID, user.Dept.ID, user.Description, user.Gender,
			user.Status, user.Locked, user.Modifier.ID,
			user.ID, user.Ts)
	}
	// Check the number of rows affected by SQL update operation.
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("User.Edit db.Exec(uSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	updateNum, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("User.Edit uStmt.exec() res.RowsAffected() failed", zap.Error(err))
		tx.Rollback()
		return
	}

	// If the update operation affects equals zero,
	// it indicates that another user has already updated that row.
	if updateNum == 0 {
		resStatus = i18n.StatusOtherEdit
		zap.L().Info("User.Edit other edit")
		_ = tx.Rollback()
		return
	}

	// Delete existing records from the sysuserrole table.
	dSql := "delete from sysuserrole where userid=$1"
	dStmt, err := tx.Prepare(dSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("User.Edit delete sysuserrole tx.Prepare(dSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer dStmt.Close()

	_, err = dStmt.Exec(user.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("User.Edit delete antduserrole exec failed", zap.Error(err))
		tx.Rollback()
		return
	}

	// Insert new records into the sysuserole table.
	iSql := "insert into sysuserrole(user_id,role_id,ts) values($1,$2,now())"
	insertStmt, err := tx.Prepare(iSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("User.Edit insert into sysuserrole tx.prepare(isql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer insertStmt.Close()

	for _, item := range user.Roles {
		_, err = insertStmt.Exec(user.ID, item.ID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("User.Edit into sysuserrole insertStmt.Exec failed", zap.Error(err))
			tx.Rollback()
			return
		}
	}
	// Delete the user from local cache
	user.DelFromLocalCache()
	return i18n.StatusOK, nil
}

// User updates via personal center.
func (user *User) ModifyProfile() (resStatus i18n.ResKey, err error) {
	//Check if the opertor and user are the same person.
	if user.ID != user.Modifier.ID {
		resStatus = i18n.StatusProfileOnlySelf
		return
	}
	// Update the record in the sysuser table.
	sqlStr := `update sysuser set mobile=$1, email=$2, fileid=$3,description=$4, modifytime=current_timestamp,modifierid=$5,
	ts=current_timestamp where id=$6 and dr=0 and ts=$7 returning modifytime,ts,modifierid`
	err = db.QueryRow(sqlStr, user.Mobile, user.Email, user.Avatar.ID, user.Description, user.Modifier.ID,
		user.ID, user.Ts).Scan(&user.ModifyDate, &user.Ts, &user.Modifier.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			resStatus = i18n.StatusOtherEdit
			return
		}
		resStatus = i18n.StatusInternalError
		zap.L().Error("User.ModityProfile db.QueryRow failed", zap.Error(err))
		return
	}
	// Delete from the local cache
	user.DelFromLocalCache()
	return i18n.StatusOK, nil
}

// Get user list
func GetUsers() (users []User, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	users = make([]User, 0)
	sqlStr := `select id,code,name, COALESCE(mobile,'') as mobile,COALESCE(email,'') as email,
	fileid,isoperator,positionid,deptid,COALESCE(description,''),
	gender,status,locked,systemflag,createtime,
	creatorid,modifytime,modifierid,dr,ts 
	from sysuser where dr=0`
	rows, err := db.Query(sqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetUsers db.Query failed", zap.Error(err))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Code, &user.Name, &user.Mobile, &user.Email,
			&user.Avatar.ID, &user.IsOperator, &user.Position.ID, &user.Dept.ID, &user.Description,
			&user.Gender, &user.Status, &user.Locked, &user.SystemFlag, &user.CreateDate,
			&user.Creator.ID, &user.ModifyDate, &user.Modifier.ID, &user.Dr, &user.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetUsers row.Next() failed", zap.Error(err))
			return
		}
		// Get user assigned roles
		resStatus, err = user.GetUserRolesByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}

		// Get user avatar detail.
		if user.Avatar.ID > 0 {
			resStatus, err = user.Avatar.GetFileInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		// Get Position detail.
		if user.Position.ID > 0 {
			resStatus, err = user.Position.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		// Get Department detail.
		if user.Dept.ID > 0 {
			resStatus, err = user.Dept.GetSimpDeptInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		// Get Creator deatil.
		if user.Creator.ID > 0 {
			resStatus, err = user.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		// Get modifier detail.
		if user.Modifier.ID > 0 {
			resStatus, err = user.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		users = append(users, user)
	}

	return
}

// Delete user
func (user *User) Delete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the user id referenced.
	resStatus, err = user.CheckIsUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update the sysuser table with a deletion flag.
	sqlStr := `update sysuser set dr=1,ts=current_timestamp,modifierid=$1,modifytime=current_timestamp 
	where id = $2 and ts=$3 and dr=0`
	res, err := db.Exec(sqlStr, user.Modifier.ID, user.ID, user.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("User.delete db.Exec failed", zap.Error(err))
		return
	}
	// Check the number of rows affected by the SQL update operation.
	delNum, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteUsers Code"+user.Code+" res.RowsAffected() failed", zap.Error(err))
		return
	}
	// If the update operation affects equals zero,
	// it indicates that another user has already updated that row.
	if delNum == 0 {
		resStatus = i18n.StatusOtherEdit
		zap.L().Info("DeleteUsers Code" + user.Code + " failed,the doc has been refreshed")
		return
	}
	// Delete from local cache
	user.DelFromLocalCache()
	return
}

// Batch delete users
func DeleteUsers(users *[]User, modifyUserID int32) (resStatus i18n.ResKey, err error) {
	// Begin a database transaction.
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteUsers db.begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Pre-Processing for update the sysuser table deletion flag.
	sqlStr := "update sysuser set dr = 1,ts = current_timestamp,modifyuserid=$1,modify_time=current_timestamp where id = $2 and ts=$3 and dr=0"
	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteUsers tx.prepare failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer stmt.Close()

	for _, u := range *users {
		// Check if the user id is refrenced
		resStatus, err = u.CheckIsUsed()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		// Write data
		res, err := stmt.Exec(modifyUserID, u.ID, u.Ts)
		if err != nil {
			zap.L().Error("DeleteUsers Code:"+u.Code+"delete failed", zap.Error(err))
			tx.Rollback()
			return i18n.StatusInternalError, err
		}
		// Check the number of rows affected by the write operation.
		delNum, err := res.RowsAffected()
		if err != nil {
			zap.L().Error("DeleteUsers Code"+u.Code+" res.RowsAffected() failed", zap.Error(err))
			tx.Rollback()
			return i18n.StatusInternalError, err
		}
		// if the number of affected rows is less than one,
		// it indicates that another user has already updated that row.
		if delNum < 1 {
			zap.L().Info("DeleteUsers Code" + u.Code + "  failed,the doc has been refreshed")
			tx.Rollback()
			return i18n.StatusOtherEdit, err
		}

		// Delete from local cache
		u.DelFromLocalCache()
	}

	return i18n.StatusOK, nil
}

// Delete Person Archive from local Cache
func (user *User) DelFromLocalCache() {
	number, _, _ := cache.Get(pub.Person, user.ID)
	if number > 0 {
		cache.Del(pub.Person, user.ID)
	}
}

// Check if the user id is refrenced.
func (user *User) CheckIsUsed() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Define the items to be checked.
	checkItems := []ArchiveCheckUsed{
		{
			Description:    "Refrenced by the role creator",
			SqlStr:         "select count(id) from sysrole where dr=0 and creatorid=$1",
			UsedReturnCode: i18n.StatusRoleCreateUsed,
		},
		{
			Description:    "Refrenced by the role modifier",
			SqlStr:         "select count(id) from sysrole where dr = 0 and modifierid=$1",
			UsedReturnCode: i18n.StatusRoleModifyUsed,
		},
		{
			Description:    "Refrenced by the user creator",
			SqlStr:         "select count(id) from sysuser where dr = 0 and creatorid=$1",
			UsedReturnCode: i18n.StatusUserCreateUsed,
		},
		{
			Description:    "Refrenced by the user modifier",
			SqlStr:         "select count(id) from sysuser where dr = 0 and modifierid=$1",
			UsedReturnCode: i18n.StatusUserModifyUsed,
		},
		/* 		{
		   			Description:    "被指令单创建人引用",
		   			SqlStr:         "select count(id) from workorder_h where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusWOCreateUsed,
		   		},
		   		{
		   			Description:    "被指令单修改人引用",
		   			SqlStr:         "select count(id) from workorder_h where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusWOModifyUsed,
		   		},
		   		{
		   			Description:    "被指令单确认人引用",
		   			SqlStr:         "select count(id) from workorder_h where dr = 0 and confirmuserid=$1",
		   			UsedReturnCode: i18n.StatusWOConfirmUsed,
		   		},
		   		{
		   			Description:    "被执行单创建人引用",
		   			SqlStr:         "select count(id) from executedoc_h where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusEOCreateUsed,
		   		},
		   		{
		   			Description:    "被执行单修改人引用",
		   			SqlStr:         "select count(id) from executedoc_h where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusEOModifyUsed,
		   		},
		   		{
		   			Description:    "被执行单确认人引用",
		   			SqlStr:         "select count(id) from executedoc_h where dr = 0 and confirmuserid=$1",
		   			UsedReturnCode: i18n.StatusEOConfirmUsed,
		   		},
		   		{
		   			Description:    "被问题处理单创建人引用",
		   			SqlStr:         "select count(id) from disposedoc where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusIRFCreateUsed,
		   		},
		   		{
		   			Description:    "被问题处理单修改人引用",
		   			SqlStr:         "select count(id) from disposedoc where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusIRFModifyUsed,
		   		},
		   		{
		   			Description:    "被问题处理单确认人引用",
		   			SqlStr:         "select count(id) from disposedoc where dr = 0 and confirmuserid=$1",
		   			UsedReturnCode: i18n.StatusIRFConfirmUsed,
		   		},
		   		{
		   			Description:    "被文档类别创建人引用",
		   			SqlStr:         "select count(id) from documentclass where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusDCCreateUsed,
		   		},
		   		{
		   			Description:    "被文档类别修改人引用",
		   			SqlStr:         "select count(id) from documentclass where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusDCModifyUsed,
		   		},
		   		{
		   			Description:    "被文档创建人引用",
		   			SqlStr:         "select count(id) from document where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusDocumentCreateUsed,
		   		},
		   		{
		   			Description:    "被文档修改人引用",
		   			SqlStr:         "select count(id) from document where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusDocumentModifyUsed,
		   		},
		   		{
		   			Description:    "被培训课程创建人引用",
		   			SqlStr:         "select count(id) from traincourse where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusTCCreateUsed,
		   		},
		   		{
		   			Description:    "被培训课程修改人引用",
		   			SqlStr:         "select count(id) from traincourse where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusTCModifyUsed,
		   		},
		   		{
		   			Description:    "被培训记录主讲人引用",
		   			SqlStr:         "select count(id) from trainrecord_h where dr = 0 and lecturer_id=$1",
		   			UsedReturnCode: i18n.StatusTRLecturerUsed,
		   		},
		   		{
		   			Description:    "被培训记录学员引用",
		   			SqlStr:         "select count(id) from trainrecord_b where dr = 0 and student_id=$1",
		   			UsedReturnCode: i18n.StatusTRStudentUsed,
		   		},
		   		{
		   			Description:    "被培训记录创建人引用",
		   			SqlStr:         "select count(id) from trainrecord_h where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusTRCreateUsed,
		   		},
		   		{
		   			Description:    "被培训记录修改人引用",
		   			SqlStr:         "select count(id) from trainrecord_h where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusTRModifyUsed,
		   		},
		   		{
		   			Description:    "被培训记录确认人引用",
		   			SqlStr:         "select count(id) from trainrecord_h where dr = 0 and confirmuserid=$1",
		   			UsedReturnCode: i18n.StatusTRConfirmUsed,
		   		},
		   		{
		   			Description:    "被岗位定额创建人引用",
		   			SqlStr:         "select count(id) from lpaquota_h where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusPQCreateUsed,
		   		},
		   		{
		   			Description:    "被岗位定额修改人引用",
		   			SqlStr:         "select count(id) from lpaquota_h where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusPQModifyUsed,
		   		},
		   		{
		   			Description:    "被岗位定额确认人引用",
		   			SqlStr:         "select count(id) from lpaquota_h where dr = 0 and confirmuserid=$1",
		   			UsedReturnCode: i18n.StatusPQConfirmUsed,
		   		},
		   		{
		   			Description:    "被劳保用品发放单领用人引用",
		   			SqlStr:         "select count(id) from lpaissuedoc_b where dr = 0 and recipient_id=$1",
		   			UsedReturnCode: i18n.StatusPPEIFRecipientUsed,
		   		},
		   		{
		   			Description:    "被劳保用品发放单确认人引用",
		   			SqlStr:         "select count(id) from lpaissuedoc_h where dr = 0 and confirmuserid=$1",
		   			UsedReturnCode: i18n.StatusPPEIFConfirmUsed,
		   		},
		   		{
		   			Description:    "被劳保用品发放单修改人引用",
		   			SqlStr:         "select count(id) from lpaissuedoc_h where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusPPEIFModifyUsed,
		   		},
		   		{
		   			Description:    "被劳保用品发放单创建人引用",
		   			SqlStr:         "select count(id) from lpaissuedoc_h where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusPPEIFCreateUsed,
		   		},
		   		{
		   			Description:    "被岗位档案创建人引用",
		   			SqlStr:         "select count(id) from operatingpost where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusPositionCreateUsed,
		   		},
		   		{
		   			Description:    "被岗位档案修改人引用",
		   			SqlStr:         "select count(id) from operatingpost where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusPositionModifyUsed,
		   		},
		   		{
		   			Description:    "被现场档案类别创建人引用",
		   			SqlStr:         "select count(id) from sceneitemclass where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusCSCreateUsed,
		   		},
		   		{
		   			Description:    "被现场档案类别修改人引用",
		   			SqlStr:         "select count(id) from sceneitemclass where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusCSModifyUsed,
		   		},
		   		{
		   			Description:    "被现场档案创建人引用",
		   			SqlStr:         "select count(id) from sceneitem where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusCSCreateUsed,
		   		},
		   		{
		   			Description:    "被现场档案修改人引用",
		   			SqlStr:         "select count(id) from sceneitem where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusCSModifyUsed,
		   		},
		   		{
		   			Description:    "被自定义档案类别创建人引用",
		   			SqlStr:         "select count(id) from userdefineclass where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusUDCCreateUsed,
		   		},
		   		{
		   			Description:    "被自定义档案类别修改人引用",
		   			SqlStr:         "select count(id) from userdefineclass where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusUDCModifyUsed,
		   		},
		   		{
		   			Description:    "被自定义档案创建人引用",
		   			SqlStr:         "select count(id) from userdefinedoc where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusUDCreateUsed,
		   		},
		   		{
		   			Description:    "被自定义档案修改人引用",
		   			SqlStr:         "select count(id) from userdefinedoc where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusUDModifyUsed,
		   		},
		   		{
		   			Description:    "被执行项目类别创建人引用",
		   			SqlStr:         "select count(id) from exectiveitemclass where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusEPCCreateUsed,
		   		},
		   		{
		   			Description:    "被执行项目类别修改人引用",
		   			SqlStr:         "select count(id) from exectiveitemclass where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusEPCModifyUsed,
		   		},
		   		{
		   			Description:    "被执行项目创建人引用",
		   			SqlStr:         "select count(id) from exectiveitem where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusEPCreateUsed,
		   		},
		   		{
		   			Description:    "被执行项目修改人引用",
		   			SqlStr:         "select count(id) from exectiveitem where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusEPModifyUsed,
		   		},
		   		{
		   			Description:    "被风险等级创建人引用",
		   			SqlStr:         "select count(id) from risklevel where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusRLCreateUsed,
		   		},
		   		{
		   			Description:    "被风险等级修改人引用",
		   			SqlStr:         "select count(id) from risklevel where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusRLModifyUsed,
		   		},
		   		{
		   			Description:    "被劳保用品创建人引用",
		   			SqlStr:         "select count(id) from laborprotection where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusPPECreateUsed,
		   		},
		   		{
		   			Description:    "被劳保用品修改人引用",
		   			SqlStr:         "select count(id) from laborprotection where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusPPEModifyUsed,
		   		},
		   		{
		   			Description:    "被执行模板创建人引用",
		   			SqlStr:         "select count(id) from exectivetemplate_h where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusEPTCreateUsed,
		   		},
		   		{
		   			Description:    "被执行模板修改人引用",
		   			SqlStr:         "select count(id) from exectivetemplate_h where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusEPTModifyUsed,
		   		},

		   		{
		   			Description:    "被部门档案负责人引用",
		   			SqlStr:         "select count(id) from department where dr = 0 and leader=$1",
		   			UsedReturnCode: i18n.StatusDeptLeaderUsed,
		   		},
		   		{
		   			Description:    "被部门档案创建人引用",
		   			SqlStr:         "select count(id) from department where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusDeptCreateUsed,
		   		},
		   		{
		   			Description:    "被部门档案修改人引用",
		   			SqlStr:         "select count(id) from department where dr = 0 and modifyuserid=$1",
		   			UsedReturnCode: i18n.StatusDeptModifyUsed,
		   		},
		   		{
		   			Description:    "被现场档案负责人引用",
		   			SqlStr:         "select count(id) from sceneitem where dr = 0 and respperson_id=$1",
		   			UsedReturnCode: i18n.StatusCSRespUsed,
		   		},
		   		{
		   			Description:    "执行项目默认值引用",
		   			SqlStr:         `select count(id) as usednum from exectiveitem where resulttypeid = '510' and dr=0 and defaultvalue=cast($1 as varchar)`,
		   			UsedReturnCode: i18n.StatusEPDefaultUsed,
		   		},
		   		{
		   			Description:    "执行项目错误值引用",
		   			SqlStr:         `select count(id) as usednum from exectiveitem where resulttypeid = '510' and dr=0 and errorvalue=cast($1 as varchar)`,
		   			UsedReturnCode: i18n.StatusEPErrorUsed,
		   		},
		   		{
		   			Description:    "执行模板默认值引用",
		   			SqlStr:         `select count(id) from exectivetemplate_b where eid_id in (select id from exectiveitem where resulttypeid='510' and dr=0) and dr=0 and defaultvalue=CAST($1 as varchar)`,
		   			UsedReturnCode: i18n.StatusEPTDefaultUsed,
		   		},
		   		{
		   			Description:    "执行模板错误值引用",
		   			SqlStr:         `select count(id) from exectivetemplate_b where eid_id in (select id from exectiveitem where resulttypeid='510' and dr=0) and dr=0 and errorvalue=CAST($1 as varchar)`,
		   			UsedReturnCode: i18n.StatusEPTErrorUsed,
		   		},
		   		{
		   			Description:    "被指令单表体执行人引用",
		   			SqlStr:         "select count(id) from workorder_b where dr = 0 and ep_id=$1",
		   			UsedReturnCode: i18n.StatusWOEpUsed,
		   		},
		   		{
		   			Description:    "被执行单表头执行人引用",
		   			SqlStr:         "select count(id) from executedoc_h where dr = 0 and ep_id=$1",
		   			UsedReturnCode: i18n.StatusEOEpUsed,
		   		},
		   		{
		   			Description:    "被执行单表体处理人引用",
		   			SqlStr:         "select count(id) from executedoc_b where dr = 0 and hp_id=$1",
		   			UsedReturnCode: i18n.StatusEODpUsed,
		   		},
		   		{
		   			Description:    "被执行单执行值引用",
		   			SqlStr:         `select count(id) from executedoc_b where eid_id in (select id from exectiveitem where resulttypeid='510' and dr=0) and dr=0 and exectivevalue=CAST($1 as varchar)`,
		   			UsedReturnCode: i18n.StatusEOValueUsed,
		   		},
		   		{
		   			Description:    "被执行单错误值引用",
		   			SqlStr:         `select count(id) from executedoc_b where eid_id in (select id from exectiveitem where resulttypeid='510' and dr=0) and dr=0 and errorvalue=CAST($1 as varchar)`,
		   			UsedReturnCode: i18n.StatusEOErrorUsed,
		   		},
		   		{
		   			Description:    "被问题处理单执行人引用",
		   			SqlStr:         "select count(id) from disposedoc where dr = 0 and ep_id=$1",
		   			UsedReturnCode: i18n.StatusIRFEpUsed,
		   		},
		   		{
		   			Description:    "被问题处理单处理人引用",
		   			SqlStr:         "select count(id) from disposedoc where dr = 0 and dp_id=$1",
		   			UsedReturnCode: i18n.StatusIRFDpUsed,
		   		},
		   		{
		   			Description:    "被执行单批注接收人引用",
		   			SqlStr:         "select count(id) from executedoc_comment where dr = 0 and sendto_id=$1",
		   			UsedReturnCode: i18n.StatusEOCommentUsed,
		   		},
		   		{
		   			Description:    "被执行单批注人引用",
		   			SqlStr:         "select count(id) from executedoc_comment where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusEOCommentUsed,
		   		},
		   		{
		   			Description:    "被执行单审阅人引用",
		   			SqlStr:         "select count(id) from executedoc_review where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusEOCommentUsed,
		   		},
		   		{
		   			Description:    "被执行单审阅人引用",
		   			SqlStr:         "select count(id) from executedoc_review where dr = 0 and createuserid=$1",
		   			UsedReturnCode: i18n.StatusEOCommentUsed,
		   		}, */
	}

	// Check item by item
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, user.ID).Scan(&usedNum)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("User.CheckIsUsed "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return
}

// Change user password
func (pcp *ParamChangePwd) ChangePassword() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the values of the NewPassword and ConfirmNewPwd fields are the same.
	if pcp.NewPassword != pcp.ConfirmNewPwd {
		return i18n.StatusPasswordDisaccord, nil
	}
	// Check if the old password is correct
	var oldPassword string
	sqlStr := "select password from sysuser where id=$1"
	err = db.QueryRow(sqlStr, pcp.UserID).Scan(&oldPassword)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("ParmChangePwd.ChangePassword db.QueryRow failed:", zap.Error(err))
		resStatus = i18n.StatusErrorUnknow
		return
	}
	if err == sql.ErrNoRows {
		zap.L().Error("ParmChangePwd.ChangePassword db.QueryRow failed:", zap.Error(err))
		resStatus = i18n.StatusUserNotExist
		return
	}
	if oldPassword != encryptPassword(pcp.Password) {
		zap.L().Info("ParmChangePwd.ChangePassword invalid password.")
		resStatus = i18n.StatusInvalidPassword
		return
	}

	// change password
	sqlStr = "update sysuser set password=$1 where id=$2"
	newPwd := encryptPassword(pcp.NewPassword)
	_, err = db.Exec(newPwd, pcp.UserID)
	if err != nil {
		zap.L().Error("ChangePassword update exec failed:", zap.Error(err))
		resStatus = i18n.StatusErrorUnknow
		return
	}
	return
}
