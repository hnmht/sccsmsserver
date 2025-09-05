package pg

func initEPA() (isFinish bool, err error) {
	return
}

/* // Execution Project Archive struct
type EPA struct {
	ID               int32              `db:"id" json:"id"`
	Code             string             `db:"code" json:"code"`
	Name             string             `db:"name" json:"name"`
	EPC              SimpEPC            `db:"epcid" json:"epc"`
	Description      string             `db:"description" json:"description"`
	Status           int16              `db:"status" json:"status"`
	ResultType       ScDataType         `db:"resulttypeid" json:"resultType"`
	UDCID            UserDefineCategory `db:"udcid" json:"udc"`
	DefaultValue     string             `db:"defaultvalue" json:"defaultValue"`
	DefaultValueDisp string             `db:"defaultvaluedisp" json:"defaultValueDisp"`
	IsCheckError     int16              `db:"ischeckerror" json:"isCheckError"`
	ErrorValue       string             `db:"errorvalue" json:"errorValue"`
	ErrorValueDisp   string             `db:"errorvaluedisp" json:"errorValueDisp"`
	IsRequireFile    int16              `db:"isrequirefile" json:"isRequireFile"`
	IsOnsitePhoto    int16              `db:"isonsitephoto" json:"isOnsitePhoto"`
	RiskLevel        RiskLevel          `db:"risklevelid" json:"riskLevel"`
	CreateDate       time.Time          `db:"createtime" json:"createDate"`
	Creator          Person             `db:"creatorid" json:"creator"`
	ModifyDate       time.Time          `db:"modifytime" json:"modifyDate"`
	Modifier         Person             `db:"modifierid" json:"modifier"`
	Ts               time.Time          `db:"ts" json:"ts"`
	Dr               int16              `db:"dr" json:"dr"`
}

// Execution Project Archive front-end cache
type EIDCahce struct {
	QueryTs      time.Time `json:"queryTs"`
	ResultNumber int32     `json:"resultNumber"`
	DelItems     []EPA     `json:"delItems"`
	UpdateItems  []EPA     `json:"updateItems"`
	NewItems     []EPA     `json:"newItems"`
	ResultTs     time.Time `json:"resultTs"`
}

// Get EPA List
func GetEPAList() (eidList []EPA, resStatus i18n.ResKey, err error) {
	eidList = make([]EPA, 0)
	sqlStr := `select id,code,name,epcid,description,
	status,resulttypeid,
	udcid,defaultvalue,defaultvaluedisp,ischeckerror,errorvalue,errorvaluedisp,isrequirefile,isonsitephoto,
	risklevelid,createtime,creatorid,modifytime,modifierid,
	ts,dr
	from exectiveitem where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetEPAList db.Query() failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	var resultNum int32
	//从查询记录中获取数据
	for rows.Next() {
		resultNum++
		var eid EPA
		err = rows.Scan(&eid.ID, &eid.Code, &eid.Name, &eid.EPC.ID, &eid.Description, &eid.Status, &eid.ResultType.ID,
			&eid.UDCID.ID, &eid.DefaultValue, &eid.DefaultValueDisp, &eid.IsCheckError, &eid.ErrorValue, &eid.ErrorValueDisp, &eid.IsRequireFile, &eid.IsOnsitePhoto,
			&eid.RiskLevel.ID, &eid.CreateDate, &eid.Creator.ID, &eid.ModifyDate, &eid.Modifier.ID,
			&eid.Ts, &eid.Dr)
		if err != nil {
			zap.L().Error("GetEPAList rows.Next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}

		//填充执行项目档案类别信息
		if eid.EPC.ID > 0 {
			resStatus, err = eid.EPC.GetEPCInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return
			}
		}

		//填充执行结果类型数据
		eid.ResultType.GetDataTypeInfoByID()

		//填充用户自定义档案类别
		if eid.UDCID.ID > 0 {
			resStatus, err = eid.UDCID.GetUDCInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return
			}
		}
		//填充风险等级
		if eid.RiskLevel.ID > 0 {
			resStatus, err = eid.RiskLevel.GetRLInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return
			}
		}
		//填充创建人信息
		if eid.Creator.ID > 0 {
			resStatus, err = eid.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充更新人信息
		if eid.Modifier.ID > 0 {
			resStatus, err = eid.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//追加数组
		eidList = append(eidList, eid)
	}
	//如果获取的列表数目为0
	if resultNum == 0 {
		resStatus = i18n.StatusResNoData
		return
	}

	resStatus = i18n.StatusOK
	return
}

// GetEIDCache 获取执行项目档案缓存
func (eidc *EIDCahce) GetEIDCache() (resStatus i18n.ResKey, err error) {
	eidc.DelItems = make([]EPA, 0)
	eidc.NewItems = make([]EPA, 0)
	eidc.UpdateItems = make([]EPA, 0)
	//查询执行项目档案最新缓存
	sqlStr := `select ts from exectiveitem where ts > $1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, eidc.QueryTs).Scan(&eidc.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			eidc.ResultNumber = 0
			eidc.ResultTs = eidc.QueryTs
			resStatus = i18n.StatusOK
			return
		}
		zap.L().Error("EIDCahce.GetEIDCache query latest ts failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	//查询所有大于QueryTs的数据
	sqlStr = `select id,code,name,epcid,description,status,resulttypeid,
	udcid,defaultvalue,defaultvaluedisp,ischeckerror,errorvalue,errorvaluedisp,isrequirefile,isonsitephoto,
	risklevelid,createtime,creatorid,modifytime,modifierid,
	ts,dr
	from exectiveitem where ts > $1 order by ts desc`
	rows, err := db.Query(sqlStr, eidc.QueryTs)
	if err != nil {
		zap.L().Error("EIDCache.GetEIDCache db.Query failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	for rows.Next() {
		var eid EPA
		err = rows.Scan(&eid.ID, &eid.Code, &eid.Name, &eid.EPC.ID, &eid.Description, &eid.Status, &eid.ResultType.ID,
			&eid.UDCID.ID, &eid.DefaultValue, &eid.DefaultValueDisp, &eid.IsCheckError, &eid.ErrorValue, &eid.ErrorValueDisp, &eid.IsRequireFile, &eid.IsOnsitePhoto,
			&eid.RiskLevel.ID, &eid.CreateDate, &eid.Creator.ID, &eid.ModifyDate, &eid.Modifier.ID,
			&eid.Ts, &eid.Dr)
		if err != nil {
			zap.L().Error("GetEIDCache rows.Next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}

		//填充执行项目档案类别信息
		if eid.EPC.ID > 0 {
			resStatus, err = eid.EPC.GetInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return
			}
		}
		//填充执行结果类型数据
		eid.ResultType.GetDataTypeInfoByID()

		//填充用户自定义档案类别
		if eid.UDCID.ID > 0 {
			resStatus, err = eid.UDCID.GetUDCInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return
			}
		}
		//填充风险等级
		if eid.RiskLevel.ID > 0 {
			resStatus, err = eid.RiskLevel.GetRLInfoByID()
			if err != nil || resStatus != i18n.StatusOK {
				return
			}
		}
		//填充创建人信息
		if eid.Creator.ID > 0 {
			resStatus, err = eid.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充更新人信息
		if eid.Modifier.ID > 0 {
			resStatus, err = eid.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		if eid.Dr == 0 {
			if eid.CreateDate.Before(eidc.QueryTs) || eid.CreateDate.Equal(eidc.QueryTs) {
				eidc.ResultNumber++
				eidc.UpdateItems = append(eidc.UpdateItems, eid)
			} else {
				eidc.ResultNumber++
				eidc.NewItems = append(eidc.NewItems, eid)
			}
		} else {
			if eid.CreateDate.Before(eidc.QueryTs) || eid.CreateDate.Equal(eidc.QueryTs) {
				eidc.ResultNumber++
				eidc.DelItems = append(eidc.DelItems, eid)
			}
		}
	}
	return i18n.StatusOK, nil
}

// EPA Add 增加执行项目
func (eid *EPA) Add() (resStatus i18n.ResKey, err error) {
	//检查编码是否重复
	resStatus, err = eid.CheckCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	//预处理
	sqlStr := `insert into exectiveitem(code,name,epcid,description,status,
	resulttypeid,udcid,defaultvalue,defaultvaluedisp,ischeckerror,
	errorvalue,errorvaluedisp,isrequirefile,isonsitephoto,risklevelid,
	creatorid)
	values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
	returning id`
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPA.Add db.Prepare failed", zap.Error(err))
		return
	}
	defer stmt.Close()

	//写入
	err = stmt.QueryRow(eid.Code, eid.Name, eid.EPC.ID, eid.Description, eid.Status,
		eid.ResultType.ID, eid.UDCID.ID, eid.DefaultValue, eid.DefaultValueDisp, eid.IsCheckError,
		eid.ErrorValue, eid.ErrorValueDisp, eid.IsRequireFile, eid.IsOnsitePhoto, eid.RiskLevel.ID,
		eid.Creator.ID).Scan(&eid.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPA.Add stmt.QueryRow failed", zap.Error(err))
		return
	}

	return
}

// EPA Edit 修改执行项目
func (eid *EPA) Edit() (resStatus i18n.ResKey, err error) {
	//检查执行项目编码是否重复
	resStatus, err = eid.CheckCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	//检查是否被引用
	var isUsed bool = false
	statusIsUsed, err := eid.CheckUsed()
	if err != nil {
		return
	}
	if statusIsUsed != i18n.StatusOK {
		isUsed = true
	}
	//检查是否修改了结果类型
	var oldResultTypeID int32
	checkSql := "select resulttypeid from exectiveitem where id=$1 and dr=0"
	err = db.QueryRow(checkSql, eid.ID).Scan(&oldResultTypeID)
	if err != nil {
		zap.L().Error("EPA.Edit db.QueryRow(checkSql) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	//已经被引用的执行项目不能修改结果类型
	if isUsed && oldResultTypeID != eid.ResultType.ID {
		resStatus = i18n.StatusEIDChangeResultType
		return
	}

	//预处理
	sqlStr := `update exectiveitem set code = $1,name = $2,epcid =$3,description=$4,status=$5,
	resulttypeid =$6,udcid =$7,defaultvalue =$8,defaultvaluedisp = $9,ischeckerror = $10,
	errorvalue = $11,errorvaluedisp = $12,isrequirefile =$13,isonsitephoto = $14,risklevelid=$15,
	modifytime = current_timestamp,modifierid = $16,ts=current_timestamp
	where id=$17 and dr = 0 and ts=$18`
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		zap.L().Error("EPA.Edit db.Prepare failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer stmt.Close()

	//写入
	res, err := stmt.Exec(eid.Code, eid.Name, eid.EPC.ID, eid.Description, eid.Status,
		eid.ResultType.ID, eid.UDCID.ID, eid.DefaultValue, eid.DefaultValueDisp, eid.IsCheckError,
		eid.ErrorValue, eid.ErrorValueDisp, eid.IsRequireFile, eid.IsOnsitePhoto, eid.RiskLevel.ID,
		eid.Modifier.ID,
		eid.ID, eid.Ts)
	if err != nil {
		zap.L().Error("EPA.Edit stmt.Exec failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	//检查更新的行数
	updateNumber, err := res.RowsAffected()
	if err != nil {
		zap.L().Error("EPA.Edit res.RowsAffected failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	if updateNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		return
	}
	//从localCache删除
	eid.DelFromLocalCache()
	//返回
	return
}

// ScDataType GetDataTypeInfoByID 根据ID获取数据类型信息
func (sct *ScDataType) GetDataTypeInfoByID() (resStatus i18n.ResKey, err error) {
	var s = ScDataTypeList[sct.ID]
	sct.TypeCode = s.TypeCode
	sct.TypeName = s.TypeName
	sct.DataType = s.DataType
	sct.FrontDb = s.FrontDb
	sct.InputMode = s.InputMode
	return
}

// EPA CheckCodeExist 检查执行项目编码是否存在
func (eid *EPA) CheckCodeExist() (resStatus i18n.ResKey, err error) {
	var count int32
	sqlStr := `select count(id) from exectiveitem where dr=0 and code=$1 and id <> $2`
	err = db.QueryRow(sqlStr, eid.Code, eid.ID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPA.CheckCodeExist queryRow failed", zap.Error(err))
		return
	}

	if count > 0 {
		resStatus = i18n.StatusEIDCodeExist
		return
	}

	resStatus = i18n.StatusOK
	return
}

// EPA Delete 删除执行项目档案
func (eid *EPA) Delete() (resStatus i18n.ResKey, err error) {
	//检查是否被引用
	resStatus, err = eid.CheckUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	//准备
	sqlStr := `update exectiveitem set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp where id=$2 and dr=0 and ts=$3`
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPA.Delete db.Prepare failed", zap.Error(err))
		return
	}
	defer stmt.Close()
	//执行删除操作
	res, err := stmt.Exec(eid.Modifier.ID, eid.ID, eid.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPA.Delete stmt.Exec failed", zap.Error(err))
		return
	}
	//检查删除操作影响的行数
	affetced, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("EPA.Delete res.RowsAffected failed", zap.Error(err))
		return
	}

	if affetced < 1 {
		resStatus = i18n.StatusOtherEdit
		zap.L().Info("ExectiveItemClass.Delete Other user edit")
		return
	}
	//从localCache删除
	eid.DelFromLocalCache()

	return
}

// DeleteEIDs 批量删除执行项目档案
func DeleteEIDs(eids *[]EPA, modifyUserID int32) (resStatus i18n.ResKey, err error) {
	//开始执行事务
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteEids db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	delSqlStr := `update exectiveitem set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp where id=$2 and dr=0 and ts=$3`
	//准备
	stmt, err := tx.Prepare(delSqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteEIDs tx.Prepare failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer stmt.Close()
	//逐项执行删除
	for _, eid := range *eids {
		//检查是否被引用
		resStatus, err = eid.CheckUsed()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
		//执行删除操作
		res, err := stmt.Exec(modifyUserID, eid.ID, eid.Ts)
		if err != nil {
			zap.L().Error("DeleteEIDs stmt.Exec failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			tx.Rollback()
			return resStatus, err
		}

		//检查操作影响的行数
		affectedRows, err := res.RowsAffected()
		if err != nil {
			zap.L().Error("DeleteEIDs res.RowsAffected failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			tx.Rollback()
			return resStatus, err
		}

		if affectedRows < 1 {
			resStatus = i18n.StatusOtherEdit
			zap.L().Info("DeleteEPCs" + eid.Name + "other user editing")
			tx.Rollback()
			return resStatus, nil
		}
		//从localCache删除
		eid.DelFromLocalCache()
	}

	return
}

// EPA GetInfoByID() 根据id获取执行项目档案详情
func (eid *EPA) GetInfoByID() (resStatus i18n.ResKey, err error) {
	//从localcache获取
	number, b, _ := cache.Get(pub.EID, eid.ID)
	if number > 0 {
		json.Unmarshal(b, &eid)
		resStatus = i18n.StatusOK

		return
	}
	//从数据库获取
	sqlStr := `select code,name,epcid,description,status,
	resulttypeid,udcid,defaultvalue,defaultvaluedisp,ischeckerror,
	errorvalue,errorvaluedisp,isrequirefile,isonsitephoto,risklevelid,
	createtime,creatorid,modifytime,modifierid,ts,dr
	from exectiveitem where id = $1`
	err = db.QueryRow(sqlStr, eid.ID).Scan(&eid.Code, &eid.Name, &eid.EPC.ID, &eid.Description, &eid.Status,
		&eid.ResultType.ID, &eid.UDCID.ID, &eid.DefaultValue, &eid.DefaultValueDisp, &eid.IsCheckError,
		&eid.ErrorValue, &eid.ErrorValueDisp, &eid.IsRequireFile, &eid.IsOnsitePhoto, &eid.RiskLevel.ID,
		&eid.CreateDate, &eid.Creator.ID, &eid.ModifyDate, &eid.Modifier.ID, &eid.Ts, &eid.Dr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("eid.GetInfoByID db.QueryRow failed", zap.Error(err))
		return
	}
	//填充执行项目档案类别信息
	if eid.EPC.ID > 0 {
		resStatus, err = eid.EPC.GetEPCInfoByID()
		if err != nil || resStatus != i18n.StatusOK {
			return
		}
	}

	//填充执行结果类型数据
	eid.ResultType.GetDataTypeInfoByID()

	//填充用户自定义档案类别
	if eid.UDCID.ID > 0 {
		resStatus, err = eid.UDCID.GetUDCInfoByID()
		if err != nil || resStatus != i18n.StatusOK {
			return
		}
	}
	//填充风险等级
	if eid.RiskLevel.ID > 0 {
		resStatus, err = eid.RiskLevel.GetRLInfoByID()
		if err != nil || resStatus != i18n.StatusOK {
			return
		}
	}
	//填充创建人信息
	if eid.Creator.ID > 0 {
		resStatus, err = eid.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//填充更新人信息
	if eid.Modifier.ID > 0 {
		resStatus, err = eid.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//写入localcache
	eidB, _ := json.Marshal(eid)
	cache.Set(pub.EID, eid.ID, eidB)

	return
}

// EPA.DelFromLocalCache 从localcache删除
func (eid *EPA) DelFromLocalCache() {
	number, _, _ := cache.Get(pub.EID, eid.ID) //判断是否存在于本地缓存中
	if number > 0 {                            //如果存在于本地缓存中则直接删除
		cache.Del(pub.EID, eid.ID)
	}
}

// EPA CheckUsed() 检查是否被其他项目引用
func (eid *EPA) CheckUsed() (resStatus i18n.ResKey, err error) {
	//检查项目
	checkItems := []ScDocCheckUsed{
		{
			Description:    "被执行模板引用",
			SqlStr:         `select count(id) as usednumber from exectivetemplate_b where dr=0 and eid_id=$1`,
			UsedReturnCode: i18n.StatusEITUsed,
		},
		{
			Description:    "被执行单引用",
			SqlStr:         `select count(id) as usednumber from executedoc_b where dr=0 and eid_id=$1`,
			UsedReturnCode: i18n.StatusEDUsed,
		},
		{
			Description:    "被问题处理单单引用",
			SqlStr:         `select count(id) as usednumber from disposedoc where dr=0 and eid_id=$1`,
			UsedReturnCode: i18n.StatusDDUsed,
		},
	}
	//检查项目
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, eid.ID).Scan(&usedNum)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("EPA.CheckIsUsed "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return
}
*/
