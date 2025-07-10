package pg

import (
	"sccsmsserver/i18n"

	"go.uber.org/zap"
)

// System menu struct
type SystemMenu struct {
	ID             int32       `json:"ID"`
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
