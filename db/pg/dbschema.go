package pg

import "go.uber.org/zap"

// Database table schema description struct
type table struct {
	TableName      string
	Description    string
	CreateSQL      string
	CreateIndexSQL string
	AddFromVersion string
	InitFunc       func() (bool, error)
}

// All database table information
var tables []table = []table{
	{
		TableName:   "sysinfo",
		Description: "system information",
		CreateSQL: `create table sysinfo (			
		dbid bigint,
		serialnumber varchar(64) default '',
		macarray varchar(1024) default '',
		machinehash varchar(512) default '',
		privatekey varchar(2048),
		publickey varchar(2048),
		starttime timestamp with time zone,
		endtime timestamp  with time zone,
		dbversion varchar(16),
		isFinish boolean DEFAULT false,
		registerflag smallint default 0,
		organizationid bigint default 0,
		organizationcode varchar(64) default '',
		organizationname varchar(2048) default '',
		contactperson varchar(32) default '',
		contacttitle varchar(32) default '',
		phone varchar(32) default '',
		email varchar(32) default '',
		registertime varchar(20) default ''
		);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initSysInfo,
	},
	{
		TableName:   "sysmenu",
		Description: "System menus table",
		CreateSQL: `
			create table sysmenu (
			autoid serial NOT NULL,
			id int NOT NULL,
			fatherid int,
			title varchar(64),
			path varchar(256),
			icon varchar(64),
			component varchar(128),
			selected bool default false,
			indeterminate bool default false,
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default current_timestamp,
			PRIMARY KEY(autoid,id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initSysMenu,
	},
	{
		TableName:   "sysrole",
		Description: "system Role",
		CreateSQL: `
			create table sysrole (
			id serial NOT NULL,
			rolename varchar(64) not null,
			description varchar(256),
			systemflag smallint DEFAULT 0,
			alluserflag smallint DEFAULT 0,
			createtime timestamp with time zone default current_timestamp,
			creatorid int DEFAULT 0,
			confirmtime timestamp with time zone default current_timestamp,
			confirmerid int DEFAULT 0,
			modifytime timestamp with time zone default current_timestamp,
			modifierid int DEFAULT 0,
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default current_timestamp,
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initSysrole,
	},
	{
		TableName:   "sysuser",
		Description: "User",
		CreateSQL: `
			create table sysuser (
			id serial NOT NULL,
			code varchar(32) NOT NULL,
			name varchar(64) NOT NULL,
			password varchar(64) NOT NULL,
			mobile varchar(32) default '',
			email varchar(64) default '',
			isoperator smallint DEFAULT 1,
			positionid int  DEFAULT 0,			
			fileid int DEFAULT 0,
			deptid int DEFAULT 0,
			description varchar(256) default '',
			gender smallint DEFAULT 0,
			locked smallint DEFAULT 0,
			status smallint DEFAULT 0,			
			systemflag smallint DEFAULT 0,	
			createtime timestamp  with time zone default CURRENT_TIMESTAMP,
			creatorid int DEFAULT 0,
			modifytime timestamp  with time zone default CURRENT_TIMESTAMP,
			modifierid int DEFAULT 0,
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initSysUser,
	},
	{
		TableName:   "sysuserrole",
		Description: "User and Role Mapping",
		CreateSQL: `
			create table sysuserrole (
			id serial NOT NULL,
			userid int,
			roleid int,
			createtime timestamp  with time zone default current_timestamp,
			creatorid int DEFAULT 0,
			modifytime timestamp  with time zone default current_timestamp,
			modifierid int DEFAULT 0,
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default current_timestamp,
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initSysUserRole,
	},
	{
		TableName:   "sysrolemenu",
		Description: "Role and Menu Mapping Table",
		CreateSQL: `
			create table sysrolemenu (
			id serial NOT NULL,
			roleid int,
			menuid int,
			selected bool default true,
			indeterminate bool,
			createtime timestamp  with time zone default current_timestamp,
			creatorid int DEFAULT 0,
			modifytime timestamp  with time zone default current_timestamp,
			modifierid int DEFAULT 0,
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default current_timestamp,
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initSysRoleMenu,
	},
	{
		TableName:   "sysloginfault",
		Description: "User authentication failed record",
		CreateSQL: `
			create table sysloginfault (
			id serial NOT NULL,
			userid int DEFAULT 0,
			usercode varchar(32), 
			clientip varchar(32),
			useragent varchar(256),
			type smallint DEFAULT 0,
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default current_timestamp,
			PRIMARY KEY (id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       genericInitTable,
	},
	{
		TableName:   "filelist",
		Description: "File information Record",
		CreateSQL: `
			create table filelist (
			id serial NOT NULL,
			filehash varchar(64),
			miniofilename varchar(256),
			originfilename varchar(256),
			filekey int default 0,
			filetype varchar(64),
			isimage smallint default 0,
			model varchar(128),
			longitude numeric,
			latitude numeric,
			size int,
			fileurl varchar(256),
			datetimeoriginal varchar(12) default '',
			uploaddate timestamp with time zone,
			source varchar(20) default 'browser',
			createtime timestamp  with time zone default current_timestamp,
			creatorid int default 0,
			ts timestamp with time zone default current_timestamp,
			dr smallint  DEFAULT 0,
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       genericInitTable,
	},
	{
		TableName:   "department",
		Description: "Department Master data",
		CreateSQL: `
			create table department (
			id serial NOT NULL,
			code varchar(64), 
			name varchar(128),
			fatherid int default 0,
			leader int default 0,
			description varchar(256),
			status smallint DEFAULT 0,
			createtime timestamp with time zone default current_timestamp,
			creatorid int DEFAULT 0,
			modifytime timestamp with time zone default current_timestamp,
			modifierid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initDepartment,
	},
	{
		TableName:   "position",
		Description: "Position Master Data",
		CreateSQL: `
			create table position(
			id serial NOT NULL,
			name varchar(128),
			description varchar(256),
			status smallint DEFAULT 0,
			createtime timestamp with time zone default current_timestamp,
			creatorid int DEFAULT 0,
			modifytime timestamp with time zone default current_timestamp,
			modifierid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.2.0",
		InitFunc:       initPosition,
	},
	{
		TableName:   "csc",
		Description: "Construction Site Category",
		CreateSQL: `
			create table csc (
			id serial NOT NUll,
			name varchar(128),
			description varchar(256),
			fatherid int default 0,
			status smallint default 0,
			createtime timestamp with time zone default current_timestamp,
			creatorid int DEFAULT 0,
			modifytime timestamp with time zone default current_timestamp,
			modifierid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initCSC,
	},
	{
		TableName:   "cs",
		Description: "Construction Site Master Data",
		CreateSQL: `
			create table cs (
			id serial NOT NULL,
			code varchar(64),
			name varchar(64),
			description varchar(256),
			cscid int default 0,				
			subdeptid int default 0,
			respdeptid int default 0,
			resppersonid int default 0,
			status smallint DEFAULT 0,
			finishflag smallint DEFAULT 0,
			finishdate varchar(16),
			longitude numeric,
			latitude numeric,
			udf1 int default 0,
			udf2 int default 0,
			udf3 int default 0,
			udf4 int default 0,
			udf5 int default 0,
			udf6 int default 0,
			udf7 int default 0,
			udf8 int default 0,
			udf9 int default 0,
			udf10 int default 0,
			createtime timestamp with time zone default current_timestamp,
			creatorid int DEFAULT 0,
			modifytime timestamp with time zone default current_timestamp,
			modifierid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initCS,
	},
	{
		TableName:   "cso",
		Description: "Construction Site options",
		CreateSQL: `
			create table cso (
			id int,
			code varchar(64),
			name varchar(64),
			displayname varchar(64),
			udcid int default 0,
			defaultvalueid int default 0,
			enable smallint default 0,
			createtime timestamp with time zone default current_timestamp,
			creatorid int DEFAULT 0,
			modifytime timestamp with time zone default current_timestamp,
			modifierid int DEFAULT 0,
			dr smallint default 0,	
			ts timestamp with time zone default current_timestamp,
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initCSO,
	},
	{
		TableName:   "udc",
		Description: "User-define Category",
		CreateSQL: `
			create table udc (
			id serial NOT NULL,
			name varchar(128),
			description varchar(256),
			islevel smallint default 0,
			status smallint DEFAULT 0,
			createtime timestamp with time zone default current_timestamp,
			creatorid int DEFAULT 0,
			modifytime timestamp with time zone default current_timestamp,
			modifierid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       genericInitTable,
	},
	{
		TableName:   "ud",
		Description: "User-define Master Data",
		CreateSQL: `
			create table ud (
			id serial NOT NULL,
			udcid int default 0,
			code varchar(128), 
			name varchar(128),
			description varchar(256),
			fatherid int default 0,
			status smallint DEFAULT 0,
			createtime timestamp with time zone default current_timestamp,
			creatorid int DEFAULT 0,
			modifytime timestamp with time zone default current_timestamp,
			modifierid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       genericInitTable,
	},
	{
		TableName:   "epc",
		Description: "Execution Project Category Master Data",
		CreateSQL: `
			create table exectiveitemclass (
			id serial NOT NUll,
			classname varchar(128),
			description varchar(256),
			fatherid int default 0,
			status smallint default 0,
			createtime timestamp with time zone default current_timestamp,
			creatorid int DEFAULT 0,
			modifytime timestamp with time zone default current_timestamp,
			modifierid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initEPC,
	},
	{
		TableName:   "ep",
		Description: "Execution Project",
		CreateSQL: `
			create table ep (
			id serial NOT NUll,
			code varchar(128),
			name varchar(128),
			epcid int default 0,
			description varchar(2048),
			status smallint default 0,
			resulttypeid int default 0,
			udcid int default 0,
			defaultvalue varchar(1024),
			defaultvaluedisp varchar(1024),
			ischeckerror smallint default 0,
			errorvalue varchar(1024),
			errorvaluedisp varchar(1024),
			isrequirefile smallint default 0,
			isonsitephoto smallint default 0,
			risklevelid int default 0,
			createtime timestamp with time zone default current_timestamp,
			creatorid int DEFAULT 0,
			modifytime timestamp with time zone default current_timestamp,
			modifierid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       genericInitTable,
	},
	{
		TableName:   "epth",
		Description: "Execution Project Template Header",
		CreateSQL: `
			create table ept_h (
			id serial NOT NUll,
			code varchar(128),
			name varchar(128),
			description varchar(2048),
			status smallint default 0,
			allowaddrow smallint default 0,
			allowdelrow smallint default 0,
			create_time timestamp with time zone default current_timestamp,
			creatorid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifierid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       genericInitTable,
	},
	{ //执行模板档案表体表
		TableName:   "ept_b",
		Description: "Execution Project Template Body",
		CreateSQL: `
			create table ept_b (
			id serial NOT NUll,
			hid int default 0,
			rownumber int default 0,
			eid_id int default 0,
			allowdelrow smallint default 0,
			description varchar(2048),
			defaultvalue varchar(1024),
			defaultvaluedisp varchar(1024),
			ischeckerror smallint default 0,
			errorvalue varchar(1024),
			errorvaluedisp varchar(1024),
			isrequirefile smallint default 0,
			isonsitephoto smallint default 0,
			risklevel_id int default 0,
			create_time timestamp with time zone default current_timestamp,
			creatorid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifierid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       genericInitTable,
	},
}

// Generic database table initialization function.
// Tables that don't require initializaiton use this funciton.
func genericInitTable() (isFinish bool, err error) {
	return true, nil
}

// Generic database table record check function
func genericCheckRecord(tableName, sqlStr string) (hasRecord, isFinish bool, err error) {
	var rowNum int
	isFinish = true
	hasRecord = false
	err = db.QueryRow(sqlStr).Scan(&rowNum)
	if err != nil {
		isFinish = false
		zap.L().Error("genericCheckRecord checking "+tableName+" failed", zap.Error(err))
		return
	}
	if rowNum > 0 {
		hasRecord = true
		zap.L().Warn(tableName + " table has already exist data")
		return
	}
	return
}
