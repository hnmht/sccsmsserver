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

// Initialize cs table
func initCS() (isFinish bool, err error) {
	return
}

// Construction Site struct
type ConstructionSite struct {
	ID          int32              `db:"id" json:"id"`
	Code        string             `db:"code" json:"code"`
	Name        string             `db:"name" json:"name"`
	Description string             `db:"description" json:"description"`
	Csc         SimpCSC            `db:"cscid" json:"csc"`
	Department  SimpDept           `db:"subdeptid" json:"subDept"`
	RespDept    SimpDept           `db:"respdeptid" json:"respDept"`
	RespPerson  Person             `db:"resppersonid" json:"respPerson"`
	Status      int16              `db:"status" json:"status"`
	FinishFlag  int16              `db:"finishflag" json:"finishFlag"`
	FinishDate  string             `db:"finishdate" json:"finishDate"`
	Longitude   float64            `db:"longitude" json:"longitude"`
	Latitude    float64            `db:"latitude" json:"latitude"`
	Udf1        UserDefinedArchive `db:"udf1" json:"udf1"`
	Udf2        UserDefinedArchive `db:"udf2" json:"udf2"`
	Udf3        UserDefinedArchive `db:"udf3" json:"udf3"`
	Udf4        UserDefinedArchive `db:"udf4" json:"udf4"`
	Udf5        UserDefinedArchive `db:"udf5" json:"udf5"`
	Udf6        UserDefinedArchive `db:"udf6" json:"udf6"`
	Udf7        UserDefinedArchive `db:"udf7" json:"udf7"`
	Udf8        UserDefinedArchive `db:"udf8" json:"udf8"`
	Udf9        UserDefinedArchive `db:"udf9" json:"udf9"`
	Udf10       UserDefinedArchive `db:"udf10" json:"udf10"`
	CreateDate  time.Time          `db:"createtime" json:"createDate"`
	Creator     Person             `db:"creatorid" json:"creator"`
	ModifyDate  time.Time          `db:"modifytime" json:"modifyDate"`
	Modifier    Person             `db:"modifierid" json:"modifier"`
	Ts          time.Time          `db:"ts" json:"ts"`
	Dr          int16              `db:"dr" json:"dr"`
}

// Construction Site Front-end Cache
type ConstructionSiteCache struct {
	QueryTs      time.Time          `json:"queryTs"`
	ResultNumber int32              `json:"resultNumber"`
	DelItems     []ConstructionSite `json:"delItems"`
	UpdateItems  []ConstructionSite `json:"updateItems"`
	NewItems     []ConstructionSite `json:"newItems"`
	ResultTs     time.Time          `json:"resultTs"`
}

// Get Construction Site master data list
func GetCSs() (css []ConstructionSite, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	css = make([]ConstructionSite, 0)
	// Retrieve data from cs table
	sqlStr := `select id,code,name,description,cscid,
	subdeptid,respdeptid,resppersonid,status,finishflag,
	finishdate,	longitude,latitude,createtime,creatorid,
	udf1,udf2,udf3,udf4,udf5,
	udf6,udf7,udf8,udf9,udf10,
	modifytime,modifierid,ts,dr
	from cs where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetCSs db.Query failed", zap.Error(err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cs ConstructionSite
		err = rows.Scan(&cs.ID, &cs.Code, &cs.Name, &cs.Description, &cs.Csc.ID,
			&cs.Department.ID, &cs.RespDept.ID, &cs.RespPerson.ID, &cs.Status, &cs.FinishFlag,
			&cs.FinishDate, &cs.Longitude, &cs.Latitude, &cs.CreateDate, &cs.Creator.ID,
			&cs.Udf1.ID, &cs.Udf2.ID, &cs.Udf3.ID, &cs.Udf4.ID, &cs.Udf5.ID,
			&cs.Udf6.ID, &cs.Udf7.ID, &cs.Udf8.ID, &cs.Udf9.ID, &cs.Udf10.ID,
			&cs.ModifyDate, &cs.Modifier.ID, &cs.Ts, &cs.Dr)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetCSs rows.Next failed", zap.Error(err))
			return
		}
		// Get details
		resStatus, err = cs.GetAttachInfo()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		css = append(css, cs)
	}
	return
}

// Get latest Construction Site front-end cache
func (csc *ConstructionSiteCache) GetCSCache() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	csc.DelItems = make([]ConstructionSite, 0)
	csc.NewItems = make([]ConstructionSite, 0)
	csc.UpdateItems = make([]ConstructionSite, 0)
	// Get the latest timestamp from the cs table
	sqlStr := `select ts from cs where ts>$1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, csc.QueryTs).Scan(&csc.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			csc.ResultNumber = 0
			csc.ResultTs = csc.QueryTs
			resStatus = i18n.StatusOK
			return
		}
		zap.L().Error("ConstructionSiteCache.GetCSCache db.QueryRow failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	// Retrieve all data greater than latest timestamp from cs table
	sqlStr = `select id,code,name,description,cscid,
	subdeptid,respdeptid,resppersonid,status,finishflag,
	finishdate,longitude,latitude,createtime,creatorid,
	udf1,udf2,udf3,	udf4,udf5,
	udf6,udf7,udf8,udf9,udf10,
	modifytime,modifierid,ts,dr
	from cs where ts>$1 order by ts desc`
	rows, err := db.Query(sqlStr, csc.QueryTs)
	if err != nil {
		zap.L().Error("ConstructionSiteCache.GetCSCache db.Query failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cs ConstructionSite
		err = rows.Scan(&cs.ID, &cs.Code, &cs.Name, &cs.Description, &cs.Csc.ID,
			&cs.Department.ID, &cs.RespDept.ID, &cs.RespPerson.ID, &cs.Status, &cs.FinishFlag,
			&cs.FinishDate, &cs.Longitude, &cs.Latitude, &cs.CreateDate, &cs.Creator.ID,
			&cs.Udf1.ID, &cs.Udf2.ID, &cs.Udf3.ID, &cs.Udf4.ID, &cs.Udf5.ID,
			&cs.Udf6.ID, &cs.Udf7.ID, &cs.Udf8.ID, &cs.Udf9.ID, &cs.Udf10.ID,
			&cs.ModifyDate, &cs.Modifier.ID, &cs.Ts, &cs.Dr)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("ConstructionSiteCache.GetCSCache rows.Next failed", zap.Error(err))
			return
		}
		// Get details
		resStatus, err = cs.GetAttachInfo()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}

		if cs.Dr == 0 {
			if cs.CreateDate.Before(csc.QueryTs) || cs.CreateDate.Equal(csc.QueryTs) {
				csc.ResultNumber++
				csc.UpdateItems = append(csc.UpdateItems, cs)
			} else {
				csc.ResultNumber++
				csc.NewItems = append(csc.NewItems, cs)
			}
		} else {
			if cs.CreateDate.Before(csc.QueryTs) || cs.CreateDate.Equal(csc.QueryTs) {
				csc.ResultNumber++
				csc.DelItems = append(csc.DelItems, cs)
			}
		}
	}
	return
}

// Get CS master data details
func (cs *ConstructionSite) GetAttachInfo() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get Construction Site Category details
	if cs.Csc.ID > 0 {
		resStatus, err = cs.Csc.GetSCSCInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Department details
	if cs.Department.ID > 0 {
		resStatus, err = cs.Department.GetSimpDeptInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Responsible  Department details
	if cs.RespDept.ID > 0 {
		resStatus, err = cs.RespDept.GetSimpDeptInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Responsible Person details
	if cs.RespPerson.ID > 0 {
		resStatus, err = cs.RespPerson.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 1 details
	if cs.Udf1.ID > 0 {
		resStatus, err = cs.Udf1.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 2 details
	if cs.Udf2.ID > 0 {
		resStatus, err = cs.Udf2.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 3 details
	if cs.Udf3.ID > 0 {
		resStatus, err = cs.Udf3.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 4 details
	if cs.Udf4.ID > 0 {
		resStatus, err = cs.Udf4.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 51 details
	if cs.Udf5.ID > 0 {
		resStatus, err = cs.Udf5.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 6 details
	if cs.Udf6.ID > 0 {
		resStatus, err = cs.Udf6.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 7 details
	if cs.Udf7.ID > 0 {
		resStatus, err = cs.Udf7.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 8 details
	if cs.Udf8.ID > 0 {
		resStatus, err = cs.Udf8.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 9 details
	if cs.Udf9.ID > 0 {
		resStatus, err = cs.Udf9.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 10 details
	if cs.Udf10.ID > 0 {
		resStatus, err = cs.Udf10.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Creator details
	if cs.Creator.ID > 0 {
		resStatus, err = cs.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Modifier details
	if cs.Modifier.ID > 0 {
		resStatus, err = cs.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}

	return
}

// Get Construction Site details by ID
func (cs *ConstructionSite) GetInfoByID() (resStatus i18n.ResKey, err error) {
	// Get detail from cache
	number, b, _ := cache.Get(pub.CS, cs.ID)
	if number > 0 {
		json.Unmarshal(b, &cs)
		resStatus = i18n.StatusOK
		return
	}
	// If CS information not in cache, retrieve it from the cs table
	sqlStr := `select code,name,description,cscid,subdeptid,
	respdeptid,resppersonid,status,finishflag,finishdate,
	longitude,latitude,
	udf1,udf2,udf3,udf4,udf5,
	udf6,udf7,udf8,udf9,udf10,
	createtime,creatorid,modifytime,modifierid,ts,
	dr 
	from cs where id=$1`
	err = db.QueryRow(sqlStr, cs.ID).Scan(
		&cs.Code, &cs.Name, &cs.Description, &cs.Csc.ID, &cs.Department.ID,
		&cs.RespDept.ID, &cs.RespPerson.ID, &cs.Status, &cs.FinishFlag, &cs.FinishDate,
		&cs.Longitude, &cs.Latitude,
		&cs.Udf1.ID, &cs.Udf2.ID, &cs.Udf3.ID, &cs.Udf4.ID, &cs.Udf5.ID,
		&cs.Udf6.ID, &cs.Udf7.ID, &cs.Udf8.ID, &cs.Udf9.ID, &cs.Udf10.ID,
		&cs.CreateDate, &cs.Creator.ID, &cs.ModifyDate, &cs.Modifier.ID, &cs.Ts,
		&cs.Dr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ConstructionSite.GetInfoByID db.QueryRow failed", zap.Error(err))
		return
	}

	// Get detail
	resStatus, err = cs.GetAttachInfo()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Write into cache
	csB, _ := json.Marshal(cs)
	cache.Set(pub.CS, cs.ID, csB)

	return
}

// Add Construction Site
func (cs *ConstructionSite) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the CS Code exist
	resStatus, err = cs.CheckCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Write data into the cs table
	sqlStr := `insert into cs(code,name,description,cscid,subdeptid,
	respdeptid,resppersonid,status,finishflag,finishdate,
	longitude,latitude,udf1,udf2,udf3,
	udf4,udf5,udf6,udf7,udf8,
	udf9,udf10,creatorid,modifierid)
	values($1,$2,$3,$4,$5,
	$6,$7,$8,$9,$10,
	$11,$12,$13,$14,$15,
	$16,$17,$18,$19,$20,
	$21,$22,$23,$24)
	returning id`
	_, err = db.Exec(sqlStr, cs.Code, cs.Name, cs.Description, cs.Csc.ID, cs.Department.ID,
		cs.RespDept.ID, cs.RespPerson.ID, cs.Status, cs.FinishFlag, cs.FinishDate,
		cs.Longitude, cs.Latitude, cs.Udf1.ID, cs.Udf2.ID, cs.Udf3.ID,
		cs.Udf4.ID, cs.Udf5.ID, cs.Udf6.ID, cs.Udf7.ID, cs.Udf8.ID,
		cs.Udf9.ID, cs.Udf10.ID, cs.Creator.ID, cs.Modifier.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ConstructionSite.Add db.Exec failed", zap.Error(err))
		return
	}
	return
}

// Edit Construction Site
func (cs *ConstructionSite) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the CS code exists
	resStatus, err = cs.CheckCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update record in the cs table
	sqlStr := `update cs set code=$1,name=$2,description=$3,cscid=$4,subdeptid=$5,
	respdeptid=$6,resppersonid=$7,status=$8,finishflag=$9,finishdate=$10,
	longitude=$11,latitude=$12,udf1=$13,udf2=$14,udf3=$15,
	udf4=$16,udf5=$17,udf6=$18,udf7=$19,udf8=$20,
	udf9=$21,udf10=$22,modifierid=$23,
	modifytime=current_timestamp,ts=current_timestamp
	where id=$24 and ts=$25 and dr=0`
	res, err := db.Exec(sqlStr, cs.Code, cs.Name, cs.Description, cs.Csc.ID, cs.Department.ID,
		cs.RespDept.ID, cs.RespPerson.ID, cs.Status, cs.FinishFlag, cs.FinishDate,
		cs.Longitude, cs.Latitude, cs.Udf1.ID, cs.Udf2.ID, cs.Udf3.ID,
		cs.Udf4.ID, cs.Udf5.ID, cs.Udf6.ID, cs.Udf7.ID, cs.Udf8.ID,
		cs.Udf9.ID, cs.Udf10.ID, cs.Modifier.ID,
		cs.ID, cs.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ConstructionSite.Edit db.Exec failed", zap.Error(err))
		return
	}
	// Check the number of rows affected by SQL update operation
	affected, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ConstructionSite.Edit res.RowsAffected failed", zap.Error(err))
		return
	}
	// If the number of affected rows is less than one,
	// it means that someone else already updated the record.
	if affected < 1 {
		zap.L().Info("ConstructionSite.Edit failed,Other user are Editing")
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Delete from cache
	cs.DelFromLocalCache()
	return
}

// Delte Construction Site master data
func (cs *ConstructionSite) Delete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the CS ID is refrenced
	resStatus, err = cs.CheckUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update the delection field for this data in the cs table
	sqlStr := `update cs set dr=1,modifierid=$1,modifytime=current_timestamp,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	res, err := db.Exec(sqlStr, cs.Modifier.ID, cs.ID, cs.Ts)
	if err != nil {
		zap.L().Error("ConstructionSite.Delete db.Exec failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	// Check the number of rows affected by the SQL update operation
	affected, err := res.RowsAffected()
	if err != nil {
		zap.L().Error("ConstructionSite.Delete  res.RowsAffected failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// If the number of affected rows is less than one,
	// it means that someone else has already updated the record.
	if affected < 1 {
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Delete from cache
	cs.DelFromLocalCache()

	return i18n.StatusOK, nil
}

// Batch delete Construction Sites
func DeleteCSs(css *[]ConstructionSite, modifyUserId int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Begin a databese transaction
	tx, err := db.Begin()
	if err != nil {
		zap.L().Error("DeleteCSs db.Begin failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer tx.Commit()

	delSqlStr := `update cs set dr=1,modifierid=$1,modifytime=current_timestamp,ts=current_timestamp
	where id=$2 and dr=0 and ts=$3`
	// Update operation preprocessing
	stmt, err := tx.Prepare(delSqlStr)
	if err != nil {
		zap.L().Error("DeleteCSs tx.Prepare failed", zap.Error(err))
		resStatus = i18n.StatusOK
		tx.Rollback()
		return
	}
	defer stmt.Close()

	for _, cs := range *css {
		// Check if the Construction Site ID is refrenced.
		resStatus, err = cs.CheckUsed()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
		// Write the deletion field for this record in cs table
		res, err1 := stmt.Exec(modifyUserId, cs.ID, cs.Ts)
		if err != nil {
			zap.L().Error("DeleteCSs stmt.Exec failed", zap.Error(err1))
			resStatus = i18n.StatusInternalError
			tx.Rollback()
			return resStatus, err1
		}

		// Check the number of the rows affected by SQL update opeartion
		affected, err2 := res.RowsAffected()
		if err2 != nil {
			zap.L().Error("DeleteCSs check res.RowsAffected failed", zap.Error(err))
			_ = tx.Rollback()
			return i18n.StatusInternalError, err2
		}
		// If the number of affected rows is less than one,
		// it means that someone else has alreay updated the record.
		if affected < 1 {
			zap.L().Info("DeleteCSs other edit")
			_ = tx.Rollback()
			return i18n.StatusOtherEdit, nil
		}
		// Delete from cache
		cs.DelFromLocalCache()
	}
	return i18n.StatusOK, nil
}

// Check if the Construction Site code exists
func (cs *ConstructionSite) CheckCodeExist() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var count int32
	sqlStr := `select count(id) from cs where dr=0 and cscid=$1 and code=$2 and id <>$3`
	err = db.QueryRow(sqlStr, cs.Csc.ID, cs.Code, cs.ID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ConstructionSite.CheckCodeExist db.QueryRow failed", zap.Error(err))
		return
	}

	if count > 0 {
		resStatus = i18n.StatusCSCodeExist
		return
	}
	return
}

// Delete Construction Site from cache
func (cs *ConstructionSite) DelFromLocalCache() {
	number, _, _ := cache.Get(pub.CS, cs.ID)
	if number > 0 {
		cache.Del(pub.CS, cs.ID)
	}
}

// Check Construction Site ID is refrenced
func (cs *ConstructionSite) CheckUsed() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Defined the items that need to be checked
	checkItems := []ArchiveCheckUsed{
		/* {
			Description:    "Refrenced by Work Order",
			SqlStr:         `select count(id) as usednumber from workorder_b where dr=0 and si_id=$1`,
			UsedReturnCode: i18n.StatusWOUsed,
		},
		{
			Description:    "被执行单引用",
			SqlStr:         `select count(id) as usednumber from executedoc_h where dr=0 and si_id=$1`,
			UsedReturnCode: i18n.StatusEDUsed,
		},
		{
			Description:    "被问题处理单单引用",
			SqlStr:         `select count(id) as usednumber from disposedoc where dr=0 and si_id=$1`,
			UsedReturnCode: i18n.StatusDDUsed,
		}, */
	}
	// Check item by item
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, cs.ID).Scan(&usedNum)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("ConstructionSite.CheckUsed  "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return
}
