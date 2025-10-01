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

// Position Master Data
type Position struct {
	ID          int32     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Creator     Person    `db:"creatorid" json:"creator"`
	ModifyDate  time.Time `db:"modifytime" json:"modifyDate"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	Dr          int16     `db:"dr" json:"dr"`
	Ts          time.Time `db:"ts" json:"ts"`
}

// Latest Position Master Data for front end cache
type PositionCache struct {
	QueryTs      time.Time  `json:"queryTs"`
	ResultNumber int32      `json:"resultNumber"`
	DelItems     []Position `json:"delItems"`
	UpdateItems  []Position `json:"updateItems"`
	NewItems     []Position `json:"newItems"`
	ResultTs     time.Time  `json:"resultTs"`
}

// Initialize postion table
func initPosition() (isFinish bool, err error) {
	// Step 1: Check if a record exists for the default position
	sqlStr := "select count(id) as rownum from position where id=10000"
	hasRecord, isFinish, err := genericCheckRecord("position", sqlStr)
	// Step 2: Exit if the record exists or an error occurs
	if hasRecord || !isFinish || err != nil {
		return
	}
	// Step 3: Insert a record for the system default positon "Default position" into the position table.
	sqlStr = `insert into position(id,name,description,creatorid) 
	values(10000,'Default position','System pre-set position',10000)`
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initPosition insert initvalue failed", zap.Error(err))
		return isFinish, err
	}
	return
}

// Get Position information by ID
func (p *Position) GetInfoByID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get Postion information from cache
	number, b, _ := cache.Get(pub.Position, p.ID)
	if number > 0 {
		json.Unmarshal(b, &p)
		return
	}
	// If Position information isn't in cache, retrieve it from database.
	sqlStr := `select name,description,status,createtime,creatorid,
	modifytime,modifierid,ts,dr
	from position 
	where id = $1`
	err = db.QueryRow(sqlStr, p.ID).Scan(&p.Name, &p.Description, &p.Status, &p.CreateDate, &p.Creator.ID,
		&p.ModifyDate, &p.Modifier.ID, &p.Ts, &p.Dr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Position.GetInfoByID db.QueryRow failed", zap.Error(err))
		return
	}
	// Get creator information.
	if p.Creator.ID > 0 {
		resStatus, err = p.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Modifier information.
	if p.Modifier.ID > 0 {
		resStatus, err = p.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Write into cache
	pB, _ := json.Marshal(p)
	cache.Set(pub.Position, p.ID, pB)

	return
}

// Get Position Master Data list
func GetPositionList() (ps []Position, resStatus i18n.ResKey, err error) {
	ps = make([]Position, 0)
	resStatus = i18n.StatusOK
	// Retrieve data from position table
	sqlStr := `select id,name,description,status,createtime,
	creatorid,modifytime,modifierid,ts,dr
	from position 
	where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetPositionList failed from database", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p Position
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.CreateDate,
			&p.Creator.ID, &p.ModifyDate, &p.Modifier.ID, &p.Ts, &p.Dr)
		if err != nil {
			zap.L().Error("GetPositionList from rows failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}

		// Get Creator details
		if p.Creator.ID > 0 {
			resStatus, err = p.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		// Get Modifier details
		if p.Modifier.ID > 0 {
			resStatus, err = p.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Append
		ps = append(ps, p)
	}

	return
}

// Get latest positon master data front-end cache
func (pc *PositionCache) GetOPsCache() (resStatus i18n.ResKey, err error) {
	pc.DelItems = make([]Position, 0)
	pc.NewItems = make([]Position, 0)
	pc.UpdateItems = make([]Position, 0)
	resStatus = i18n.StatusOK
	// Get the latest timestamp from the position table
	sqlStr := `select ts from position where ts > $1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, pc.QueryTs).Scan(&pc.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			pc.ResultNumber = 0
			pc.ResultTs = pc.QueryTs
			return
		}
		zap.L().Error("PositionCache.GetOPsCache query latest ts failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Retrieve all data greater than the QueryTs
	sqlStr = `select id,name,description,status,createtime,
	creatorid,modifytime,modifierid,ts,dr
	from position 
	where ts > $1 order by ts desc`
	rows, err := db.Query(sqlStr, pc.QueryTs)
	if err != nil {
		zap.L().Error("PositionCache.GetOPsCache get Cache from database failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p Position
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.CreateDate,
			&p.Creator.ID, &p.ModifyDate, &p.Modifier.ID, &p.Ts, &p.Dr)
		if err != nil {
			zap.L().Error("PositionCache.GetOPsCache rows.next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get Creator details
		if p.Creator.ID > 0 {
			resStatus, err = p.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get modifier details
		if p.Modifier.ID > 0 {
			resStatus, err = p.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		if p.Dr == 0 {
			if p.CreateDate.Before(pc.QueryTs) || p.CreateDate.Equal(pc.QueryTs) {
				pc.ResultNumber++
				pc.UpdateItems = append(pc.UpdateItems, p)
			} else {
				pc.ResultNumber++
				pc.NewItems = append(pc.NewItems, p)
			}
		} else {
			if p.CreateDate.Before(pc.QueryTs) || p.CreateDate.Equal(pc.QueryTs) {
				pc.ResultNumber++
				pc.DelItems = append(pc.DelItems, p)
			}
		}
	}

	return
}

// Add Position
func (p *Position) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the position name exists
	resStatus, err = p.CheckNameExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Insert a record to positon table
	sqlStr := `insert into position(name,description,status,creatorid) 
	values($1,$2,$3,$4) 
	returning id`
	err = db.QueryRow(sqlStr, p.Name, p.Description, p.Status, p.Creator.ID).Scan(&p.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Position.Add db.QueryRow failed", zap.Error(err))
		return
	}

	return
}

// Edit position
func (p *Position) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the position name exists
	resStatus, err = p.CheckNameExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update the record in the position table
	sqlStr := `update position set 
		name=$1,description=$2,status=$3,modifierid=$4,modifytime=current_timestamp,
		ts=current_timestamp 
		where id=$5 and ts=$6 and dr=0`
	res, err := db.Exec(sqlStr, p.Name, p.Description, p.Status, p.Modifier.ID,
		p.ID, p.Ts)
	if err != nil {
		zap.L().Error("Position.Edit db.exec failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Get the  number of rows affected by  the SQL statement update
	affected, err := res.RowsAffected()
	if err != nil {
		zap.L().Error("Position.Edit  get res.RowsAffected failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// If the number of affected rows is less than one,
	// it means that someone else has already modified the record
	if affected < 1 {
		zap.L().Info("Position.Edit failed,Other user are Editing")
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Delete from local cache
	p.DelFromLocalCache()
	return
}

// Delete Position
func (p *Position) Delete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the position id refrenced
	resStatus, err = p.CheckUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update the record in the position table
	sqlStr := `update position set dr=1,modifierid=$1,modifytime=current_timestamp,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`

	res, err := db.Exec(sqlStr, p.Modifier.ID, p.ID, p.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Position.Delete db.exec failed", zap.Error(err))
		return
	}
	// Get the number of rows affected by the SQL statement update
	effected, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Position.Delete res.RowsAffected failed", zap.Error(err))
		return
	}
	// If the number of affected rows is less than one,
	// it means that someone else has already modified the record.
	if effected < 1 {
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Delete from local cache
	p.DelFromLocalCache()

	return
}

// Chcek the position name exists
func (p *Position) CheckNameExist() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var count int32
	sqlStr := `select count(id) from position 
	where dr=0 and name=$1 and id <> $2`
	err = db.QueryRow(sqlStr, p.Name, p.ID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Position.CheckNameExist query failed", zap.Error(err))
		return
	}
	if count > 0 {
		resStatus = i18n.StatusPositionNameExist
		return
	}
	return
}

// Batch delete position
func DeleteOPs(ops *[]Position, modifyUserId int32) (statusCode i18n.ResKey, err error) {
	statusCode = i18n.StatusOK
	// begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		statusCode = i18n.StatusInternalError
		zap.L().Error("DeleteOPs db.begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	delSqlStr := `update position set dr=1,modifierid=$1,modifytime=current_timestamp,ts=current_timestamp 
		where id=$2 and dr=0 and ts=$3`
	stmt, err := tx.Prepare(delSqlStr)
	if err != nil {
		statusCode = i18n.StatusInternalError
		zap.L().Error("DeleteOPs tx.Prepare failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer stmt.Close()

	for _, p := range *ops {
		// Check if the postion id  referenced
		statusCode, err = p.CheckUsed()
		if statusCode != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
		// Update the record in postion table
		result, err1 := stmt.Exec(modifyUserId, p.ID, p.Ts)
		if err1 != nil {
			zap.L().Error("DeleteOPs stmt.exec failed", zap.Error(err))
			tx.Rollback()
			return i18n.StatusInternalError, err1
		}
		// Check the number of rows affected by SQL statement update
		affected, err2 := result.RowsAffected()
		if err2 != nil {
			zap.L().Error("DeleteOPs check RowsAffected failed", zap.Error(err))
			tx.Rollback()
			return i18n.StatusInternalError, err2
		}
		if affected < 1 {
			zap.L().Info("DeleteOPs other edit")
			_ = tx.Rollback()
			return i18n.StatusOtherEdit, nil
		}
		// Delete from local cache
		p.DelFromLocalCache()
	}
	return
}

// Delete postion from local cache
func (p *Position) DelFromLocalCache() {
	number, _, _ := cache.Get(pub.Position, p.ID)
	if number > 0 {
		cache.Del(pub.Position, p.ID)
	}
}

// Check if the position ID is refrenced
func (p *Position) CheckUsed() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Define the items to be checked.
	checkItems := []ArchiveCheckUsed{
		{
			Description:    "Refrenced by user",
			SqlStr:         "select count(id) as usednum from sysuser where dr=0 and positionid=$1",
			UsedReturnCode: i18n.StatusUserUsed,
		},
		{
			Description:    "Referenced by Personal Protective Equipment",
			SqlStr:         "select count(id) as usednum from ppequotas_h where dr=0 and positionid=$1",
			UsedReturnCode: i18n.StatusPQPositionUsed,
		},
	}
	// Check item by item
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, p.ID).Scan(&usedNum)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("Position.CheckIsUsed "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return
}
