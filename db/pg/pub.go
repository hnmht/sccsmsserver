package pg

import (
	"sccsmsserver/pub"

	"go.uber.org/zap"
)

// Data references check struct
type DataReferenceCheck struct {
	Description    string
	SqlStr         string
	UsedReturnCode pub.ResStatus
}

// Data query conditions
type QueryParams struct {
	QueryString string `json:"querystring"`
}

// Data query pagination
type PagingQueryParams struct {
	QueryString string `json:"querystring"`
	Page        int32  `json:"page"`
	PerPage     int32  `json:"perpage"`
}

// Server information struct
type ServerInfo struct {
	DbID         int64                `db:"dbid" json:"dbid,string"`          //数据库id
	SerialNumber string               `db:"serialnumber" json:"serialnumber"` //服务器主板序列号
	MacArray     string               `db:"macarray" json:"macarray"`         //服务器网卡mac地址
	MachineHash  string               `db:"machinehash" json:"machinehash"`   //服务器硬件Hash
	MachineID    string               `db:"machineid" json:"machineid"`       //服务器ID
	PublicKey    string               `db:"publickey" json:"publickey"`       //服务器公钥
	DbVersion    string               `db:"dbversion" json:"dbversion"`       //数据结构版本
	Organization pub.OrganizationInfo `json:"organization"`                   //组织信息
	ServerSoft   pub.ServerSoftInfo   `json:"serversoft"`                     //服务器软件信息
}

// Server public information
var ServerPubInfo ServerInfo

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
