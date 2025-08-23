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

// User-Defined Archive
type UserDefinedArchive struct {
	ID          int32              `db:"id" json:"id"`
	UDC         UserDefineCategory `db:"udcid" json:"udc"`
	Code        string             `db:"code" json:"code"`
	Name        string             `db:"name" json:"name"`
	Description string             `db:"description" json:"description"`
	FatherId    int32              `db:"fatherid" json:"fatherID"`
	Status      int16              `db:"status" json:"status"`
	CreateDate  time.Time          `db:"createtime" json:"createDate"`
	Creator     Person             `db:"creatorid" json:"creator"`
	Modifier    Person             `db:"modifierid" json:"modifier"`
	ModifyDate  time.Time          `db:"modifytime" json:"modifyDate"`
	Ts          time.Time          `db:"ts" json:"ts"`
	Dr          int16              `db:"dr" json:"dr"`
}

// User-defined Archive front-end cache
type UDACache struct {
	QueryTs      time.Time            `json:"queryTs"`
	ResultNumber int32                `json:"resultNumber"`
	DelItems     []UserDefinedArchive `json:"delItems"`
	UpdateItems  []UserDefinedArchive `json:"updateItems"`
	NewItems     []UserDefinedArchive `json:"newItems"`
	ResultTs     time.Time            `json:"resultTs"`
}

// Get all User-defined archives under the UDC
func GetUDAList(udc *UserDefineCategory) (udaList []UserDefinedArchive, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	udaList = make([]UserDefinedArchive, 0)
	// Retrieve data from the uda table
	sqlStr := `select id,udcid,code,name,description,
	fatherid,status,createtime,creatorid,modifierid,
	modifytime,ts,dr 
	from uda where dr=0 and udcid=$1`
	rows, err := db.Query(sqlStr, udc.ID)
	if err != nil {
		zap.L().Error("GetUDAList failed from database", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()
	// Get UDA master data row by row
	for rows.Next() {
		var uda UserDefinedArchive
		err = rows.Scan(&uda.ID, &uda.UDC.ID, &uda.Code, &uda.Name, &uda.Description,
			&uda.FatherId, &uda.Status, &uda.CreateDate, &uda.Creator.ID, &uda.Modifier.ID,
			&uda.ModifyDate, &uda.Ts, &uda.Dr)
		if err != nil {
			zap.L().Error("GetUDAList from rows failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}

		// Change UDC
		uda.UDC = *udc
		// Get Creator detail
		if uda.Creator.ID > 0 {
			resStatus, err = uda.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier deatail
		if uda.Modifier.ID > 0 {
			resStatus, err = uda.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Append uda to udaList array
		udaList = append(udaList, uda)
	}
	return
}

// Get All User-defined Archive list
func GetUDAAll() (udaAll []UserDefinedArchive, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	udaAll = make([]UserDefinedArchive, 0)
	// Retrieve data from the uda table
	sqlStr := `select id,udcid,code,name,description,
	fatherID,status,createtime,creatorid,modifierid,
	modifytime,ts,dr 
	from uda where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetUDAAll failed from database", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()
	// Get UDA master data row by row
	for rows.Next() {
		var uda UserDefinedArchive
		err = rows.Scan(&uda.ID, &uda.UDC.ID, &uda.Code, &uda.Name, &uda.Description,
			&uda.FatherId, &uda.Status, &uda.CreateDate, &uda.Creator.ID, &uda.Modifier.ID,
			&uda.ModifyDate, &uda.Ts, &uda.Dr)
		if err != nil {
			zap.L().Error("GetUDAAll from rows failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get UDC detail
		if uda.UDC.ID > 0 {
			resStatus, err = uda.UDC.GetUDCInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get creator detail
		if uda.Creator.ID > 0 {
			resStatus, err = uda.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier detail
		if uda.Modifier.ID > 0 {
			resStatus, err = uda.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Append uda to udaAll array
		udaAll = append(udaAll, uda)
	}
	return
}

// Get latest User-defined Archice master data front-end cache
func (udac *UDACache) GetUDACache() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	udac.DelItems = make([]UserDefinedArchive, 0)
	udac.NewItems = make([]UserDefinedArchive, 0)
	udac.UpdateItems = make([]UserDefinedArchive, 0)
	// Get the latest  timestamp from uda table
	sqlStr := `select ts from uda where ts > $1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, udac.QueryTs).Scan(&udac.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			udac.ResultNumber = 0
			udac.ResultTs = udac.QueryTs
			return
		}
		zap.L().Error("UDACache.GetUDACache query latest ts failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	// Retrieve all data greater than the latest timestamp
	sqlStr = `select id,udcid,code,name,description,
	fatherid,status,createtime,creatorid,modifierid,
	modifytime,ts,dr 
	from uda where ts > $1 order by ts desc`
	rows, err := db.Query(sqlStr, udac.QueryTs)
	if err != nil {
		zap.L().Error("UDACache.GetUDDCahce get cache from database failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	for rows.Next() {
		var uda UserDefinedArchive
		err = rows.Scan(&uda.ID, &uda.UDC.ID, &uda.Code, &uda.Name, &uda.Description,
			&uda.FatherId, &uda.Status, &uda.CreateDate, &uda.Creator.ID, &uda.Modifier.ID,
			&uda.ModifyDate, &uda.Ts, &uda.Dr)
		if err != nil {
			zap.L().Error("UDACache.GetUDACache rows.Next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get UDC detail
		if uda.UDC.ID > 0 {
			resStatus, err = uda.UDC.GetUDCInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return
			}
		}
		// Get creator detail
		if uda.Creator.ID > 0 {
			resStatus, err = uda.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get modifier detail
		if uda.Modifier.ID > 0 {
			resStatus, err = uda.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		if uda.Dr == 0 {
			if uda.CreateDate.Before(udac.QueryTs) || uda.CreateDate.Equal(udac.QueryTs) {
				udac.ResultNumber++
				udac.UpdateItems = append(udac.UpdateItems, uda)
			} else {
				udac.ResultNumber++
				udac.NewItems = append(udac.NewItems, uda)
			}
		} else {
			if uda.CreateDate.Before(udac.QueryTs) || uda.CreateDate.Equal(udac.QueryTs) {
				udac.ResultNumber++
				udac.DelItems = append(udac.DelItems, uda)
			}
		}
	}
	return
}

// Add UDA
func (uda *UserDefinedArchive) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the UDA name exists
	resStatus, err = uda.CheckCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Insert a record into  the uda table
	sqlStr := `insert into uda(udcid,code,name,description,fatherID,
	status,creatorid)
	values($1,$2,$3,$4,$5,$6,$7) returning id`
	err = db.QueryRow(sqlStr, uda.UDC.ID, uda.Code, uda.Name, uda.Description, uda.FatherId,
		uda.Status, uda.Creator.ID).Scan(&uda.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("UserDefinedArchive.add db.QueryRow failed:", zap.Error(err))
		return
	}
	return
}

// Edit UDA
func (uda *UserDefinedArchive) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the UDA name exists
	resStatus, err = uda.CheckCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update the record in the uda table
	sqlStr := `update uda set
	udcid=$1,code=$2,name=$3,description=$4,fatherID=$5,
	status=$6,modifierid=$7,modifytime=current_timestamp,ts=current_timestamp
	where id=$8 and ts=$9 and dr=0`
	res, err := db.Exec(sqlStr, uda.UDC.ID, uda.Code, uda.Name, uda.Description, uda.FatherId,
		uda.Status,
		uda.Modifier.ID, uda.ID, uda.Ts)
	if err != nil {
		zap.L().Error("UserDefinedArchive.Edit db.Exec failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Get the number of rows affected by the SQL update operation
	affected, err := res.RowsAffected()
	if err != nil {
		zap.L().Error("UserDefinedArchive.Edit get res.RowsAffected failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// If the number of affected  rows is less than one,
	// it means that someone else has already modified the record.
	if affected < 1 {
		zap.L().Info("UserDefinedArchive.Edit failed,Other user are Editing")
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Delete from cache
	uda.DelFromLocalCache()
	return
}

// Delete UDA
func (uda *UserDefinedArchive) Delete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the UDA is refrenced
	resStatus, err = uda.CheckUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update the deletion flag for this record
	sqlStr := `update uda set dr=1,modifierid=$1,modifytime=current_timestamp,ts=current_timestamp
	where id=$2 and dr=0 and ts=$3`
	res, err := db.Exec(sqlStr, uda.Modifier.ID, uda.ID, uda.Ts)
	if err != nil {
		zap.L().Error("UserDefinedArchive.Delete db.Exec failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	// Check the number of rows affected by the SQL update operation
	affected, err := res.RowsAffected()
	if err != nil {
		zap.L().Error("UserDefinedArchive.Delete res.RowsAffected failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// if the number of affected rows is less than one,
	// it means that someone else has already modified the record.
	if affected < 1 {
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Delete from cache
	uda.DelFromLocalCache()
	return
}

// Batch delete UDA
func DeleteUDAs(udas *[]UserDefinedArchive, modifyUserId int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Start a databse transaction
	tx, err := db.Begin()
	if err != nil {
		zap.L().Error("DeleteUDAs db.begin failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer tx.Commit()

	delSqlStr := `update uda set dr=1,modifierid=$1,modifytime=current_timestamp,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	// Prepare to update
	stmt, err := tx.Prepare(delSqlStr)
	if err != nil {
		zap.L().Error("DeleteUDAs tx.Prepare failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer stmt.Close()

	for _, uda := range *udas {
		// Check if the UDA is refrenced.
		resStatus, err = uda.CheckUsed()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
		// Update the deletion flag for this record
		res, err1 := stmt.Exec(modifyUserId, uda.ID, uda.Ts)
		if err1 != nil {
			zap.L().Error("DeleteUDAs stmt.Exec failed:", zap.Error(err))
			tx.Rollback()
			resStatus = i18n.StatusInternalError
			return resStatus, err1
		}

		// Check the number of rows affected by the SQL update operation
		affected, err2 := res.RowsAffected()
		if err2 != nil {
			zap.L().Error("DeleteUDAs check res.RowsAffected failed", zap.Error(err))
			_ = tx.Rollback()
			return i18n.StatusInternalError, err2
		}

		if affected < 1 {
			zap.L().Info("DeleteUDAs other edit")
			_ = tx.Rollback()
			return i18n.StatusOtherEdit, nil
		}
		// Delete from cache
		uda.DelFromLocalCache()
	}

	return
}

// Check if the UDA Code exists
func (uda *UserDefinedArchive) CheckCodeExist() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var count int32
	sqlStr := `select count(id) from uda 
	where dr=0 and udcid = $1 and code=$2 and id <> $3`
	err = db.QueryRow(sqlStr, uda.UDC.ID, uda.Code, uda.ID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("UserDefinedArchive.CheckCodeExist QueryRow failed", zap.Error(err))
		return
	}
	if count > 0 {
		resStatus = i18n.StatusUDACodeExist
		return
	}
	return
}

// Check if the UDA ID is refrenced
func (uda *UserDefinedArchive) CheckUsed() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	checkItems := []ArchiveCheckUsed{
		/* {
			Description:    "执行项目档案默认值",
			SqlStr:         `select count(id) as usednumber from exectiveitem where resulttypeid='550' and dr=0 and defaultvalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEIDDefaultUsed,
		},
		{
			Description:    "执行项目档案错误值",
			SqlStr:         `select count(id) as usednumber from exectiveitem where resulttypeid='550' and dr=0 and errorvalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEIDErrorUsed,
		},
		{
			Description:    "执行模板默认值",
			SqlStr:         `select count(id) from exectivetemplate_b where eid_id in (select id from exectiveitem where resulttypeid='550' and dr=0) and dr=0 and defaultvalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEITDefaultUsed,
		},
		{
			Description:    "执行模板错误值",
			SqlStr:         `select count(id) from exectivetemplate_b where eid_id in (select id from exectiveitem where resulttypeid='550' and dr=0) and dr=0 and errorvalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEITErrorUsed,
		},
		{
			Description:    "被执行单执行值引用",
			SqlStr:         `select count(id) from executedoc_b where eid_id in (select id from exectiveitem where resulttypeid='550' and dr=0) and dr=0 and exectivevalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEDValueUsed,
		},
		{
			Description:    "被执行单错误值引用",
			SqlStr:         `select count(id) from executedoc_b where eid_id in (select id from exectiveitem where resulttypeid='550' and dr=0) and dr=0 and errorvalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEDErrorUsed,
		},
		{
			Description:    "被现场档案引用",
			SqlStr:         "select count(id) as usednum from sceneitem where dr=0 and (udf1=$1 or udf2=$1 or udf3=$1 or udf4=$1 or udf5=$1 or udf6=$1 or udf7=$1 or udf8=$1 or udf9=$1 or udf10=$1)",
			UsedReturnCode: i18n.StatusSIUsed,
		},
		{
			Description:    "被现场档案默认值引用",
			SqlStr:         "select count(id) as usednum from sceneitemoption where dr=0  and defaultvalue_id = $1",
			UsedReturnCode: i18n.StatusSIOUsed,
		}, */
	}
	// Check item by item
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, uda.ID).Scan(&usedNum)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("UserDefinedArchive.CheckIsUsed "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return i18n.StatusOK, nil
}

// Get UDA information by ID
func (uda *UserDefinedArchive) GetInfoByID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get detail from cache
	number, b, _ := cache.Get(pub.UDA, uda.ID)
	if number > 0 {
		json.Unmarshal(b, &uda)
		return
	}
	// If the UDA's information isn't in cache,
	// retrieve it from database.
	sqlStr := `select udcid,code,name,description,fatherID,
	status,	createtime,creatorid,modifierid,modifytime,
	ts,dr 
	from uda where id=$1`
	err = db.QueryRow(sqlStr, uda.ID).Scan(&uda.UDC.ID, &uda.Code, &uda.Name, &uda.Description, &uda.FatherId,
		&uda.Status, &uda.CreateDate, &uda.Creator.ID, &uda.Modifier.ID, &uda.ModifyDate,
		&uda.Ts, &uda.Dr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("UserDefinedArchive.GetInfoByID db.QueryRow failed:", zap.Error(err))
		return
	}
	// Get User-defined Category deatil
	if uda.UDC.ID > 0 {
		resStatus, err = uda.UDC.GetUDCInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Creator detail
	if uda.Creator.ID > 0 {
		resStatus, err = uda.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Modifier detail
	if uda.Modifier.ID > 0 {
		resStatus, err = uda.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Write into cache
	uddB, _ := json.Marshal(uda)
	cache.Set(pub.UDA, uda.ID, uddB)

	return i18n.StatusOK, nil
}

// Delete UDA from cache
func (uda *UserDefinedArchive) DelFromLocalCache() {
	number, _, _ := cache.Get(pub.UDA, uda.ID)
	if number > 0 {
		cache.Del(pub.UDA, uda.ID)
	}
}
