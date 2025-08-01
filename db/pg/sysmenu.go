package pg

import (
	"sccsmsserver/i18n"

	"go.uber.org/zap"
)

// System menu struct
type SystemMenu struct {
	ID             int32       `json:"id"`
	FatherID       int32       `json:"fatherID"`
	Title          i18n.ResKey `json:"title"`
	Path           string      `json:"path"`
	Icon           string      `json:"icon"`
	Component      string      `json:"component"`
	Selected       bool        `json:"selected"`
	Indeterminate  bool        `json:"indeterminate"`
	AddFromVersion string      `json:"addFromVersion"`
}

// Menu object slice
type SystemMenus []SystemMenu

// Menu object struct
type MenuItem struct {
	SystemMenu
	Children []MenuItem `json:"children"`
}

// initialize sysmenu table
func initSysMenu() (isFinish bool, err error) {
	// Step 1: Check if a record exists for the sysmenu table
	sqlStr := "select count(id) as rownum from sysmenu where dr=0"
	hasRecord, isFinish, err := genericCheckRecord("sysmenu", sqlStr)
	// Step 2: Exit if the record exists or an error occurs
	if hasRecord || !isFinish || err != nil {
		return
	}
	// Step 3: Prepare to insert system menus into the sysmenu table.
	sqlStr = `insert into sysmenu(id,fatherid,title,path,icon,
		component,selected,indeterminate) 
		values($1,$2,$3,$4,$5,
		$6,$7,$8)`
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysMenu db.Prepare failed:", zap.Error(err))
		return isFinish, err
	}
	defer stmt.Close()
	// Step 4: Write system menu data entry by entry
	for _, menu := range SysFunctionList {
		_, err = stmt.Exec(menu.ID, menu.FatherID, menu.Title, menu.Path, menu.Icon, menu.Component, menu.Selected, menu.Indeterminate)
		if err != nil {
			isFinish = false
			zap.L().Error("initSysMenu Failed to write the "+string(menu.Title)+" menu to the sysmenu table:", zap.Error(err))
			return isFinish, err
		}
	}
	return
}

// Generate Menu Tree
func (m *SystemMenus) ProcessToTree(pid int32, level int32) []MenuItem {
	var menuTree []MenuItem
	if level == 10 {
		return menuTree
	}

	list := m.FindChildren(pid)

	if len(list) == 0 {
		return menuTree
	}

	for _, v := range list {
		child := m.ProcessToTree(v.ID, level+1)
		menuTree = append(menuTree, MenuItem{v, child})
	}

	return menuTree
}

// Find submenus
func (m *SystemMenus) FindChildren(pid int32) []SystemMenu {
	child := []SystemMenu{}
	for _, v := range *m {
		if v.FatherID == pid {
			child = append(child, v)
		}
	}
	return child
}

// Get menu tree
func GetMenuTree() (trees []MenuItem, err error) {
	systemMenus, err := GetMenuList()
	trees = systemMenus.ProcessToTree(0, 0)
	return
}

// Get menu list from the database
func GetMenuList() (menus SystemMenus, err error) {
	sqlStr := "select id,fatherid, title,path,icon,component from sysmenu"
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetMenuList db.Query failed:", zap.Error(err))
		return
	}
	var menu SystemMenu
	for rows.Next() {
		if err = rows.Scan(&menu.ID, &menu.FatherID, &menu.Title, &menu.Path, &menu.Icon, &menu.Component); err != nil {
			zap.L().Error("GetMenuList rows.Scan failed:", zap.Error(err))
			return
		}
		menus = append(menus, menu)
	}
	return
}
