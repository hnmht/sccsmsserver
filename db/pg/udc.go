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

// User-defined Category struct
type UserDefineCategory struct {
	ID          int32     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	IsLevel     int16     `db:"islevel" json:"isLevel"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Creator     Person    `db:"Creatorid" json:"creator"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	ModifyDate  time.Time `db:"modifytime" json:"modifyDate"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

// User-defined Category front-end cache struct
type UDCCache struct {
	QueryTs      time.Time            `json:"queryTS"`
	ResultNumber int32                `json:"resultNumber"`
	DelUDCs      []UserDefineCategory `json:"delItems"`
	UpdateUDCs   []UserDefineCategory `json:"updateItems"`
	NewUDCs      []UserDefineCategory `json:"newItems"`
	ResultTs     time.Time            `json:"resultTs"`
}

// Get User-defined Categroy master data list
func GetUDCList() (udcs []UserDefineCategory, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	udcs = make([]UserDefineCategory, 0)
	sqlStr := `select a.id,a.name,a.description,a.islevel,a.status,
	a.createtime,a.Creatorid,a.modifytime,a.modifierid,a.ts,
	a.dr
	from udc a
	where a.dr=0 order by a.ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetUDCList db.Query failed: ", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	for rows.Next() {
		var udc UserDefineCategory
		err = rows.Scan(&udc.ID, &udc.Name, &udc.Description, &udc.IsLevel, &udc.Status,
			&udc.CreateDate, &udc.Creator.ID, &udc.ModifyDate, &udc.Modifier.ID, &udc.Ts,
			&udc.Dr)
		if err != nil {
			zap.L().Error("GetUDCList rows.Next failed:", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get creator detail
		if udc.Creator.ID > 0 {
			resStatus, err = udc.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get modifier detail
		if udc.Modifier.ID > 0 {
			resStatus, err = udc.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Append
		udcs = append(udcs, udc)
	}

	return
}

// Get latest UDC front-end cache
func (udcc *UDCCache) GetUDCsCache() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	udcc.DelUDCs = make([]UserDefineCategory, 0)
	udcc.NewUDCs = make([]UserDefineCategory, 0)
	udcc.UpdateUDCs = make([]UserDefineCategory, 0)
	// Retrieve latest timestamp from udc table
	sqlStr := `select ts from udc where ts > $1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, udcc.QueryTs).Scan(&udcc.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			udcc.ResultNumber = 0
			udcc.ResultTs = udcc.QueryTs
			resStatus = i18n.StatusOK
			return
		}
		zap.L().Error("UDCCache.GetUDCsCache db.QueryRow failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	// Retrieve all data greater than the latest timestamp from usc table
	sqlStr = `select a.id,a.name,a.description,a.islevel,a.status,
	a.createtime,a.Creatorid,a.modifytime,a.modifierid,a.ts,
	a.dr 
	from udc a
	where a.ts > $1 order by a.ts desc`
	rows, err := db.Query(sqlStr, udcc.QueryTs)
	if err != nil {
		zap.L().Error("UDCCache.GetUDCsCache db.Query failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	for rows.Next() {
		var udc UserDefineCategory
		err = rows.Scan(&udc.ID, &udc.Name, &udc.Description, &udc.IsLevel, &udc.Status,
			&udc.CreateDate, &udc.Creator.ID, &udc.ModifyDate, &udc.Modifier.ID, &udc.Ts,
			&udc.Dr)
		if err != nil {
			zap.L().Error("UDCCache.GetUDCsCache rows.next failed:", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get Creator detail
		if udc.Creator.ID > 0 {
			resStatus, err = udc.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get modifier detail
		if udc.Modifier.ID > 0 {
			resStatus, err = udc.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		if udc.Dr == 0 {
			if udc.CreateDate.Before(udcc.QueryTs) || udc.CreateDate.Equal(udcc.QueryTs) {
				udcc.ResultNumber++
				udcc.UpdateUDCs = append(udcc.UpdateUDCs, udc)
			} else {
				udcc.ResultNumber++
				udcc.NewUDCs = append(udcc.NewUDCs, udc)
			}
		} else {
			if udc.CreateDate.Before(udcc.QueryTs) || udc.CreateDate.Equal(udcc.QueryTs) {
				udcc.ResultNumber++
				udcc.DelUDCs = append(udcc.DelUDCs, udc)
			}
		}
	}
	return
}

// Add User-defined Category
func (udfc *UserDefineCategory) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the udc name exist
	resStatus, err = udfc.CheckNameExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Write a record to the udc table
	sqlStr := `insert into udc(name,description,status,Creatorid) values($1,$2,$3,$4) returning id`
	err = db.QueryRow(sqlStr, udfc.Name, udfc.Description, udfc.Status, udfc.Creator.ID).Scan(&udfc.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("UserDefineCategory Add stmt.QueryRow failed", zap.Error(err))
		return
	}
	return
}

// Get User-defined Category detail by ID
func (udc *UserDefineCategory) GetUDCInfoByID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get UDC from cache
	number, b, _ := cache.Get(pub.UDC, udc.ID)
	if number > 0 {
		json.Unmarshal(b, &udc)
		return
	}
	// If UDC master data is't in cache, retrieve it from database
	sqlStr := `select a.name,a.description,a.islevel,a.status,a.createtime,
	a.Creatorid,a.modifytime,a.modifierid,a.ts,a.dr
	from udc a
	where a.id = $1`
	err = db.QueryRow(sqlStr, udc.ID).Scan(&udc.Name, &udc.Description, &udc.IsLevel, &udc.Status, &udc.CreateDate,
		&udc.Creator.ID, &udc.ModifyDate, &udc.Modifier.ID, &udc.Ts, &udc.Dr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("UserDefineCategory.GetUDCInfoByID db.QueryRow failed:", zap.Error(err))
		return
	}
	// Get creator detail
	if udc.Creator.ID > 0 {
		resStatus, err = udc.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get modifier detail
	if udc.Modifier.ID > 0 {
		resStatus, err = udc.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Write it into cache
	udcB, _ := json.Marshal(udc)
	cache.Set(pub.UDC, udc.ID, udcB)
	return
}

// Edit User-defined Category
func (udfc *UserDefineCategory) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the UDC name exists.
	resStatus, err = udfc.CheckNameExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update record in udc table
	sqlStr := `update udc set 
	name=$1,description=$2,islevel=$3,status=$4,modifierid=$5,
	modifytime=current_timestamp,ts=current_timestamp
	where id=$6 and ts=$7 and dr=0`
	res, err := db.Exec(sqlStr, udfc.Name, udfc.Description, udfc.IsLevel, udfc.Status, udfc.Modifier.ID,
		udfc.ID, udfc.Ts)
	if err != nil {
		zap.L().Error("UserDefineCategory.Edit db.exec failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Check the number of rows affected by update operation.
	affected, err := res.RowsAffected()
	if err != nil {
		zap.L().Error("UserDefineCategory.Edit  get res.RowsAffected failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// If the number of affected rows is less than 1,
	// it means that someone else has already updated the record.
	if affected < 1 {
		zap.L().Info("UserDefineCategory.Edit failed: another user is updating this data.")
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Delete from cache
	udfc.DelFromCache()
	return
}

// Delete User-define Category master data
func (udfc *UserDefineCategory) Delete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the data is referenced.
	resStatus, err = udfc.CheckUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update the deletion field for this data in the udc table.
	sqlStr := `update udc set dr=1,modifierid=$1,modifytime=current_timestamp,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`

	res, err := db.Exec(sqlStr, udfc.Modifier.ID, udfc.ID, udfc.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("UserDefineCategory.Delete stmt.exec failed", zap.Error(err))
		return
	}
	// Check the number of rows affected by update operation.
	affected, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("UserDefineCategory.Delete res.RowsAffected failed", zap.Error(err))
		return
	}
	if affected < 1 {
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Delete from local cache
	udfc.DelFromCache()
	return
}

// Check if the UDC name exists
func (udfc *UserDefineCategory) CheckNameExist() (resStatus i18n.ResKey, err error) {
	var count int32
	resStatus = i18n.StatusOK
	sqlStr := "select count(id) from udc where dr=0 and name=$1 and id <> $2"
	err = db.QueryRow(sqlStr, udfc.Name, udfc.ID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("UserDefineCategory.CheckNameExist query failed", zap.Error(err))
		return
	}
	if count > 0 {
		resStatus = i18n.StatusUDCNameExist
		return
	}
	return
}

// Batch delete UDC
func DeleteUDCs(udcs *[]UserDefineCategory, modifyUserId int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteUDCs db.begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	delSqlStr := `update udc set dr=1,modifierid=$1,modifytime=current_timestamp,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	stmt, err := tx.Prepare(delSqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteUDCs tx.Prepare failed", zap.Error(err))
		_ = tx.Rollback()
		return
	}
	defer stmt.Close()

	for _, udc := range *udcs {
		// Check the UDC ID is refrenced
		resStatus, err = udc.CheckUsed()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
		// Update the deletion field for this data in the udc table.
		result, err1 := stmt.Exec(modifyUserId, udc.ID, udc.Ts)
		if err1 != nil {
			zap.L().Error("DeleteUDCs stmt.exec failed", zap.Error(err))
			_ = tx.Rollback()
			return i18n.StatusInternalError, err1
		}
		// Check the number of rows affected by update operation
		affected, err2 := result.RowsAffected()
		if err2 != nil {
			zap.L().Error("DeleteUDCs check RowsAffected failed", zap.Error(err))
			_ = tx.Rollback()
			return i18n.StatusInternalError, err2
		}
		if affected < 1 {
			zap.L().Info("DeleteUDCs other edit")
			_ = tx.Rollback()
			return i18n.StatusOtherEdit, nil
		}
		// Delete item from cache
		udc.DelFromCache()
	}
	return
}

// Delete from cache
func (udc *UserDefineCategory) DelFromCache() {
	number, _, _ := cache.Get(pub.UDC, udc.ID)
	if number > 0 {
		cache.Del(pub.UDC, udc.ID)
	}
}

// Check if the UDC ID is refrenced
func (udc *UserDefineCategory) CheckUsed() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Define the items to be checked
	checkItems := []ArchiveCheckUsed{
		{
			Description:    "Refrenced by User-define Archive",
			SqlStr:         "select count(id) as usednum from ud where dr=0  and udcid = $1",
			UsedReturnCode: i18n.StatusUDUsed,
		},
		/* {
			Description:    "被执行项目默认值引用",
			SqlStr:         `select count(id) as usednum from exectiveitem where resulttypeid = '530' and dr=0 and defaultvalue=cast($1 as varchar)`,
			UsedReturnCode: i18n.StatusEIDDefaultUsed,
		},
		{
			Description:    "被执行项目错误值引用",
			SqlStr:         `select count(id) as usednum from exectiveitem where resulttypeid = '530' and dr=0 and errorvalue=cast($1 as varchar)`,
			UsedReturnCode: i18n.StatusEIDErrorUsed,
		},
		{
			Description:    "被执行模板默认值引用",
			SqlStr:         `select count(id) from exectivetemplate_b where eid_id in (select id from exectiveitem where resulttypeid='530' and dr=0) and dr=0 and defaultvalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEITDefaultUsed,
		},
		{
			Description:    "被执行模板错误值引用",
			SqlStr:         `select count(id) from exectivetemplate_b where eid_id in (select id from exectiveitem where resulttypeid='530' and dr=0) and dr=0 and errorvalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEITErrorUsed,
		},
		{
			Description:    "被执行单执行值引用",
			SqlStr:         `select count(id) from executedoc_b where eid_id in (select id from exectiveitem where resulttypeid='530' and dr=0) and dr=0 and exectivevalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEDValueUsed,
		},
		{
			Description:    "被执行单错误值引用",
			SqlStr:         `select count(id) from executedoc_b where eid_id in (select id from exectiveitem where resulttypeid='530' and dr=0) and dr=0 and errorvalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEDErrorUsed,
		},
		{
			Description:    "被现场档案自定义项引用",
			SqlStr:         "select count(id) as usednum from sceneitemoption where dr=0  and udc_id = $1",
			UsedReturnCode: i18n.StatusSIOUsed,
		}, */
	}
	// Check item by item
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, udc.ID).Scan(&usedNum)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("UserDefineCategory.CheckIsUsed "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return
}
