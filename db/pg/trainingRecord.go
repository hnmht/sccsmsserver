package pg

import (
	"sccsmsserver/i18n"
	"sccsmsserver/setting"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Training Record Struct
type TrainingRecord struct {
	HID          int32               `db:"id" json:"id"`
	BillNumber   string              `db:"billnumber" json:"billNumber"`
	BillDate     time.Time           `db:"billdate" json:"billDate"`
	Department   SimpDept            `db:"deptid" json:"department"`
	Description  string              `db:"description" json:"description"`
	Lecturer     Person              `db:"lecturerid" json:"lecturer"`
	TrainingDate time.Time           `db:"trainingdate" json:"trainingDate"`
	TC           TC                  `db:"tcid" json:"tc"`
	StartTime    time.Time           `db:"starttime" json:"startTime"`
	EndTime      string              `db:"endtime" json:"endTime"`
	ClassHour    float64             `db:"classhour" json:"classHour"`
	IsExam       int16               `db:"isExam" json:"isExamine"`
	HFiles       []VoucherFile       `json:"hFiles"`
	Body         []TrainingRecordRow `json:"body"`
	Status       int16               `db:"status" json:"status"` // 0 free 1 confirmed 2 executing 3 completed
	CreateDate   time.Time           `db:"createtime" json:"createDate"`
	Creator      Person              `db:"creatorid" json:"creator"`
	ConfirmDate  time.Time           `db:"confirmtime" json:"confirmDate"`
	Confirmer    Person              `db:"confirmerid" json:"confirmer"`
	ModifyDate   time.Time           `db:"modifytime" json:"modifyDate"`
	Modifier     Person              `db:"modifierid" json:"modifier"`
	Ts           time.Time           `db:"ts" json:"ts"`
	Dr           int16               `db:"dr" json:"dr"`
}

// Training Record Row struct
type TrainingRecordRow struct {
	BID          int32         `db:"id" json:"id"`
	HID          int32         `db:"hid" json:"hid"`
	RowNumber    int32         `db:"rownumber" json:"rownNmber"`
	Student      Person        `db:"studentid" json:"student"`
	PositionName string        `db:"positionname" json:"positionName"`
	DeptName     string        `db:"name" json:"deptName"`
	StartTime    time.Time     `db:"starttime" json:"startTime"`
	EndTime      time.Time     `db:"endtime" json:"endTime"`
	ClassHour    float64       `db:"classhour" json:"classHour"`
	Description  string        `db:"description" json:"description"`
	ExamRes      int16         `db:"examres" json:"examRes"`
	ExamScore    float64       `db:"examscore" json:"examScore"`
	Status       int16         `db:"status" json:"status"`
	BFiles       []VoucherFile `json:"files"`
	CreateDate   time.Time     `db:"createtime" json:"createDate"`
	Creator      Person        `db:"creatorid" json:"creator"`
	ConfirmDate  time.Time     `db:"confirmtime" json:"confirmDate"`
	Confirmer    Person        `db:"confirmerid" json:"confirmer"`
	ModifyDate   time.Time     `db:"modifytime" json:"modifyDate"`
	Modifier     Person        `db:"modifierid" json:"modifier"`
	Ts           time.Time     `db:"ts" json:"ts"`
	Dr           int16         `db:"dr" json:"dr"`
}

// Taught Lessons Report
type TaughtLessonsReport struct {
	HID                    int32     `json:"hid"`
	BillNumber             string    `json:"billNumber"`
	BillDate               time.Time `json:"billDate"`
	ID                     int32     `json:"deptID"`
	DeptCode               string    `json:"deptCode"`
	DeptName               string    `json:"deptName"`
	Description            string    `json:"description"`
	LecturerID             int32     `json:"lecturerID"`
	LecturerCode           string    `json:"lecturerCode"`
	LecturerName           string    `json:"lecturerName"`
	TrainingDate           time.Time `json:"trainingDate"`
	TCID                   int32     `json:"tcID"`
	TCCode                 string    `json:"tcCode"`
	TCName                 string    `json:"tcName"`
	StartTime              time.Time `json:"startTime"`
	EndTime                time.Time `json:"endTime"`
	ClassHour              float64   `json:"classHour"`
	IsExam                 int16     `json:"isExam"`
	StudentNumber          int32     `json:"studentNumber"`
	QualifiedNumber        int32     `json:"qualifiedNumber"`
	DisqualificationNumber int32     `json:"disqualificationNumber"`
	Status                 int16     `json:"status"`
	CreatorID              int32     `json:"creatorID"`
	CreatorCode            string    `json:"creatorCode"`
	CreatorName            string    `json:"creatorName"`
}

// Recived Training Report
type RecivedTrainingReport struct {
	HID                 int32     `json:"hid"`
	BID                 int32     `json:"bid"`
	BillNumber          string    `json:"billNumber"`
	BillDate            time.Time `json:"billDate"`
	ID                  int32     `json:"deptID"`
	DeptCode            string    `json:"deptCode"`
	DeptName            string    `json:"deptName"`
	LecturerID          int32     `json:"lecturerID"`
	LecturerCode        string    `json:"lecturerCode"`
	LecturerName        string    `json:"lecturerName"`
	TCID                int32     `json:"tcID"`
	TCCode              string    `json:"tcCode"`
	TCName              string    `json:"tcName"`
	StartTime           time.Time `json:"startTime"`
	EndTime             time.Time `json:"endTime"`
	TCClassHour         float64   `json:"tcClassHour"`
	IsExam              int16     `json:"isExam"`
	HStatus             int16     `json:"hStatus"`
	HDescription        string    `json:"hDescription"`
	StudentID           int32     `json:"studentID"`
	StudentCode         string    `json:"studentCode"`
	StudentName         string    `json:"studentName"`
	StudentPositionName string    `json:"studentPositionName"`
	StudentDeptName     string    `json:"studentDeptName"`
	SignStartTime       time.Time `json:"signStartTime"`
	SignEndTime         time.Time `json:"signEndTime"`
	BClassHour          float64   `json:"bClassHour"`
	BDescription        string    `json:"bDescription"`
	ExamRes             int16     `json:"examRes"`
	ExamScore           float64   `json:"examScore"`
	BStatus             int16     `json:"bStatus"`
	CreatorID           int32     `json:"creatorID"`
	CreatorCode         string    `json:"creatorCode"`
	CreatorName         string    `json:"creatorName"`
}

// Add Training Record
func (tr *TrainingRecord) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check the number of body rows, zero is not allowed
	if len(tr.Body) == 0 {
		resStatus = i18n.StatusVoucherNoBody
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Add db.Begin failed:", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Get the latest serial number
	billNo, resStatus, err := GetLatestSerialNo(tx, "TR", tr.BillDate.Format("060102"))
	if resStatus != i18n.StatusOK || err != nil {
		tx.Rollback()
		return
	}
	tr.BillNumber = billNo
	// Insert data into the trainingrecord_h table
	headSql := `insert into trainingrecord_h(billnumber,billdate,deptid,description,lecturerid,
		trainingdate,tcid,starttime,endtime,classhour,
		isexam,status,creatorid) 
		values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) 
		returning id`
	err = tx.QueryRow(headSql, tr.BillNumber, tr.BillDate, tr.Department.ID, tr.Description, tr.Lecturer.ID,
		tr.TrainingDate, tr.TC.ID, tr.StartTime, tr.EndTime, tr.ClassHour,
		tr.IsExam, tr.Status, tr.Creator.ID).Scan(&tr.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Add tx.QeuryRow(headSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	// Prepare insert the header attachment record into the trainingrecord_file table
	headFileSql := `insert into trainingrecord_file(billhid,fileid,creatorid) 
	values($1,$2,$3) returning id`
	headFileStmt, err := tx.Prepare(headFileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Add tx.Prepare(headFileSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	defer headFileStmt.Close()
	// Insert the header attachement record row by row
	for _, hFile := range tr.HFiles {
		err = headFileStmt.QueryRow(tr.HID, hFile.File.ID, tr.Creator.ID).Scan(&hFile.ID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("TrainingRecord.Add headFileStmt.QueryRow failed:", zap.Error(err))
			tx.Rollback()
			return
		}
	}

	// Perpare insert data into the trainingrecord_b table
	bodySql := `insert into trainingrecord_b(hid,rownumber,studentid,positionname,name,
		starttime,endtime,classhour,description,examres,
		examscore,status,creatorid)
		values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) 
		returning id`
	bodyStmt, err := tx.Prepare(bodySql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Add tx.Prepare(bodySql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	defer bodyStmt.Close()
	// Prepare insert the body attachment into the trainingrecord_file table
	fileSql := `insert into trainingrecord_file(billbid,billhid,fileid,creatorid) 
	values($1,$2,$3,$4) returning id`
	fileStmt, err := tx.Prepare(fileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainRecordAdd tx.Prepare(fileSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	defer fileStmt.Close()
	// Insert Traning Record row Data into the trainingrecord_b table row by row
	for _, row := range tr.Body {
		err = bodyStmt.QueryRow(tr.HID, row.RowNumber, row.Student.ID, row.PositionName, row.DeptName,
			row.StartTime, row.EndTime, row.ClassHour, row.Description, row.ExamRes,
			row.ExamScore, row.Status, tr.Creator.ID).Scan(&row.BID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("TrainingRecord.Add bodyStmt.QueryRow failed:", zap.Error(err))
			tx.Rollback()
			return
		}
		// Insert Training Record row Attachment data into trainingrecord_file item by item
		if len(row.BFiles) > 0 {
			for _, file := range row.BFiles {
				err = fileStmt.QueryRow(row.BID, tr.HID, file.File.ID, tr.Creator.ID).Scan(&file.ID)
				if err != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("TrainRecordAdd fileStmt.QueryRow failed:", zap.Error(err))
					tx.Rollback()
					return
				}
			}
		}
	}
	return
}

// Get Training Record details by HID
func (tr *TrainingRecord) GetDetailByHID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the bill has been deleted
	var rowNumber int32
	checkSql := `select count(id) as rownumber from trainingrecord_h where id=$1 and dr=0`
	err = db.QueryRow(checkSql, tr.HID).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.GetDetailByHID db.QueryRow(checkSql) failed:", zap.Error(err))
		return
	}
	if rowNumber < 1 {
		resStatus = i18n.StatusDataDeleted
		return
	}
	// Get body details
	resStatus, err = tr.FillBody()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}

	return
}

// Get Training Record body details
func (tr *TrainingRecord) FillBody() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Retrieve Training Record body rows from the trainingrecord_b table
	bodySql := `select id,hid,rownumber,studentid,positionname,
		name,starttime,endtime,classhour,description,
		examres,examscore,status,createtime,creatorid,
		confirmtime,confirmerid,modifytime,modifierid,dr,
		ts from trainingrecord_b
		where hid=$1 and dr=0 order by rownumber asc`
	bodyRows, err := db.Query(bodySql, tr.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.FillBody db.Query(bodySql) failed:", zap.Error(err))
		return
	}
	defer bodyRows.Close()
	var bodyRowNumber int32
	// Extract data row by row
	for bodyRows.Next() {
		bodyRowNumber++
		var trr TrainingRecordRow
		err = bodyRows.Scan(&trr.BID, &trr.HID, &trr.RowNumber, &trr.Student.ID, &trr.PositionName,
			&trr.DeptName, &trr.StartTime, &trr.EndTime, &trr.ClassHour, &trr.Description,
			&trr.ExamRes, &trr.ExamScore, &trr.Status, &trr.CreateDate, &trr.Creator.ID,
			&trr.ConfirmDate, &trr.Confirmer.ID, &trr.ModifyDate, &trr.Modifier.ID, &trr.Dr,
			&trr.Ts)
		if err != nil {
			zap.L().Error("TrainingRecord.FillBody bodyRows.scan failed:", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get student details
		if trr.Student.ID > 0 {
			resStatus, err = trr.Student.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		// Get Creator details
		if trr.Creator.ID > 0 {
			resStatus, err = trr.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Comfirmer details
		if trr.Confirmer.ID > 0 {
			resStatus, err = trr.Confirmer.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier details
		if trr.Modifier.ID > 0 {
			resStatus, err = trr.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get row attachments
		trr.BFiles, resStatus, err = GetTRRFiles(trr.BID)
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		tr.Body = append(tr.Body, trr)
	}

	return
}

// Get Training Record row attachments
func GetTRRFiles(bid int32) (voucherFiles []VoucherFile, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	voucherFiles = make([]VoucherFile, 0)
	// Get Training Record row Attachments from trainingrecord_file table
	attachSql := `select id,billbid,billhid,fileid,createtime,
	creatorid,modifytime,modifierid,dr,ts 
	from trainingrecord_file where billbid=$1 and dr=0`
	fileRows, err := db.Query(attachSql, bid)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetTRRFiles db.query(attachsql) failed:", zap.Error(err))
		return
	}
	defer fileRows.Close()
	// Extract data row by row
	for fileRows.Next() {
		var f VoucherFile
		fileErr := fileRows.Scan(&f.ID, &f.BillBID, &f.BillHID, &f.File.ID, &f.CreateDate,
			&f.Creator.ID, &f.ModifyDate, &f.Modifier.ID, &f.Dr, &f.Ts)
		if fileErr != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetTRRFiles fileRows.Scan failed:", zap.Error(fileErr))
			return
		}
		// Get file details
		if f.File.ID > 0 {
			resStatus, err = f.File.GetFileInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get file creator details
		if f.Creator.ID > 0 {
			resStatus, err = f.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get file modifier details
		if f.Modifier.ID > 0 {
			resStatus, err = f.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		voucherFiles = append(voucherFiles, f)
	}
	return
}

// Get Training Record header files by HID
func GetTRHFiles(hid int32) (voucherFiles []VoucherFile, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	voucherFiles = make([]VoucherFile, 0)
	// Retrieve Training Record header file from trainingrecord_file table
	attachSql := `select id,billbid,billhid,fileid,createtime,
	creatorid,modifytime,modifierid,dr,ts 
	from trainingrecord_file where billhid=$1 and dr=0 and billbid=0`
	fileRows, err := db.Query(attachSql, hid)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetTRHFiles db.query(attachsql) failed:", zap.Error(err))
		return
	}
	defer fileRows.Close()
	// Extract data row by row
	for fileRows.Next() {
		var f VoucherFile
		fileErr := fileRows.Scan(&f.ID, &f.BillBID, &f.BillHID, &f.File.ID, &f.CreateDate,
			&f.Creator.ID, &f.ModifyDate, &f.Modifier.ID, &f.Dr, &f.Ts)
		if fileErr != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetTRHFiles fileRows.Scan failed:", zap.Error(fileErr))
			return
		}
		// Get file details
		if f.File.ID > 0 {
			resStatus, err = f.File.GetFileInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get file creator details
		if f.Creator.ID > 0 {
			resStatus, err = f.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get file modifier details
		if f.Modifier.ID > 0 {
			resStatus, err = f.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		voucherFiles = append(voucherFiles, f)
	}
	return
}

// Fill Training Record header details
func (tr *TrainingRecord) FillHead() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get Department details
	if tr.Department.ID > 0 {
		resStatus, err = tr.Department.GetSimpDeptInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Lecturer details
	if tr.Lecturer.ID > 0 {
		resStatus, err = tr.Lecturer.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Training Course details
	if tr.TC.ID > 0 {
		resStatus, err = tr.TC.GetDetailByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}

	// Get Creator details
	if tr.Creator.ID > 0 {
		resStatus, err = tr.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Confimer details
	if tr.Confirmer.ID > 0 {
		resStatus, err = tr.Confirmer.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Modifier details
	if tr.Modifier.ID > 0 {
		resStatus, err = tr.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Header Attaments details
	tr.HFiles, resStatus, err = GetTRHFiles(tr.HID)
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	return
}

// Get Training Record List
func GetTRList(queryString string) (trs []TrainingRecord, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	trs = make([]TrainingRecord, 0)
	var build strings.Builder
	// Concatenate SQL string for check
	build.WriteString(`select count(h.id) as rownumber
	from trainingrecord_h as h
	left join department on h.deptid = department.id
	left join tc on h.tcid = tc.id 
	left join sysuser as lecturer on h.lecturerid = lecturer.id
	left join sysuser as creator on h.creatorid = creator.id
	left join sysuser as modifier on h.modifierid = modifier.id
	where (h.dr = 0) `)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	checkSql := build.String()
	// Check
	var rowNumber int32
	err = db.QueryRow(checkSql).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetTRList db.QueryRow(checkSql) failed:", zap.Error(err))
		return
	}
	if rowNumber == 0 {
		resStatus = i18n.StatusResNoData
		return
	}
	if rowNumber > setting.Conf.PqConfig.MaxRecord {
		resStatus = i18n.StatusOverRecord
		return
	}
	build.Reset()

	// Concatenate SQL String for data retriveal
	build.WriteString(`select h.id,h.billnumber,h.billdate,h.deptid,h.description,
	h.lecturerid,h.trainingdate,h.tcid,h.starttime,h.endtime,
	h.classhour,h.isexam,h.status,h.createtime,h.creatorid,
	h.confirmtime,h.confirmerid,h.modifytime,h.modifierid,h.dr,
	h.ts 
	from trainingrecord_h as h
	left join department on h.deptid = department.id
	left join tc on h.tcid = tc.id 
	left join sysuser as lecturer on h.lecturerid = lecturer.id
	left join sysuser as creator on h.creatorid = creator.id
	left join sysuser as modifier on h.modifierid = modifier.id
	where (h.dr = 0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	headSql := build.String()
	// Retrieve Training Record from database
	headRows, err := db.Query(headSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetTRList db.Query failed:", zap.Error(err))
		return
	}
	defer headRows.Close()
	// Extract data row by row
	for headRows.Next() {
		var tr TrainingRecord
		err = headRows.Scan(&tr.HID, &tr.BillNumber, &tr.BillDate, &tr.Department.ID, &tr.Description,
			&tr.Lecturer.ID, &tr.TrainingDate, &tr.TC.ID, &tr.StartTime, &tr.EndTime,
			&tr.ClassHour, &tr.IsExam, &tr.Status, &tr.CreateDate, &tr.Creator.ID,
			&tr.ConfirmDate, &tr.Confirmer.ID, &tr.ModifyDate, &tr.Modifier.ID, &tr.Dr,
			&tr.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetTRList headRows.Next failed:", zap.Error(err))
			return
		}
		// Get Header details
		resStatus, err = tr.FillHead()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		trs = append(trs, tr)
	}
	return
}

// Edit Training Record
func (tr *TrainingRecord) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check the number of body rows, zero cannot allowed
	if len(tr.Body) == 0 {
		resStatus = i18n.StatusVoucherNoBody
		return
	}
	// Check if the Creator and Modifier are same person
	if tr.Creator.ID != tr.Modifier.ID {
		resStatus = i18n.StatusVoucherOnlyCreateEdit
		return
	}

	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Edit db.Begin() failed:", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Update the Training Record in the trainingrecord_h table
	trHeadSql := `update trainingrecord_h set billdate=$1,deptid=$2,description=$3, lecturerid=$4,trainingdate=$5,
	tcid=$6,starttime=$7,endtime=$8,classhour=$9,isexam=$10,
	modifytime=current_timestamp,modifierid=$11,ts=current_timestamp  
	where id=$12 and dr=0 and status=0 and ts=$13`
	trHeadRes, err := tx.Exec(trHeadSql, &tr.BillDate, &tr.Department.ID, &tr.Description, &tr.Lecturer.ID, &tr.TrainingDate,
		&tr.TC.ID, &tr.StartTime, &tr.EndTime, &tr.ClassHour, &tr.IsExam,
		&tr.Modifier.ID,
		&tr.HID, &tr.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Edit tx.Exec(tritHeadSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	headUpdateNumber, err := trHeadRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Edit EditHeadRes.RowsAffected failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	if headUpdateNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Perpare Update the header attachments record
	updateHFileSql := `update trainingrecord_file set modifytime=current_timestamp,modifierid=$1,dr=$2,ts=current_timestamp
	where id=$3 and dr=0 and ts=$4 and billbid=0`
	updateHFileStmt, err := tx.Prepare(updateHFileSql)
	if err != nil {
		zap.L().Error("TrainingRecord.Edit tx.Prepare(updateHFileSql) failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer updateHFileStmt.Close()
	// Perpare Add the header attachments record
	addHFileSql := `insert into trainingrecord_file(billhid,fileid,creatorid) 
	values($1,$2,$3) returning id`
	addHFileStmt, err := tx.Prepare(addHFileSql)
	if err != nil {
		zap.L().Error("TrainingRecord.Edit tx.Prepare(addHFileSql) failed::", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer addHFileStmt.Close()
	// Add Or Edit header attachments
	if len(tr.HFiles) > 0 {
		for _, hFile := range tr.HFiles {
			if hFile.ID == 0 { // If ID value is zero, it means it is a newly file
				addHFileErr := addHFileStmt.QueryRow(tr.HID, hFile.File.ID, tr.Modifier.ID).Scan(&hFile.ID)
				if addHFileErr != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("TrainingRecord.Edit old row addHFileStmt.QueryRow failed:", zap.Error(addHFileErr))
					tx.Rollback()
					return resStatus, addHFileErr
				}
			} else { // If ID value is not zero, it means it is a file that needs to be modified
				updateHFileRes, updateHFileErr := updateHFileStmt.Exec(tr.Modifier.ID, hFile.Dr, hFile.ID, hFile.Ts)
				if updateHFileErr != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("TrainingRecord.Edit old row updateHFileStmt.Exec() failed:", zap.Error(updateHFileErr))
					tx.Rollback()
					return resStatus, updateHFileErr
				}
				updateHFileNumber, updateHFileEffectErr := updateHFileRes.RowsAffected()
				if updateHFileEffectErr != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("TrainingRecord.Edit old row updateHFileRes.RowsAffected() failed:", zap.Error(updateHFileEffectErr))
					tx.Rollback()
					return resStatus, updateHFileEffectErr
				}

				if updateHFileNumber < 1 {
					resStatus = i18n.StatusOtherEdit
					tx.Rollback()
					return
				}
			}
		}
	}

	// Prepare update Training Record body Row
	updateRowSql := `update trainingrecord_b set rownumber=$1,studentid=$2,positionname=$3,name=$4,starttime=$5,
	endtime=$6,classhour=$7,description=$8,examres=$9,examscore=$10,	
	modifytime=current_timestamp,modifierid=$11,ts=current_timestamp,dr=$12 
	where id=$13 and ts=$14 and status=0 and dr=0`
	updateRowStmt, err := tx.Prepare(updateRowSql)
	if err != nil {
		zap.L().Error("TrainingRecord.Edit tx.Prepare(updateRowSql) failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer updateRowStmt.Close()
	// Prepare Add Training Record body row
	addRowSql := `insert into trainingrecord_b(hid,rownumber,studentid,positionname,name,
		starttime,endtime,classhour,description,examres,
		examscore,creatorid)
		values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) 
		returning id`
	addRowStmt, err := tx.Prepare(addRowSql)
	if err != nil {
		zap.L().Error("TrainingRecord.Edit tx.Prepare(addRowSql) failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer addRowStmt.Close()
	// Prepare Update row attachment
	updateFileSql := `update trainingrecord_file set modifytime=current_timestamp,modifierid=$1,dr=$2,ts=current_timestamp 
		where id=$3 and dr=0 and ts=$4`
	updateFileStmt, err := tx.Prepare(updateFileSql)
	if err != nil {
		zap.L().Error("TrainingRecord.Edit tx.Prepare(updateFileSql) failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer updateFileStmt.Close()
	// Prepare add row attachment
	addFileSql := `insert into trainingrecord_file(billbid,billhid,fileid,creatorid) 
		values($1,$2,$3,$4) returning id`
	addFileStmt, err := tx.Prepare(addFileSql)
	if err != nil {
		zap.L().Error("TrainingRecord.Edit tx.Prepare(addFileSql) failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer addFileStmt.Close()
	// Update data row by row
	for _, row := range tr.Body {
		// Check row status, any value than 0 is not permitted
		if row.Status != 0 {
			resStatus = i18n.StatusVoucherNoFree
			tx.Rollback()
			return
		}
		if row.BID == 0 { // If the BID value is 0, it means it is a newly row
			addRowErr := addRowStmt.QueryRow(tr.HID, row.RowNumber, row.Student.ID, row.PositionName, row.DeptName,
				row.StartTime, row.EndTime, row.ClassHour, row.Description, row.ExamRes,
				row.ExamScore, tr.Modifier.ID).Scan(&row.BID)
			if addRowErr != nil {
				zap.L().Error("TrainingRecord.Edit addRowStmt.QueryRow() failed:", zap.Error(addRowErr))
				resStatus = i18n.StatusInternalError
				tx.Rollback()
				return resStatus, addRowErr
			}

			// Insert row attachments in the trainingrecord_file table
			if len(row.BFiles) > 0 {
				for _, file := range row.BFiles {
					addFileErr := addFileStmt.QueryRow(row.BID, tr.HID, file.File.ID, tr.Creator.ID).Scan(&file.ID)
					if addFileErr != nil {
						resStatus = i18n.StatusInternalError
						zap.L().Error("TrainingRecord.Edit new row addFileStmt.QueryRow failed:", zap.Error(err))
						tx.Rollback()
						return resStatus, addFileErr
					}
				}
			}

		} else { // If the BID is not 0, it means it is a row that needs to be modified
			updateRowRes, updateRowErr := updateRowStmt.Exec(row.RowNumber, row.Student.ID, row.PositionName, row.DeptName, row.StartTime,
				row.EndTime, row.ClassHour, row.Description, row.ExamRes, row.ExamScore,
				tr.Modifier.ID, row.Dr,
				row.BID, row.Ts)
			if updateRowErr != nil {
				zap.L().Error("TrainingRecord.Edit updateRowStmt.Exec() failed:", zap.Error(updateRowErr))
				resStatus = i18n.StatusInternalError
				tx.Rollback()
				return
			}
			// Check the number of rows affected by SQL Statement
			updateRowNumber, errUpdateEffect := updateRowRes.RowsAffected()
			if errUpdateEffect != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("TrainingRecord.Edit updateRowRes.RowsAffected failed:", zap.Error(errUpdateEffect))
				tx.Rollback()
				return resStatus, errUpdateEffect
			}
			if updateRowNumber < 1 {
				resStatus = i18n.StatusOtherEdit
				tx.Rollback()
				return
			}

			// Process attachments of the row
			if len(row.BFiles) > 0 {
				for _, file := range row.BFiles {
					if file.ID == 0 { // if the ID value is 0, it means it is a newly file
						addFileErr := addFileStmt.QueryRow(row.BID, tr.HID, file.File.ID, tr.Modifier.ID).Scan(&file.ID)
						if addFileErr != nil {
							resStatus = i18n.StatusInternalError
							zap.L().Error("TrainingRecord.Edit old row addFileStmt.QueryRow failed:", zap.Error(addFileErr))
							tx.Rollback()
							return resStatus, addFileErr
						}
					} else { // If the ID value is not 0, it means it is a file that needs to be modified
						updateFileRes, updateFileErr := updateFileStmt.Exec(tr.Modifier.ID, file.Dr, file.ID, file.Ts)
						if updateFileErr != nil {
							resStatus = i18n.StatusInternalError
							zap.L().Error("TrainingRecord.Edit old row updateFileRes.Exec() failed:", zap.Error(updateFileErr))
							tx.Rollback()
							return resStatus, updateFileErr
						}
						updateFileNumber, updateFileEffectErr := updateFileRes.RowsAffected()
						if updateFileEffectErr != nil {
							resStatus = i18n.StatusInternalError
							zap.L().Error("TrainingRecord.Edit old row updateFileRes.RowsAffected() failed:", zap.Error(updateFileEffectErr))
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
		}
	}
	return
}

// Deleete Training Record
func (tr *TrainingRecord) Delete(operatorID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get Training Record by HID
	resStatus, err = tr.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Chcek the data status
	if tr.Status != 0 { // Only value is zero allow delete
		resStatus = i18n.StatusVoucherNoFree
		return
	}
	// Check if the creator and Operator are same person
	if tr.Creator.ID != operatorID {
		resStatus = i18n.StatusVoucherOnlyCreateEdit
		return
	}

	// Begin a database transction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Delete db.Begin failed:", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Modify the Training Record header's deletion flag
	delHeadSql := `update trainingrecord_h set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	delHeadRes, err := tx.Exec(delHeadSql, operatorID, tr.HID, tr.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Delete tx.Exec(delHeadSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by the SQL statement
	delHeadNumber, err := delHeadRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Delete delHeadRes.RowsAffected() failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	if delHeadNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Prepare Modify the Training Record header's attachments deleteion flag
	delHFileSql := `update trainingrecord_file set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and billhid=$3 and ts=$4 and billbid=0`
	delHFileStmt, err := tx.Prepare(delHFileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Delete tx.Prepare(delFileSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delHFileStmt.Close()

	if len(tr.HFiles) > 0 {
		// Wirte data into table row by row
		for _, hFile := range tr.HFiles {
			delHFileRes, delHFileErr := delHFileStmt.Exec(operatorID, hFile.ID, tr.HID, hFile.Ts)
			if delHFileErr != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("TrainingRecord.Delete delHFileStmt.Exec() failed:", zap.Error(delHFileErr))
				tx.Rollback()
				return resStatus, delHFileErr
			}
			// Check the number of rows affected by SQL statement
			delHFileNumber, delHFileEffectedErr := delHFileRes.RowsAffected()
			if delHFileEffectedErr != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("TrainingRecord.Delete delHFileRes.RowsAffected() failed:", zap.Error(delHFileEffectedErr))
				tx.Rollback()
				return resStatus, delHFileEffectedErr
			}
			if delHFileNumber < 1 {
				resStatus = i18n.StatusOtherEdit
				tx.Rollback()
				return
			}
		}
	}

	// Prepare Modify the Training Record body Row's deletion flag
	delRowSql := `update trainingrecord_b set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	// Prepare modify the Training Record body row Attachments's deletion flag
	delFileSql := `update trainingrecord_file set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and billbid=$3 and ts=$4`
	delRowStmt, err := tx.Prepare(delRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Delete tx.Prepare(delRowSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delRowStmt.Close()
	delFileStmt, err := tx.Prepare(delFileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Delete tx.Prepare(delFileSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delFileStmt.Close()
	// Process data row by row
	for _, row := range tr.Body {
		// Check the row status, any value than 0 is not permitted
		if row.Status != 0 {
			resStatus = i18n.StatusVoucherNoFree
			tx.Rollback()
			return
		}
		// Modify row data
		delRowRes, errDelRow := delRowStmt.Exec(operatorID, row.BID, row.Ts)
		if errDelRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("TrainingRecord.Delete delRowStmt.Exec() failed:", zap.Error(errDelRow))
			tx.Rollback()
			return resStatus, errDelRow
		}
		// Chcek the number of rows affected by SQL statement
		delRowNumber, errDelRow := delRowRes.RowsAffected()
		if errDelRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("TrainingRecord.Delete delRowRes.RowsAffected() failed:", zap.Error(errDelRow))
			tx.Rollback()
			return resStatus, errDelRow
		}
		if delRowNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}

		if len(row.BFiles) > 0 {
			// Modify attachments deletion flag item by item
			for _, file := range row.BFiles {
				delFileRes, delFileErr := delFileStmt.Exec(operatorID, file.ID, row.BID, file.Ts)
				if delFileErr != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("TrainingRecord.Delete delFileStmt.Exec() failed:", zap.Error(delFileErr))
					tx.Rollback()
					return resStatus, delFileErr
				}
				// Check the number of rows affected by SQL statement
				delFileNumber, delFileErr := delFileRes.RowsAffected()
				if delFileErr != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("TrainingRecord.Delete delFileRes.RowsAffected() failed:", zap.Error(delFileErr))
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
	}
	return
}

// Confirm Training Record
func (tr *TrainingRecord) Confirm(operatorID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get Training Record details
	resStatus, err = tr.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check the Training Record Status
	if tr.Status != 0 { // Only value is zero allowed
		resStatus = i18n.StatusVoucherNoFree
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Confirm db.Begin failed:", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Update the Training Record header in the trainingrecord_h table
	confirmHeadSql := `update trainingrecord_h set status=1,confirmtime=current_timestamp,confirmerid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and status=0 and ts=$3`
	headRes, err := tx.Exec(confirmHeadSql, operatorID, tr.HID, tr.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Confirm tx.Exec(confirmHeadSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	confirmHeadNumber, err := headRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Confirm headRes.RowsAffected failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	if confirmHeadNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}

	// Prepare modify Training Record body rows in the trainingrecord_b table
	confirmRowSql := `update trainingrecord_b set status=1,confirmtime=current_timestamp,confirmerid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and status=0 and ts=$3`
	rowStmt, err := tx.Prepare(confirmRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.Confirm tx.Prepare(confirmRowSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	defer rowStmt.Close()

	// modify data row by row
	for _, row := range tr.Body {
		// Check the row status
		if row.Status != 0 { // Only status value is zero allowed
			resStatus = i18n.StatusVoucherNoFree
			tx.Rollback()
			return
		}
		confirmRowRes, errConfirmRow := rowStmt.Exec(operatorID, row.BID, row.Ts)
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("TrainingRecord.Confirm rowStmt.Exec failed:", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}

		confirmRowNumber, errConfirmRow := confirmRowRes.RowsAffected()
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("TrainingRecord.Confirm confirmRowRes.RowsAffected failed:", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}
		if confirmRowNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}
	}

	return
}

// UnConfirm Training Record
func (tr *TrainingRecord) UnConfirm(operatorID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get the Training Record details
	resStatus, err = tr.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check the data status
	if tr.Status != 1 { // Must be 1
		resStatus = i18n.StatusVoucherNoConfirm
		return
	}
	// Check the confirmer and operator are not same person
	if tr.Confirmer.ID != operatorID {
		resStatus = i18n.StatusVoucherCancelConfirmSelf
		return
	}
	// Check body rows status
	var noConfirmRowNumber int32
	for _, row := range tr.Body {
		if row.Status > 1 {
			noConfirmRowNumber++
		}
	}
	if noConfirmRowNumber > 0 {
		resStatus = i18n.StatusTRBodyNoConfirm
		return
	}

	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.UnConfirm db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Update the header record in the trainingrecord_h table
	confirmHeadSql := `update trainingrecord_h set status=0,confirmerid=0,ts=current_timestamp 
	where id=$1 and dr=0 and status=1 and ts=$2`
	headRes, err := tx.Exec(confirmHeadSql, tr.HID, tr.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.TrainingRecord.UnConfirm tx.Exec(confirmHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	confirmHeadNumber, err := headRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.UnConfirm headRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if confirmHeadNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Prepare update the body data
	confirmRowSql := `update trainingrecord_b set status=0,confirmerid=0,ts=current_timestamp 
	where id=$1 and dr=0 and status=1 and ts=$2`
	rowStmt, err := tx.Prepare(confirmRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainingRecord.UnConfirm tx.Prepare(confirmRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer rowStmt.Close()
	// Update the body  data row by row
	for _, row := range tr.Body {
		// Check the row status
		if row.Status != 1 { // Must be 1
			resStatus = i18n.StatusVoucherNoConfirm
			tx.Rollback()
			return
		}
		confirmRowRes, errConfirmRow := rowStmt.Exec(row.BID, row.Ts)
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("TrainingRecord.UnConfirm rowStmt.Exec failed", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}

		confirmRowNumber, errConfirmRow := confirmRowRes.RowsAffected()
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("TrainingRecord.UnConfirm confirmRowRes.RowsAffected failed", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}
		if confirmRowNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}
	}
	return
}

// Get Taught Lesson Report
func GetTaughtLessonsReport(queryString string) (tlrs []TaughtLessonsReport, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	tlrs = make([]TaughtLessonsReport, 0)
	var build strings.Builder
	// Concatenate the SQL strings for check
	build.WriteString(`select count(h.id) as hid 
	from trainingrecord_h as h
	left join department as dept on h.deptid=dept.id
	left join sysuser as lecturer on h.lecturerid=lecturer.id
	left join tc on h.tcid=tc.id
	left join sysuser as creator on h.creatorid=creator.id
	where (h.dr=0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	checkSql := build.String()
	// Check
	var rowNumber int32
	err = db.QueryRow(checkSql).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetTaughtLessonsReport db.QueryRow(checkSql) failed", zap.Error(err))
		return
	}

	if rowNumber == 0 {
		resStatus = i18n.StatusResNoData
		return
	}
	if rowNumber > setting.Conf.PqConfig.MaxRecord {
		resStatus = i18n.StatusOverRecord
		return
	}
	build.Reset()

	// Concatenate the SQL string for data retrieval
	build.WriteString(`select h.id as hid,
	h.billnumber as billnumber,
	h.billdate as billdate,
	h.deptid as deptid,
	coalesce(dept.code,'') as code,
	coalesce(dept.name,'') as name,
	h.description as description,
	h.lecturerid as lecturerid,
	coalesce(lecturer.code,'') as lecturercode,
	coalesce(lecturer.name,'') as lecturername,
	h.trainingdate as trainingdate,
	h.tcid as tcid,
	coalesce(tc.code,'') as tccode,
	coalesce(tc.name,'') as tcname,
	h.starttime as starttime,
	h.endtime as endtime,
	h.classhour as classhour,
	h.isexam as isexam,
	(select count(b.id) as studentnumber from trainingrecord_b as b where b.dr=0 and b.hid=h.id),
	(select count(b.id) as qualifiednumber from trainingrecord_b as b where b.dr=0 and b.hid=h.id and b.examres=1),
	(select count(b.id) as disqualificationnumber from trainingrecord_b as b where b.dr=0 and b.hid=h.id and b.examres=0),
	h.status as status,
	h.creatorid as creatorid,
	coalesce(creator.code,'') as creatorCode,
	coalesce(creator.name,'') as creatorName
	from trainingrecord_h as h
	left join department as dept on h.deptid=dept.id
	left join sysuser as lecturer on h.lecturerid=lecturer.id
	left join tc on h.tcid=tc.id
	left join sysuser as creator on h.creatorid=creator.id
	where (h.dr=0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	repSql := build.String()
	// Get Report from database
	glRep, err := db.Query(repSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetTaughtLessonsReport db.Query failed", zap.Error(err))
		return
	}
	defer glRep.Close()

	// Extract data row by row
	for glRep.Next() {
		var glr TaughtLessonsReport
		err = glRep.Scan(&glr.HID, &glr.BillNumber, &glr.BillDate, &glr.ID, &glr.DeptCode,
			&glr.DeptName, &glr.Description, &glr.LecturerID, &glr.LecturerCode, &glr.LecturerName,
			&glr.TrainingDate, &glr.TCID, &glr.TCCode, &glr.TCName, &glr.StartTime,
			&glr.EndTime, &glr.ClassHour, &glr.IsExam, &glr.StudentNumber, &glr.QualifiedNumber,
			&glr.DisqualificationNumber, &glr.Status, &glr.CreatorID, &glr.CreatorCode, &glr.CreatorName)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetTaughtLessonsReport glRep.Next() glRep.Scan failed", zap.Error(err))
			return
		}
		tlrs = append(tlrs, glr)
	}

	return
}

// Get Recived Training Report
func GetRecivedTrainingReport(queryString string) (rtrs []RecivedTrainingReport, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	rtrs = make([]RecivedTrainingReport, 0)
	var build strings.Builder
	// Concatenate the SQL string for check
	build.WriteString(`select count(b.id) as rowcount 
	from trainingrecord_b as b
	left join trainingrecord_h as h on h.id = b.hid
	left join sysuser as student on b.studentid=student.id
	left join department as dept on h.deptid=dept.id
	left join sysuser as lecturer on h.lecturerid=lecturer.id
	left join tc on h.tcid=tc.id
	left join sysuser as creator on h.creatorid=creator.id
	where (h.dr=0 and b.dr=0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	checkSql := build.String()
	// Check
	var rowNumber int32
	err = db.QueryRow(checkSql).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetTaughtLessonsReport db.QueryRow(checkSql) failed", zap.Error(err))
		return
	}

	if rowNumber == 0 {
		resStatus = i18n.StatusResNoData
		return
	}
	if rowNumber > setting.Conf.PqConfig.MaxRecord {
		resStatus = i18n.StatusOverRecord
		return
	}
	build.Reset()

	// Concatenate the SQL string for data retrieval
	build.WriteString(`select h.id as hid,
	b.id as bid,
	h.billnumber as billnumber,
	h.billdate as billdate,
	h.deptid as deptid,
	coalesce(dept.code,'') as code,
	coalesce(dept.name,'') as name,
	h.lecturerid as lecturerid,
	coalesce(lecturer.code,'') as lecturercode,
	coalesce(lecturer.name,'') as lecturername,
	h.tcid as tcid,
	coalesce(tc.code,'') as tccode,
	coalesce(tc.name,'') as tcname,
	h.starttime as starttime,
	h.endtime as endtime,
	h.classhour as tcclasshour,
	h.isexam as isexam,
	h.status as hstatus,
	h.description as hdescription,
	b.studentid as studentid,
	coalesce(student.code,'') as studentcode,
	coalesce(student.name,'') as studentname,
	b.positionname as studentpositionname,
	b.deptname as studentdeptname,
	b.starttime as signstartime,
	b.endtime as signendtime,
	b.classhour as bclasshour,
	b.description as bdscription,
	b.examres as examres,
	b.examscore as examscore,
	b.status as bstatus,
	h.creatorid as creatorid,
	coalesce(creator.code,'') as creatorCode,
	coalesce(creator.name,'') as creatorName
	from trainingrecord_b as b
	left join trainingrecord_h as h on h.id = b.hid
	left join sysuser as student on b.studentid=student.id
	left join department as dept on h.deptid=dept.id
	left join sysuser as lecturer on h.lecturerid=lecturer.id
	left join tc on h.tcid=tc.id
	left join sysuser as creator on h.creatorid=creator.id
	where (h.dr=0 and b.dr=0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	repSql := build.String()
	// Get Report data from database
	rtRep, err := db.Query(repSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetTaughtLessonsReport db.Query failed", zap.Error(err))
		return
	}
	defer rtRep.Close()

	// Extract data row by row
	for rtRep.Next() {
		var rtr RecivedTrainingReport
		err = rtRep.Scan(&rtr.HID, &rtr.BID, &rtr.BillNumber, &rtr.BillDate, &rtr.ID,
			&rtr.DeptCode, &rtr.DeptName, &rtr.LecturerID, &rtr.LecturerCode, &rtr.LecturerName,
			&rtr.TCID, &rtr.TCCode, &rtr.TCName, &rtr.StartTime, &rtr.EndTime,
			&rtr.TCClassHour, &rtr.IsExam, &rtr.HStatus, &rtr.HDescription, &rtr.StudentID,
			&rtr.StudentCode, &rtr.StudentName, &rtr.StudentPositionName, &rtr.StudentDeptName, &rtr.SignStartTime,
			&rtr.SignEndTime, &rtr.BClassHour, &rtr.BDescription, &rtr.ExamRes, &rtr.ExamScore,
			&rtr.BStatus, &rtr.CreatorID, &rtr.CreatorCode, &rtr.CreatorName)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetTaughtLessonsReport glRep.Next() glRep.Scan failed", zap.Error(err))
			return
		}

		rtrs = append(rtrs, rtr)
	}

	return
}
