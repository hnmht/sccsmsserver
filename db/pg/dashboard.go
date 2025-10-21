package pg

import (
	"sccsmsserver/i18n"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// Give Work Order Count struct
// 0 free 1 confirmed 2 executing 3 completed
type GiveWO struct {
	FreeCount      int32 `json:"freeCount"`
	ConfirmedCount int32 `json:"confirmedCount"`
	ExecutingCount int32 `json:"exectutingCount"`
	CompletedCount int32 `json:"completedCount"`
}

// Recive Work Order Count struct
type ReciveWO struct {
	Count           int32 `json:"count"`
	UnFinishedCount int32 `json:"unFinishedCount"`
}

// Discovered Issue Count struct
type DiscoveredIssue struct {
	Count    int32 `json:"count"`
	Finished int32 `json:"finished"`
}

// Process the Issue struct
type ProcessIssue struct {
	CompletedCount  int32 `json:"completedCount"`
	UnFinishedCount int32 `json:"unFinishedCount"`
}

// Issue Item detail struct
type IssueItem struct {
	EOBID       int32  `json:"eoBID"`       // Execution Order Row ID
	EPAID       int32  `json:"epaID"`       // Execution Project ID
	EPACode     string `json:"epaCode"`     // Execution Project Code
	EPAName     string `json:"epaName"`     // Execution Project Name
	EPCID       int32  `json:"epcID"`       // Execution Project Category ID
	EPCName     string `json:"epcName"`     // Execution Project Category Name
	CSAID       int32  `json:"csaID"`       // Construction Site ID
	CSACode     string `json:"csaCode"`     // Construction Site Code
	CSAName     string `json:"csaName"`     // Construction Site Name
	CSCID       int32  `json:"cscID"`       // Construction Site Category ID
	CSCName     string `json:"cscName"`     // Construction Site Category Name
	RespID      int32  `json:"respID"`      // Response Person ID
	RespCode    string `json:"respCode"`    // Response Person Code
	RespName    string `json:"respName"`    // Resopnse Person Name
	RLID        int32  `json:"rlID"`        // Risk Level ID
	RLName      string `json:"rlName"`      // Risk Level Name
	RLColor     string `json:"rlColor"`     // Risk Level Color
	IsRectify   string `json:"isRectify"`   // Is On-Site correction performed
	IsFinish    string `json:"isfinish"`    // Is handle
	CreatorID   int32  `json:"creatorID"`   // Fixer ID
	CreatorCode string `json:"creatorCode"` // Fixer Code
	CreatorName string `json:"creatorName"` // Fixer Name
}

// User Reviewed of the Execution Order Record
type ReviewedEORecord struct {
	ID             int32  `json:"id"`
	HID            int32  `json:"hid"`
	BillNumber     string `json:"billNumber"`
	StartTime      string `json:"startTime"`
	EndTime        string `json:"endTime"`
	ConsumeSeconds int32  `json:"consumeSeconds"`
	CSAID          int32  `json:"csaID"`
	CSACode        string `json:"csaCode"`
	CSAName        string `json:"csaName"`
	CreatorID      int32  `json:"creatorID"`
	CreateUSerCode string `json:"creatorCode"`
	CreatorName    string `json:"creatorName"`
}

// User's Execution Order Reviewed by other User Record
type BeReviewedItem struct {
	ID             int32  `json:"id"`
	HID            int32  `json:"hid"`
	BillNumber     string `json:"billNumber"`
	StartTime      string `json:"startTime"`
	EndTime        string `json:"endTime"`
	ConsumeSeconds int32  `json:"consumeSeconds"`
	CSAID          int32  `json:"csaID"`
	CSACode        string `json:"csaCode"`
	CSAName        string `json:"csaName"`
	ReviewerID     int32  `json:"reviewerID"`
	ReviewerCode   string `json:"reviewerCode"`
	ReviewerName   string `json:"reviewerName"`
}

// DashBoard Data struct
type DashBoardData struct {
	StartDate       string             `json:"startDate"`
	EndDate         string             `json:"endDate"`
	GiveWO          GiveWO             `json:"giveWo"`
	ReciveWO        ReciveWO           `json:"reciveWo"`
	DiscoveredIssue DiscoveredIssue    `json:"discoveredIssue"`
	ProcessIssue    ProcessIssue       `json:"processIssue"`
	IssueItems      []IssueItem        `json:"issueItems"`
	ReviewedItems   []ReviewedEORecord `json:"reviewedItems"`
	BeReviewedItems []BeReviewedItem   `json:"beReviewedItems"`
}

// Risk Count struct
type RiskCount struct {
	OccYear     string    `json:"occYear"`
	OccMonth    string    `json:"occMonth"`
	OccWeek     string    `json:"occWeek"`
	OccDay      string    `json:"occDay"`
	RiskLevel   RiskLevel `json:"riskLevel"`
	TotalNumber int32     `json:"totalNumber"`
}

// Risk Trend Data struct
type RiskTrendData struct {
	StartDate  string      `json:"startDate"`
	EndDate    string      `json:"endDate"`
	RiskTrends []RiskCount `json:"riskTrends"`
}

// Get Dashboard Data
func (dd *DashBoardData) Get(userID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	resStatus, err = dd.GiveWO.Get(userID, dd.StartDate, dd.EndDate)
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	resStatus, err = dd.ReciveWO.Get(userID, dd.StartDate, dd.EndDate)
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	resStatus, err = dd.DiscoveredIssue.Get(userID, dd.StartDate, dd.EndDate)
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	resStatus, err = dd.ProcessIssue.Get(userID, dd.StartDate, dd.EndDate)
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	dd.IssueItems, resStatus, err = GetIssueItems(dd.StartDate, dd.EndDate)
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	dd.ReviewedItems, resStatus, err = GetReviewedRecords(userID, dd.StartDate, dd.EndDate)
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	dd.BeReviewedItems, resStatus, err = GetBeReviewedItems(userID, dd.StartDate, dd.EndDate)
	if resStatus != i18n.StatusOK || err != nil {
		return
	}

	return
}

// Statistics on Work Orders Issued by Users
func (gw *GiveWO) Get(userID int32, startDate string, endDate string) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var build strings.Builder
	// Concatenate SQL strings
	build.WriteString(` and (b.creatorid=`)
	build.WriteString(strconv.Itoa(int(userID)))
	if startDate != "" {
		build.WriteString(" and b.starttime>=")
		build.WriteString(startDate)
	}
	if endDate != "" {
		build.WriteString(" and b.starttime<=")
		build.WriteString(endDate)
	}
	build.WriteString(`)`)
	sqlString := build.String()
	// Get the count of Work Order Rows in "free" status
	build.Reset()
	build.WriteString(`select count(b.id) as freecount 
	from workorder_b as b
	left join workorder_h as h on b.hid = h.id
	where (b.dr=0 and h.dr=0 and b.status=0)`)
	build.WriteString(sqlString)
	freeSql := build.String()
	err = db.QueryRow(freeSql).Scan(&gw.FreeCount)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GiveWO.Get db.QueryRow(freesql) failed", zap.Error(err))
		return
	}
	build.Reset()

	// Get the count of Work Order Rows in "confirmed" status
	build.WriteString(`select count(b.id) as confirmedcount 
	from workorder_b as b
	left join workorder_h as h on b.hid = h.id
	where (b.dr=0 and h.dr=0 and b.status=1)`)
	build.WriteString(sqlString)
	confirmSql := build.String()
	err = db.QueryRow(confirmSql).Scan(&gw.ConfirmedCount)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GiveWO.Get db.QueryRow(confirmSql) failed", zap.Error(err))
		return
	}
	build.Reset()

	// Get the count of Work Order Rows in "executing" status
	build.WriteString(`select count(b.id) as executingcount 
	from workorder_b as b
	left join workorder_h as h on b.hid = h.id
	where (b.dr=0 and h.dr=0 and b.status=2)`)
	build.WriteString(sqlString)
	executeSql := build.String()
	err = db.QueryRow(executeSql).Scan(&gw.ExecutingCount)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GiveWO.Get db.QueryRow(executeSql) failed", zap.Error(err))
		return
	}
	build.Reset()
	// Get the count of Work Order rows in "completed" status
	build.WriteString(`select count(b.id) as completedcount 
	from workorder_b as b
	left join workorder_h as h on b.hid = h.id
	where (b.dr=0 and h.dr=0 and b.status=3)`)
	build.WriteString(sqlString)
	finishedSql := build.String()
	err = db.QueryRow(finishedSql).Scan(&gw.CompletedCount)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GiveWO.Get db.QueryRow(finishedSql) failed", zap.Error(err))
		return
	}
	build.Reset()

	return
}

// Statistics on Work Order Received by User
func (rw *ReciveWO) Get(userID int32, startDate string, endDate string) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var build strings.Builder
	// Concatenate SQL strings
	build.WriteString(` and (b.executorid=`)
	build.WriteString(strconv.Itoa(int(userID)))
	if startDate != "" {
		build.WriteString(" and b.startTime>=")
		build.WriteString(startDate)
	}
	if endDate != "" {
		build.WriteString(" and b.startTime<=")
		build.WriteString(endDate)
	}
	build.WriteString(`)`)
	sqlString := build.String()
	build.Reset()
	// Statistics the number of Work Order rows Recieved by User
	build.WriteString(`select count(b.id) as count
	from workorder_b as b
	left join workorder_h as h on b.hid = h.id
	where (b.dr=0 and h.dr=0 and b.status>0)`)
	build.WriteString(sqlString)
	allSql := build.String()
	err = db.QueryRow(allSql).Scan(&rw.Count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ReviveWOItem.Get db.QueryRow(allSql) failed", zap.Error(err))
		return
	}
	build.Reset()

	// Count of Work Order rows recieved but not executed by user
	build.WriteString(`select count(b.id) as count
	from workorder_b as b
	left join workorder_h as h on b.hid = h.id
	where (b.dr=0 and h.dr=0 and b.status=1)`)
	build.WriteString(sqlString)
	unFinishedSql := build.String()
	err = db.QueryRow(unFinishedSql).Scan(&rw.UnFinishedCount)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ReviveWOItem.Get db.QueryRow(unFinishedSql) failed", zap.Error(err))
		return
	}
	build.Reset()

	return i18n.StatusOK, nil
}

// Statistics on Issues Discovered by User
func (dp *DiscoveredIssue) Get(userID int32, startDate string, endDate string) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var build strings.Builder
	// Concatenate SQL strings
	build.WriteString(` and (b.creatorid=`)
	build.WriteString(strconv.Itoa(int(userID)))
	if startDate != "" {
		build.WriteString(" and h.startTime>=")
		build.WriteString(startDate)
	}
	if endDate != "" {
		build.WriteString(" and h.startTime<=")
		build.WriteString(endDate)
	}
	build.WriteString(`)`)
	sqlString := build.String()
	build.Reset()
	// Total Count of Issues Discovered by Users
	build.WriteString(`select count(b.id) as count  
	from executionorder_b as b
	left join executionorder_h as h on b.hid = h.id
	where (b.dr=0 and h.dr=0 and b.isissue=1)`)
	build.WriteString(sqlString)
	allSql := build.String()
	err = db.QueryRow(allSql).Scan(&dp.Count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DiscoveredIssue.Get db.QueryRow(allSql) failed", zap.Error(err))
		return
	}
	build.Reset()
	// Total Count of Issues Discovered and Already Resolved by User
	build.WriteString(`select count(b.id) as count
	from executionorder_b as b
	left join executionorder_h as h on b.hid = h.id
	where (b.dr=0 and h.dr=0 and b.isissue=1) and (b.isRectify = 1 or b.isfinish=1)`)
	build.WriteString(sqlString)
	finishedSql := build.String()
	err = db.QueryRow(finishedSql).Scan(&dp.Finished)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DiscoveredIssue.Get db.QueryRow(finishedSql) failed", zap.Error(err))
		return
	}
	build.Reset()

	return
}

// Statistics on Issues where the User is the Issue Owner
func (dp *ProcessIssue) Get(userID int32, startDate string, endDate string) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var build strings.Builder
	// Concatenate SQL string
	build.WriteString(` and (b.issueownerid=`)
	build.WriteString(strconv.Itoa(int(userID)))
	if startDate != "" {
		build.WriteString(" and b.handlestarttime>=")
		build.WriteString(startDate)
	}
	if endDate != "" {
		build.WriteString(" and b.handlestarttime<=")
		build.WriteString(endDate)
	}
	build.WriteString(`)`)
	sqlString := build.String()
	build.Reset()
	// Summary of the number of issues resolved by the user as the responsible person
	build.WriteString(`select count(b.id) as count
	from executionorder_b as b
	left join executionorder_h as h on b.hid = h.id
	where (b.dr=0 and h.dr=0 and b.ishandle=1 and b.isfinish=1)`)
	build.WriteString(sqlString)
	finishedSql := build.String()
	err = db.QueryRow(finishedSql).Scan(&dp.CompletedCount)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ProcessIssue.Get db.QueryRow(finishedSql) failed", zap.Error(err))
		return
	}
	build.Reset()
	// Summary of the number of unresolved issues where the user is the responsible person
	build.WriteString(`select count(b.id) as count
	from executionorder_b as b
	left join executionorder_h as h on b.hid = h.id
	where (b.dr=0 and h.dr=0 and b.ishandle=1 and b.isfinish=0)`)
	build.WriteString(sqlString)
	unFinishedSql := build.String()
	err = db.QueryRow(unFinishedSql).Scan(&dp.UnFinishedCount)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ProcessIssue.Get db.QueryRow(unFinishedSql) failed", zap.Error(err))
		return
	}
	build.Reset()

	return
}

// Get Issue Items List
func GetIssueItems(startDate string, endDate string) (iis []IssueItem, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	iis = make([]IssueItem, 0)
	var build strings.Builder
	// Concatenate SQL strings
	if startDate != "" {
		build.WriteString(" and h.startTime>=")
		build.WriteString(startDate)
	}
	if endDate != "" {
		build.WriteString(" and h.startTime<=")
		build.WriteString(endDate)
	}
	sqlString := build.String()
	build.Reset()

	// Retrieve Issue Items from database
	build.WriteString(`select b.id as eoBID, 
	b.epaid as epaID,
	epa.code as epacode,
	epa.name as epaname,
	epa.epcid as epcid,
	epc.name as epcname,
	h.csaid as csaid,
	csa.code as csacode,
	csa.name as csaname,
	csa.cscid as cscid,
	csc.name as cscname,
	csa.resppersonid as respid,
	respperson.code as respcode,
	respperson.name as respname,
	rl.id as rlid,
	rl.name as rlname,
	rl.color as rlcolor,
	b.isrectify as isRectify,
	b.isfinish as isfinish,
	h.creatorid as creatorid,
	creator.code as creatorcode,
	creator.name as creatorname
	from executionorder_b as b
	left join executionorder_h as h on b.hid = h.id
	left join sysuser as creator on h.creatorID = creator.id
	left join epa on b.epaid = epa.id
	left join risklevel as rl on b.risklevelid = rl.id
	left join epc on epa.epcid = epc.id
	left join csa on h.csaid = csa.id
	left join sysuser as respperson on csa.resppersonid = respperson.id
	left join csc on csa.cscid = csc.id
	where (b.dr=0 and h.dr=0 and b.isissue=1)`)
	build.WriteString(sqlString)

	sqlStr := build.String()
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetIssueItems db.Query(sqlStr) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()
	// Extract IssueItem row by row
	for rows.Next() {
		var ii IssueItem
		err = rows.Scan(&ii.EOBID, &ii.EPAID, &ii.EPACode, &ii.EPAName, &ii.EPCID,
			&ii.EPCName, &ii.CSAID, &ii.CSACode, &ii.CSAName, &ii.CSCID,
			&ii.CSCName, &ii.RespID, &ii.RespCode, &ii.RespName, &ii.RLID,
			&ii.RLName, &ii.RLColor, &ii.IsRectify, &ii.IsFinish, &ii.CreatorID,
			&ii.CreatorCode, &ii.CreatorName)
		if err != nil {
			zap.L().Error("GetIssueItems rows.Scan failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		iis = append(iis, ii)
	}
	return
}

// Get User Reviewed Execution Order records
func GetReviewedRecords(userID int32, startDate string, endDate string) (rrs []ReviewedEORecord, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	rrs = make([]ReviewedEORecord, 0)
	var build strings.Builder
	// Concatenate SQL strings
	build.WriteString(` and (r.creatorID=`)
	build.WriteString(strconv.Itoa(int(userID)))
	if startDate != "" {
		build.WriteString(" and r.startTime>=")
		build.WriteString(startDate)
	}
	if endDate != "" {
		build.WriteString(" and r.startTime<=")
		build.WriteString(endDate)
	}
	build.WriteString(`)`)
	sqlString := build.String()
	build.Reset()
	// Retrieve User Reviewed records from database
	build.WriteString(`select r.id as id,
	r.hid as hid,
	r.billnumber as billnumber,
	r.starttime as starttime,
	r.endtime as endtime,
	r.consumeseconds as consumeseconds,
	h.csaid as csaid,
	csa.code as csacode,
	csa.name as csaname,
	h.creatorid as creatorid,
	creator.code as creatorcode,
	creator.name as creatorname
	from executionorder_review as r
	left join executionorder_h as h on r.hid = h.id
	left join csa on h.csaid = csa.id
	left join sysuser as creator on h.creatorid = creator.id
	where (r.dr=0 and h.dr=0)`)
	build.WriteString(sqlString)
	sqlStr := build.String()

	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetReviewedRecords db.Query(sqlStr) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	// Extract Reviewed Records row by row
	for rows.Next() {
		var rr ReviewedEORecord
		err = rows.Scan(&rr.ID, &rr.HID, &rr.BillNumber, &rr.StartTime, &rr.EndTime,
			&rr.ConsumeSeconds, &rr.CSAID, &rr.CSACode, &rr.CSAName, &rr.CreatorID,
			&rr.CreateUSerCode, &rr.CreatorName)
		if err != nil {
			zap.L().Error("GetReviewedRecords rows.Scan failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		rrs = append(rrs, rr)
	}
	return
}

// Statistics of Records where the Execution Order Filled by the user Reviewed by others
func GetBeReviewedItems(userID int32, startDate string, endDate string) (brs []BeReviewedItem, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	brs = make([]BeReviewedItem, 0)
	var build strings.Builder
	// Concatenate SQL strings
	build.WriteString(` and (h.creatorID=`)
	build.WriteString(strconv.Itoa(int(userID)))
	if startDate != "" {
		build.WriteString(" and r.startTime>=")
		build.WriteString(startDate)
	}
	if endDate != "" {
		build.WriteString(" and r.startTime<=")
		build.WriteString(endDate)
	}
	build.WriteString(`)`)
	sqlString := build.String()
	build.Reset()
	// Get the Reviewed Records from database
	build.WriteString(`select r.id as id,
	r.hid as hid,
	r.billnumber as billnumber,
	r.starttime as starttime,
	r.endtime as endtime,
	r.consumeseconds as consumeseconds,
	h.csaid as csaid,
	csa.code as csacode,
	csa.name as csaname,
	r.creatorid as reviewerid,
	reviewer.code as reviewercode,
	reviewer.name as reviewername
	from executionorder_review as r
	left join executionorder_h as h on r.hid = h.id
	left join csa as csa on h.csaid = csa.id
	left join sysuser as reviewer on r.creatorid = reviewer.id
	where (r.dr=0 and h.dr=0) `)
	build.WriteString(sqlString)
	sqlStr := build.String()

	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetBeReviewedItems db.Query(sqlStr) failed", zap.Error(err))
		resStatus = i18n.StatusOK
		return
	}
	defer rows.Close()
	// Extract Reviewed Records row by row
	for rows.Next() {
		var br BeReviewedItem
		err = rows.Scan(&br.ID, &br.HID, &br.BillNumber, &br.StartTime, &br.EndTime,
			&br.ConsumeSeconds, &br.CSAID, &br.CSACode, &br.CSAName, &br.ReviewerID,
			&br.ReviewerCode, &br.ReviewerName)
		if err != nil {
			zap.L().Error("GetBeReviewedItems rows.Scan failed", zap.Error(err))
			resStatus = i18n.StatusOK
			return
		}
		brs = append(brs, br)
	}
	return
}

// Summarize Risk Records
func GetRiskRecords(startDate string, endDate string) (rcs []RiskCount, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	rcs = make([]RiskCount, 0)
	var build strings.Builder
	// Concatenate SQL strings
	if startDate != "" {
		build.WriteString(" and h.billdate>=")
		build.WriteString(startDate)
	}
	if endDate != "" {
		build.WriteString(" and h.billdate<='")
		build.WriteString(endDate)
		build.WriteString("'")
	}

	build.WriteString("	group by occYear,occMonth,occDay,rlid")
	build.WriteString(" order by occDay")
	conString := build.String()
	build.Reset()

	// Summarize Risk Records by Occurrence date
	build.WriteString(`
	select EXTRACT(YEAR from h.billdate) as occyear,
	EXTRACT(MONTH from h.billdate) as occmonth,
	EXTRACT(WEEK from h.billdate) as occweek,
	h.billdate as occday,
	b.risklevelid as rlid,
	count(b.id) as itemnumber
	from executionorder_b as b
	left join executionorder_h as h on b.hid = h.id
	left join riskLevel as rl on b.risklevelid = rl.id	
	where (b.dr=0 and h.dr=0 and b.isissue=1)
	`)
	build.WriteString(conString)
	sqlStr := build.String()
	// Retrieve Records from database
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetRiskRecords db.Query(sqlStr) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()
	// Extract Records row by row
	for rows.Next() {
		var rc RiskCount
		err = rows.Scan(&rc.OccYear, &rc.OccMonth, &rc.OccWeek, &rc.OccDay, &rc.RiskLevel.ID, &rc.TotalNumber)
		if err != nil {
			zap.L().Error("GetRiskRecords rows.Scan failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get Risk Level Details
		if rc.RiskLevel.ID > 0 {
			resStatus, err := rc.RiskLevel.GetRLInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return rcs, resStatus, err
			}
		}
		rcs = append(rcs, rc)
	}
	return rcs, i18n.StatusOK, nil
}

// Get Risk Trend data
func (rtd *RiskTrendData) Get() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	rtd.RiskTrends, resStatus, err = GetRiskRecords(rtd.StartDate, rtd.EndDate)
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	return
}
