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

// Construction Site Category
type CSC struct {
	ID          int32     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Father      SimpCSC   `db:"fatherid" json:"fatherID"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Creator     Person    `db:"creatorid" json:"creator"`
	ModifyDate  time.Time `db:"modifyTime" json:"modifyDate"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

// Simple Construction Site Category
type SimpCSC struct {
	ID          int32     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	FatherID    int32     `db:"fatherid" json:"fatherID"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Creator     Person    `db:"creatorid" json:"creator"`
	ModifyDate  time.Time `db:"modifyTime" json:"modifyDate"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

// Simple Construction Site Category Front-end Cache
type SimpSICCache struct {
	QueryTs      time.Time `json:"queryTs"`
	ResultNumber int32     `json:"resultNumber"`
	DelItems     []SimpCSC `json:"delItems"`
	UpdateItems  []SimpCSC `json:"updateItems"`
	NewItems     []SimpCSC `json:"newItems"`
	ResultTs     time.Time `json:"resultTs"`
}

// Initialize csc table
func initCSC() (isFinish bool, err error) {
	return true, nil
}

// Get the SimpCSC information by ID
func (scsc *SimpCSC) GetSCSCInfoByID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get simplify information from cache
	number, b, _ := cache.Get(pub.SimpCSC, scsc.ID)
	if number > 0 {
		err = json.Unmarshal(b, &scsc)
		if err != nil {
			zap.L().Error("SimpCSC.GetSCSCInfoByID json.Unmarshal failed:", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		return
	}
	// If Simplify information isn't in cache, retrieve it from database.
	sqlStr := `select name,description,fatherid,status,
	createtime,creatorid,modifyTime,modifierid,ts,dr
	from csc where id=$1`
	err = db.QueryRow(sqlStr, scsc.ID).Scan(&scsc.Name, &scsc.Description, &scsc.FatherID, &scsc.Status,
		&scsc.CreateDate, &scsc.Creator.ID, &scsc.ModifyDate, &scsc.Modifier.ID, &scsc.Ts, &scsc.Dr)
	if err != nil {
		zap.L().Error("GetSCSCInfoByID db.QueryRow from database failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Get Creator information
	if scsc.Creator.ID > 0 {
		resStatus, err = scsc.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Modifier information
	if scsc.Modifier.ID > 0 {
		resStatus, err = scsc.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Write into cache
	ssicB, _ := json.Marshal(scsc)
	cache.Set(pub.SimpCSC, scsc.ID, ssicB)

	return
}

// Get Constructor Site Category master data list
func GetCSCList() (cscs []CSC, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	cscs = make([]CSC, 0)
	// Retrieve data list from database
	sqlStr := `select id,name,description,fatherid,status,
	createtime,creatorid,modifyTime,modifierid,ts,dr 
	from csc
	where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetCSCList db.Query failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	// Get CSC information row by row
	for rows.Next() {
		var csc CSC
		err = rows.Scan(&csc.ID, &csc.Name, &csc.Description, &csc.Father.ID, &csc.Status,
			&csc.CreateDate, &csc.Creator.ID, &csc.ModifyDate, &csc.Modifier.ID, &csc.Ts, &csc.Dr)
		if err != nil {
			zap.L().Error("GetCSCList row.Next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}

		// Get Father Category information
		if csc.Father.ID > 0 {
			resStatus, err = csc.Father.GetSCSCInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return
			}
		}

		// Get Creator information
		if csc.Creator.ID > 0 {
			resStatus, err = csc.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier information
		if csc.Modifier.ID > 0 {
			resStatus, err = csc.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Append csc to cscs
		cscs = append(cscs, csc)
	}
	return
}

// Get Simple Constructor Site Category master data list
func GetSimpCSCList() (scscs []SimpCSC, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	scscs = make([]SimpCSC, 0)
	// Retrieve Simple Constructor Site Category list from database
	sqlStr := `select id,name,description,fatherid,status,
	createtime,creatorid,modifyTime,modifierid,ts,
	dr 
	from csc
	where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetSimpCSCList db.Query failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()
	// Get SimpCSC information row by row
	for rows.Next() {
		var scsc SimpCSC
		err = rows.Scan(&scsc.ID, &scsc.Name, &scsc.Description, &scsc.FatherID, &scsc.Status,
			&scsc.CreateDate, &scsc.Creator.ID, &scsc.ModifyDate, &scsc.Modifier.ID, &scsc.Ts,
			&scsc.Dr)
		if err != nil {
			zap.L().Error("GetSimpCSCList row.Next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}

		// Get Creator information
		if scsc.Creator.ID > 0 {
			resStatus, err = scsc.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier information
		if scsc.Modifier.ID > 0 {
			resStatus, err = scsc.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Append SimpCSC to scscs
		scscs = append(scscs, scsc)
	}
	return
}

// Get latest Simple CSC master data for front-end cache
func (scscc *SimpSICCache) GetSimpCSCCache() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	scscc.DelItems = make([]SimpCSC, 0)
	scscc.NewItems = make([]SimpCSC, 0)
	scscc.UpdateItems = make([]SimpCSC, 0)
	// Get latest timestamp from csc table
	sqlStr := `select ts from csc where ts > $1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, scscc.QueryTs).Scan(&scscc.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			scscc.ResultNumber = 0
			scscc.ResultTs = scscc.QueryTs
			resStatus = i18n.StatusOK
			return
		}
		zap.L().Error("SimpCSCCache.GetSimpCSCCache query latest ts failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	// Retrieve all data greater than the laterst timestamp
	sqlStr = `select id,name,description,fatherid,status,
	createtime,creatorid,modifyTime,modifierid,ts,
	dr 
	from csc
	where ts > $1 order by ts desc`
	rows, err := db.Query(sqlStr, scscc.QueryTs)
	if err != nil {
		zap.L().Error("SimpCSCCache.GetSimpCSCCache get cache from database failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	for rows.Next() {
		var scsc SimpCSC
		err = rows.Scan(&scsc.ID, &scsc.Name, &scsc.Description, &scsc.FatherID, &scsc.Status,
			&scsc.CreateDate, &scsc.Creator.ID, &scsc.ModifyDate, &scsc.Modifier.ID, &scsc.Ts,
			&scsc.Dr)
		if err != nil {
			zap.L().Error("SimpCSCCache.GetSimpCSCCache rows.Next() failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}

		// Get Creator information
		if scsc.Creator.ID > 0 {
			resStatus, err = scsc.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier information
		if scsc.Modifier.ID > 0 {
			resStatus, err = scsc.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		if scsc.Dr == 0 {
			if scsc.CreateDate.Before(scscc.QueryTs) || scsc.CreateDate.Equal(scscc.QueryTs) {
				scscc.ResultNumber++
				scscc.UpdateItems = append(scscc.UpdateItems, scsc)
			} else {
				scscc.ResultNumber++
				scscc.NewItems = append(scscc.NewItems, scsc)
			}
		} else {
			if scsc.CreateDate.Before(scscc.QueryTs) || scsc.CreateDate.Equal(scscc.QueryTs) {
				scscc.ResultNumber++
				scscc.DelItems = append(scscc.DelItems, scsc)
			}
		}
	}
	return
}

// Check if the CSC name exist
func (csc *CSC) CheckNameExist() (resStatus i18n.ResKey, err error) {
	var count int32
	resStatus = i18n.StatusOK
	sqlStr := `select count(id) from csc where dr=0 and name=$1 and id <> $2`
	err = db.QueryRow(sqlStr, csc.Name, csc.ID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("CSC.CheckNameExist query failed", zap.Error(err))
		return
	}
	if count > 0 {
		resStatus = i18n.StatusCSCNameExist
		return
	}
	return
}

// Add CSC
func (csc *CSC) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the csc name exists
	resStatus, err = csc.CheckNameExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Write data into the csc table
	sqlStr := `insert into csc(name,description,fatherid,status,creatorid) 
	values($1,$2,$3,$4,$5) 
	returning id`
	err = db.QueryRow(sqlStr, csc.Name, csc.Description, csc.Father.ID, csc.Status, csc.Creator.ID).Scan(&csc.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("CSC.Add stmt.QueryRow failed", zap.Error(err))
		return
	}
	return
}

// Edit CSC
func (csc *CSC) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the parent CSC is compliant.
	if csc.Father.ID > 0 {
		if csc.ID == csc.Father.ID {
			resStatus = i18n.StatusCSCFatherSelf
			return
		}
		// The parent csc cannot be in a circular dependency
		cscs, res, err1 := GetSimpCSCList()
		if resStatus != i18n.StatusOK || err1 != nil {
			return res, err1
		}
		childrens := FindSimpSICChildrens(cscs, csc.ID)
		var number int32
		for _, child := range childrens {
			if child.ID == csc.Father.ID {
				number++
			}
		}

		if number > 0 {
			resStatus = i18n.StatusCSCFatherCircle
			return
		}
	}
	// Check if the csc name eists
	resStatus, err = csc.CheckNameExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update data in the csc table
	sqlStr := `update csc set
	name=$1,description=$2,fatherid=$3,status=$4,modifyTime=current_timestamp,
	modifierid=$5,ts=current_timestamp
	where id=$6 and dr=0 and ts=$7`
	res, err := db.Exec(sqlStr, csc.Name, csc.Description, csc.Father.ID, csc.Status,
		csc.Modifier.ID, csc.ID, csc.Ts)
	if err != nil {
		zap.L().Error("CSC.Edit stmt.Exec failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Check the number of rows updated by SQL
	updateNumber, err := res.RowsAffected()
	if err != nil {
		zap.L().Error("CSC.Edit res.RowsAffected failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	if updateNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Delete from localcache
	csc.DelFromLocalCache()
	return
}

// Check if the CSC have been refrenced
func (csc *CSC) CheckIsUsed() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Define the items to be checked.
	checkItems := []ArchiveCheckUsed{
		{
			Description:    "Refrenced by subCategory",
			SqlStr:         `select count(id) as usedNum from csc where fatherid=$1 and dr=0`,
			UsedReturnCode: i18n.StatusCSCLowLevelExist,
		},
		{
			Description:    "Refrenced by Construction Site",
			SqlStr:         `select count(id) as usednum from cs where cscid = $1 and dr=0`,
			UsedReturnCode: i18n.StatusCSUsed,
		},
		/*
			{
				Description:    "被执行项目默认值引用",
				SqlStr:         `select count(id) as usednum from cs where resulttypeid = '525' and dr=0 and defaultvalue=cast($1 as varchar)`,
				UsedReturnCode: i18n.StatusEIDDefaultUsed,
			},
			{
				Description:    "被执行项目错误值引用",
				SqlStr:         `select count(id) as usednum from exectiveitem where resulttypeid = '525' and dr=0 and errorvalue=cast($1 as varchar)`,
				UsedReturnCode: i18n.StatusEIDErrorUsed,
			},

				{
					Description:    "被执行模板默认值引用",
					SqlStr:         `select count(id) from exectivetemplate_b where eid_id in (select id from exectiveitem where resulttypeid='525' and dr=0) and dr=0 and defaultvalue=CAST($1 as varchar)`,
					UsedReturnCode: i18n.StatusEITDefaultUsed,
				},
				{
					Description:    "被执行模板错误值引用",
					SqlStr:         `select count(id) from exectivetemplate_b where eid_id in (select id from exectiveitem where resulttypeid='525' and dr=0) and dr=0 and errorvalue=CAST($1 as varchar)`,
					UsedReturnCode: i18n.StatusEITErrorUsed,
				},
				{
					Description:    "被执行单执行值引用",
					SqlStr:         `select count(id) from executedoc_b where eid_id in (select id from exectiveitem where resulttypeid='525' and dr=0) and dr=0 and exectivevalue=CAST($1 as varchar)`,
					UsedReturnCode: i18n.StatusEDValueUsed,
				},
				{
					Description:    "被执行单错误值引用",
					SqlStr:         `select count(id) from executedoc_b where eid_id in (select id from exectiveitem where resulttypeid='525' and dr=0) and dr=0 and errorvalue=CAST($1 as varchar)`,
					UsedReturnCode: i18n.StatusEDErrorUsed,
				}, */
	}
	// Check item by item
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, csc.ID).Scan(&usedNum)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("CSC.CheckIsUsed "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return
}

// Delete Construction Site Cateory
func (csc *CSC) Delete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if it's refrenced by other data.
	resStatus, err = csc.CheckIsUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update the deletion flag for the record
	sqlStr := `update csc set dr=1,modifyTime=current_timestamp,modifierid=$1,ts=current_timestamp where id=$2 and dr=0 and ts=$3`
	res, err := db.Exec(sqlStr, csc.Modifier.ID, csc.ID, csc.Ts)
	if err != nil {
		zap.L().Error("CSC.Delete stmt.Exec failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Check the number of the rows affected by the update operation.
	affected, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("CSC.Delete res.RowsAffedted failed", zap.Error(err))
		return
	}

	if affected < 1 {
		resStatus = i18n.StatusOtherEdit
		zap.L().Info("CSC.Delete Other User Edit")
		return
	}
	// Delete from the local cache
	csc.DelFromLocalCache()
	return
}

// Batch delete cscs
func DeleteCSCs(cscs *[]CSC, modifierID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Start a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteCSCs db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	delSqlStr := `update csc set dr=1,modifyTime=current_timestamp,modifierid=$1,ts=current_timestamp where id=$2 and dr=0 and ts=$3`

	stmt, err := tx.Prepare(delSqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteCSCs Delete prepare failed", zap.Error(err))
		tx.Rollback()
		return resStatus, err
	}
	defer stmt.Close()

	for _, csc := range *cscs {
		// Check if it's refrenced by other data
		resStatus, err = csc.CheckIsUsed()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
		// Write into database table
		res, err := stmt.Exec(modifierID, csc.ID, csc.Ts)
		if err != nil {
			zap.L().Error("DeleteCSCs stmt.Exec failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			tx.Rollback()
			return resStatus, err
		}
		// Check the number of the rows bye the updated operation.
		affectRows, err := res.RowsAffected()
		if err != nil {
			zap.L().Error("DeleteCSCs res.RowsAffected failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			tx.Rollback()
			return resStatus, err
		}

		if affectRows < 1 {
			resStatus = i18n.StatusOtherEdit
			zap.L().Info("DeleteCSCs: " + csc.Name + " other user eidting")
			tx.Rollback()
			return resStatus, nil
		}
		// Delete from local cache
		csc.DelFromLocalCache()
	}
	return
}

// Delete CSC from local cache
func (csc *CSC) DelFromLocalCache() {
	number, _, _ := cache.Get(pub.SimpCSC, csc.ID)
	if number > 0 {
		cache.Del(pub.SimpCSC, csc.ID)
	}
}

// Find all subcategories by ID
func FindSimpSICChildrens(scscs []SimpCSC, id int32) []SimpCSC {
	childrens := make([]SimpCSC, 0)
	for _, scsc := range scscs {
		if scsc.FatherID == id {
			childrens = append(childrens, scsc)
			child := FindSimpSICChildrens(scscs, scsc.ID)
			childrens = append(childrens, child...)
		}
	}
	return childrens
}
