package pg

import "go.uber.org/zap"

// System functionality list
var SysFunctionList SystemMenus = SystemMenus{
	SystemMenu{ID: 1, FatherID: 0, Title: "MenuDashboard", Path: "/private/dashboard", Icon: "Home", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 10, FatherID: 0, Title: "MenuCalendar", Path: "/private/calendar", Icon: "CalendarMonth", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 15, FatherID: 0, Title: "MenuMessage", Path: "/private/message", Icon: "Message", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 20, FatherID: 0, Title: "MenuAddressBook", Path: "/private/addressBook", Icon: "ContactPhone", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 30, FatherID: 0, Title: "MenuCSM", Path: "/private/constructionSiteManagement", Icon: "Streetview", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 110, FatherID: 30, Title: "MenuWO", Path: "/private/workOrder/workOrder", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 210, FatherID: 30, Title: "MenuEO", Path: "/private/execute/executionOrder", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 220, FatherID: 30, Title: "MenuEOReview", Path: "/private/execute/EOReview", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 310, FatherID: 30, Title: "MenuIRF", Path: "/private/problem/issueResolutionForm", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 410, FatherID: 30, Title: "MenuWOStatus", Path: "/private/reports/WOStatus", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 420, FatherID: 30, Title: "MenuEOStatus", Path: "/private/reports/EOStatus", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 430, FatherID: 30, Title: "MenuIRFStatus", Path: "/private/reports/IRFStatus", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 500, FatherID: 0, Title: "MenuDM", Path: "/private/documentManagement", Icon: "Inventory", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 510, FatherID: 500, Title: "MenuDC", Path: "/private/document/category", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 520, FatherID: 500, Title: "MenuDocumentUpload", Path: "/private/document/upload", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 530, FatherID: 500, Title: "MenuDocumentFind", Path: "/private/document/find", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 600, FatherID: 0, Title: "MenuTM", Path: "/private/trainingManagement", Icon: "School", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 610, FatherID: 600, Title: "MenuTC", Path: "/private/training/course", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 620, FatherID: 600, Title: "MenuTR", Path: "/private/training/record", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 630, FatherID: 600, Title: "MenuTS", Path: "/private/training/teachingStatistics", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 640, FatherID: 600, Title: "MenuTPS", Path: "/private/train/participationStatistics", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 700, FatherID: 0, Title: "MenuPPEM", Path: "/private/personalProtectiveEquipmentManagement", Icon: "Masks", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 710, FatherID: 700, Title: "MenuPQ", Path: "/private/ppe/quota", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 720, FatherID: 700, Title: "MenuPPEWizard", Path: "/private/ppe/wizard", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 730, FatherID: 700, Title: "MenuPPEIF", Path: "/private/ppe/ppeIssuanceForm", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 740, FatherID: 700, Title: "MenuPPES", Path: "/private/ppe/ppeStatistics", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1000, FatherID: 0, Title: "MenuMD", Path: "/private/masterData", Icon: "Article", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1010, FatherID: 1000, Title: "MenuDepartment", Path: "/private/masterData/department", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1011, FatherID: 1000, Title: "MenuPosition", Path: "/private/masterData/position", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1016, FatherID: 1000, Title: "MenuCSC", Path: "/private/masterData/constructionSiteCategory", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1020, FatherID: 1000, Title: "MenuCSA", Path: "/private/masterData/constructionSiteArchive", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1030, FatherID: 1000, Title: "MenuUDC", Path: "/private/masterData/userDefinedCategory", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1040, FatherID: 1000, Title: "MenuUDA", Path: "/private/masterData/userDefineArchive", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1050, FatherID: 1000, Title: "MenuEPC", Path: "/private/masterData/executionProjectCategory", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1060, FatherID: 1000, Title: "MenuEP", Path: "/private/masterData/executionProject", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1070, FatherID: 1000, Title: "MenuRL", Path: "/private/masterData/riskLevel", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.1.0"},
	SystemMenu{ID: 1080, FatherID: 1000, Title: "MenuPPE", Path: "/private/masterData/personalProtectiveEquipment", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1100, FatherID: 0, Title: "MenuTemplate", Path: "/private/template", Icon: "FormatListNumbered", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 1110, FatherID: 1100, Title: "MenuEPT", Path: "/private/template/executionProjectTemplate", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9000, FatherID: 0, Title: "MenuPermission", Path: "/private/permission", Icon: "People", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9010, FatherID: 9000, Title: "MenuRole", Path: "/private/permission/role", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9020, FatherID: 9000, Title: "MenuUser", Path: "/private/permission/user", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9030, FatherID: 9000, Title: "MenuPA", Path: "/private/permission/permissionAssignment", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9040, FatherID: 9000, Title: "MenuOU", Path: "/private/permission/onlineUser", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9100, FatherID: 0, Title: "MenuSettings", Path: "/private/options", Icon: "Settings", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9110, FatherID: 9100, Title: "MenuCSO", Path: "/private/options/constructionSiteOptions", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9130, FatherID: 9100, Title: "MenuLPS", Path: "/private/options/landingPageSetup", Icon: "", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.1.0"},
	SystemMenu{ID: 9910, FatherID: 0, Title: "MenuProfile", Path: "/private/my/profile", Icon: "ManageAccounts", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9920, FatherID: 0, Title: "MenuAbout", Path: "/private/my/about", Icon: "Info", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
}

// SeaCloud Data Type List
var ScDataTypeList map[int32]ScDataType = map[int32]ScDataType{
	301: {ID: 301, TypeCode: "ScTextInput", TypeName: "Text", DataType: "string", FrontDb: "", InputMode: "Input"},
	302: {ID: 302, TypeCode: "ScNumberInput", TypeName: "Number", DataType: "number", FrontDb: "", InputMode: "Input"},
	306: {ID: 306, TypeCode: "ScDateInput", TypeName: "Date", DataType: "string", FrontDb: "", InputMode: "Input"},
	307: {ID: 307, TypeCode: "ScDateTimeInput", TypeName: "DateTime", DataType: "string", FrontDb: "", InputMode: "Input"},
	401: {ID: 401, TypeCode: "ScSelectGender", TypeName: "Gender", DataType: "int16", FrontDb: "", InputMode: "Select"},
	404: {ID: 404, TypeCode: "ScSelectYesOrNo", TypeName: "Bool", DataType: "int16", FrontDb: "", InputMode: "Select"},
	510: {ID: 510, TypeCode: "ScPersonSelect", TypeName: "Person", DataType: "Person", FrontDb: "person", InputMode: "Select"},
	520: {ID: 520, TypeCode: "ScDeptSelect", TypeName: "Department", DataType: "SimpDept", FrontDb: "department", InputMode: "Select"},
	525: {ID: 525, TypeCode: "ScCSCSelect", TypeName: "Construction Site Category", DataType: "CSC", FrontDb: "csc", InputMode: "Select"},
	530: {ID: 530, TypeCode: "ScUDCSelect", TypeName: "User-define Category", DataType: "UDC", FrontDb: "udc", InputMode: "Select"},
	540: {ID: 540, TypeCode: "ScEPCSelect", TypeName: "Execution Project Category", DataType: "EPC", FrontDb: "epc", InputMode: "Select"},
	550: {ID: 550, TypeCode: "ScUDASelect", TypeName: "User-define Archive", DataType: "UDA", FrontDb: "uda", InputMode: "Select"},
}

// Default User Permissions List
var PublicFunctionList SystemMenus = SystemMenus{
	SystemMenu{ID: 1, FatherID: 0, Title: "MenuDashboard", Path: "/private/dashboard", Icon: "Home", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 10, FatherID: 0, Title: "MenuCalendar", Path: "/private/calendar", Icon: "CalendarMonth", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 15, FatherID: 0, Title: "MenuMessage", Path: "/private/message", Icon: "Message", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 20, FatherID: 0, Title: "MenuAddressBook", Path: "/private/addressBook", Icon: "ContactPhone", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9910, FatherID: 0, Title: "MenuProfile", Path: "/private/my/profile", Icon: "ManageAccounts", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
	SystemMenu{ID: 9920, FatherID: 0, Title: "MenuAbout", Path: "/private/my/about", Icon: "Info", Component: "", Selected: false, Indeterminate: false, AddFromVersion: "1.0.0"},
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
