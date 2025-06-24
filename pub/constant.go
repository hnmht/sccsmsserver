package pub

import "time"

// Master Data Type
type DataType string

// Sea&Cloud Construction Site management System Software info
type ServerSoftInfo struct {
	ScSoftName      string `json:"scsoftname"`
	ScServerVersion string `json:"scserverversion"`
	ScServerState   string `json:"scserverstate"`
	MaxUserNumber   int32  `json:"maxusernumber"`
	Author          string `json:"author"`
}

// Software Information
var SoftInfo ServerSoftInfo = ServerSoftInfo{
	ScSoftName:      "Sea&Cloud Construction Site management System Backend",
	ScServerVersion: "1.0.0",
	ScServerState:   "Standard",
	MaxUserNumber:   100000,
	Author:          "Haitao Meng",
}

// Cache expiration time
const CacheExpiration = 2 * time.Hour

// New user default password
const DefaultPassword = "sc@123"

// Database Schema version
const DbVersion = "1.0.0"

// Md5 secret
const Secret = "Sea&Cloud comes from a character in both my wife's and my names."

// Default Language
const DefaultLanguage = "en_us"

// Master Data Type
const (
	User       DataType = "user"       // User Profile
	Person     DataType = "person"     // Simple User Profile
	Department DataType = "department" // Department Profile
	SimpDept   DataType = "simpdept"   // Simple Department Profile
	File       DataType = "file"       // File metadata
	CSC        DataType = "csc"        // Construction site Category
	SimpOSAC   DataType = "simposac"   // Simple Construction site Category
	CS         DataType = "cs"         // Construction Site Master Data
	UDDT       DataType = "uddt"       // User-defined Data Type
	UDD        DataType = "udd"        // User-defined Data
	EPC        DataType = "epc"        // Execution Project Category
	SimpEPC    DataType = "simpepc"    // Simple Execution Project Category
	EP         DataType = "ep"         // Execution Project Master Data
	EPT        DataType = "ept"        // Execution Project Template
	EPTHead    DataType = "epthead"    // Execution Project Template Header
	EPTBody    DataType = "eptbody"    // Execution Project Template Body
	CFCS       DataType = "cfcs"       // Custom Fields for Construction Site
	ACFCS      DataType = "acfcs"      // All custom fileds for Construction Site
	RL         DataType = "rl"         // Risk Level Master Data
	DC         DataType = "dc"         // Document Category
	SimpDC     DataType = "simpdc"     // Simple Document Category
	Document   DataType = "document"   // Document Archive
	Pos        DataType = "pos"        // Position Master Data
	TC         DataType = "tc"         // Training Course Master Data
	PPE        DataType = "ppe"        // Personal Protective Equipment
	IPBlack    DataType = "ipblack"    // IP Address Blacklist
)
