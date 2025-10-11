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
	SourceType       string              `db:"sourcetype" json:"sourceType"` //Source Type: UA:User Add WO: Work Order
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
	HandleStartTime    string           `db:"handlestarttime" json:"handleStartTime"`
	HandleEndTime      string           `db:"handleendtime" json:"handleEndTime"`
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
	Department         SimpDept         `db:"h.dept_id" json:"department"`
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
	SendTo     Person    `db:"sendto_id" json:"sendTo"`
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

// Execution Order  Comments Params
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

// GetReferEDs 获取参照执行单列表
func GetReferEDs(queryString string) (reds []ReferExecutionOrder, resStatus i18n.ResKey, err error) {
	reds = make([]ReferExecutionOrder, 0)
	var build strings.Builder
	//拼接检查sql
	build.WriteString(`select count(b.id) as rownumber
	from executedoc_b as b
	left join executedoc_h as h on b.hid = h.id
	left join exectiveitem as epa on b.epaid = epa.id
	left join sysuser as issueowner on b.issueownerid = issueowner.id
	left join sysuser as epuser on h.executorid = epuser.id
	left join sysuser as creator on h.creatorid = creator.id
	left join department as dept on h.dept_id = dept.id
	where (b.ishandle=1 and b.dr = 0 and b.isfinish=0 and b.status=1)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	checkSql := build.String()
	//检查
	var rowNumber int32
	err = db.QueryRow(checkSql).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEDRefer db.QueryRow(checkSql) failed", zap.Error(err))
		return
	}
	if rowNumber == 0 { //如果查询数据量为0
		resStatus = i18n.StatusResNoData
		return
	}
	if rowNumber > setting.Conf.PqConfig.MaxRecord { //查询数据量大于最大记录数
		resStatus = i18n.StatusOverRecord
		return
	}
	build.Reset() //清空build
	//拼接正式sql
	build.WriteString(`select b.id,b.hid,b.rownumber,b.epaid,b.executionvalue,
	b.executionvaluedisp,b.description,b.ishandle,b.issueownerid,b.handlestarttime,
	b.handleendtime,b.status,b.risklevelid, b.isfinish,b.dr,
	b.ts,h.billnumber,h.billdate,h.dept_id,h.csaid,
	h.executorid 
	from executedoc_b as b
	left join executedoc_h as h on b.hid = h.id
	left join exectiveitem as epa on b.epaid = epa.id
	left join sysuser as issueowner on b.issueownerid = issueowner.id
	left join sysuser as epuser on h.executorid = epuser.id
	left join sysuser as creator on h.creatorid = creator.id
	left join department as dept on h.dept_id = dept.id
	where (b.ishandle=1 and b.dr = 0 and b.isfinish=0 and (b.status=1 or b.status=2))`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	refSql := build.String()
	//获取执行单参照列表
	edRef, err := db.Query(refSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEDRefer db.Query failed", zap.Error(err))
		return
	}
	defer edRef.Close()

	//提取数据
	for edRef.Next() {
		var red ReferExecutionOrder
		err = edRef.Scan(&red.BID, &red.HID, &red.RowNumber, &red.EPA.ID, &red.ExecutionValue,
			&red.ExecutionValueDisp, &red.Description, &red.IsHandle, &red.IssueOwner.ID, &red.HandleStartTime,
			&red.HandleEndTime, &red.Status, &red.RiskLevel.ID, &red.IsFinish, &red.Dr,
			&red.Ts, &red.BillNumber, &red.BillDate, &red.Department.ID, &red.CSA.ID,
			&red.Executor.ID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetEDRefer edRef.Next() edRef.Scan() failed", zap.Error(err))
			return
		}
		//填充执行项目信息
		if red.EPA.ID > 0 {
			resStatus, err = red.EPA.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充部门信息
		if red.Department.ID > 0 {
			resStatus, err = red.Department.GetSimpDeptInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充风险等级信息
		if red.RiskLevel.ID > 0 {
			resStatus, err = red.RiskLevel.GetRLInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充现场档案信息
		if red.CSA.ID > 0 {
			resStatus, err = red.CSA.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充执行人信息
		if red.Executor.ID > 0 {
			resStatus, err = red.Executor.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充处理人信息
		if red.IssueOwner.ID > 0 {
			resStatus, err = red.IssueOwner.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//获取附件
		red.EOFiles, resStatus, err = GetEDRFiles(red.BID)
		if resStatus != i18n.StatusOK || err != nil {
			return
		}

		reds = append(reds, red)
	}

	return reds, i18n.StatusOK, nil
}

// GetEDList 获取执行单列表
func GetEDList(queryString string) (eos []ExecutionOrder, resStatus i18n.ResKey, err error) {
	var build strings.Builder
	//拼接检查sql
	build.WriteString(`select count(h.id) as rownumber
	from executedoc_h as h
	left join department on h.dept_id = department.id
	left join sysuser as creator on h.creatorid = creator.id
	left join sysuser as modifier on h.modifierid = modifier.id
	where (h.dr = 0) `)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	checkSql := build.String()
	//检查
	var rowNumber int32
	err = db.QueryRow(checkSql).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEDList db.QueryRow(checkSql) failed", zap.Error(err))
		return
	}
	if rowNumber == 0 { //如果查询数据量为0
		resStatus = i18n.StatusResNoData
		return
	}
	if rowNumber > setting.Conf.PqConfig.MaxRecord { //查询数据量大于最大记录数
		resStatus = i18n.StatusOverRecord
		return
	}
	build.Reset() //清空build

	//拼接正式sql
	build.WriteString(`select h.id,h.billnumber,h.billdate,h.dept_id,h.description,
	h.status,h.sourcetype,h.sourcebillnumber,h.sourcehid,h.sourcerownumber,
	h.sourcebid,h.starttime,h.endtime,h.csaid,h.executorid,
	h.eptid,h.allowaddrow,h.allowdelrow,h.createtime,h.creatorid,
	h.confirmtime,h.confirmerid,h.modifytime,h.modifierid,h.dr,
	h.ts 
	from executedoc_h as h
	left join department on h.dept_id = department.id
	left join sysuser as creator on h.creatorid = creator.id
	left join sysuser as modifier on h.modifierid = modifier.id
	where (h.dr = 0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	headSql := build.String()
	//获取指令单列表
	headRows, err := db.Query(headSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEDList db.Query failed", zap.Error(err))
		return
	}
	defer headRows.Close()
	//提取数据
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
			zap.L().Error("GetEDList headRows.Next failed", zap.Error(err))
			return
		}
		//填充信息
		resStatus, err = eo.FillHead()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		eos = append(eos, eo)
	}
	resStatus = i18n.StatusOK
	return
}

// GetEDReviewList 获取审阅执行单列表
func GetEDReviewList(queryString string, useID int32) (eos []ExecutionOrder, resStatus i18n.ResKey, err error) {
	eos = make([]ExecutionOrder, 0)
	var build strings.Builder
	//拼接检查sql
	build.WriteString(`select count(h.id) as rownumber
	from executedoc_h as h
	left join department on h.dept_id = department.id
	left join sysuser as creator on h.creatorid = creator.id
	left join sysuser as modifier on h.modifierid = modifier.id
	where (h.dr = 0) `)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	checkSql := build.String()
	//检查
	var rowNumber int32
	err = db.QueryRow(checkSql).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEDList db.QueryRow(checkSql) failed", zap.Error(err))
		return
	}
	if rowNumber == 0 { //如果查询数据量为0
		resStatus = i18n.StatusResNoData
		return
	}
	if rowNumber > setting.Conf.PqConfig.MaxRecord { //查询数据量大于最大记录数
		resStatus = i18n.StatusOverRecord
		return
	}
	build.Reset() //清空build

	//拼接正式sql
	build.WriteString(`select h.id,h.billnumber,h.billdate,h.dept_id,h.description,
	h.status,h.sourcetype,h.sourcebillnumber,h.sourcehid,h.sourcerownumber,
	h.sourcebid,h.starttime,h.endtime,h.csaid,h.executorid,
	h.eptid,h.allowaddrow,h.allowdelrow,h.createtime,h.creatorid,
	h.confirmtime,h.confirmerid,h.modifytime,h.modifierid,h.dr,
	h.ts,
	(select count(b.id) as errnumber from executedoc_b as b where b.hid = h.id and b.dr=0 and b.isissue=1),
	(select count(r.id) as reviewednumber from executedoc_review as r where r.hid = h.id and r.dr=0 and r.creatorid=$1),
	(select coalesce( sum(r.consumeseconds),0) as reviewedseconds  from executedoc_review as r where r.hid = h.id and r.dr=0 and r.creatorid=$1)
	from executedoc_h as h
	left join department on h.dept_id = department.id
	left join sysuser as creator on h.creatorid = creator.id
	left join sysuser as modifier on h.modifierid = modifier.id
	where (h.dr = 0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	headSql := build.String()
	//获取指令单列表
	headRows, err := db.Query(headSql, useID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEDList db.Query failed", zap.Error(err))
		return
	}
	defer headRows.Close()
	//提取数据
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
			zap.L().Error("GetEDList headRows.Next failed", zap.Error(err))
			return
		}
		//填充信息
		resStatus, err = eo.FillHead()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		eos = append(eos, eo)
	}

	resStatus = i18n.StatusOK
	return
}

// GetEDRReviesListPaging 获取审阅执行单分页列表
func GetEDRReviewListPaging(con PagingQueryParams, userID int32) (edsp EOListPaging, resStatus i18n.ResKey, err error) {
	edsp.EOs = make([]ExecutionOrder, 0)
	var build strings.Builder
	//拼接检查sql
	build.WriteString(`select count(h.id) as rownumber
	from executedoc_h as h
	left join department on h.dept_id = department.id
	left join sysuser as creator on h.creatorid = creator.id
	left join sysuser as modifier on h.modifierid = modifier.id
	where (h.dr = 0) `)
	if con.QueryString != "" {
		build.WriteString(" and (")
		build.WriteString(con.QueryString)
		build.WriteString(")")
	}
	checkSql := build.String()

	//检查
	err = db.QueryRow(checkSql).Scan(&edsp.Count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEDList db.QueryRow(checkSql) failed", zap.Error(err))
		return
	}
	if edsp.Count == 0 { //如果查询数据量为0
		resStatus = i18n.StatusResNoData
		return
	}
	if edsp.Count > setting.Conf.PqConfig.MaxRecord { //查询数据量大于最大记录数
		resStatus = i18n.StatusOverRecord
		return
	}
	//重新计算分页
	if con.PerPage > edsp.Count {
		con.Page = 0
	} else {
		var totalPage = int32(math.Ceil(float64(edsp.Count) / float64(con.PerPage)))
		if (con.Page + 1) > totalPage {
			con.Page = totalPage - 1
		}
	}
	build.Reset() //清空build

	//拼接正式sql
	build.WriteString(`select h.id,h.billnumber,h.billdate,h.dept_id,h.description,
	h.status,h.sourcetype,h.sourcebillnumber,h.sourcehid,h.sourcerownumber,
	h.sourcebid,h.starttime,h.endtime,h.csaid,h.executorid,
	h.eptid,h.allowaddrow,h.allowdelrow,h.createtime,h.creatorid,
	h.confirmtime,h.confirmerid,h.modifytime,h.modifierid,h.dr,
	h.ts,
	(select count(b.id) as errnumber from executedoc_b as b where b.hid = h.id and b.dr=0 and b.isissue=1),
	(select count(r.id) as reviewednumber from executedoc_review as r where r.hid = h.id and r.dr=0 and r.creatorid=$1),
	(select coalesce( sum(r.consumeseconds),0) as reviewedseconds  from executedoc_review as r where r.hid = h.id and r.dr=0 and r.creatorid=$1)
	from executedoc_h as h
	left join department on h.dept_id = department.id
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

	//获取指令单列表
	headRows, err := db.Query(headSql, userID, con.PerPage, con.Page*con.PerPage)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEDList db.Query failed", zap.Error(err))
		return
	}
	defer headRows.Close()
	//提取数据
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
			zap.L().Error("GetEDList headRows.Next failed", zap.Error(err))
			return
		}
		//填充信息
		resStatus, err = eo.FillHead()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		edsp.EOs = append(edsp.EOs, eo)
	}
	edsp.Page = con.Page
	edsp.PerPage = con.PerPage

	resStatus = i18n.StatusOK
	return
}

// ExecutionOrder.FillHeadStruct 填充表头结构体
func (eo *ExecutionOrder) FillHead() (resStatus i18n.ResKey, err error) {
	//填充部门信息
	if eo.Department.ID > 0 {
		resStatus, err = eo.Department.GetSimpDeptInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//填充现场档案
	if eo.CSA.ID > 0 {
		resStatus, err = eo.CSA.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//填充执行人
	if eo.Executor.ID > 0 {
		resStatus, err = eo.Executor.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//填充执行模板
	if eo.EPT.HID > 0 {
		resStatus, err = eo.EPT.GetEPTHeaderByHid()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//填充创建人
	if eo.Creator.ID > 0 {
		resStatus, err = eo.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//填充确认人
	if eo.Confirmer.ID > 0 {
		resStatus, err = eo.Confirmer.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//填充修改人
	if eo.Modifier.ID > 0 {
		resStatus, err = eo.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}

	return i18n.StatusOK, nil
}

// GetEDRFiles 获取执行单表体附件
func GetEDRFiles(bid int32) (voucherFiles []VoucherFile, resStatus i18n.ResKey, err error) {
	voucherFiles = make([]VoucherFile, 0) //解决返回文件为空问题
	attachSql := `select id,billbid,billhid,fileid,createtime,
	creatorid,modifytime,modifierid,dr,ts 
	from executedoc_file where billbid=$1 and dr=0`
	//填充附件
	fileRows, err := db.Query(attachSql, bid)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEDRFiles db.query(attachsql) failed", zap.Error(err))
		return
	}
	defer fileRows.Close()

	for fileRows.Next() {
		var f VoucherFile
		fileErr := fileRows.Scan(&f.ID, &f.BillBID, &f.BillHID, &f.File.ID, &f.CreateDate,
			&f.Creator.ID, &f.ModifyDate, &f.Modifier.ID, &f.Dr, &f.Ts)
		if fileErr != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetEDRFiles fileRows.Scan failed", zap.Error(fileErr))
			return
		}
		//填充文件信息
		if f.File.ID > 0 {
			resStatus, err = f.File.GetFileInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充创建人
		if f.Creator.ID > 0 {
			resStatus, err = f.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充更新人
		if f.Modifier.ID > 0 {
			resStatus, err = f.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		voucherFiles = append(voucherFiles, f)
	}
	resStatus = i18n.StatusOK
	return
}

// ExecutionOrder.FillBody 填充表体
func (eo *ExecutionOrder) FillBody() (resStatus i18n.ResKey, err error) {
	bodySql := `select id,hid,rownumber,epaid,allowdelrow,
	executionvalue,executionvaluedisp,description,epadescription,ischeckerror,
	errorvalue,errorvaluedisp,isrequirefile,isonsitephoto,isissue,
	isrectify,ishandle,issueownerid,handlestarttime,handleendtime,
	status,isfromept,risklevelid,createtime,creatorid,
	confirmtime,confirmerid,modifytime,modifierid,dr,
	ts from executedoc_b
	where hid=$1 and dr=0 order by rownumber asc`
	bodyRows, err := db.Query(bodySql, eo.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.FillBody db.Query(bodySql) failed", zap.Error(err))
		return
	}
	defer bodyRows.Close()

	var bodyRowNumber int32

	for bodyRows.Next() {
		bodyRowNumber++
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
		//填充风险等级
		if edr.RiskLevel.ID > 0 {
			resStatus, err = edr.RiskLevel.GetRLInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充后续问题处理人
		if edr.IssueOwner.ID > 0 {
			resStatus, err = edr.IssueOwner.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		//填充制单人
		if edr.Creator.ID > 0 {
			resStatus, err = edr.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充确认人
		if edr.Confirmer.ID > 0 {
			resStatus, err = edr.Confirmer.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充修改人
		if edr.Modifier.ID > 0 {
			resStatus, err = edr.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充附件
		edr.Files, resStatus, err = GetEDRFiles(edr.BID)
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		eo.Body = append(eo.Body, edr)
	}

	return i18n.StatusOK, nil
}

// ExecutionOrder.GetDetailByHID 根据hid获取执行单详情
func (eo *ExecutionOrder) GetDetailByHID() (resStatus i18n.ResKey, err error) {
	//检查单据是否已经被删除
	var rowNumber int32
	checkSql := `select count(id) as rownumber from executedoc_h where id=$1 and dr=0`
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
	billNo, resStatus, err := GetLatestSerialNo(tx, "EO", eo.BillDate.Format("20060102"))
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	eo.BillNumber = billNo

	//增加表头项目
	headSql := `insert into executedoc_h(billnumber,billdate,dept_id,description,status,
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

	//表体预处理
	bodySql := `insert into executedoc_b(hid,rownumber,epaid,allowdelrow,executionvalue,
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
	//附件插入预处理
	fileSql := `insert into executedoc_file(billbid,billhid,fileid,creatorid) 
	values($1,$2,$3,$4) returning id`
	fileStmt, err := tx.Prepare(fileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Add tx.Prepare(fileSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer fileStmt.Close()
	for _, row := range eo.Body {
		var isFinish int16
		if row.IsIssue == 1 && row.IsRectify == 1 {
			isFinish = 1
		}
		//写入表体行
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
		//写入附件记录
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

	//如果数据来自于指令单,需要回写指令单status为执行态
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
	return i18n.StatusOK, nil
}

// ExecutionOrder.Edit 编辑执行单
func (eo *ExecutionOrder) Edit() (resStatus i18n.ResKey, err error) {
	//检查表体行数
	if len(eo.Body) == 0 {
		resStatus = i18n.StatusVoucherNoBody
		return
	}
	//检查创建人和编辑人是否为同一人
	if eo.Creator.ID != eo.Modifier.ID {
		resStatus = i18n.StatusVoucherOnlyCreateEdit
		return
	}

	//创建事务
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Edit db.Begin() failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	//修改表头项目
	editHeadSql := `update executedoc_h set billdate=$1,dept_id=$2,description=$3,starttime=$4,endtime=$5,
	csaid=$6,executorid=$7,modifytime=current_timestamp,modifierid=$8,ts=current_timestamp  
	where id=$9 and dr=0 and status=0 and ts=$10`
	//表头修改写入
	editHeadRes, err := tx.Exec(editHeadSql, &eo.BillDate, &eo.Department.ID, &eo.Description, &eo.StartTime, &eo.EndTime,
		&eo.CSA.ID, &eo.Executor.ID, &eo.Modifier.ID,
		&eo.HID, &eo.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Edit tx.Exec(editHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	//检查表头修改的行数
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

	//修改表体内容
	updateRowSql := `update executedoc_b set epaid=$1, allowdelrow=$2,executionvalue=$3,executionvaluedisp=$4,description=$5,
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
	//新增行准备
	addRowSql := `insert into executedoc_b(hid,rownumber,epaid,allowdelrow,executionvalue,
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
	//修改文件准备
	updateFileSql := `update executedoc_file set modifytime=current_timestamp,modifierid=$1,dr=$2,ts=current_timestamp
	where id=$3 and dr=0 and ts=$4`
	updateFileStmt, err := tx.Prepare(updateFileSql)
	if err != nil {
		zap.L().Error("ExecutionOrder.Edit tx.Prepare(updateFileSql) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer updateFileStmt.Close()
	//增加文件准备
	addFileSql := `insert into executedoc_file(billbid,billhid,fileid,creatorid) 
	values($1,$2,$3,$4) returning id`
	addFileStmt, err := tx.Prepare(addFileSql)
	if err != nil {
		zap.L().Error("ExecutionOrder.Edit tx.Prepare(addFileSql) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer addFileStmt.Close()
	//写入行数据
	for _, row := range eo.Body {
		//检查表体行状态
		if row.Status != 0 {
			resStatus = i18n.StatusVoucherNoFree
			tx.Rollback()
			return
		}
		var isFinish int16
		if row.IsIssue == 1 && row.IsRectify == 1 {
			isFinish = 1
		}
		if row.BID == 0 { //新增的行
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

			//写入附件记录
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

		} else { //原有需要更新的行
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
			//检查更新行数
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

			//处理附件
			if len(row.Files) > 0 {
				for _, file := range row.Files {
					if file.ID == 0 { //新增附件
						addFileErr := addFileStmt.QueryRow(row.BID, eo.HID, file.File.ID, eo.Modifier.ID).Scan(&file.ID)
						if addFileErr != nil {
							resStatus = i18n.StatusInternalError
							zap.L().Error("ExecutionOrder.Edit old row addFileStmt.QueryRow failed", zap.Error(addFileErr))
							tx.Rollback()
							return resStatus, addFileErr
						}
					} else { //原有附件
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
	return i18n.StatusOK, nil
}

// ExecutionOrder.Delete 删除执行单
func (eo *ExecutionOrder) Delete(modifyUserId int32) (resStatus i18n.ResKey, err error) {
	//获取单据详情
	resStatus, err = eo.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	//检查单据状态
	if eo.Status != 0 { //非自由态单据不允许删除
		resStatus = i18n.StatusVoucherNoFree
		return
	}
	//检查创建人和删除人是否为同一人
	if eo.Creator.ID != modifyUserId {
		resStatus = i18n.StatusVoucherOnlyCreateEdit
		return
	}

	//创建事务
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Delete db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	//删除表头
	delHeadSql := `update executedoc_h set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	delHeadRes, err := tx.Exec(delHeadSql, modifyUserId, eo.HID, eo.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Delete tx.Exec(delHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	//检查表头删除行数
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

	//删除表体
	delRowSql := `update executedoc_b set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	delFileSql := `update executedoc_file set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and billbid=$3 and ts=$4`
	//删除表体预处理
	delRowStmt, err := tx.Prepare(delRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Delete tx.Prepare(delRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delRowStmt.Close()
	//删除文件预处理
	delFileStmt, err := tx.Prepare(delFileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Delete tx.Prepare(delFileSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delFileStmt.Close()
	//表体删除写入
	for _, row := range eo.Body {
		//检查表体状态
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
		//检查删除影响行数
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

		if len(row.Files) > 0 { //如果存在附件
			for _, file := range row.Files {
				delFileRes, delFileErr := delFileStmt.Exec(modifyUserId, file.ID, row.BID, file.Ts)
				if delFileErr != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("ExecutionOrder.Delete delFileStmt.Exec() failed", zap.Error(delFileErr))
					tx.Rollback()
					return resStatus, delFileErr
				}
				//检查删除影响行数
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

	//如果数据来自于指令单,需要回写指令单
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

	return i18n.StatusOK, nil
}

// ExecutionOrder.Confirm 确认执行单
func (eo *ExecutionOrder) Confirm(confirmUserID int32) (resStatus i18n.ResKey, err error) {
	//获取单据详情
	resStatus, err = eo.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	//检查单据状态
	if eo.Status != 0 { //非自由态单据不允许确认
		resStatus = i18n.StatusVoucherNoFree
		return
	}
	//创建事务
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Confirm db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	//表头确认
	confirmHeadSql := `update executedoc_h set status=1,confirmtime=current_timestamp,confirmerid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and status=0 and ts=$3`
	headRes, err := tx.Exec(confirmHeadSql, confirmUserID, eo.HID, eo.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Confirm tx.Exec(confirmHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	//检查表头更新行数
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

	//确认表体行
	confirmRowSql := `update executedoc_b set status=1,confirmtime=current_timestamp,confirmerid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and status=0 and ts=$3`
	//预处理
	rowStmt, err := tx.Prepare(confirmRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.Confirm tx.Prepare(confirmRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer rowStmt.Close()
	//表体确认写入
	for _, row := range eo.Body {
		//检查表体行状态
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

	//如果数据来自于指令单,需要回写指令单
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

	return i18n.StatusOK, nil
}

// ExecutionOrder.CancelConfirm 取消确认执行单
func (eo *ExecutionOrder) CancelConfirm(confirmUserID int32) (resStatus i18n.ResKey, err error) {
	//获取单据详情
	resStatus, err = eo.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	//检查单据状态
	if eo.Status != 1 { //非确认状态单据不允许取消确认
		resStatus = i18n.StatusVoucherNoConfirm
		return
	}
	if eo.Confirmer.ID != confirmUserID {
		resStatus = i18n.StatusVoucherCancelConfirmSelf
		return
	}
	//检查表体是否有非确认状态行
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

	//创建事务
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.CancelConfirm db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	//表头取消确认
	confirmHeadSql := `update executedoc_h set status=0,confirmerid=0,ts=current_timestamp 
	where id=$1 and dr=0 and status=1 and ts=$2`
	headRes, err := tx.Exec(confirmHeadSql, eo.HID, eo.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.ExecutionOrder.CancelConfirm tx.Exec(confirmHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	//检查表头更新行数
	confirmHeadNumber, err := headRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.CancelConfirm headRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if confirmHeadNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	//取消确认表体行
	confirmRowSql := `update executedoc_b set status=0,confirmerid=0,ts=current_timestamp 
	where id=$1 and dr=0 and status=1 and ts=$2`
	//预处理
	rowStmt, err := tx.Prepare(confirmRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrder.CancelConfirm tx.Prepare(confirmRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer rowStmt.Close()
	//表体确认写入
	for _, row := range eo.Body {
		//检查表体行状态
		if row.Status != 1 {
			resStatus = i18n.StatusVoucherNoConfirm
			tx.Rollback()
			return
		}
		confirmRowRes, errConfirmRow := rowStmt.Exec(row.BID, row.Ts)
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("ExecutionOrder.CancelConfirm rowStmt.Exec failed", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}

		confirmRowNumber, errConfirmRow := confirmRowRes.RowsAffected()
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("ExecutionOrder.CancelConfirm confirmRowRes.RowsAffected failed", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}
		if confirmRowNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}
	}

	//如果数据来自于指令单,需要回写指令单
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
	return i18n.StatusOK, nil
}

// ExecutionOrderRow.Dispose 执行单表体行处理
func (edr *ExecutionOrderRow) Dispose() (resStatus i18n.ResKey, err error) {
	rowSql := `update executedoc_b set status=2,ts=current_timestamp,isfinish=$1,irfid=$2,irfnumber=$3  
	where id=$4 and hid=$5 and ts=$6 and dr=0 and status=1 and isfinish=0`
	//修改执行单行
	rowUpdateRes, err := db.Exec(rowSql, edr.IsFinish, edr.IRFID, edr.IRFNumber, edr.BID, edr.HID, edr.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrderRow.Dispose  db.Exec(rowSql) failed", zap.Error(err))
		return
	}
	//检查修改的指令单行数
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

// ExecutionOrderRow.CancelDispose 执行单表体行取消处理
func (edr *ExecutionOrderRow) CancelDispose() (resStatus i18n.ResKey, err error) {
	rowSql := `update executedoc_b set status=1,ts=current_timestamp,isfinish=$1,irfid=$2,irfnumber=$3  
	where id=$4 and hid=$5 and dr=0 and status=2 and isfinish=1`
	//修改执行单行
	rowUpdateRes, err := db.Exec(rowSql, edr.IsFinish, edr.IRFID, edr.IRFNumber, edr.BID, edr.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrderRow.CancelDispose  db.Exec(rowSql) failed", zap.Error(err))
		return
	}
	//检查修改的指令单行数
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
	return i18n.StatusOK, nil
}

// ExecutionOrderRow.Complete 执行单表体行完成
func (edr *ExecutionOrderRow) Complete() (resStatus i18n.ResKey, err error) {
	rowSql := `update executedoc_b set status=3,ts=current_timestamp   
	where id=$1 and hid=$2 and dr=0 and status=2 and isfinish=1`
	//修改执行单行
	rowUpdateRes, err := db.Exec(rowSql, edr.BID, edr.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrderRow.Complete  db.Exec(rowSql) failed", zap.Error(err))
		return
	}
	//检查修改的指令单行数
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
	return i18n.StatusOK, nil
}

// ExecutionOrderRow.CancelComplete 执行单表体行完成
func (edr *ExecutionOrderRow) CancelComplete() (resStatus i18n.ResKey, err error) {
	rowSql := `update executedoc_b set status=2,ts=current_timestamp   
	where id=$1 and hid=$2 and dr=0 and status=3 and isfinish=1`
	//修改执行单行
	rowUpdateRes, err := db.Exec(rowSql, edr.BID, edr.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrderRow.CancelComplete  db.Exec(rowSql) failed", zap.Error(err))
		return
	}
	//检查修改的指令单行数
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
	return i18n.StatusOK, nil
}

// ExectueDocCommit.Add 增加批注
func (edc *ExecutionOrderComment) Add() (resStatus i18n.ResKey, err error) {
	sqlStr := `insert into executedoc_comment(bid,hid,billnumber,rownumber,sendto_id,
	content,creatorid,sendtime,createtime)
	values($1,$2,$3,$4,$5,$6,$7,
	to_char(current_timestamp,'YYYYMMDDHH24MI'),current_timestamp) returning id`
	err = db.QueryRow(sqlStr, edc.BID, edc.HID, edc.BillNumber, edc.RowNUmber, edc.SendTo.ID,
		edc.Content, edc.Creator.ID).Scan(&edc.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrderComment.Add db.QueryRow(sqlStr) failed", zap.Error(err))
		return
	}
	resStatus = i18n.StatusOK
	return
}

// ExecutionOrderReview.Add 增加审阅记录
func (edrr *ExecutionOrderReview) Add() (resStatus i18n.ResKey, err error) {
	sqlStr := `insert into executedoc_review(hid,billnumber,starttime,endtime,consumeseconds,creatorid) 
	values($1,$2,$3,$4,$5,$6) returning id`
	err = db.QueryRow(sqlStr, edrr.HID, edrr.BillNumber, edrr.StartTime, edrr.EndTime, edrr.ConsumeSeconds, edrr.Creator.ID).Scan(&edrr.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionOrderReview.Add db.QueryRow(sqlStr) failed", zap.Error(err))
		return
	}
	resStatus = i18n.StatusOK
	return
}

// EOCommentsParams.Get 获取单据批注列表
func (cs *EOCommentsParams) Get() (resStatus i18n.ResKey, err error) {
	cs.Comments = make([]ExecutionOrderComment, 0) //解决返回列表为null
	sqlStr := `select id,bid,hid,billnumber,rownumber,
	sendto_id,isread,readtime,content,sendtime,
	createtime,creatorid,dr,ts from executedoc_comment where dr=0 and hid=$1`
	rows, err := db.Query(sqlStr, cs.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EOCommentsParams.Get db.Query failed", zap.Error(err))
		return
	}
	defer rows.Close()

	//提取数据
	for rows.Next() {
		var edc ExecutionOrderComment
		err = rows.Scan(&edc.ID, &edc.BID, &edc.HID, &edc.BillNumber, &edc.RowNUmber,
			&edc.SendTo.ID, &edc.IsRead, &edc.ReadTime, &edc.Content, &edc.SendTime,
			&edc.CreateDate, &edc.Creator.ID, &edc.Dr, &edc.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("EOCommentsParams.Get rows.Next failed", zap.Error(err))
			return
		}
		//填充发送人信息
		if edc.SendTo.ID > 0 {
			resStatus, err = edc.SendTo.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充发送人信息
		if edc.Creator.ID > 0 {
			resStatus, err = edc.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		cs.Comments = append(cs.Comments, edc)
	}

	return i18n.StatusOK, nil
}

// EOReviewsParams.Get 获取审阅记录模型
func (rs *EOReviewsParams) Get() (resStatus i18n.ResKey, err error) {
	rs.Reviews = make([]ExecutionOrderReview, 0)
	sqlStr := `select id,hid,billnumber,starttime,endtime,
	consumeseconds,createtime,creatorid,dr,ts 
	from executedoc_review where dr=0 and hid=$1`
	rows, err := db.Query(sqlStr, rs.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EOReviewsParams.Get db.Query failed", zap.Error(err))
		return
	}
	defer rows.Close()

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
	return i18n.StatusOK, nil
}
