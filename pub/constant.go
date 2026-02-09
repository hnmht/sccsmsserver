package pub

import "time"

// Master Data Type
type DataType string

// Sea&Cloud Construction Site management System Software info
type ServerSoftInfo struct {
	ScSoftName      string `json:"scSoftName"`
	ScServerVersion string `json:"scServerVersion"`
	ScServerState   string `json:"scServerState"`
	MaxUserNumber   int32  `json:"maxUserNumber"`
	Author          string `json:"author"`
}

// User Organization Information struct
type OrganizationInfo struct {
	RegisterFlag     int16  `db:"registerflag" json:"registerFlag"`
	OrganizationID   int64  `db:"organizationid" json:"organizationID,string"`
	OrganizationCode string `db:"organizationcode" json:"organizationCode"`
	OrganizationName string `db:"organizationname" json:"organizationName"`
	ContactPerson    string `db:"contactperson" json:"contactPerson"`
	ContactTitle     string `db:"contacttitle" json:"contactTitle"`
	Phone            string `db:"phone" json:"phone"`
	Email            string `db:"email" json:"email"`
	RegisterTime     string `db:"registertime" json:"registerTime"`
}

// Default Organization  Information
var DefaultOrg OrganizationInfo = OrganizationInfo{
	RegisterFlag:     0,
	OrganizationID:   1,
	OrganizationCode: "Unknown Code",
	OrganizationName: "Unknown Company",
	ContactPerson:    "Unknow",
	ContactTitle:     "Unknown",
	Phone:            "Unknown Phone",
	Email:            "unknown@email.com",
	RegisterTime:     "20060102150405",
}

// Software Information
var SoftInfo ServerSoftInfo = ServerSoftInfo{
	ScSoftName:      "Sea&Cloud Construction Site management System Backend",
	ScServerVersion: "1.0.0",
	ScServerState:   "Standard",
	MaxUserNumber:   100000,
	Author:          "Haitao Meng",
}

// API service path
const APIPath = "/api/v1"

// New user default password
const DefaultPassword string = "sc@123"

// Database Schema version
const DbVersion = "1.0.0"

// Md5 secret
var Md5Secret = []byte("Sea&Cloud comes from a character in both my wife's and my names.")

// JWT secret
var TokenSecret = []byte("This is JWT secret")

// Master Data Type
const (
	User       DataType = "user"       // User Profile
	Person     DataType = "person"     // Simple User Profile
	Department DataType = "department" // Department Profile
	SimpDept   DataType = "simpdept"   // Simple Department Profile
	File       DataType = "file"       // File metadata
	FileHash   DataType = "filehash"   // File metadata by hash
	CSC        DataType = "csc"        // Construction site Category
	SimpCSC    DataType = "simpcsc"    // Simple Construction site Category
	CSA        DataType = "csa"        // Construction Site Archive
	UDC        DataType = "udc"        // User-defined Data Category
	UDA        DataType = "uda"        // User-defined Archive Master Data
	EPC        DataType = "epc"        // Execution Project Category
	SimpEPC    DataType = "simpepc"    // Simple Execution Project Category
	EPA        DataType = "epa"        // Execution Project Archive
	EPT        DataType = "ept"        // Execution Project Template
	EPTHead    DataType = "epthead"    // Execution Project Template Header
	EPTBody    DataType = "eptbody"    // Execution Project Template Body
	CFCS       DataType = "cfcs"       // Custom Fields for Construction Site
	ACFCS      DataType = "acfcs"      // All custom fileds for Construction Site
	RL         DataType = "rl"         // Risk Level Master Data
	DC         DataType = "dc"         // Document Category
	SimpDC     DataType = "simpdc"     // Simple Document Category
	Document   DataType = "document"   // Document Archive
	Position   DataType = "position"   // Position Master Data
	TC         DataType = "tc"         // Training Course Master Data
	PPE        DataType = "ppe"        // Personal Protective Equipment
	IPBlack    DataType = "ipblack"    // IP Address Blacklist
)

// Valid values for the "clientType" request header
var ValidClientTypes = [2]string{"sceneweb", "scenemob"}

// Minio File URL expiration time
const FileURLExpireTime = 24 * time.Hour

// Cache expiration time
const CacheExpiration = 2 * time.Hour

// Token Expire Duration
const TokenExpireDuration = 2 * time.Hour

// Token Issuer
const TokenIssuer = "SeaCloud"

// Event Background Colors
var EventBackgroundColors = [4]string{"orange", "blue", "green", "grey"}

// Token About to Expire seconds
const TokenAboutToExpirtSeconds = float64(60)
