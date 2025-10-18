package pg

import (
	"sccsmsserver/i18n"
	"sccsmsserver/setting"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Issue Resolution Form struct
type IssueResolutionForm struct {
	ID                 int32            `db:"id" json:"id"`                                 //ID
	BillNumber         string           `db:"billnumber" json:"billnumber"`                 //单据编号
	BillDate           time.Time        `db:"billdate" json:"billdate"`                     //单据日期
	CSA                ConstructionSite `db:"csaid" json:"csa"`                             //现场
	EPA                ExecutionProject `db:"epaid" json:"eid"`                             //执行项目
	ExecutionValue     string           `db:"executionvalue" json:"executionvalue"`         //执行值
	ExecutionValueDisp string           `db:"executionvaluedisp" json:"executionvaluedisp"` //执行值显示
	Executor           Person           `db:"executorid" json:"executor"`                   //执行人
	Department         SimpDept         `db:"deptid" json:"department"`                     //部门
	Fixer              Person           `db:"fixerid" json:"fixer"`                         //处理人
	IsFinish           int16            `db:"isfinish" json:"isfinish"`                     //是否处理完成
	StartTime          time.Time        `db:"starttime" json:"starttime"`                   //开始时间
	EndTime            time.Time        `db:"endtime" json:"endtime"`                       //结束时间
	EODescription      string           `db:"eodescription" json:"eodescription"`           //问题说明
	Description        string           `db:"description" json:"description"`               //说明
	Status             int16            `db:"status" json:"status"`                         //状态
	SourceType         string           `db:"sourcetype" json:"sourcetype"`                 //来源单据类型 ED:执行单
	SourceBillNumber   string           `db:"sourcebillnumber" json:"sourcebillnumber"`     //来源单据号
	SourceHID          int32            `db:"sourcehid" json:"sourcehid"`                   //来源单据表头ID
	SourceRowNumber    int32            `db:"sourcerownumber" json:"sourcerownumber"`       //来源单据行号
	SourceBID          int32            `db:"sourcebid" json:"sourcebid"`                   //来源单据表体ID
	RiskLevel          RiskLevel        `db:"risklevelid" json:"risklevel"`                 //风险等级
	SourceRowTs        time.Time        `json:"sourcerowts"`                                //来源单据表体ts
	IssueFiles         []VoucherFile    `json:"issueFiles"`                                 //问题附件
	FixFiles           []VoucherFile    `json:"fixFiles"`                                   //处理结果附件
	CreateDate         time.Time        `db:"createtime" json:"createdate"`                 //创建日期
	Creator            Person           `db:"creatorid" json:"createuser"`                  //创建人
	ConfirmDate        time.Time        `db:"confirmtime" json:"confirmdate"`               //确认时间
	Confirmer          Person           `db:"confirmerid" json:"confirmuser"`               //确认人
	ModifyDate         time.Time        `db:"modifytime" json:"modifydate"`                 //修改日期
	Modifier           Person           `db:"modifierid" json:"modifyuser"`                 //修改人
	Ts                 time.Time        `db:"ts" json:"ts"`                                 //时间戳
	Dr                 int16            `db:"dr" json:"dr"`                                 //删除标志
}

// Add Issue Resolution Form
func (irf *IssueResolutionForm) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK

	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("IssueResolutionForm.Add db.Begin() failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Get the latest serial number
	irf.BillNumber, resStatus, err = GetLatestSerialNo(tx, "IRF", irf.BillDate.Format("060102"))
	if resStatus != i18n.StatusOK || err != nil {
		tx.Rollback()
		return
	}
	// Insert into IRF data to
	addSql := `insert into issueresolutionform(billnumber,billdate,csaid,epaid,executionvalue,
	executionvaluedisp,executorid,deptid,fixerid,isfinish,
	starttime,endtime,eodescription,description,status,
	sourcetype,sourcebillnumber,sourcehid,sourcerownumber,sourcebid,
	risklevelid,creatorid)
	values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)  
	returning id`
	err = tx.QueryRow(addSql, irf.BillNumber, irf.BillDate, irf.CSA.ID, irf.EPA.ID, irf.ExecutionValue,
		irf.ExecutionValueDisp, irf.Executor.ID, irf.Department.ID, irf.Fixer.ID, irf.IsFinish,
		irf.StartTime, irf.EndTime, irf.EODescription, irf.Description, irf.Status,
		irf.SourceType, irf.SourceBillNumber, irf.SourceHID, irf.SourceRowNumber, irf.SourceBID,
		irf.RiskLevel.ID, irf.Creator.ID).Scan(&irf.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("IssueResolutionForm.Add tx.QueryRow(addsql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Add Attachments
	if len(irf.FixFiles) > 0 {
		// Prepare write the file record to the issueresolutionform_file table
		fileSql := `insert into issueresolutionform_file(billbid,fileid,creatorid) 
	            values($1,$2,$3) returning id`
		fileStmt, fileErr := tx.Prepare(fileSql)
		if fileErr != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("IssueResolutionForm.Add tx.Prepare(fileSql) failed", zap.Error(fileErr))
			tx.Rollback()
			return resStatus, fileErr
		}
		defer fileStmt.Close()

		// Write the file record to database item by item
		for _, f := range irf.FixFiles {
			fileErr = fileStmt.QueryRow(irf.ID, f.File.ID, irf.Creator.ID).Scan(&f.ID)
			if fileErr != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("IssueResolutionForm.Add fileStmt.QueryRow(fileSql) failed", zap.Error(fileErr))
				tx.Rollback()
				return resStatus, fileErr
			}
		}
	}
	// Write back the Execution Order
	if irf.SourceBID > 0 {
		edr := new(ExecutionOrderRow)
		edr.BID = irf.SourceBID
		edr.HID = irf.SourceHID
		edr.IsFinish = 1
		edr.Ts = irf.SourceRowTs
		edr.IRFID = irf.ID
		edr.IRFNumber = irf.BillNumber
		resStatus, err = edr.Dispose()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
	}
	return
}

// IssueResolutionForm.Edit 编辑问题处理单
func (irf *IssueResolutionForm) Edit() (resStatus i18n.ResKey, err error) {
	//检查创建人和编辑人是否为同一人
	if irf.Creator.ID != irf.Modifier.ID {
		resStatus = i18n.StatusVoucherOnlyCreateEdit
		return
	}
	//创建事务
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("IssueResolutionForm.Edit db.Begin() failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	//修改单据
	editSql := `update issueresolutionform set billdate=$1,deptid=$2,fixerid=$3,isfinish=$4,starttime=$5,
	endtime=$6,	description=$7,modifytime=current_timestamp,modifierid=$8,ts=current_timestamp 
	where id=$9 and dr=0 and status=0 and ts=$10`
	editRes, err := tx.Exec(editSql, irf.BillDate, irf.Department.ID, irf.Fixer.ID, irf.IsFinish, irf.StartTime,
		irf.EndTime, irf.Description, irf.Modifier.ID,
		irf.ID, irf.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("IssueResolutionForm.Edit tx.Exec(editSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}

	//检查表头修改的行数
	updateNumber, err := editRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("IssueResolutionForm.Edit editRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if updateNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}

	//修改文件准备
	updateFileSql := `update issueresolutionform_file set modifytime=current_timestamp,modifierid=$1,dr=$2,ts=current_timestamp
	where id=$3 and dr=0 and ts=$4`
	updateFileStmt, err := tx.Prepare(updateFileSql)
	if err != nil {
		zap.L().Error("IssueResolutionForm.Edit tx.Prepare(updateFileSql) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer updateFileStmt.Close()
	//增加文件准备
	addFileSql := `insert into issueresolutionform_file(billbid,fileid,creatorid) 
	values($1,$2,$3) returning id`
	addFileStmt, err := tx.Prepare(addFileSql)
	if err != nil {
		zap.L().Error("IssueResolutionForm.Edit tx.Prepare(addFileSql) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer addFileStmt.Close()

	//文件写入
	if len(irf.FixFiles) > 0 {
		for _, file := range irf.FixFiles {
			if file.ID == 0 { //新增附件
				addFileErr := addFileStmt.QueryRow(irf.ID, file.File.ID, irf.Modifier.ID).Scan(&file.ID)
				if addFileErr != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("IssueResolutionForm.Edit old row addFileStmt.QueryRow failed", zap.Error(addFileErr))
					tx.Rollback()
					return resStatus, addFileErr
				}
			} else { //原有附件
				updateFileRes, updateFileErr := updateFileStmt.Exec(irf.Modifier.ID, file.Dr, file.ID, file.Ts)
				if updateFileErr != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("IssueResolutionForm.Edit old row updateFileRes.Exec() failed", zap.Error(updateFileErr))
					tx.Rollback()
					return resStatus, updateFileErr
				}
				updateFileNumber, updateFileEffectErr := updateFileRes.RowsAffected()
				if updateFileEffectErr != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("IssueResolutionForm.Edit old row updateFileRes.RowsAffected() failed", zap.Error(updateFileEffectErr))
					tx.Rollback()
					return resStatus, updateFileEffectErr
				}
				if updateFileNumber < 1 {
					resStatus = i18n.StatusOtherEdit
					tx.Rollback()
					return
				}
			}
		}
	}

	return i18n.StatusOK, nil
}

// IssueResolutionForm.Delete 删除问题处理单
func (irf *IssueResolutionForm) Delete(modifyUserId int32) (resStatus i18n.ResKey, err error) {
	//检查单据状态
	if irf.Status != 0 {
		resStatus = i18n.StatusVoucherNoFree
		return
	}

	//检查创建人和删除人是否为同一人
	if irf.Creator.ID != modifyUserId {
		resStatus = i18n.StatusVoucherOnlyCreateEdit
		return
	}

	//创建事务
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("IssueResolutionForm.Delete db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	//删除
	delSql := `update issueresolutionform set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	delRes, err := tx.Exec(delSql, modifyUserId, irf.ID, irf.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("IssueResolutionForm.Delete tx.Exec(delHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	//检查删除行数
	delNumber, err := delRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("IssueResolutionForm.Delete delRes.RowsAffected() failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if delNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}

	//删除附件
	if len(irf.FixFiles) > 0 {
		delFileSql := `update issueresolutionform_file set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
		where id=$2 and dr=0 and billbid=$3 and ts=$4`
		//删除文件预处理
		delFileStmt, delFileErr := tx.Prepare(delFileSql)
		if delFileErr != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("IssueResolutionForm.Delete tx.Prepare(delFileSql) failed", zap.Error(delFileErr))
			tx.Rollback()
			return resStatus, delFileErr
		}
		defer delFileStmt.Close()

		for _, row := range irf.FixFiles {
			delFileRes, delFileErr := delFileStmt.Exec(modifyUserId, row.ID, row.BillBID, row.Ts)
			if delFileErr != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("IssueResolutionForm.Delete delFileStmt.Exec() failed", zap.Error(delFileErr))
				tx.Rollback()
				return resStatus, delFileErr
			}
			//检查删除影响行数
			delFileNumber, delFileErr := delFileRes.RowsAffected()
			if delFileErr != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("IssueResolutionForm.Delete delFileRes.RowsAffected() failed", zap.Error(delFileErr))
				tx.Rollback()
				return resStatus, delFileErr
			}
			if delFileNumber < 1 {
				resStatus = i18n.StatusOtherEdit
				tx.Rollback()
				return
			}
		}
	}
	//回写执行单
	if irf.SourceBID > 0 {
		edr := new(ExecutionOrderRow)
		edr.BID = irf.SourceBID
		edr.HID = irf.SourceHID
		edr.IsFinish = 0
		edr.IRFID = 0
		edr.IRFNumber = ""
		resStatus, err = edr.CancelDispose()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
	}

	return i18n.StatusOK, nil
}

// IssueResolutionForm.Confirm 确认问题处理单
func (irf *IssueResolutionForm) Confirm(confirmUserID int32) (resStatus i18n.ResKey, err error) {
	//检查单据状态
	if irf.Status != 0 { //非自由态单据不允许确认
		resStatus = i18n.StatusVoucherNoFree
		return
	}
	//创建事务
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("IssueResolutionForm.Confirm db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	//单据确认
	sqlStr := `update issueresolutionform set status=1,confirmtime=current_timestamp,confirmerid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and status=0 and ts=$3`
	confirmRes, err := tx.Exec(sqlStr, confirmUserID, irf.ID, irf.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("IssueResolutionForm.Confirm tx.Exec(sqlStr) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	//检查更新行数
	updateNumber, err := confirmRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("IssueResolutionForm.Confirm confirmRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if updateNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}

	//回写执行单
	if irf.SourceBID > 0 {
		edr := new(ExecutionOrderRow)
		edr.BID = irf.SourceBID
		edr.HID = irf.SourceHID
		edr.IsFinish = 1
		edr.IRFID = irf.ID
		edr.IRFNumber = irf.BillNumber
		resStatus, err = edr.Complete()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
	}

	return
}

// IssueResolutionForm.CancelConfirm 取消确认问题处理单
func (irf *IssueResolutionForm) CancelConfirm(confirmUserID int32) (resStatus i18n.ResKey, err error) {
	//检查单据状态
	if irf.Status != 1 { //非确认态单据不允许确认
		resStatus = i18n.StatusVoucherNoConfirm
		return
	}
	if irf.Confirmer.ID != confirmUserID {
		resStatus = i18n.StatusVoucherCancelConfirmSelf
		return
	}
	//创建事务
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("IssueResolutionForm.CancelConfirm db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	//单据取消确认
	sqlStr := `update issueresolutionform set status=0,confirmerid=0,ts=current_timestamp 
	where id=$1 and dr=0 and status=1 and ts=$2`
	confirmRes, err := tx.Exec(sqlStr, irf.ID, irf.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("IssueResolutionForm.CancelConfirm tx.Exec(sqlStr) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	//检查更新行数
	updateNumber, err := confirmRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("IssueResolutionForm.Confirm confirmRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if updateNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}

	//回写执行单
	if irf.SourceBID > 0 {
		edr := new(ExecutionOrderRow)
		edr.BID = irf.SourceBID
		edr.HID = irf.SourceHID
		edr.IsFinish = 1
		edr.IRFID = irf.ID
		edr.IRFNumber = irf.BillNumber
		resStatus, err = edr.CancelComplete()
		if resStatus != i18n.StatusOK || err != nil {
			tx.Rollback()
			return
		}
	}
	return i18n.StatusOK, nil
}

// GetDDList 获取处理单列表
func GetDDList(queryString string) (dds []IssueResolutionForm, resStatus i18n.ResKey, err error) {
	var build strings.Builder
	//拼接检查sql
	build.WriteString(`select count(b.id) as rownumber 
	from issueresolutionform as b
	left join csa as si on b.csaid = si.id
	left join epa as eid on b.epaid = eid.id
	left join sysuser as executor on b.executorid = executor.id
	left join sysuser as fixer on b.fixerid = fixer.id
	left join department as dept on b.deptid = dept.id
	where (b.dr=0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	checkSql := build.String()
	//检查
	var rowNumber int32
	err = db.QueryRow(checkSql).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetDDList db.QueryRow(checkSql) failed", zap.Error(err))
		return
	}
	if rowNumber == 0 { //如果查询数据量为0
		resStatus = i18n.StatusResNoData
		return
	}
	if rowNumber > setting.Conf.PqConfig.MaxRecord { //查询数据量大于最大记录数
		resStatus = i18n.StatusOverRecord
		return
	}
	build.Reset() //清空build
	//拼接正式sql
	build.WriteString(`select b.id,b.billnumber,b.billdate,b.csaid,b.epaid,
	b.executionvalue,b.executionvaluedisp,b.executorid,b.deptid,b.fixerid,
	b.isfinish,b.starttime,b.endtime,b.eodescription,b.description,
	b.status,b.sourcetype,b.sourcebillnumber,b.sourcehid,b.sourcerownumber,
	b.sourcebid,b.risklevelid,b.createtime,b.creatorid,confirmtime,
	confirmerid,b.modifytime,b.modifierid,b.dr,b.ts 
	from issueresolutionform as b
	left join csa as si on b.csaid = si.id
	left join epa as eid on b.epaid = eid.id
	left join sysuser as executor on b.executorid = executor.id
	left join sysuser as fixer on b.fixerid = fixer.id
	left join department as dept on b.deptid = dept.id
	where (b.dr=0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	ddsSql := build.String()
	//获取指令单列表
	ddsRows, err := db.Query(ddsSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetDDList db.Query failed", zap.Error(err))
		return
	}
	defer ddsRows.Close()
	//提取数据
	for ddsRows.Next() {
		var irf IssueResolutionForm
		err = ddsRows.Scan(&irf.ID, &irf.BillNumber, &irf.BillDate, &irf.CSA.ID, &irf.EPA.ID,
			&irf.ExecutionValue, &irf.ExecutionValueDisp, &irf.Executor.ID, &irf.Department.ID, &irf.Fixer.ID,
			&irf.IsFinish, &irf.StartTime, &irf.EndTime, &irf.EODescription, &irf.Description,
			&irf.Status, &irf.SourceType, &irf.SourceBillNumber, &irf.SourceHID, &irf.SourceRowNumber,
			&irf.SourceBID, &irf.RiskLevel.ID, &irf.CreateDate, &irf.Creator.ID, &irf.ConfirmDate,
			&irf.Confirmer.ID, &irf.ModifyDate, &irf.Modifier.ID, &irf.Dr, &irf.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetDDList ddsRows.Scan() failed", zap.Error(err))
			return
		}
		//填充现场信息
		if irf.CSA.ID > 0 {
			resStatus, err = irf.CSA.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充执行项目信息
		if irf.EPA.ID > 0 {
			resStatus, err = irf.EPA.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充执行人信息
		if irf.Executor.ID > 0 {
			resStatus, err = irf.Executor.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充部门信息
		if irf.Department.ID > 0 {
			resStatus, err = irf.Department.GetSimpDeptInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充处理人信息
		if irf.Fixer.ID > 0 {
			resStatus, err = irf.Fixer.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充创建人信息
		if irf.Creator.ID > 0 {
			resStatus, err = irf.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充风险等级信息
		if irf.RiskLevel.ID > 0 {
			resStatus, err = irf.RiskLevel.GetRLInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充确认人信息
		if irf.Confirmer.ID > 0 {
			resStatus, err = irf.Confirmer.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充更新人信息
		if irf.Modifier.ID > 0 {
			resStatus, err = irf.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		//获取执行单附件
		irf.IssueFiles, resStatus, err = GetEORowFiles(irf.SourceBID)
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		//获取处理单附件
		irf.FixFiles, resStatus, err = GetDDFiles(irf.ID)
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		dds = append(dds, irf)
	}
	resStatus = i18n.StatusOK

	return
}

// GetDDFiles 获取处理单附件
func GetDDFiles(bid int32) (voucherFiles []VoucherFile, resStatus i18n.ResKey, err error) {
	voucherFiles = make([]VoucherFile, 0) //解决返回文件为空问题
	attachSql := `select id,billbid,billhid,fileid,createtime,
	creatorid,modifytime,modifierid,dr,ts 
	from issueresolutionform_file where billbid=$1 and dr=0`
	//填充附件
	fileRows, err := db.Query(attachSql, bid)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetDDFiles db.query(attachsql) failed", zap.Error(err))
		return
	}
	defer fileRows.Close()

	for fileRows.Next() {
		var f VoucherFile
		fileErr := fileRows.Scan(&f.ID, &f.BillBID, &f.BillHID, &f.File.ID, &f.CreateDate,
			&f.Creator.ID, &f.ModifyDate, &f.Modifier.ID, &f.Dr, &f.Ts)
		if fileErr != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetDDFiles fileRows.Scan failed", zap.Error(fileErr))
			return
		}
		//填充文件信息
		if f.File.ID > 0 {
			resStatus, err = f.File.GetFileInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充创建人
		if f.Creator.ID > 0 {
			resStatus, err = f.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		//填充更新人
		if f.Modifier.ID > 0 {
			resStatus, err = f.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		voucherFiles = append(voucherFiles, f)
	}

	resStatus = i18n.StatusOK

	return
}
