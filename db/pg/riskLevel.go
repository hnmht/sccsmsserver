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

// Risk Level struct
type RiskLevel struct {
	ID          int32     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Color       string    `db:"color" json:"color"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Creator     Person    `db:"creatorid" json:"creator"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	ModifyDate  time.Time `db:"modifytime" json:"modifyDate"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

// Risk Level front-end cache struct
type RLCache struct {
	QueryTs      time.Time   `json:"queryTs"`
	ResultNumber int32       `json:"resultNumber"`
	DelRLs       []RiskLevel `json:"delItems"`
	UpdateRLs    []RiskLevel `json:"updateItems"`
	NewRLs       []RiskLevel `json:"newItems"`
	ResultTs     time.Time   `json:"resultTs"`
}

// Initialize risklevel table
func initRiskLevel() (isFinish bool, err error) {
	// Step 1: Check if a record exists for the default Risk Level.
	sqlStr := "select count(id) as rownum from risklevel where dr=0"
	// Step 2: Exit if the record exists or an error occurs.
	hasRecord, isFinish, err := genericCheckRecord("risklevel", sqlStr)
	if hasRecord || !isFinish || err != nil {
		return
	}
	// Step 3: Insert default Risk Level records.
	sqlStrs := []string{
		"insert into risklevel(id,name,description,color,creatorid) values(1,'Major Risk','System pre-set','red',10000)",
		"insert into risklevel(id,name,description,color,creatorid) values(2,'Significant Risk','System pre-set','orange',10000)",
		"insert into risklevel(id,name,description,color,creatorid) values(3,'General Risk','System pre-set','yellow',10000)",
		"insert into risklevel(id,name,description,color,creatorid) values(4,'Low Risk','System pre-set','blue',10000)",
		"insert into risklevel(id,name,description,color,creatorid) values(5,'No Risk','System pre-set','white',10000)",
	}

	for _, t := range sqlStrs {
		_, err = db.Exec(t)
		if err != nil {
			isFinish = false
			zap.L().Error("initRiskLevel insert default data:"+t+" failed.", zap.Error(err))
			return
		}
	}
	return
}

// Get Risk Level List
func GetRLList() (rls []RiskLevel, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	rls = make([]RiskLevel, 0)
	// Retrieve Risk Level list from risklevel table
	sqlStr := `select id,name,description,color,status,
			createtime,creatorid,modifytime,modifierid,ts,
			dr 
			from risklevel 
			where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetRLList db.Query failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	for rows.Next() {
		var rl RiskLevel
		err = rows.Scan(&rl.ID, &rl.Name, &rl.Description, &rl.Color, &rl.Status,
			&rl.CreateDate, &rl.Creator.ID, &rl.ModifyDate, &rl.Modifier.ID, &rl.Ts,
			&rl.Dr)
		if err != nil {
			zap.L().Error("GetRLList rows.Scan failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get creator detail
		if rl.Creator.ID > 0 {
			resStatus, err = rl.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier detail
		if rl.Modifier.ID > 0 {
			resStatus, err = rl.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Append rl to rls slice
		rls = append(rls, rl)
	}

	return
}

// Get Risk Level front-end cache
func (rlc *RLCache) GetRLsCache() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	rlc.DelRLs = make([]RiskLevel, 0)
	rlc.NewRLs = make([]RiskLevel, 0)
	rlc.UpdateRLs = make([]RiskLevel, 0)
	// Retrieve the latest timestamp from the risklevel table
	sqlStr := `select ts from risklevel where ts > $1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, rlc.QueryTs).Scan(&rlc.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			rlc.ResultNumber = 0
			rlc.ResultTs = rlc.QueryTs
			resStatus = i18n.StatusOK
			return
		}
		zap.L().Error("RLCache.GetRLCsCache query latest ts failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	// Retrieve all data greator than the QueryTs
	sqlStr = `select a.id,a.name,a.description,a.color,a.status,
		a.createtime,a.creatorid,a.modifytime,a.modifierid,a.ts,
		a.dr
		from risklevel a
		where a.ts > $1 order by a.ts desc`
	rows, err := db.Query(sqlStr, rlc.QueryTs)
	if err != nil {
		zap.L().Error("RLCache.GetRLCsCache gdb.Query failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	// Extract data from result set
	for rows.Next() {
		var rl RiskLevel
		err = rows.Scan(&rl.ID, &rl.Name, &rl.Description, &rl.Color, &rl.Status,
			&rl.CreateDate, &rl.Creator.ID, &rl.ModifyDate, &rl.Modifier.ID, &rl.Ts,
			&rl.Dr)
		if err != nil {
			zap.L().Error("RLCache.GetRLCsCache rows.next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get Creator detail
		if rl.Creator.ID > 0 {
			resStatus, err = rl.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier detail
		if rl.Modifier.ID > 0 {
			resStatus, err = rl.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		if rl.Dr == 0 {
			if rl.CreateDate.Before(rlc.QueryTs) || rl.CreateDate.Equal(rlc.QueryTs) {
				rlc.ResultNumber++
				rlc.UpdateRLs = append(rlc.UpdateRLs, rl)
			} else {
				rlc.ResultNumber++
				rlc.NewRLs = append(rlc.NewRLs, rl)
			}
		} else {
			if rl.CreateDate.Before(rlc.QueryTs) || rl.CreateDate.Equal(rlc.QueryTs) {
				rlc.ResultNumber++
				rlc.DelRLs = append(rlc.DelRLs, rl)
			}
		}
	}

	return
}

// Add Risk Level
func (rl *RiskLevel) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the Risk Level name exists
	resStatus, err = rl.CheckNameExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Add data to the risklevel table
	sqlStr := `insert into risklevel(name,description,color,status,creatorid) 
	values($1,$2,$3,$4,$5) 
	returning id`
	err = db.QueryRow(sqlStr, rl.Name, rl.Description, rl.Color, rl.Status, rl.Creator.ID).Scan(&rl.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("RiskLevel.Add  stmt.QueryRow failed", zap.Error(err))
		return
	}
	return
}

// Check if the Risk Level name exist
func (rl *RiskLevel) CheckNameExist() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var count int32
	sqlStr := "select count(id) from risklevel where dr=0 and name=$1 and id <> $2"
	err = db.QueryRow(sqlStr, rl.Name, rl.ID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("RiskLevel.CheckNameExist query failed", zap.Error(err))
		return
	}
	if count > 0 {
		resStatus = i18n.StatusRLNameExist
		return
	}
	return
}

// Get Risk Level information by ID
func (rl *RiskLevel) GetRLInfoByID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get detail from cache
	number, b, _ := cache.Get(pub.RL, rl.ID)
	if number > 0 {
		json.Unmarshal(b, &rl)
		return
	}
	// If the Risk Level not in cache, then retrieve it from database
	sqlStr := `select a.name,a.description,a.color,a.status,a.createtime,
	a.creatorid,a.modifytime,a.modifierid,a.ts,a.dr
	from risklevel a
	where a.id = $1`

	err = db.QueryRow(sqlStr, rl.ID).Scan(&rl.Name, &rl.Description, &rl.Color, &rl.Status, &rl.CreateDate,
		&rl.Creator.ID, &rl.ModifyDate, &rl.Modifier.ID, &rl.Ts, &rl.Dr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("RiskLevel.GetRLInfoByID  db.QueryRow failed", zap.Error(err))
		return
	}
	// Get Creator detail
	if rl.Creator.ID > 0 {
		resStatus, err = rl.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Modifier detail
	if rl.Modifier.ID > 0 {
		resStatus, err = rl.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Write into cache
	rlB, _ := json.Marshal(rl)
	cache.Set(pub.RL, rl.ID, rlB)

	return
}

// Modify Risk Level
func (rl *RiskLevel) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the Risk Level name exists
	resStatus, err = rl.CheckNameExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Modify record in the risklevel table
	sqlStr := `update risklevel set 
	name=$1,description=$2,color=$3,status=$4,modifierid=$5,
	modifytime=current_timestamp,ts=current_timestamp 
	where id=$6 and ts=$7 and dr=0`
	res, err := db.Exec(sqlStr, rl.Name, rl.Description, rl.Color, rl.Status, rl.Modifier.ID,
		rl.ID, rl.Ts)
	if err != nil {
		zap.L().Error("RiskLevel.Edit db.exec failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Check the number of rows affected by the SQL update operation
	affected, err := res.RowsAffected()
	if err != nil {
		zap.L().Error("RiskLevel.Edit  get res.RowsAffected failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// If the number of rows effected by the update operation is less than one,
	// it means someone else has already updated the data
	if affected < 1 {
		zap.L().Info("RiskLevel.Edit failed,Other user are Editing")
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Delete from cache
	rl.DelFromLocalCache()

	return
}

// Delete Risk Level
func (rl *RiskLevel) Delete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the Risk Level id refrenced
	resStatus, err = rl.CheckUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update the delete flag for this record in risklevel table
	sqlStr := `update risklevel set dr=1,modifierid=$1,modifytime=current_timestamp,ts=current_timestamp where id=$2 and dr=0 and ts=$3`
	res, err := db.Exec(sqlStr, rl.Modifier.ID, rl.ID, rl.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("RiskLevel.Delete db.exec failed", zap.Error(err))
		return
	}

	// Check the number of rows affected by SQL update operation
	affected, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("RiskLevel.Delete res.RowsAffected failed", zap.Error(err))
		return
	}
	// if the number of rows affected by SQL update operation is less than one,
	// it means that someone else already updated the data.
	if affected < 1 {
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Delete from cache
	rl.DelFromLocalCache()

	return
}

// Batch delete Risk Level
func DeleteRLs(rls *[]RiskLevel, modifyUserId int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteRLs db.begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Prepare a SQL statement for execution
	delSqlStr := "update risklevel set dr=1,modifierid=$1,modifytime=current_timestamp,ts=current_timestamp where id=$2 and dr=0 and ts=$3"
	stmt, err := tx.Prepare(delSqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteRLs tx.Prepare failed", zap.Error(err))
		_ = tx.Rollback()
		return
	}
	defer stmt.Close()
	// Update the delete flag for each record one by one.
	for _, rl := range *rls {
		// Check if the Risk Level Id is referenced
		resStatus, err = rl.CheckUsed()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
		// Execute the update operation
		result, err1 := stmt.Exec(modifyUserId, rl.ID, rl.Ts)
		if err1 != nil {
			zap.L().Error("DeleteRLs stmt.exec failed", zap.Error(err))
			_ = tx.Rollback()
			return i18n.StatusInternalError, err1
		}
		// Check if the number of rows affected by the Update operation
		affected, err2 := result.RowsAffected()
		if err2 != nil {
			zap.L().Error("DeleteRLs check RowsAffected failed", zap.Error(err))
			_ = tx.Rollback()
			return i18n.StatusInternalError, err2
		}
		// If the number of rows affected by update operation less than 1,
		// it means that someone else already updated the record.
		if affected < 1 {
			zap.L().Info("DeleteRLs other edit")
			_ = tx.Rollback()
			return i18n.StatusOtherEdit, nil
		}
		// Delete from cache
		rl.DelFromLocalCache()
	}
	return
}

// Delete Risk Level from cache
func (rl *RiskLevel) DelFromLocalCache() {
	number, _, _ := cache.Get(pub.RL, rl.ID)
	if number > 0 {
		cache.Del(pub.RL, rl.ID)
	}
}

// Check the Risk Levelv is reference
func (rl *RiskLevel) CheckUsed() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Define a list of items to be checked.
	checkItems := []ArchiveCheckUsed{
		{
			Description:    "Referenced by Execution Project Archive",
			SqlStr:         `select count(id) as usednum from epa where  dr=0 and risklevelid=$1`,
			UsedReturnCode: i18n.StatusEPAUsed,
		},
		{
			Description:    "被执行模板引用",
			SqlStr:         `select count(id) from exectivetemplate_b where  dr=0 and risklevelid=$1`,
			UsedReturnCode: i18n.StatusEPTUsed,
		},
		/*	{
			 		Description:    "被执行单引用",
					SqlStr:         `select count(id) from executedoc_b where  dr=0 and risklevelid=$1`,
					UsedReturnCode: i18n.StatusEDUsed,
				},
				{
					Description:    "被问题处理单引用",
					SqlStr:         "select count(id) as usednum from disposedoc where dr=0  and risklevelid=$1",
					UsedReturnCode: i18n.StatusDDUsed,
				}, */
	}

	// Check one by one
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, rl.ID).Scan(&usedNum)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("RiskLevel.CheckIsUsed "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return
}
