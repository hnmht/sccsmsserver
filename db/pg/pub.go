package pg

import (
	"sccsmsserver/pub"

	"github.com/denisbrodbeck/machineid"
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
	DbID         int64              `db:"dbid" json:"dbid,string"`          //数据库id
	SerialNumber string             `db:"serialnumber" json:"serialnumber"` //服务器主板序列号
	MacArray     string             `db:"macarray" json:"macarray"`         //服务器网卡mac地址
	MachineHash  string             `db:"machinehash" json:"machinehash"`   //服务器硬件Hash
	MachineID    string             `db:"machineid" json:"machineid"`       //服务器ID
	PublicKey    string             `db:"publickey" json:"publickey"`       //服务器公钥
	DbVersion    string             `db:"dbversion" json:"dbversion"`       //数据结构版本
	ServerSoft   pub.ServerSoftInfo `json:"serversoft"`                     //服务器软件信息
}

var ServerPubInfo ServerInfo

// Initialize Current Server information
func (info *ServerInfo) Init() (err error) {
	// Get MachineID
	info.MachineID, err = machineid.ProtectedID(pub.Secret)
	if err != nil {
		zap.L().Error("ServerInfo.Init machineid.ProtectedID() failed")
		return
	}

	return
}
