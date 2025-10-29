package pg

import (
	"sccsmsserver/i18n"
	"sccsmsserver/setting"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Work Order Status Report struct
type WorkOrderReport struct {
	WoBID           int32     `json:"woBID"`
	WoHID           int32     `json:"woHID"`
	WoBillDate      time.Time `json:"woBillDate"`
	WoBillNumber    string    `json:"woBillNumber"`
	WoRowNumber     int32     `json:"woRowNumber"`
	CSAID           int32     `json:"csaID"`
	CSACode         string    `json:"csaCode"`
	CSAName         string    `json:"csaName"`
	RespPersonID    int32     `json:"respPersonID"`
	RespPersonCode  string    `json:"respPersonCode"`
	RespPersonName  string    `json:"respPersonName"`
	RespDeptID      int32     `json:"respDeptID"`
	RespDeptCode    string    `json:"respDeptCode"`
	RespDeptName    string    `json:"respDeptName"`
	ExecutorID      int32     `json:"executorID"`
	ExecutorCode    string    `json:"executorCode"`
	ExecutorName    string    `json:"executorName"`
	WorDescription  string    `json:"worDescription"`
	EPTID           int32     `json:"eptID"`
	EPTCode         string    `json:"eptCode"`
	EPTName         string    `json:"eptName"`
	WoStartTime     time.Time `json:"woStartTime"`
	WoEndTime       time.Time `json:"woEndTime"`
	WorStatus       int16     `json:"worStatus"`
	WoCreateDate    time.Time `json:"woCreateDate"`
	WoCreatorID     int32     `json:"woCreatorID"`
	WoCreatorCode   string    `json:"woCreatorCode"`
	WoCreatorName   string    `json:"woCreatorName"`
	WoConfirmDate   string    `json:"woConfirmDate"`
	WoConfirmerID   int32     `json:"woConfirmerID"`
	WoConfirmerCode string    `json:"woConfirmerCode"`
	WoConfirmerName string    `json:"woConfirmerName"`
	WoDeptID        int32     `json:"woDeptID"`
	WoDeptCode      string    `json:"woDeptCode"`
	WoDeptName      string    `json:"woDeptName"`
	WoDescription   string    `json:"woDescription"`
	WoStatus        int16     `json:"woStatus"`
	WoWorkDate      time.Time `json:"woWorkDate"`
	EoHID           int32     `json:"eoHID"`
	EoNumber        string    `json:"eoNumber"`
	EoCreatorID     int32     `json:"eoCreatorID"`
	EoCreatorCode   string    `json:"eoCreatorCode"`
	EoCreatorName   string    `json:"eoCreatorName"`
	EoBillDate      time.Time `json:"eoBillDate"`
	EoStartTime     time.Time `json:"eoStartTime"`
	EoEndTIme       time.Time `json:"eoEndTime"`
	EoHStatus       int16     `json:"eoHStatus"`
	Udf1Code        string    `json:"udf1Code"`
	Udf1Name        string    `json:"udf1Name"`
	Udf2Code        string    `json:"udf2Code"`
	Udf2Name        string    `json:"udf2Name"`
	Udf3Code        string    `json:"udf3Code"`
	Udf3Name        string    `json:"udf3Name"`
	Udf4Code        string    `json:"udf4Code"`
	Udf4Name        string    `json:"udf4Name"`
	Udf5Code        string    `json:"udf5Code"`
	Udf5Name        string    `json:"udf5Name"`
	Udf6Code        string    `json:"udf6Code"`
	Udf6Name        string    `json:"udf6Name"`
	Udf7Code        string    `json:"udf7Code"`
	Udf7Name        string    `json:"udf7Name"`
	Udf8Code        string    `json:"udf8Code"`
	Udf8Name        string    `json:"udf8Name"`
	Udf9Code        string    `json:"udf9Code"`
	Udf9Name        string    `json:"udf9Name"`
	Udf10Code       string    `json:"udf10Code"`
	Udf10Name       string    `json:"udf10Name"`
}

// Execution Order Status Report struct
type ExecutionOrderReport struct {
	BID                int32     `json:"bid"`
	HID                int32     `json:"hid"`
	BillNumber         string    `json:"billNumber"`
	RowNumber          int32     `json:"rowNumber"`
	BillDate           time.Time `json:"billDate"`
	HDeptID            int32     `json:"hDeptID"`
	HDeptCode          string    `json:"hDeptCode"`
	HDeptName          string    `json:"hDeptName"`
	HDescription       string    `json:"hDescription"`
	HStatus            int16     `json:"hStatus"`
	SourceType         string    `json:"sourceType"`
	SourceHID          string    `json:"sourceHID"`
	SourceBillnumber   string    `json:"sourceBillNumber"`
	SourceRowNumber    int32     `json:"sourceRowNumber"`
	SourceBID          int32     `json:"sourceBID"`
	HStartTime         time.Time `json:"hStartTime"`
	HEndTime           time.Time `json:"hEndTime"`
	CSAID              int32     `json:"csaID"`
	CSACode            string    `json:"csaCode"`
	CSAName            string    `json:"csaName"`
	CSCID              int32     `json:"cscID"`
	ExecutorID         int32     `json:"executorID"`
	ExecutorCode       string    `json:"executorCode"`
	ExecutorName       string    `json:"executorName"`
	EPTID              int32     `json:"eptID"`
	EPTCode            string    `json:"eptCode"`
	EPTName            string    `json:"eptName"`
	EPAID              int32     `json:"epaID"`
	EPACode            string    `json:"epaCode"`
	EPAName            string    `json:"epaName"`
	RLID               int32     `json:"rlID"`
	RLName             string    `json:"rlName"`
	RLColor            string    `json:"rlColor"`
	ExecutionValue     string    `json:"executionValue"`
	ExecutionValueDIsp string    `json:"executionValueDisp"`
	BDescription       string    `json:"bDescription"`
	IsCheckError       int16     `json:"isCheckError"`
	IsRequireFile      int16     `json:"isRequireFile"`
	IsOnsitePhoto      int16     `json:"isOnSitePhoto"`
	IsIssue            int16     `json:"isIssue"`
	IsRectify          int16     `json:"isRectify"`
	IsHandle           int16     `json:"isHandle"`
	IssueOwnerID       int32     `json:"issueOwnerID"`
	IssueOwnerCode     string    `json:"issueOwnerCode"`
	IssueOwnerName     string    `json:"issueOwnerName"`
	HandleStartTime    time.Time `json:"handleStartTime"`
	HandleEndTime      time.Time `json:"handleEndTime"`
	BStatus            int16     `json:"bStatus"`
	IsFromEPT          int16     `json:"isFromEpt"`
	IsFinish           int16     `json:"isFinish"`
	IRFID              int32     `json:"irfID"`
	IRFNumber          string    `json:"irfNumber"`
	CreateDate         time.Time `json:"createDate"`
	CreatorID          int32     `json:"creatorID"`
	CreatorCode        string    `json:"creatorCode"`
	CreatorName        string    `json:"creatorName"`
	ConfirmDate        time.Time `json:"confirmDate"`
	ConfirmerID        int32     `json:"confirmerID"`
	ConfirmerCode      string    `json:"confirmCode"`
	ConfirmerName      string    `json:"confirmName"`
	Udf1Code           string    `json:"udf1Code"`
	Udf1Name           string    `json:"udf1Name"`
	Udf2Code           string    `json:"udf2Code"`
	Udf2Name           string    `json:"udf2Name"`
	Udf3Code           string    `json:"udf3Code"`
	Udf3Name           string    `json:"udf3Name"`
	Udf4Code           string    `json:"udf4Code"`
	Udf4Name           string    `json:"udf4Name"`
	Udf5Code           string    `json:"udf5Code"`
	Udf5Name           string    `json:"udf5Name"`
	Udf6Code           string    `json:"udf6Code"`
	Udf6Name           string    `json:"udf6Name"`
	Udf7Code           string    `json:"udf7Code"`
	Udf7Name           string    `json:"udf7Name"`
	Udf8Code           string    `json:"udf8Code"`
	Udf8Name           string    `json:"udf8Name"`
	Udf9Code           string    `json:"udf9Code"`
	Udf9Name           string    `json:"udf9Name"`
	Udf10Code          string    `json:"udf10Code"`
	Udf10Name          string    `json:"udf10Name"`
}

// Issue Resolution Form Status Report struct
type IssueResolutionFormReport struct {
	EOBID              int32         `json:"eoBID"`
	EOHID              int32         `json:"eoHID"`
	EOBillNumber       string        `json:"eoBillNumber"`
	EORowNumber        int32         `json:"eoRowNumber"`
	EOBillDate         time.Time     `json:"eoBillDate"`
	EOHDeptID          int32         `json:"eoHDeptID"`
	EOHDeptCode        string        `json:"eoHDeptCode"`
	EOHDeptName        string        `json:"eoHDeptName"`
	CSAID              int32         `json:"csaID"`
	CSACode            string        `json:"csaCode"`
	CSAName            string        `json:"csaName"`
	CSCID              int32         `json:"cscID"`
	ExecutorID         int32         `json:"executorID"`
	ExecutorCode       string        `json:"executorCode"`
	ExecutorName       string        `json:"executorName"`
	EPAID              int32         `json:"epaID"`
	EPACode            string        `json:"epaCode"`
	EPAName            string        `json:"epaName"`
	RLID               int32         `json:"rlID"`
	RLName             string        `json:"rlName"`
	RLColor            string        `json:"rlColor"`
	ExecutionValue     string        `json:"executionValue"`
	ExecutionValueDIsp string        `json:"executionValueDisp"`
	EOBDescription     string        `json:"eoBDescription"`
	IsIssue            int16         `json:"isIssue"`
	IsRectify          int16         `json:"isRectify"`
	IsHandle           int16         `json:"isHandle"`
	IssueOwnerID       int32         `json:"issueOwnerID"`
	IssueOwnerCode     string        `json:"issueOwnerCode"`
	IssueOwnerName     string        `json:"issueOwnerName"`
	EOBStartTime       time.Time     `json:"eoBStartTime"`
	EOBEndTime         time.Time     `json:"eoBEndTime"`
	IsFinish           int16         `json:"isFinish"`
	IRFID              int32         `json:"irfID"`
	IRFBillNumber      string        `json:"irfBillNumber"`
	IRFBillDate        time.Time     `json:"irfBillDate"`
	HandlerID          int32         `json:"handlerID"`
	HandlerCode        string        `json:"handlerCode"`
	HandlerName        string        `json:"handlerName"`
	IRFStartTime       time.Time     `json:"irfStartTime"`
	IRFEndTime         time.Time     `json:"irfEndTime"`
	IRFDescription     string        `json:"irfDescription"`
	IRFStatus          int16         `json:"irfStatus"`
	CreatorID          int32         `json:"creatorID"`
	CreatorCode        string        `json:"creatorCode"`
	CreatorName        string        `json:"creatorName"`
	ConfirmerID        int32         `json:"confirmerID"`
	ConfirmerCode      string        `json:"confirmerCode"`
	ConfirmerName      string        `json:"confirmerName"`
	EORFiles           []VoucherFile `json:"eorFiles"`
	IRFFiles           []VoucherFile `json:"irfFiles"`
	Udf1Code           string        `json:"udf1Code"`
	Udf1Name           string        `json:"udf1Name"`
	Udf2Code           string        `json:"udf2Code"`
	Udf2Name           string        `json:"udf2Name"`
	Udf3Code           string        `json:"udf3Code"`
	Udf3Name           string        `json:"udf3Name"`
	Udf4Code           string        `json:"udf4Code"`
	Udf4Name           string        `json:"udf4Name"`
	Udf5Code           string        `json:"udf5Code"`
	Udf5Name           string        `json:"udf5Name"`
	Udf6Code           string        `json:"udf6Code"`
	Udf6Name           string        `json:"udf6Name"`
	Udf7Code           string        `json:"udf7Code"`
	Udf7Name           string        `json:"udf7Name"`
	Udf8Code           string        `json:"udf8Code"`
	Udf8Name           string        `json:"udf8Name"`
	Udf9Code           string        `json:"udf9Code"`
	Udf9Name           string        `json:"udf9Name"`
	Udf10Code          string        `json:"udf10Code"`
	Udf10Name          string        `json:"udf10Name"`
}

// Get Work Order Status Report
func GetWorkOrderReport(queryString string) (wors []WorkOrderReport, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	wors = make([]WorkOrderReport, 0)
	var build strings.Builder
	// Concatenate the SQL string for inspection
	build.WriteString(`select count(b.id) as rownumber
	from workorder_b as b
	left join workorder_h as h on b.hid = h.id
	left join executionorder_h as eoh on b.eoid = eoh.id
	left join csa on b.csaid = csa.id
	left join uda as udf1 on csa.udf1 = udf1.id
	left join uda as udf2 on csa.udf2 = udf2.id
	left join uda as udf3 on csa.udf3 = udf3.id
	left join uda as udf4 on csa.udf4 = udf4.id
	left join uda as udf5 on csa.udf5 = udf5.id
	left join uda as udf6 on csa.udf6 = udf6.id
	left join uda as udf7 on csa.udf7 = udf7.id
	left join uda as udf8 on csa.udf8 = udf8.id
	left join uda as udf9 on csa.udf9 = udf9.id
	left join uda as udf10 on csa.udf10 = udf10.id
	left join sysuser as executor on b.executorid = executor.id
	left join sysuser as respperson on csa.resppersonid = respperson.id
	left join sysuser as acturalep on eoh.creatorid = acturalep.id
	left join sysuser as creator on b.creatorid = creator.id
	left join sysuser as confirmer on b.confirmerid = confirmer.id
	left join ept_h on b.eptid = ept_h.id
	left join department as hdept on h.deptid = hdept.id
	left join department as respdept on csa.respdeptid = respdept.id
	where (b.dr=0) `)
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
		zap.L().Error("GetWorkOrderReport db.QueryRow(checkSql) failed", zap.Error(err))
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

	// Concatenate the SQL for data retrieval
	build.WriteString(`select 
	b.id as wobid,
	b.hid as wohid,
	h.billdate as wobilldate,
	h.billnumber as wobillnumber,
	b.rownumber as worownumber,
	b.csaid,
	csa.code as csacode,
	csa.name as csaname,
	coalesce(csa.resppersonid,0) as resppersonid,
	coalesce(respperson.code,'') as resppersoncode,
	coalesce(respperson.name,'') as resppersonname,
	coalesce(csa.respdeptid,0) as respdeptid,
	coalesce(respdept.code,'') as respdeptcode,
	coalesce(respdept.name,'') as respdeptname,
	b.executorid as executorid,
	executor.code as executorcode,
	executor.name as executorname,
	b.description as wordescription,
	b.eptid as eptid,
	ept_h.code as eptcode,
	ept_h.name as eptname,
	b.starttime as wostarttime,
	b.endtime as woendtime,
	b.status as worstatus,
	b.createtime as wocreatedate,
	b.creatorid as wocreateuserid,
	creator.code as wocreatorcode,
	creator.name as wocreatorname,
	b.confirmtime as confirmdate,
	coalesce(b.confirmerid,0) as woconfirmerid,
	coalesce(confirmer.code,'') as woconfirmercode,
	coalesce(confirmer.name,'') as woconfirmername,
	h.deptid as wodeptid,
	coalesce(hdept.code,'') as wodeptcode,
	coalesce(hdept.name,'') as wodeptname,
	h.description as wodescription,
	h.status as wotatus,
	h.workdate as woworkdate,
	b.eoid as eoid,
	b.eonumber as eonumber,
	coalesce(eoh.creatorid,0) as eocreatorid,
	coalesce(acturalep.code,'') as eocreatorcode,
	coalesce(acturalep.name,'') as eocreatorname,
	coalesce(eoh.billdate,to_timestamp(0)) as eobilldate,
	coalesce(eoh.starttime,to_timestamp(0)) as eobstarttime,
	coalesce(eoh.endtime,to_timestamp(0)) as eoendtime,
	coalesce(eoh.status,4) as eohstatus,
	coalesce(udf1.name,'') as udf1name,
	coalesce(udf1.code,'') as udf1code,
	coalesce(udf2.name,'') as udf2name,
	coalesce(udf2.code,'') as udf2code,
	coalesce(udf3.name,'') as udf3name,
	coalesce(udf3.code,'') as udf3code,
	coalesce(udf4.name,'') as udf4name,
	coalesce(udf4.code,'') as udf4code,
	coalesce(udf5.name,'') as udf5name,
	coalesce(udf5.code,'') as udf5code,
	coalesce(udf6.name,'') as udf6name,
	coalesce(udf6.code,'') as udf6code,
	coalesce(udf7.name,'') as udf7name,
	coalesce(udf7.code,'') as udf7code,
	coalesce(udf8.name,'') as udf8name,
	coalesce(udf8.code,'') as udf8code,
	coalesce(udf9.name,'') as udf9name,
	coalesce(udf9.code,'') as udf9code,
	coalesce(udf10.name,'') as udf10name,
	coalesce(udf10.code,'') as udf10code
	from workorder_b as b
	left join workorder_h as h on b.hid = h.id
	left join executionorder_h as eoh on b.eoid = eoh.id
	left join csa as csa on b.csaid = csa.id
	left join uda as udf1 on csa.udf1 = udf1.id
	left join uda as udf2 on csa.udf2 = udf2.id
	left join uda as udf3 on csa.udf3 = udf3.id
	left join uda as udf4 on csa.udf4 = udf4.id
	left join uda as udf5 on csa.udf5 = udf5.id
	left join uda as udf6 on csa.udf6 = udf6.id
	left join uda as udf7 on csa.udf7 = udf7.id
	left join uda as udf8 on csa.udf8 = udf8.id
	left join uda as udf9 on csa.udf9 = udf9.id
	left join uda as udf10 on csa.udf10 = udf10.id
	left join sysuser as executor on b.executorid = executor.id
	left join sysuser as respperson on csa.resppersonid = respperson.id
	left join sysuser as acturalep on eoh.creatorid = acturalep.id
	left join sysuser as creator on b.creatorid = creator.id
	left join sysuser as confirmer on b.confirmerid = confirmer.id
	left join ept_h on b.eptid = ept_h.id
	left join department as hdept on h.deptid = hdept.id
	left join department as respdept on csa.respdeptid = respdept.id
	where (b.dr=0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	repSql := build.String()
	// Retrieve Work Order status list from database
	woRep, err := db.Query(repSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetWorkOrderReport db.Query failed", zap.Error(err))
		return
	}
	defer woRep.Close()
	// Extract data row by row
	for woRep.Next() {
		var wor WorkOrderReport
		err = woRep.Scan(&wor.WoBID, &wor.WoHID, &wor.WoBillDate, &wor.WoBillNumber, &wor.WoRowNumber,
			&wor.CSAID, &wor.CSACode, &wor.CSAName, &wor.RespPersonID, &wor.RespPersonCode,
			&wor.RespPersonName, &wor.RespDeptID, &wor.RespDeptCode, &wor.RespDeptName, &wor.ExecutorID,
			&wor.ExecutorCode, &wor.ExecutorName, &wor.WorDescription, &wor.EPTID, &wor.EPTCode,
			&wor.EPTName, &wor.WoStartTime, &wor.WoEndTime, &wor.WorStatus, &wor.WoCreateDate,
			&wor.WoCreatorID, &wor.WoCreatorCode, &wor.WoCreatorName, &wor.WoConfirmDate, &wor.WoConfirmerID,
			&wor.WoConfirmerCode, &wor.WoConfirmerName, &wor.WoDeptID, &wor.WoDeptCode, &wor.WoDeptName,
			&wor.WoDescription, &wor.WoStatus, &wor.WoWorkDate, &wor.EoHID, &wor.EoNumber,
			&wor.EoCreatorID, &wor.EoCreatorCode, &wor.EoCreatorName, &wor.EoBillDate, &wor.EoStartTime,
			&wor.EoEndTIme, &wor.EoHStatus,
			&wor.Udf1Name, &wor.Udf1Code, &wor.Udf2Name, &wor.Udf2Code, &wor.Udf3Name, &wor.Udf3Code, &wor.Udf4Name, &wor.Udf4Code, &wor.Udf5Name, &wor.Udf5Code,
			&wor.Udf6Name, &wor.Udf6Code, &wor.Udf7Name, &wor.Udf7Code, &wor.Udf8Name, &wor.Udf8Code, &wor.Udf9Name, &wor.Udf9Code, &wor.Udf10Name, &wor.Udf10Code)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetWorkOrderReport woRep.Next() woRep.Scan failed", zap.Error(err))
			return
		}

		wors = append(wors, wor)
	}

	return
}

// Get Execution Order status report
func GetExecutionOrderReport(queryString string) (eors []ExecutionOrderReport, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	eors = make([]ExecutionOrderReport, 0)
	var build strings.Builder
	// Concatenate the SQL string for check
	build.WriteString(`select count(b.id) as rownumber 
	from executionorder_b as b
	left join executionorder_h as h on b.hid = h.id
	left join department as dept on h.deptid = dept.id
	left join epa on b.epaid = epa.id
	left join risklevel as rl on b.risklevelid = rl.id
	left join sysuser as creator on b.creatorid = creator.id
	left join sysuser as confirmer on b.confirmerid = confirmer.id
	left join sysuser as issueowner on b.issueownerid = issueowner.id
	left join sysuser as executor on h.executorid = executor.id
	left join ept_h as ept_h on h.eptid = ept_h.id
	left join csa as csa on h.csaid = csa.id
	left join uda as udf1 on csa.udf1 = udf1.id
	left join uda as udf2 on csa.udf2 = udf2.id
	left join uda as udf3 on csa.udf3 = udf3.id
	left join uda as udf4 on csa.udf4 = udf4.id
	left join uda as udf5 on csa.udf5 = udf5.id
	left join uda as udf6 on csa.udf6 = udf6.id
	left join uda as udf7 on csa.udf7 = udf7.id
	left join uda as udf8 on csa.udf8 = udf8.id
	left join uda as udf9 on csa.udf9 = udf9.id
	left join uda as udf10 on csa.udf10 = udf10.id
	where (b.dr=0)`)
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
		zap.L().Error("GetExecutionOrderReport db.QueryRow(checkSql) failed", zap.Error(err))
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

	// Concatenate the SQL string for data retrieve
	build.WriteString(`select b.id as bid,
	b.hid as hid,
	h.billnumber as billnumber,
	b.rownumber as rownumber,
	h.billdate as billdate,
	h.deptid as hdeptid,
	coalesce(dept.code,'') as hdeptcode,
	coalesce(dept.name,'') as hdeptname,
	h.description as hdescription,
	h.status as hstatus,
	h.sourcetype as sourcetype,
	h.sourcehid as sourcehid,
	h.sourcebillnumber as sourcebillnumber,
	h.sourcerownumber as sourcerownumber,
	h.sourcebid sourcebid,
	h.starttime as hstarttime,
	h.endtime as hendtime,
	h.csaid as siid,
	coalesce(csa.code,'') as csacode,
	coalesce(csa.name,'') as csaname,
	coalesce(csa.cscid,0) as cscid,
	h.executorid as executorid,
	coalesce(executor.code,'') as executorcode,
	coalesce(executor.name,'') as executorname,
	h.eptid as eptid,
	coalesce(ept_h.code,'') as eptcode,
	coalesce(ept_h.name,'') as eptname,
	b.epaid as epaid,
	coalesce(epa.code,'') as epacode,
	coalesce(epa.name,'') as epaname,
	b.risklevelid as rlid,
	coalesce(rl.name,'') as rlname,
	coalesce(rl.color,'white') as rlcolor,
	b.executionvalue as executionvalue,
	b.executionvaluedisp as executionvaluedisp,
	b.description as bdescription,
	b.ischeckerror as ischeckerr,
	b.isrequirefile as isrequirefile,
	b.isonsitephoto as isonsitephoto,
	b.isissue as isissue,
	b.isrectify as isrectify,
	b.ishandle as ishandle,
	case b.ishandle when 1 then b.issueownerid else 0 end issueownerid,
	case b.ishandle when 1 then issueowner.code else '' end issueownercode,
	case b.ishandle when 1 then issueowner.name else '' end issueownername,
	case b.ishandle when 1 then b.handlestarttime else to_timestamp(0) end handlestarttime,
	case b.ishandle when 1 then b.handleendtime else to_timestamp(0) end handleendtime,
	b.status as bstatus,
	b.isfromept as isfromept,
	b.isfinish as isfinish,
	b.irfid as irfid,
	b.irfnumber as irfnumber,
	b.createtime as createdate,
	b.creatorid as creatorid,
	coalesce(creator.code,'') as creatorcode,
	coalesce(creator.name,'') as creatorname,
	b.confirmtime as confirmdate,
	b.confirmerid as confirmerid,
	coalesce(confirmer.code,'') as confirmercode,
	coalesce(confirmer.name,'') as confirmername,
	coalesce(udf1.name,'') as udf1name,
	coalesce(udf1.code,'') as udf1code,
	coalesce(udf2.name,'') as udf2name,
	coalesce(udf2.code,'') as udf2code,
	coalesce(udf3.name,'') as udf3name,
	coalesce(udf3.code,'') as udf3code,
	coalesce(udf4.name,'') as udf4name,
	coalesce(udf4.code,'') as udf4code,
	coalesce(udf5.name,'') as udf5name,
	coalesce(udf5.code,'') as udf5code,
	coalesce(udf6.name,'') as udf6name,
	coalesce(udf6.code,'') as udf6code,
	coalesce(udf7.name,'') as udf7name,
	coalesce(udf7.code,'') as udf7code,
	coalesce(udf8.name,'') as udf8name,
	coalesce(udf8.code,'') as udf8code,
	coalesce(udf9.name,'') as udf9name,
	coalesce(udf9.code,'') as udf9code,
	coalesce(udf10.name,'') as udf10name,
	coalesce(udf10.code,'') as udf10code
	from executionorder_b as b
	left join executionorder_h as h on b.hid = h.id
	left join department as dept on h.deptid = dept.id
	left join epa as epa on b.epaid = epa.id
	left join risklevel as rl on b.risklevelid = rl.id
	left join sysuser as creator on b.creatorid = creator.id
	left join sysuser as confirmer on b.confirmerid = confirmer.id
	left join sysuser as issueowner on b.issueownerid = issueowner.id
	left join sysuser as executor on h.executorid = executor.id
	left join ept_h as ept_h on h.eptid = ept_h.id
	left join csa as csa on h.csaid = csa.id
	left join uda as udf1 on csa.udf1 = udf1.id
	left join uda as udf2 on csa.udf2 = udf2.id
	left join uda as udf3 on csa.udf3 = udf3.id
	left join uda as udf4 on csa.udf4 = udf4.id
	left join uda as udf5 on csa.udf5 = udf5.id
	left join uda as udf6 on csa.udf6 = udf6.id
	left join uda as udf7 on csa.udf7 = udf7.id
	left join uda as udf8 on csa.udf8 = udf8.id
	left join uda as udf9 on csa.udf9 = udf9.id
	left join uda as udf10 on csa.udf10 = udf10.id
	where (b.dr=0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	repSql := build.String()
	// Retrieve Execution Order Reports from database
	eoRep, err := db.Query(repSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetExecutionOrderReport db.Query failed", zap.Error(err))
		return
	}
	defer eoRep.Close()

	// Extract data row by row
	for eoRep.Next() {
		var eor ExecutionOrderReport
		err = eoRep.Scan(&eor.BID, &eor.HID, &eor.BillNumber, &eor.RowNumber, &eor.BillDate,
			&eor.HDeptID, &eor.HDeptCode, &eor.HDeptName, &eor.HDescription, &eor.HStatus,
			&eor.SourceType, &eor.SourceHID, &eor.SourceBillnumber, &eor.SourceRowNumber, &eor.SourceBID,
			&eor.HStartTime, &eor.HEndTime, &eor.CSAID, &eor.CSACode, &eor.CSAName,
			&eor.CSCID, &eor.ExecutorID, &eor.ExecutorCode, &eor.ExecutorName, &eor.EPTID,
			&eor.EPTCode, &eor.EPTName, &eor.EPAID, &eor.EPACode, &eor.EPAName, &eor.RLID, &eor.RLName, &eor.RLColor,
			&eor.ExecutionValue, &eor.ExecutionValueDIsp, &eor.BDescription, &eor.IsCheckError, &eor.IsRequireFile,
			&eor.IsOnsitePhoto, &eor.IsIssue, &eor.IsRectify, &eor.IsHandle, &eor.IssueOwnerID,
			&eor.IssueOwnerCode, &eor.IssueOwnerName, &eor.HandleStartTime, &eor.HandleEndTime, &eor.BStatus,
			&eor.IsFromEPT, &eor.IsFinish, &eor.IRFID, &eor.IRFNumber, &eor.CreateDate,
			&eor.CreatorID, &eor.CreatorCode, &eor.CreatorName, &eor.ConfirmDate, &eor.ConfirmerID,
			&eor.ConfirmerCode, &eor.ConfirmerName,
			&eor.Udf1Name, &eor.Udf1Code, &eor.Udf2Name, &eor.Udf2Code, &eor.Udf3Name,
			&eor.Udf3Code, &eor.Udf4Name, &eor.Udf4Code, &eor.Udf5Name, &eor.Udf5Code,
			&eor.Udf6Name, &eor.Udf6Code, &eor.Udf7Name, &eor.Udf7Code, &eor.Udf8Name,
			&eor.Udf8Code, &eor.Udf9Name, &eor.Udf9Code, &eor.Udf10Name, &eor.Udf10Name)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetExecutionOrderReport eoRep.Next() eoRef.Scan failed", zap.Error(err))
			return
		}
		eors = append(eors, eor)
	}
	return
}

// Get Issue Resolution Form Report
func GetIssueResolutionFormReport(queryString string) (irfs []IssueResolutionFormReport, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	irfs = make([]IssueResolutionFormReport, 0)
	var build strings.Builder
	// Concatenate the SQL string for
	build.WriteString(`select count(b.id) as rowcount
	from executionorder_b as b
	left join executionorder_h as h on b.hid = h.id
	left join issueresolutionform as irf on b.irfid = irf.id
	left join sysuser as handler on irf.handlerid = handler.id
	left join department as dept on h.deptid = dept.id
	left join epa as epa on b.epaid = epa.id
	left join risklevel as rl on b.risklevelid = rl.id
	left join sysuser as creator on irf.creatorid = creator.id
	left join sysuser as confirmer on irf.confirmerid = confirmer.id
	left join sysuser as issueowner on b.issueownerid = issueowner.id
	left join sysuser as executor on h.executorid = executor.id
	left join ept_h as ept_h on h.eptid = ept_h.id
	left join csa as csa on h.csaid = csa.id
	left join uda as udf1 on csa.udf1 = udf1.id
	left join uda as udf2 on csa.udf2 = udf2.id
	left join uda as udf3 on csa.udf3 = udf3.id
	left join uda as udf4 on csa.udf4 = udf4.id
	left join uda as udf5 on csa.udf5 = udf5.id
	left join uda as udf6 on csa.udf6 = udf6.id
	left join uda as udf7 on csa.udf7 = udf7.id
	left join uda as udf8 on csa.udf8 = udf8.id
	left join uda as udf9 on csa.udf9 = udf9.id
	left join uda as udf10 on csa.udf10 = udf10.id
	where (b.dr=0 and b.isissue = 1)`)
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
		zap.L().Error("GetIssueResolutionFormReport db.QueryRow(checkSql) failed", zap.Error(err))
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

	// Concatenate the SQL string for data retrieve
	build.WriteString(`select b.id as edbid,
	b.hid as edhid,
	h.billnumber as edbillnumber,
	b.rownumber as edrownumber,
	h.billdate as eobilldate,
	h.deptid as edhdeptid,
	coalesce(dept.code,'') as edhdeptcode,
	coalesce(dept.name,'') as edhdeptname,
	h.csaid as siid,
	coalesce(csa.code,'') as csacode,
	coalesce(csa.name,'') as csaname,
	coalesce(csa.cscid,0) as cscid,
	h.executorid as executorid,
	coalesce(executor.code,'') as executorcode,
	coalesce(executor.name,'') as executorname,
	b.epaid as epaid,
	coalesce(epa.code,'') as epacode,
	coalesce(epa.name,'') as epaname,
	b.risklevelid as rlid,
	coalesce(rl.name,'') as rlname,
	coalesce(rl.color,'white') as rlcolor,
	b.executionvalue as executionvalue,
	b.executionvaluedisp as executionvaluedisp,
	b.description as edbdescription,
	b.isissue as isissue,
	b.isrectify as isrectify,
	b.ishandle as ishandle,
	case b.ishandle when 1 then b.issueownerid else 0 end issueownerid,
	case b.ishandle when 1 then issueowner.code else '' end issueownercode,
	case b.ishandle when 1 then issueowner.name else '' end issueownername,
	case b.ishandle when 1 then b.handlestarttime else to_timestamp(0) end eobstarttime,
	case b.ishandle when 1 then b.handleendtime else to_timestamp(0) end eobendtime,
	b.isfinish as isfinish,
	b.irfid as irfid,
	coalesce(irf.billnumber,'') as irfbillnumber,
	coalesce(irf.billdate,to_timestamp(0)) as irfbilldate,
	coalesce(irf.handlerid,0) as handlerid,
	coalesce(handler.code,'') as handlercode,
	coalesce(handler.name,'') as handlername,
	coalesce(irf.starttime,to_timestamp(0)) as irfstarttime,
	coalesce(irf.endtime,to_timestamp(0)) as irfendtime,
	coalesce(irf.description,'') as irfdescription,
	coalesce(irf.status,0) as irfstatus,
	coalesce(irf.creatorid,0) as creatorid,
	coalesce(creator.code,'') as creatorcode,
	coalesce(creator.name,'') as creatorname,
	coalesce(irf.confirmerid,0) as confirmerid,
	coalesce(confirmer.code,'') as confirmercode,
	coalesce(confirmer.name,'') as confirmername,
	coalesce(udf1.name,'') as udf1name,
	coalesce(udf1.code,'') as udf1code,
	coalesce(udf2.name,'') as udf2name,
	coalesce(udf2.code,'') as udf2code,
	coalesce(udf3.name,'') as udf3name,
	coalesce(udf3.code,'') as udf3code,
	coalesce(udf4.name,'') as udf4name,
	coalesce(udf4.code,'') as udf4code,
	coalesce(udf5.name,'') as udf5name,
	coalesce(udf5.code,'') as udf5code,
	coalesce(udf6.name,'') as udf6name,
	coalesce(udf6.code,'') as udf6code,
	coalesce(udf7.name,'') as udf7name,
	coalesce(udf7.code,'') as udf7code,
	coalesce(udf8.name,'') as udf8name,
	coalesce(udf8.code,'') as udf8code,
	coalesce(udf9.name,'') as udf9name,
	coalesce(udf9.code,'') as udf9code,
	coalesce(udf10.name,'') as udf10name,
	coalesce(udf10.code,'') as udf10code
	from executionorder_b as b
	left join executionorder_h as h on b.hid = h.id
	left join issueresolutionform as irf on b.irfid = irf.id
	left join sysuser as handler on irf.handlerid = handler.id
	left join department as dept on h.deptid = dept.id
	left join epa as epa on b.epaid = epa.id
	left join risklevel as rl on b.risklevelid = rl.id
	left join sysuser as creator on irf.creatorid = creator.id
	left join sysuser as confirmer on irf.confirmerid = confirmer.id
	left join sysuser as issueowner on b.issueownerid = issueowner.id
	left join sysuser as executor on h.executorid = executor.id
	left join ept_h as ept_h on h.eptid = ept_h.id
	left join csa as csa on h.csaid = csa.id
	left join uda as udf1 on csa.udf1 = udf1.id
	left join uda as udf2 on csa.udf2 = udf2.id
	left join uda as udf3 on csa.udf3 = udf3.id
	left join uda as udf4 on csa.udf4 = udf4.id
	left join uda as udf5 on csa.udf5 = udf5.id
	left join uda as udf6 on csa.udf6 = udf6.id
	left join uda as udf7 on csa.udf7 = udf7.id
	left join uda as udf8 on csa.udf8 = udf8.id
	left join uda as udf9 on csa.udf9 = udf9.id
	left join uda as udf10 on csa.udf10 = udf10.id
	where (b.dr=0 and b.isissue = 1)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	repSql := build.String()
	// Retrieve IRF data from database
	irfRep, err := db.Query(repSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetIssueResolutionFormReport db.Query failed", zap.Error(err))
		return
	}
	defer irfRep.Close()

	// Extra data row by row
	for irfRep.Next() {
		var irf IssueResolutionFormReport
		err = irfRep.Scan(&irf.EOBID, &irf.EOHID, &irf.EOBillNumber, &irf.EORowNumber, &irf.EOBillDate,
			&irf.EOHDeptID, &irf.EOHDeptCode, &irf.EOHDeptName, &irf.CSAID, &irf.CSACode,
			&irf.CSAName, &irf.CSCID, &irf.ExecutorID, &irf.ExecutorCode, &irf.ExecutorName,
			&irf.EPAID, &irf.EPACode, &irf.EPAName, &irf.RLID, &irf.RLName,
			&irf.RLColor, &irf.ExecutionValue, &irf.ExecutionValueDIsp, &irf.EOBDescription, &irf.IsIssue,
			&irf.IsRectify, &irf.IsHandle, &irf.IssueOwnerID, &irf.IssueOwnerCode, &irf.IssueOwnerName,
			&irf.EOBStartTime, &irf.EOBEndTime, &irf.IsFinish, &irf.IRFID, &irf.IRFBillNumber,
			&irf.IRFBillDate, &irf.HandlerID, &irf.HandlerCode, &irf.HandlerName, &irf.IRFStartTime,
			&irf.IRFEndTime, &irf.IRFDescription, &irf.IRFStatus, &irf.CreatorID, &irf.CreatorCode,
			&irf.CreatorName, &irf.ConfirmerID, &irf.ConfirmerCode, &irf.ConfirmerName, &irf.Udf1Name,
			&irf.Udf1Code, &irf.Udf2Name, &irf.Udf2Code, &irf.Udf3Name, &irf.Udf3Code,
			&irf.Udf4Name, &irf.Udf4Code, &irf.Udf5Name, &irf.Udf5Code, &irf.Udf6Name,
			&irf.Udf6Code, &irf.Udf7Name, &irf.Udf7Code, &irf.Udf8Name, &irf.Udf8Code,
			&irf.Udf9Name, &irf.Udf9Code, &irf.Udf10Name, &irf.Udf10Name)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetIssueResolutionFormReport irfRep.Next()  irfRep.Scan failed", zap.Error(err))
			return
		}

		// Get Execution Order Row attachments
		irf.EORFiles, resStatus, err = GetEORowFiles(irf.EOBID)
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		// Get Issue Resolution Form attachments
		if irf.IRFID > 0 {
			irf.IRFFiles, resStatus, err = GetIRFFiles(irf.IRFID)
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		irfs = append(irfs, irf)
	}
	return
}
