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
		TableName:   "sysrole",
		Description: "角色表",
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
