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
		TableName:   "i18n",
		Description: "Internationalization",
		CreateSQL: `create table i18n (
			id serial NOT NULL,
			language varchar(16) default '',
			name varchar(128) default '',
			weekfirstday varchar(10) default '',
			shortdateformat varchar(20) default '',
			longdateformat varchar(20) default '',
			shorttimeformat varchar(20) default '',
			longtimeformat varchar(20) default '',
			timezone varchar(20) default 'UTC',
			createtime timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			confirmtime timestamp with time zone default current_timestamp,
			confirmuserid int DEFAULT 0,
			modifytime timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initI18n,
	},
	{
		TableName:   "sysmsg",
		Description: "System Message",
		CreateSQL: `create table sysmsg (
			id serial NOT NULL,
			code int default 0,
			content varchar(2048) default '',
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initSysMsg,
	},
	{
		TableName:   "sysmsg_t",
		Description: "System Message Translations",
		CreateSQL: `create table sysmsg_t (
			id serial NOT NULL,
			code int default 0,
			defaultcontent varchar(2048) default '',
			language varchar(10) default '',
			content varchar(2048) default '',
			createtime timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			confirmtime timestamp with time zone default current_timestamp,
			confirmuserid int DEFAULT 0,
			modifytime timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initSysMsgTranslate,
	},
	{
		TableName:   "logicmsg",
		Description: "Business Logic Message",
		CreateSQL: `create table logicmsg (
			id serial NOT NULL,
			code int default 0,
			content varchar(2048) default '',
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       genericInitTable,
	},
	{
		TableName:   "logicmsg_t",
		Description: "Business Logic Message",
		CreateSQL: `create table logicmsg_t (
			id serial NOT NULL,
			code int default 0,
			language varchar(10) default '',
			content varchar(2048) default '',
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       genericInitTable,
	},
}

// Generic database table initialization function. Tables that don't require initializaiton use this funciton.
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
