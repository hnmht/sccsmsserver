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

// Initialize csa table
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
	EndFlag     int16              `db:"endflag" json:"endFlag"`
	EndDate     time.Time          `db:"enddate" json:"endDate"`
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
	// Retrieve data from csa table
	sqlStr := `select id,code,name,description,cscid,
	subdeptid,respdeptid,resppersonid,status,endflag,
	enddate,longitude,latitude,createtime,creatorid,
	udf1,udf2,udf3,udf4,udf5,
	udf6,udf7,udf8,udf9,udf10,
	modifytime,modifierid,ts,dr
	from csa where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetCSs db.Query failed", zap.Error(err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var csa ConstructionSite
		err = rows.Scan(&csa.ID, &csa.Code, &csa.Name, &csa.Description, &csa.Csc.ID,
			&csa.Department.ID, &csa.RespDept.ID, &csa.RespPerson.ID, &csa.Status, &csa.EndFlag,
			&csa.EndDate, &csa.Longitude, &csa.Latitude, &csa.CreateDate, &csa.Creator.ID,
			&csa.Udf1.ID, &csa.Udf2.ID, &csa.Udf3.ID, &csa.Udf4.ID, &csa.Udf5.ID,
			&csa.Udf6.ID, &csa.Udf7.ID, &csa.Udf8.ID, &csa.Udf9.ID, &csa.Udf10.ID,
			&csa.ModifyDate, &csa.Modifier.ID, &csa.Ts, &csa.Dr)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetCSs rows.Next failed", zap.Error(err))
			return
		}
		// Get details
		resStatus, err = csa.GetAttachInfo()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		css = append(css, csa)
	}
	return
}

// Get latest Construction Site front-end cache
func (csc *ConstructionSiteCache) GetCSCache() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	csc.DelItems = make([]ConstructionSite, 0)
	csc.NewItems = make([]ConstructionSite, 0)
	csc.UpdateItems = make([]ConstructionSite, 0)
	// Get the latest timestamp from the csa table
	sqlStr := `select ts from csa where ts>$1 order by ts desc limit(1)`
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

	// Retrieve all data greater than the QueryTs from csa table
	sqlStr = `select id,code,name,description,cscid,
	subdeptid,respdeptid,resppersonid,status,endflag,
	enddate,longitude,latitude,createtime,creatorid,
	udf1,udf2,udf3,	udf4,udf5,
	udf6,udf7,udf8,udf9,udf10,
	modifytime,modifierid,ts,dr
	from csa where ts>$1 order by ts desc`
	rows, err := db.Query(sqlStr, csc.QueryTs)
	if err != nil {
		zap.L().Error("ConstructionSiteCache.GetCSCache db.Query failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	for rows.Next() {
		var csa ConstructionSite
		err = rows.Scan(&csa.ID, &csa.Code, &csa.Name, &csa.Description, &csa.Csc.ID,
			&csa.Department.ID, &csa.RespDept.ID, &csa.RespPerson.ID, &csa.Status, &csa.EndFlag,
			&csa.EndDate, &csa.Longitude, &csa.Latitude, &csa.CreateDate, &csa.Creator.ID,
			&csa.Udf1.ID, &csa.Udf2.ID, &csa.Udf3.ID, &csa.Udf4.ID, &csa.Udf5.ID,
			&csa.Udf6.ID, &csa.Udf7.ID, &csa.Udf8.ID, &csa.Udf9.ID, &csa.Udf10.ID,
			&csa.ModifyDate, &csa.Modifier.ID, &csa.Ts, &csa.Dr)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("ConstructionSiteCache.GetCSCache rows.Next failed", zap.Error(err))
			return
		}
		// Get details
		resStatus, err = csa.GetAttachInfo()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}

		if csa.Dr == 0 {
			if csa.CreateDate.Before(csc.QueryTs) || csa.CreateDate.Equal(csc.QueryTs) {
				csc.ResultNumber++
				csc.UpdateItems = append(csc.UpdateItems, csa)
			} else {
				csc.ResultNumber++
				csc.NewItems = append(csc.NewItems, csa)
			}
		} else {
			if csa.CreateDate.Before(csc.QueryTs) || csa.CreateDate.Equal(csc.QueryTs) {
				csc.ResultNumber++
				csc.DelItems = append(csc.DelItems, csa)
			}
		}
	}
	return
}

// Get CS master data details
func (csa *ConstructionSite) GetAttachInfo() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get Construction Site Category details
	if csa.Csc.ID > 0 {
		resStatus, err = csa.Csc.GetSCSCInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Department details
	if csa.Department.ID > 0 {
		resStatus, err = csa.Department.GetSimpDeptInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Responsible  Department details
	if csa.RespDept.ID > 0 {
		resStatus, err = csa.RespDept.GetSimpDeptInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Responsible Person details
	if csa.RespPerson.ID > 0 {
		resStatus, err = csa.RespPerson.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 1 details
	if csa.Udf1.ID > 0 {
		resStatus, err = csa.Udf1.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 2 details
	if csa.Udf2.ID > 0 {
		resStatus, err = csa.Udf2.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 3 details
	if csa.Udf3.ID > 0 {
		resStatus, err = csa.Udf3.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 4 details
	if csa.Udf4.ID > 0 {
		resStatus, err = csa.Udf4.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 51 details
	if csa.Udf5.ID > 0 {
		resStatus, err = csa.Udf5.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 6 details
	if csa.Udf6.ID > 0 {
		resStatus, err = csa.Udf6.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 7 details
	if csa.Udf7.ID > 0 {
		resStatus, err = csa.Udf7.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 8 details
	if csa.Udf8.ID > 0 {
		resStatus, err = csa.Udf8.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 9 details
	if csa.Udf9.ID > 0 {
		resStatus, err = csa.Udf9.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get User-defined filed 10 details
	if csa.Udf10.ID > 0 {
		resStatus, err = csa.Udf10.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Creator details
	if csa.Creator.ID > 0 {
		resStatus, err = csa.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Modifier details
	if csa.Modifier.ID > 0 {
		resStatus, err = csa.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}

	return
}

// Get Construction Site details by ID
func (csa *ConstructionSite) GetInfoByID() (resStatus i18n.ResKey, err error) {
	// Get detail from cache
	number, b, _ := cache.Get(pub.CSA, csa.ID)
	if number > 0 {
		json.Unmarshal(b, &csa)
		resStatus = i18n.StatusOK
		return
	}
	// If CS information not in cache, retrieve it from the csa table
	sqlStr := `select code,name,description,cscid,subdeptid,
	respdeptid,resppersonid,status,endflag,enddate,
	longitude,latitude,
	udf1,udf2,udf3,udf4,udf5,
	udf6,udf7,udf8,udf9,udf10,
	createtime,creatorid,modifytime,modifierid,ts,
	dr 
	from csa where id=$1`
	err = db.QueryRow(sqlStr, csa.ID).Scan(
		&csa.Code, &csa.Name, &csa.Description, &csa.Csc.ID, &csa.Department.ID,
		&csa.RespDept.ID, &csa.RespPerson.ID, &csa.Status, &csa.EndFlag, &csa.EndDate,
		&csa.Longitude, &csa.Latitude,
		&csa.Udf1.ID, &csa.Udf2.ID, &csa.Udf3.ID, &csa.Udf4.ID, &csa.Udf5.ID,
		&csa.Udf6.ID, &csa.Udf7.ID, &csa.Udf8.ID, &csa.Udf9.ID, &csa.Udf10.ID,
		&csa.CreateDate, &csa.Creator.ID, &csa.ModifyDate, &csa.Modifier.ID, &csa.Ts,
		&csa.Dr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ConstructionSite.GetInfoByID db.QueryRow failed", zap.Error(err))
		return
	}

	// Get detail
	resStatus, err = csa.GetAttachInfo()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Write into cache
	csaB, _ := json.Marshal(csa)
	cache.Set(pub.CSA, csa.ID, csaB)
	return
}

// Add Construction Site
func (csa *ConstructionSite) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK

	// Check if the EndDate field is a zero value to prevent writing
	// a zero value to the dateabase.
	if csa.EndDate.IsZero() {
		csa.EndDate = time.Now()
		// For CSA where work has been suspended,
		// the suspended date must be filled in.
		if csa.EndFlag == 1 {
			resStatus = i18n.StatusCSAEndDateRequired
			return
		}
	}
	// Check if the CS Code exist
	resStatus, err = csa.CheckCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Write data into the csa table
	sqlStr := `insert into csa(code,name,description,cscid,subdeptid,
	respdeptid,resppersonid,status,endflag,enddate,
	longitude,latitude,udf1,udf2,udf3,
	udf4,udf5,udf6,udf7,udf8,
	udf9,udf10,creatorid,modifierid)
	values($1,$2,$3,$4,$5,
	$6,$7,$8,$9,$10,
	$11,$12,$13,$14,$15,
	$16,$17,$18,$19,$20,
	$21,$22,$23,$24)
	returning id`
	_, err = db.Exec(sqlStr, csa.Code, csa.Name, csa.Description, csa.Csc.ID, csa.Department.ID,
		csa.RespDept.ID, csa.RespPerson.ID, csa.Status, csa.EndFlag, csa.EndDate,
		csa.Longitude, csa.Latitude, csa.Udf1.ID, csa.Udf2.ID, csa.Udf3.ID,
		csa.Udf4.ID, csa.Udf5.ID, csa.Udf6.ID, csa.Udf7.ID, csa.Udf8.ID,
		csa.Udf9.ID, csa.Udf10.ID, csa.Creator.ID, csa.Modifier.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ConstructionSite.Add db.Exec failed", zap.Error(err))
		return
	}
	return
}

// Edit Construction Site
func (csa *ConstructionSite) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the EndDate field is a zero value to prevent writing
	// a zero value to the dateabase.
	if csa.EndDate.IsZero() {
		csa.EndDate = time.Now()
		// For CSA where work has been suspended,
		// the suspended date must be filled in.
		if csa.EndFlag == 1 {
			resStatus = i18n.StatusCSAEndDateRequired
			return
		}
	}
	// Check if the CS code exists
	resStatus, err = csa.CheckCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update record in the csa table
	sqlStr := `update csa set code=$1,name=$2,description=$3,cscid=$4,subdeptid=$5,
	respdeptid=$6,resppersonid=$7,status=$8,endflag=$9,enddate=$10,
	longitude=$11,latitude=$12,udf1=$13,udf2=$14,udf3=$15,
	udf4=$16,udf5=$17,udf6=$18,udf7=$19,udf8=$20,
	udf9=$21,udf10=$22,modifierid=$23,
	modifytime=current_timestamp,ts=current_timestamp
	where id=$24 and ts=$25 and dr=0`
	res, err := db.Exec(sqlStr, csa.Code, csa.Name, csa.Description, csa.Csc.ID, csa.Department.ID,
		csa.RespDept.ID, csa.RespPerson.ID, csa.Status, csa.EndFlag, csa.EndDate,
		csa.Longitude, csa.Latitude, csa.Udf1.ID, csa.Udf2.ID, csa.Udf3.ID,
		csa.Udf4.ID, csa.Udf5.ID, csa.Udf6.ID, csa.Udf7.ID, csa.Udf8.ID,
		csa.Udf9.ID, csa.Udf10.ID, csa.Modifier.ID,
		csa.ID, csa.Ts)
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
	csa.DelFromLocalCache()
	return
}

// Delte Construction Site master data
func (csa *ConstructionSite) Delete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the CS ID is refrenced
	resStatus, err = csa.CheckUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update the delection field for this data in the csa table
	sqlStr := `update csa set dr=1,modifierid=$1,modifytime=current_timestamp,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	res, err := db.Exec(sqlStr, csa.Modifier.ID, csa.ID, csa.Ts)
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
	csa.DelFromLocalCache()

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

	delSqlStr := `update csa set dr=1,modifierid=$1,modifytime=current_timestamp,ts=current_timestamp
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

	for _, csa := range *css {
		// Check if the Construction Site ID is refrenced.
		resStatus, err = csa.CheckUsed()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
		// Write the deletion field for this record in csa table
		res, err1 := stmt.Exec(modifyUserId, csa.ID, csa.Ts)
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
		csa.DelFromLocalCache()
	}
	return i18n.StatusOK, nil
}

// Check if the Construction Site code exists
func (csa *ConstructionSite) CheckCodeExist() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var count int32
	sqlStr := `select count(id) from csa where dr=0 and cscid=$1 and code=$2 and id <>$3`
	err = db.QueryRow(sqlStr, csa.Csc.ID, csa.Code, csa.ID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ConstructionSite.CheckCodeExist db.QueryRow failed", zap.Error(err))
		return
	}

	if count > 0 {
		resStatus = i18n.StatusCSACodeExist
		return
	}
	return
}

// Delete Construction Site from cache
func (csa *ConstructionSite) DelFromLocalCache() {
	number, _, _ := cache.Get(pub.CSA, csa.ID)
	if number > 0 {
		cache.Del(pub.CSA, csa.ID)
	}
}

// Check Construction Site ID is refrenced
func (csa *ConstructionSite) CheckUsed() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Defined the items that need to be checked
	checkItems := []ArchiveCheckUsed{
		{
			Description:    "Refrenced by Work Order",
			SqlStr:         `select count(id) as usednumber from workorder_b where dr=0 and csaid=$1`,
			UsedReturnCode: i18n.StatusWOUsed,
		},
		{
			Description:    "Refreneced by Execution Order Header",
			SqlStr:         `select count(id) as usednumber from executionorder_h where dr=0 and csaid=$1`,
			UsedReturnCode: i18n.StatusEOUsed,
		},
		{
			Description:    "Refrenece by Issue Resolution Form",
			SqlStr:         `select count(id) as usednumber from issueresolutionform where dr=0 and csaid=$1`,
			UsedReturnCode: i18n.StatusIRFUsed,
		},
	}
	// Check item by item
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, csa.ID).Scan(&usedNum)
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
