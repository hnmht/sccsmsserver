package pg

import (
	"sccsmsserver/i18n"
	"sccsmsserver/setting"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Work Order Header struct
type WorkOrder struct {
	HID         int32          `db:"id" json:"id"`
	BillNumber  string         `db:"billnumber" json:"billNumber"`
	BillDate    time.Time      `db:"billdate" json:"billDate"`
	Department  SimpDept       `db:"deptid" json:"department"`
	Description string         `db:"description" json:"description"`
	Status      int16          `db:"status" json:"status"`
	WorkDate    time.Time      `db:"workdate" json:"workDate"`
	Body        []WorkOrderRow `json:"body"`
	CreateDate  time.Time      `db:"createtime" json:"createDate"`
	Creator     Person         `db:"creatorid" json:"creator"`
	ConfirmDate time.Time      `db:"confirmtime" json:"confirmDate"`
	Confirmer   Person         `db:"confirmerid" json:"confirmer"`
	ModifyDate  time.Time      `db:"modifytime" json:"modifyDate"`
	Modifier    Person         `db:"modifierid" json:"modifier"`
	Ts          time.Time      `db:"ts" json:"ts"`
	Dr          int16          `db:"dr" json:"dr"`
}

// Work Order Body Row struct
type WorkOrderRow struct {
	BID          int32            `db:"id" json:"id"`
	HID          int32            `db:"hid" json:"hid"`
	RowNumber    int32            `db:"rownumber" json:"rowNumber"`
	CSA          ConstructionSite `db:"csaid" json:"csa"`
	Executor     Person           `db:"executorid" json:"executor"`
	Description  string           `db:"description" json:"description"`
	EPT          EPT              `db:"eptid" json:"ept"`
	StartTime    time.Time        `db:"starttime" json:"startTime"`
	EndTime      time.Time        `db:"endtime" json:"endTime"`
	Status       int16            `db:"status" json:"status"`
	EOID         int32            `db:"eoid" json:"eoID"`
	EONumber     string           `db:"eonumber" json:"eoNumber"`
	CreateDate   time.Time        `db:"createtime" json:"createDate"`
	Creator      Person           `db:"creatorid" json:"creator"`
	ConfirmDate  time.Time        `db:"confirmtime" json:"confirmDate"`
	Confirmer    Person           `db:"confirmerid" json:"confirmer"`
	ModifyDate   time.Time        `db:"modifytime" json:"modifyDate"`
	Modifier     Person           `db:"modifierid" json:"modifier"`
	Ts           time.Time        `db:"ts" json:"ts"`
	Dr           int16            `db:"dr" json:"dr"`
	BillNumber   string           `json:"billNumber"`
	BillDate     string           `json:"billDate"`
	Department   SimpDept         `json:"department"`
	HDescription string           `json:"headerDescription"`
	WorkDate     string           `json:"workDate"`
}

// Get the list of Work Order to be executed
func GetWORefer(queryString string) (wors []WorkOrderRow, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	wors = make([]WorkOrderRow, 0)
	var build strings.Builder
	// Concatenate the SQL for inspection.
	build.WriteString(`select count(b.id) as rownumber
	from workorder_b as b
	left join workorder_h as h on b.hid = h.id
	left join ept_h as epth on b.eptid = epth.id
	where (b.dr=0 and h.dr=0 and b.status=1)`)
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
		zap.L().Error("GetWORefer db.QueryRow(checkSql) failed", zap.Error(err))
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
	build.WriteString(`select b.id,b.hid,b.rownumber,b.csaid,b.executorid,
	b.description as bdescription,b.eptid,epth.code,epth.name,epth.description as eptdescription,
	epth.allowaddrow,epth.allowdelrow,b.starttime,b.endtime,b.status,		
	b.eoid,b.eonumber,b.createtime,b.creatorid,b.confirmtime,
	b.confirmerid,b.modifytime,b.modifierid,b.ts,b.dr,
	h.billnumber,h.billdate,h.deptid,h.description as hdescription,h.workdate
	from workorder_b as b
	left join workorder_h as h on b.hid = h.id
	left join ept_h as epth on b.eptid = epth.id
	where (b.dr=0 and h.dr=0 and b.status=1)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	refSql := build.String()
	// Get Work Order List
	woRef, err := db.Query(refSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetWORefer db.Query failed", zap.Error(err))
		return
	}
	defer woRef.Close()
	// Get the Work Order row by row
	for woRef.Next() {
		var wor WorkOrderRow
		err = woRef.Scan(&wor.BID, &wor.HID, &wor.RowNumber, &wor.CSA.ID, &wor.Executor.ID,
			&wor.Description, &wor.EPT.HID, &wor.EPT.Code, &wor.EPT.Name, &wor.EPT.Description,
			&wor.EPT.AllowAddRow, &wor.EPT.AllowDelRow, &wor.StartTime, &wor.EndTime, &wor.Status,
			&wor.EOID, &wor.EONumber, &wor.CreateDate, &wor.Creator.ID, &wor.ConfirmDate,
			&wor.Confirmer.ID, &wor.ModifyDate, &wor.Modifier.ID, &wor.Ts, &wor.Dr,
			&wor.BillNumber, &wor.BillDate, &wor.Department.ID, &wor.HDescription, &wor.WorkDate)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetWORefer woRef.Next() woRef.Scan failed", zap.Error(err))
			return
		}
		// Get Construction Site detail
		if wor.CSA.ID > 0 {
			resStatus, err = wor.CSA.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Executor detail
		if wor.Executor.ID > 0 {
			resStatus, err = wor.Executor.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				zap.L().Warn("GetWORefer  wor.Executor.GetPersonInfoByID anomaly")
				return
			}
		}
		// Get Creator detail
		if wor.Creator.ID > 0 {
			resStatus, err = wor.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				zap.L().Warn("GetWORefer wor.Creator.GetPersonInfoByID anomaly")
				return
			}
		}
		// Get Confirmer detail
		if wor.Confirmer.ID > 0 {
			resStatus, err = wor.Confirmer.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				zap.L().Warn("GetWORefer   wor.Confirmer.GetPersonInfoByID anomaly")
				return
			}
		}
		// Get Modifier detail
		if wor.Modifier.ID > 0 {
			resStatus, err = wor.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				zap.L().Warn("GetWORefer  wor.Modifier.GetPersonInfoByID anomaly")
				return
			}
		}
		// Get Department detail
		if wor.Department.ID > 0 {
			resStatus, err = wor.Department.GetSimpDeptInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				zap.L().Warn("GetWORefer  wor.Department.GetSimpDeptInfoByID anomaly")
				return
			}
		}
		wors = append(wors, wor)
	}
	return
}

// Get Work Order List
func GetWOList(queryString string) (wos []WorkOrder, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	wos = make([]WorkOrder, 0)
	var build strings.Builder
	// Concatenate the SQL for inspection
	build.WriteString(`select count(workorder_h.id) as rownumber
	from workorder_h
	left join department on workorder_h.deptid = department.id
	left join sysuser as creator on workorder_h.creatorid = creator.id
	left join sysuser as modifier on workorder_h.modifierid = modifier.id
	where (workorder_h.dr = 0)`)
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
		zap.L().Error("GetWOList db.QueryRow(checkSql) failed", zap.Error(err))
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
	build.WriteString(`select workorder_h.id,workorder_h.billnumber,workorder_h.billdate,workorder_h.deptid,workorder_h.description,
	workorder_h.status,workorder_h.workdate,workorder_h.createtime,workorder_h.creatorid,workorder_h.confirmtime,
	workorder_h.confirmerid,workorder_h.modifytime,workorder_h.modifierid,workorder_h.dr,workorder_h.ts
	from workorder_h
	left join department on workorder_h.deptid = department.id
	left join sysuser as creator on workorder_h.creatorid = creator.id
	left join sysuser as modifier on workorder_h.modifierid = modifier.id
	where (workorder_h.dr = 0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	build.WriteString(`order by workorder_h.ts desc`)
	headSql := build.String()
	// Get Work Order List
	headRows, err := db.Query(headSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetWOList db.Query failed", zap.Error(err))
		return
	}
	defer headRows.Close()
	// Retrieve data row by row
	for headRows.Next() {
		var wo WorkOrder
		err = headRows.Scan(&wo.HID, &wo.BillNumber, &wo.BillDate, &wo.Department.ID, &wo.Description,
			&wo.Status, &wo.WorkDate, &wo.CreateDate, &wo.Creator.ID, &wo.ConfirmDate,
			&wo.Confirmer.ID, &wo.ModifyDate, &wo.Modifier.ID, &wo.Dr, &wo.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetWOList headRows.Next failed", zap.Error(err))
			return
		}
		// Get Department detail
		if wo.Department.ID > 0 {
			resStatus, err = wo.Department.GetSimpDeptInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Creator detail
		if wo.Creator.ID > 0 {
			resStatus, err = wo.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier detail
		if wo.Modifier.ID > 0 {
			resStatus, err = wo.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Confirmer detail
		if wo.Confirmer.ID > 0 {
			resStatus, err = wo.Confirmer.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		wos = append(wos, wo)
	}
	resStatus = i18n.StatusOK
	return
}

// Get Work Order details by ID
func (wo *WorkOrder) GetDetailByHID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the Work Order has been modified
	var rowNumber int32
	checkSql := `select count(id) as rownumber from workorder_h where id=$1 and dr=0`
	err = db.QueryRow(checkSql, wo.HID).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetDetailByHID db.QueryRow(checkSql) failed", zap.Error(err))
		return
	}

	if rowNumber < 1 {
		resStatus = i18n.StatusDataDeleted
		return
	}
	// Get Department details
	if wo.Department.ID > 0 {
		resStatus, err = wo.Department.GetSimpDeptInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Creator details
	if wo.Creator.ID > 0 {
		resStatus, err = wo.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Modifier details
	if wo.Modifier.ID > 0 {
		resStatus, err = wo.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Confirmer details
	if wo.Confirmer.ID > 0 {
		resStatus, err = wo.Confirmer.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	var bodyRowNumber int32
	// Get Body row details
	bodySql := `select id,hid,rownumber,csaid,executorid,
	description,eptid,starttime,endtime,status,
	eoid,createtime,creatorid,confirmtime,confirmerid,
	modifytime,modifierid,ts,dr
	from workorder_b
	where dr=0 and hid=$1 order by rownumber asc`
	bodyRows, err := db.Query(bodySql, wo.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.GetDetailByHID db.Query(bodySql) failed", zap.Error(err))
		return
	}
	defer bodyRows.Close()
	// Get body row by row
	for bodyRows.Next() {
		bodyRowNumber++
		var wor WorkOrderRow
		err = bodyRows.Scan(&wor.BID, &wor.HID, &wor.RowNumber, &wor.CSA.ID, &wor.Executor.ID,
			&wor.Description, &wor.EPT.HID, &wor.StartTime, &wor.EndTime, &wor.Status,
			&wor.EOID, &wor.CreateDate, &wor.Creator.ID, &wo.ConfirmDate, &wo.Confirmer.ID,
			&wor.ModifyDate, &wor.Modifier.ID, &wor.Ts, &wor.Dr)
		// Get Construction Site Achive details
		if wor.CSA.ID > 0 {
			resStatus, err = wor.CSA.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Executor deatails
		if wor.Executor.ID > 0 {
			resStatus, err = wor.Executor.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Execution Project Template Details.
		// in order to save resources, change it to retrieve from the frontend cache
		// Get Creator details
		if wor.Creator.ID > 0 {
			resStatus, err = wor.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier details
		if wor.Confirmer.ID > 0 {
			resStatus, err = wor.Confirmer.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Confirmer details
		if wor.Confirmer.ID > 0 {
			resStatus, err = wor.Confirmer.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		wo.Body = append(wo.Body, wor)
	}
	return
}

// Add Work Order
func (wo *WorkOrder) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check the number of rows in the body,
	// It cannot be zero
	if len(wo.Body) == 0 {
		resStatus = i18n.StatusVoucherNoBody
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.Add db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Get the latest Serial Number
	billNo, resStatus, err := GetLatestSerialNo(tx, "WO", wo.BillDate.Format("20060102"))
	if resStatus != i18n.StatusOK || err != nil {
		tx.Rollback()
		return
	}
	wo.BillNumber = billNo
	// Write the header content to the database
	headSql := `insert into workorder_h(billnumber,billdate,deptid,description,status,
	workdate,creatorid) values($1,$2,$3,$4,$5,$6,$7) returning id`
	err = tx.QueryRow(headSql, wo.BillNumber, wo.BillDate, wo.Department.ID, wo.Description, wo.Status,
		wo.WorkDate, wo.Creator.ID).Scan(&wo.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.Add tx.QueryRow Failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Prepare Write the body content to the database
	bodySql := `insert into workorder_b(hid,rownumber,csaid,executorid,description,
	eptid,starttime,endtime,status,creatorid)
	values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) returning id`
	bodyStmt, err := tx.Prepare(bodySql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.Add tx.Prepare(bodySql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer bodyStmt.Close()
	// Write data to database row by row
	for _, row := range wo.Body {
		err = bodyStmt.QueryRow(wo.HID, row.RowNumber, row.CSA.ID, row.Executor.ID, row.Description,
			row.EPT.HID, row.StartTime, row.EndTime, row.Status, wo.Creator.ID).Scan(&row.BID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("WorkOrder.Add bodyStmt.QueryRow falied", zap.Error(err))
			tx.Rollback()
			return
		}
	}
	return
}

// Edit the  Work Order content
func (wo *WorkOrder) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check the number of rows in the body, it cannot be zero
	if len(wo.Body) == 0 {
		resStatus = i18n.StatusVoucherNoBody
		return
	}
	// Check if the creator and the modifier are the same person
	if wo.Creator.ID != wo.Modifier.ID {
		resStatus = i18n.StatusVoucherOnlyCreateEdit
		return
	}

	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.Edit db.Begin() failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Update the Work Order header content in the database
	editHeadSql := `update workorder_h set billdate=$1,deptid=$2,description=$3,status=$4,workdate=$5,
	modifytime=current_timestamp,modifierid=$6,ts=current_timestamp
	where id=$7 and dr=0 and status=0 and ts=$8`
	editHeadRes, err := tx.Exec(editHeadSql, wo.BillDate, wo.Department.ID, wo.Description, wo.Status, wo.WorkDate,
		wo.Modifier.ID,
		wo.HID, wo.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.Edit tx.Exec(editHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by the SQL statement
	headUpdateNumber, err := editHeadRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.Edit EditHeadRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// If the number of effected rows is less than 1,
	// it indicates that someone else has modified the data.
	if headUpdateNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Prepare to update or insert body rows
	updateRowSql := `update workorder_b set hid=$1,rownumber=$2,csaid=$3,executorid=$4,description=$5,
	eptid=$6,starttime=$7,endtime=$8,status=$9,modifytime=current_timestamp,modifierid=$10,
	ts=current_timestamp,dr=$11  
	where id=$12 and ts=$13 and status=0 and dr=0 and eoid=0`
	addRowSql := `insert into workorder_b(hid,rownumber,csaid,executorid,description,
	eptid,starttime,endtime,status,creatorid,modifierid)
	values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) returning id`
	// Prepare to update row
	updateRowStmt, err := tx.Prepare(updateRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.Edit tx.Prepare(updateRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer updateRowStmt.Close()
	// Prepare to add row
	addRowStmt, err := tx.Prepare(addRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.Edit tx.Prepare(addRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer addRowStmt.Close()
	// Write body rows into the database
	for _, row := range wo.Body {
		if row.BID == 0 { // If the row.BID is o, it means the row is new
			err = addRowStmt.QueryRow(wo.HID, row.RowNumber, row.CSA.ID, row.Executor.ID, row.Description,
				row.EPT.HID, row.StartTime, row.EndTime, row.Status, wo.Modifier.ID, wo.Modifier.ID).Scan(&row.BID)
			if err != nil {
				zap.L().Error("WorkOrder.Edit addRowStmt.QueryRow failed", zap.Error(err))
				resStatus = i18n.StatusInternalError
				tx.Rollback()
				return
			}
		} else { //If bid is non-zero, it means the row needs to be modified
			updateRowRes, errUpdate := updateRowStmt.Exec(wo.HID, row.RowNumber, row.CSA.ID, row.Executor.ID, row.Description,
				row.EPT.HID, row.StartTime, row.EndTime, row.Status, wo.Modifier.ID,
				row.Dr, row.BID, row.Ts)
			if errUpdate != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("WorkOrder.Edit updateRowStmt.Exec failed", zap.Error(errUpdate))
				tx.Rollback()
				return resStatus, errUpdate
			}
			// Check the number of rows affected by the SQL statement
			updateRowNumber, errUpdateEffect := updateRowRes.RowsAffected()
			if errUpdateEffect != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("WorkOrder.Edit updateRowRes.RowsAffected failed", zap.Error(errUpdateEffect))
				tx.Rollback()
				return resStatus, errUpdateEffect
			}
			// If the number of affected rows less than one,
			// it means that some one else has modified the row.
			if updateRowNumber < 1 {
				resStatus = i18n.StatusOtherEdit
				tx.Rollback()
				return
			}
		}
	}

	return
}

// Delete Work Order
func (wo *WorkOrder) Delete(modifyUserId int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get Work Order details by HID
	resStatus, err = wo.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check the voucher status,  only voucher in the free state can be deleted
	if wo.Status != 0 { //
		resStatus = i18n.StatusVoucherNoFree
		return
	}
	// Check if the creator and the modifier are the same person
	// only creator can  delete the voucher.
	if wo.Creator.ID != modifyUserId {
		resStatus = i18n.StatusVoucherOnlyCreateEdit
		return
	}
	// Begin a database the transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.Delete db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Update the header deletion flag in the workorder_h table
	delHeadSql := `update workorder_h set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	delHeadRes, err := tx.Exec(delHeadSql, modifyUserId, wo.HID, wo.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.Delete tx.Exec(delHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by the SQL statement
	delHeadNumber, err := delHeadRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.Delete delHeadRes.RowsAffected() failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if delHeadNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}

	// Update the body row deletion flag in the workorder_b table
	delRowSql := `update workorder_b set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	// SQL statement preparation
	delRowStmt, err := tx.Prepare(delRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.Delete tx.Prepare(delRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delRowStmt.Close()
	// Write the deletion flag row by row
	for _, row := range wo.Body {
		// Check the Work Order Body Row status
		if row.Status != 0 {
			resStatus = i18n.StatusVoucherNoFree
			tx.Rollback()
			return
		}
		delRowRes, errDelRow := delRowStmt.Exec(modifyUserId, row.BID, row.Ts)
		if errDelRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("WorkOrder.Delete delRowStmt.Exec() failed", zap.Error(errDelRow))
			tx.Rollback()
			return resStatus, errDelRow
		}
		// Check the number of rows affected by SQL statement
		delRowNumber, errDelRow := delRowRes.RowsAffected()
		if errDelRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("WorkOrder.Delete delRowRes.RowsAffected() failed", zap.Error(errDelRow))
			tx.Rollback()
			return resStatus, errDelRow
		}
		if delRowNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}
	}

	return
}

// Batch Delete Work Order
func DeleteWOs(wos *[]WorkOrder, modifyUserID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteWOs db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// SQL statement preparation
	// SQL for updating the header deletion flag
	delHeadSql := `update workorder_h set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	delHeadStmt, err := tx.Prepare(delHeadSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteWOs tx.Prepare(delHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delHeadStmt.Close()
	// SQL for updating the body row deletion flag
	delRowSql := `update workorder_b set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	delRowStmt, err := tx.Prepare(delRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteWOs tx.Prepare(delRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delRowStmt.Close()
	// Write deletion flag item by item
	for _, wo := range *wos {
		// Get Work Order details
		resStatus, err = wo.GetDetailByHID()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
		// Check the Work Order status
		if wo.Status != 0 {
			resStatus = i18n.StatusVoucherNoFree
			tx.Rollback()
			return
		}
		// Check the creator and the modifier are the same person
		if wo.Creator.ID != modifyUserID {
			resStatus = i18n.StatusVoucherOnlyCreateEdit
			tx.Rollback()
			return
		}

		// Write the deletion flag into the workorder_h table
		delHeadRes, errDelHead := delHeadStmt.Exec(modifyUserID, wo.HID, wo.Ts)
		if errDelHead != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("DeleteWOs delHeadStmt.Exec failed", zap.Error(errDelHead))
			tx.Rollback()
			return resStatus, errDelHead
		}
		// Check the number of rows effected by SQL statement
		delHeadNumber, errCheck := delHeadRes.RowsAffected()
		if errCheck != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("DeleteWOs delHeadRes.RowsAffected() failed", zap.Error(err))
			tx.Rollback()
			return resStatus, errCheck
		}
		if delHeadNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}

		// Write the deletion flag to the workorder_b table
		for _, row := range wo.Body {
			// Check the row status
			if wo.Status != 0 {
				resStatus = i18n.StatusVoucherNoFree
				tx.Rollback()
				return
			}
			delRowRes, errDelRow := delRowStmt.Exec(modifyUserID, row.BID, row.Ts)
			if errDelRow != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("DeleteWOs delRowStmt.Exec() failed", zap.Error(errDelRow))
				tx.Rollback()
				return resStatus, errDelRow
			}
			// Check the number of rows affected by SQL statement
			delRowNumber, errDelRow := delRowRes.RowsAffected()
			if errDelRow != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("DeleteWOs delRowRes.RowsAffected() failed", zap.Error(errDelRow))
				tx.Rollback()
				return resStatus, errDelRow
			}
			if delRowNumber < 1 {
				resStatus = i18n.StatusOtherEdit
				tx.Rollback()
				return
			}
		}
	}
	return
}

// Confirm Work Order
func (wo *WorkOrder) Confirm(confirmUserID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get Work Order details
	resStatus, err = wo.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check the Work Order status
	// the WO status must be free
	if wo.Status != 0 {
		resStatus = i18n.StatusVoucherNoFree
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.Confirm db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Update the header confirmation flag
	confirmHeadSql := `update workorder_h set status=1,confirmtime=current_timestamp,confirmerid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and status=0 and ts=$3`
	headRes, err := tx.Exec(confirmHeadSql, confirmUserID, wo.HID, wo.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorKOrder.Confirm tx.Exec(confirmHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	confirmHeadNumber, err := headRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.Confirm headRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if confirmHeadNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}

	// Prepare SQL statement for update the body row confirmation flag
	confirmRowSql := `update workorder_b set status=1,confirmtime=current_timestamp,confirmerid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and status=0 and ts=$3`
	rowStmt, err := tx.Prepare(confirmRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.Confirm tx.Prepare(confirmRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer rowStmt.Close()

	// Update the body confirmation flag row by row.
	for _, row := range wo.Body {
		// Check the body row status
		if row.Status != 0 {
			resStatus = i18n.StatusVoucherNoFree
			tx.Rollback()
			return
		}
		confirmRowRes, errConfirmRow := rowStmt.Exec(confirmUserID, row.BID, row.Ts)
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("WorkOrder.Confirm rowStmt.Exec failed", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}

		confirmRowNumber, errConfirmRow := confirmRowRes.RowsAffected()
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("WorkOrder.Confirm confirmRowRes.RowsAffected failed", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}
		if confirmRowNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}
	}

	return
}

// Unconfirm Work Order
func (wo *WorkOrder) UnConfirm(confirmUserID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get the Work Order details
	resStatus, err = wo.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check the Work Order status,
	// Only confirmed voucher can be unconfirmed
	if wo.Status != 1 {
		resStatus = i18n.StatusVoucherNoConfirm
		return
	}
	// Check if the confirmer and the uncomfirmer are the same person
	if wo.Confirmer.ID != confirmUserID {
		resStatus = i18n.StatusVoucherCancelConfirmSelf
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.UnConfirm db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Update the header confirmation flag
	confirmHeadSql := `update workorder_h set status=0,confirmerid=0,confirmtime=to_timestamp(0),ts=current_timestamp 
	where id=$1 and dr=0 and status=1 and ts=$2`
	headRes, err := tx.Exec(confirmHeadSql, wo.HID, wo.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorKOrder.UnConfirm tx.Exec(confirmHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	confirmHeadNumber, err := headRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.UnConfirm headRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if confirmHeadNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}

	// Prepare the SQL statement for update the body row confirmation flag
	confirmRowSql := `update workorder_b set status=0,confirmerid=0,confirmtime=to_timestamp(0),ts=current_timestamp 
	where id=$1 and dr=0 and status=1 and ts=$2`
	rowStmt, err := tx.Prepare(confirmRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrder.UnConfirm tx.Prepare(confirmRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer rowStmt.Close()

	// Write confirmation flag to the workorder_b row by row
	for _, row := range wo.Body {
		// Check the row status
		if row.Status != 1 {
			resStatus = i18n.StatusVoucherNoConfirm
			tx.Rollback()
			return
		}
		confirmRowRes, errConfirmRow := rowStmt.Exec(row.BID, row.Ts)
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("WorkOrder.UnConfirm rowStmt.Exec failed", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}

		confirmRowNumber, errConfirmRow := confirmRowRes.RowsAffected()
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("WorkOrder.UnConfirm confirmRowRes.RowsAffected failed", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}
		if confirmRowNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}
	}

	return
}

// Execute the Work Order
func (wor *WorkOrderRow) Execute() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrderRow.Execute db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Update the Work Order body row
	rowSql := `update workorder_b set status=2,eoid = $1,eonumber=$2,ts=current_timestamp 
	where id=$3 and hid=$4 and ts=$5 and dr=0 and status=1`
	rowUpdateRes, err := tx.Exec(rowSql, wor.EOID, wor.EONumber, wor.BID, wor.HID, wor.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrderRow.Execute tx.Exec(rowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affect by the SQL statement
	rowUpdateNumber, err := rowUpdateRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrderRow.Execute rowUpdateRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if rowUpdateNumber < 1 {
		resStatus = i18n.StatusWOOtherEdit
		zap.L().Info("WorkOrderRow.Execute row OtherEdit")
		tx.Rollback()
		return
	}

	// Check the Work Order header status
	var headStatus int16
	headStatusSql := `select status from workorder_h where id=$1 limit 1`
	err = tx.QueryRow(headStatusSql, wor.HID).Scan(&headStatus)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrderRow.Execute tx.QueryRow(headStatusSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// If work order header status is confirmed, then change it to executing
	if headStatus == 1 {
		headSql := `update workorder_h set status=2,ts=current_timestamp where id=$1 and dr=0 and status=1`
		headUpdateRes, headErr := tx.Exec(headSql, wor.HID)
		if headErr != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("WorkOrderRow.Execute tx.Exec(headSql) failed", zap.Error(headErr))
			tx.Rollback()
			return resStatus, headErr
		}
		// Check the number of rows affected by SQL statement
		headUpdateNumber, effectErr := headUpdateRes.RowsAffected()
		if effectErr != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("WorkOrderRow.Execute headUpdateRes.RowsAffected() failed", zap.Error(effectErr))
			tx.Rollback()
			return resStatus, effectErr
		}
		if headUpdateNumber < 1 {
			resStatus = i18n.StatusWOOtherEdit
			zap.L().Info("WorkOrderRow.Execute header OtherEdit")
			tx.Rollback()
			return
		}
	}

	return
}

// Cancel Execute Work Order
func (wor *WorkOrderRow) CancelExecute() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrderRow.CancelExecute db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Update the Work Order body Row
	rowSql := `update workorder_b set status=1,eoid = 0,eonumber='',ts=current_timestamp 
	where id=$1 and hid=$2 and dr=0 and status=2`
	rowUpdateRes, err := tx.Exec(rowSql, wor.BID, wor.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrderRow.CancelExecute tx.Exec(rowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	rowUpdateNumber, err := rowUpdateRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrderRow.CancelExecute rowUpdateRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if rowUpdateNumber < 1 {
		resStatus = i18n.StatusWOOtherEdit
		zap.L().Info("WorkOrderRow.CancelExecute row OtherEdit")
		tx.Rollback()
		return
	}

	// Check if the are any rows, excluding the current one, with a status greater than or equal executing status
	var checkNumber int32
	checkSql := `select count(id) from workorder_b where dr = 0 and id <> $1 and hid=$2 and dr=0 and status >= 2`
	err = tx.QueryRow(checkSql, wor.BID, wor.HID).Scan(&checkNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrderRow.CancelExecute tx.QueryRow(checkSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}

	// If there are no rows with a status greater than "confirmed", then modify the header status to "confirmed"
	if checkNumber == 0 {
		headSql := `update workorder_h set status=1,ts=current_timestamp where id=$1 and dr=0 and status=2`
		headUpdateRes, headErr := tx.Exec(headSql, wor.HID)
		if headErr != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("WorkOrderRow.CancelExecute tx.Exec(headSql) failed", zap.Error(headErr))
			tx.Rollback()
			return resStatus, headErr
		}
		// Check the number of rows affected by SQL statement
		headUpdateNumber, effectErr := headUpdateRes.RowsAffected()
		if effectErr != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("WorkOrderRow.CancelExecute headUpdateRes.RowsAffected() failed", zap.Error(effectErr))
			tx.Rollback()
			return resStatus, effectErr
		}

		if headUpdateNumber < 1 {
			resStatus = i18n.StatusWOOtherEdit
			zap.L().Info("WorkOrderRow.CancelExecute header OtherEdit")
			tx.Rollback()
			return
		}

	}
	return
}

// Completed Work Order
func (wor *WorkOrderRow) Complete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrderRow.Complete db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Update the Work Order body row status to "completed"
	rowSql := `update workorder_b set status=3,ts=current_timestamp where id=$1 and hid=$2 and dr=0 and status=2`
	rowUpdateRes, err := tx.Exec(rowSql, wor.BID, wor.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrderRow.Complete tx.Exec(rowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	rowUpdateNumber, err := rowUpdateRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrderRow.Complete rowUpdateRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if rowUpdateNumber < 1 {
		resStatus = i18n.StatusWOOtherEdit
		zap.L().Info("WorkOrderRow.Complete row OtherEdit")
		tx.Rollback()
		return
	}
	// Check if this Work Order has any other rows in an uncompleted status
	var checkNumber int32
	checkSql := `select count(id) from workorder_b where dr = 0 and id <> $1 and hid=$2 and dr=0 and status < 3`
	err = tx.QueryRow(checkSql, wor.BID, wor.HID).Scan(&checkNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrderRow.Complete tx.QueryRow(checkSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// If all row status are "completed", modify the header status to "completed"
	if checkNumber == 0 {
		headSql := `update workorder_h set status=3,ts=current_timestamp where id=$1 and dr=0 and status=2`
		headUpdateRes, headErr := tx.Exec(headSql, wor.HID)
		if headErr != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("WorkOrderRow.Complete tx.Exec(headSql) failed", zap.Error(headErr))
			tx.Rollback()
			return resStatus, headErr
		}
		// Check the number of rows affected by SQL statement
		headUpdateNumber, effectErr := headUpdateRes.RowsAffected()
		if effectErr != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("WorkOrderRow.Complete headUpdateRes.RowsAffected() failed", zap.Error(effectErr))
			tx.Rollback()
			return resStatus, effectErr
		}
		if headUpdateNumber < 1 {
			resStatus = i18n.StatusWOOtherEdit
			zap.L().Info("WorkOrderRow.Complete header OtherEdit")
			tx.Rollback()
			return
		}
	}

	return
}

// Cancel the Work Order completed status
func (wor *WorkOrderRow) CancelComplete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrderRow.CancelComplete db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Update the Work Order body row status
	rowSql := `update workorder_b set status=2,ts=current_timestamp where id=$1 and hid=$2 and dr=0 and status=3`
	rowUpdateRes, err := tx.Exec(rowSql, wor.BID, wor.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrderRow.CancelComplete tx.Exec(rowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	rowUpdateNumber, err := rowUpdateRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrderRow.CancelComplete rowUpdateRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if rowUpdateNumber < 1 {
		resStatus = i18n.StatusWOOtherEdit
		zap.L().Info("WorkOrderRow.CancelComplete row OtherEdit")
		tx.Rollback()
	}

	// Check the Work Order Header status
	var headStatus int16
	headStatusSql := `select status from workorder_h where id=$1 and dr=0 limit 1`
	err = tx.QueryRow(headStatusSql, wor.HID).Scan(&headStatus)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("WorkOrderRow.CancelComplete tx.QueryRow(headStatusSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}

	// if the Work Order Header status is "completed", then chang it to "executing"
	if headStatus == 3 {
		headSql := `update workorder_h set status=2,ts=current_timestamp where id=$1 and dr=0 and status=3`
		headUpdateRes, headErr := tx.Exec(headSql, wor.HID)
		if headErr != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("WorkOrderRow.CancelComplete tx.Exec(headSql) failed", zap.Error(headErr))
			tx.Rollback()
			return resStatus, headErr
		}
		// Check the number of rows affected by SQL statement
		headUpdateNumber, effectErr := headUpdateRes.RowsAffected()
		if effectErr != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("WorkOrderRow.CancelComplete headUpdateRes.RowsAffected() failed", zap.Error(effectErr))
			tx.Rollback()
			return resStatus, effectErr
		}

		if headUpdateNumber < 1 {
			resStatus = i18n.StatusWOOtherEdit
			zap.L().Info("WorkOrderRow.CancelComplete header OtherEdit")
			tx.Rollback()
			return
		}
	}
	return
}
