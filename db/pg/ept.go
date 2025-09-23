package pg

import (
	"database/sql"
	"sccsmsserver/i18n"
	"time"

	"go.uber.org/zap"
)

// Execution Project Template header struct
type EPT struct {
	HID         int32     `db:"id" json:"id"`
	Code        string    `db:"code" json:"code"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Status      int16     `db:"status" json:"status"`
	AllowAddRow int16     `db:"allowaddrow" json:"allowAddRow"`
	AllowDelRow int16     `db:"allowdelrow" json:"allowDelRow"`
	Body        []EPTRow  `json:"body"`
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Creator     Person    `db:"creatorid" json:"creator"`
	ModifyDate  time.Time `db:"modifytime" json:"modifyDate"`
	Modifier    Person    `db:"modifierid" json:"modifierid"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

// Eexcution Project Template Row struct
type EPTRow struct {
	BID              int32            `db:"id" json:"id"`
	HID              int32            `db:"hid" json:"hid"`
	RowNumber        int32            `db:"rownumber" json:"rowNumber"`
	EP               ExecutionProject `db:"epaid" json:"epa"`
	AllowDelRow      int16            `db:"allowdelrow" json:"allowDelRow"`
	Description      string           `db:"description" json:"description"`
	DefaultValue     string           `db:"defaultvalue" json:"defaultValue"`
	DefaultValueDisp string           `db:"defaultvaluedisp" json:"defaultValueDisp"`
	IsCheckError     int16            `db:"ischeckerror" json:"isCheckError"`
	ErrorValue       string           `db:"errorvalue" json:"errorValue"`
	ErrorValueDisp   string           `db:"errorvaluedisp" json:"errorValueDisp"`
	IsRequireFile    int16            `db:"isrequirefile" json:"isRequireFile"`
	IsOnsitePhoto    int16            `db:"isonsitephoto" json:"isOnsitePhoto"`
	RiskLevel        RiskLevel        `db:"risklevelid" json:"riskLevel"`
	CreateDate       time.Time        `db:"createtime" json:"createDate"`
	Creator          Person           `db:"creatorid" json:"creator"`
	ModifyDate       time.Time        `db:"modifytime" json:"modifyDate"`
	Modifier         Person           `db:"modifierid" json:"modifier"`
	Ts               time.Time        `db:"ts" json:"ts"`
	Dr               int16            `db:"dr" json:"dr"`
}

// Execution Project Template Front-end cache struct
type EPTCache struct {
	QueryTs      time.Time `json:"queryTs"`
	ResultNumber int32     `json:"resultNumber"`
	DelItems     []EPT     `json:"delTtems"`
	UpdateItems  []EPT     `json:"updateItems"`
	NewItems     []EPT     `json:"newItems"`
	ResultTs     time.Time `json:"resultTs"`
}

// Get the latest Execution Project Template for front-end cache
func (eptc *EPTCache) GetEPTCahce() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	eptc.DelItems = make([]EPT, 0)
	eptc.NewItems = make([]EPT, 0)
	eptc.UpdateItems = make([]EPT, 0)
	// Check if there is any data newer than QueryTs
	sqlStr := `select ts from ept_h where ts > $1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, eptc.QueryTs).Scan(&eptc.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows { // No data newer than QueryTs
			eptc.ResultNumber = 0
			eptc.ResultTs = eptc.QueryTs
			resStatus = i18n.StatusOK
			return
		}
		zap.L().Error("EPTCache.GetEPTCahce query latest ts failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Retrieve all data newer than QueryTs
	sqlStr = `select id,code,name,description,status,allowaddrow,allowdelrow,
	createtime,creatorid,modifytime,modifierid,dr,ts 
	from ept_h where ts > $1 order by ts desc`
	headers, err := db.Query(sqlStr, eptc.QueryTs)
	if err != nil {
		zap.L().Error("EPTCache.GetEPTCache db.Query failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer headers.Close()

	for headers.Next() {
		var ept EPT
		err = headers.Scan(&ept.HID, &ept.Code, &ept.Name, &ept.Description, &ept.Status, &ept.AllowAddRow, &ept.AllowDelRow,
			&ept.CreateDate, &ept.Creator.ID, &ept.ModifyDate, &ept.Modifier.ID, &ept.Dr, &ept.Ts)
		if err != nil {
			zap.L().Error("GetEPTCache headers.Next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get Creator info
		if ept.Creator.ID > 0 {
			resStatus, err = ept.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier info
		if ept.Modifier.ID > 0 {
			resStatus, err = ept.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Body info
		if ept.Dr == 0 {
			ept.Body, resStatus, err = getEPTBody(ept.HID)
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Classify data into New, Update, Delete based on CreateDate and Dr
		if ept.Dr == 0 {
			if ept.CreateDate.Before(eptc.QueryTs) || ept.CreateDate.Equal(eptc.QueryTs) {
				eptc.ResultNumber++
				eptc.UpdateItems = append(eptc.UpdateItems, ept)
			} else {
				eptc.ResultNumber++
				eptc.NewItems = append(eptc.NewItems, ept)
			}
		} else {
			if ept.CreateDate.Before(eptc.QueryTs) || ept.CreateDate.Equal(eptc.QueryTs) {
				eptc.ResultNumber++
				eptc.DelItems = append(eptc.DelItems, ept)
			}
		}
	}

	return
}

// Get Execution Project Template list
func GetEPTList() (eptList []EPT, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	eptList = make([]EPT, 0)
	// Get EPT header list from database
	headerSql := `select id,code,name,description,status,
	allowaddrow,allowdelrow,createtime,creatorid,modifytime,
	modifierid,dr,ts 
	from ept_h where dr=0 order by ts desc`
	headerRows, err := db.Query(headerSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEPTsList db.Query(headerSql) failed", zap.Error(err))
		return
	}
	defer headerRows.Close()

	// Fill in header information
	for headerRows.Next() {
		var ept EPT
		err = headerRows.Scan(&ept.HID, &ept.Code, &ept.Name, &ept.Description, &ept.Status,
			&ept.AllowAddRow, &ept.AllowDelRow, &ept.CreateDate, &ept.Creator.ID, &ept.ModifyDate,
			&ept.Modifier.ID, &ept.Dr, &ept.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetEPTsList headerRow.Next Scan EPT failed", zap.Error(err))
			return
		}
		// Fill in creator information
		if ept.Creator.ID > 0 {
			resStatus, err = ept.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Fill in modifierid information
		if ept.Modifier.ID > 0 {
			resStatus, err = ept.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Fill in body information
		ept.Body, resStatus, err = getEPTBody(ept.HID)
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		// Append to eptList
		eptList = append(eptList, ept)
	}

	return
}

// Get Execution Project Template body by HID
func getEPTBody(hid int32) (rows []EPTRow, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get EPT body row from database
	bodySql := `select id,hid,rownumber,epaid,allowdelrow,
	description,defaultvalue,defaultvaluedisp,ischeckerror,errorvalue,
	errorvaluedisp,isrequirefile,isonsitephoto,risklevelid,createtime,
	creatorid,modifytime,modifierid,dr,ts
	from ept_b where dr=0 and hid=$1 order by rownumber asc`
	var bodyRowNumber = 0
	bRows, err := db.Query(bodySql, hid)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetEPTBody db.Query(bodysql) failed", zap.Error(err))
		return
	}
	defer bRows.Close()
	// Fill in body information
	for bRows.Next() {
		bodyRowNumber++
		var eptRow EPTRow
		err = bRows.Scan(&eptRow.BID, &eptRow.HID, &eptRow.RowNumber, &eptRow.EP.ID, &eptRow.AllowDelRow,
			&eptRow.Description, &eptRow.DefaultValue, &eptRow.DefaultValueDisp, &eptRow.IsCheckError, &eptRow.ErrorValue,
			&eptRow.ErrorValueDisp, &eptRow.IsRequireFile, &eptRow.IsOnsitePhoto, &eptRow.RiskLevel.ID, &eptRow.CreateDate,
			&eptRow.Creator.ID, &eptRow.ModifyDate, &eptRow.Modifier.ID, &eptRow.Dr, &eptRow.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetEPTBody bRows.Next Scan EitRow failed", zap.Error(err))
			return
		}

		// Fill in execution project details
		if eptRow.EP.ID > 0 {
			resStatus, err = eptRow.EP.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Fill in risk level details
		if eptRow.RiskLevel.ID > 0 {
			resStatus, err = eptRow.RiskLevel.GetRLInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Fill in creator details
		if eptRow.Creator.ID > 0 {
			resStatus, err = eptRow.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Fill in modifier details
		if eptRow.Modifier.ID > 0 {
			resStatus, err = eptRow.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		rows = append(rows, eptRow)
	}

	// Check if there is no body row
	if bodyRowNumber == 0 {
		resStatus = i18n.StatusVoucherNoBody
		return
	}

	return
}

// Get Execution Project Template header by HID
func (ept *EPT) GetEPTHeaderByHid() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Retrieve header information from database
	headerSql := `select code,name,description,status,allowaddrow,
	allowdelrow,createtime,creatorid,modifytime,modifierid,
	dr,ts 
	from ept_h 
	where dr=0 and id=$1`
	err = db.QueryRow(headerSql, ept.HID).Scan(&ept.Code, &ept.Name, &ept.Description, &ept.Status, &ept.AllowAddRow,
		&ept.AllowDelRow, &ept.CreateDate, &ept.Creator.ID, &ept.ModifyDate, &ept.Modifier.ID,
		&ept.Dr, &ept.Ts)
	if err != nil {
		zap.L().Error("EPT.GetEPTHead db.QueryRow failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Fill in Creator info
	if ept.Creator.ID > 0 {
		resStatus, err = ept.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Fill in Modifier info
	if ept.Modifier.ID > 0 {
		resStatus, err = ept.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}

	return
}

// Check if the Execution Project Template code already exists
func (ept *EPT) CheckCodeExist() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var count int32
	sqlStr := `select count(id) from ept_h where dr=0 and code=$1 and id<>$2`
	err = db.QueryRow(sqlStr, ept.Code, ept.HID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPT.CheckCodeExist db.QueryRow failed", zap.Error(err))
		return
	}
	if count > 0 {
		resStatus = i18n.StatusEPTCodeExist
		return
	}

	return
}

// Add a new Execution Project Template
func (ept *EPT) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the code already exists
	resStatus, err = ept.CheckCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check if the body row count is zero, which is not allowed
	if len(ept.Body) == 0 {
		resStatus = i18n.StatusVoucherNoBody
		return
	}
	// Create a transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPT.Add db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Write header information into ept_h table
	addHeaderSql := `insert into ept_h(code,name,description,
	status,allowaddrow,allowdelrow,creatorid)
	values($1,$2,$3,$4,$5,$6,$7) returning id`
	err = tx.QueryRow(addHeaderSql, ept.Code, ept.Name, ept.Description, ept.Status,
		ept.AllowAddRow, ept.AllowDelRow, ept.Creator.ID).Scan(&ept.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPT tx.QueryRow(addHeaderSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Prepare to write body rows into ept_b table
	addBodySql := `insert into ept_b(hid,rownumber,epaid,allowdelrow,description,
		defaultvalue,defaultvaluedisp,ischeckerror,errorvalue,errorvaluedisp,
		isrequirefile,isonsitephoto,risklevelid,creatorid) 
		values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) 
		returning id`
	bodyStmt, err := tx.Prepare(addBodySql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPT addBodySql tx.Prepare failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer bodyStmt.Close()

	for _, row := range ept.Body {
		err = bodyStmt.QueryRow(ept.HID, row.RowNumber, row.EP.ID, row.AllowDelRow, row.Description,
			row.DefaultValue, row.DefaultValueDisp, row.IsCheckError, row.ErrorValue, row.ErrorValueDisp,
			row.IsRequireFile, row.IsOnsitePhoto, row.RiskLevel.ID, ept.Creator.ID).Scan(&row.BID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("EPT bodyStmt.QueryRow failed", zap.Error(err))
			tx.Rollback()
			return
		}
	}
	return
}

// Modify an existing Execution Project Template
func (ept *EPT) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the code already exists
	resStatus, err = ept.CheckCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check if the body row count is zero, which is not allowed
	if len(ept.Body) == 0 {
		resStatus = i18n.StatusVoucherNoBody
		return
	}
	// Create a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPT.Edit db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Modify header information in ept_h table
	editHeadSql := `update ept_h set code=$1,name=$2, description=$3,status=$4,
	allowaddrow=$5,allowdelrow=$6,modifytime=current_timestamp,modifierid=$7,ts=current_timestamp 
	where id=$8 and dr = 0 and ts=$9`
	editHeaderRes, err := tx.Exec(editHeadSql, ept.Code, ept.Name, ept.Description, ept.Status,
		ept.AllowAddRow, ept.AllowDelRow, ept.Modifier.ID,
		ept.HID, ept.Ts)
	if err != nil {
		zap.L().Error("EPT.Edit tx.Exec(editHeadSql) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	// Check the number of affected rows
	headerUpdatedNumber, err := editHeaderRes.RowsAffected()
	if err != nil {
		zap.L().Error("EPT.Edit editHeaderRes.RowsAffected failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	if headerUpdatedNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Prepare to update existing rows and add new rows in ept_b table
	updateRowSql := `update ept_b set hid=$1,rownumber=$2,epaid=$3,allowdelrow=$4,description=$5,
	defaultvalue=$6,defaultvaluedisp=$7,ischeckerror=$8,errorvalue=$9,errorvaluedisp=$10,
	isrequirefile=$11,isonsitephoto=$12,risklevelid=$13,modifierid=$14,modifytime=current_timestamp,
	ts=current_timestamp,dr=$15 
	where id=$16 and ts=$17 and dr=0`
	addRowSql := `insert into ept_b(hid,rownumber,epaid,allowdelrow,description,
		defaultvalue,defaultvaluedisp,ischeckerror,errorvalue,errorvaluedisp,
		isrequirefile,isonsitephoto,risklevelid,creatorid,modifierid) 
	values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) 
	returning id`
	// Prepare update statement
	updateRowStmt, err := tx.Prepare(updateRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPT.Edit updateRowStmt tx.Prepare failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer updateRowStmt.Close()
	// Prepare add statement
	addRowStmt, err := tx.Prepare(addRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPT.Edit addRowStmt tx.Prepare failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer addRowStmt.Close()

	// Update existing rows and add new rows in the body
	for _, eptRow := range ept.Body {
		if eptRow.BID == 0 { // if BID is 0, it's a new row, need to add
			errAddRow := addRowStmt.QueryRow(ept.HID, eptRow.RowNumber, eptRow.EP.ID, eptRow.AllowDelRow, eptRow.Description,
				eptRow.DefaultValue, eptRow.DefaultValueDisp, eptRow.IsCheckError, eptRow.ErrorValue, eptRow.ErrorValueDisp,
				eptRow.IsRequireFile, eptRow.IsOnsitePhoto, eptRow.RiskLevel.ID, ept.Modifier.ID, ept.Modifier.ID).Scan(&eptRow.BID)
			if errAddRow != nil {
				zap.L().Error("EPT.Edit addRowStmt.QueryRow failed", zap.Error(errAddRow))
				resStatus = i18n.StatusInternalError
				tx.Rollback()
				return resStatus, errAddRow
			}
		} else { // otherwise, it's an existing row, need to update
			updateRowRes, errUpdate := updateRowStmt.Exec(ept.HID, eptRow.RowNumber, eptRow.EP.ID, eptRow.AllowDelRow, eptRow.Description,
				eptRow.DefaultValue, eptRow.DefaultValueDisp, eptRow.IsCheckError, eptRow.ErrorValue, eptRow.ErrorValueDisp,
				eptRow.IsRequireFile, eptRow.IsOnsitePhoto, eptRow.RiskLevel.ID, ept.Modifier.ID,
				eptRow.Dr,
				eptRow.BID, eptRow.Ts)
			if errUpdate != nil {
				zap.L().Error("EPT.Edit updateRowStmt.QueryRow failed", zap.Error(errUpdate))
				resStatus = i18n.StatusInternalError
				tx.Rollback()
				return resStatus, errUpdate
			}
			// Check the number of affected rows
			updateRowNumber, errUpdateEffec := updateRowRes.RowsAffected()
			if errUpdateEffec != nil {
				zap.L().Error("EPT.Edit updateRowRes.RowsAffected failed", zap.Error(errUpdateEffec))
				resStatus = i18n.StatusInternalError
				tx.Rollback()
				return resStatus, errUpdate
			}
			// If no rows were affected, it means the row was edited by someone else
			if updateRowNumber < 1 {
				resStatus = i18n.StatusOtherEdit
				tx.Rollback()
				return
			}
		}
	}
	return
}

// Delete an existing Execution Project Template
func (ept *EPT) Delete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the Execution Project Template is referenced by other documents
	resStatus, err = ept.CheckUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Create a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPT.Delete db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Update the header to mark it as deleted
	delHeadSql := `update ept_h set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp
	where id=$2 and dr=0 and ts=$3`
	delHeadRes, err := tx.Exec(delHeadSql, ept.Modifier.ID, ept.HID, ept.Ts)
	if err != nil {
		zap.L().Error("EPT.Delete delHeadStmt.Exec failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	// Check the number of affected rows
	delHeadNum, err := delHeadRes.RowsAffected()
	if err != nil {
		zap.L().Error("EPT.Delete delHeadRes.RowsAffected failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	// If no rows were affected, it means the header was edited by someone else
	if delHeadNum < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Prepare to delete body rows
	delRowSql := `update ept_b set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp
	where id=$2 and dr=0 and ts=$3`
	delRowStmt, err := tx.Prepare(delRowSql)
	if err != nil {
		zap.L().Error("EPT.Delete tx.Prepare(delRowSql) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer delRowStmt.Close()
	// Write deletion for each body row
	for _, row := range ept.Body {
		delRowRes, errDelRow := delRowStmt.Exec(ept.Modifier.ID, row.BID, row.Ts)
		if errDelRow != nil {
			zap.L().Error("EPT.Delete delRowStmt.Exec failed", zap.Error(errDelRow))
			resStatus = i18n.StatusInternalError
			tx.Rollback()
			return resStatus, errDelRow
		}
		// Check the number of affected rows
		updateRowNumber, errUpdateEffc := delRowRes.RowsAffected()
		if errUpdateEffc != nil {
			zap.L().Error("EPT.Delete delRowRes.RowsAffected failed", zap.Error(errUpdateEffc))
			resStatus = i18n.StatusInternalError
			tx.Rollback()
			return resStatus, errUpdateEffc
		}
		// If no rows were affected, it means the row was edited by someone else
		if updateRowNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}
	}
	return
}

// Batch delete Execution Project Templates
func DeleteEPTs(eits *[]EPT, modifyUserID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteEPTs db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// SQL to delete header data
	delHeadSql := `update ept_h set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp
	where id=$2 and dr=0 and ts=$3`
	// SQL to delete body row data
	delRowSql := `update ept_b set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp
	where id=$2 and dr=0 and ts=$3`
	// Prepare to delete header
	delHeadStmt, err := tx.Prepare(delHeadSql)
	if err != nil {
		zap.L().Error("DeleteEPTs tx.Prepare(delHeadSql) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer delHeadStmt.Close()
	// Prepare to delete body rows
	delRowStmt, err := tx.Prepare(delRowSql)
	if err != nil {
		zap.L().Error("DeleteEPTs tx.Prepare(delRowSql) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer delRowStmt.Close()
	// Delete each Execution Project Template
	for _, ept := range *eits {
		// Check if the Execution Project Template is referenced by other documents
		resStatus, err = ept.CheckUsed()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
		// Delete header
		delHeadRes, errH := delHeadStmt.Exec(modifyUserID, ept.HID, ept.Ts)
		if errH != nil {
			zap.L().Error("DeleteEPTs delHeadStmt.Exec failed", zap.Error(errH))
			resStatus = i18n.StatusInternalError
			tx.Rollback()
			return resStatus, errH
		}
		// Check the number of affected rows
		delHeadNum, errH := delHeadRes.RowsAffected()
		if errH != nil {
			zap.L().Error("DeleteEPTs delHeadRes.RowsAffected failed", zap.Error(errH))
			resStatus = i18n.StatusInternalError
			tx.Rollback()
			return resStatus, errH
		}
		// If no rows were affected, it means the header was edited by someone else
		if delHeadNum < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}

		// Delete body rows
		for _, row := range ept.Body {
			delRowRes, errDelRow := delRowStmt.Exec(modifyUserID, row.BID, row.Ts)
			if errDelRow != nil {
				zap.L().Error("DeleteEPTs delRowStmt.Exec failed", zap.Error(errDelRow))
				resStatus = i18n.StatusInternalError
				tx.Rollback()
				return resStatus, errDelRow
			}
			// Check the number of affected rows
			updateRowNumber, errUpdateEffc := delRowRes.RowsAffected()
			if errUpdateEffc != nil {
				zap.L().Error("DeleteEPTs delRowRes.RowsAffected failed", zap.Error(errUpdateEffc))
				resStatus = i18n.StatusInternalError
				tx.Rollback()
				return resStatus, errUpdateEffc
			}
			// If no rows were affected, it means the row was edited by someone else
			if updateRowNumber < 1 {
				resStatus = i18n.StatusOtherEdit
				tx.Rollback()
				return
			}
		}
	}

	return
}

// Check if the Execution Project Template is referenced by other documents
func (ept *EPT) CheckUsed() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Define the items to check for references
	checkItems := []ArchiveCheckUsed{
		{
			Description:    "Referenced by Work Orders",
			SqlStr:         `select count(id) as usednumber from workorder_b where dr=0 and eptid=$1`,
			UsedReturnCode: i18n.StatusWOUsed,
		},
		{
			Description:    "Referenced by Execution Orders",
			SqlStr:         `select count(id) as usednumber from executionorder_b where dr=0 and eptid=$1`,
			UsedReturnCode: i18n.StatusEOUsed,
		},
	}
	// Check each item
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, ept.HID).Scan(&usedNum)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("EPT.CheckUsed "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return
}
