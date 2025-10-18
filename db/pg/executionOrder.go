package pg

import (
	"math"
	"sccsmsserver/i18n"
	"sccsmsserver/setting"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Execution Order Struct
type ExecutionOrder struct {
	HID              int32               `db:"id" json:"id"`
	BillNumber       string              `db:"billnumber" json:"billNumber"`
	BillDate         time.Time           `db:"billdate" json:"billDate"`
	Department       SimpDept            `db:"deptid" json:"department"`
	Description      string              `db:"description" json:"description"`
	Status           int16               `db:"status" json:"status"`         //0 free 1 confirmed 2 executing 3 completed
	SourceType       string              `db:"sourcetype" json:"sourceType"` //Source Type: di:Direct Input wo: Work Order
	SourceBillNumber string              `db:"sourcebillnumber" json:"sourceBillNumber"`
	SourceHid        int32               `db:"sourcehid" json:"sourceHID"`
	SourceRowNumber  int32               `db:"sourcerownumber" json:"sourceRowNumber"`
	SourceBid        int32               `db:"sourcebid" json:"sourceBID"`
	SourceRowTs      time.Time           `json:"sourceRowTs"`
	StartTime        time.Time           `db:"starttime" json:"startTime"`
	EndTime          time.Time           `db:"endtime" json:"endTime"`
	CSA              ConstructionSite    `db:"csaid" json:"csa"`
	Executor         Person              `db:"executorid" json:"executor"`
	EPT              EPT                 `db:"eptid" json:"ept"`
	AllowAddRow      int16               `db:"allowaddrow" json:"allowAddRow"`
	AllowDelRow      int16               `db:"allowdelrow" json:"allowDelRow"`
	Body             []ExecutionOrderRow `json:"body"`
	IssueNumber      int32               `json:"IssueNumber"`
	ReviewedNumber   int16               `json:"reviewedNumber"`
	ReviewedSeconds  int32               `json:"reviewedSeconds"`
	CreateDate       time.Time           `db:"createtime" json:"createDate"`
	Creator          Person              `db:"creatorid" json:"creator"`
	ConfirmDate      time.Time           `db:"confirmtime" json:"confirmDate"`
	Confirmer        Person              `db:"confirmerid" json:"confirmer"`
	ModifyDate       time.Time           `db:"modifytime" json:"modifyDate"`
	Modifier         Person              `db:"modifierid" json:"modifier"`
	Ts               time.Time           `db:"ts" json:"ts"`
	Dr               int16               `db:"dr" json:"dr"`
}

// Execution Order Row struct
type ExecutionOrderRow struct {
	BID                int32            `db:"id" json:"id"`
	HID                int32            `db:"hid" json:"hid"`
	RowNumber          int32            `db:"rownumber" json:"rowNumber"`
	EPA                ExecutionProject `db:"epaid" json:"epa"`
	AllowDelRow        int16            `db:"allowdelrow" json:"allowDelRow"`
	ExecutionValue     string           `db:"executionvalue" json:"executionValue"`
	ExecutionValueDisp string           `db:"executionvaluedisp" json:"executionValueDisp"`
	Files              []VoucherFile    `json:"files"`
	Description        string           `db:"description" json:"description"`
	EpaDescription     string           `db:"epadescription" json:"epaDescription"`
	IsCheckError       int16            `db:"ischeckerror" json:"isCheckError"`
	ErrorValue         string           `db:"errorvalue" json:"errorValue"`
	ErrorValueDisp     string           `db:"errorvaluedisp" json:"errorValueDisp"`
	IsRequireFile      int16            `db:"isrequirefile" json:"isRequireFile"`
	IsOnsitePhoto      int16            `db:"isonsitephoto" json:"isOnSitePhoto"`
	IsIssue            int16            `db:"isissue" json:"isIssue"`
	IsRectify          int16            `db:"isrectify" json:"isRectify"` // On-Site correction performed
	IsHandle           int16            `db:"ishandle" json:"isHandle"`   // 0 No 1 Yes
	IssueOwner         Person           `db:"issueownerid" json:"issueOwner"`
	HandleStartTime    time.Time        `db:"handlestarttime" json:"handleStartTime"`
	HandleEndTime      time.Time        `db:"handleendtime" json:"handleEndTime"`
	Status             int16            `db:"status" json:"status"`
	IsFromEPT          int16            `db:"isfromept" json:"isFromEpt"`
	IsFinish           int16            `db:"isfinish" json:"isFinish"`
	IRFID              int32            `db:"irfid" json:"irfID"`
	IRFNumber          string           `db:"irfnumber" json:"irfNumber"`
	RiskLevel          RiskLevel        `db:"risklevelid" json:"riskLevel"`
	CreateDate         time.Time        `db:"createtime" json:"createDate"`
	Creator            Person           `db:"creatorid" json:"creator"`
	ConfirmDate        time.Time        `db:"confirmtime" json:"confirmDate"`
	Confirmer          Person           `db:"confirmerid" json:"confirmer"`
	ModifyDate         time.Time        `db:"modifytime" json:"modifyDate"`
	Modifier           Person           `db:"modifierid" json:"modifier"`
	Ts                 time.Time        `db:"ts" json:"ts"`
	Dr                 int16            `db:"dr" json:"dr"`
}

// Execution Order Record for Reference by Downstreams Voucher
type ReferExecutionOrder struct {
	BID                int32            `db:"b.id" json:"id"`
	HID                int32            `db:"b.hid" json:"hid"`
	RowNumber          int32            `db:"b.rownumber" json:"rowNumber"`
	EPA                ExecutionProject `db:"b.epaid" json:"epa"`
	ExecutionValue     string           `db:"b.executionvalue" json:"executionValue"`
	ExecutionValueDisp string           `db:"b.executionvaluedisp" json:"executionValueDisp"`
	Description        string           `db:"b.description" json:"description"`
	EOFiles            []VoucherFile    `json:"eoFiles"`
	IsHandle           int16            `db:"b.ishandle" json:"isHandle"` //0 No 1 Yes
	IssueOwner         Person           `db:"b.issueownerid" json:"issueOwner"`
	HandleStartTime    string           `db:"b.handlestarttime" json:"handleStartTime"`
	HandleEndTime      string           `db:"b.handleendtime" json:"handleEndTime"`
	Status             int16            `db:"b.status" json:"status"`
	RiskLevel          RiskLevel        `db:"b.risklevelid" json:"riskLevel"`
	Dr                 int16            `db:"b.dr" json:"dr"`
	Ts                 time.Time        `db:"b.ts" json:"ts"`
	IsFinish           int16            `db:"b.isfinish" json:"isFinish"`
	BillNumber         string           `db:"h.billnumber" json:"billNumber"`
	BillDate           time.Time        `db:"h.billdate" json:"billDate"`
	Department         SimpDept         `db:"h.deptid" json:"department"`
	CSA                ConstructionSite `db:"h.csaid" json:"csa"`
	Executor           Person           `db:"h.executorid" json:"executor"`
}

// Execution Order Comment Record struct
type ExecutionOrderComment struct {
	ID         int32     `db:"id" json:"id"`
	HID        int32     `db:"hid" json:"hid"`
	BID        int32     `db:"bid" json:"bid"`
	RowNUmber  int32     `db:"rownumber" json:"rowNumber"`
	BillNumber string    `db:"billnumber" json:"billNumber"`
	SendTo     Person    `db:"sendtoid" json:"sendTo"`
	IsRead     int16     `db:"isread" json:"isRead"`
	ReadTime   time.Time `db:"readtime" json:"readTime"`
	Content    string    `db:"content" json:"content"`
	SendTime   time.Time `db:"sendtime" json:"sendTime"`
	CreateDate time.Time `db:"createtime" json:"createDate"`
	Creator    Person    `db:"creatorid" json:"creator"`
	ModifyDate time.Time `db:"modifytime" json:"modifyDate"`
	Modifier   Person    `db:"modifierid" json:"modifier"`
	Ts         time.Time `db:"ts" json:"ts"`
	Dr         int16     `db:"dr" json:"dr"`
}

// Execution Order Review Record struct
type ExecutionOrderReview struct {
	ID             int32     `db:"id" json:"id"`
	HID            int32     `db:"hid" json:"hid"`
	BillNumber     string    `db:"billnumber" json:"billNumber"`
	StartTime      time.Time `db:"starttime" json:"startTime"`
	EndTime        time.Time `db:"endtime" json:"endTime"`
	ConsumeSeconds int32     `db:"consumeseconds" json:"consumeSeconds"`
	CreateDate     time.Time `db:"createtime" json:"createDate"`
	Creator        Person    `db:"creatorid" json:"creator"`
	Ts             time.Time `db:"ts" json:"ts"`
	Dr             int16     `db:"dr" json:"dr"`
}

// Execution Order Comments Params
type EOCommentsParams struct {
	HID      int32                   `json:"hid"`
	Comments []ExecutionOrderComment `json:"comments"`
}

// Execution Order Reviews Params
type EOReviewsParams struct {
	HID     int32                  `json:"hid"`
	Reviews []ExecutionOrderReview `json:"reviews"`
}

// Execution Order List Paging struct
type EOListPaging struct {
	EOs     []ExecutionOrder `json:"eos"`
	Count   int32            `json:"count"`
	Page    int32            `json:"page"`
	PerPage int32            `json:"perPage"`
}

// Get the list of execution orders to be referenced
func GetReferEOs(queryString string) (reos []ReferExecutionOrder, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	reos = make([]ReferExecutionOrder, 0)
	var build strings.Builder
	// Assemble the SQL for checking
	build.WriteString(`select count(b.id) as rownumber
	from executionorder_b as b
	left join executionorder_h as h on b.hid = h.id
	left join epa as epa on b.epaid = epa.id
	left join sysuser as issueowner on b.issueownerid = issueowner.id
	left join sysuser as epuser on h.executorid = epuser.id
	left join sysuser as creator on h.creatorid = creator.id
	left join department as dept on h.deptid = dept.id
	where (b.ishandle=1 and b.dr = 0 and b.isfinish=0 and b.status=1)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	checkSql := build.String()
	// Check
	var rowNumber int32
	err = db.QueryRow(checkSql).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetReferEOs db.QueryRow(checkSql) failed", zap.Error(err))
		return
	}
	if rowNumber == 0 {
		resStatus = i18n.StatusResNoData
		return
	}
	if rowNumber > setting.Conf.PqConfig.MaxRecord {
		resStatus = i18n.StatusOverRecord
		return
	}
	build.Reset()
	// Assemble the SQL for data retrieval
	build.WriteString(`select b.id,b.hid,b.rownumber,b.epaid,b.executionvalue,
	b.executionvaluedisp,b.description,b.ishandle,b.issueownerid,b.handlestarttime,
	b.handleendtime,b.status,b.risklevelid, b.isfinish,b.dr,
	b.ts,h.billnumber,h.billdate,h.deptid,h.csaid,
	h.executorid 
	from executionorder_b as b
	left join executionorder_h as h on b.hid = h.id
	left join epa as epa on b.epaid = epa.id
	left join sysuser as issueowner on b.issueownerid = issueowner.id
	left join sysuser as epuser on h.executorid = epuser.id
	left join sysuser as creator on h.creatorid = creator.id
	left join department as dept on h.deptid = dept.id
	where (b.ishandle=1 and b.dr = 0 and b.isfinish=0 and (b.status=1 or b.status=2))`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	refSql := build.String()
	// Retrieve the list of Execution Orders to be referenced
	edRef, err := db.Query(refSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetReferEOs db.Query failed", zap.Error(err))
		return
	}
	defer edRef.Close()
	// Extract data row by row
	for edRef.Next() {
		var reo ReferExecutionOrder
		err = edRef.Scan(&reo.BID, &reo.HID, &reo.RowNumber, &reo.EPA.ID, &reo.ExecutionValue,
			&reo.ExecutionValueDisp, &reo.Description, &reo.IsHandle, &reo.IssueOwner.ID, &reo.HandleStartTime,
			&reo.HandleEndTime, &reo.Status, &reo.RiskLevel.ID, &reo.IsFinish, &reo.Dr,
			&reo.Ts, &reo.BillNumber, &reo.BillDate, &reo.Department.ID, &reo.CSA.ID,
			&reo.Executor.ID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetReferEOs edRef.Next() edRef.Scan() failed", zap.Error(err))
			return
		}
		// Get Execution Project details
		if reo.EPA.ID > 0 {
			resStatus, err = reo.EPA.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Department details
		if reo.Department.ID > 0 {
			resStatus, err = reo.Department.GetSimpDeptInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Risk Level details
		if reo.RiskLevel.ID > 0 {
			resStatus, err = reo.RiskLevel.GetRLInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Construction Site details
		if reo.CSA.ID > 0 {
			resStatus, err = reo.CSA.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Executor deatils
		if reo.Executor.ID > 0 {
			resStatus, err = reo.Executor.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Issue Owner details
		if reo.IssueOwner.ID > 0 {
			resStatus, err = reo.IssueOwner.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get files details
		reo.EOFiles, resStatus, err = GetEORowFiles(reo.BID)
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		reos = append(reos, reo)
	}

	return
}

// Get Execution Order list
func GetEOList(queryString string) (eos []ExecutionOrder, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var build strings.Builder
	// Assemble the SQL for checking
	build.WriteString(`select count(h.id) as rownumber
	from executionorder_h as h
	left join department on h.deptid = department.id
	left join sysuser as creator on h.creatorid = creator.id
	left join sysuser as modifier on h.modifierid = modifier.id
	where (h.dr = 0) `)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	checkSql := build.String()
	// Check
	var rowNumber int32
	err = db.QueryRow(checkSql).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEOList db.QueryRow(checkSql) failed", zap.Error(err))
		return
	}
	if rowNumber == 0 {
		resStatus = i18n.StatusResNoData
		return
	}
	if rowNumber > setting.Conf.PqConfig.MaxRecord {
		resStatus = i18n.StatusOverRecord
		return
	}
	build.Reset()

	// Assemble the SQL for data retrieval
	build.WriteString(`select h.id,h.billnumber,h.billdate,h.deptid,h.description,
	h.status,h.sourcetype,h.sourcebillnumber,h.sourcehid,h.sourcerownumber,
	h.sourcebid,h.starttime,h.endtime,h.csaid,h.executorid,
	h.eptid,h.allowaddrow,h.allowdelrow,h.createtime,h.creatorid,
	h.confirmtime,h.confirmerid,h.modifytime,h.modifierid,h.dr,
	h.ts 
	from executionorder_h as h
	left join department on h.deptid = department.id
	left join sysuser as creator on h.creatorid = creator.id
	left join sysuser as modifier on h.modifierid = modifier.id
	where (h.dr = 0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	headSql := build.String()
	// Retrieve Execution Order list from database
	headRows, err := db.Query(headSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEOList db.Query failed", zap.Error(err))
		return
	}
	defer headRows.Close()
	// Extract data row by row
	for headRows.Next() {
		var eo ExecutionOrder
		err = headRows.Scan(&eo.HID, &eo.BillNumber, &eo.BillDate, &eo.Department.ID, &eo.Description,
			&eo.Status, &eo.SourceType, &eo.SourceBillNumber, &eo.SourceHid, &eo.SourceRowNumber,
			&eo.SourceBid, &eo.StartTime, &eo.EndTime, &eo.CSA.ID, &eo.Executor.ID,
			&eo.EPT.HID, &eo.AllowAddRow, &eo.AllowDelRow, &eo.CreateDate, &eo.Creator.ID,
			&eo.ConfirmDate, &eo.Confirmer.ID, &eo.ModifyDate, &eo.Modifier.ID, &eo.Dr,
			&eo.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetEOList headRows.Next failed", zap.Error(err))
			return
		}
		// Get Execution Order Header details
		resStatus, err = eo.FillHead()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		eos = append(eos, eo)
	}
	return
}

// Get the list of Execution Order to be reviewed
func GetEOReviewList(queryString string, useID int32) (eos []ExecutionOrder, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	eos = make([]ExecutionOrder, 0)
	var build strings.Builder
	// Assemble the SQL for check
	build.WriteString(`select count(h.id) as rownumber
	from executionorder_h as h
	left join department on h.deptid = department.id
	left join sysuser as creator on h.creatorid = creator.id
	left join sysuser as modifier on h.modifierid = modifier.id
	where (h.dr = 0) `)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	checkSql := build.String()
	// Check the number of rows
	var rowNumber int32
	err = db.QueryRow(checkSql).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEOReviewList db.QueryRow(checkSql) failed", zap.Error(err))
		return
	}
	if rowNumber == 0 {
		resStatus = i18n.StatusResNoData
		return
	}
	if rowNumber > setting.Conf.PqConfig.MaxRecord {
		resStatus = i18n.StatusOverRecord
		return
	}
	build.Reset()

	// Assemble the SQL for data retrieve
	build.WriteString(`select h.id,h.billnumber,h.billdate,h.deptid,h.description,
	h.status,h.sourcetype,h.sourcebillnumber,h.sourcehid,h.sourcerownumber,
	h.sourcebid,h.starttime,h.endtime,h.csaid,h.executorid,
	h.eptid,h.allowaddrow,h.allowdelrow,h.createtime,h.creatorid,
	h.confirmtime,h.confirmerid,h.modifytime,h.modifierid,h.dr,
	h.ts,
	(select count(b.id) as errnumber from executionorder_b as b where b.hid = h.id and b.dr=0 and b.isissue=1),
	(select count(r.id) as reviewednumber from executionorder_review as r where r.hid = h.id and r.dr=0 and r.creatorid=$1),
	(select coalesce( sum(r.consumeseconds),0) as reviewedseconds  from executionorder_review as r where r.hid = h.id and r.dr=0 and r.creatorid=$1)
	from executionorder_h as h
	left join department on h.deptid = department.id
	left join sysuser as creator on h.creatorid = creator.id
	left join sysuser as modifier on h.modifierid = modifier.id
	where (h.dr = 0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	headSql := build.String()
	// Retrieve Execution Order from database
	headRows, err := db.Query(headSql, useID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEOReviewList db.Query failed", zap.Error(err))
		return
	}
	defer headRows.Close()

	for headRows.Next() {
		var eo ExecutionOrder
		err = headRows.Scan(&eo.HID, &eo.BillNumber, &eo.BillDate, &eo.Department.ID, &eo.Description,
			&eo.Status, &eo.SourceType, &eo.SourceBillNumber, &eo.SourceHid, &eo.SourceRowNumber,
			&eo.SourceBid, &eo.StartTime, &eo.EndTime, &eo.CSA.ID, &eo.Executor.ID,
			&eo.EPT.HID, &eo.AllowAddRow, &eo.AllowDelRow, &eo.CreateDate, &eo.Creator.ID,
			&eo.ConfirmDate, &eo.Confirmer.ID, &eo.ModifyDate, &eo.Modifier.ID, &eo.Dr,
			&eo.Ts, &eo.IssueNumber, &eo.ReviewedNumber, &eo.ReviewedSeconds)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetEOReviewList headRows.Next failed", zap.Error(err))
			return
		}
		// Get Execution Order Header details
		resStatus, err = eo.FillHead()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		eos = append(eos, eo)
	}
	return
}

// Get the list of Execution Orders to be reviewed by pagination
func GetEOReviewListPagination(con PagingQueryParams, userID int32) (edsp EOListPaging, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	edsp.EOs = make([]ExecutionOrder, 0)
	var build strings.Builder
	// Assemble the SQL for checking
	build.WriteString(`select count(h.id) as rownumber
	from executionorder_h as h
	left join department on h.deptid = department.id
	left join sysuser as creator on h.creatorid = creator.id
	left join sysuser as modifier on h.modifierid = modifier.id
	where (h.dr = 0) `)
	if con.QueryString != "" {
		build.WriteString(" and (")
		build.WriteString(con.QueryString)
		build.WriteString(")")
	}
	checkSql := build.String()

	// Check
	err = db.QueryRow(checkSql).Scan(&edsp.Count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEOReviewListPagination db.QueryRow(checkSql) failed", zap.Error(err))
		return
	}
	if edsp.Count == 0 {
		resStatus = i18n.StatusResNoData
		return
	}
	if edsp.Count > setting.Conf.PqConfig.MaxRecord {
		resStatus = i18n.StatusOverRecord
		return
	}
	// Recalculate pagination
	if con.PerPage > edsp.Count {
		con.Page = 0
	} else {
		var totalPage = int32(math.Ceil(float64(edsp.Count) / float64(con.PerPage)))
		if (con.Page + 1) > totalPage {
			con.Page = totalPage - 1
		}
	}
	build.Reset()

	// Assemble the SQL for data retrieve
	build.WriteString(`select h.id,h.billnumber,h.billdate,h.deptid,h.description,
	h.status,h.sourcetype,h.sourcebillnumber,h.sourcehid,h.sourcerownumber,
	h.sourcebid,h.starttime,h.endtime,h.csaid,h.executorid,
	h.eptid,h.allowaddrow,h.allowdelrow,h.createtime,h.creatorid,
	h.confirmtime,h.confirmerid,h.modifytime,h.modifierid,h.dr,
	h.ts,
	(select count(b.id) as errnumber from executionorder_b as b where b.hid = h.id and b.dr=0 and b.isissue=1),
	(select count(r.id) as reviewednumber from executionorder_review as r where r.hid = h.id and r.dr=0 and r.creatorid=$1),
	(select coalesce( sum(r.consumeseconds),0) as reviewedseconds  from executionorder_review as r where r.hid = h.id and r.dr=0 and r.creatorid=$1)
	from executionorder_h as h
	left join department on h.deptid = department.id
	left join sysuser as creator on h.creatorid = creator.id
	left join sysuser as modifier on h.modifierid = modifier.id
	where (h.dr = 0)`)
	if con.QueryString != "" {
		build.WriteString(" and (")
		build.WriteString(con.QueryString)
		build.WriteString(")")
	}
	build.WriteString(" order by h.id limit $2 offset $3")
	headSql := build.String()

	// Get Execution Order from database
	headRows, err := db.Query(headSql, userID, con.PerPage, con.Page*con.PerPage)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEOReviewListPagination db.Query failed", zap.Error(err))
		return
	}
	defer headRows.Close()
	// Extract data row by row
	for headRows.Next() {
		var eo ExecutionOrder
		err = headRows.Scan(&eo.HID, &eo.BillNumber, &eo.BillDate, &eo.Department.ID, &eo.Description,
			&eo.Status, &eo.SourceType, &eo.SourceBillNumber, &eo.SourceHid, &eo.SourceRowNumber,
			&eo.SourceBid, &eo.StartTime, &eo.EndTime, &eo.CSA.ID, &eo.Executor.ID,
			&eo.EPT.HID, &eo.AllowAddRow, &eo.AllowDelRow, &eo.CreateDate, &eo.Creator.ID,
			&eo.ConfirmDate, &eo.Confirmer.ID, &eo.ModifyDate, &eo.Modifier.ID, &eo.Dr,
			&eo.Ts, &eo.IssueNumber, &eo.ReviewedNumber, &eo.ReviewedSeconds)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetEOReviewListPagination headRows.Next failed", zap.Error(err))
			return
		}
		// Get Execution Order details
		resStatus, err = eo.FillHead()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		edsp.EOs = append(edsp.EOs, eo)
	}
	edsp.Page = con.Page
	edsp.PerPage = con.PerPage

	return
}

// Fill in the detailed information of the Execution Order Header
func (eo *ExecutionOrder) FillHead() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get Department details
	if eo.Department.ID > 0 {
		resStatus, err = eo.Department.GetSimpDeptInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Construction Site details
	if eo.CSA.ID > 0 {
		resStatus, err = eo.CSA.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Executor details
	if eo.Executor.ID > 0 {
		resStatus, err = eo.Executor.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Execution Project Template details
	if eo.EPT.HID > 0 {
		resStatus, err = eo.EPT.GetEPTHeaderByHid()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Creator details
	if eo.Creator.ID > 0 {
		resStatus, err = eo.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Confirmer details
	if eo.Confirmer.ID > 0 {
		resStatus, err = eo.Confirmer.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Modifier details
	if eo.Modifier.ID > 0 {
		resStatus, err = eo.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	return
}

// Get the Execution Order Row attachments
func GetEORowFiles(bid int32) (voucherFiles []VoucherFile, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	voucherFiles = make([]VoucherFile, 0)
	// Retrieve the attachments from executionorder_file table
	attachSql := `select id,billbid,billhid,fileid,createtime,
	creatorid,modifytime,modifierid,dr,ts 
	from executionorder_file where billbid=$1 and dr=0`
	fileRows, err := db.Query(attachSql, bid)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEORowFiles db.query(attachsql) failed", zap.Error(err))
		return
	}
	defer fileRows.Close()
	// Extract data row by row
	for fileRows.Next() {
		var f VoucherFile
		fileErr := fileRows.Scan(&f.ID, &f.BillBID, &f.BillHID, &f.File.ID, &f.CreateDate,
			&f.Creator.ID, &f.ModifyDate, &f.Modifier.ID, &f.Dr, &f.Ts)
		if fileErr != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetEORowFiles fileRows.Scan failed", zap.Error(fileErr))
			return
		}
		// get File details
		if f.File.ID > 0 {
			resStatus, err = f.File.GetFileInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Creator details
		if f.Creator.ID > 0 {
			resStatus, err = f.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier details
		if f.Modifier.ID > 0 {
			resStatus, err = f.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		voucherFiles = append(voucherFiles, f)
	}

	return
}

// Fill in the detailed information of the Execution Order body
func (eo *ExecutionOrder) FillBody() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Retrieve data from database
	bodySql := `select id,hid,rownumber,epaid,allowdelrow,
	executionvalue,executionvaluedisp,description,epadescription,ischeckerror,
	errorvalue,errorvaluedisp,isrequirefile,isonsitephoto,isissue,
	isrectify,ishandle,issueownerid,handlestarttime,handleendtime,
	status,isfromept,risklevelid,createtime,creatorid,
	confirmtime,confirmerid,modifytime,modifierid,dr,
	ts from executionorder_b
	where hid=$1 and dr=0 order by rownumber asc`
	bodyRows, err := db.Query(bodySql, eo.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.FillBody db.Query(bodySql) failed", zap.Error(err))
		return
	}
	defer bodyRows.Close()
	// Extract data row by row
	for bodyRows.Next() {
		var edr ExecutionOrderRow
		err = bodyRows.Scan(&edr.BID, &edr.HID, &edr.RowNumber, &edr.EPA.ID, &edr.AllowDelRow,
			&edr.ExecutionValue, &edr.ExecutionValueDisp, &edr.Description, &edr.EpaDescription, &edr.IsCheckError,
			&edr.ErrorValue, &edr.ErrorValueDisp, &edr.IsRequireFile, &edr.IsOnsitePhoto, &edr.IsIssue,
			&edr.IsRectify, &edr.IsHandle, &edr.IssueOwner.ID, &edr.HandleStartTime, &edr.HandleEndTime,
			&edr.Status, &edr.IsFromEPT, &edr.RiskLevel.ID, &edr.CreateDate, &edr.Creator.ID,
			&edr.ConfirmDate, &edr.Confirmer.ID, &edr.ModifyDate, &edr.Modifier.ID, &edr.Dr,
			&edr.Ts)
		if err != nil {
			zap.L().Error("ExecutionOrder.FillBody bodyRows.scan failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get Risk Level details
		if edr.RiskLevel.ID > 0 {
			resStatus, err = edr.RiskLevel.GetRLInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get IssueOwner details
		if edr.IssueOwner.ID > 0 {
			resStatus, err = edr.IssueOwner.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		// Get Creator details
		if edr.Creator.ID > 0 {
			resStatus, err = edr.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Confirmer details
		if edr.Confirmer.ID > 0 {
			resStatus, err = edr.Confirmer.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier details
		if edr.Modifier.ID > 0 {
			resStatus, err = edr.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Attachments
		edr.Files, resStatus, err = GetEORowFiles(edr.BID)
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		eo.Body = append(eo.Body, edr)
	}

	return i18n.StatusOK, nil
}

// Get Execution Order details by HID
func (eo *ExecutionOrder) GetDetailByHID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the Execution Order has already been deleted
	var rowNumber int32
	checkSql := `select count(id) as rownumber from executionorder_h where id=$1 and dr=0`
	err = db.QueryRow(checkSql, eo.HID).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.GetDetailByHID db.QueryRow(checkSql) failed", zap.Error(err))
		return
	}
	if rowNumber < 1 {
		resStatus = i18n.StatusDataDeleted
		return
	}
	// Get the Execution Order Body details
	resStatus, err = eo.FillBody()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}

	return
}

// Add Execution Order
func (eo *ExecutionOrder) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check the number of body rows, zero is not allowed
	if len(eo.Body) == 0 {
		resStatus = i18n.StatusVoucherNoBody
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Add db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Get the latest serial number
	billNo, resStatus, err := GetLatestSerialNo(tx, "EO", eo.BillDate.Format("060102"))
	if resStatus != i18n.StatusOK || err != nil {
		tx.Rollback()
		return
	}
	eo.BillNumber = billNo

	// Insert data into the executionorder_h table
	headSql := `insert into executionorder_h(billnumber,billdate,deptid,description,status,
	sourcetype,sourcebillnumber,sourcehid,sourcerownumber,sourcebid,
	starttime,endtime,csaid,executorid,eptid,
	allowaddrow,allowdelrow,creatorid) 
	values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) 
	returning id`
	err = tx.QueryRow(headSql, eo.BillNumber, eo.BillDate, eo.Department.ID, eo.Description, eo.Status,
		eo.SourceType, eo.SourceBillNumber, eo.SourceHid, eo.SourceRowNumber, eo.SourceBid,
		eo.StartTime, eo.EndTime, eo.CSA.ID, eo.Executor.ID, eo.EPT.HID,
		eo.AllowAddRow, eo.AllowDelRow, eo.Creator.ID).Scan(&eo.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Add tx.QeuryRow(headSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}

	// Prepare write the body row to the executionorder_b table
	bodySql := `insert into executionorder_b(hid,rownumber,epaid,allowdelrow,executionvalue,
		executionvaluedisp,description,epadescription,ischeckerror,errorvalue,
		errorvaluedisp,	isrequirefile,isonsitephoto,isissue,isrectify, 
		ishandle,issueownerid,handlestarttime,handleendtime,status,
		isfromept, isfinish,risklevelid,creatorid)
		values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24) 
		returning id`
	bodyStmt, err := tx.Prepare(bodySql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Add tx.Prepare(bodySql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer bodyStmt.Close()
	// Prepare write the body row attachment to the executionorder_file table
	fileSql := `insert into executionorder_file(billbid,billhid,fileid,creatorid) 
	values($1,$2,$3,$4) returning id`
	fileStmt, err := tx.Prepare(fileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Add tx.Prepare(fileSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer fileStmt.Close()
	// Write data to the database row by row
	for _, row := range eo.Body {
		var isFinish int16
		if row.IsIssue == 1 && row.IsRectify == 1 {
			isFinish = 1
		}
		// Write row data to the executionorder_b table
		err = bodyStmt.QueryRow(eo.HID, row.RowNumber, row.EPA.ID, row.AllowDelRow, row.ExecutionValue,
			row.ExecutionValueDisp, row.Description, row.EpaDescription, row.IsCheckError, row.ErrorValue,
			row.ErrorValueDisp, row.IsRequireFile, row.IsOnsitePhoto, row.IsIssue, row.IsRectify,
			row.IsHandle, row.IssueOwner.ID, row.HandleStartTime, row.HandleEndTime, row.Status,
			row.IsFromEPT, isFinish, row.RiskLevel.ID, eo.Creator.ID).Scan(&row.BID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("ExecutionOrder.Add bodyStmt.QueryRow failed", zap.Error(err))
			tx.Rollback()
			return
		}
		// Write row attachments to the executionorder_file table
		if len(row.Files) > 0 {
			for _, file := range row.Files {
				err = fileStmt.QueryRow(row.BID, eo.HID, file.File.ID, eo.Creator.ID).Scan(&file.ID)
				if err != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("ExecutionOrder.Add fileStmt.QueryRow failed", zap.Error(err))
					tx.Rollback()
					return
				}
			}
		}
	}

	// If the data comes from the Work Order,
	// then the work Order status needs to be written back as 2 (executing)
	if eo.SourceBid > 0 {
		wor := new(WorkOrderRow)
		wor.BID = eo.SourceBid
		wor.HID = eo.SourceHid
		wor.Ts = eo.SourceRowTs
		wor.EOID = eo.HID
		wor.EONumber = eo.BillNumber

		resStatus, err = wor.Execute()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
	}

	return
}

// Edit Execution Order
func (eo *ExecutionOrder) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check the number of rows in the Execution Order body
	if len(eo.Body) == 0 {
		resStatus = i18n.StatusVoucherNoBody
		return
	}
	// Check if the creator and modifier are the same person
	if eo.Creator.ID != eo.Modifier.ID {
		resStatus = i18n.StatusVoucherOnlyCreateEdit
		return
	}

	// Create a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Edit db.Begin() failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Modify Construction Order Header in the executionorder_h table
	editHeadSql := `update executionorder_h set billdate=$1,deptid=$2,description=$3,starttime=$4,endtime=$5,
	csaid=$6,executorid=$7,modifytime=current_timestamp,modifierid=$8,ts=current_timestamp  
	where id=$9 and dr=0 and status=0 and ts=$10`
	editHeadRes, err := tx.Exec(editHeadSql, &eo.BillDate, &eo.Department.ID, &eo.Description, &eo.StartTime, &eo.EndTime,
		&eo.CSA.ID, &eo.Executor.ID, &eo.Modifier.ID,
		&eo.HID, &eo.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Edit tx.Exec(editHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	headUpdateNumber, err := editHeadRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Edit EditHeadRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if headUpdateNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}

	// Prepare modify the Execution Order Row in the executionorder_b table
	updateRowSql := `update executionorder_b set epaid=$1, allowdelrow=$2,executionvalue=$3,executionvaluedisp=$4,description=$5,
	epadescription=$6,	ischeckerror=$7,errorvalue=$8,errorvaluedisp=$9,isrequirefile=$10,
	isonsitephoto=$11,isissue=$12,isrectify=$13,ishandle=$14,issueownerid=$15,
	handlestarttime=$16,handleendtime=$17, status=$18,isfromept=$19,risklevelid=$20,
	modifytime=current_timestamp,modifierid=$21,ts=current_timestamp,dr=$22,isFinish=$23
	where id=$24 and ts=$25 and status=0 and dr=0`
	updateRowStmt, err := tx.Prepare(updateRowSql)
	if err != nil {
		zap.L().Error("ExecutionOrder.Edit tx.Prepare(updateRowSql) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer updateRowStmt.Close()
	// Prepare Add the Execution Order row in the execution_b table
	addRowSql := `insert into executionorder_b(hid,rownumber,epaid,allowdelrow,executionvalue,
		executionvaluedisp,description,epadescription,ischeckerror,errorvalue,
		errorvaluedisp,	isrequirefile,isonsitephoto,isissue,isrectify,
		ishandle,status,issueownerid,handlestarttime,handleendtime,
		isfromept,isfinish,risklevelid,creatorid) 
		values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24) 
	returning id`
	addRowStmt, err := tx.Prepare(addRowSql)
	if err != nil {
		zap.L().Error("ExecutionOrder.Edit tx.Prepare(addRowSql) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer addRowStmt.Close()
	// Prepare modify the Execution Order Row attachments in the executionorder_file table
	updateFileSql := `update executionorder_file set modifytime=current_timestamp,modifierid=$1,dr=$2,ts=current_timestamp
	where id=$3 and dr=0 and ts=$4`
	updateFileStmt, err := tx.Prepare(updateFileSql)
	if err != nil {
		zap.L().Error("ExecutionOrder.Edit tx.Prepare(updateFileSql) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer updateFileStmt.Close()
	// Prepare Add the Execution Order Row attachments in the executionorder_file table
	addFileSql := `insert into executionorder_file(billbid,billhid,fileid,creatorid) 
	values($1,$2,$3,$4) returning id`
	addFileStmt, err := tx.Prepare(addFileSql)
	if err != nil {
		zap.L().Error("ExecutionOrder.Edit tx.Prepare(addFileSql) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer addFileStmt.Close()
	// Write the Execution Order Row Data into executionorder_b row by row
	for _, row := range eo.Body {
		// Check the Execution Order Row status
		if row.Status != 0 {
			resStatus = i18n.StatusVoucherNoFree
			tx.Rollback()
			return
		}
		var isFinish int16
		if row.IsIssue == 1 && row.IsRectify == 1 {
			isFinish = 1
		}

		if row.BID == 0 { // If the HID value is 0, it menas it is a newly row
			addRowErr := addRowStmt.QueryRow(eo.HID, row.RowNumber, row.EPA.ID, row.AllowDelRow, row.ExecutionValue,
				row.ExecutionValueDisp, row.Description, row.EpaDescription, row.IsCheckError, row.ErrorValue,
				row.ErrorValueDisp, row.IsRequireFile, row.IsOnsitePhoto, row.IsIssue, row.IsRectify,
				row.IsHandle, row.Status, row.IssueOwner.ID, row.HandleStartTime, row.HandleEndTime,
				row.IsFromEPT, isFinish, row.RiskLevel.ID, eo.Modifier.ID).Scan(&row.BID)
			if addRowErr != nil {
				zap.L().Error("ExecutionOrder.Edit addRowStmt.QueryRow() failed", zap.Error(addRowErr))
				resStatus = i18n.StatusInternalError
				tx.Rollback()
				return resStatus, addRowErr
			}

			// Add the row attachments records
			if len(row.Files) > 0 {
				for _, file := range row.Files {
					addFileErr := addFileStmt.QueryRow(row.BID, eo.HID, file.File.ID, eo.Creator.ID).Scan(&file.ID)
					if addFileErr != nil {
						resStatus = i18n.StatusInternalError
						zap.L().Error("ExecutionOrder.Edit new row addFileStmt.QueryRow failed", zap.Error(err))
						tx.Rollback()
						return resStatus, addFileErr
					}
				}
			}

		} else { // If the HID value is not 0, it means it is a row that needs to be modified
			updateRowRes, updateRowErr := updateRowStmt.Exec(row.EPA.ID, row.AllowDelRow, row.ExecutionValue, row.ExecutionValueDisp, row.Description,
				row.EpaDescription, row.IsCheckError, row.ErrorValue, row.ErrorValueDisp, row.IsRequireFile,
				row.IsOnsitePhoto, row.IsIssue, row.IsRectify, row.IsHandle, row.IssueOwner.ID,
				row.HandleStartTime, row.HandleEndTime, row.Status, row.IsFromEPT, row.RiskLevel.ID,
				eo.Modifier.ID, row.Dr, isFinish,
				row.BID, row.Ts)
			if updateRowErr != nil {
				zap.L().Error("ExecutionOrder.Edit updateRowStmt.Exec() failed", zap.Error(updateRowErr))
				resStatus = i18n.StatusInternalError
				tx.Rollback()
				return
			}
			// Chcek the number of rows affected by SQL statement
			updateRowNumber, errUpdateEffect := updateRowRes.RowsAffected()
			if errUpdateEffect != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("ExecutionOrder.Edit updateRowRes.RowsAffected failed", zap.Error(errUpdateEffect))
				tx.Rollback()
				return resStatus, errUpdateEffect
			}
			if updateRowNumber < 1 {
				resStatus = i18n.StatusOtherEdit
				tx.Rollback()
				return
			}

			// Handle the row attachment
			if len(row.Files) > 0 {
				for _, file := range row.Files {
					if file.ID == 0 { // If the file.ID value is 0, it means it is a newly file.
						addFileErr := addFileStmt.QueryRow(row.BID, eo.HID, file.File.ID, eo.Modifier.ID).Scan(&file.ID)
						if addFileErr != nil {
							resStatus = i18n.StatusInternalError
							zap.L().Error("ExecutionOrder.Edit old row addFileStmt.QueryRow failed", zap.Error(addFileErr))
							tx.Rollback()
							return resStatus, addFileErr
						}
					} else { // If the file.ID value is not o, it means it is a file that needs to be modified
						updateFileRes, updateFileErr := updateFileStmt.Exec(eo.Modifier.ID, file.Dr, file.ID, file.Ts)
						if updateFileErr != nil {
							resStatus = i18n.StatusInternalError
							zap.L().Error("ExecutionOrder.Edit old row updateFileRes.Exec() failed", zap.Error(updateFileErr))
							tx.Rollback()
							return resStatus, updateFileErr
						}
						updateFileNumber, updateFileEffectErr := updateFileRes.RowsAffected()
						if updateFileEffectErr != nil {
							resStatus = i18n.StatusInternalError
							zap.L().Error("ExecutionOrder.Edit old row updateFileRes.RowsAffected() failed", zap.Error(updateFileEffectErr))
							tx.Rollback()
							return resStatus, updateFileEffectErr
						}

						if updateFileNumber < 1 {
							resStatus = i18n.StatusOtherEdit
							tx.Rollback()
							return
						}
					}
				}
			}
		}
	}
	return
}

// Delete Execution Order
func (eo *ExecutionOrder) Delete(modifyUserId int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get the Execution Order Header details
	resStatus, err = eo.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check the Execution order status
	if eo.Status != 0 { // Status must be 0
		resStatus = i18n.StatusVoucherNoFree
		return
	}
	// Check if the creator and modifier are the same person
	if eo.Creator.ID != modifyUserId {
		resStatus = i18n.StatusVoucherOnlyCreateEdit
		return
	}

	// Create a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Delete db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Modify the Execution Order Header delete flag to 1 in the executionorder_b table
	delHeadSql := `update executionorder_h set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	delHeadRes, err := tx.Exec(delHeadSql, modifyUserId, eo.HID, eo.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Delete tx.Exec(delHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	delHeadNumber, err := delHeadRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Delete delHeadRes.RowsAffected() failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if delHeadNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}

	// Prepare modify the Execution Order Rows delete flag to 1 in the executionorder_b table
	delRowSql := `update executionorder_b set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	delRowStmt, err := tx.Prepare(delRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Delete tx.Prepare(delRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delRowStmt.Close()
	// Prepare modify the Execution Order Row Attachments delete flag to 1 in the executionorder_file table
	delFileSql := `update executionorder_file set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and billbid=$3 and ts=$4`
	delFileStmt, err := tx.Prepare(delFileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Delete tx.Prepare(delFileSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delFileStmt.Close()
	// Write the modified content to the database row by row
	for _, row := range eo.Body {
		// Check the Execution Order Row status
		if row.Status != 0 {
			resStatus = i18n.StatusVoucherNoFree
			tx.Rollback()
			return
		}
		delRowRes, errDelRow := delRowStmt.Exec(modifyUserId, row.BID, row.Ts)
		if errDelRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("ExecutionOrder.Delete delRowStmt.Exec() failed", zap.Error(errDelRow))
			tx.Rollback()
			return resStatus, errDelRow
		}
		// Check the number of rows affected by SQL statement
		delRowNumber, errDelRow := delRowRes.RowsAffected()
		if errDelRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("ExecutionOrder.Delete delRowRes.RowsAffected() failed", zap.Error(errDelRow))
			tx.Rollback()
			return resStatus, errDelRow
		}
		if delRowNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}

		// Handle the Execution Order attaments
		if len(row.Files) > 0 {
			for _, file := range row.Files {
				delFileRes, delFileErr := delFileStmt.Exec(modifyUserId, file.ID, row.BID, file.Ts)
				if delFileErr != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("ExecutionOrder.Delete delFileStmt.Exec() failed", zap.Error(delFileErr))
					tx.Rollback()
					return resStatus, delFileErr
				}

				delFileNumber, delFileErr := delFileRes.RowsAffected()
				if delFileErr != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("ExecutionOrder.Delete delFileRes.RowsAffected() failed", zap.Error(delFileErr))
					tx.Rollback()
					return resStatus, delFileErr
				}
				if delFileNumber < 1 {
					resStatus = i18n.StatusOtherEdit
					tx.Rollback()
					return
				}
			}
		}
	}
	// If the data comes from the Work Order,
	// then the work Order status needs to be written back as 1 (confirmed)
	if eo.SourceBid > 0 {
		wor := new(WorkOrderRow)
		wor.BID = eo.SourceBid
		wor.HID = eo.SourceHid
		wor.Ts = eo.SourceRowTs
		wor.EOID = eo.HID
		wor.EONumber = eo.BillNumber

		resStatus, err = wor.CancelExecute()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
	}

	return
}

// Confirm Execution Order
func (eo *ExecutionOrder) Confirm(confirmUserID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get the Execution Order details
	resStatus, err = eo.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check the Execution Order status
	if eo.Status != 0 { // Must be 0
		resStatus = i18n.StatusVoucherNoFree
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Confirm db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Write the confirmation information to the executionorder_h table
	confirmHeadSql := `update executionorder_h set status=1,confirmtime=current_timestamp,confirmerid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and status=0 and ts=$3`
	headRes, err := tx.Exec(confirmHeadSql, confirmUserID, eo.HID, eo.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Confirm tx.Exec(confirmHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	confirmHeadNumber, err := headRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Confirm headRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if confirmHeadNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}

	// Prepare write the confirmation information to the executionorder_b table
	confirmRowSql := `update executionorder_b set status=1,confirmtime=current_timestamp,confirmerid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and status=0 and ts=$3`
	rowStmt, err := tx.Prepare(confirmRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Confirm tx.Prepare(confirmRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer rowStmt.Close()
	// Write the confirmation information row by row
	for _, row := range eo.Body {
		// Check the Execution Order rows status
		if row.Status != 0 {
			resStatus = i18n.StatusVoucherNoFree
			tx.Rollback()
			return
		}
		confirmRowRes, errConfirmRow := rowStmt.Exec(confirmUserID, row.BID, row.Ts)
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("ExecutionOrder.Confirm rowStmt.Exec failed", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}

		confirmRowNumber, errConfirmRow := confirmRowRes.RowsAffected()
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("ExecutionOrder.Confirm confirmRowRes.RowsAffected failed", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}
		if confirmRowNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}
	}

	// If the data comes from the Work Order,
	// then the work Order status needs to be written back as 3 (completed)
	if eo.SourceBid > 0 {
		wor := new(WorkOrderRow)
		wor.BID = eo.SourceBid
		wor.HID = eo.SourceHid
		wor.EOID = eo.HID
		wor.EONumber = eo.BillNumber

		resStatus, err = wor.Complete()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
	}

	return
}

// UnConfirm Execution Order
func (eo *ExecutionOrder) UnConfirm(confirmUserID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get the Execution Order details
	resStatus, err = eo.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check the Execution Order status
	if eo.Status != 1 { // Must be 1
		resStatus = i18n.StatusVoucherNoConfirm
		return
	}
	// Check the Execution Order Confirmer and Unconfirmer are the same person
	if eo.Confirmer.ID != confirmUserID {
		resStatus = i18n.StatusVoucherCancelConfirmSelf
		return
	}
	// Check the Execution Order Rows status
	var noConfirmRowNumber int32
	for _, row := range eo.Body {
		if row.Status > 1 {
			noConfirmRowNumber++
		}
	}
	if noConfirmRowNumber > 0 {
		resStatus = i18n.StatusEOBodyNoConfirm
		return
	}

	// Create a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.UnConfirm db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Write the un-confirmation information to the executionorder_h table
	confirmHeadSql := `update executionorder_h set status=0,confirmerid=0,confirmtime=to_timestamp(0),ts=current_timestamp 
	where id=$1 and dr=0 and status=1 and ts=$2`
	headRes, err := tx.Exec(confirmHeadSql, eo.HID, eo.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.ExecutionOrder.UnConfirm tx.Exec(confirmHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	confirmHeadNumber, err := headRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.UnConfirm headRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if confirmHeadNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Prepare write the un-confirmation information to the executionorder_b table
	confirmRowSql := `update executionorder_b set status=0,confirmerid=0,confirmtime=to_timestamp(0),ts=current_timestamp 
	where id=$1 and dr=0 and status=1 and ts=$2`
	rowStmt, err := tx.Prepare(confirmRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.UnConfirm tx.Prepare(confirmRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer rowStmt.Close()
	// Write the un-confirmation information to the database row by row
	for _, row := range eo.Body {
		// Check the Execution Rows status
		if row.Status != 1 {
			resStatus = i18n.StatusVoucherNoConfirm
			tx.Rollback()
			return
		}
		confirmRowRes, errConfirmRow := rowStmt.Exec(row.BID, row.Ts)
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("ExecutionOrder.UnConfirm rowStmt.Exec failed", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}

		confirmRowNumber, errConfirmRow := confirmRowRes.RowsAffected()
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("ExecutionOrder.UnConfirm confirmRowRes.RowsAffected failed", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}
		if confirmRowNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}
	}

	// If the data comes from the Work Order,
	// then the work Order status needs to be written back as 2 (executing)
	if eo.SourceBid > 0 {
		wor := new(WorkOrderRow)
		wor.BID = eo.SourceBid
		wor.HID = eo.SourceHid
		wor.EOID = eo.HID
		wor.EONumber = eo.BillNumber

		resStatus, err = wor.CancelComplete()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
	}
	return
}

// Update the status after handle the Execution Order Row
func (edr *ExecutionOrderRow) Dispose() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Write the Issue Resolution Form information to the executionorder_b table
	rowSql := `update executionorder_b set status=2,ts=current_timestamp,isfinish=$1,irfid=$2,irfnumber=$3  
	where id=$4 and hid=$5 and ts=$6 and dr=0 and status=1 and isfinish=0`
	rowUpdateRes, err := db.Exec(rowSql, edr.IsFinish, edr.IRFID, edr.IRFNumber, edr.BID, edr.HID, edr.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrderRow.Dispose  db.Exec(rowSql) failed", zap.Error(err))
		return
	}
	// Check the number of rows affected by SQL statement
	rowUpdateNumber, err := rowUpdateRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrderRow.Dispose rowUpdateRes.RowsAffected failed", zap.Error(err))
		return
	}
	if rowUpdateNumber < 1 {
		resStatus = i18n.StatusWOOtherEdit
		zap.L().Info("ExecutionOrderRow.Dispose row OtherEdit")
		return
	}
	return
}

// Update the status after cancel handle the Execution Order Row
func (edr *ExecutionOrderRow) CancelDispose() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Clear the Issue Resolution Form information in the executionorder_b table
	rowSql := `update executionorder_b set status=1,ts=current_timestamp,isfinish=$1,irfid=$2,irfnumber=$3  
	where id=$4 and hid=$5 and dr=0 and status=2 and isfinish=1`
	rowUpdateRes, err := db.Exec(rowSql, edr.IsFinish, edr.IRFID, edr.IRFNumber, edr.BID, edr.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrderRow.CancelDispose  db.Exec(rowSql) failed", zap.Error(err))
		return
	}
	// Check the number of rows affected by SQL statement
	rowUpdateNumber, err := rowUpdateRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrderRow.CancelDispose rowUpdateRes.RowsAffected failed", zap.Error(err))
		return
	}
	if rowUpdateNumber < 1 {
		resStatus = i18n.StatusWOOtherEdit
		zap.L().Info("ExecutionOrderRow.CancelDispose row OtherEdit")
		return
	}
	return
}

// Complete the Execution Order Row (Update the status after confirm the Issue Resolutin Form)
func (edr *ExecutionOrderRow) Complete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Modify the status in the executionorder_b table
	rowSql := `update executionorder_b set status=3,ts=current_timestamp   
	where id=$1 and hid=$2 and dr=0 and status=2 and isfinish=1`
	rowUpdateRes, err := db.Exec(rowSql, edr.BID, edr.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrderRow.Complete  db.Exec(rowSql) failed", zap.Error(err))
		return
	}
	// Check the number of rows affected by SQL statement
	rowUpdateNumber, err := rowUpdateRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrderRow.Complete rowUpdateRes.RowsAffected failed", zap.Error(err))
		return
	}
	if rowUpdateNumber < 1 {
		resStatus = i18n.StatusWOOtherEdit
		zap.L().Info("ExecutionOrderRow.Complete row OtherEdit")
		return
	}
	return
}

// Cancel complete the Execution Order Row (update the status after unconfirm the Issue Resolution Form)
func (edr *ExecutionOrderRow) CancelComplete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Modify the status in the executionorder_b table
	rowSql := `update executionorder_b set status=2,ts=current_timestamp   
	where id=$1 and hid=$2 and dr=0 and status=3 and isfinish=1`
	rowUpdateRes, err := db.Exec(rowSql, edr.BID, edr.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrderRow.CancelComplete  db.Exec(rowSql) failed", zap.Error(err))
		return
	}
	// Check the number of rows affected by SQL statement
	rowUpdateNumber, err := rowUpdateRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrderRow.CancelComplete rowUpdateRes.RowsAffected failed", zap.Error(err))
		return
	}
	if rowUpdateNumber < 1 {
		resStatus = i18n.StatusWOOtherEdit
		zap.L().Info("ExecutionOrderRow.CancelComplete row OtherEdit")
		return
	}
	return
}

// Add Execution Order Comment
func (eoc *ExecutionOrderComment) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Insert comment content into the executionorder_comment table
	sqlStr := `insert into executionorder_comment(bid,hid,billnumber,rownumber,sendtoid,
	content,creatorid)
	values($1,$2,$3,$4,$5,$6,$7) 
	returning id`
	err = db.QueryRow(sqlStr, eoc.BID, eoc.HID, eoc.BillNumber, eoc.RowNUmber, eoc.SendTo.ID,
		eoc.Content, eoc.Creator.ID).Scan(&eoc.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrderComment.Add db.QueryRow(sqlStr) failed", zap.Error(err))
		return
	}
	return
}

// Add Execution Order Review Record
func (edrr *ExecutionOrderReview) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Insert review content into the executionorder_comment table
	sqlStr := `insert into executionorder_review(hid,billnumber,starttime,endtime,consumeseconds,creatorid) 
	values($1,$2,$3,$4,$5,$6) returning id`
	err = db.QueryRow(sqlStr, edrr.HID, edrr.BillNumber, edrr.StartTime, edrr.EndTime, edrr.ConsumeSeconds, edrr.Creator.ID).Scan(&edrr.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrderReview.Add db.QueryRow(sqlStr) failed", zap.Error(err))
		return
	}
	return
}

// Get Execution Order Components
func (cs *EOCommentsParams) Get() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	cs.Comments = make([]ExecutionOrderComment, 0)
	// Get the comments list from executionorder_comment table
	sqlStr := `select id,bid,hid,billnumber,rownumber,
	sendtoid,isread,readtime,content,sendtime,
	createtime,creatorid,dr,ts 
	from executionorder_comment 
	where dr=0 and hid=$1`
	rows, err := db.Query(sqlStr, cs.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EOCommentsParams.Get db.Query failed", zap.Error(err))
		return
	}
	defer rows.Close()

	// Extract comments data row by row
	for rows.Next() {
		var eoc ExecutionOrderComment
		err = rows.Scan(&eoc.ID, &eoc.BID, &eoc.HID, &eoc.BillNumber, &eoc.RowNUmber,
			&eoc.SendTo.ID, &eoc.IsRead, &eoc.ReadTime, &eoc.Content, &eoc.SendTime,
			&eoc.CreateDate, &eoc.Creator.ID, &eoc.Dr, &eoc.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("EOCommentsParams.Get rows.Next failed", zap.Error(err))
			return
		}
		// Get SendTo Person details
		if eoc.SendTo.ID > 0 {
			resStatus, err = eoc.SendTo.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Creator details
		if eoc.Creator.ID > 0 {
			resStatus, err = eoc.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		cs.Comments = append(cs.Comments, eoc)
	}

	return
}

// Get Execution Order Review Records list
func (rs *EOReviewsParams) Get() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	rs.Reviews = make([]ExecutionOrderReview, 0)
	// Get the review records from the executionorder_review table
	sqlStr := `select id,hid,billnumber,starttime,endtime,
	consumeseconds,createtime,creatorid,dr,ts 
	from executionorder_review 
	where dr=0 and hid=$1`
	rows, err := db.Query(sqlStr, rs.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EOReviewsParams.Get db.Query failed", zap.Error(err))
		return
	}
	defer rows.Close()
	// Extract records row by row
	for rows.Next() {
		var er ExecutionOrderReview
		err = rows.Scan(&er.ID, &er.HID, &er.BillNumber, &er.StartTime, &er.EndTime,
			&er.ConsumeSeconds, &er.CreateDate, &er.Creator.ID, &er.Dr, &er.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("EOReviewsParams.Get db.Query failed", zap.Error(err))
			return
		}
		if er.Creator.ID > 0 {
			resStatus, err = er.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		rs.Reviews = append(rs.Reviews, er)
	}
	return
}
