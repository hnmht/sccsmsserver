package pg

import "go.uber.org/zap"

// System functionality list
var SysFunctionList SystemMenus = SystemMenus{
	SystemMenu{ID: 1, FatherID: 0, Title: "MenuDashboard", Path: "/private/dashboard", Icon: "Home", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 10, FatherID: 0, Title: "MenuCalendar", Path: "/private/calendar", Icon: "CalendarMonth", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 15, FatherID: 0, Title: "MenuMessage", Path: "/private/message", Icon: "Message", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 20, FatherID: 0, Title: "MenuAddressBook", Path: "/private/addressBook", Icon: "ContactPhone", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 30, FatherID: 0, Title: "现场管理", Path: "/private/scene", Icon: "Streetview", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 110, FatherID: 30, Title: "指令单", Path: "/private/workOrder/workOrderDoc", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 210, FatherID: 30, Title: "执行单", Path: "/private/execute/executeDoc", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 220, FatherID: 30, Title: "执行单审阅", Path: "/private/execute/executeDocReview", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 310, FatherID: 30, Title: "问题处理单", Path: "/private/problem/disposeDoc", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 410, FatherID: 30, Title: "指令单执行统计", Path: "/private/reports/workOrderStat", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 420, FatherID: 30, Title: "执行单统计", Path: "/private/reports/executeDocStat", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 430, FatherID: 30, Title: "问题处理单统计", Path: "/private/reports/problemDisposeStat", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 500, FatherID: 0, Title: "文档管理", Path: "/private/document", Icon: "Inventory", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 510, FatherID: 500, Title: "文档类别", Path: "/private/document/class", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 520, FatherID: 500, Title: "上传文档", Path: "/private/document/upload", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 530, FatherID: 500, Title: "查阅文档", Path: "/private/document/lookup", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 600, FatherID: 0, Title: "培训管理", Path: "/private/train", Icon: "School", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 610, FatherID: 600, Title: "培训课程", Path: "/private/train/course", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 620, FatherID: 600, Title: "培训记录", Path: "/private/train/record", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 630, FatherID: 600, Title: "授课查询", Path: "/private/train/givelessons", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 640, FatherID: 600, Title: "受训查询", Path: "/private/train/receivetraining", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 700, FatherID: 0, Title: "劳保管理", Path: "/private/lpa", Icon: "Masks", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 710, FatherID: 700, Title: "岗位定额", Path: "/private/lpa/quota", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 720, FatherID: 700, Title: "发放向导", Path: "/private/lpa/wizard", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 730, FatherID: 700, Title: "发放单", Path: "/private/lpa/issuedvoucher", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 740, FatherID: 700, Title: "发放查询", Path: "/private/lpa/issuedquery", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 1000, FatherID: 0, Title: "档案", Path: "/private/archive", Icon: "Article", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1010, FatherID: 1000, Title: "部门档案", Path: "/private/archive/department", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1011, FatherID: 1000, Title: "岗位档案", Path: "/private/archive/operatingpost", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 1016, FatherID: 1000, Title: "现场档案类别", Path: "/private/archive/sceneItemClass", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1020, FatherID: 1000, Title: "现场档案", Path: "/private/archive/sceneItem", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1030, FatherID: 1000, Title: "自定义档案类别", Path: "/private/archive/userDefineClass", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1040, FatherID: 1000, Title: "自定义档案", Path: "/private/archive/userDefine", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1050, FatherID: 1000, Title: "执行项目类别", Path: "/private/archive/exectiveItemClass", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1060, FatherID: 1000, Title: "执行项目", Path: "/private/archive/exectiveItem", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1070, FatherID: 1000, Title: "风险等级", Path: "/private/archive/riskLevel", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.1.0"},
	SystemMenu{ID: 1080, FatherID: 1000, Title: "劳保用品档案", Path: "/private/archive/laborProtection", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 1100, FatherID: 0, Title: "模板", Path: "/private/template", Icon: "FormatListNumbered", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1110, FatherID: 1100, Title: "执行模板", Path: "/private/template/execItemTemplate", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9000, FatherID: 0, Title: "权限", Path: "/private/permission", Icon: "People", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9010, FatherID: 9000, Title: "角色管理", Path: "/private/permission/role", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9020, FatherID: 9000, Title: "用户管理", Path: "/private/permission/user", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9030, FatherID: 9000, Title: "权限分配", Path: "/private/permission/permissionAssignment", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9040, FatherID: 9000, Title: "在线用户", Path: "/private/permission/onlineUser", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.2.0"},
	SystemMenu{ID: 9100, FatherID: 0, Title: "设置", Path: "/private/options", Icon: "Settings", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9110, FatherID: 9100, Title: "现场档案自定义项", Path: "/private/options/sceneItemOption", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9120, FatherID: 9100, Title: "授权许可", Path: "/private/options/register", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9130, FatherID: 9100, Title: "首页定义", Path: "/private/options/landingPageSetup", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.1.0"},
	SystemMenu{ID: 9910, FatherID: 0, Title: "个人中心", Path: "/private/my/profile", Icon: "ManageAccounts", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9920, FatherID: 0, Title: "关于", Path: "/private/my/about", Icon: "Info", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
}

// Default User Permissions List
var PublicFunctionList SystemMenus = SystemMenus{
	SystemMenu{ID: 1, FatherID: 0, Title: "首页", Path: "/private/dashboard", Icon: "Home", Component: "", Selected: true, Indeterminate: false},
	SystemMenu{ID: 10, FatherID: 0, Title: "日程", Path: "/private/calendar", Icon: "CalendarMonth", Component: "", Selected: false, Indeterminate: false},
	SystemMenu{ID: 15, FatherID: 0, Title: "消息", Path: "/private/message", Icon: "Message", Component: "", Selected: false, Indeterminate: false},
	SystemMenu{ID: 9910, FatherID: 0, Title: "个人中心", Path: "/private/my/profile", Icon: "ManageAccounts", Component: "", Selected: true, Indeterminate: false},
	SystemMenu{ID: 9920, FatherID: 0, Title: "关于", Path: "/private/my/about", Icon: "Info", Component: "", Selected: true, Indeterminate: false},
}

// Create database table
func createTable() (isFinish bool, err error) {
	var rowNum int
	var sqlStr string
	isFinish = true
	// Step 1: create table
	for _, table := range tables {
		// Check if table exists
		sqlStr = "select count(tablename) from pg_tables where tablename=$1"
		err = db.QueryRow(sqlStr, table.TableName).Scan(&rowNum)
		if err != nil {
			isFinish = false
			zap.L().Error("createTable check table "+table.TableName+" exist failed", zap.Error(err))
			return
		}

		if rowNum <= 0 {
			_, err = db.Exec(table.CreateSQL)
			if err != nil {
				isFinish = false
				zap.L().Error(" createTable Create table "+table.TableName+" failed", zap.Error(err))
				return
			}
			zap.L().Info("Table " + table.TableName + " created successfully.")
		} else {
			zap.L().Warn("Table " + table.TableName + " already exists.")
		}
	}
	zap.L().Info("All database tables created successfully.")
	//Setp 2: Initialize table record
	for _, table := range tables {
		isFinish, err = table.InitFunc()
		if err != nil {
			return
		}
		zap.L().Info("Table " + table.TableName + " initialized successfully.")
	}
	// Step 3: Modify the isfinish value in the sysinfo tableW
	sqlStr = "update sysinfo set endtime=now(), isfinish=TRUE"
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("createTable db.Exec update the isfinish failed:", zap.Error(err))
		return isFinish, err
	}
	zap.L().Info("All database tables initialized successfully.")
	return
}

// Check if the database initialization is complete
func checkDbInit() (isFinish bool, err error) {
	var rowNum int
	// Check if the sysinfo table exists
	sqlStr := "select count(tablename) as rownum from pg_tables where tablename='sysinfo'"
	err = db.QueryRow(sqlStr).Scan(&rowNum)

	// If an error occurs, return incomplete
	if err != nil {
		zap.L().Error("checkDbInit db.QueryRow check sysinfo table exists failed:", zap.Error(err))
		return false, err
	}

	// If the sysinfo table does not exist, return incomplete
	if rowNum <= 0 {
		return false, nil
	}

	// if the sysinfo table exist, proceed to check if it contains any data
	sqlStr = "select count(isFinish) from sysinfo"
	err = db.QueryRow(sqlStr).Scan(&rowNum)
	if err != nil {
		zap.L().Error("checkDbInit db.QueryRow check sysinfo row count failed:", zap.Error(err))
		return false, err
	}
	// If the sysinfo tablse has less than or equal to zero rows,or more than one row,return incomplete
	if rowNum <= 0 || rowNum > 1 {
		return false, nil
	}
	// If the sysinfo table has exactly one row, the proceed to check the content of the isfinish field
	sqlStr = "select isfinish from sysinfo"
	err = db.QueryRow(sqlStr).Scan(&isFinish)
	if err != nil {
		zap.L().Error("checkDbInit db.QueryRow query sysinfo table isfinish field failed:", zap.Error(err))
		return false, err
	}
	return isFinish, nil
}
