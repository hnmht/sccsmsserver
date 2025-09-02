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

func initEPC() (isFinish bool, err error) {
	return
}

// Execution Project Category struct
type EPC struct {
	ID          int32     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Father      SimpEPC   `db:"fatherid" json:"fatherID"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createdate"`
	Creator     Person    `db:"creator" json:"creator"`
	ModifyDate  time.Time `db:"modifytime" json:"modifydate"`
	Modifier    Person    `db:"modifier" json:"modifier"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

// Simple Execution Project Category struct
type SimpEPC struct {
	ID          int32     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	FatherID    int32     `db:"fatherid" json:"fatherID"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createdate"`
	Creator     Person    `db:"creator" json:"creator"`
	ModifyDate  time.Time `db:"modifytime" json:"modifydate"`
	Modifier    Person    `db:"modifier" json:"modifier"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

// Simple EPC front-end cache
type SimpEPCCache struct {
	QueryTs      time.Time `json:"queryTs"`
	ResultNumber int32     `json:"resultNumber"`
	DelItems     []SimpEPC `json:"delItems"`
	UpdateItems  []SimpEPC `json:"updateItems"`
	NewItems     []SimpEPC `json:"newItems"`
	ResultTs     time.Time `json:"resultTs"`
}

// Get the SimpEPC information by ID
func (sepc *SimpEPC) GetInfoByID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get SimpEPC infomation from cache
	number, b, _ := cache.Get(pub.SimpEPC, sepc.ID)
	if number > 0 {
		json.Unmarshal(b, &sepc)
		return
	}
	// Retrieve SimpEPC information from the epc table
	sqlStr := `select name,description,fatherid,status,createtime,
	creator,modifytime,modifier,ts,dr
	from epc where id=$1`
	err = db.QueryRow(sqlStr, sepc.ID).Scan(&sepc.Name, &sepc.Description, &sepc.FatherID, &sepc.Status, &sepc.CreateDate,
		&sepc.Creator.ID, &sepc.ModifyDate, &sepc.Modifier.ID, &sepc.Ts, &sepc.Dr)
	if err != nil {
		zap.L().Error("SimpEPC.GetInfoByID db.QueryRow from database failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Get Creator information
	if sepc.Creator.ID > 0 {
		resStatus, err = sepc.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Modifier information
	if sepc.Modifier.ID > 0 {
		resStatus, err = sepc.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Write into cache
	sepcB, _ := json.Marshal(sepc)
	cache.Set(pub.SimpEPC, sepc.ID, sepcB)

	return
}

// Get Eexcution Project Archive list
func GetEPCList() (epcs []EPC, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	epcs = make([]EPC, 0)
	// Retrive EPC list from the epc table
	sqlStr := `select id,name,description,fatherid,status,
	createtime,creator,modifytime,modifier,ts,
	dr
	from epc 
	where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetEPCList db.Query failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()
	// Get EPC information row by row
	for rows.Next() {
		var epc EPC
		err = rows.Scan(&epc.ID, &epc.Name, &epc.Description, &epc.Father.ID, &epc.Status,
			&epc.CreateDate, &epc.Creator.ID, &epc.ModifyDate, &epc.Modifier.ID, &epc.Ts, &epc.Dr)
		if err != nil {
			zap.L().Error("GetEPCList row.Next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get Father information
		if epc.Father.ID > 0 {
			resStatus, err = epc.Father.GetInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return
			}
		}
		// Get Creator information
		if epc.Creator.ID > 0 {
			resStatus, err = epc.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier information
		if epc.Modifier.ID > 0 {
			resStatus, err = epc.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Append epc to epcs
		epcs = append(epcs, epc)
	}
	return
}

// Get Simple Execution Project Category list
func GetSimpEPCList() (sepcs []SimpEPC, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	sepcs = make([]SimpEPC, 0)
	// Retrieve SimpEPC list from the epc table
	sqlStr := `select id,name,description,fatherid,status,
	createtime,creator,modifytime,modifier,ts,
	dr 
	from epc 
	where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetSimpEPCList db.Query failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()
	// Get SimpEPC information row by row
	for rows.Next() {
		var sepc SimpEPC
		err = rows.Scan(&sepc.ID, &sepc.Name, &sepc.Description, &sepc.FatherID, &sepc.Status,
			&sepc.CreateDate, &sepc.Creator.ID, &sepc.ModifyDate, &sepc.Modifier.ID, &sepc.Ts,
			&sepc.Dr)
		if err != nil {
			zap.L().Error("GetSimpEPCList row.Next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get Creator information
		if sepc.Creator.ID > 0 {
			resStatus, err = sepc.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier information
		if sepc.Modifier.ID > 0 {
			resStatus, err = sepc.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Append sepc to sepcs
		sepcs = append(sepcs, sepc)
	}
	return
}

// Get latest Simple Execution Project Category for front-end cache
func (sepcc *SimpEPCCache) GetSimpEPCCache() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	sepcc.DelItems = make([]SimpEPC, 0)
	sepcc.NewItems = make([]SimpEPC, 0)
	sepcc.UpdateItems = make([]SimpEPC, 0)
	// Get the latest timestamp from epc table
	sqlStr := `select ts from epc where ts > $1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, sepcc.QueryTs).Scan(&sepcc.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			sepcc.ResultNumber = 0
			sepcc.ResultTs = sepcc.QueryTs
			return
		}
		zap.L().Error("SimpEPCCache.GetSimpEPCCache query latest ts failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	// Retrieve all data greater than latest timestamp
	sqlStr = `select id,name,description,fatherid,status,
	createtime,creator,modifytime,modifier,ts,
	dr 
	from epc 
	where ts > $1 order by ts desc`
	rows, err := db.Query(sqlStr, sepcc.QueryTs)
	if err != nil {
		zap.L().Error("SimpEPCCache.GetSimpEPCCache get cache from database failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	for rows.Next() {
		var sepc SimpEPC
		err = rows.Scan(&sepc.ID, &sepc.Name, &sepc.Description, &sepc.FatherID, &sepc.Status,
			&sepc.CreateDate, &sepc.Creator.ID, &sepc.ModifyDate, &sepc.Modifier.ID, &sepc.Ts,
			&sepc.Dr)
		if err != nil {
			zap.L().Error("SimpEPCCache.GetSimpEPCCache rows.Next() failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}

		// Get Creator information
		if sepc.Creator.ID > 0 {
			resStatus, err = sepc.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier information
		if sepc.Modifier.ID > 0 {
			resStatus, err = sepc.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		if sepc.Dr == 0 { // The data has not been deleted yet
			if sepc.CreateDate.Before(sepcc.QueryTs) || sepc.CreateDate.Equal(sepcc.QueryTs) { // The data has been modified
				sepcc.ResultNumber++
				sepcc.UpdateItems = append(sepcc.UpdateItems, sepc)
			} else { // Newly added data
				sepcc.ResultNumber++
				sepcc.NewItems = append(sepcc.NewItems, sepc)
			}
		} else { // The data has been deleted yet
			if sepc.CreateDate.Before(sepcc.QueryTs) || sepc.CreateDate.Equal(sepcc.QueryTs) { // Old deleted data
				sepcc.ResultNumber++
				sepcc.DelItems = append(sepcc.DelItems, sepc)
			}
			//Deleted new data, no need to process
		}
	}
	return
}

// Check if the epc name exists
func (epc *EPC) CheckNameExist() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var count int32
	sqlStr := `select count(id) from epc where dr=0 and name=$1 and id <> $2`
	err = db.QueryRow(sqlStr, epc.Name, epc.ID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPC.CheckNameExist query failed", zap.Error(err))
		return
	}
	if count > 0 {
		resStatus = i18n.StatusEPCNameExist
		return
	}
	return
}

// Add Execution Project Archive
func (epc *EPC) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the EPC name exists
	resStatus, err = epc.CheckNameExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Write data into the epc table
	sqlStr := `insert into epc(name,description,fatherid,status,creator)
	values($1,$2,$3,$4,$5)
	returning id`
	err = db.QueryRow(sqlStr, epc.Name, epc.Description, epc.Father.ID, epc.Status, epc.Creator.ID).Scan(&epc.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPC.Add stmt.QueryRow failed", zap.Error(err))
		return
	}

	return
}

// Edit Execution Project Archive
func (epc *EPC) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the parent EPC is compliant
	if epc.Father.ID > 0 {
		// The parent category cannot be its own
		if epc.ID == epc.Father.ID {
			resStatus = i18n.StatusEPCFatherSelf
			return
		}
		// Check for a circular dependency in the parent category
		sepcs, res, err1 := GetSimpEPCList()
		if resStatus != i18n.StatusOK || err != nil {
			return res, err1
		}
		childrens := FindSimpEPCChildrens(sepcs, epc.ID)
		var number int32
		for _, child := range childrens {
			if child.ID == epc.Father.ID {
				number++
			}
		}

		if number > 0 {
			resStatus = i18n.StatusEPCFatherCircle
			return
		}
	}
	// Check if the EPC name exists
	resStatus, err = epc.CheckNameExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Write the updates to the epc table
	sqlStr := `update epc set
	name=$1,description=$2,fatherid=$3,status=$4,modifytime=current_timestamp,
	modifier=$5,ts=current_timestamp 
	where id=$6 and dr=0 and ts=$7`

	res, err := db.Exec(sqlStr, epc.Name, epc.Description, epc.Father.ID, epc.Status,
		epc.Modifier.ID, epc.ID, epc.Ts)
	if err != nil {
		zap.L().Error("EPC.Edit stmt.Exec failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Check the number of affected rows by SQL update
	updateNumber, err := res.RowsAffected()
	if err != nil {
		zap.L().Error("EPC.Edit res.RowsAffected failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	if updateNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Delete from cache
	epc.DelFromLocalCache()
	return
}

// Check if the EPC ID is refrenced
func (epc *EPC) CheckIsUsed() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Define the items to be checked
	checkItems := []ArchiveCheckUsed{
		{
			Description:    "Refrenced by a sub-category",
			SqlStr:         `select count(id) as usedNum from epc where fatherid=$1 and dr=0`,
			UsedReturnCode: i18n.StatusEPCLowLevelExist,
		},
		{
			Description:    "Refrenced by Execution Project Archive",
			SqlStr:         `select count(id) as usednum from epa where epcid = $1 and dr=0`,
			UsedReturnCode: i18n.StatusEPAUsed,
		},
		{
			Description:    "Refrenced by EPA Default Value ",
			SqlStr:         `select count(id) as usednum from exectiveitem where resulttypeid = '540' and dr=0 and defaultvalue=cast($1 as varchar)`,
			UsedReturnCode: i18n.StatusEPDefaultUsed,
		},
		{
			Description:    "Refrenced by EPA Error Value",
			SqlStr:         `select count(id) as usednum from exectiveitem where resulttypeid = '540' and dr=0 and errorvalue=cast($1 as varchar)`,
			UsedReturnCode: i18n.StatusEPErrorUsed,
		},
		/* {
			Description:    "被执行模板默认值引用",
			SqlStr:         `select count(id) from exectivetemplate_b where eid_id in (select id from exectiveitem where resulttypeid='540' and dr=0) and dr=0 and defaultvalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEITDefaultUsed,
		},
		{
			Description:    "被执行模板错误值引用",
			SqlStr:         `select count(id) from exectivetemplate_b where eid_id in (select id from exectiveitem where resulttypeid='540' and dr=0) and dr=0 and errorvalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEITErrorUsed,
		},
		{
			Description:    "被执行单执行值引用",
			SqlStr:         `select count(id) from executedoc_b where eid_id in (select id from exectiveitem where resulttypeid='540' and dr=0) and dr=0 and exectivevalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEDValueUsed,
		},
		{
			Description:    "被执行单错误值引用",
			SqlStr:         `select count(id) from executedoc_b where eid_id in (select id from exectiveitem where resulttypeid='540' and dr=0) and dr=0 and errorvalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEDErrorUsed,
		}, */
	}
	// Item by item check
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, epc.ID).Scan(&usedNum)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("EPC.CheckIsUsed "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return
}

// EPC Delete 删除执行项目类别档案
func (epc *EPC) Delete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the EPC is refrenced
	resStatus, err = epc.CheckIsUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Write the updates to the epc table
	sqlStr := `update epc set dr=1,modifytime=current_timestamp,modifier=$1,ts=current_timestamp where id=$2 and dr=0 and ts=$3`
	res, err := db.Exec(sqlStr, epc.Modifier.ID, epc.ID, epc.Ts)
	if err != nil {
		zap.L().Error("EPC.Delete stmt.Exec failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Check the number of affected rows by SQL update
	affected, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPC.Delete res.RowsAffedted failed", zap.Error(err))
		return
	}

	if affected < 1 {
		resStatus = i18n.StatusOtherEdit
		zap.L().Info("EPC.Delete Other User Edit")
		return
	}
	// Delete from cache
	epc.DelFromLocalCache()

	return
}

// Batch delete EPCs
func DeleteEPCs(epcs *[]EPC, modifyUserID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteEPCs db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	delSqlStr := `update epc set dr=1,modifytime=current_timestamp,modifier=$1,ts=current_timestamp where id=$2 and dr=0 and ts=$3`
	stmt, err := tx.Prepare(delSqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteEPCs Delete prepare failed", zap.Error(err))
		tx.Rollback()
		return resStatus, err
	}
	defer stmt.Close()

	for _, epc := range *epcs {
		// Check if the EPC is refrenced
		resStatus, err = epc.CheckIsUsed()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
		// Write delete flag to the epc table
		res, err := stmt.Exec(modifyUserID, epc.ID, epc.Ts)
		if err != nil {
			zap.L().Error("DeleteEPCs stmt.Exec failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			tx.Rollback()
			return resStatus, err
		}
		// Check the number of affected rows by SQL update
		affectRows, err := res.RowsAffected()
		if err != nil {
			zap.L().Error("DeleteEPCs res.RowsAffected failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			tx.Rollback()
			return resStatus, err
		}
		if affectRows < 1 {
			resStatus = i18n.StatusOtherEdit
			zap.L().Info("DeleteEPCs " + epc.Name + " other user eidting")
			tx.Rollback()
			return resStatus, nil
		}
		// Delete from cache
		epc.DelFromLocalCache()
	}
	return
}

// Delete the EPC from cache
func (epc *EPC) DelFromLocalCache() {
	number, _, _ := cache.Get(pub.SimpEPC, epc.ID)
	if number > 0 {
		cache.Del(pub.SimpEPC, epc.ID)
	}
}

// Find all subcategories by ID
func FindSimpEPCChildrens(sepcs []SimpEPC, id int32) []SimpEPC {
	childrens := make([]SimpEPC, 0)
	for _, sepc := range sepcs {
		if sepc.FatherID == id {
			childrens = append(childrens, sepc)
			child := FindSimpEPCChildrens(sepcs, sepc.ID)
			childrens = append(childrens, child...)
		}
	}
	return childrens
}
