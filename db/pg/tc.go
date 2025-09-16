package pg

/* // TrainCourse 培训课程结构体
type TrainCourse struct {
	ID          int32         `db:"id" json:"id"`                   //ID
	Code        string        `db:"code" json:"code"`               //编码
	Name        string        `db:"name" json:"name"`               //名称
	ClassHour   float64       `db:"classhour" json:"classhour"`     //课时
	IsExamine   int16         `db:"isexamine" json:"isexamine"`     //是否考核
	Description string        `db:"description" json:"description"` //说明
	Files       []VoucherFile `json:"files"`                        //文件附件
	CreateDate  time.Time     `db:"create_time" json:"createdate"`  //创建日期
	CreateUser  Person        `db:"createuserid" json:"createuser"`
	ModifyDate  time.Time     `db:"modify_time" json:"modifydate"`
	ModifyUser  Person        `db:"modifyuserid" json:"modifyuser"`
	Ts          time.Time     `db:"ts" json:"ts"` //时间戳
	Dr          int16         `db:"dr" json:"dr"` //删除标志
}

// TCCache 培训课程缓存
type TCCache struct {
	QueryTs      time.Time     `json:"queryts"`
	ResultNumber int32         `json:"resultnum"`
	DelItems     []TrainCourse `json:"delitems"`
	UpdateItems  []TrainCourse `json:"updateitems"`
	NewItems     []TrainCourse `json:"newitems"`
	ResultTs     time.Time     `json:"resultts"`
}

// TrainCourse.Add 增加培训课程
func (tc *TrainCourse) Add() (resStatus pub.ResStatus, err error) {
	//检查名称是否重复
	resStatus, err = tc.CheckNameExist()
	if resStatus != pub.StatusOK || err != nil {
		return
	}
	//创建事务
	tx, err := db.Begin()
	if err != nil {
		resStatus = pub.StatusInternalError
		zap.L().Error("TrainCourse.Add db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	headSql := `insert into traincourse(name,classhour,isexamine,description,createuserid)
	 	values($1,$2,$3,$4,$5)
		returning id`
	err = tx.QueryRow(headSql, tc.Name, tc.ClassHour, tc.IsExamine, tc.Description, tc.CreateUser.UserID).Scan(&tc.ID)
	if err != nil {
		resStatus = pub.StatusInternalError
		zap.L().Error("TrainCourse.Add tx.QueryRow failed", zap.Error(err))
		tx.Rollback()
		return
	}

	//插入文件表
	fileSql := `insert into traincourse_file(billhid,fileid,createuserid)
	values($1,$2,$3)
	returning id`
	fileStmt, err := tx.Prepare(fileSql)
	if err != nil {
		resStatus = pub.StatusInternalError
		zap.L().Error("TrainCourse.Add  tx.Prepare(fileSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer fileStmt.Close()

	for _, file := range tc.Files {
		err = fileStmt.QueryRow(tc.ID, file.File.FileId, tc.CreateUser.UserID).Scan(&file.ID)
		if err != nil {
			resStatus = pub.StatusInternalError
			zap.L().Error("TrainCourse.Add  fileStmt.QueryRow failed", zap.Error(err))
			tx.Rollback()
			return
		}
	}
	return pub.StatusOK, nil
}

// TrainCourse.Edit 修改培训课程
func (tc *TrainCourse) Edit() (resStatus pub.ResStatus, err error) {
	//检查创建人和修改人是否同为一人
	if tc.CreateUser.UserID != tc.ModifyUser.UserID {
		resStatus = pub.StatusOtherEdit
		return
	}
	//检查名称是否重复
	resStatus, err = tc.CheckNameExist()
	if resStatus != pub.StatusOK || err != nil {
		return
	}
	//创建事务
	tx, err := db.Begin()
	if err != nil {
		resStatus = pub.StatusInternalError
		zap.L().Error("TrainCourse.Edit db.Begin() failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	//修改文档内容
	editDocSql := `update traincourse set code=$1,name=$2,classhour=$3,isexamine=$4,description=$5,
		modify_time=current_timestamp,modifyuserid=$6,ts=current_timestamp
		where id=$7 and dr=0 and ts=$8`
	editDocRes, err := tx.Exec(editDocSql, &tc.Code, &tc.Name, &tc.ClassHour, &tc.IsExamine, &tc.Description,
		&tc.ModifyUser.UserID,
		&tc.ID, &tc.Ts)
	if err != nil {
		resStatus = pub.StatusInternalError
		zap.L().Error("TrainCourse.Edit tx.Exec(editDocSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	//检查文档修改的行数
	updateNumber, err := editDocRes.RowsAffected()
	if err != nil {
		resStatus = pub.StatusInternalError
		zap.L().Error("TrainCourse.Edit editDocRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if updateNumber < 1 {
		resStatus = pub.StatusOtherEdit
		tx.Rollback()
		return
	}

	//修改文件内容
	updateFileSql := `update TrainCourse_file set modifyuserid=$1,modify_time=current_timestamp,dr=$2,ts=current_timestamp
		where id=$3 and billhid=$4 and dr=0 and ts=$5`
	addFileSql := `insert into TrainCourse_file(billhid,fileid,createuserid) values($1,$2,$3) returning id`
	//更新文件准备
	updateFileStmt, err := tx.Prepare(updateFileSql)
	if err != nil {
		zap.L().Error("TrainCourse.Edit tx.Prepare(updateFileSql) failed", zap.Error(err))
		resStatus = pub.StatusInternalError
		tx.Rollback()
		return
	}
	defer updateFileStmt.Close()
	//增加文件准备
	addFileStmt, err := tx.Prepare(addFileSql)
	if err != nil {
		zap.L().Error("TrainCourse.Edit tx.Prepare(addFileStmt) failed", zap.Error(err))
		resStatus = pub.StatusInternalError
		tx.Rollback()
		return
	}
	defer addFileStmt.Close()

	for _, file := range tc.Files {
		if file.ID != 0 { //修改原有文件
			updateFileRes, updateFileErr := updateFileStmt.Exec(tc.ModifyUser.UserID, file.Dr, file.ID, tc.ID, file.Ts)
			if updateFileErr != nil {
				resStatus = pub.StatusInternalError
				zap.L().Error("TrainCourse.Edit  updateFileRes.Exec() failed", zap.Error(updateFileErr))
				tx.Rollback()
				return resStatus, updateFileErr
			}
			updateFileNumber, updateFileEffectErr := updateFileRes.RowsAffected()
			if updateFileEffectErr != nil {
				resStatus = pub.StatusInternalError
				zap.L().Error("TrainCourse.EditupdateFileRes.RowsAffected failed", zap.Error(updateFileEffectErr))
				tx.Rollback()
				return resStatus, updateFileEffectErr
			}
			if updateFileNumber < 1 {
				resStatus = pub.StatusOtherEdit
				tx.Rollback()
				return
			}
		} else { //新增文件
			addFileErr := addFileStmt.QueryRow(tc.ID, file.File.FileId, tc.ModifyUser.UserID).Scan(&file.ID)
			if addFileErr != nil {
				resStatus = pub.StatusInternalError
				zap.L().Error("Document.Edit addFileStmt.QueryRow failed", zap.Error(addFileErr))
				tx.Rollback()
				return resStatus, addFileErr
			}
		}
	}
	//从缓存删除
	tc.DelFromLocalCache()

	return pub.StatusOK, nil
}

// TrainCourse.CheckNameExist 检查名称是否存在
func (tc *TrainCourse) CheckNameExist() (resStatus pub.ResStatus, err error) {
	var count int32
	sqlStr := `select count(id) from traincourse where dr=0 and name=$1 and id <> $2`
	err = db.QueryRow(sqlStr, tc.Name, tc.ID).Scan(&count)
	if err != nil {
		resStatus = pub.StatusInternalError
		zap.L().Error("TrainCourse.CheckNameExist query failed", zap.Error(err))
		return
	}
	if count > 0 {
		resStatus = pub.StatusTCNameExist
		return
	}

	resStatus = pub.StatusOK
	return
}

// TrainCourse.DelFromLocalCache 从localCache删除
func (tc *TrainCourse) DelFromLocalCache() {
	number, _, _ := cache.Get(pub.TC, tc.ID)
	if number > 0 {
		cache.Del(pub.TC, tc.ID)
	}
}

// Document.Delete 删除培训课程
func (tc *TrainCourse) Delete(modifyUserId int32) (resStatus pub.ResStatus, err error) {
	//检查是否被引用
	resStatus, err = tc.CheckIsUsed()
	if resStatus != pub.StatusOK || err != nil {
		return
	}
	//检查删除人和创建人是否同为一人
	if tc.CreateUser.UserID != modifyUserId {
		resStatus = pub.StatusOtherEdit
		return
	}
	//创建事务
	tx, err := db.Begin()
	if err != nil {
		resStatus = pub.StatusInternalError
		zap.L().Error("TrainCourse.Delete db.Begin() failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	//删除
	delDocSql := `update traincourse set dr=1,modify_time=current_timestamp,modifyuserid=$1,ts=current_timestamp
		where id=$2 and dr=0 and ts=$3`
	delDocRes, err := tx.Exec(delDocSql, modifyUserId, tc.ID, tc.Ts)
	if err != nil {
		resStatus = pub.StatusInternalError
		zap.L().Error("TrainCourse.Delete tx.Exec(delDocSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	//检查删除效果
	delDocNumber, err := delDocRes.RowsAffected()
	if err != nil {
		resStatus = pub.StatusInternalError
		zap.L().Error("TrainCourse.Delete delDocRes.RowsAffected() failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if delDocNumber < 1 {
		resStatus = pub.StatusOtherEdit
		tx.Rollback()
		return
	}
	//删除附件
	delFileSql := `update traincourse_file set dr=1,modify_time=current_timestamp,modifyuserid=$1,ts=current_timestamp
		where id=$2 and dr=0 and billhid=$3 and ts=$4`
	//预处理
	delFileStmt, err := tx.Prepare(delFileSql)
	if err != nil {
		resStatus = pub.StatusInternalError
		zap.L().Error("TrainCourse.Delete tx.Prepare(delFileSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delFileStmt.Close()

	for _, file := range tc.Files {
		delFileRes, errDelFile := delFileStmt.Exec(modifyUserId, file.ID, file.BIllHid, file.Ts)
		if errDelFile != nil {
			resStatus = pub.StatusInternalError
			zap.L().Error("TrainCourse.Delete delFileStmt.Exec failed", zap.Error(errDelFile))
			tx.Rollback()
			return
		}
		delFileNumber, errDelEff := delFileRes.RowsAffected()
		if errDelEff != nil {
			resStatus = pub.StatusInternalError
			zap.L().Error("TrainCourse.Delete delFileRes.RowsAffected failed", zap.Error(errDelEff))
			tx.Rollback()
			return
		}
		if delFileNumber < 1 {
			resStatus = pub.StatusOtherEdit
			tx.Rollback()
			return
		}
	}

	//从缓存删除
	tc.DelFromLocalCache()
	return pub.StatusOK, nil
}

// DeleteTCs 批量删除课程
func DeleteTCs(tcs *[]TrainCourse, modifyUserID int32) (resStatus pub.ResStatus, err error) {
	//开始执行事务
	tx, err := db.Begin()
	if err != nil {
		resStatus = pub.StatusInternalError
		zap.L().Error("DeleteTCs db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	//删除文档准备
	delDocSql := `update traincourse set dr=1,modify_time=current_timestamp,modifyuserid=$1,ts=current_timestamp
		where id=$2 and dr=0 and ts=$3`
	docStmt, err := tx.Prepare(delDocSql)
	if err != nil {
		resStatus = pub.StatusInternalError
		zap.L().Error("DeleteTCs tx.Prepare(delDocSql) failed", zap.Error(err))
		tx.Rollback()
		return resStatus, err
	}
	defer docStmt.Close()
	//删除文档文件准备
	delFileSql := `update traincourse_file set dr=1,modify_time=current_timestamp,modifyuserid=$1,ts=current_timestamp
		where id=$2 and dr=0 and billhid=$3 and ts=$4`
	fileStmt, err := tx.Prepare(delFileSql)
	if err != nil {
		resStatus = pub.StatusInternalError
		zap.L().Error("DeleteTCs tx.Prepare(delFileSql) failed", zap.Error(err))
		tx.Rollback()
		return resStatus, err
	}
	defer fileStmt.Close()

	for _, tc := range *tcs {
		//检查是否被引用
		resStatus, err = tc.CheckIsUsed()
		if resStatus != pub.StatusOK || err != nil {
			tx.Rollback()
			return
		}
		//检查删除人和创建人是否同为一人
		if tc.CreateUser.UserID != modifyUserID {
			resStatus = pub.StatusOtherEdit
			tx.Rollback()
			return
		}
		//执行删除操作
		delDocRes, errDelDoc := docStmt.Exec(modifyUserID, tc.ID, tc.Ts)
		if errDelDoc != nil {
			resStatus = pub.StatusInternalError
			zap.L().Error("DeleteTCs docStmt.Exec failed", zap.Error(errDelDoc))
			tx.Rollback()
			return
		}
		//检查删除效果
		delDocNumber, errDelDocEff := delDocRes.RowsAffected()
		if errDelDocEff != nil {
			resStatus = pub.StatusInternalError
			zap.L().Error("DeleteTCs delDocRes.RowsAffected failed", zap.Error(errDelDocEff))
			tx.Rollback()
			return
		}
		if delDocNumber < 1 {
			resStatus = pub.StatusOtherEdit
			tx.Rollback()
			return
		}

		//删除文件
		for _, file := range tc.Files {
			delFileRes, errDelFile := fileStmt.Exec(modifyUserID, file.ID, file.BIllHid, file.Ts)
			if errDelFile != nil {
				resStatus = pub.StatusInternalError
				zap.L().Error("DeleteTCs fileStmt.Exec failed", zap.Error(errDelFile))
				tx.Rollback()
				return
			}
			delFileNumber, errDelEff := delFileRes.RowsAffected()
			if errDelEff != nil {
				resStatus = pub.StatusInternalError
				zap.L().Error("DeleteTCs delFileRes.RowsAffected failed", zap.Error(errDelEff))
				tx.Rollback()
				return
			}
			if delFileNumber < 1 {
				resStatus = pub.StatusOtherEdit
				tx.Rollback()
				return
			}
		}
		//从缓存删除
		tc.DelFromLocalCache()
	}

	return pub.StatusOK, nil
}

// TrainCourse.GetDetailByID 根据id获取详情
func (tc *TrainCourse) GetDetailByID() (resStatus pub.ResStatus, err error) {
	//从localcache中获取
	number, b, _ := cache.Get(pub.TC, tc.ID)
	if number > 0 {
		json.Unmarshal(b, &tc)
		resStatus = pub.StatusOK
		return
	}
	//从数据库获取
	sqlStr := `select code,name,classhour,isexamine,description,
		create_time,createuserid,modify_time,modifyuserid,ts,
		dr
		from traincourse where id=$1`
	err = db.QueryRow(sqlStr, tc.ID).Scan(&tc.Code, &tc.Name, &tc.ClassHour, &tc.IsExamine, &tc.Description,
		&tc.CreateDate, &tc.CreateUser.UserID, &tc.ModifyDate, &tc.ModifyUser.UserID, &tc.Ts,
		&tc.Dr)
	if err != nil {
		zap.L().Error("TrainCourse.GetDetailByID db.QueryRow  failed", zap.Error(err))
		resStatus = pub.StatusInternalError
		return
	}
	//获取创建人详情
	if tc.CreateUser.UserID > 0 {
		resStatus, err = tc.CreateUser.GetPersonInfoByID()
		if resStatus != pub.StatusOK || err != nil {
			return
		}
	}
	//获取更新人详情
	if tc.ModifyUser.UserID > 0 {
		resStatus, err = tc.ModifyUser.GetPersonInfoByID()
		if resStatus != pub.StatusOK || err != nil {
			return
		}
	}

	//获取附件
	filesStr := `select id,billhid,fileid,create_time,createuserid,
		modify_time,modifyuserid,ts,dr
		from traincourse_file
		where dr=0 and billhid=$1`
	rows, err := db.Query(filesStr, tc.ID)
	if err != nil {
		zap.L().Error("traincourse.GetDetailByID db.Query(files) failed", zap.Error(err))
		resStatus = pub.StatusInternalError
		return
	}

	for rows.Next() {
		var file VoucherFile
		err = rows.Scan(&file.ID, &file.BIllHid, &file.File.FileId, &file.CreateDate, &file.CreateUser.UserID,
			&file.ModifyDate, &file.ModifyUser.UserID, &file.Ts, &file.Dr)
		if err != nil {
			zap.L().Error("traincourse.GetDetailByID db.Query(files)  rows.Scan  failed", zap.Error(err))
			resStatus = pub.StatusInternalError
			return
		}

		//填充file
		if file.File.FileId > 0 {
			resStatus, err = file.File.GetFileInfoByID()
			if resStatus != pub.StatusOK || err != nil {
				return
			}
		}
		//填充创建人
		if file.CreateUser.UserID > 0 {
			resStatus, err = file.CreateUser.GetPersonInfoByID()
			if resStatus != pub.StatusOK || err != nil {
				return
			}
		}
		//填充更新人
		if file.ModifyUser.UserID > 0 {
			resStatus, err = file.ModifyUser.GetPersonInfoByID()
			if resStatus != pub.StatusOK || err != nil {
				return
			}
		}

		tc.Files = append(tc.Files, file)
	}

	//写入localcache
	tcB, _ := json.Marshal(tc)
	cache.Set(pub.Document, tc.ID, tcB)

	return pub.StatusOK, nil
}

// GetTCList 获取培训课程列表
func GetTCList() (tcs []TrainCourse, resStatus pub.ResStatus, err error) {
	tcs = make([]TrainCourse, 0)
	//从数据库获取列表记录
	sqlStr := `select id from traincourse where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("GetTCList db.Query failed", zap.Error(err))
		resStatus = pub.StatusInternalError
		return
	}
	defer rows.Close()

	//记录从rows中获取的数据行数
	var rowsNum int32
	//从查询记录中逐行获取数据
	for rows.Next() {
		rowsNum++
		var tc TrainCourse
		err = rows.Scan(&tc.ID)
		if err != nil {
			zap.L().Error("GetTCList row.Next failed", zap.Error(err))
			resStatus = pub.StatusInternalError
			return
		}
		//获取详情
		resStatus, err = tc.GetDetailByID()
		if resStatus != pub.StatusOK || err != nil {
			return
		}

		//追加数组
		tcs = append(tcs, tc)
	}

	//如果获取的数据行数等于0,则表示没有数据
	if rowsNum == 0 {
		resStatus = pub.StatusResNoData
		return
	}

	resStatus = pub.StatusOK
	return
}

// GetSimpTCCache 获取课程前端缓存
func (tcc *TCCache) GetTCCache() (resStatus pub.ResStatus, err error) {
	tcc.DelItems = make([]TrainCourse, 0)
	tcc.NewItems = make([]TrainCourse, 0)
	tcc.UpdateItems = make([]TrainCourse, 0)
	//查询课程档案最新ts
	sqlStr := `select ts from traincourse where ts > $1 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, tcc.QueryTs).Scan(&tcc.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			tcc.ResultNumber = 0
			tcc.ResultTs = tcc.QueryTs
			resStatus = pub.StatusOK
			return
		}
		zap.L().Error("TCCache.GetTCCache query latest ts failed", zap.Error(err))
		resStatus = pub.StatusInternalError
		return
	}

	//查询所有大于QueryTs的数据

	sqlStr = `select id
	from traincourse
	where ts > $1 order by ts desc`
	rows, err := db.Query(sqlStr, tcc.QueryTs)
	if err != nil {
		zap.L().Error("TCCache.GetTCCache get cache from database failed", zap.Error(err))
		resStatus = pub.StatusInternalError
		return
	}
	defer rows.Close()

	//从查询结果中提取数据
	for rows.Next() {
		var tc TrainCourse

		err = rows.Scan(&tc.ID)
		if err != nil {
			zap.L().Error("TCCache.GetTCCache rows.Next() failed", zap.Error(err))
			resStatus = pub.StatusInternalError
			return
		}
		resStatus, err = tc.GetDetailByID()
		if resStatus != pub.StatusOK || err != nil {
			return
		}


		if tc.Dr == 0 {
			//档案还没有被删除
			if tc.CreateDate.Before(tcc.QueryTs) || tc.CreateDate.Equal(tcc.QueryTs) {
				//QueryTS之前增加的,说明已经被修改过了
				tcc.ResultNumber++
				tcc.UpdateItems = append(tcc.UpdateItems, tc)
			} else {
				//QueryTs之后新增的
				tcc.ResultNumber++
				tcc.NewItems = append(tcc.NewItems, tc)
			}
		} else {
			//档案已经被删除
			if tc.CreateDate.Before(tcc.QueryTs) || tc.CreateDate.Equal(tcc.QueryTs) {
				//queryTs之前增加的,表明是QueryTS之后删除的
				tcc.ResultNumber++
				tcc.DelItems = append(tcc.DelItems, tc)
			}
			//QueryTs之后增加并且之后删除的,不需要进行处理
		}
	}

	return pub.StatusOK, nil
}

// TrainCourse.CheckIsUsed 检查培训课程是否被引用
func (tc *TrainCourse) CheckIsUsed() (resStatus pub.ResStatus, err error) {
	//创建数据被引用检查类型切片
	checkItems := []ScDocCheckUsed{
		{
			Description:    "被培训记录引用",
			SqlStr:         `select count(id) as usedNum from trainrecord_h where tc_id=$1 and dr=0`,
			UsedReturnCode: pub.StatusTRUsed,
		},
	}
	//检查项目
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, tc.ID).Scan(&usedNum)
		if err != nil {
			resStatus = pub.StatusInternalError
			zap.L().Error("TrainCourse.CheckIsUsed "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return pub.StatusOK, nil
} */
