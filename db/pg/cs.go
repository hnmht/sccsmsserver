package pg

func initCS() (isFinish bool, err error) {
	return
}

/*
// SceneItem 现场档案模型
type SceneItem struct {
	ID          int32         `db:"id" json:"id"`
	Code        string        `db:"code" json:"code"`
	Name        string        `db:"name" json:"name"`
	Description string        `db:"description" json:"description"`
	ItemClass   SimpCSC       `db:"class_id" json:"itemclass"`       //所属分类
	Department  SimpDept      `db:"subdept_id" json:"subdept"`       //所属部门
	RespDept    SimpDept      `db:"respdept_id" json:"respdept"`     //负责部门
	RespPerson  Person        `db:"respperson_id" json:"respperson"` //负责人
	Status      int16         `db:"status" json:"status"`            //0 在用 1 停用
	FinishFlag  int16         `db:"finishflag" json:"finishflag"`    //完工标志
	FinishDate  string        `db:"finishdate" json:"finishdate"`    //完工日期
	Longitude   float64       `db:"longitude" json:"longitude"`      //经度
	Latitude    float64       `db:"latitude" json:"latitude"`        //纬度
	Udf1        UserDefineDoc `db:"udf1" json:"udf1"`                //用户自定义档案1
	Udf2        UserDefineDoc `db:"udf2" json:"udf2"`
	Udf3        UserDefineDoc `db:"udf3" json:"udf3"`
	Udf4        UserDefineDoc `db:"udf4" json:"udf4"`
	Udf5        UserDefineDoc `db:"udf5" json:"udf5"`
	Udf6        UserDefineDoc `db:"udf6" json:"udf6"`
	Udf7        UserDefineDoc `db:"udf7" json:"udf7"`
	Udf8        UserDefineDoc `db:"udf8" json:"udf8"`
	Udf9        UserDefineDoc `db:"udf9" json:"udf9"`
	Udf10       UserDefineDoc `db:"udf10" json:"udf10"`
	CreateDate  time.Time     `db:"create_time" json:"createdate"`
	CreateUser  Person        `db:"createuserid" json:"createuser"`
	ModifyDate  time.Time     `db:"modify_time" json:"modifydate"`
	ModifyUser  Person        `db:"modifyuserid" json:"modifyuser"`
	Ts          time.Time     `db:"ts" json:"ts"`
	Dr          int16         `db:"dr" json:"dr"` //删除标志
}

// SceneItemOption 现场档案自定义项模型
type SceneItemOption struct {
	ID           int32           `db:"id" json:"id"`
	Code         string          `db:"code" json:"code"`
	Name         string          `db:"name" json:"name"`
	DisplayName  string          `db:"displayname" json:"displayname"`
	UDC          UserDefineClass `db:"udc_id" json:"udc"`
	DefaultValue UserDefineDoc   `db:"defaultvalue_id" json:"defaultvalue"`
	Enable       int16           `db:"enable" json:"enable"`
	IsModify     int16           `json:"ismodify"`
	CreateDate   time.Time       `db:"create_time" json:"createdate"`
	CreateUser   Person          `db:"createuserid" json:"createuser"`
	ModifyDate   time.Time       `db:"modify_time" json:"modifydate"`
	ModifyUser   Person          `db:"modifyuserid" json:"modifyuser"`
	Ts           time.Time       `db:"ts" json:"ts"`
	Dr           int16           `db:"dr" json:"dr"` //删除标志
}

// SceneItemCache 现场档案缓存
type SceneItemCache struct {
	QueryTs      time.Time   `json:"queryts"`     //查询ts
	ResultNumber int32       `json:"resultnum"`   //获取数据数量
	DelItems     []SceneItem `json:"delitems"`    //已经删除的档案
	UpdateItems  []SceneItem `json:"updateitems"` //已经更新的档案
	NewItems     []SceneItem `json:"newitems"`    //新增加的档案
	ResultTs     time.Time   `json:"resultts"`    //返回结果的ts
}

// SceneItemOptionCache 现场档案选项缓存
type SceneItemOptionCache struct {
	QueryTs      time.Time         `json:"queryts"`     //查询ts
	ResultNumber int32             `json:"resultnum"`   //获取数据数量
	DelItems     []SceneItemOption `json:"delitems"`    //已经删除的档案
	UpdateItems  []SceneItemOption `json:"updateitems"` //已经更新的档案
	NewItems     []SceneItemOption `json:"newitems"`    //新增加的档案
	ResultTs     time.Time         `json:"resultts"`    //返回结果的ts
}

// GetSIs 获取现场档案列表
func GetSIs() (sis []SceneItem, resStatus i18n.ResKey, err error) {
	sis = make([]SceneItem, 0)
	sqlStr := `select id,code,name,description,class_id,subdept_id,
	respdept_id,respperson_id,status,finishflag,finishdate,
	longitude,latitude,udf1,udf2,udf3,
	udf4,udf5,udf6,udf7,udf8,
	udf9,udf10,create_time,createuserid,modify_time,
	modifyuserid,ts,dr
	from sceneitem where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetSis db.Query failed", zap.Error(err))
		return
	}
	defer rows.Close()

	var rowsNum int32
	for rows.Next() {
		rowsNum++
		var si SceneItem
		err = rows.Scan(&si.ID, &si.Code, &si.Name, &si.Description, &si.ItemClass.ID, &si.Department.ID,
			&si.RespDept.ID, &si.RespPerson.ID, &si.Status, &si.FinishFlag, &si.FinishDate,
			&si.Longitude, &si.Latitude, &si.Udf1.ID, &si.Udf2.ID, &si.Udf3.ID,
			&si.Udf4.ID, &si.Udf5.ID, &si.Udf6.ID, &si.Udf7.ID, &si.Udf8.ID,
			&si.Udf9.ID, &si.Udf10.ID, &si.CreateDate, &si.CreateUser.ID, &si.ModifyDate,
			&si.ModifyUser.ID, &si.Ts, &si.Dr)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetSIs rows.Next failed", zap.Error(err))
			return
		}

		//获取附加信息
		resStatus, err = si.GetAttachInfo()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		sis = append(sis, si)
	}
	//如果获取的数据行数等于0,则表示没有数据
	if rowsNum == 0 {
		resStatus = i18n.StatusResNoData
		return
	}

	resStatus = i18n.StatusOK
	return
}

// GetSICache 获取现场档案缓存
func (sic *SceneItemCache) GetSICache() (resStatus i18n.ResKey, err error) {
	sic.DelItems = make([]SceneItem, 0)
	sic.NewItems = make([]SceneItem, 0)
	sic.UpdateItems = make([]SceneItem, 0)
	//查询现场档案最新缓存
	sqlStr := `select ts from sceneitem where ts>$1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, sic.QueryTs).Scan(&sic.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			sic.ResultNumber = 0
			sic.ResultTs = sic.QueryTs
			resStatus = i18n.StatusOK
			return
		}
		zap.L().Error("SceneItemCache.GetSICache db.QueryRow failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	//查询所有大于QueryTs的数据
	sqlStr = `select id,code,name,description,class_id,subdept_id,
	respdept_id,respperson_id,status,finishflag,finishdate,
	longitude,latitude,udf1,udf2,udf3,
	udf4,udf5,udf6,udf7,udf8,
	udf9,udf10,create_time,createuserid,modify_time,
	modifyuserid,ts,dr
	from sceneitem where ts>$1 order by ts desc`
	rows, err := db.Query(sqlStr, sic.QueryTs)
	if err != nil {
		zap.L().Error("SceneItemCache.GetSICache db.Query failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()

	for rows.Next() {
		var si SceneItem
		err = rows.Scan(&si.ID, &si.Code, &si.Name, &si.Description, &si.ItemClass.ID, &si.Department.ID,
			&si.RespDept.ID, &si.RespPerson.ID, &si.Status, &si.FinishFlag, &si.FinishDate,
			&si.Longitude, &si.Latitude, &si.Udf1.ID, &si.Udf2.ID, &si.Udf3.ID,
			&si.Udf4.ID, &si.Udf5.ID, &si.Udf6.ID, &si.Udf7.ID, &si.Udf8.ID,
			&si.Udf9.ID, &si.Udf10.ID, &si.CreateDate, &si.CreateUser.ID, &si.ModifyDate,
			&si.ModifyUser.ID, &si.Ts, &si.Dr)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("SceneItemCache.GetSICache rows.Next failed", zap.Error(err))
			return
		}
		//获取附加信息
		resStatus, err = si.GetAttachInfo()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}

		if si.Dr == 0 {
			if si.CreateDate.Before(sic.QueryTs) || si.CreateDate.Equal(sic.QueryTs) {
				sic.ResultNumber++
				sic.UpdateItems = append(sic.UpdateItems, si)
			} else {
				sic.ResultNumber++
				sic.NewItems = append(sic.NewItems, si)
			}
		} else {
			if si.CreateDate.Before(sic.QueryTs) || si.CreateDate.Equal(sic.QueryTs) {
				sic.ResultNumber++
				sic.DelItems = append(sic.DelItems, si)
			}
		}
	}
	return i18n.StatusOK, nil
}

// SceneItem.GetAttachInfo 获取现场档案附加信息
func (si *SceneItem) GetAttachInfo() (resStatus i18n.ResKey, err error) {
	//获取现场档案类别
	if si.ItemClass.ID > 0 {
		resStatus, err = si.ItemClass.GetSSICInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//获取Department详情
	if si.Department.ID > 0 {
		resStatus, err = si.Department.GetSimpDeptInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//获取RespDept详情
	if si.RespDept.ID > 0 {
		resStatus, err = si.RespDept.GetSimpDeptInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//获取RespPerson详情
	if si.RespPerson.ID > 0 {
		resStatus, err = si.RespPerson.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//获取Udf1详情
	if si.Udf1.ID > 0 {
		resStatus, err = si.Udf1.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//获取Udf2详情
	if si.Udf2.ID > 0 {
		resStatus, err = si.Udf2.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//获取Udf3详情
	if si.Udf3.ID > 0 {
		resStatus, err = si.Udf3.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//获取Udf4详情
	if si.Udf4.ID > 0 {
		resStatus, err = si.Udf4.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//获取Udf5详情
	if si.Udf5.ID > 0 {
		resStatus, err = si.Udf5.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//获取Udf6详情
	if si.Udf6.ID > 0 {
		resStatus, err = si.Udf6.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//获取Udf7详情
	if si.Udf7.ID > 0 {
		resStatus, err = si.Udf7.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//获取Udf8详情
	if si.Udf8.ID > 0 {
		resStatus, err = si.Udf8.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//获取Udf9详情
	if si.Udf9.ID > 0 {
		resStatus, err = si.Udf9.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//获取Udf10详情
	if si.Udf10.ID > 0 {
		resStatus, err = si.Udf10.GetInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//获取CreateUser详情
	if si.CreateUser.ID > 0 {
		resStatus, err = si.CreateUser.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	//获取ModifyUser详情
	if si.ModifyUser.ID > 0 {
		resStatus, err = si.ModifyUser.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}

	return i18n.StatusOK, nil
}

// SceneItem.GetInfoByID 根据ID获取现场档案详情
func (si *SceneItem) GetInfoByID() (resStatus i18n.ResKey, err error) {
	//从localcache获取
	number, b, _ := cache.Get(i18n.SI, si.ID)
	if number > 0 {
		json.Unmarshal(b, &si)
		resStatus = i18n.StatusOK
		return
	}
	//从数据库获取
	sqlStr := `select code,name,description,class_id,subdept_id,
	respdept_id,respperson_id,status,finishflag,finishdate,
	longitude,latitude,udf1,udf2,udf3,
	udf4,udf5,udf6,udf7,udf8,
	udf9,udf10,create_time,createuserid,modify_time,
	modifyuserid,ts,dr
	from sceneitem where id=$1`
	err = db.QueryRow(sqlStr, si.ID).Scan(
		&si.Code, &si.Name, &si.Description, &si.ItemClass.ID, &si.Department.ID,
		&si.RespDept.ID, &si.RespPerson.ID, &si.Status, &si.FinishFlag, &si.FinishDate,
		&si.Longitude, &si.Latitude, &si.Udf1.ID, &si.Udf2.ID, &si.Udf3.ID,
		&si.Udf4.ID, &si.Udf5.ID, &si.Udf6.ID, &si.Udf7.ID, &si.Udf8.ID,
		&si.Udf9.ID, &si.Udf10.ID, &si.CreateDate, &si.CreateUser.ID, &si.ModifyDate,
		&si.ModifyUser.ID, &si.Ts, &si.Dr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("SceneItem.GetInfoByID db.QueryRow failed", zap.Error(err))
		return
	}

	//获取附加信息
	resStatus, err = si.GetAttachInfo()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	//写入localcache
	siB, _ := json.Marshal(si)
	cache.Set(i18n.SI, si.ID, siB)

	return i18n.StatusOK, nil
}

// SceneItem.Add 增加现场档案
func (si *SceneItem) Add() (resStatus i18n.ResKey, err error) {
	//检查现场档案编码是否存在
	resStatus, err = si.CheckCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	//向数据库写入现场档案
	sqlStr := `insert into sceneitem(code,name,description,class_id,subdept_id,
	respdept_id,respperson_id,status,finishflag,finishdate,
	longitude,latitude,udf1,udf2,udf3,
	udf4,udf5,udf6,udf7,udf8,
	udf9,udf10,createuserid,modifyuserid)
	values($1,$2,$3,$4,$5,
	$6,$7,$8,$9,$10,
	$11,$12,$13,$14,$15,
	$16,$17,$18,$19,$20,
	$21,$22,$23,$24)
	returning id`
	_, err = db.Exec(sqlStr, si.Code, si.Name, si.Description, si.ItemClass.ID, si.Department.ID,
		si.RespDept.ID, si.RespPerson.ID, si.Status, si.FinishFlag, si.FinishDate,
		si.Longitude, si.Latitude, si.Udf1.ID, si.Udf2.ID, si.Udf3.ID,
		si.Udf4.ID, si.Udf5.ID, si.Udf6.ID, si.Udf7.ID, si.Udf8.ID,
		si.Udf9.ID, si.Udf10.ID, si.CreateUser.ID, si.ModifyUser.ID)

	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("SceneItem.Add db.Exec failed", zap.Error(err))
		return
	}
	return i18n.StatusOK, nil
}

// SceneItem.Edit 修改现场档案
func (si *SceneItem) Edit() (resStatus i18n.ResKey, err error) {
	//检查现场档案编码是否存在
	resStatus, err = si.CheckCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	//向数据库写入现场档案修改
	sqlStr := `update sceneitem set code=$1,name=$2,description=$3,class_id=$4,subdept_id=$5,
	respdept_id=$6,respperson_id=$7,status=$8,finishflag=$9,finishdate=$10,
	longitude=$11,latitude=$12,udf1=$13,udf2=$14,udf3=$15,
	udf4=$16,udf5=$17,udf6=$18,udf7=$19,udf8=$20,
	udf9=$21,udf10=$22,modifyuserid=$23,
	modify_time=current_timestamp,ts=current_timestamp
	where id=$24 and ts=$25 and dr=0`
	res, err := db.Exec(sqlStr, si.Code, si.Name, si.Description, si.ItemClass.ID, si.Department.ID,
		si.RespDept.ID, si.RespPerson.ID, si.Status, si.FinishFlag, si.FinishDate,
		si.Longitude, si.Latitude, si.Udf1.ID, si.Udf2.ID, si.Udf3.ID,
		si.Udf4.ID, si.Udf5.ID, si.Udf6.ID, si.Udf7.ID, si.Udf8.ID,
		si.Udf9.ID, si.Udf10.ID, si.ModifyUser.ID,
		si.ID, si.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("SceneItem.Edit db.Exec failed", zap.Error(err))
		return
	}
	//检查更新影响的行数
	effected, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("SceneItem.Edit res.RowsAffected failed", zap.Error(err))
		return
	}
	//如果更新影响的行数为0,则说明其他人正在操作该档案
	if effected < 1 {
		zap.L().Info("SceneItem.Edit failed,Other user are Editing")
		resStatus = i18n.StatusOtherEdit
		return
	}
	//从localCache删除
	si.DelFromLocalCache()

	return i18n.StatusOK, nil
}

// SceneItem.CheckCodeExist 检查现场档案编码是否存在
func (si *SceneItem) CheckCodeExist() (resStatus i18n.ResKey, err error) {
	var count int32
	sqlStr := `select count(id) from sceneitem where dr=0 and class_id=$1 and code=$2 and id <>$3`
	err = db.QueryRow(sqlStr, si.ItemClass.ID, si.Code, si.ID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("SceneItem.CheckCodeExist db.QueryRow failed", zap.Error(err))
		return
	}

	if count > 0 {
		resStatus = i18n.StatusSICodeExist
		return
	}

	return i18n.StatusOK, nil
}

// GetSceneItemOptions 获取现场档案自定义项模型
func GetSIOs() (sciOptions []SceneItemOption, resStatus i18n.ResKey, err error) {
	sqlStr := `select id,code,name,displayname,udc_id,
	defaultvalue_id,enable,create_time,createuserid,modify_time,
	modifyuserid,ts,dr
	from sceneitemoption order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetSceneItemOptions db.Query failed", zap.Error(err))
		return
	}
	defer rows.Close()
	var rowNumber int32
	for rows.Next() {
		rowNumber++
		var option SceneItemOption
		err = rows.Scan(&option.ID, &option.Code, &option.Name, &option.DisplayName, &option.UDC.ClassID,
			&option.DefaultValue.ID, &option.Enable, &option.CreateDate, &option.CreateUser.ID, &option.ModifyDate,
			&option.ModifyUser.ID, &option.Ts, &option.Dr)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetSceneItemOptions rows.Next failed", zap.Error(err))
			return
		}
		//获取自定义档案类别详情
		if option.UDC.ClassID > 0 {
			resStatus, err = option.UDC.GetUDCInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//获取默认值详情
		if option.DefaultValue.ID > 0 {
			resStatus, err = option.DefaultValue.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//获取更新人信息
		if option.ModifyUser.ID > 0 {
			resStatus, err = option.ModifyUser.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		//获取是否能够修改信息isModify
		resStatus, err = option.CheckIsModify()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		sciOptions = append(sciOptions, option)
	}
	//如果获取的数据行数等于0,则表示没有数据
	if rowNumber == 0 {
		resStatus = i18n.StatusResNoData
		return
	}
	resStatus = i18n.StatusOK

	return
}

// SceneItemOption.GetSIOCache 获取现场档案选项缓存
func (sioc *SceneItemOptionCache) GetSIOCache() (resStatus i18n.ResKey, err error) {
	//查询现场档案选项最新缓存
	sqlStr := `select ts from sceneitemoption where ts > $1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, sioc.QueryTs).Scan(&sioc.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			sioc.ResultNumber = 0
			sioc.ResultTs = sioc.QueryTs
			resStatus = i18n.StatusOK
			return
		}
		zap.L().Error("SceneItemOption.GetSIOCache db.QueryRow failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	//查询所有大于QueryTs的数据
	sqlStr = `select id,code,name,displayname,udc_id,
	defaultvalue_id,enable,create_time,createuserid,modify_time,
	modifyuserid,ts,dr
	from sceneitemoption where ts>$1 order by ts desc`
	rows, err := db.Query(sqlStr, sioc.QueryTs)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("SceneItemOption.GetSIOCache  db.Query failed", zap.Error(err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var option SceneItemOption
		err = rows.Scan(&option.ID, &option.Code, &option.Name, &option.DisplayName, &option.UDC.ClassID,
			&option.DefaultValue.ID, &option.Enable, &option.CreateDate, &option.CreateUser.ID, &option.ModifyDate,
			&option.ModifyUser.ID, &option.Ts, &option.Dr)
		if err != nil {
			zap.L().Error("SceneItemOption.GetSIOCache rows.Next failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		//获取自定义档案类别详情
		if option.UDC.ClassID > 0 {
			resStatus, err = option.UDC.GetUDCInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//获取默认值详情
		if option.DefaultValue.ID > 0 {
			resStatus, err = option.DefaultValue.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//获取更新人信息
		if option.ModifyUser.ID > 0 {
			resStatus, err = option.ModifyUser.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//获取是否能够修改信息isModify
		resStatus, err = option.CheckIsModify()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}

		if option.Dr == 0 {
			if option.CreateDate.Before(sioc.QueryTs) || option.CreateDate.Equal(sioc.QueryTs) {
				sioc.ResultNumber++
				sioc.UpdateItems = append(sioc.UpdateItems, option)
			} else {
				sioc.ResultNumber++
				sioc.NewItems = append(sioc.NewItems, option)
			}
		} else {
			if option.CreateDate.Before(sioc.QueryTs) || option.CreateDate.Equal(sioc.QueryTs) {
				sioc.ResultNumber++
				sioc.DelItems = append(sioc.DelItems, option)
			}
		}
	}

	return i18n.StatusOK, nil
}

// SceneItemOption.CheckIsModify 获取是否能够修改信息
func (sio *SceneItemOption) CheckIsModify() (resStatus i18n.ResKey, err error) {
	var usedNumber int32
	sqlStr := `select count(id) as usednum from sceneitem where dr=0 and ` + sio.Code + `>0`
	err = db.QueryRow(sqlStr).Scan(&usedNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("SceneItemOption.CheckIsModify db.QueryRow failed", zap.Error(err))
		return
	}

	if usedNumber > 0 { //如果已经使用过,则设置为否
		sio.IsModify = 1
	} else {
		sio.IsModify = 0
	}

	return i18n.StatusOK, nil
}

// SceneItemOption.Edit 修改现场档案自定义项
func (sio *SceneItemOption) Edit() (resStatus i18n.ResKey, err error) {
	sqlStr := `update sceneitemoption set displayname=$1,udc_id=$2,defaultvalue_id=$3 ,enable=$4,
	modify_time=current_timestamp, modifyuserid = $5,ts=current_timestamp where id=$6 and ts=$7`
	//写入
	res, err := db.Exec(sqlStr, sio.DisplayName, sio.UDC.ClassID, sio.DefaultValue.ID, sio.Enable,
		sio.ModifyUser.ID, sio.ID, sio.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("SceneItemOption.Edit db.Exec failed", zap.Error(err))
		return
	}

	//检查更新的行数
	updateNumber, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("SceneItemOption.Edit  res.RowsAffected failed", zap.Error(err))
		return
	}
	if updateNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		return
	}

	return i18n.StatusOK, nil
}

// SceneItem.Delete 删除现场档案
func (si *SceneItem) Delete() (resStatus i18n.ResKey, err error) {
	//检查现场档案是否被引用
	resStatus, err = si.CheckUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}

	//删除现场档案
	sqlStr := `update sceneitem set dr=1,modifyuserid=$1,modify_time=current_timestamp,ts=current_timestamp
	where id=$2 and dr=0 and ts=$3`
	res, err := db.Exec(sqlStr, si.ModifyUser.ID, si.ID, si.Ts)
	if err != nil {
		zap.L().Error("SceneItem.Delete db.Exec failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	//检查删除操作影响的行数
	effected, err := res.RowsAffected()
	if err != nil {
		zap.L().Error("SceneItem.Delete  res.RowsAffected failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	if effected < 1 {
		resStatus = i18n.StatusOtherEdit
		return
	}
	//从localCache删除
	si.DelFromLocalCache()

	return i18n.StatusOK, nil
}

// DeleteSIs 批量删除现场档案
func DeleteSIs(sis *[]SceneItem, modifyUserId int32) (resStatus i18n.ResKey, err error) {
	//开始执行事务
	tx, err := db.Begin()
	if err != nil {
		zap.L().Error("DeleteSIs db.Begin failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer tx.Commit()

	delSqlStr := `update sceneitem set dr=1,modifyuserid=$1,modify_time=current_timestamp,ts=current_timestamp
	where id=$2 and dr=0 and ts=$3`
	//删除操作预处理
	stmt, err := tx.Prepare(delSqlStr)
	if err != nil {
		zap.L().Error("DeleteSIs tx.Prepare failed", zap.Error(err))
		resStatus = i18n.StatusOK
		tx.Rollback()
		return
	}
	defer stmt.Close()

	for _, si := range *sis {
		//检查现场档案是否被引用
		resStatus, err = si.CheckUsed()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
		//执行删除操作
		res, err1 := stmt.Exec(modifyUserId, si.ID, si.Ts)
		if err != nil {
			zap.L().Error("DeleteSIs stmt.Exec failed", zap.Error(err1))
			resStatus = i18n.StatusInternalError
			tx.Rollback()
			return resStatus, err1
		}

		//检查删除操作影响的行数
		affected, err2 := res.RowsAffected()
		if err2 != nil {
			zap.L().Error("DeleteSIs check res.RowsAffected failed", zap.Error(err))
			_ = tx.Rollback()
			return i18n.StatusInternalError, err2
		}
		if affected < 1 {
			zap.L().Info("DeleteSIs other edit")
			_ = tx.Rollback()
			return i18n.StatusOtherEdit, nil
		}
		//从localCache删除
		si.DelFromLocalCache()
	}
	return i18n.StatusOK, nil
}

// SceneItem.DelFromLocalCache 从localCache删除
func (si *SceneItem) DelFromLocalCache() {
	number, _, _ := cache.Get(i18n.SI, si.ID) //判断是否存在于本地缓存中
	if number > 0 {                           //如果存在于本地缓存中则直接删除
		cache.Del(i18n.SI, si.ID)
	}
}

// SceneItem.CheckUsed 检查现场档案是否被引用
func (si *SceneItem) CheckUsed() (resStatus i18n.ResKey, err error) {
	checkItems := []ScDocCheckUsed{
		{
			Description:    "被指令单引用",
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
		},
	}
	//检查项目
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, si.ID).Scan(&usedNum)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("SceneItem.CheckUsed  "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return i18n.StatusOK, nil
}
*/
