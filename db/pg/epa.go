package pg

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sccsmsserver/cache"
	"sccsmsserver/i18n"
	"sccsmsserver/pub"
	"time"

	"go.uber.org/zap"
)

func initEP() (isFinish bool, err error) {
	return
}

// Execution Project struct
type ExecutionProject struct {
	ID               int32              `db:"id" json:"id"`
	Code             string             `db:"code" json:"code"`
	Name             string             `db:"name" json:"name"`
	EPC              SimpEPC            `db:"epcid" json:"epc"`
	Description      string             `db:"description" json:"description"`
	Status           int16              `db:"status" json:"status"`
	ResultType       ScDataType         `db:"resulttypeid" json:"resultType"`
	UDC              UserDefineCategory `db:"udcid" json:"udc"`
	DefaultValue     string             `db:"defaultvalue" json:"defaultValue"`
	DefaultValueDisp string             `db:"defaultvaluedisp" json:"defaultValueDisp"`
	IsCheckError     int16              `db:"ischeckerror" json:"isCheckError"`
	ErrorValue       string             `db:"errorvalue" json:"errorValue"`
	ErrorValueDisp   string             `db:"errorvaluedisp" json:"errorValueDisp"`
	IsRequireFile    int16              `db:"isrequirefile" json:"isRequireFile"`
	IsOnsitePhoto    int16              `db:"isonsitephoto" json:"isOnSitePhoto"`
	RiskLevel        RiskLevel          `db:"risklevelid" json:"riskLevel"`
	CreateDate       time.Time          `db:"createtime" json:"createDate"`
	Creator          Person             `db:"creatorid" json:"creator"`
	ModifyDate       time.Time          `db:"modifytime" json:"modifyDate"`
	Modifier         Person             `db:"modifierid" json:"modifier"`
	Ts               time.Time          `db:"ts" json:"ts"`
	Dr               int16              `db:"dr" json:"dr"`
}

// Execution Project  front-end cache
type EPCache struct {
	QueryTs      time.Time          `json:"queryTs"`
	ResultNumber int32              `json:"resultNumber"`
	DelItems     []ExecutionProject `json:"delItems"`
	UpdateItems  []ExecutionProject `json:"updateItems"`
	NewItems     []ExecutionProject `json:"newItems"`
	ResultTs     time.Time          `json:"resultTs"`
}

// Get Execution Project List
func GetEPList() (epList []ExecutionProject, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	epList = make([]ExecutionProject, 0)
	// Retrieve data from epa table
	sqlStr := `select id,code,name,epcid,description,
	status,resulttypeid,udcid,defaultvalue,defaultvaluedisp,
	ischeckerror,errorvalue,errorvaluedisp,isrequirefile,isonsitephoto,
	risklevelid,createtime,creatorid,modifytime,modifierid,
	ts,dr 
	from epa 
	where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetEPList db.Query() failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()
	for rows.Next() {
		var ep ExecutionProject
		err = rows.Scan(&ep.ID, &ep.Code, &ep.Name, &ep.EPC.ID, &ep.Description,
			&ep.Status, &ep.ResultType.ID, &ep.UDC.ID, &ep.DefaultValue, &ep.DefaultValueDisp,
			&ep.IsCheckError, &ep.ErrorValue, &ep.ErrorValueDisp, &ep.IsRequireFile, &ep.IsOnsitePhoto,
			&ep.RiskLevel.ID, &ep.CreateDate, &ep.Creator.ID, &ep.ModifyDate, &ep.Modifier.ID,
			&ep.Ts, &ep.Dr)
		if err != nil {
			zap.L().Error("GetEPList rows.Next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get EPC detail
		if ep.EPC.ID > 0 {
			resStatus, err = ep.EPC.GetInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return
			}
		}
		// Get Result Type detail
		resStatus, err = ep.ResultType.GetDataTypeInfoByID()
		if err != nil || resStatus != i18n.StatusOK {
			return
		}
		// Get UDC detail
		if ep.UDC.ID > 0 {
			resStatus, err = ep.UDC.GetUDCInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return
			}
		}
		// Get Risk Level detail
		if ep.RiskLevel.ID > 0 {
			resStatus, err = ep.RiskLevel.GetRLInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return
			}
		}
		// Get Creator detail
		if ep.Creator.ID > 0 {
			resStatus, err = ep.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier detail
		if ep.Modifier.ID > 0 {
			resStatus, err = ep.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Append ep to EPList slice
		epList = append(epList, ep)
	}

	return
}

// Get Execution Project front-end cache
func (epac *EPCache) GetEPCache() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	epac.DelItems = make([]ExecutionProject, 0)
	epac.NewItems = make([]ExecutionProject, 0)
	epac.UpdateItems = make([]ExecutionProject, 0)
	// Query the latest timestamp from the epa table where the filed ts is greater than QueryTs
	sqlStr := `select ts from epa where ts > $1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, epac.QueryTs).Scan(&epac.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			epac.ResultNumber = 0
			epac.ResultTs = epac.QueryTs
			resStatus = i18n.StatusOK
			return
		}
		zap.L().Error("EPCache.GetEPCache query latest ts failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Retrieve all data greater than QueryTs
	sqlStr = `select id,code,name,epcid,description,
	status,resulttypeid,udcid,defaultvalue,defaultvaluedisp,
	ischeckerror,errorvalue,errorvaluedisp,isrequirefile,isonsitephoto,
	risklevelid,createtime,creatorid,modifytime,modifierid,
	ts,dr 
	from epa 
	where ts > $1 order by ts desc`
	rows, err := db.Query(sqlStr, epac.QueryTs)
	if err != nil {
		zap.L().Error("EPCache.GetEPCache db.Query failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()
	// Extract data from the database query results
	for rows.Next() {
		var ep ExecutionProject
		err = rows.Scan(&ep.ID, &ep.Code, &ep.Name, &ep.EPC.ID, &ep.Description,
			&ep.Status, &ep.ResultType.ID, &ep.UDC.ID, &ep.DefaultValue, &ep.DefaultValueDisp,
			&ep.IsCheckError, &ep.ErrorValue, &ep.ErrorValueDisp, &ep.IsRequireFile, &ep.IsOnsitePhoto,
			&ep.RiskLevel.ID, &ep.CreateDate, &ep.Creator.ID, &ep.ModifyDate, &ep.Modifier.ID,
			&ep.Ts, &ep.Dr)
		if err != nil {
			zap.L().Error("GetEPCache rows.Next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get EPC detail
		if ep.EPC.ID > 0 {
			resStatus, err = ep.EPC.GetInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return
			}
		}
		// Get Result Type detail
		resStatus, err = ep.ResultType.GetDataTypeInfoByID()
		if err != nil || resStatus != i18n.StatusOK {
			return
		}
		// Get UDC detail
		if ep.UDC.ID > 0 {
			resStatus, err = ep.UDC.GetUDCInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return
			}
		}
		// Get RIsk Level detail
		if ep.RiskLevel.ID > 0 {
			resStatus, err = ep.RiskLevel.GetRLInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return
			}
		}
		// Get creator detail
		if ep.Creator.ID > 0 {
			resStatus, err = ep.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get modifier detail
		if ep.Modifier.ID > 0 {
			resStatus, err = ep.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		if ep.Dr == 0 {
			if ep.CreateDate.Before(epac.QueryTs) || ep.CreateDate.Equal(epac.QueryTs) {
				epac.ResultNumber++
				epac.UpdateItems = append(epac.UpdateItems, ep)
			} else {
				epac.ResultNumber++
				epac.NewItems = append(epac.NewItems, ep)
			}
		} else {
			if ep.CreateDate.Before(epac.QueryTs) || ep.CreateDate.Equal(epac.QueryTs) {
				epac.ResultNumber++
				epac.DelItems = append(epac.DelItems, ep)
			}
		}
	}

	return
}

// Add Execution Project master data
func (ep *ExecutionProject) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the EP code exists
	resStatus, err = ep.CheckCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Add data to the epa table
	sqlStr := `insert into epa(code,name,epcid,description,status,
	resulttypeid,udcid,defaultvalue,defaultvaluedisp,ischeckerror,
	errorvalue,errorvaluedisp,isrequirefile,isonsitephoto,risklevelid,
	creatorid)
	values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
	returning id`
	err = db.QueryRow(sqlStr, ep.Code, ep.Name, ep.EPC.ID, ep.Description, ep.Status,
		ep.ResultType.ID, ep.UDC.ID, ep.DefaultValue, ep.DefaultValueDisp, ep.IsCheckError,
		ep.ErrorValue, ep.ErrorValueDisp, ep.IsRequireFile, ep.IsOnsitePhoto, ep.RiskLevel.ID,
		ep.Creator.ID).Scan(&ep.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionProject.Add db.QueryRow failed", zap.Error(err))
		return
	}

	return
}

// Edit Execution Project master data
func (ep *ExecutionProject) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the EP code exists
	resStatus, err = ep.CheckCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check if the EP id is refrenced
	var isUsed bool = false
	statusIsUsed, err := ep.CheckUsed()
	if err != nil {
		return
	}
	if statusIsUsed != i18n.StatusOK {
		isUsed = true
	}
	// Check if the Result Type has been modified
	var oldResultTypeID int32
	checkSql := "select resulttypeid from epa where id=$1 and dr=0"
	err = db.QueryRow(checkSql, ep.ID).Scan(&oldResultTypeID)
	if err != nil {
		zap.L().Error("ExecutionProject.Edit db.QueryRow(checkSql) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// The ResultType of an already applied EP  cannot be modified
	if isUsed && oldResultTypeID != ep.ResultType.ID {
		resStatus = i18n.StatusEPChangeResultType
		return
	}

	// Modify record in the epa table
	sqlStr := `update epa set code=$1,name=$2,epcid=$3,description=$4,status=$5,
	resulttypeid=$6,udcid=$7,defaultvalue=$8,defaultvaluedisp=$9,ischeckerror=$10,
	errorvalue=$11,errorvaluedisp=$12,isrequirefile=$13,isonsitephoto=$14,risklevelid=$15,
	modifytime=current_timestamp,modifierid=$16,ts=current_timestamp
	where id=$17 and dr = 0 and ts=$18`
	res, err := db.Exec(sqlStr, ep.Code, ep.Name, ep.EPC.ID, ep.Description, ep.Status,
		ep.ResultType.ID, ep.UDC.ID, ep.DefaultValue, ep.DefaultValueDisp, ep.IsCheckError,
		ep.ErrorValue, ep.ErrorValueDisp, ep.IsRequireFile, ep.IsOnsitePhoto, ep.RiskLevel.ID,
		ep.Modifier.ID,
		ep.ID, ep.Ts)
	if err != nil {
		zap.L().Error("ExecutionProject.Edit db.Exec failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Check if the number of rows affected by SQL update operation
	updateNumber, err := res.RowsAffected()
	if err != nil {
		zap.L().Error("ExecutionProject.Edit res.RowsAffected failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// If the number of affected rows is less than 1.
	// it means someone else already modified the record.
	if updateNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Delete from cache
	ep.DelFromLocalCache()

	return
}

// Get ScDataType Detail by id
func (sct *ScDataType) GetDataTypeInfoByID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get Value
	s, exists := ScDataTypeList[sct.ID]
	if exists {
		sct.TypeCode = s.TypeCode
		sct.TypeName = s.TypeName
		sct.DataType = s.DataType
		sct.FrontDb = s.FrontDb
		sct.InputMode = s.InputMode
	} else {
		resStatus = i18n.StatusInternalError
		err = errors.New("ScDataType not found")
		zap.L().Error(fmt.Sprintf(`%s%d%s`, "ScDateType GetDataTypeInfoByID Get ID=", sct.ID, "Failed:"), zap.Error(err))
		return
	}
	return
}

// Check if the Execution Project Code exists
func (ep *ExecutionProject) CheckCodeExist() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var count int32
	sqlStr := `select count(id) from epa where dr=0 and code=$1 and id <> $2`
	err = db.QueryRow(sqlStr, ep.Code, ep.ID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionProject.CheckCodeExist queryRow failed", zap.Error(err))
		return
	}

	if count > 0 {
		resStatus = i18n.StatusEPCodeExist
		return
	}

	return
}

// Delete Execution Project Master Data
func (ep *ExecutionProject) Delete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the EP is refrenced
	resStatus, err = ep.CheckUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update the delete flag of this EP in the epa table
	sqlStr := `update epa set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp where id=$2 and dr=0 and ts=$3`
	res, err := db.Exec(sqlStr, ep.Modifier.ID, ep.ID, ep.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionProject.Delete db.Exec failed", zap.Error(err))
		return
	}
	// Check the number of rows affected by SQL udpate operation
	affetced, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ExecutionProject.Delete res.RowsAffected failed", zap.Error(err))
		return
	}
	// If the number of affected rows if less than 1,
	// it means sonmeone else already modified this record.
	if affetced < 1 {
		resStatus = i18n.StatusOtherEdit
		zap.L().Info("ExectiveItemClass.Delete Other user edit")
		return
	}
	// Delete from cache
	ep.DelFromLocalCache()
	return
}

// Batch delete Execution Project master data
func DeleteEPs(epas *[]ExecutionProject, modifyUserID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteEids db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Prepare a SQL statement execution
	delSqlStr := `update epa set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp where id=$2 and dr=0 and ts=$3`
	stmt, err := tx.Prepare(delSqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteEPs tx.Prepare failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer stmt.Close()
	// Update the delete flag for each record one by one
	for _, ep := range *epas {
		// Check if the EP is referenced
		resStatus, err = ep.CheckUsed()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
		// Update the delete flag of this EP in epa table
		res, err := stmt.Exec(modifyUserID, ep.ID, ep.Ts)
		if err != nil {
			zap.L().Error("DeleteEPs stmt.Exec failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			tx.Rollback()
			return resStatus, err
		}

		// Check the number of rows affected by SQL update opeartion
		affectedRows, err := res.RowsAffected()
		if err != nil {
			zap.L().Error("DeleteEPs res.RowsAffected failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			tx.Rollback()
			return resStatus, err
		}
		// If the number of affected rows less than one,
		// it means that someone else has already modified this record.
		if affectedRows < 1 {
			resStatus = i18n.StatusOtherEdit
			zap.L().Info("DeleteEPs" + ep.Name + "has alreday modified by other.")
			tx.Rollback()
			return resStatus, nil
		}
		// Delete from cache
		ep.DelFromLocalCache()
	}

	return
}

// Get Execution Project detail by ID
func (ep *ExecutionProject) GetInfoByID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get detail from cache
	number, b, _ := cache.Get(pub.EPA, ep.ID)
	if number > 0 {
		json.Unmarshal(b, &ep)
		return
	}
	// If the EP not in cache, retrieve it from epa table
	sqlStr := `select code,name,epcid,description,status,
	resulttypeid,udcid,defaultvalue,defaultvaluedisp,ischeckerror,
	errorvalue,errorvaluedisp,isrequirefile,isonsitephoto,risklevelid,
	createtime,creatorid,modifytime,modifierid,ts,dr
	from epa where id = $1`
	err = db.QueryRow(sqlStr, ep.ID).Scan(&ep.Code, &ep.Name, &ep.EPC.ID, &ep.Description, &ep.Status,
		&ep.ResultType.ID, &ep.UDC.ID, &ep.DefaultValue, &ep.DefaultValueDisp, &ep.IsCheckError,
		&ep.ErrorValue, &ep.ErrorValueDisp, &ep.IsRequireFile, &ep.IsOnsitePhoto, &ep.RiskLevel.ID,
		&ep.CreateDate, &ep.Creator.ID, &ep.ModifyDate, &ep.Modifier.ID, &ep.Ts, &ep.Dr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ep.GetInfoByID db.QueryRow failed", zap.Error(err))
		return
	}
	// Get EPC detail
	if ep.EPC.ID > 0 {
		resStatus, err = ep.EPC.GetInfoByID()
		if err != nil || resStatus != i18n.StatusOK {
			return
		}
	}
	// Get ResultTpe detail
	resStatus, err = ep.ResultType.GetDataTypeInfoByID()
	if err != nil || resStatus != i18n.StatusOK {
		return
	}
	// Get UDC detail
	if ep.UDC.ID > 0 {
		resStatus, err = ep.UDC.GetUDCInfoByID()
		if err != nil || resStatus != i18n.StatusOK {
			return
		}
	}
	// Get Risk Level detail
	if ep.RiskLevel.ID > 0 {
		resStatus, err = ep.RiskLevel.GetRLInfoByID()
		if err != nil || resStatus != i18n.StatusOK {
			return
		}
	}
	// Get Creator detail
	if ep.Creator.ID > 0 {
		resStatus, err = ep.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Modifier detail
	if ep.Modifier.ID > 0 {
		resStatus, err = ep.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Write into cache
	epB, _ := json.Marshal(ep)
	cache.Set(pub.EPA, ep.ID, epB)

	return
}

// Delete Execution Project from cache
func (ep *ExecutionProject) DelFromLocalCache() {
	number, _, _ := cache.Get(pub.EPA, ep.ID)
	if number > 0 {
		cache.Del(pub.EPA, ep.ID)
	}
}

// Check if the Execution Project is referenced
func (ep *ExecutionProject) CheckUsed() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Define a list of items to be checked
	checkItems := []ArchiveCheckUsed{
		{
			Description:    "Referenced by Execution Project Template ",
			SqlStr:         `select count(id) as usednumber from ept_b where dr=0 and epaid=$1`,
			UsedReturnCode: i18n.StatusEPAUsed,
		},
		{
			Description:    "Referecnced by Execution Order",
			SqlStr:         `select count(id) as usednumber from executionorder_b where dr=0 and epaid=$1`,
			UsedReturnCode: i18n.StatusEOUsed,
		},
		{
			Description:    "Referenced by Issue Resolution Form",
			SqlStr:         `select count(id) as usednumber from issueresolutionform where dr=0 and epaid=$1`,
			UsedReturnCode: i18n.StatusIRFUsed,
		},
	}
	// Check one by one
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, ep.ID).Scan(&usedNum)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("ExecutionProject.CheckIsUsed "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return
}
