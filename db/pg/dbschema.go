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
			usercode varchar(32) NOT NULL,
			username varchar(64) NOT NULL,
			password varchar(64) NOT NULL,
			mobile varchar(32) default '',
			email varchar(64) default '',
			isoperator smallint DEFAULT 1,
			position_id int  DEFAULT 0,			
			fileid int DEFAULT 0,
			deptid int DEFAULT 0,
			description varchar(256) default '',
			gender smallint DEFAULT 0,
			locked smallint DEFAULT 0,
			status smallint DEFAULT 0,			
			systemflag smallint DEFAULT 0,	
			createtime timestamp  with time zone default CURRENT_TIMESTAMP,
			createuserid int DEFAULT 0,
			modifytime timestamp  with time zone default CURRENT_TIMESTAMP,
			modifyuserid int DEFAULT 0,
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
			createuserid int DEFAULT 0,
			modifytime timestamp  with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
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
			createuserid int DEFAULT 0,
			modifytime timestamp  with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default current_timestamp,
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initSysRoleMenu,
	},
	{ //登录认证失败记录
		TableName:   "sysloginfault",
		Description: "登录认证失败记录",
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
