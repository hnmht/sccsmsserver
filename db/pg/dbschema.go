package pg

import (
	"sccsmsserver/pkg/environment"
	"sccsmsserver/pkg/mysf"
	"sccsmsserver/pkg/security"
	"sccsmsserver/pub"

	"go.uber.org/zap"
)

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
		InitFunc:       initSysInfo,
	},
	{ //角色表
		TableName:   "sysrole",
		Description: "角色表",
		CreateSQL: `
			create table sysrole (
			id serial NOT NULL,
			rolename varchar(64) not null,
			description varchar(256),
			systemflag smallint DEFAULT 0,
			alluserflag smallint DEFAULT 0,
			create_time timestamp  with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp  with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default current_timestamp,
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initSysrole,
	},
	{ //用户表
		TableName:   "sysuser",
		Description: "用户表",
		CreateSQL: `
			create table sysuser (
			id serial NOT NULL,
			usercode varchar(32) NOT NULL,
			username varchar(64) NOT NULL,
			password varchar(64) NOT NULL,
			mobile varchar(32) default '',
			email varchar(64) default '',
			isoperator smallint DEFAULT 1,
			op_id int  DEFAULT 0,			
			file_id int DEFAULT 0,
			dept_id int DEFAULT 0,
			description varchar(256) default '',
			gender smallint DEFAULT 0,
			locked smallint DEFAULT 0,
			status smallint DEFAULT 0,			
			systemflag smallint DEFAULT 0,	
			create_time timestamp  with time zone default CURRENT_TIMESTAMP,
			createuserid int DEFAULT 0,
			modify_time timestamp  with time zone default CURRENT_TIMESTAMP,
			modifyuserid int DEFAULT 0,
			dr smallint  DEFAULT 0,
			ts timestamp with time zone default CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initSysUser,
	},
	{ //登录认证失败记录
		TableName:   "sysloginfault",
		Description: "登录认证失败记录",
		CreateSQL: `
			create table sysloginfault (
			id serial NOT NULL,
			user_id int DEFAULT 0,
			usercode varchar(32), 
			clientip varchar(32),
			useragent varchar(256),
			type smallint DEFAULT 0,
			ts timestamp with time zone default current_timestamp,
			PRIMARY KEY (id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //系统菜单表
		TableName:   "sysmenu",
		Description: "系统菜单表",
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
			ts timestamp with time zone default current_timestamp,
			PRIMARY KEY(autoid,id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initSysMenu,
	},
	{ //用户角色对照表
		TableName:   "sysuserrole",
		Description: "用户角色对照表",
		CreateSQL: `
			create table sysuserrole (
			id serial NOT NULL,
			user_id int,
			role_id int,
			create_time timestamp  with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp  with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			ts timestamp with time zone default current_timestamp,
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initSysUserRole,
	},
	{ //角色权限表
		TableName:   "sysrolemenu",
		Description: "角色权限表",
		CreateSQL: `
			create table sysrolemenu (
			id serial NOT NULL,
			role_id int,
			menu_id int,
			selected bool default true,
			indeterminate bool,
			create_time timestamp  with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp  with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			ts timestamp with time zone default current_timestamp,
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initSysRoleMenu,
	},
	{ //文件记录表
		TableName:   "filelist",
		Description: "文件记录表",
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
			create_time timestamp  with time zone default current_timestamp,
			createuserid int default 0,
			ts timestamp with time zone default current_timestamp,
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //部门档案表
		TableName:   "department",
		Description: "部门档案表",
		CreateSQL: `
			create table department (
			id serial NOT NULL,
			deptcode varchar(64), 
			deptname varchar(128),
			fatherid int default 0,
			leader int default 0,
			description varchar(256),
			status smallint DEFAULT 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initDepartment,
	},
	{ //岗位档案表
		TableName:   "operatingpost",
		Description: "岗位档案表",
		CreateSQL: `
			create table operatingpost(
			id serial NOT NULL,
			name varchar(128),
			description varchar(256),
			status smallint DEFAULT 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.2.0",
		InitFunc:       initOperatingPost,
	},
	{ //现场档案类别表
		TableName:   "sceneitemclass",
		Description: "现场档案类别表",
		CreateSQL: `
			create table sceneitemclass (
			id serial NOT NUll,
			classname varchar(128),
			description varchar(256),
			fatherid int default 0,
			status smallint default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initSceneItemClass,
	},
	{ //现场档案表
		TableName:   "sceneitem",
		Description: "现场档案表",
		CreateSQL: `
			create table sceneitem (
			id serial NOT NULL,
			code varchar(64),
			name varchar(64),
			description varchar(256),
			class_id int default 0,				
			subdept_id int default 0,
			respdept_id int default 0,
			respperson_id int default 0,
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
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //现场档案选项表
		TableName:   "sceneitemoption",
		Description: "现场档案选项表",
		CreateSQL: `
			create table sceneitemoption (
			id int,
			code varchar(64),
			name varchar(64),
			displayname varchar(64),
			udc_id int default 0,
			defaultvalue_id int default 0,
			enable smallint default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,	
			ts timestamp with time zone default current_timestamp,
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initSceneItemOption,
	},
	{ //自定义档案类别表
		TableName:   "userdefineclass",
		Description: "自定义档案类别表",
		CreateSQL: `
			create table userdefineclass (
			id serial NOT NULL,
			classname varchar(128),
			description varchar(256),
			islevel smallint default 0,
			status smallint DEFAULT 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},

	{ //自定义档案
		TableName:   "userdefinedoc",
		Description: "自定义档案",
		CreateSQL: `
			create table userdefinedoc (
			id serial NOT NULL,
			class_id int default 0,
			doccode varchar(128), 
			docname varchar(128),
			description varchar(256),
			fatherid int default 0,
			status smallint DEFAULT 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //执行项目类别表
		TableName:   "exectiveitemclass",
		Description: "执行项目类别表",
		CreateSQL: `
			create table exectiveitemclass (
			id serial NOT NUll,
			classname varchar(128),
			description varchar(256),
			fatherid int default 0,
			status smallint default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       initExectiveItemClass,
	},
	{ //执行项目表
		TableName:   "exectiveitem",
		Description: "执行项目表",
		CreateSQL: `
			create table exectiveitem (
			id serial NOT NUll,
			itemcode varchar(128),
			itemname varchar(128),
			class_id int default 0,
			description varchar(2048),
			status smallint default 0,
			resulttypeid int default 0,
			udc_id int default 0,
			defaultvalue varchar(1024),
			defaultvaluedisp varchar(1024),
			ischeckerror smallint default 0,
			errorvalue varchar(1024),
			errorvaluedisp varchar(1024),
			isrequirefile smallint default 0,
			isonsitephoto smallint default 0,
			risklevel_id int default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //执行模板档案表头表
		TableName:   "exectivetemplate_h",
		Description: "执行模板档案表头表",
		CreateSQL: `
			create table exectivetemplate_h (
			id serial NOT NUll,
			templatecode varchar(128),
			templatename varchar(128),
			description varchar(2048),
			status smallint default 0,
			allowaddrow smallint default 0,
			allowdelrow smallint default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //执行模板档案表体表
		TableName:   "exectivetemplate_b",
		Description: "执行模板档案表体表",
		CreateSQL: `
			create table exectivetemplate_b (
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
			createuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //单据流水号表
		TableName:   "serialno",
		Description: "单据流水号表",
		CreateSQL: `
			create table serialno (
			id serial NOT NULL,
			datestring varchar(8) NOT NULL,
			vouchertype varchar(4) NOT NULL,
			serialno int default 0,
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //指令单表头表
		TableName:   "workorder_h",
		Description: "指令单表头表",
		CreateSQL: `
			create table workorder_h (
			id serial NOT NUll,
			billnumber varchar(20),
			billdate varchar(10),
			dept_id int default 0,
			description varchar(256),
			status smallint default 0,
			workdate varchar(10),
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			confirm_time timestamp with time zone default current_timestamp,
			confirmuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
		);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //指令单表体表
		TableName:   "workorder_b",
		Description: "指令单表体表",
		CreateSQL: `
			create table workorder_b (
			id serial NOT NUll,
			hid int default 0,
			rownumber int default 0,
			si_id int default 0,
			ep_id int default 0,
			description varchar(256),
			eit_id int default 0,
			evnumber varchar(20) default '',
			starttime varchar(20),
			endtime varchar(20),
			status smallint default 0,
			ev_id int default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			confirm_time timestamp with time zone default current_timestamp,
			confirmuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //执行单表头表
		TableName:   "executedoc_h",
		Description: "执行单表头表",
		CreateSQL: `create table executedoc_h (
			id serial NOT NUll,
			billnumber varchar(20),
			billdate varchar(10),
			dept_id int default 0,
			description varchar(256),
			status smallint default 0,
			sourcetype varchar(8),
			sourcebillnumber varchar(20),
			source_hid int default 0,
			sourcerownumber int default 0,
			source_bid int default 0,
			starttime varchar(20),
			endtime varchar(20),
			si_id int default 0,
			ep_id int default 0,
			eit_id int default 0,
			allowaddrow smallint default 0,
			allowdelrow smallint default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			confirm_time timestamp with time zone default current_timestamp,
			confirmuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
		);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //执行单表体表
		TableName:   "executedoc_b",
		Description: "执行单表体表",
		CreateSQL: `create table executedoc_b (
			id serial NOT NUll,
			hid int default 0,
			rownumber int default 0,
			eid_id int default 0,
			allowdelrow smallint default 0,
			exectivevalue varchar(1024),
			exectivevaluedisp varchar(1024),				
			description varchar(256),
			eiddescription varchar(2048),
			ischeckerror smallint default 0,
			errorvalue varchar(1024),
			errorvaluedisp varchar(1024),
			isrequirefile smallint default 0,
			isonsitephoto smallint default 0,
			iserr smallint default 0,
			isrectify smallint default 0,
			ishandle smallint default 0,
			hp_id int default 0,
			handlestarttime varchar(20) default '',
			handleendtime varchar(20) default '',
			status smallint default 0,
			isfromeit smallint default 0,
			isfinish smallint default 0,
			dd_id int default 0,
			ddnumber varchar(20) default '',
			risklevel_id int default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			confirm_time timestamp with time zone default current_timestamp,
			confirmuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
		);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //执行单文件表
		TableName:   "executedoc_file",
		Description: "执行单附件表",
		CreateSQL: `create table executedoc_file(
			id serial NOT NUll,
			billbid int default 0,
			billhid int default 0,
			fileid int default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,				
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
		);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //执行单批注表
		TableName:   "executedoc_comment",
		Description: "执行单批注表",
		CreateSQL: `create table executedoc_comment(
			id serial NOT NUll,
			bid int default 0,
			hid int default 0,
			billnumber varchar(20),
			rownumber int default 0,
			sendto_id int default 0,
			isread smallint default 0,
			readtime varchar(20) default '',
			content varchar(512) default '',
			sendtime varchar(20) default '',
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,				
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,	
			PRIMARY KEY(id)
		);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //执行单审阅记录表
		TableName:   "executedoc_review",
		Description: "执行单审阅记录表",
		CreateSQL: `create table executedoc_review(
			id serial NOT NUll,
			hid int default 0,
			billnumber varchar(20),
			starttime varchar(17),
			endtime varchar(17),
			consumeseconds int default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,	
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,	
			PRIMARY KEY(id)
		);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //问题处理单表头表
		TableName:   "disposedoc",
		Description: "问题处理单表",
		CreateSQL: `create table disposedoc(
			id serial NOT NUll,
			billnumber varchar(20),
			billdate varchar(10),
			si_id int default 0,
			eid_id int default 0,
			exectivevalue varchar(1024),
			exectivevaluedisp varchar(1024),
			ep_id int default 0,
			dept_id int default 0,
			dp_id int default 0,
			isfinish smallint default 0,
			starttime varchar(20),
			endtime varchar(20),
			eddescription varchar(256),
			description varchar(256),				
			status smallint default 0,
			sourcetype varchar(8),
			sourcebillnumber varchar(20),
			source_hid int default 0,
			sourcerownumber int default 0,
			source_bid int default 0,
			risklevel_id int default 0,				
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			confirm_time timestamp with time zone default current_timestamp,
			confirmuserid int DEFAULT 0,				
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
		);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //问题处理单文件表
		TableName:   "disposedoc_file",
		Description: "问题处理单附件表",
		CreateSQL: `create table disposedoc_file(
			id serial NOT NUll,
			billbid int default 0,
			billhid int default 0,
			fileid int default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,				
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
		);`,
		AddFromVersion: "1.0.0",
		InitFunc:       simpleInitTable,
	},
	{ //首页内容定义表
		TableName:   "landingpage",
		Description: "首页内容定义表",
		CreateSQL: `create table if not exists landingpage(
			sysnamedisp varchar(64),
			introtext varchar(256),
			file_id int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			ts timestamp with time zone default current_timestamp
		);`,
		AddFromVersion: "1.1.0",
		InitFunc:       initLandingPage,
	},
	{ //风险等级表
		TableName:   "risklevel",
		Description: "风险等级表",
		CreateSQL: `create table if not exists risklevel(
			id serial NOT NUll,
			name varchar(128),
			description varchar(512),
			color varchar(128),
			status smallint default 0, 
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,				
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.1.0",
		InitFunc:       initRiskLevel,
	},
	{ //文档类别表
		TableName:   "documentclass",
		Description: "文档类别表",
		CreateSQL: `
			create table documentclass (
			id serial NOT NUll,
			classname varchar(128),
			description varchar(256),
			fatherid int default 0,
			status smallint default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.2.0",
		InitFunc:       initDocumentClass,
	},
	{ //文档表
		TableName:   "document",
		Description: "文档表",
		CreateSQL: `
			create table document (
			id serial NOT NUll,
			dc_id int default 0,
			name varchar(256) default '',
			edition varchar(256) default '',
			author varchar(256) default '',
			uploaddate varchar(20) default to_char(current_timestamp,'YYYYMMDD'),
			releasedate varchar(20) default to_char(current_timestamp,'YYYYMMDD'),
			tags varchar(1024) default '',
			description varchar(2048) default '',			
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.2.0",
		InitFunc:       simpleInitTable,
	},
	{ //文档文件表
		TableName:   "document_file",
		Description: "文档文件表",
		CreateSQL: `create table document_file(
			id serial NOT NUll,
			billbid int default 0,
			billhid int default 0,
			fileid int default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,				
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
		);`,
		AddFromVersion: "1.2.0",
		InitFunc:       simpleInitTable,
	},
	{ //培训课程表
		TableName:   "traincourse",
		Description: "培训课程表",
		CreateSQL: `
			create table traincourse (
			id serial NOT NUll,
			code varchar(256) default '',
			name varchar(256) default '',
			classhour numeric default 0,
			isexamine smallint default 1,
			description varchar(2048) default '',			
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.2.0",
		InitFunc:       simpleInitTable,
	},
	{ //培训课程表
		TableName:   "traincourse_file",
		Description: "培训课程附件表",
		CreateSQL: `
			create table traincourse_file(
			id serial NOT NUll,
			billbid int default 0,
			billhid int default 0,
			fileid int default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,				
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.2.0",
		InitFunc:       simpleInitTable,
	},
	{ //培训记录表头表
		TableName:   "trainrecord_h",
		Description: "培训记录表头表",
		CreateSQL: `create table trainrecord_h (
			id serial NOT NUll,
			billnumber varchar(20),
			billdate varchar(10),
			dept_id int default 0,
			description varchar(256),
			lecturer_id int default 0,
			traindate varchar(20),
			tc_id int default 0,
			starttime varchar(20),
			endtime varchar(20),
			classhour numeric default 0,
			isexamine smallint default 0,
			status smallint default 0,					
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			confirm_time timestamp with time zone default current_timestamp,
			confirmuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
		);`,
		AddFromVersion: "1.2.0",
		InitFunc:       simpleInitTable,
	},
	{ //培训记录表体表
		TableName:   "trainrecord_b",
		Description: "培训记录表体表",
		CreateSQL: `
			create table trainrecord_b (
			id serial NOT NUll,
			hid int default 0,
			rownumber int default 0,
			student_id int default 0,
			opname varchar(128) default '',
			deptname varchar(128) default '',
			starttime varchar(20),
			endtime varchar(20),
			classhour numeric default 0,
			description varchar(256),
			examineres smallint default 0,
			examinescore numeric default 0,
			status smallint default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			confirm_time timestamp with time zone default current_timestamp,
			confirmuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.2.0",
		InitFunc:       simpleInitTable,
	},
	{ //培训记录文件表
		TableName:   "trainrecord_file",
		Description: "培训记录附件表",
		CreateSQL: `create table trainrecord_file(
			id serial NOT NUll,
			billbid int default 0,
			billhid int default 0,
			fileid int default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,				
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
		);`,
		AddFromVersion: "1.2.0",
		InitFunc:       simpleInitTable,
	},
	{ //劳保用品档案表
		TableName:   "laborprotection",
		Description: "劳保用品档案表",
		CreateSQL: `
			create table laborprotection (
			id serial NOT NUll,
			code varchar(256) default '',
			name varchar(256) default '',
			model varchar(256) default '',
			unit varchar(256) default 'pcs',
			description varchar(2048) default '',			
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.2.0",
		InitFunc:       simpleInitTable,
	},
	{ //劳保用品定额主表
		TableName:   "lpaquota_h",
		Description: "劳保用品定额主表",
		CreateSQL: `
			create table lpaquota_h (
			id serial NOT NUll,
			billdate varchar(10),
			op_id int default 0,
			period varchar(20) default '',
			description varchar(2048) default '',
			status smallint default 0,			
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			confirm_time timestamp with time zone default current_timestamp,
			confirmuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.2.0",
		InitFunc:       simpleInitTable,
	},
	{ //劳保用品定额子表
		TableName:   "lpaquota_b",
		Description: "劳保用品定额子表",
		CreateSQL: `
			create table lpaquota_b (
			id serial NOT NUll,
			hid int default 0,
			rownumber int default 0,
			lp_id int default 0,
			quantity numeric default 0,
			description varchar(2048) default '',
			status smallint default 0,			
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			confirm_time timestamp with time zone default current_timestamp,
			confirmuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.2.0",
		InitFunc:       simpleInitTable,
	},
	{ //劳保用品发放单表头表
		TableName:   "lpaissuedoc_h",
		Description: "劳保用品发放单表头表",
		CreateSQL: `create table lpaissuedoc_h (
			id serial NOT NUll,
			billnumber varchar(20),
			billdate varchar(10),
			dept_id int default 0,
			description varchar(256) default '',
			period varchar(10) default '',			
			startdate varchar(20),
			enddate varchar(20),
			sourcetype varchar(8) default 'UA',
			status smallint default 0,					
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			confirm_time timestamp with time zone default current_timestamp,
			confirmuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
		);`,
		AddFromVersion: "1.2.0",
		InitFunc:       simpleInitTable,
	},
	{ //劳保用品发放单表体表
		TableName:   "lpaissuedoc_b",
		Description: "劳保用品发放单表体表",
		CreateSQL: `
			create table lpaissuedoc_b (
			id serial NOT NUll,
			hid int default 0,
			rownumber int default 0,
			recipient_id int default 0,
			opname varchar(128) default '',
			deptname varchar(128) default '',
			lp_id int default 0,
			quantity numeric default 0,
			description varchar(256),
			status smallint default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,
			confirm_time timestamp with time zone default current_timestamp,
			confirmuserid int DEFAULT 0,
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
			);`,
		AddFromVersion: "1.2.0",
		InitFunc:       simpleInitTable,
	},
	{ //劳保用品发放单文件表
		TableName:   "lpaissuedoc_file",
		Description: "劳保用品发放单文件表",
		CreateSQL: `create table lpaissuedoc_file(
			id serial NOT NUll,
			billbid int default 0,
			billhid int default 0,
			fileid int default 0,
			create_time timestamp with time zone default current_timestamp,
			createuserid int DEFAULT 0,				
			modify_time timestamp with time zone default current_timestamp,
			modifyuserid int DEFAULT 0,
			dr smallint default 0,			
			ts timestamp with time zone default current_timestamp,				
			PRIMARY KEY(id)
		);`,
		AddFromVersion: "1.2.0",
		InitFunc:       simpleInitTable,
	},
}

// 初始化sysinfo表
func initSysInfo() (isFinish bool, err error) {
	var rowNum int
	var sqlStr string
	isFinish = true
	//1.1 查询sysinfo表中的数据行数,sysinfo表中应当有且仅有一条数据
	sqlStr = "select count(isfinish) from sysinfo"
	err = db.QueryRow(sqlStr).Scan(&rowNum)
	if err != nil {
		isFinish = false
		zap.L().Error("查询sysinfo表记录数据时出现错误", zap.Error(err))
		return isFinish, err
	}

	//1.2 如果sysinfo表中数据数量超过1条则退出创建数据库表操作
	if rowNum > 1 {
		isFinish = false
		zap.L().Error("sysinfo表中数据数量超过1条", zap.Error(err))
		return isFinish, err
	}

	//1.3 如果sysinfo表中存在一条数据,则删除该条数据
	if rowNum == 1 {
		sqlStr = "delete from sysinfo"
		_, err = db.Exec(sqlStr)
		if err != nil {
			isFinish = false
			zap.L().Error("删除sysinfo系统旧数据时出现错误", zap.Error(err))
			return isFinish, err
		}
	}

	//1.4 向sysinfo表中插入数据
	//获取macArray
	macArray, err := environment.GetMacArray()
	if err != nil {
		isFinish = false
		zap.L().Error("获取macArray出现错误", zap.Error(err))
		return isFinish, err
	}

	//获取主板编号
	serialNumber, err := environment.GetSerialNumber()
	if err != nil {
		isFinish = false
		zap.L().Error("获取主板序列号时出现错误", zap.Error(err))
		return isFinish, err
	}

	//获取machineHash
	machineHash, err := environment.GetMachineHash(macArray, serialNumber)
	if err != nil {
		isFinish = false
		zap.L().Error("获取machineHash时出错", zap.Error(err))
		return isFinish, err
	}

	//生成RSA privatekey 和 publickey
	privateKey, publicKey, err := security.GenRsaKey(2048)
	if err != nil {
		isFinish = false
		zap.L().Error("生成rsa key时出错", zap.Error(err))
		return isFinish, err
	}

	//生成数据库唯一id
	dbID := mysf.GenID()

	sqlInsert := `insert into sysinfo(dbid,serialnumber,macarray,machinehash,privatekey,publickey,dbversion,starttime) values($1,$2,$3,$4,$5,$6,$7,now())`
	_, err = db.Exec(sqlInsert, dbID, serialNumber, macArray, machineHash, privateKey, publicKey, pub.DbVersion)
	if err != nil {
		isFinish = false
		zap.L().Error("向sysinfo中插入开始数据时错误", zap.Error(err))
		return isFinish, err
	}

	return
}

// 初始化sysrole表
func initSysrole() (isFinish bool, err error) {
	//检查role表中是否已经存在预置数据systemadmin
	sqlStr := "select count(id) as rownum from sysrole where id=10000"
	hasRecord, isFinish, err := checkRecord("sysrole sysadmin", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//插入sysAdmin角色
	sqlStr = "insert into sysrole(id,rolename,description,systemflag,alluserflag) values(10000,'systemadmin','系统预置角色',1,0)"
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("检查角色表中插入预置数据sysadmin出现错误", zap.Error(err))
		return isFinish, err
	}

	//2.3.2 检查role表中是否已经存在预置数据public
	sqlStr = "select count(id) as rownum from sysrole where id=10001"
	hasRecord, isFinish, err = checkRecord("sysrole public", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//插入public角色
	sqlStr = "insert into sysrole(id,rolename,description,systemflag,alluserflag) values(10001,'public','系统预置角色',1,1)"
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("检查角色表中插入预置数据public出现错误", zap.Error(err))
		return isFinish, err
	}
	return
}

// 初始化sysuser表
func initSysUser() (isFinish bool, err error) {
	//检查sysuser中是否存在预置数据admin
	sqlStr := "select count(id) as rownum from sysuser where id=10000"
	hasRecord, isFinish, err := checkRecord("sysuser", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//如果sysuser表中没有预置数据admin则向表中插入预置数据
	sqlStr = "insert into sysuser(id,username,password,create_time,description,systemflag,usercode,createuserid) values(10000,'系统管理员',$1,now(),'系统预置',1,'admin',10000)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("准备向sysuser表中插入预置admin数据时出错", zap.Error(err))
		return isFinish, err
	}
	defer stmt.Close()
	pwd := EncryptPassword(pub.DefaultPassword)
	_, err = stmt.Exec(pwd)
	if err != nil {
		isFinish = false
		zap.L().Error("向sysuser表中插入预置admin数据时出错", zap.Error(err))
		return isFinish, err
	}
	return
}

// 系统菜单表初始化
func initSysMenu() (isFinish bool, err error) {
	// 检查系统菜单表中是否存在数据
	sqlStr := "select count(id) as rownum from sysmenu"
	hasRecord, isFinish, err := checkRecord("sysmenu", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//如果不存在数据则插入预置数据
	sqlStr = "insert into sysmenu(id,fatherid,title,path,icon,component,selected,indeterminate) values($1,$2,$3,$4,$5,$6,$7,$8)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("准备向sysmenu中插入初始数据时出错", zap.Error(err))
		return isFinish, err
	}
	defer stmt.Close()
	for _, menu := range SysFunctionList {
		_, err = stmt.Exec(menu.ID, menu.FatherID, menu.Title, menu.Path, menu.Icon, menu.Component, menu.Selected, menu.Indeterminate)
		if err != nil {
			isFinish = false
			zap.L().Error("向sysmenu中插入初始数据'"+menu.Title+"'时出错", zap.Error(err))
			return isFinish, err
		}
	}
	return
}

// 用户角色对照表初始化
func initSysUserRole() (isFinish bool, err error) {
	//检查用户角色对照表中是否存在数据
	sqlStr := "select count(id) as rownum from sysuserrole where user_id=10000"
	hasRecord, isFinish, err := checkRecord("sysuserrole", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//如果不存在预置数据则插入预置数据
	//1 插入10000 用户 10000角色
	sqlStr = "insert into sysuserrole(user_id,role_id,ts) values(10000,10000,now())"
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("向sysuserrole插入系统管理员角色时出错", zap.Error(err))
		return isFinish, err
	}
	//2 插入10000 用户 10001角色
	sqlStr = "insert into sysuserrole(user_id,role_id,ts) values(10000,10001,now())"
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("向sysuserrole插入普通用户角色时出错", zap.Error(err))
		return isFinish, err
	}
	return
}

// 角色权限表初始化
func initSysRoleMenu() (isFinish bool, err error) {
	// 检查角色权限表中是否有记录
	sqlStr := "select count(id) as rownum from sysrolemenu where role_id=10000"
	hasRecord, isFinish, err := checkRecord("sysrolemenu", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//如果角色权限表中没有记录则插入预置数据
	// 1 插入systemadmin默认权限
	sqlStr1 := "insert into sysrolemenu(role_id,menu_id,selected,indeterminate,ts) values(10000,$1,true,false,now())"
	stmt1, err := db.Prepare(sqlStr1)
	if err != nil {
		isFinish = false
		zap.L().Error("准备向sysrolemenu中插入sysadmin角色初始数据时错误", zap.Error(err))
		return isFinish, err
	}
	defer stmt1.Close()
	for _, menu := range SysFunctionList {
		_, err = stmt1.Exec(menu.ID)
		if err != nil {
			isFinish = false
			zap.L().Error("向sysrolemenu中插入sysadmin角色初始数据"+menu.Title+"时出错", zap.Error(err))
			return isFinish, err
		}
	}
	//2 插入public角色默认权限
	sqlStr2 := "insert into sysrolemenu(role_id,menu_id,selected,indeterminate,ts) values(10001,$1,true,false,now())"
	stmt2, err := db.Prepare(sqlStr2)
	if err != nil {
		isFinish = false
		zap.L().Error("准备向sysrolemenu中插入public角色初始数据时错误", zap.Error(err))
		return isFinish, err
	}
	defer stmt2.Close()
	for _, menu := range PublicFunctionList {
		_, err = stmt2.Exec(menu.ID)
		if err != nil {
			isFinish = false
			zap.L().Error("向sysrolemenu中插入public角色初始数据"+menu.Title+"时出错", zap.Error(err))
			return isFinish, err
		}
	}
	return
}

// 部门档案表初始化
func initDepartment() (isFinish bool, err error) {
	//检查部门档案表中是否存在记录
	sqlStr := "select count(id) as rownum from department where id=10000"
	hasRecord, isFinish, err := checkRecord("department", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//如果表中没有数据则插入预置数据
	sqlStr = "insert into department(id,deptcode,deptname,description,createuserid) values(10000,'default','预置部门','系统预置部门',10000)"
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initDepartment insert initvalue failed", zap.Error(err))
		return isFinish, err
	}
	return
}

// 初始化岗位档案
func initOperatingPost() (isFinish bool, err error) {
	//检查岗位档案表中是否存在记录
	sqlStr := "select count(id) as rownum from operatingpost where id=10000"
	hasRecord, isFinish, err := checkRecord("operatingpost", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//如果表中没有数据则插入预置数据
	sqlStr = "insert into operatingpost(id,name,description,createuserid) values(10000,'预置岗位','系统预置岗位',10000)"
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initOperatingPost insert initvalue failed", zap.Error(err))
		return isFinish, err
	}
	return
}

// 现场档案分类表初始化
func initSceneItemClass() (isFinish bool, err error) {
	//检查执行项目分类表是否存在记录
	sqlStr := "select count(id) as rownum from sceneitemclass where id=1"
	hasRecord, isFinish, err := checkRecord("sceneitemclass", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//表中没有数据则插入预置数据
	sqlStr = `insert into sceneitemclass(id,classname,description,createuserid) values(10000,'默认分类','系统预置',10000)`
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initSceneItemClass insert initvalue failed", zap.Error(err))
		return isFinish, err
	}
	return
}

// 执行项目分类表初始化
func initExectiveItemClass() (isFinish bool, err error) {
	//检查执行项目分类表是否存在记录
	sqlStr := "select count(id) as rownum from exectiveitemclass where id=1"
	hasRecord, isFinish, err := checkRecord("exectiveitemclass", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//表中没有数据则插入预置数据
	sqlStr = `insert into exectiveitemclass(id,classname,description,createuserid) values(10000,'默认分类','系统预置',10000)`
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initExectiveItemClass insert initvalue failed", zap.Error(err))
		return isFinish, err
	}
	return
}

// 初始化首页内容定义表
func initLandingPage() (isFinish bool, err error) {
	//检查首页内容定义表是否存在记录
	sqlStr := "select count(file_id) as rownum from landingpage"
	hasRecord, isFinish, err := checkRecord("landingpage", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}

	//表中没有数据则插入预置数据
	sqlStr = `insert into landingpage(sysnamedisp,introtext,file_id,modifyuserid) 
		values('SeaCloud现场管理系统',
		'一套实用有效的企业安全生产信息化系统,包含现场管理、文档管理、培训管理、劳保用品管理四大模块，帮助企业有效落实安全生产措施.'
		,0,10000);`
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("landingpage insert initvalue failed", zap.Error(err))
		return isFinish, err
	}
	return
}

// 初始化风险等级表
func initRiskLevel() (isFinish bool, err error) {
	//检查首页内容定义表是否存在记录
	sqlStr := "select count(id) as rownum from risklevel"
	hasRecord, isFinish, err := checkRecord("risklevel", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//插入预置数据
	sqlStrs := []string{
		"insert into risklevel(id,name,description,color,createuserid) values(1,'重大风险','系统预置','red',10000)",
		"insert into risklevel(id,name,description,color,createuserid) values(2,'较大风险','系统预置','orange',10000)",
		"insert into risklevel(id,name,description,color,createuserid) values(3,'一般风险','系统预置','yellow',10000)",
		"insert into risklevel(id,name,description,color,createuserid) values(4,'低风险','系统预置','blue',10000)",
		"insert into risklevel(id,name,description,color,createuserid) values(5,'无风险','系统预置','white',10000)",
	}
	for _, t := range sqlStrs {
		_, err = db.Exec(t)
		if err != nil {
			isFinish = false
			zap.L().Error("initRiskLevel insert default data:"+t+" failed.", zap.Error(err))
			return
		}
	}
	return
}

// 文档分类表初始化
func initDocumentClass() (isFinish bool, err error) {
	//检查执行项目分类表是否存在记录
	sqlStr := "select count(id) as rownum from documentclass where id=1"
	hasRecord, isFinish, err := checkRecord("documentclass", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//表中没有数据则插入预置数据
	sqlStr = `insert into documentclass(id,classname,description,createuserid) values(10000,'默认分类','系统预置',10000)`
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initDocumentClass insert initvalue failed", zap.Error(err))
		return isFinish, err
	}
	return
}

// 通用初始化函数
func simpleInitTable() (isFinish bool, err error) {
	return true, nil
}

// 通用表内是否有记录
func checkRecord(tableName, sqlStr string) (hasRecord, isFinish bool, err error) {
	var rowNum int
	isFinish = true
	hasRecord = false
	err = db.QueryRow(sqlStr).Scan(&rowNum)
	if err != nil {
		isFinish = false
		zap.L().Error(tableName+" checkRecord failed", zap.Error(err))
		return
	}
	if rowNum > 0 { //存在数据
		hasRecord = true
		zap.L().Warn(tableName + " table has already exist data")
		return
	}
	return
}
