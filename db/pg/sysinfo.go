package pg

import (
	"sccsmsserver/pkg/environment"
	"sccsmsserver/pkg/mysf"
	"sccsmsserver/pkg/security"
	"sccsmsserver/pub"

	"go.uber.org/zap"
)

// Server information struct
type ServerInfo struct {
	DbID         int64                `db:"dbid" json:"dbID,string"`          //数据库id
	SerialNumber string               `db:"serialnumber" json:"serialNumber"` //服务器主板序列号
	MacArray     string               `db:"macarray" json:"macArray"`         //服务器网卡mac地址
	MachineHash  string               `db:"machinehash" json:"machineHash"`   //服务器硬件Hash
	MachineID    string               `db:"machineid" json:"machineID"`       //服务器ID
	PublicKey    string               `db:"publickey" json:"publicKey"`       //服务器公钥
	DbVersion    string               `db:"dbversion" json:"dbVersion"`       //数据结构版本
	Organization pub.OrganizationInfo `json:"organization"`                   //组织信息
	ServerSoft   pub.ServerSoftInfo   `json:"serverSoft"`                     //服务器软件信息
}

// Server public information
var ServerPubInfo ServerInfo

// Initialize the sysinfo table
func initSysInfo() (isFinish bool, err error) {
	var rowNum int
	var sqlStr string
	isFinish = true
	// Step 1: Query the row count from sysinfo table. There should be exactly one row in this table.
	sqlStr = "select count(isfinish) from sysinfo"
	err = db.QueryRow(sqlStr).Scan(&rowNum)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysInfo db.QueryRow failed:", zap.Error(err))
		return
	}

	// Step 2: If the sysinfo table contains more than one row, Exit the database table creation operation
	if rowNum > 1 {
		isFinish = false
		zap.L().Error("initSysInfo query the sysinfo table records failed:  contains more than one row.")
		return
	}

	// Step 3: If there is a record in the sysinfo table, delete it.
	if rowNum == 1 {
		sqlStr = "delete from sysinfo"
		_, err = db.Exec(sqlStr)
		if err != nil {
			isFinish = false
			zap.L().Error("initSysInfo db.Exec delete old data failed:", zap.Error(err))
			return
		}
	}

	// Step 4: Insert data into the sysinfo table
	// Step 4.1： Get the machines's network adapter MAC address
	macArray, err := environment.GetMacArray()
	if err != nil {
		isFinish = false
		zap.L().Error("initSysInfo environment.GetMacArray failed:", zap.Error(err))
		return
	}

	// Step 4.2: Get the machine's motherboard serial number
	serialNumber, err := environment.GetSerialNumber()
	if err != nil {
		isFinish = false
		zap.L().Error("initSysInfo environment.GetSerialNumber failed:", zap.Error(err))
		return
	}

	// Step 4.3: Get the machine's hash value
	machineHash, err := environment.GetMachineHash(macArray, serialNumber)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysInfo environment.GetMachineHash failed:", zap.Error(err))
		return
	}

	// Step 4.4: Generate RSA public and private keys
	privateKey, publicKey, err := security.GenRsaKey(2048)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysInfo security.GenRsaKey failed:", zap.Error(err))
		return
	}

	// Step 4.5: Generate a unique database ID.
	dbID := mysf.GenID()

	// Step 4.6: Insert data into the sysinfo table.
	sqlInsert := `insert into sysinfo(dbid,serialnumber,macarray,machinehash,privatekey,
		publickey,starttime,dbversion,registerflag,organizationid,organizationcode,
		organizationname,contactperson,contacttitle,phone,email,
		registertime) values($1,$2,$3,$4,$5,
		$6,now(),$7,$8,$9,$10,
		$11,$12,$13,$14,$15,
		$16)`
	_, err = db.Exec(sqlInsert, dbID, serialNumber, macArray, machineHash, privateKey,
		publicKey, pub.DbVersion, pub.DefaultOrg.RegisterFlag, pub.DefaultOrg.OrganizationID, pub.DefaultOrg.OrganizationCode,
		pub.DefaultOrg.OrganizationName, pub.DefaultOrg.ContactPerson, pub.DefaultOrg.ContactTitle, pub.DefaultOrg.Phone, pub.DefaultOrg.Email,
		pub.DefaultOrg.RegisterTime)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysInfo db.Exec(sqlInsert) insert data failed:", zap.Error(err))
		return
	}
	return
}

// Initialize Current Server information
func (info *ServerInfo) Init() (err error) {
	/*// Get the machineID
	 info.MachineID, err = machineid.ProtectedID(pub.Secret)
	if err != nil {
		zap.L().Error("ServerInfo.Init machineid.ProtectedID() failed")
		return
	}

	// Get the machines's network adapter MAC address
	info.MacArray, err = environment.GetMacArray()
	if err != nil {
		zap.L().Error("SysInfo.Init GetMacArray failed:", zap.Error(err))
		return
	}
	// Get the machine's motherboard serial number
	info.SerialNumber, err = environment.GetSerialNumber()
	if err != nil {
		zap.L().Error("SysInfo.Init GetSerialNumber failed:", zap.Error(err))
		return
	}

	// Get the machine's hash value
	info.MachineHash, err = environment.GetMachineHash(info.MacArray, info.SerialNumber)
	if err != nil {
		zap.L().Error("SysInfo.Init GetMachineHash failed", zap.Error(err))
		return
	}
	*/
	// Retrieve server public information from the sysinfo table.
	sqlStr := `select dbid,dbversion,publickey,organizationid,organizationcode,
	organizationname,contactperson,contacttitle,phone,email,
	registertime 
	from sysinfo limit(1)`
	err = db.QueryRow(sqlStr).Scan(&info.DbID, &info.DbVersion, &info.PublicKey, &info.Organization.OrganizationID, &info.Organization.OrganizationCode,
		&info.Organization.OrganizationName, &info.Organization.ContactPerson, &info.Organization.ContactTitle, &info.Organization.Phone, &info.Organization.Email,
		&info.Organization.RegisterTime)
	if err != nil {
		zap.L().Error("ServerInfo.Init db.QueryRow failed:", zap.Error(err))
		return
	}

	info.ServerSoft = pub.SoftInfo
	zap.L().Info("Server infomaition initilized successfully.")
	return
}
