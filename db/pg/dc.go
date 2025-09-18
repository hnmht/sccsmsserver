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

// Document Category Master Data
type DC struct {
	ID          int32     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Father      SimpDC    `db:"fatherid" json:"fatherID"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Creator     Person    `db:"creatorid" json:"creator"`
	ModifyDate  time.Time `db:"modifytime" json:"modifyDate"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

// Simplify Document category
type SimpDC struct {
	ID          int32     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	FatherID    int32     `db:"fatherid" json:"fatherID"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Creator     Person    `db:"creatorid" json:"creator"`
	ModifyDate  time.Time `db:"modifytime" json:"modifyDate"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

// Simple Document Category front-end cache
type SimpDCCache struct {
	QueryTs      time.Time `json:"queryTs"`
	ResultNumber int32     `json:"resultNumber"`
	DelItems     []SimpDC  `json:"delItems"`
	UpdateItems  []SimpDC  `json:"updateItems"`
	NewItems     []SimpDC  `json:"newItems"`
	ResultTs     time.Time `json:"resultTs"`
}

// Initialize Document Category table
func initDocumentCategory() (isFinish bool, err error) {
	// Step 1: Check if a record exists for the defualt Document Category.
	sqlStr := "select count(id) as rownum from dc where id=1"
	// Step 2: Exit if the record exists or an error occurs.
	hasRecord, isFinish, err := genericCheckRecord("dc", sqlStr)
	if hasRecord || !isFinish || err != nil {
		return
	}
	// Step 3: Insert a record for the system default DC "Default Category" into the
	sqlStr = `insert into dc(id,name,description,creatorid) values(10000,'Default Category','System Pre-Set',10000)`
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initDocumentCategory insert initvalue failed", zap.Error(err))
		return isFinish, err
	}
	return
}

// Get Simple Document Category information by ID
func (sdc *SimpDC) GetSDCInfoByID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get SimpDC information from cache
	number, b, _ := cache.Get(pub.SimpDC, sdc.ID)
	if number > 0 {
		json.Unmarshal(b, &sdc)
		resStatus = i18n.StatusOK
		return
	}
	// If data not in cache, get from the database
	sqlStr := `select name,description,fatherid,status,createtime,
	creatorid,modifytime,modifierid,ts,dr
	from dc where id=$1`
	err = db.QueryRow(sqlStr, sdc.ID).Scan(&sdc.Name, &sdc.Description, &sdc.FatherID, &sdc.Status, &sdc.CreateDate,
		&sdc.Creator.ID, &sdc.ModifyDate, &sdc.Modifier.ID, &sdc.Ts, &sdc.Dr)
	if err != nil {
		zap.L().Error("GetSDCInfoByID db.QueryRow failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Get creator detail
	if sdc.Creator.ID > 0 {
		resStatus, err = sdc.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get modifier detail
	if sdc.Modifier.ID > 0 {
		resStatus, err = sdc.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Write in cache
	sdcB, _ := json.Marshal(sdc)
	cache.Set(pub.SimpDC, sdc.ID, sdcB)

	return
}

// Get Document Category list
func GetDCList() (dcs []DC, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	dcs = make([]DC, 0)
	// Retrieve Document Category list from the dc table
	sqlStr := `select id,name,description,fatherid,status,
	createtime,creatorid,modifytime,modifierid,ts,dr
	from dc
	where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetDCList db.Query failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()
	// Extract item by item from the returned rows
	for rows.Next() {
		var dc DC
		err = rows.Scan(&dc.ID, &dc.Name, &dc.Description, &dc.Father.ID, &dc.Status,
			&dc.CreateDate, &dc.Creator.ID, &dc.ModifyDate, &dc.Modifier.ID, &dc.Ts, &dc.Dr)
		if err != nil {
			zap.L().Error("GetDCList row.Next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get Father DC detail
		if dc.Father.ID > 0 {
			resStatus, err = dc.Father.GetSDCInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return
			}
		}
		// Get Creator detail
		if dc.Creator.ID > 0 {
			resStatus, err = dc.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get modifier detail
		if dc.Modifier.ID > 0 {
			resStatus, err = dc.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Append to slice
		dcs = append(dcs, dc)
	}

	return
}

// Get Simple Document Category List
func GetSimpDCList() (sdcs []SimpDC, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	sdcs = make([]SimpDC, 0)
	// Retrieve Document Category list from the dc table
	sqlStr := `select id,name,description,fatherid,status,
	createtime,creatorid,modifytime,modifierid,ts,dr
	from dc
	where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetSimpDCList db.Query failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()
	// Extract item by item from the returned rows
	for rows.Next() {
		var sdc SimpDC
		err = rows.Scan(&sdc.ID, &sdc.Name, &sdc.Description, &sdc.FatherID, &sdc.Status,
			&sdc.CreateDate, &sdc.Creator.ID, &sdc.ModifyDate, &sdc.Modifier.ID, &sdc.Ts, &sdc.Dr)
		if err != nil {
			zap.L().Error("GetSimpDCList row.Next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get creator detail
		if sdc.Creator.ID > 0 {
			resStatus, err = sdc.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get modifier deatail
		if sdc.Modifier.ID > 0 {
			resStatus, err = sdc.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Append to slice
		sdcs = append(sdcs, sdc)
	}
	return
}

// Get latest SimpDC front-end cache
func (sdcc *SimpDCCache) GetSimpDCCache() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	sdcc.DelItems = make([]SimpDC, 0)
	sdcc.NewItems = make([]SimpDC, 0)
	sdcc.UpdateItems = make([]SimpDC, 0)
	// Query the latest timestamp from the dc table that is greater than QueryTs
	sqlStr := `select ts from dc where ts > $1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, sdcc.QueryTs).Scan(&sdcc.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			sdcc.ResultNumber = 0
			sdcc.ResultTs = sdcc.QueryTs
			resStatus = i18n.StatusOK
			return
		}
		zap.L().Error("SimpDCCache.GetSimpDCCache query latest ts failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	// Retrieve all data that timestamp greater than QueryTs
	sqlStr = `select id,name,description,fatherid,status,
	createtime,creatorid,modifytime,modifierid,ts,dr
	from dc
	where ts > $1 order by ts desc`
	rows, err := db.Query(sqlStr, sdcc.QueryTs)
	if err != nil {
		zap.L().Error("SimpDCCache.GetSimpDCCache get cache from database failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	// Extract data item by item from the returned rows
	for rows.Next() {
		var sdc SimpDC
		err = rows.Scan(&sdc.ID, &sdc.Name, &sdc.Description, &sdc.FatherID, &sdc.Status,
			&sdc.CreateDate, &sdc.Creator.ID, &sdc.ModifyDate, &sdc.Modifier.ID, &sdc.Ts, &sdc.Dr)
		if err != nil {
			zap.L().Error("SimpDCCache.GetSimpDCCache rows.Next() failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}

		// Get creator detail
		if sdc.Creator.ID > 0 {
			resStatus, err = sdc.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get modifier detail
		if sdc.Modifier.ID > 0 {
			resStatus, err = sdc.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		if sdc.Dr == 0 {
			if sdc.CreateDate.Before(sdcc.QueryTs) || sdc.CreateDate.Equal(sdcc.QueryTs) {
				sdcc.ResultNumber++
				sdcc.UpdateItems = append(sdcc.UpdateItems, sdc)
			} else {
				sdcc.ResultNumber++
				sdcc.NewItems = append(sdcc.NewItems, sdc)
			}
		} else {
			if sdc.CreateDate.Before(sdcc.QueryTs) || sdc.CreateDate.Equal(sdcc.QueryTs) {
				sdcc.ResultNumber++
				sdcc.DelItems = append(sdcc.DelItems, sdc)
			}
		}
	}
	return
}

// Check if the Document Category name exists
func (dc *DC) CheckNameExist() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var count int32
	sqlStr := `select count(id) from dc where dr=0 and name=$1 and id <> $2`
	err = db.QueryRow(sqlStr, dc.Name, dc.ID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DC.CheckNameExist query failed", zap.Error(err))
		return
	}
	if count > 0 {
		resStatus = i18n.StatusDCNameExist
		return
	}
	return
}

// Add Document Category
func (dc *DC) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the DC name exists
	resStatus, err = dc.CheckNameExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Inser record into the dc table
	sqlStr := `insert into dc(name,description,fatherid,status,creatorid)
	values($1,$2,$3,$4,$5)
	returning id`
	err = db.QueryRow(sqlStr, dc.Name, dc.Description, dc.Father.ID, dc.Status, dc.Creator.ID).Scan(&dc.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DC.Add stmt.QueryRow failed", zap.Error(err))
		return
	}
	return
}

// Edit Document Category
func (dc *DC) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check whether the superior category is compliant
	// If the superior category ID is great than 0,
	// it means s superior category exists.
	if dc.Father.ID > 0 {
		if dc.ID == dc.Father.ID {
			resStatus = i18n.StatusDCFatherSelf
			return
		}
		// Check if the superior category is circular.
		dcs, res, err1 := GetSimpDCList()
		if resStatus != i18n.StatusOK || err1 != nil {
			return res, err1
		}
		childrens := FindSimpDCChildrens(dcs, dc.ID)
		var number int32
		for _, child := range childrens {
			if child.ID == dc.Father.ID {
				number++
			}
		}

		if number > 0 {
			resStatus = i18n.StatusDCFatherSelf
			return
		}
	}
	// Check if the DC name exists
	resStatus, err = dc.CheckNameExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update the record in the dc table
	sqlStr := `update dc set
	name=$1,description=$2,fatherid=$3,status=$4,
	modifytime=current_timestamp,modifierid=$5,ts=current_timestamp
	where id=$6 and dr=0 and ts=$7`
	res, err := db.Exec(sqlStr, dc.Name, dc.Description, dc.Father.ID, dc.Status,
		dc.Modifier.ID, dc.ID, dc.Ts)
	if err != nil {
		zap.L().Error("DC.Edit db.Exec failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Get the number of rows affected by SQL update operation.
	affected, err := res.RowsAffected()
	if err != nil {
		zap.L().Error("DC.Edit res.RowsAffected failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// If the number of affected rows is less than one,
	// it means that someone else has already modified the record.
	if affected < 1 {
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Delete from cache
	dc.DelFromLocalCache()

	return
}

// Check if the DC id is referenced
func (dc *DC) CheckIsUsed() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Define the item to be checked
	checkItems := []ArchiveCheckUsed{
		{
			Description:    "Referenced by Low level DC",
			SqlStr:         `select count(id) as usedNum from dc where fatherid=$1 and dr=0`,
			UsedReturnCode: i18n.StatusDCLowLevelExist,
		},
		/* {
			Description:    "被文档引用",
			SqlStr:         `select count(id) as usedNum from document where dc_id=$1 and dr=0`,
			UsedReturnCode: i18n.StatusDocumentUsed,
		}, */
	}
	// Check item by item
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, dc.ID).Scan(&usedNum)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("DC.CheckIsUsed "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return
}

// Delete Document Category
func (dc *DC) Delete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the category is referenced
	resStatus, err = dc.CheckIsUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update the record in the dc table
	sqlStr := `update dc set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp where id=$2 and dr=0 and ts=$3`
	res, err := db.Exec(sqlStr, dc.Modifier.ID, dc.ID, dc.Ts)
	if err != nil {
		zap.L().Error("DC.Delete stmt.Exec failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Check the number of rows affected by SQL update operation.
	affected, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DC.Delete res.RowsAffedted failed", zap.Error(err))
		return
	}

	if affected < 1 {
		resStatus = i18n.StatusOtherEdit
		zap.L().Info("DC.Delete Other User Edit")
		return
	}
	// Delete from cache
	dc.DelFromLocalCache()
	return
}

// Delete DC from cache
func (dc *DC) DelFromLocalCache() {
	number, _, _ := cache.Get(pub.SimpDC, dc.ID)
	if number > 0 {
		cache.Del(pub.SimpDC, dc.ID)
	}
}

// Batch Delete Document Category
func DeleteDCs(dcs *[]DC, modifyUserID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteDCs db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Prepare update SQL statement
	delSqlStr := `update dc set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp where id=$2 and dr=0 and ts=$3`
	stmt, err := tx.Prepare(delSqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteDCs Delete prepare failed", zap.Error(err))
		tx.Rollback()
		return resStatus, err
	}
	defer stmt.Close()

	for _, dc := range *dcs {
		// Check if the category is referenced
		resStatus, err = dc.CheckIsUsed()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}

		// Execute to update
		res, err := stmt.Exec(modifyUserID, dc.ID, dc.Ts)
		if err != nil {
			zap.L().Error("DeleteDCs stmt.Exec failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			tx.Rollback()
			return resStatus, err
		}
		// Check the number of rows affected by the SQL update operation.
		affectRows, err := res.RowsAffected()
		if err != nil {
			zap.L().Error("DeleteDCs res.RowsAffected failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			tx.Rollback()
			return resStatus, err
		}

		if affectRows < 1 {
			resStatus = i18n.StatusOtherEdit
			zap.L().Info("DeleteDCs: " + dc.Name + " other user eidting")
			tx.Rollback()
			return resStatus, nil
		}
		// Delete from cache
		dc.DelFromLocalCache()
	}

	return
}

// Find all child Document Category based on the Document Category
func FindSimpDCChildrens(sdcs []SimpDC, id int32) []SimpDC {
	childrens := make([]SimpDC, 0)
	for _, sdc := range sdcs {
		if sdc.FatherID == id {
			childrens = append(childrens, sdc)
			child := FindSimpDCChildrens(sdcs, sdc.ID)
			childrens = append(childrens, child...)
		}
	}
	return childrens
}
