package pg

import (
	"database/sql"
	"sccsmsserver/i18n"
	"time"

	"go.uber.org/zap"
)

// Construction Site Option
type ConstructionSiteOption struct {
	ID           int32              `db:"id" json:"id"`
	Code         string             `db:"code" json:"code"`
	Name         string             `db:"name" json:"name"`
	DisplayName  string             `db:"displayname" json:"displayName"`
	UDC          UserDefineCategory `db:"udcid" json:"udc"`
	DefaultValue UserDefinedArchive `db:"defaultvalueid" json:"defaultValue"`
	Enable       int16              `db:"enable" json:"enable"`
	IsModify     int16              `json:"isModify"`
	CreateDate   time.Time          `db:"createtime" json:"createDate"`
	Creator      Person             `db:"creatorid" json:"creator"`
	ModifyDate   time.Time          `db:"modifytime" json:"modifyDate"`
	Modifier     Person             `db:"modifierid" json:"modifier"`
	Ts           time.Time          `db:"ts" json:"ts"`
	Dr           int16              `db:"dr" json:"dr"`
}

// Construction Site Options Front-end Cache
type SceneItemOptionCache struct {
	QueryTs      time.Time                `json:"queryTs"`
	ResultNumber int32                    `json:"resultNumber"`
	DelItems     []ConstructionSiteOption `json:"delItems"`
	UpdateItems  []ConstructionSiteOption `json:"updateItems"`
	NewItems     []ConstructionSiteOption `json:"newItems"`
	ResultTs     time.Time                `json:"resultTs"`
}

// Initialize cso table
func initCSO() (isFinish bool, err error) {
	return
}

// Get Construction Site Options
func GetCSOs() (scOptions []ConstructionSiteOption, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	scOptions = make([]ConstructionSiteOption, 0)
	// Retrieve Construction Site Options from cso table
	sqlStr := `select id,code,name,displayname,udcid,
	defaultvalueid,enable,createtime,creatorid,modifytime,
	modifierid,ts,dr
	from cso order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetCSOs db.Query failed", zap.Error(err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var option ConstructionSiteOption
		err = rows.Scan(&option.ID, &option.Code, &option.Name, &option.DisplayName, &option.UDC.ID,
			&option.DefaultValue.ID, &option.Enable, &option.CreateDate, &option.Creator.ID, &option.ModifyDate,
			&option.Modifier.ID, &option.Ts, &option.Dr)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetCSOs rows.Next failed", zap.Error(err))
			return
		}
		// Get UDC details
		if option.UDC.ID > 0 {
			resStatus, err = option.UDC.GetUDCInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get DefaultValue details
		if option.DefaultValue.ID > 0 {
			resStatus, err = option.DefaultValue.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier details
		if option.Modifier.ID > 0 {
			resStatus, err = option.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Check if the value is allowed to be updated
		resStatus, err = option.CheckIsModify()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		scOptions = append(scOptions, option)
	}

	return
}

// Get Construction Site Options front-end cache
func (csoc *SceneItemOptionCache) GetCSOCache() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get the latest timestamp from cso table
	sqlStr := `select ts from cso where ts > $1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, csoc.QueryTs).Scan(&csoc.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			csoc.ResultNumber = 0
			csoc.ResultTs = csoc.QueryTs
			resStatus = i18n.StatusOK
			return
		}
		zap.L().Error("ConstructionSiteOption.GetCSOCache db.QueryRow failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Retrieve all data greater than the latest timestamp from cso table
	sqlStr = `select id,code,name,displayname,udcid,
	defaultvalueid,enable,createtime,creatorid,modifytime,
	modifierid,ts,dr
	from cso where ts>$1 order by ts desc`
	rows, err := db.Query(sqlStr, csoc.QueryTs)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ConstructionSiteOption.GetCSOCache db.Query failed", zap.Error(err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var option ConstructionSiteOption
		err = rows.Scan(&option.ID, &option.Code, &option.Name, &option.DisplayName, &option.UDC.ID,
			&option.DefaultValue.ID, &option.Enable, &option.CreateDate, &option.Creator.ID, &option.ModifyDate,
			&option.Modifier.ID, &option.Ts, &option.Dr)
		if err != nil {
			zap.L().Error("ConstructionSiteOption.GetCSOCache rows.Next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get UDC details
		if option.UDC.ID > 0 {
			resStatus, err = option.UDC.GetUDCInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Default Value details
		if option.DefaultValue.ID > 0 {
			resStatus, err = option.DefaultValue.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier details
		if option.Modifier.ID > 0 {
			resStatus, err = option.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Check if the value is allowed to be updated
		resStatus, err = option.CheckIsModify()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}

		if option.Dr == 0 {
			if option.CreateDate.Before(csoc.QueryTs) || option.CreateDate.Equal(csoc.QueryTs) {
				csoc.ResultNumber++
				csoc.UpdateItems = append(csoc.UpdateItems, option)
			} else {
				csoc.ResultNumber++
				csoc.NewItems = append(csoc.NewItems, option)
			}
		} else {
			if option.CreateDate.Before(csoc.QueryTs) || option.CreateDate.Equal(csoc.QueryTs) {
				csoc.ResultNumber++
				csoc.DelItems = append(csoc.DelItems, option)
			}
		}
	}

	return
}

// Check if the value is allowed to be updated
func (cso *ConstructionSiteOption) CheckIsModify() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var usedNumber int32
	sqlStr := `select count(id) as usednum from cs where dr=0 and ` + cso.Code + `>0`
	err = db.QueryRow(sqlStr).Scan(&usedNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ConstructionSiteOption.CheckIsModify db.QueryRow failed", zap.Error(err))
		return
	}

	if usedNumber > 0 {
		cso.IsModify = 1
	} else {
		cso.IsModify = 0
	}
	return
}

// Edit Construction Site Option
func (cso *ConstructionSiteOption) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Update the record in the cso table
	sqlStr := `update cso set displayname=$1,udcid=$2,defaultvalueid=$3 ,enable=$4,
	modifytime=current_timestamp, modifierid=$5,ts=current_timestamp 
	where id=$6 and ts=$7`
	res, err := db.Exec(sqlStr, cso.DisplayName, cso.UDC.ID, cso.DefaultValue.ID, cso.Enable,
		cso.Modifier.ID, cso.ID, cso.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ConstructionSiteOption.Edit db.Exec failed", zap.Error(err))
		return
	}

	// Check the number of rows affected by the SQL update operation
	affected, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("ConstructionSiteOption.Edit  res.RowsAffected failed", zap.Error(err))
		return
	}
	// if the number of updated rows is less than one,
	// it means that someone else already updated the record.
	if affected < 1 {
		resStatus = i18n.StatusOtherEdit
		return
	}

	return
}
