package pub

//DocType 基本档案类型
type DocType string

const (
	User       DocType = "user"       // User Archive
	Person     DocType = "person"     // Simple User Archive
	Department DocType = "department" // Department Archive
	SimpDept   DocType = "simpdept"   // Simple Department Profile
	File       DocType = "file"       // File Archive
	SIC        DocType = "sic"        // On-Site Archive Category
	SimpSIC    DocType = "simpsic"    // Simple On-Site Archive Category
	SI         DocType = "si"         // On-Site Archive
	UDC        DocType = "udc"        // User-defined Archive Category
	UDD        DocType = "udd"        // User-defined Archive
	EIC        DocType = "eic"        // Execution Project Category
	SimpEIC    DocType = "simpeic"    // Simple Execution Project Category
	EID        DocType = "eid"        // Execution Project
	EIT        DocType = "eit"        // Execution Project Template
	EITHead    DocType = "eithead"    // Execution Project Template Header
	EITBody    DocType = "eitbody"    // Execution Project Template Body
	SIO        DocType = "sio"        // Custom Fields for On-site Archive
	SIOS       DocType = "sios"       // All custom fileds for On-site Archive
	RL         DocType = "rl"         // Risk Level
	DC         DocType = "dc"         // Document Category
	SimpDC     DocType = "simpdc"     // Simple Document Category
	Document   DocType = "document"   // Document Archive
	OP         DocType = "op"         // Job Profile
	TC         DocType = "tc"         // Training Course
	LP         DocType = "lp"         // Labor Protection Equipment Archive
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
	// On-Site archive category (10500-10599)
	StatusSICNameExist     ResStatus = 10500
	StatusSICFatherSelf    ResStatus = 10501
	StatusSICFatherCircle  ResStatus = 10502
	StatusSICLowLevelExist ResStatus = 10503
	// On-Site archive (10600-10699)
	StatusSICodeExist ResStatus = 10600
	// User-defined archive category (10700-10799)
	StatusUDCNameExist ResStatus = 10700
	//  User-defined archive (10800-10899)
	StatusUDDCodeExist ResStatus = 10800
	// Execution project category (10900-10999)
	StatusEICNameExist     ResStatus = 10900
	StatusEICFatherSelf    ResStatus = 10901
	StatusEICFatherCircle  ResStatus = 10902
	StatusEICLowLevelExist ResStatus = 10903
	// Execution project (11000-11099)
	StatusEIDCodeExist        ResStatus = 11000
	StatusEIDChangeResultType ResStatus = 11001
	// Execution project template (11100-11199)
	StatusEITCodeExist ResStatus = 11100
	// File Archive (11200-11299)
	StatusFileOpenFailed   ResStatus = 11200
	StatusFileUploadFailed ResStatus = 11201
	StatusFileGetUrlFailed ResStatus = 11202
	StatusFileNotExist     ResStatus = 11203
	// Work Order (11300-11399)
	StatusWOOtherEdit ResStatus = 11300
	// Execution Order (11400-11499)
	StatusEDBodyNoConfirm       ResStatus = 11400
	StatusErrorDisposeCompleted ResStatus = 11401
	// Message (11500-11599)
	StatusMsgOnlyReadSelf ResStatus = 11500
	// Risk Level（11600-11699)
	StatusRLNameExist ResStatus = 11600
	// Document Archive Category (11700-11799)
	StatusDCNameExist     ResStatus = 11700
	StatusDCLowLevelExist ResStatus = 11701
	StatusDCFatherSelf    ResStatus = 11702
	// Document Archive (11800-11899)
	StatusDocumentNoFile ResStatus = 11800
	// Job Profile (11900-11999)
	StatusOPNameExist ResStatus = 11900
	// Training Course (12000-12099)
	StatusTCNameExist ResStatus = 12000
	// Labor Protection Equipment Archive (12100-12199)
	StatusLPCodeExist ResStatus = 12100
	// Labor Protection Supplies Post Quota (12200-12199)
	StatusOPQuotaExist ResStatus = 12200
	// Referenced （80000-89999）
	StatusUDDUsed            ResStatus = 80000
	StatusEIDUsed            ResStatus = 80001
	StatusEIDDefaultUsed     ResStatus = 80002
	StatusEIDErrorUsed       ResStatus = 80003
	StatusEITDefaultUsed     ResStatus = 80004
	StatusEITErrorUsed       ResStatus = 80005
	StatusEITUsed            ResStatus = 80006
	StatusUserUsed           ResStatus = 80007
	StatusDeptLeaderUsed     ResStatus = 80008
	StatusSIUsed             ResStatus = 80009
	StatusWOUsed             ResStatus = 80010
	StatusEDUsed             ResStatus = 80011
	StatusEDValueUsed        ResStatus = 80012
	StatusEDErrorUsed        ResStatus = 80013
	StatusDDUsed             ResStatus = 80014
	StatusSIRespUsed         ResStatus = 80015
	StatusWOEpUsed           ResStatus = 80016
	StatusEDEpUsed           ResStatus = 80017
	StatusEDDpUsed           ResStatus = 80018
	StatusEDCommentUsed      ResStatus = 80019
	StatusEDSendToUsed       ResStatus = 80020
	StatusEDReviewUsed       ResStatus = 80021
	StatusDDEpUsed           ResStatus = 80022
	StatusDDDpUsed           ResStatus = 80023
	StatusSIOUsed            ResStatus = 80024
	StatusDocumentUsed       ResStatus = 80025
	StatusTRUsed             ResStatus = 80026
	StatusWOCreateUsed       ResStatus = 80027
	StatusWOModifyUsed       ResStatus = 80028
	StatusWOConfirmUsed      ResStatus = 80029
	StatusEDCreateUsed       ResStatus = 80030
	StatusEDModifyUsed       ResStatus = 80031
	StatusEDConfirmUsed      ResStatus = 80032
	StatusDDCreateUsed       ResStatus = 80033
	StatusDDModifyUsed       ResStatus = 80034
	StatusDDConfirmUsed      ResStatus = 80035
	StatusDCCreateUsed       ResStatus = 80036
	StatusDCModifyUsed       ResStatus = 80037
	StatusDocumentCreateUsed ResStatus = 80038
	StatusDocumentModifyUsed ResStatus = 80039
	StatusTCCreateUsed       ResStatus = 80040
	StatusTCModifyUsed       ResStatus = 80041
	StatusTRLecturerUsed     ResStatus = 80042
	StatusTRStudentUsed      ResStatus = 80043
	StatusTRCreateUsed       ResStatus = 80044
	StatusTRModifyUsed       ResStatus = 80045
	StatusTRConfirmUsed      ResStatus = 80046
	StatusOPQuotaCreateUsed  ResStatus = 80047
	StatusOPQuotaModifyUsed  ResStatus = 80048
	StatusOPQuotaConfirmUsed ResStatus = 80049
	StatusLIDCreateUsed      ResStatus = 80050
	StatusLIDModifyUsed      ResStatus = 80051
	StatusLIDConfirmUsed     ResStatus = 80052
	StatusLIDRecipientUsed   ResStatus = 80053
	StatusDeptCreateUsed     ResStatus = 80054
	StatusDeptModifyUsed     ResStatus = 80055
	StatusOPCreateUsed       ResStatus = 80056
	StatusOPModifyUsed       ResStatus = 80057
	StatusSICCreateUsed      ResStatus = 80058
	StatusSICModifyUsed      ResStatus = 80059
	StatusSICreateUsed       ResStatus = 80060
	StatusSIModifyUsed       ResStatus = 80061
	StatusUDCCreateUsed      ResStatus = 80062
	StatusUDCModifyUsed      ResStatus = 80063
	StatusUDDCreateUsed      ResStatus = 80064
	StatusUDDModifyUsed      ResStatus = 80065
	StatusEICCreateUsed      ResStatus = 80066
	StatusEICModifyUsed      ResStatus = 80067
	StatusEIDCreateUsed      ResStatus = 80068
	StatusEIDModifyUsed      ResStatus = 80069
	StatusRLCreateUsed       ResStatus = 80070
	StatusRLModifyUsed       ResStatus = 80071
	StatusLPCreateUsed       ResStatus = 80072
	StatusLPModifyUsed       ResStatus = 80073
	StatusEITCreateUsed      ResStatus = 80074
	StatusEITModifyUsed      ResStatus = 80075
	StatusUserCreateUsed     ResStatus = 80076
	StatusUserModifyUsed     ResStatus = 80077
	StatusRoleCreateUsed     ResStatus = 80078
	StatusRoleModifyUsed     ResStatus = 80079
	StatusTRDeptUsed         ResStatus = 80080
	StatusLIDDeptUsed        ResStatus = 80081
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
	StatusSICNameExist:     "On-Site archive category name exists",
	StatusSICFatherSelf:    "The parent category cannot be itself",
	StatusSICFatherCircle:  "Loop detected in parent category hierarchy",
	StatusSICLowLevelExist: "Sub-category exist",
	// On-site Archive (10600-10699)
	StatusSICodeExist: "On-site archive code already exist",
	//User-defined Archive Category (10700-10799)
	StatusUDCNameExist: "User-defined archive category name already exist",
	// User-defined Archive (10800-10899)
	StatusUDDCodeExist: "User-defined archive code already exist",
	// Execution Project Category(10900-10999)
	StatusEICNameExist:     "Execution project category name already exists",
	StatusEICFatherSelf:    "The parent category cannot be itself",
	StatusEICFatherCircle:  "Loop detected in parent category hierarchy",
	StatusEICLowLevelExist: "Sub-category exists",
	// Execution Project(11000-11099)
	StatusEIDCodeExist:        "Execution project code already exists",
	StatusEIDChangeResultType: "Referenced items cannot have their result type modified",
	// Execution Project Template (11100-11199)
	StatusEITCodeExist: "Execution project template code already exists",
	// File Archive (11200-11299)
	StatusFileOpenFailed:   "Failed to open file",
	StatusFileUploadFailed: "File upload failed",
	StatusFileGetUrlFailed: "Failed to get file URL",
	StatusFileNotExist:     "File does not exist",
	// Work Order (11300-11399)
	StatusWOOtherEdit: "Work Order has been modified by another user",
	// Execution Order (11400-11499)
	StatusEDBodyNoConfirm:       "There are unconfirmed row in the execution order body",
	StatusErrorDisposeCompleted: "The issue has been resolved",
	// Message (11500-11599)
	StatusMsgOnlyReadSelf: "Messages can only be read by the recipient themselves",
	// Risk Level (11600-11699)
	StatusRLNameExist: "Resk level name already exists",
	// Document Archive Category (11700-11799)
	StatusDCNameExist:     "Document Archive Category name already exists",
	StatusDCLowLevelExist: "Sub-category exists",
	StatusDCFatherSelf:    "The parent category cannot be itself",
	// Document Archive (11800-11899)
	StatusDocumentNoFile: "The Document Archive has no attachments",
	// Job Profile (11900-11999)(11900-11999)
	StatusOPNameExist: "Job Profile name already exists",
	// Training Course (12000-12099)
	StatusTCNameExist: "Training Course name already exists",
	// Labor Protection Equipment Archive (12100-12199)
	StatusLPCodeExist: "Labor Protection Equipment Archive code already exists",
	// Labor Protection Supplies Post Quota (12200-12199)
	StatusOPQuotaExist: "A labor protection supplies post quota for the same period already exists",
	// Referenced （80000-89999）
	StatusUDDUsed:            "Referenced by User Defined Archive",
	StatusEIDUsed:            "Referenced by Eexcution Project",
	StatusEIDDefaultUsed:     "Referenced by Execution Project default value",
	StatusEIDErrorUsed:       "Referenced by Execution Project error value",
	StatusEITDefaultUsed:     "Referenced by Execution Project Template default values",
	StatusEITErrorUsed:       "Referenced by Execution Project Template error values",
	StatusEITUsed:            "Referenced by Execution Project Template",
	StatusUserUsed:           "Referenced by User Archive",
	StatusDeptLeaderUsed:     "Referenced by Department Archive",
	StatusSIUsed:             "Referenced by On-site Archive",
	StatusWOUsed:             "Referenced by Work Order",
	StatusEDUsed:             "Referenced by Execution Order",
	StatusEDValueUsed:        `Referenced by Execution Order value`,
	StatusEDErrorUsed:        "Referenced by Execution Order error value",
	StatusDDUsed:             "Referenced by Issue Handling Form",
	StatusSIRespUsed:         "Referenced by On-site Archive owner",
	StatusWOEpUsed:           "Referenced by Execution Order executor",
	StatusEDEpUsed:           "Referenced by Execution order head executor",
	StatusEDDpUsed:           "Referenced by Execution Order body handler",
	StatusEDCommentUsed:      "Referenced by Execution Order annotator",
	StatusEDSendToUsed:       "Referenced by Execution Order annotation message recipient",
	StatusEDReviewUsed:       "Referenced by Execution Order reviewer",
	StatusDDEpUsed:           "Referenced by Issue Handling Form executor",
	StatusDDDpUsed:           "Referenced by Issue Handling Form handler",
	StatusSIOUsed:            "Referenced by On-site Archive options",
	StatusDocumentUsed:       "Referenced by Document Archive",
	StatusTRUsed:             "Referenced by Training Record",
	StatusWOCreateUsed:       "Referenced by Work Order creator",
	StatusWOModifyUsed:       "Referenced by Work Order editor",
	StatusWOConfirmUsed:      "Referenced by Work Order confirmer",
	StatusEDCreateUsed:       "Referenced by Execution Order creator",
	StatusEDModifyUsed:       "Referenced by Execution Order editor",
	StatusEDConfirmUsed:      "Referenced by Execution Order confirmer",
	StatusDDCreateUsed:       "Referenced by Issue Handling Form creator",
	StatusDDModifyUsed:       "Referenced by Issue Handling Form editor",
	StatusDDConfirmUsed:      "Referenced by Issue Handling Form confirmer",
	StatusDCCreateUsed:       "Referenced by Document Archive Category creator",
	StatusDCModifyUsed:       "Referenced by Document Archive Category editor",
	StatusDocumentCreateUsed: "Referenced by Document Archive creator",
	StatusDocumentModifyUsed: "Referenced by Document Archive editor",
	StatusTCCreateUsed:       "Referenced by Training Course creator",
	StatusTCModifyUsed:       "Referenced by Training Course editor",
	StatusTRLecturerUsed:     "Referenced by Training Record speaker",
	StatusTRStudentUsed:      "Referenced by Training Record trainee",
	StatusTRCreateUsed:       "Referenced by Training Record creator",
	StatusTRModifyUsed:       "Referenced by Training Record editor",
	StatusTRConfirmUsed:      "Referenced by Training Record confirmer",
	StatusOPQuotaCreateUsed:  "Referenced by Labor Protection Equipment Archive creator",
	StatusOPQuotaModifyUsed:  "Referenced by Labor Protection Equipment Archive editor",
	StatusOPQuotaConfirmUsed: "Referenced by Labor Protection Equipment Archive confirmer",
	StatusLIDCreateUsed:      "Referenced by PPE Issuance Form creator",
	StatusLIDModifyUsed:      "Referenced by PPE Issuance Form editor",
	StatusLIDConfirmUsed:     "Referenced by PPE Issuance Form confirmer",
	StatusLIDRecipientUsed:   "Referenced by PPE Issuance Form recipient",
	StatusDeptCreateUsed:     "Referenced by Department Arctive creator",
	StatusDeptModifyUsed:     "Referenced by Department Arctive editor",
	StatusOPCreateUsed:       "Referenced by Job Profile creator",
	StatusOPModifyUsed:       "Referenced by Job Profile editor",
	StatusSICCreateUsed:      "Referenced by On-site Archive Category creator",
	StatusSICModifyUsed:      "Referenced by On-site Archive Category editor",
	StatusSICreateUsed:       "Referenced by On-site Archive creator",
	StatusSIModifyUsed:       "Referenced by On-site Archive editor",
	StatusUDCCreateUsed:      "Referenced by User-defined Archive Category creator",
	StatusUDCModifyUsed:      "Referenced by User-defined Archive Category editor",
	StatusUDDCreateUsed:      "Referenced by User-defined Archive creator",
	StatusUDDModifyUsed:      "Referenced by User-defined Archive editor",
	StatusEICCreateUsed:      "Referenced by Execution Project Category creator",
	StatusEICModifyUsed:      "Referenced by Execution Project Category editor",
	StatusEIDCreateUsed:      "Referenced by Execution Project creator",
	StatusEIDModifyUsed:      "Referenced by Execution Project editor",
	StatusRLCreateUsed:       "Referenced by Risk Level creator",
	StatusRLModifyUsed:       "Referenced by Risk Level editor",
	StatusLPCreateUsed:       "Referenced by Labor Protection Equipment Archive creator",
	StatusLPModifyUsed:       "Referenced by Labor Protection Equipment Archive editor",
	StatusEITCreateUsed:      "Referenced by Execution Project Template creator",
	StatusEITModifyUsed:      "Referenced by Execution Project Template editor",
	StatusUserCreateUsed:     "Referenced by User Archive creator",
	StatusUserModifyUsed:     "Referenced by User Archive editor",
	StatusRoleCreateUsed:     "Referenced by Role Archive creator",
	StatusRoleModifyUsed:     "Referenced by Role Archive editor",
	StatusTRDeptUsed:         "Referenced by Training Records head department",
	StatusLIDDeptUsed:        "Referenced by the issuing department in the PPE Issuance Form",
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
