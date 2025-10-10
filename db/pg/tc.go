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

// Traning course master data struct
type TC struct {
	ID          int32         `db:"id" json:"id"`
	Code        string        `db:"code" json:"code"`
	Name        string        `db:"name" json:"name"`
	ClassHour   float64       `db:"classhour" json:"classHour"`
	IsExamine   int16         `db:"isexamine" json:"isExamine"`
	Description string        `db:"description" json:"description"`
	Status      int16         `db:"status" json:"status"`
	Files       []VoucherFile `json:"files"`
	CreateDate  time.Time     `db:"createtime" json:"createDate"`
	Creator     Person        `db:"creatorid" json:"creator"`
	ModifyDate  time.Time     `db:"modifytime" json:"modifyDate"`
	Modifier    Person        `db:"modifierid" json:"modifier"`
	Ts          time.Time     `db:"ts" json:"ts"`
	Dr          int16         `db:"dr" json:"dr"`
}

// Traning Course front-end cache struct
type TCCache struct {
	QueryTs      time.Time `json:"queryTs"`
	ResultNumber int32     `json:"resultNumber"`
	DelItems     []TC      `json:"delItems"`
	UpdateItems  []TC      `json:"updateItems"`
	NewItems     []TC      `json:"newItems"`
	ResultTs     time.Time `json:"resultTs"`
}

// Add Traning Course
func (tc *TC) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the name already exists
	resStatus, err = tc.CheckNameExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Create a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TC.Add db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Insert the main record to the TC table
	headSql := `insert into tc(name,classhour,isexamine,description,status,creatorid)
	 	values($1,$2,$3,$4,$5,$6)
		returning id`
	err = tx.QueryRow(headSql, tc.Name, tc.ClassHour, tc.IsExamine, tc.Description, tc.Status, tc.Creator.ID).Scan(&tc.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TC.Add tx.QueryRow failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Prepare to insert attachments to the TC_file table
	fileSql := `insert into tc_file(billhid,fileid,creatorid)
	values($1,$2,$3)
	returning id`
	fileStmt, err := tx.Prepare(fileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TC.Add  tx.Prepare(fileSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer fileStmt.Close()
	// Insert files one by one
	for _, file := range tc.Files {
		err = fileStmt.QueryRow(tc.ID, file.File.ID, tc.Creator.ID).Scan(&file.ID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("TC.Add  fileStmt.QueryRow failed", zap.Error(err))
			tx.Rollback()
			return
		}
	}
	return
}

// Edit Traning Course
func (tc *TC) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the modifier and creator are the same person
	if tc.Creator.ID != tc.Modifier.ID {
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Check if the name already exists
	resStatus, err = tc.CheckNameExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Create a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TC.Edit db.Begin() failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Update the main record in the TC table
	editDocSql := `update tc set code=$1,name=$2,classhour=$3,isexamine=$4,description=$5,
		status=$6,modifytime=current_timestamp,modifierid=$7,ts=current_timestamp
		where id=$8 and dr=0 and ts=$9`
	editDocRes, err := tx.Exec(editDocSql, &tc.Code, &tc.Name, &tc.ClassHour, &tc.IsExamine, &tc.Description,
		&tc.Status, &tc.Modifier.ID,
		&tc.ID, &tc.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TC.Edit tx.Exec(editDocSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Get the number of rows affected by the SQL execution
	updatedRows, err := editDocRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TC.Edit editDocRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// If no rows were affected, it indicates that the data was modified by another user
	if updatedRows < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Update or add attachments in the TC_file table
	// If ID is not 0, update the existing attachment;
	updateFileSql := `update tc_file set modifierid=$1,modifytime=current_timestamp,dr=$2,ts=current_timestamp
		where id=$3 and billhid=$4 and dr=0 and ts=$5`
	// If ID is 0, add a new attachment
	addFileSql := `insert into tc_file(billhid,fileid,creatorid) values($1,$2,$3) returning id`
	// Prepare to update attachments
	updateFileStmt, err := tx.Prepare(updateFileSql)
	if err != nil {
		zap.L().Error("TC.Edit tx.Prepare(updateFileSql) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer updateFileStmt.Close()
	// Prepare to add attachments
	addFileStmt, err := tx.Prepare(addFileSql)
	if err != nil {
		zap.L().Error("TC.Edit tx.Prepare(addFileStmt) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer addFileStmt.Close()
	// Process attachments one by one
	for _, file := range tc.Files {
		// Update existing attachment
		if file.ID != 0 {
			updateFileRes, updateFileErr := updateFileStmt.Exec(tc.Modifier.ID, file.Dr, file.ID, tc.ID, file.Ts)
			if updateFileErr != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("TC.Edit  updateFileRes.Exec() failed", zap.Error(updateFileErr))
				tx.Rollback()
				return resStatus, updateFileErr
			}
			updateFileNumber, updateFileEffectErr := updateFileRes.RowsAffected()
			if updateFileEffectErr != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("TC.EditupdateFileRes.RowsAffected failed", zap.Error(updateFileEffectErr))
				tx.Rollback()
				return resStatus, updateFileEffectErr
			}
			if updateFileNumber < 1 {
				resStatus = i18n.StatusOtherEdit
				tx.Rollback()
				return
			}
		} else {
			// Add new attachment
			addFileErr := addFileStmt.QueryRow(tc.ID, file.File.ID, tc.Modifier.ID).Scan(&file.ID)
			if addFileErr != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("Document.Edit addFileStmt.QueryRow failed", zap.Error(addFileErr))
				tx.Rollback()
				return resStatus, addFileErr
			}
		}
	}
	// Delete form cache
	tc.DelFromCache()
	return
}

// Check if the training course name already exists
func (tc *TC) CheckNameExist() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var count int32
	sqlStr := `select count(id) from tc where dr=0 and name=$1 and id <> $2`
	err = db.QueryRow(sqlStr, tc.Name, tc.ID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TC.CheckNameExist query failed", zap.Error(err))
		return
	}
	if count > 0 {
		resStatus = i18n.StatusTCNameExist
		return
	}
	return
}

// Delete TC from cache
func (tc *TC) DelFromCache() {
	number, _, _ := cache.Get(pub.TC, tc.ID)
	if number > 0 {
		cache.Del(pub.TC, tc.ID)
	}
}

// Delete Traning Course
func (tc *TC) Delete(modifyUserId int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the training course is referenced
	resStatus, err = tc.CheckIsUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check if the modifier and creator are the same person
	if tc.Creator.ID != modifyUserId {
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TC.Delete db.Begin() failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Update the main record in the TC table to mark it as deleted
	delDocSql := `update tc set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp
		where id=$2 and dr=0 and ts=$3`
	delDocRes, err := tx.Exec(delDocSql, modifyUserId, tc.ID, tc.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TC.Delete tx.Exec(delDocSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by the SQL execution
	deletedRows, err := delDocRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TC.Delete delDocRes.RowsAffected() failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// If no rows were affected,
	// it indicates that the data was modified by another user
	if deletedRows < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Mark attachments in the TC_file table as deleted
	delFileSql := `update tc_file set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp
		where id=$2 and dr=0 and billhid=$3 and ts=$4`
	delFileStmt, err := tx.Prepare(delFileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TC.Delete tx.Prepare(delFileSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delFileStmt.Close()
	// Process attachments one by one
	for _, file := range tc.Files {
		delFileRes, errDelFile := delFileStmt.Exec(modifyUserId, file.ID, file.BillHID, file.Ts)
		if errDelFile != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("TC.Delete delFileStmt.Exec failed", zap.Error(errDelFile))
			tx.Rollback()
			return
		}
		delFileNumber, errDelEff := delFileRes.RowsAffected()
		if errDelEff != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("TC.Delete delFileRes.RowsAffected failed", zap.Error(errDelEff))
			tx.Rollback()
			return
		}
		if delFileNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}
	}
	// Delete from cache
	tc.DelFromCache()
	return
}

// Batch delete Traning Courses
func DeleteTCs(tcs *[]TC, modifyUserID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	if len(*tcs) == 0 {
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteTCs db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Prepare to delete main records
	delDocSql := `update tc set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp
		where id=$2 and dr=0 and ts=$3`
	docStmt, err := tx.Prepare(delDocSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteTCs tx.Prepare(delDocSql) failed", zap.Error(err))
		tx.Rollback()
		return resStatus, err
	}
	defer docStmt.Close()
	// Prepare to delete attachments
	delFileSql := `update tc_file set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp
		where id=$2 and dr=0 and billhid=$3 and ts=$4`
	fileStmt, err := tx.Prepare(delFileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteTCs tx.Prepare(delFileSql) failed", zap.Error(err))
		tx.Rollback()
		return resStatus, err
	}
	defer fileStmt.Close()
	// Process each training course one by one
	for _, tc := range *tcs {
		// Check if the training course is referenced
		resStatus, err = tc.CheckIsUsed()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
		// Check if the modifier and creator are the same person
		if tc.Creator.ID != modifyUserID {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}
		// Delete the main record
		delDocRes, errDelDoc := docStmt.Exec(modifyUserID, tc.ID, tc.Ts)
		if errDelDoc != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("DeleteTCs docStmt.Exec failed", zap.Error(errDelDoc))
			tx.Rollback()
			return
		}
		// Check the number of rows affected by the SQL execution
		deletedRows, errDelDocEff := delDocRes.RowsAffected()
		if errDelDocEff != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("DeleteTCs delDocRes.RowsAffected failed", zap.Error(errDelDocEff))
			tx.Rollback()
			return
		}
		// If no rows were affected, it indicates that the data was modified by another user
		if deletedRows < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}

		// Delete attachments one by one
		for _, file := range tc.Files {
			delFileRes, errDelFile := fileStmt.Exec(modifyUserID, file.ID, file.BillHID, file.Ts)
			if errDelFile != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("DeleteTCs fileStmt.Exec failed", zap.Error(errDelFile))
				tx.Rollback()
				return
			}
			delFileNumber, errDelEff := delFileRes.RowsAffected()
			if errDelEff != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("DeleteTCs delFileRes.RowsAffected failed", zap.Error(errDelEff))
				tx.Rollback()
				return
			}
			if delFileNumber < 1 {
				resStatus = i18n.StatusOtherEdit
				tx.Rollback()
				return
			}
		}
		// Delete from cache
		tc.DelFromCache()
	}
	return
}

// Get Traning Course details by ID
func (tc *TC) GetDetailByID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get from cache
	number, b, _ := cache.Get(pub.TC, tc.ID)
	if number > 0 {
		json.Unmarshal(b, &tc)
		resStatus = i18n.StatusOK
		return
	}
	// If not in cache, get from database
	sqlStr := `select code,name,classhour,isexamine,description,
		status,createtime,creatorid,modifytime,modifierid,
		ts,dr
		from tc where id=$1`
	err = db.QueryRow(sqlStr, tc.ID).Scan(&tc.Code, &tc.Name, &tc.ClassHour, &tc.IsExamine, &tc.Description,
		&tc.Status, &tc.CreateDate, &tc.Creator.ID, &tc.ModifyDate, &tc.Modifier.ID,
		&tc.Ts, &tc.Dr)
	if err != nil {
		zap.L().Error("TC.GetDetailByID db.QueryRow failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Get creator details
	if tc.Creator.ID > 0 {
		resStatus, err = tc.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get modifier details
	if tc.Modifier.ID > 0 {
		resStatus, err = tc.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get attachments
	filesStr := `select id,billhid,fileid,createtime,creatorid,
		modifytime,modifierid,ts,dr
		from tc_file
		where dr=0 and billhid=$1`
	rows, err := db.Query(filesStr, tc.ID)
	if err != nil {
		zap.L().Error("tc.GetDetailByID db.Query(files) failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	for rows.Next() {
		var file VoucherFile
		err = rows.Scan(&file.ID, &file.BillHID, &file.File.ID, &file.CreateDate, &file.Creator.ID,
			&file.ModifyDate, &file.Modifier.ID, &file.Ts, &file.Dr)
		if err != nil {
			zap.L().Error("tc.GetDetailByID db.Query(files)  rows.Scan  failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get file details
		if file.File.ID > 0 {
			resStatus, err = file.File.GetFileInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get creator details
		if file.Creator.ID > 0 {
			resStatus, err = file.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get modifier details
		if file.Modifier.ID > 0 {
			resStatus, err = file.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Append to the array
		tc.Files = append(tc.Files, file)
	}
	// Save to cache
	tcB, _ := json.Marshal(tc)
	cache.Set(pub.TC, tc.ID, tcB)
	return
}

// Get Traning Course list
func GetTCList() (tcs []TC, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	tcs = make([]TC, 0)
	// Query all training course IDs
	sqlStr := `select id from tc where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetTCList db.Query failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()
	// Extract data from query results
	for rows.Next() {
		var tc TC
		err = rows.Scan(&tc.ID)
		if err != nil {
			zap.L().Error("GetTCList row.Next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get training course details
		resStatus, err = tc.GetDetailByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		// Append to the slice
		tcs = append(tcs, tc)
	}
	return
}

// Get Traning Course changes since QueryTs for front-end cache update
func (tcc *TCCache) GetTCCache() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Initialize return values
	tcc.DelItems = make([]TC, 0)
	tcc.NewItems = make([]TC, 0)
	tcc.UpdateItems = make([]TC, 0)
	// Retrieve the latest timestamp greater than QueryTs
	sqlStr := `select ts from tc where ts > $1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, tcc.QueryTs).Scan(&tcc.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			tcc.ResultNumber = 0
			tcc.ResultTs = tcc.QueryTs
			resStatus = i18n.StatusOK
			return
		}
		zap.L().Error("TCCache.GetTCCache query latest ts failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Retrieve all IDs with timestamp greater than QueryTs
	sqlStr = `select id
	from tc
	where ts > $1 order by ts desc`
	rows, err := db.Query(sqlStr, tcc.QueryTs)
	if err != nil {
		zap.L().Error("TCCache.GetTCCache get cache from database failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()
	// Process each ID one by one
	for rows.Next() {
		var tc TC
		err = rows.Scan(&tc.ID)
		if err != nil {
			zap.L().Error("TCCache.GetTCCache rows.Next() failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		resStatus, err = tc.GetDetailByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		if tc.Dr == 0 { //Record not deleted
			if tc.CreateDate.Before(tcc.QueryTs) || tc.CreateDate.Equal(tcc.QueryTs) {
				// If the reccord was created before or at queryTs, it indicates it was updated after QueryTs
				tcc.ResultNumber++
				tcc.UpdateItems = append(tcc.UpdateItems, tc)
			} else {
				// If the record was created after QueryTs, it indicates it is a new record
				tcc.ResultNumber++
				tcc.NewItems = append(tcc.NewItems, tc)
			}
		} else { // Record has been deleted
			if tc.CreateDate.Before(tcc.QueryTs) || tc.CreateDate.Equal(tcc.QueryTs) {
				// If the record was created before or at QueryTs, it indicates it was deleted after QueryTs
				tcc.ResultNumber++
				tcc.DelItems = append(tcc.DelItems, tc)
			}
			// if the record was created after QueryTs and then deleted, it is ignored
		}
	}
	return
}

// Check if the training course is referenced
func (tc *TC) CheckIsUsed() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Define check items
	checkItems := []ArchiveCheckUsed{
		{
			Description:    "Refrenced by training record",
			SqlStr:         `select count(id) as usedNum from trainingrecord_h where tcid=$1 and dr=0`,
			UsedReturnCode: i18n.StatusTRUsed,
		},
	}
	// check one by one
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, tc.ID).Scan(&usedNum)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("TC.CheckIsUsed "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return
}
