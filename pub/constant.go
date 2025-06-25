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

// User Organization Information struct
type OrganizationInfo struct {
	RegisterFlag     int16  `db:"registerflag" json:"registerflag"`
	OrganizationID   int64  `db:"organizationid" json:"organizationid,string"`
	OrganizationCode string `db:"organizationcode" json:"organizationcode"`
	OrganizationName string `db:"organizationname" json:"organizationname"`
	ContactPerson    string `db:"contactperson" json:"contactperson"`
	ContactTitle     string `db:"contacttitle" json:"contacttitle"`
	Phone            string `db:"phone" json:"phone"`
	Email            string `db:"email" json:"email"`
	RegisterTime     string `db:"registertime" json:"registertime"`
}

// Localization Struct
type Localization struct {
	Language        string `json:"language"`
	WeekFirstDay    string `json:"weekfirstday"`
	ShortDateFormat string `json:"shortdateformat"`
	LongDateFormat  string `json:"longdateformat"`
	ShortTimeFormat string `json:"shorttimeformat"`
	LongTimeFormat  string `json:"longtimeformat"`
	TimeZone        string `json:"timezone"`
	Description     string `json:"description"`
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

// Cache expiration time
const CacheExpiration = 2 * time.Hour

// New user default password
const DefaultPassword = "sc@123"

// Database Schema version
const DbVersion = "1.0.0"

// Md5 secret
const Secret = "Sea&Cloud comes from a character in both my wife's and my names."

// Default locale
var DefaultLocale = Localization{
	Language:        "en_us",
	WeekFirstDay:    "Sunday",
	ShortDateFormat: "MM/DD/YY",
	LongDateFormat:  "MM/DD/YYYY",
	ShortTimeFormat: "HH:MM AM/PM",
	LongTimeFormat:  "HH:MM:SS AM/PM",
	TimeZone:        "UTC-5",
	Description:     "System default locale",
}

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
