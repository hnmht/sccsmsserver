package pg

// 定义数据库表数据结构
type Table struct {
	TableName      string
	Description    string
	CreateSQL      string
	CreateIndexSQL string
	AddFromVersion string
	InitFunc       func() (bool, error)
}

// 数据库表信息
var Tables []Table = []Table{
	{ //系统信息表
		TableName:   "sysinfo",
		Description: "系统信息表",
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
		jobrole varchar(32) default '',
		phone varchar(32) default '',
		email varchar(32) default '',
		registertime varchar(20) default ''
		);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
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
		InitFunc:       simpleInitTable,
	},
	{
		TableName:   "sysmsg_t",
		Description: "System Message Translations",
		CreateSQL: `create table sys_msg (
			id serial NOT NULL,
			code int default 0,
			language varchar(10) default '',
			content varchar(2048) default '',
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{
		TableName:   "logicmsg",
		Description: "Business Logic Message",
		CreateSQL: `create table sysmsg (
			id serial NOT NULL,
			code int default 0,
			content varchar(2048) default '',
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{
		TableName:   "logicmsg_t",
		Description: "Business Logic Message",
		CreateSQL: `create table sysmsg (
			id serial NOT NULL,
			code int default 0,
			language varchar(10) default '',
			content varchar(2048) default '',
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
}

// 通用初始化函数
func simpleInitTable() (isFinish bool, err error) {
	return true, nil
}
