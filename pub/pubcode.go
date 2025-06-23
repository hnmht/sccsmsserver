package pub

// Basic Archive Type
type DocType string

const (
	User       DocType = "user"       // User Archive
	Person     DocType = "person"     // Simple User Archive
	Department DocType = "department" // Department Archive
	SimpDept   DocType = "simpdept"   // Simple Department Profile
	File       DocType = "file"       // File Archive
	OSAC       DocType = "osac"       // On-Site Archive Category
	SimpOSAC   DocType = "simposac"   // Simple On-Site Archive Category
	OSA        DocType = "osa"        // On-Site Archive
	UDAC       DocType = "udc"        // User-defined Archive Category
	UDA        DocType = "uda"        // User-defined Archive
	EPC        DocType = "epc"        // Execution Project Category
	SimpEPC    DocType = "simpepc"    // Simple Execution Project Category
	EPA        DocType = "epa"        // Execution Project Archive
	EPT        DocType = "ept"        // Execution Project Template
	EPTHead    DocType = "epthead"    // Execution Project Template Header
	EPTBody    DocType = "eptbody"    // Execution Project Template Body
	CFOSA      DocType = "cfosa"      // Custom Fields for On-site Archive
	ACFOSA     DocType = "acfosa"     // All custom fileds for On-site Archive
	RLA        DocType = "rla"        // Risk Level Archive
	DC         DocType = "dc"         // Document Category
	SimpDC     DocType = "simpdc"     // Simple Document Category
	Document   DocType = "document"   // Document Archive
	JPA        DocType = "jpa"        // Job Profile Archive
	TCA        DocType = "tca"        // Training Course Archive
	LPEA       DocType = "lp"         // Labor Protection Equipment Archive
	IPBlack    DocType = "ipblack"    // IP Address Blacklist
)

// Server Resopnse Message Type
type ResCode int64

const (
	CodeSuccess ResCode = 1000

	CodeInvalidParm   ResCode = 1001
	CodeServerBusy    ResCode = 1098
	CodeInternalError ResCode = 1099

	CodeInvalidToken  ResCode = 1100
	CodeNeedLogin     ResCode = 1101
	CodeClientUnknown ResCode = 1102
	CodeClientEmpty   ResCode = 1103
	CodeTokenDestroy  ResCode = 1104
	CodeLoginOther    ResCode = 1105
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess: "Success",

	CodeInvalidParm:   "Request Parameter Error",
	CodeServerBusy:    "Server is Busy",
	CodeInternalError: "Internal Server Error",

	CodeInvalidToken:  "Login credentials have expired",
	CodeNeedLogin:     "Login required",
	CodeClientUnknown: "Unknown Client",
	CodeClientEmpty:   "Client type is empty",
	CodeTokenDestroy:  "Login credentials have been revoked by the administrator",
	CodeLoginOther:    "Login credentials are invalid, the user has logged in on another client",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]

	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}

	return msg
}

// Message type returned by bussiness logic
type ResStatus int64

const (
	// success
	StatusOK ResStatus = 0
	// Concurrency issue（10001-10099）
	StatusOtherEdit   ResStatus = 10001
	StatusDataDeleted ResStatus = 10002
	// Authorization (10100-10199)
	StatusUserNotExist      ResStatus = 10100
	StatusInvalidPassword   ResStatus = 10101
	StatusPasswordDisaccord ResStatus = 10102
	StatusUserDisabled      ResStatus = 10103
	StatusUserLocked        ResStatus = 10104
	StatusOverAuthoriaztion ResStatus = 10105
	// Role Management (10200-10299)
	StatusRoleNameExist ResStatus = 10200
	StatusRoleUserExist ResStatus = 10201
	StatusRoleAuthExist ResStatus = 10202
	// User management (10300-10399)
	StatusUserIDExist     ResStatus = 10300
	StatusUserNameExist   ResStatus = 10301
	StatusUserMobileExist ResStatus = 10302
	StatusUserEmailExist  ResStatus = 10303
	StatusProfileOnlySelf ResStatus = 10304
	// Department (10400-10499)
	StatusDeptCodeExist     ResStatus = 10400
	StatusDeptFatherSelf    ResStatus = 10401
	StatusDeptFatherCircle  ResStatus = 10402
	StatusDeptLowLevelExist ResStatus = 10403
	StatusDeptNotExist      ResStatus = 10405
	// On-Site Archive Category (10500-10599)
	StatusOSACNameExist     ResStatus = 10500
	StatusOSACFatherSelf    ResStatus = 10501
	StatusOSACFatherCircle  ResStatus = 10502
	StatusOSACLowLevelExist ResStatus = 10503
	// On-Site Archive (10600-10699)
	StatusOSACodeExist ResStatus = 10600
	// User-Defined Archive Category (10700-10799)
	StatusUDACNameExist ResStatus = 10700
	//  User-Defined Archive (10800-10899)
	StatusUDACodeExist ResStatus = 10800
	// Execution Project Category (10900-10999)
	StatusEPCNameExist     ResStatus = 10900
	StatusEPCFatherSelf    ResStatus = 10901
	StatusEPCFatherCircle  ResStatus = 10902
	StatusEPCLowLevelExist ResStatus = 10903
	// Execution Project Archive (11000-11099)
	StatusEPACodeExist        ResStatus = 11000
	StatusEPAChangeResultType ResStatus = 11001
	// Execution Project Template (11100-11199)
	StatusEPTCodeExist ResStatus = 11100
	// File Archive (11200-11299)
	StatusFileOpenFailed   ResStatus = 11200
	StatusFileUploadFailed ResStatus = 11201
	StatusFileGetUrlFailed ResStatus = 11202
	StatusFileNotExist     ResStatus = 11203
	// Work Order (11300-11399)
	StatusWOOtherEdit ResStatus = 11300
	// Execution Order (11400-11499)
	StatusEOBodyNoConfirm       ResStatus = 11400
	StatusIssueDisposeCompleted ResStatus = 11401
	// Message (11500-11599)
	StatusMsgOnlyReadSelf ResStatus = 11500
	// Risk Level Archive（11600-11699)
	StatusRLANameExist ResStatus = 11600
	// Document Archive Category (11700-11799)
	StatusDACNameExist     ResStatus = 11700
	StatusDACLowLevelExist ResStatus = 11701
	StatusDACFatherSelf    ResStatus = 11702
	// Document Archive (11800-11899)
	StatusDocumentNoFile ResStatus = 11800
	// Job Profile Archive(11900-11999)
	StatusJPANameExist ResStatus = 11900
	// Training Course Archive (12000-12099)
	StatusTCANameExist ResStatus = 12000
	// Labor Protection Equipment Archive (12100-12199)
	StatusLPEACodeExist ResStatus = 12100
	// Labor Protection Supplies Post Quota (12200-12199)
	StatusLPSPQExist ResStatus = 12200
	// Referenced （80000-89999）
	StatusUDAUsed            ResStatus = 80000
	StatusEPAUsed            ResStatus = 80001
	StatusEPADefaultUsed     ResStatus = 80002
	StatusEPAErrorUsed       ResStatus = 80003
	StatusEPTDefaultUsed     ResStatus = 80004
	StatusEPTErrorUsed       ResStatus = 80005
	StatusEPTUsed            ResStatus = 80006
	StatusUserUsed           ResStatus = 80007
	StatusDeptLeaderUsed     ResStatus = 80008
	StatusOSAUsed            ResStatus = 80009
	StatusWOUsed             ResStatus = 80010
	StatusEOUsed             ResStatus = 80011
	StatusEOValueUsed        ResStatus = 80012
	StatusEOErrorUsed        ResStatus = 80013
	StatusIHFUsed            ResStatus = 80014 // Issue Handling Form
	StatusOSARespUsed        ResStatus = 80015
	StatusWOEpUsed           ResStatus = 80016
	StatusEOEpUsed           ResStatus = 80017
	StatusEODpUsed           ResStatus = 80018
	StatusEOCommentUsed      ResStatus = 80019
	StatusEOSendToUsed       ResStatus = 80020
	StatusEOReviewUsed       ResStatus = 80021
	StatusIHFEpUsed          ResStatus = 80022
	StatusIHFDpUsed          ResStatus = 80023
	StatusOSAOUsed           ResStatus = 80024 // On-Site Archive Options
	StatusDocumentUsed       ResStatus = 80025
	StatusTRUsed             ResStatus = 80026 // Training Record
	StatusWOCreateUsed       ResStatus = 80027
	StatusWOModifyUsed       ResStatus = 80028
	StatusWOConfirmUsed      ResStatus = 80029
	StatusEOCreateUsed       ResStatus = 80030
	StatusEOModifyUsed       ResStatus = 80031
	StatusEOConfirmUsed      ResStatus = 80032
	StatusIHFCreateUsed      ResStatus = 80033
	StatusIHFModifyUsed      ResStatus = 80034
	StatusIHFConfirmUsed     ResStatus = 80035
	StatusDCCreateUsed       ResStatus = 80036
	StatusDCModifyUsed       ResStatus = 80037
	StatusDocumentCreateUsed ResStatus = 80038
	StatusDocumentModifyUsed ResStatus = 80039
	StatusTCACreateUsed      ResStatus = 80040
	StatusTCAModifyUsed      ResStatus = 80041
	StatusTRLecturerUsed     ResStatus = 80042 // Training Record
	StatusTRStudentUsed      ResStatus = 80043
	StatusTRCreateUsed       ResStatus = 80044
	StatusTRModifyUsed       ResStatus = 80045
	StatusTRConfirmUsed      ResStatus = 80046
	StatusLPSPQCreateUsed    ResStatus = 80047 // Labor Protection Supplies Post Quota (12200-12199)
	StatusLPSPQModifyUsed    ResStatus = 80048
	StatusLPSPQConfirmUsed   ResStatus = 80049
	StatusPPEIFCreateUsed    ResStatus = 80050 // PPE Issuance Form
	StatusPPEIFModifyUsed    ResStatus = 80051
	StatusPPEIFConfirmUsed   ResStatus = 80052
	StatusPPEIFRecipientUsed ResStatus = 80053
	StatusDeptCreateUsed     ResStatus = 80054
	StatusDeptModifyUsed     ResStatus = 80055
	StatusJPACreateUsed      ResStatus = 80056
	StatusJPAModifyUsed      ResStatus = 80057
	StatusOSACCreateUsed     ResStatus = 80058
	StatusOSACModifyUsed     ResStatus = 80059
	StatusOSACreateUsed      ResStatus = 80060
	StatusOSAModifyUsed      ResStatus = 80061
	StatusUDACCreateUsed     ResStatus = 80062
	StatusUDACModifyUsed     ResStatus = 80063
	StatusUDACreateUsed      ResStatus = 80064
	StatusUDAModifyUsed      ResStatus = 80065
	StatusEPCCreateUsed      ResStatus = 80066
	StatusEPCModifyUsed      ResStatus = 80067
	StatusEPACreateUsed      ResStatus = 80068
	StatusEPAModifyUsed      ResStatus = 80069
	StatusRLACreateUsed      ResStatus = 80070
	StatusRLAModifyUsed      ResStatus = 80071
	StatusLPEACreateUsed     ResStatus = 80072
	StatusLPEAModifyUsed     ResStatus = 80073
	StatusEPTCreateUsed      ResStatus = 80074
	StatusEPTModifyUsed      ResStatus = 80075
	StatusUserCreateUsed     ResStatus = 80076
	StatusUserModifyUsed     ResStatus = 80077
	StatusRoleCreateUsed     ResStatus = 80078
	StatusRoleModifyUsed     ResStatus = 80079
	StatusTRDeptUsed         ResStatus = 80080
	StatusPPEIFDeptUsed      ResStatus = 80081
	// Other (90000-99999)
	StatusErrorUnknow              ResStatus = 90900
	StatusResCodeError             ResStatus = 90901
	StatusInternalError            ResStatus = 90902
	StatusResReject                ResStatus = 90903
	StatusResNoData                ResStatus = 90904
	StatusOverRecord               ResStatus = 90905
	StatusVoucherNoBody            ResStatus = 90906
	StatusVoucherNoFree            ResStatus = 90907
	StatusVoucherOnlyCreateEdit    ResStatus = 90908
	StatusVoucherNoConfirm         ResStatus = 90909
	StatusVoucherCancelConfirmSelf ResStatus = 90910
)

var resStatusCodeMsg = map[ResStatus]string{
	StatusOK: "success",
	// Concurrency issue（10001-10099）
	StatusOtherEdit:   "The data is being modified by another user",
	StatusDataDeleted: "The data has been deleted",
	// Authorization (10100-10199)
	StatusUserNotExist:      "User does not exist",
	StatusInvalidPassword:   "Incorrect password",
	StatusPasswordDisaccord: "The two entered passwords do not match",
	StatusUserDisabled:      "User account has been disabled",
	StatusUserLocked:        "User account has been locked",
	StatusOverAuthoriaztion: "The number of users exceeds the maximum authorized limit",
	// Role Management(10200-10299)
	StatusRoleNameExist: "Role name has been exist",
	StatusRoleUserExist: "The relationship between the role and the user already exists",
	StatusRoleAuthExist: "This role has already been assigned permissions",
	// User Management (10300-10399)
	StatusUserIDExist:     "User code already exists",
	StatusUserNameExist:   "User name already exists",
	StatusUserMobileExist: "User phone number already exists",
	StatusUserEmailExist:  "User email already exists",
	StatusProfileOnlySelf: "User can only modify their own information",
	// Department Archive (10400-10499)
	StatusDeptCodeExist:     "Department Code exists",
	StatusDeptFatherSelf:    "The parent department cannot be the department itself",
	StatusDeptFatherCircle:  "Loop detected in parent department hierarchy",
	StatusDeptLowLevelExist: "Sub-department exist",
	StatusDeptNotExist:      "Department does not exist",
	// On-Site Archive Category (10500-10599)
	StatusOSACNameExist:     "On-Site archive category name exists",
	StatusOSACFatherSelf:    "The parent category cannot be itself",
	StatusOSACFatherCircle:  "Loop detected in parent category hierarchy",
	StatusOSACLowLevelExist: "Sub-category exist",
	// On-site Archive (10600-10699)
	StatusOSACodeExist: "On-site archive code already exist",
	//User-defined Archive Category (10700-10799)
	StatusUDACNameExist: "User-defined archive category name already exist",
	// User-defined Archive (10800-10899)
	StatusUDACodeExist: "User-defined archive code already exist",
	// Execution Project Category(10900-10999)
	StatusEPCNameExist:     "Execution project category name already exists",
	StatusEPCFatherSelf:    "The parent category cannot be itself",
	StatusEPCFatherCircle:  "Loop detected in parent category hierarchy",
	StatusEPCLowLevelExist: "Sub-category exists",
	// Execution Project Archive(11000-11099)
	StatusEPACodeExist:        "Execution project code already exists",
	StatusEPAChangeResultType: "Referenced items cannot have their result type modified",
	// Execution Project Template (11100-11199)
	StatusEPTCodeExist: "Execution project template code already exists",
	// File Archive (11200-11299)
	StatusFileOpenFailed:   "Failed to open file",
	StatusFileUploadFailed: "File upload failed",
	StatusFileGetUrlFailed: "Failed to get file URL",
	StatusFileNotExist:     "File does not exist",
	// Work Order (11300-11399)
	StatusWOOtherEdit: "Work Order has been modified by another user",
	// Execution Order (11400-11499)
	StatusEOBodyNoConfirm:       "There are unconfirmed row in the execution order body",
	StatusIssueDisposeCompleted: "The issue has been resolved",
	// Message (11500-11599)
	StatusMsgOnlyReadSelf: "Messages can only be read by the recipient themselves",
	// Risk Level Archive (11600-11699)
	StatusRLANameExist: "Resk level name already exists",
	// Document Archive Category (11700-11799)
	StatusDACNameExist:     "Document Archive Category name already exists",
	StatusDACLowLevelExist: "Sub-category exists",
	StatusDACFatherSelf:    "The parent category cannot be itself",
	// Document Archive (11800-11899)
	StatusDocumentNoFile: "The Document Archive has no attachments",
	// Job Profile Archive (11900-11999)(11900-11999)
	StatusJPANameExist: "Job Profile name already exists",
	// Training Course Archive(12000-12099)
	StatusTCANameExist: "Training Course name already exists",
	// Labor Protection Equipment Archive (12100-12199)
	StatusLPEACodeExist: "Labor Protection Equipment Archive code already exists",
	// Labor Protection Supplies Post Quota (12200-12199)
	StatusLPSPQExist: "A labor protection supplies post quota for the same period already exists",
	// Referenced （80000-89999）
	StatusUDAUsed:            "Referenced by User Defined Archive",
	StatusEPAUsed:            "Referenced by Eexcution Project Archive",
	StatusEPADefaultUsed:     "Referenced by Execution Project default value",
	StatusEPAErrorUsed:       "Referenced by Execution Project error value",
	StatusEPTDefaultUsed:     "Referenced by Execution Project Template default values",
	StatusEPTErrorUsed:       "Referenced by Execution Project Template error values",
	StatusEPTUsed:            "Referenced by Execution Project Template",
	StatusUserUsed:           "Referenced by User Archive",
	StatusDeptLeaderUsed:     "Referenced by Department Archive",
	StatusOSAUsed:            "Referenced by On-site Archive",
	StatusWOUsed:             "Referenced by Work Order",
	StatusEOUsed:             "Referenced by Execution Order",
	StatusEOValueUsed:        `Referenced by Execution Order value`,
	StatusEOErrorUsed:        "Referenced by Execution Order error value",
	StatusIHFUsed:            "Referenced by Issue Handling Form",
	StatusOSARespUsed:        "Referenced by On-site Archive owner",
	StatusWOEpUsed:           "Referenced by Execution Order executor",
	StatusEOEpUsed:           "Referenced by Execution order head executor",
	StatusEODpUsed:           "Referenced by Execution Order body handler",
	StatusEOCommentUsed:      "Referenced by Execution Order annotator",
	StatusEOSendToUsed:       "Referenced by Execution Order annotation message recipient",
	StatusEOReviewUsed:       "Referenced by Execution Order reviewer",
	StatusIHFEpUsed:          "Referenced by Issue Handling Form executor",
	StatusIHFDpUsed:          "Referenced by Issue Handling Form handler",
	StatusOSAOUsed:           "Referenced by On-site Archive options",
	StatusDocumentUsed:       "Referenced by Document Archive",
	StatusTRUsed:             "Referenced by Training Record",
	StatusWOCreateUsed:       "Referenced by Work Order creator",
	StatusWOModifyUsed:       "Referenced by Work Order editor",
	StatusWOConfirmUsed:      "Referenced by Work Order confirmer",
	StatusEOCreateUsed:       "Referenced by Execution Order creator",
	StatusEOModifyUsed:       "Referenced by Execution Order editor",
	StatusEOConfirmUsed:      "Referenced by Execution Order confirmer",
	StatusIHFCreateUsed:      "Referenced by Issue Handling Form creator",
	StatusIHFModifyUsed:      "Referenced by Issue Handling Form editor",
	StatusIHFConfirmUsed:     "Referenced by Issue Handling Form confirmer",
	StatusDCCreateUsed:       "Referenced by Document Archive Category creator",
	StatusDCModifyUsed:       "Referenced by Document Archive Category editor",
	StatusDocumentCreateUsed: "Referenced by Document Archive creator",
	StatusDocumentModifyUsed: "Referenced by Document Archive editor",
	StatusTCACreateUsed:      "Referenced by Training Course creator",
	StatusTCAModifyUsed:      "Referenced by Training Course editor",
	StatusTRLecturerUsed:     "Referenced by Training Record speaker",
	StatusTRStudentUsed:      "Referenced by Training Record trainee",
	StatusTRCreateUsed:       "Referenced by Training Record creator",
	StatusTRModifyUsed:       "Referenced by Training Record editor",
	StatusTRConfirmUsed:      "Referenced by Training Record confirmer",
	StatusLPSPQCreateUsed:    "Referenced by Labor Protection Supplies Post Quota creator",
	StatusLPSPQModifyUsed:    "Referenced by Labor Protection Supplies Post Quota editor",
	StatusLPSPQConfirmUsed:   "Referenced by Labor Protection Supplies Post Quota confirmer",
	StatusPPEIFCreateUsed:    "Referenced by PPE Issuance Form creator",
	StatusPPEIFModifyUsed:    "Referenced by PPE Issuance Form editor",
	StatusPPEIFConfirmUsed:   "Referenced by PPE Issuance Form confirmer",
	StatusPPEIFRecipientUsed: "Referenced by PPE Issuance Form recipient",
	StatusDeptCreateUsed:     "Referenced by Department Arctive creator",
	StatusDeptModifyUsed:     "Referenced by Department Arctive editor",
	StatusJPACreateUsed:      "Referenced by Job Profile Archive creator",
	StatusJPAModifyUsed:      "Referenced by Job Profile Archive editor",
	StatusOSACCreateUsed:     "Referenced by On-site Archive Category creator",
	StatusOSACModifyUsed:     "Referenced by On-site Archive Category editor",
	StatusOSACreateUsed:      "Referenced by On-site Archive creator",
	StatusOSAModifyUsed:      "Referenced by On-site Archive editor",
	StatusUDACCreateUsed:     "Referenced by User-defined Archive Category creator",
	StatusUDACModifyUsed:     "Referenced by User-defined Archive Category editor",
	StatusUDACreateUsed:      "Referenced by User-defined Archive creator",
	StatusUDAModifyUsed:      "Referenced by User-defined Archive editor",
	StatusEPCCreateUsed:      "Referenced by Execution Project Category creator",
	StatusEPCModifyUsed:      "Referenced by Execution Project Category editor",
	StatusEPACreateUsed:      "Referenced by Execution Project Archive creator",
	StatusEPAModifyUsed:      "Referenced by Execution Project Archive editor",
	StatusRLACreateUsed:      "Referenced by Risk Level Archive creator",
	StatusRLAModifyUsed:      "Referenced by Risk Level editor",
	StatusLPEACreateUsed:     "Referenced by Labor Protection Equipment Archive creator",
	StatusLPEAModifyUsed:     "Referenced by Labor Protection Equipment Archive editor",
	StatusEPTCreateUsed:      "Referenced by Execution Project Template creator",
	StatusEPTModifyUsed:      "Referenced by Execution Project Template editor",
	StatusUserCreateUsed:     "Referenced by User Archive creator",
	StatusUserModifyUsed:     "Referenced by User Archive editor",
	StatusRoleCreateUsed:     "Referenced by Role Archive creator",
	StatusRoleModifyUsed:     "Referenced by Role Archive editor",
	StatusTRDeptUsed:         "Referenced by Training Records head department",
	StatusPPEIFDeptUsed:      "Referenced by the issuing department in the PPE Issuance Form",
	// Other(90000-99999)
	StatusErrorUnknow:              "Unknown error",
	StatusResCodeError:             "Response error",
	StatusInternalError:            "Internal server error",
	StatusResReject:                "Server denied",
	StatusResNoData:                "No data",
	StatusOverRecord:               "The number of records exceeds the maximum limit",
	StatusVoucherNoBody:            "Document body is empty",
	StatusVoucherNoFree:            "Document status is not free",
	StatusVoucherOnlyCreateEdit:    "Only the document creator can modify the document",
	StatusVoucherNoConfirm:         "Document status is not confirmed",
	StatusVoucherCancelConfirmSelf: "Only the confirmer of the document can cancel the document confirmation",
}

func (r ResStatus) Msg() string {
	msg, ok := resStatusCodeMsg[r]
	if !ok {
		msg = resStatusCodeMsg[StatusErrorUnknow]
	}
	return msg
}
