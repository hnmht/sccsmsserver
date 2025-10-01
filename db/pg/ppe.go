package pg

import (
	"database/sql"
	"encoding/json"
	"sccsmsserver/cache"
	"sccsmsserver/i18n"
	"sccsmsserver/pub"
	"time"

	"go.uber.org/zap"
)

// Personal Protective Equipment
type PPE struct {
	ID          int32     `db:"id" json:"id"`
	Code        string    `db:"code" json:"code"`
	Name        string    `db:"name" json:"name"`
	Model       string    `db:"model" json:"model"`
	Unit        string    `db:"unit" json:"unit"`
	Description string    `db:"description" json:"description"`
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Creator     Person    `db:"creatorid" json:"creator"`
	ModifyDate  time.Time `db:"modifytime" json:"modifyDate"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

// Personal Protective Equipment data from front-end cache
type PPECache struct {
	QueryTs      time.Time `json:"queryTs"`
	ResultNumber int32     `json:"resultNumber"`
	DelItems     []PPE     `json:"delItems"`
	UpdateItems  []PPE     `json:"updateItems"`
	NewItems     []PPE     `json:"newItems"`
	ResultTs     time.Time `json:"resultTs"`
}

// Get Personal Protective Equipment master data list
func GetPPEList() (ppes []PPE, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	ppes = make([]PPE, 0)
	// Retrieve data from ppe table
	sqlStr := `select id,code,name,model,unit,
		description,createtime,creatorid,modifytime,modifierid,
		ts,dr
		from ppe  
		where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetPPEList  db.Query failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	// Extract data item by item from the returned rows
	for rows.Next() {
		var ppe PPE
		err = rows.Scan(&ppe.ID, &ppe.Code, &ppe.Name, &ppe.Model, &ppe.Unit,
			&ppe.Description, &ppe.CreateDate, &ppe.Creator.ID, &ppe.ModifyDate, &ppe.Modifier.ID,
			&ppe.Ts, &ppe.Dr)
		if err != nil {
			zap.L().Error("GetPPEList from rows failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}

		// Get Creator detail
		if ppe.Creator.ID > 0 {
			resStatus, err = ppe.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		// Get Modifier detail
		if ppe.Modifier.ID > 0 {
			resStatus, err = ppe.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Append to Slice
		ppes = append(ppes, ppe)
	}

	return
}

// Get latest PPE front-end cache
func (ppec *PPECache) GetPPEsCache() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	ppec.DelItems = make([]PPE, 0)
	ppec.NewItems = make([]PPE, 0)
	ppec.UpdateItems = make([]PPE, 0)
	// Query the latest timestamp in the ppe table that is greater than QueryTs
	sqlStr := `select ts from ppe where ts > $1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, ppec.QueryTs).Scan(&ppec.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			ppec.ResultNumber = 0
			ppec.ResultTs = ppec.QueryTs
			resStatus = i18n.StatusOK
			return
		}
		zap.L().Error("PPECache.GetPPEsCache query latest ts failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	// Retrieve all data that timestamp greater than QueryTs
	sqlStr = `select id,code,name,model,unit,
		description,createtime,creatorid,modifytime,modifierid,
		ts,dr	
		from ppe 
		where ts > $1 order by ts desc`
	rows, err := db.Query(sqlStr, ppec.QueryTs)
	if err != nil {
		zap.L().Error("PPECache.GetPPEsCache  get Cache from database failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	// Extract data item by item from the returned rows
	for rows.Next() {
		var ppe PPE
		err = rows.Scan(&ppe.ID, &ppe.Code, &ppe.Name, &ppe.Model, &ppe.Unit,
			&ppe.Description, &ppe.CreateDate, &ppe.Creator.ID, &ppe.ModifyDate, &ppe.Modifier.ID,
			&ppe.Ts, &ppe.Dr)
		if err != nil {
			zap.L().Error("PPECache.GetPPEsCache rows.next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get creator detail
		if ppe.Creator.ID > 0 {
			resStatus, err = ppe.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get modifier detail
		if ppe.Modifier.ID > 0 {
			resStatus, err = ppe.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		if ppe.Dr == 0 {
			if ppe.CreateDate.Before(ppec.QueryTs) || ppe.CreateDate.Equal(ppec.QueryTs) {
				ppec.ResultNumber++
				ppec.UpdateItems = append(ppec.UpdateItems, ppe)
			} else {
				ppec.ResultNumber++
				ppec.NewItems = append(ppec.NewItems, ppe)
			}
		} else {
			if ppe.CreateDate.Before(ppec.QueryTs) || ppe.CreateDate.Equal(ppec.QueryTs) {
				ppec.ResultNumber++
				ppec.DelItems = append(ppec.DelItems, ppe)
			}
		}
	}

	return
}

// Add Personal Protective Equipment
func (ppe *PPE) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the PPE Code exist
	resStatus, err = ppe.CheckCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Insert a record to ppe table
	sqlStr := `insert into ppe(code,name,model,unit,description,
		creatorid) 
		values($1,$2,$3,$4,$5,$6) 
		returning id`
	err = db.QueryRow(sqlStr, ppe.Code, ppe.Name, ppe.Model, ppe.Unit, ppe.Description,
		ppe.Creator.ID).Scan(&ppe.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPE.Add db.QueryRow failed", zap.Error(err))
		return
	}
	return
}

// Get PPE Information by ID
func (ppe *PPE) GetInfoByID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get PPE information from cache
	number, b, _ := cache.Get(pub.PPE, ppe.ID)
	if number > 0 {
		json.Unmarshal(b, &ppe)
		resStatus = i18n.StatusOK
		return
	}
	// If PPE infromation is not in cahce, retrieve it from database
	sqlStr := `select code,name,model,unit,description,
	createtime,creatorid,modifytime,modifierid,ts,
	dr 
	from ppe
	where id = $1`
	err = db.QueryRow(sqlStr, ppe.ID).Scan(&ppe.Code, &ppe.Name, &ppe.Model, &ppe.Unit, &ppe.Description,
		&ppe.CreateDate, &ppe.Creator.ID, &ppe.ModifyDate, &ppe.Modifier.ID, &ppe.Ts,
		&ppe.Dr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPE.GetInfoByID db.QueryRow failed", zap.Error(err))
		return
	}
	// Get Creator detail
	if ppe.Creator.ID > 0 {
		resStatus, err = ppe.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Modifier detail
	if ppe.Modifier.ID > 0 {
		resStatus, err = ppe.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Write in cache
	ppeB, _ := json.Marshal(ppe)
	cache.Set(pub.PPE, ppe.ID, ppeB)

	return
}

// Modify Personal Protective Equipment
func (ppe *PPE) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the PPE code exists
	resStatus, err = ppe.CheckCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update the record in the ppe table
	sqlStr := `update ppe set 
		code=$1,name=$2,model=$3,unit=$4,description=$5,
		modifierid=$6,modifytime=current_timestamp,ts=current_timestamp 
		where id=$7 and ts=$8 and dr=0`
	res, err := db.Exec(sqlStr, ppe.Code, ppe.Name, ppe.Model, ppe.Unit, ppe.Description,
		ppe.Modifier.ID,
		ppe.ID, ppe.Ts)
	if err != nil {
		zap.L().Error("PPE.Edit db.exec failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Get the number of rows affected by the SQL statement update
	affected, err := res.RowsAffected()
	if err != nil {
		zap.L().Error("PPE.Edit  get res.RowsAffected failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// If the number of affected rows is less than one,
	// it means that someone else has already modified the record.
	if affected < 1 {
		zap.L().Info("PPE.Edit failed,Other user are Editing")
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Delete from cache
	ppe.DelFromLocalCache()

	return
}

// Delete PPE master data
func (ppe *PPE) Delete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the PPE id is refereced
	resStatus, err = ppe.CheckUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update the record in the ppe table
	sqlStr := `update ppe set dr=1,modifierid=$1,modifytime=current_timestamp,ts=current_timestamp where id=$2 and dr=0 and ts=$3`
	res, err := db.Exec(sqlStr, ppe.Modifier.ID, ppe.ID, ppe.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPE.Delete db.exec failed", zap.Error(err))
		return
	}
	// Check the number of rows affected by the SQL update statement
	affected, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPE.Delete res.RowsAffected failed", zap.Error(err))
		return
	}
	// If the number of affected rows is less than one,
	// it means that someone else has already updated the record.
	if affected < 1 {
		resStatus = i18n.StatusOtherEdit
		return
	}
	// delete from cache
	ppe.DelFromLocalCache()

	return
}

// Check if the PPE code exists
func (ppe *PPE) CheckCodeExist() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var count int32
	sqlStr := "select count(id) from ppe where dr=0 and code=$1 and id <> $2"
	err = db.QueryRow(sqlStr, ppe.Code, ppe.ID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPE.CheckCodeExist query failed", zap.Error(err))
		return
	}
	if count > 0 {
		resStatus = i18n.StatusPPECodeExist
		return
	}

	return
}

// Batch Delete PPE master data
func DeletePPEs(ppes *[]PPE, modifyUserId int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeletePPEs db.begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Prepare update SQL statement
	delSqlStr := `update ppe set dr=1,modifierid=$1,modifytime=current_timestamp,ts=current_timestamp 
		where id=$2 and dr=0 and ts=$3`
	stmt, err := tx.Prepare(delSqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeletePPEs tx.Prepare failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer stmt.Close()

	for _, ppe := range *ppes {
		// Check if the PPE id is referenced
		resStatus, err = ppe.CheckUsed()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
		// Execute the update
		result, err1 := stmt.Exec(modifyUserId, ppe.ID, ppe.Ts)
		if err1 != nil {
			zap.L().Error("DeletePPEs stmt.exec failed", zap.Error(err))
			tx.Rollback()
			return i18n.StatusInternalError, err1
		}
		// Check the number of rows affected by the Update execution.
		affected, err2 := result.RowsAffected()
		if err2 != nil {
			zap.L().Error("DeletePPEs check RowsAffected failed", zap.Error(err))
			tx.Rollback()
			return i18n.StatusInternalError, err2
		}
		if affected < 1 {
			zap.L().Info("DeletePPEs other edit")
			_ = tx.Rollback()
			return i18n.StatusOtherEdit, nil
		}
		// Delete from cache
		ppe.DelFromLocalCache()
	}
	return
}

// Delete PPE from cache
func (ppe *PPE) DelFromLocalCache() {
	number, _, _ := cache.Get(pub.PPE, ppe.ID)
	if number > 0 {
		cache.Del(pub.PPE, ppe.ID)
	}
}

// Check if the PPE id is referenced
func (ppe *PPE) CheckUsed() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Define the items to be checked
	checkItems := []ArchiveCheckUsed{
		{
			Description:    "Referenced by PPE Issuance From body",
			SqlStr:         "select count(id) as usednum from ppeissuanceform_b where dr=0 and ppeid=$1",
			UsedReturnCode: i18n.StatusPPEIFUsed,
		},
	}
	// Check item by item
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, ppe.ID).Scan(&usedNum)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPE.CheckIsUsed "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return
}
