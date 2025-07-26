package pg

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
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

// // UserInfo 获取用户信息
// func UserInfo(userID int32) (user User, resStatus i18n.ResKey, err error) {
// 	user.ID = userID
// 	resStatus, err = GetUserInfoByID(&user)
// 	if resStatus != i18n.StatusOK || err != nil {
// 		return
// 	}
// 	return
// }

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
	return i18n.StatusOK, nil
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
