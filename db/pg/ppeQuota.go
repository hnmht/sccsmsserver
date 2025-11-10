package pg

import (
	"sccsmsserver/i18n"
	"sccsmsserver/setting"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Personal Protective Equipment Quota struct
type PPEQuota struct {
	HID         int32         `db:"id" json:"id"`
	BillDate    time.Time     `db:"billdate" json:"billDate"`
	Position    Position      `db:"positionid" json:"position"`
	Period      string        `db:"period" json:"period"`
	Description string        `db:"description" json:"description"`
	Body        []PPEQuotaRow `json:"body"`
	Status      int16         `db:"status" json:"status"`
	CreateDate  time.Time     `db:"createtime" json:"createDate"`
	Creator     Person        `db:"creatorid" json:"creator"`
	ConfirmDate time.Time     `db:"confirmtime" json:"confirmDate"`
	Confirmer   Person        `db:"confirmerid" json:"confirmer"`
	ModifyDate  time.Time     `db:"modifytime" json:"modifyDate"`
	Modifier    Person        `db:"modifierid" json:"modifier"`
	Ts          time.Time     `db:"ts" json:"ts"`
	Dr          int16         `db:"dr" json:"dr"`
}

// Get Position's Personal Protective Equipment Quota Params
type PPEPositionsParams struct {
	Period    string     `json:"period"`    //周期
	Positions []Position `json:"positions"` //岗位列表
}

// Personal Protective Equipment Quota Row struct
type PPEQuotaRow struct {
	BID         int32     `db:"id" json:"id"`
	HID         int32     `db:"hid" json:"hid"`
	RowNumber   int32     `db:"rownumber" json:"rowNumber"`
	PPE         PPE       `db:"ppeid" json:"ppe"`
	Quantity    float64   `db:"quantity" json:"quantity"`
	Description string    `db:"description" json:"description"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Creator     Person    `db:"creatorid" json:"creator"`
	ConfirmDate time.Time `db:"confirmtime" json:"confirmDate"`
	Confirmer   Person    `db:"confirmerid" json:"confirmer"`
	ModifyDate  time.Time `db:"modifytime" json:"modifyDate"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

// Add Personal Protective Equipment Quota
func (pq *PPEQuota) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check the number of rows in the body, it cannot be zero
	if len(pq.Body) == 0 {
		resStatus = i18n.StatusVoucherNoBody
		return
	}
	// Check if the Position Quota Exist
	resStatus, err = pq.CheckExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Add db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Insert  the header content to the ppequota_h table
	headSql := `insert into ppequotas_h(billdate,positionid,period,description,status,
		creatorid) 
		values($1,$2,$3,$4,$5,$6) returning id`
	err = tx.QueryRow(headSql, pq.BillDate, pq.Position.ID, pq.Period, pq.Description, pq.Status,
		pq.Creator.ID).Scan(&pq.HID)

	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Add tx.QueryRow Failed", zap.Error(err))
		tx.Rollback()
		return
	}

	// Prepare insert rows to the ppequotas_b table
	bodySql := `insert into ppequotas_b(hid,rownumber,ppeid,quantity,description,
		status,creatorid)
		values($1,$2,$3,$4,$5,$6,$7) returning id`
	bodyStmt, err := tx.Prepare(bodySql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Add tx.Prepare(bodySql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer bodyStmt.Close()
	// Insert data row by row
	for _, row := range pq.Body {
		err = bodyStmt.QueryRow(pq.HID, row.RowNumber, row.PPE.ID, row.Quantity, row.Description,
			row.Status, pq.Creator.ID).Scan(&row.BID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEQuota.Add bodyStmt.QueryRow falied", zap.Error(err))
			tx.Rollback()
			return
		}
	}
	return
}

// Check if a PPE Position Quota for the same period
func (pq *PPEQuota) CheckExist() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var count int32
	sqlStr := `select count(id) from ppequotas_h where dr=0 and positionid=$1 and period=$2 and id<>$3`
	err = db.QueryRow(sqlStr, pq.Position.ID, pq.Period, pq.HID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.CheckExist query failed", zap.Error(err))
		return
	}
	if count > 0 {
		resStatus = i18n.StatusPQExist
		return
	}
	return
}

// Get Personal Protective Equipment Quota List
func GetPQList(queryString string) (pqs []PPEQuota, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	pqs = make([]PPEQuota, 0)
	var build strings.Builder
	// Concatenate the SQL for inspection
	build.WriteString(`select count(h.id) as rownumber
	from ppequotas_h h
	left join position on h.positionid=position.id
	left join sysuser as creator on h.creatorid=creator.id
	left join sysuser as modifier on h.modifierid=modifier.id
	where (h.dr = 0)`)
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
		zap.L().Error("GetPQList db.QueryRow(checkSql) failed", zap.Error(err))
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
	build.WriteString(`select h.billdate,h.id,h.positionid, h.period,h.description,
	h.status,h.createtime,h.creatorid,h.confirmtime,h.confirmerid,
	h.modifytime,h.modifierid,h.dr,h.ts 
	from ppequotas_h h
	left join position on h.positionid=position.id
	left join sysuser as creator on h.creatorid=creator.id
	left join sysuser as modifier on h.modifierid=modifier.id
	where (h.dr = 0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	build.WriteString(`order by h.ts desc`)
	headSql := build.String()
	// Retrieve the PPEQuota list from database
	headRows, err := db.Query(headSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetPQList db.Query failed", zap.Error(err))
		return
	}
	defer headRows.Close()
	// Extract data row by row
	for headRows.Next() {
		var pq PPEQuota
		err = headRows.Scan(&pq.BillDate, &pq.HID, &pq.Position.ID, &pq.Period, &pq.Description,
			&pq.Status, &pq.CreateDate, &pq.Creator.ID, &pq.ConfirmDate, &pq.Confirmer.ID,
			&pq.ModifyDate, &pq.Modifier.ID, &pq.Dr, &pq.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetPQList headRows.Next failed", zap.Error(err))
			return
		}
		// Get Position details
		if pq.Position.ID > 0 {
			resStatus, err = pq.Position.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Creator details
		if pq.Creator.ID > 0 {
			resStatus, err = pq.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier details
		if pq.Modifier.ID > 0 {
			resStatus, err = pq.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Confirmer details
		if pq.Confirmer.ID > 0 {
			resStatus, err = pq.Confirmer.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		pqs = append(pqs, pq)
	}
	return
}

// Get Personal Protective Equipment by HID
func (pq *PPEQuota) GetDetailByHID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the bill has been modified
	var rowNumber int32
	checkSql := `select count(id) as rownumber from ppequotas_h where id=$1 and dr=0`
	err = db.QueryRow(checkSql, pq.HID).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetDetailByHID db.QueryRow(checkSql) failed", zap.Error(err))
		return
	}
	if rowNumber < 1 {
		resStatus = i18n.StatusDataDeleted
		return
	}
	// Get Position details
	if pq.Position.ID > 0 {
		resStatus, err = pq.Position.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Creator details
	if pq.Creator.ID > 0 {
		resStatus, err = pq.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Modifier details
	if pq.Modifier.ID > 0 {
		resStatus, err = pq.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Confirmer details
	if pq.Confirmer.ID > 0 {
		resStatus, err = pq.Confirmer.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	var bodyRowNumber int32
	// Get the body rows details
	bodySql := `select id,hid,rownumber,ppeid,quantity,
	description,status,createtime,creatorid,confirmtime,
	confirmerid,modifytime,modifierid,ts,dr
	from ppequotas_b
	where dr=0 and hid=$1 order by rownumber asc`
	bodyRows, err := db.Query(bodySql, pq.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.GetDetailByHID db.Query(bodySql) failed", zap.Error(err))
		return
	}
	defer bodyRows.Close()
	// Extract data row by row
	for bodyRows.Next() {
		bodyRowNumber++
		var pqr PPEQuotaRow
		err = bodyRows.Scan(&pqr.BID, &pqr.HID, &pqr.RowNumber, &pqr.PPE.ID, &pqr.Quantity,
			&pqr.Description, &pqr.Status, &pqr.CreateDate, &pqr.Creator.ID, &pq.ConfirmDate,
			&pq.Confirmer.ID, &pqr.ModifyDate, &pqr.Modifier.ID, &pqr.Ts, &pqr.Dr)
		// Get PPE details
		if pqr.PPE.ID > 0 {
			resStatus, err = pqr.PPE.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Creator details
		if pqr.Creator.ID > 0 {
			resStatus, err = pqr.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier details
		if pqr.Confirmer.ID > 0 {
			resStatus, err = pqr.Confirmer.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Confirmer details
		if pqr.Confirmer.ID > 0 {
			resStatus, err = pqr.Confirmer.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		pq.Body = append(pq.Body, pqr)
	}
	return
}

// Edit Personal Protective Equipment Quota
func (pq *PPEQuota) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Checke the number of rows in the body, it cannot not be zero
	if len(pq.Body) == 0 {
		resStatus = i18n.StatusVoucherNoBody
		return
	}
	// Check if the Modifier and Creator are the same person
	if pq.Creator.ID != pq.Modifier.ID {
		resStatus = i18n.StatusVoucherOnlyCreateEdit
		return
	}
	// Check if a PPE Position Quota for the same period
	resStatus, err = pq.CheckExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Edit db.Begin() failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Update the header content in the ppequotas_h table
	editHeadSql := `update ppequotas_h set billdate=$1, positionid=$2,period=$3,description=$4,status=$5,
	modifytime=current_timestamp,modifierid=$6,ts=current_timestamp
	where id=$7 and dr=0 and status=0 and ts=$8`
	editHeadRes, err := tx.Exec(editHeadSql, pq.BillDate, pq.Position.ID, pq.Period, pq.Description, pq.Status,
		pq.Modifier.ID,
		pq.HID, pq.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Edit tx.Exec(editHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	headUpdateNumber, err := editHeadRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Edit EditHeadRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if headUpdateNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Prepare update the body content in the ppequotas_b table
	updateRowSql := `update ppequotas_b set hid=$1,rownumber=$2,ppeid=$3,quantity=$4,description=$5,
	status=$6,modifytime=current_timestamp,modifierid=$7,ts=current_timestamp,dr=$8 
	where id=$9 and ts=$10 and status=0 and dr=0`
	updateRowStmt, err := tx.Prepare(updateRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Edit tx.Prepare(updateRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer updateRowStmt.Close()
	// Prepare add body rows in the ppequotas_b table
	addRowSql := `insert into ppequotas_b(hid,rownumber,ppeid,quantity,description,
	status,creatorid,modifierid)
	values($1,$2,$3,$4,$5,$6,$7,$8) returning id`
	addRowStmt, err := tx.Prepare(addRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Edit tx.Prepare(addRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer addRowStmt.Close()
	// Processthe body data row by row
	for _, row := range pq.Body {
		if row.BID == 0 { // If the BID is 0, it means the row is new
			err = addRowStmt.QueryRow(pq.HID, row.RowNumber, row.PPE.ID, row.Quantity, row.Description,
				row.Status, pq.Modifier.ID, pq.Modifier.ID).Scan(&row.BID)
			if err != nil {
				zap.L().Error("PPEQuota.Edit addRowStmt.QueryRow failed", zap.Error(err))
				resStatus = i18n.StatusInternalError
				tx.Rollback()
				return
			}
		} else { // If the BID is non-zero, it means the row need to be modified
			updateRowRes, errUpdate := updateRowStmt.Exec(pq.HID, row.RowNumber, row.PPE.ID, row.Quantity, row.Description,
				row.Status, pq.Modifier.ID, row.Dr,
				row.BID, row.Ts)
			if errUpdate != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("PPEQuota.Edit updateRowStmt.Exec failed", zap.Error(errUpdate))
				tx.Rollback()
				return resStatus, errUpdate
			}
			updateRowNumber, errUpdateEffect := updateRowRes.RowsAffected()
			if errUpdateEffect != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("PPEQuota.Edit updateRowRes.RowsAffected failed", zap.Error(errUpdateEffect))
				tx.Rollback()
				return resStatus, errUpdateEffect
			}
			if updateRowNumber < 1 {
				resStatus = i18n.StatusOtherEdit
				tx.Rollback()
				return
			}
		}
	}
	return
}

// Delete Personal Protective Equipment Quota
func (pq *PPEQuota) Delete(operatorID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get the PPEQuota Details
	resStatus, err = pq.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check the bill status
	if pq.Status != 0 { // Cannot delete if the status is not 0
		resStatus = i18n.StatusVoucherNoFree
		return
	}
	// Check if the Creator and Operator are the same person
	if pq.Creator.ID != operatorID {
		resStatus = i18n.StatusVoucherOnlyCreateEdit
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Delete db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Update the deletion flag in the ppequotas_h table
	delHeadSql := `update ppequotas_h set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	delHeadRes, err := tx.Exec(delHeadSql, operatorID, pq.HID, pq.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Delete tx.Exec(delHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	delHeadNumber, err := delHeadRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Delete delHeadRes.RowsAffected() failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if delHeadNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Prepare to update the deletion flag in the ppequotas_b table
	delRowSql := `update ppequotas_b set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	delRowStmt, err := tx.Prepare(delRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Delete tx.Prepare(delRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delRowStmt.Close()
	// Write data row by row
	for _, row := range pq.Body {
		// Check the row status
		if row.Status != 0 { // Cannot delete if the row status value is 0
			resStatus = i18n.StatusVoucherNoFree
			tx.Rollback()
			return
		}
		delRowRes, errDelRow := delRowStmt.Exec(operatorID, row.BID, row.Ts)
		if errDelRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEQuota.Delete delRowStmt.Exec() failed", zap.Error(errDelRow))
			tx.Rollback()
			return resStatus, errDelRow
		}
		// Check the number of rows affected by SQL statement
		delRowNumber, errDelRow := delRowRes.RowsAffected()
		if errDelRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEQuota.Delete delRowRes.RowsAffected() failed", zap.Error(errDelRow))
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

// Confirm Personal Protective Equipment Quota
func (pq *PPEQuota) Confirm(operatorID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get the PPEQuota details
	resStatus, err = pq.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check the bill status
	if pq.Status != 0 { // Cannot confirm if the status is not 0
		resStatus = i18n.StatusVoucherNoFree
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Confirm db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Update the confirmation flag in the ppequotas_h table
	confirmHeadSql := `update ppequotas_h set status=1,confirmtime=current_timestamp,confirmerid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and status=0 and ts=$3`
	headRes, err := tx.Exec(confirmHeadSql, operatorID, pq.HID, pq.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Confirm tx.Exec(confirmHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	confirmHeadNumber, err := headRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Confirm headRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if confirmHeadNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Prepare to update the comfirmation flage in the ppequotas_b table
	confirmRowSql := `update ppequotas_b set status=1,confirmtime=current_timestamp,confirmerid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and status=0 and ts=$3`
	rowStmt, err := tx.Prepare(confirmRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Confirm tx.Prepare(confirmRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer rowStmt.Close()
	// Write body data row by row
	for _, row := range pq.Body {
		// Check the row status
		if row.Status != 0 { // Cannot confirm if the status is not 0
			resStatus = i18n.StatusVoucherNoFree
			tx.Rollback()
			return
		}
		confirmRowRes, errConfirmRow := rowStmt.Exec(operatorID, row.BID, row.Ts)
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEQuota.Confirm rowStmt.Exec failed", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}
		confirmRowNumber, errConfirmRow := confirmRowRes.RowsAffected()
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEQuota.Confirm confirmRowRes.RowsAffected failed", zap.Error(errConfirmRow))
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

// Unconfirm Personal Protective Equipment Quota
func (pq *PPEQuota) Unconfirm(operatorID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get the PPEQuota Details
	resStatus, err = pq.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check the bill status
	if pq.Status != 1 { // Cannot unconfirm if the status is not 1
		resStatus = i18n.StatusVoucherNoConfirm
		return
	}
	// Check the Confirmer and Operator are the same person
	if pq.Confirmer.ID != operatorID {
		resStatus = i18n.StatusVoucherCancelConfirmSelf
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Unconfirm db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Update the confirmation flag in the ppequotas_h table
	confirmHeadSql := `update ppequotas_h set status=0,confirmerid=0,confirmtime=to_timestamp(0),ts=current_timestamp 
	where id=$1 and dr=0 and status=1 and ts=$2`
	headRes, err := tx.Exec(confirmHeadSql, pq.HID, pq.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Unconfirm tx.Exec(confirmHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	confirmHeadNumber, err := headRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Unconfirm headRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if confirmHeadNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Prepare to update the confirmation flag in the ppequatos_b table
	confirmRowSql := `update ppequotas_b set status=0,confirmerid=0,confirmtime=to_timestamp(0),ts=current_timestamp 
	where id=$1 and dr=0 and status=1 and ts=$2`
	rowStmt, err := tx.Prepare(confirmRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEQuota.Unconfirm tx.Prepare(confirmRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer rowStmt.Close()

	// Wirte body data row by row
	for _, row := range pq.Body {
		// Check the row status
		if row.Status != 1 { // Cannot unconfirm if the status is not 1
			resStatus = i18n.StatusVoucherNoConfirm
			tx.Rollback()
			return
		}
		confirmRowRes, errConfirmRow := rowStmt.Exec(row.BID, row.Ts)
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEQuota.Unconfirm rowStmt.Exec failed", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}

		confirmRowNumber, errConfirmRow := confirmRowRes.RowsAffected()
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEQuota.Unconfirm confirmRowRes.RowsAffected failed", zap.Error(errConfirmRow))
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

// Retrieve the list of all position that have PPE Quotas within the same period
func (ppep *PPEPositionsParams) Get() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	ppep.Positions = make([]Position, 0)
	sqlStr := `select distinct positionid from ppequotas_h where dr=0 and status=1 and period=$1`
	rows, err := db.Query(sqlStr, ppep.Period)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEPositionsParams.Get db.Query failed:", zap.Error(err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var position Position
		err = rows.Scan(&position.ID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEPositionsParams.Get  rows.Scan failed:", zap.Error(err))
			return
		}
		if position.ID > 0 {
			resStatus, err = position.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		ppep.Positions = append(ppep.Positions, position)
	}
	return
}
