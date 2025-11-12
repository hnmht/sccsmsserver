package pg

import (
	"sccsmsserver/i18n"
	"sccsmsserver/setting"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Personal Protective Equipment Issuance Form Struct
type PPEIssuanceForm struct {
	HID         int32                `db:"id" json:"id"`
	BillNumber  string               `db:"billnumber" json:"billNumber"`
	BillDate    time.Time            `db:"billdate" json:"billDate"`
	Department  SimpDept             `db:"deptid" json:"department"`
	Description string               `db:"description" json:"description"`
	Period      string               `db:"period" json:"period"`
	StartDate   time.Time            `db:"startdate" json:"startDate"`
	EndDate     time.Time            `db:"enddate" json:"endDate"`
	HFiles      []VoucherFile        `json:"hFiles"`
	Body        []PPEIssuanceFormRow `json:"body"`
	SourceType  string               `db:"sourcetype" json:"sourceType"` // DA: Direct Add WG: Wizard Generation
	Status      int16                `db:"status" json:"status"`         // 0 Free 1 Confirmed 2 Executing 3 Completed
	CreateDate  time.Time            `db:"createtime" json:"createDate"`
	Creator     Person               `db:"creatorid" json:"creator"`
	ConfirmDate time.Time            `db:"confirmtime" json:"confirmDate"`
	Confirmer   Person               `db:"confirmerid" json:"confirmer"`
	ModifyDate  time.Time            `db:"modifytime" json:"modifyDate"`
	Modifier    Person               `db:"modifierid" json:"modifier"`
	Ts          time.Time            `db:"ts" json:"ts"`
	Dr          int16                `db:"dr" json:"dr"`
}

// Personal Protective Equipment Issuance Form Row struct
type PPEIssuanceFormRow struct {
	BID          int32         `db:"id" json:"id"`
	HID          int32         `db:"hid" json:"hid"`
	RowNumber    int32         `db:"rownumber" json:"rowNumber"`
	Recipient    Person        `db:"recipientid" json:"recipient"`
	PositionName string        `db:"positionname" json:"positionName"`
	DeptName     string        `db:"deptname" json:"deptName"`
	PPECode      string        `json:"ppeCode"`
	PPE          PPE           `db:"ppeid" json:"ppe"`
	PPEModel     string        `json:"ppeModel"`
	PPEUnit      string        `json:"ppeUnit"`
	Quantity     float64       `db:"quantity" json:"quantity"`
	Description  string        `db:"description" json:"description"`
	Status       int16         `db:"status" json:"status"` // 0 Free 1 Confirmed 2 Executing 3 Completed 4 none
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

// PPE Issuance Form Wizard Params struct
type PPEIssuanceFormWizardParams struct {
	BillDate       time.Time `json:"billDate"`
	Department     SimpDept  `json:"department"`
	Description    string    `json:"description"`
	Period         string    `json:"period"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	Creator        Person    `json:"creator"`
	GenerationType int16     `json:"generationType"` //  0: Combined Generation 1: Separate Generation
}

// PPE Issuance Form Wizard struct
type PPEIssuanceFormWizard struct {
	Params         PPEIssuanceFormWizardParams `json:"params"`
	Recipients     []Person                    `json:"recipients"`
	VoucherNumbers []string                    `json:"vouchernumbers"`
}

// PPE Issuance Form Report struct
type PPEIssuanceFormReport struct {
	HID                   int32     `json:"hid"`
	BID                   int32     `json:"bid"`
	RowNumber             int32     `json:"rowNumber"`
	RecipientID           int32     `json:"recipientID"`
	RecipientCode         string    `json:"recipientCode"`
	RecipientName         string    `json:"recipientName"`
	RecipientPositionName string    `json:"recipientPoisitionName"`
	RecipientDeptName     string    `json:"recipientDeptName"`
	PPEID                 int32     `json:"ppeID"`
	PPECode               string    `json:"ppeCode"`
	PPEName               string    `json:"ppeName"`
	PPEModel              string    `json:"ppeModel"`
	PPEUnit               string    `json:"ppeUnit"`
	Quantity              float64   `json:"quantity"`
	BDescription          string    `json:"bDescription"`
	BStatus               int16     `json:"bStatus"`
	Billnumber            string    `json:"billNumber"`
	BillDate              time.Time `json:"billDate"`
	IssuingDeptID         int32     `json:"issuingDeptID"`
	IssuingDeptCode       string    `json:"issuingDeptCode"`
	IssuingDeptName       string    `json:"issuingDeptName"`
	HDescription          string    `json:"hDescription"`
	Period                string    `json:"period"`
	StartDate             time.Time `json:"startDate"`
	EndDate               time.Time `json:"endDate"`
	SourceType            string    `json:"sourceType"`
	Hstatus               int16     `json:"hStatus"`
	CreateUserID          int32     `json:"creatorID"`
	CreateUserCode        string    `json:"creatorCode"`
	CreateUserName        string    `json:"creatorName"`
}

// Generate a PPE Issuance Form via Wizard
func (pifw *PPEIssuanceFormWizard) Generate() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	if pifw.Params.GenerationType == 0 {
		resStatus, err = pifw.CombinedGeneration()
	} else {
		resStatus, err = pifw.SeparateGeneration()
	}
	return
}

// Generate a combined PPE Issuance Form via Wizard
func (pifw *PPEIssuanceFormWizard) CombinedGeneration() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	pifw.VoucherNumbers = make([]string, 0)
	var pif PPEIssuanceForm
	// Fill in the header items of the issuance form
	pif.BillDate = pifw.Params.BillDate
	pif.Department = pifw.Params.Department
	pif.Description = pifw.Params.Description
	pif.Period = pifw.Params.Period
	pif.StartDate = pifw.Params.StartDate
	pif.EndDate = pifw.Params.EndDate
	pif.Creator = pifw.Params.Creator
	pif.SourceType = "WG"
	// Get PPE Quota from database
	ppeQuotaSQL := `select id from ppequotas_h where positionid=$1 and period=$2 and status=1 and dr=0`
	stmt, err := db.Prepare(ppeQuotaSQL)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceFormWizard.CombinedGeneration db.Prepare(ppeQuotaSQL) failed", zap.Error(err))
		return
	}
	defer stmt.Close()
	var rowNumber int32
	// Fill in the body items of the issuance form
	for _, person := range pifw.Recipients {
		var pq PPEQuota
		err = stmt.QueryRow(person.PositionID, pifw.Params.Period).Scan(&pq.HID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEIssuanceFormWizard.CombinedGeneration stmt.QueryRow failed", zap.Error(err))
			return
		}
		// Get PPE Quota details
		if pq.HID > 0 {
			resStatus, err = pq.GetDetailByHID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Generate PPE Issuance Form body rows
		for _, pqRow := range pq.Body {
			rowNumber = rowNumber + 10
			var pifr PPEIssuanceFormRow
			pifr.RowNumber = rowNumber
			pifr.Recipient = person
			pifr.PositionName = person.PositionName
			pifr.DeptName = person.DeptName
			pifr.PPE = pqRow.PPE
			pifr.Quantity = pqRow.Quantity
			pifr.Description = pqRow.Description
			pifr.Status = 0
			pifr.Creator = pifw.Params.Creator
			pif.Body = append(pif.Body, pifr)
		}
	}

	// Add PPE Issuance Form
	resStatus, err = pif.Add()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Write Bill Number back
	pifw.VoucherNumbers = append(pifw.VoucherNumbers, pif.BillNumber)

	return
}

// Generate separate PPE Issuance Forms via Wizard
func (pifw *PPEIssuanceFormWizard) SeparateGeneration() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	pifw.VoucherNumbers = make([]string, 0)
	// Get PPE Quota from database
	ppeQuotaSQL := `select id from ppequotas_h where positionid=$1 and period=$2 and status=1 and dr=0`
	stmt, err := db.Prepare(ppeQuotaSQL)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceFormWizard.SeparateGeneration db.Prepare(ppeQuotaSQL) failed", zap.Error(err))
		return
	}
	defer stmt.Close()
	// Fill in the body items of the issuance form
	for _, person := range pifw.Recipients {
		var pif PPEIssuanceForm
		// Fill in the header items of the issuance form
		pif.BillDate = pifw.Params.BillDate
		pif.Department = pifw.Params.Department
		pif.Description = person.Name + "_" + pifw.Params.Description
		pif.Period = pifw.Params.Period
		pif.StartDate = pifw.Params.StartDate
		pif.EndDate = pifw.Params.EndDate
		pif.Creator = pifw.Params.Creator
		pif.SourceType = "WG"
		var rowNumber int32
		var pq PPEQuota
		err = stmt.QueryRow(person.PositionID, pifw.Params.Period).Scan(&pq.HID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEIssuanceFormWizard.SeparateGeneration stmt.QueryRow failed", zap.Error(err))
			return
		}
		// Get PPE Quota details
		if pq.HID > 0 {
			resStatus, err = pq.GetDetailByHID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Generate PPE Issuance Form body rows
		for _, pqRow := range pq.Body {
			rowNumber = rowNumber + 10
			var pifr PPEIssuanceFormRow
			pifr.RowNumber = rowNumber
			pifr.Recipient = person
			pifr.PositionName = person.PositionName
			pifr.DeptName = person.DeptName
			pifr.PPE = pqRow.PPE
			pifr.Quantity = pqRow.Quantity
			pifr.Description = pqRow.Description
			pifr.Status = 0
			pifr.Creator = pifw.Params.Creator
			pif.Body = append(pif.Body, pifr)
		}
		// Add PPE Issuance Form
		resStatus, err = pif.Add()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		// Write Bill Number back
		pifw.VoucherNumbers = append(pifw.VoucherNumbers, pif.BillNumber)
	}
	return
}

// Add PPE Issuance Form
func (pif *PPEIssuanceForm) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check the number of rows in the body, it cannot be zero
	if len(pif.Body) == 0 {
		resStatus = i18n.StatusVoucherNoBody
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Add db.Begin failed:", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Get the latest Serial Number
	billNo, resStatus, err := GetLatestSerialNo(tx, "PIF", pif.BillDate.Format("060102"))
	if resStatus != i18n.StatusOK || err != nil {
		tx.Rollback()
		return
	}
	pif.BillNumber = billNo
	// Write the header content to the ppeissuanceform_h table
	headSql := `insert into ppeissuanceform_h(billnumber,billdate,deptid,description,period,
		startdate,enddate,sourcetype,status,creatorid) 
		values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) 
		returning id`
	err = tx.QueryRow(headSql, pif.BillNumber, pif.BillDate, pif.Department.ID, pif.Description, pif.Period,
		pif.StartDate, pif.EndDate, pif.SourceType, pif.Status, pif.Creator.ID).Scan(&pif.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Add tx.QeuryRow(headSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	// Prepare insert the header attachment record into the ppeissuanceform_file table
	headFileSql := `insert into ppeissuanceform_file(billhid,fileid,creatorid) 
	values($1,$2,$3) returning id`
	headFileStmt, err := tx.Prepare(headFileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Add tx.Prepare(headFileSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	defer headFileStmt.Close()
	// Insert the header attachments item by item
	for _, hFile := range pif.HFiles {
		err = headFileStmt.QueryRow(pif.HID, hFile.File.ID, pif.Creator.ID).Scan(&hFile.ID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEIssuanceForm.Add headFileStmt.QueryRow failed:", zap.Error(err))
			tx.Rollback()
			return
		}
	}
	// Prepare insert content of body into the ppeissuanceform_b table
	bodySql := `insert into ppeissuanceform_b(hid,rownumber,recipientid,positionname,deptname,
		ppeid,quantity,description,status,creatorid)
		values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) 
		returning id`
	bodyStmt, err := tx.Prepare(bodySql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Add tx.Prepare(bodySql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	defer bodyStmt.Close()
	// Perpare Insert  attachments of the row in the body into ppeissuanceform_file table
	fileSql := `insert into ppeissuanceform_file(billbid,billhid,fileid,creatorid) 
	values($1,$2,$3,$4) returning id`
	fileStmt, err := tx.Prepare(fileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("TrainRecordAdd tx.Prepare(fileSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	defer fileStmt.Close()
	// Write the body data to the database row by row
	for _, row := range pif.Body {
		// Row content
		err = bodyStmt.QueryRow(pif.HID, row.RowNumber, row.Recipient.ID, row.PositionName, row.DeptName,
			row.PPE.ID, row.Quantity, row.Description, row.Status, pif.Creator.ID).Scan(&row.BID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEIssuanceForm.Add bodyStmt.QueryRow failed:", zap.Error(err))
			tx.Rollback()
			return
		}
		// Row Attachments
		if len(row.BFiles) > 0 {
			for _, file := range row.BFiles {
				err = fileStmt.QueryRow(row.BID, pif.HID, file.File.ID, pif.Creator.ID).Scan(&file.ID)
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

// Get PPE Issuance Form Details by HID
func (pif *PPEIssuanceForm) GetDetailByHID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the data is deleted
	var rowNumber int32
	checkSql := `select count(id) as rownumber from ppeissuanceform_h where id=$1 and dr=0`
	err = db.QueryRow(checkSql, pif.HID).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.GetDetailByHID db.QueryRow(checkSql) failed:", zap.Error(err))
		return
	}
	if rowNumber < 1 {
		resStatus = i18n.StatusDataDeleted
		return
	}
	// Fill in the Body Rows
	resStatus, err = pif.FillBody()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	return
}

// Fill the PPE Issuance From body rows
func (pif *PPEIssuanceForm) FillBody() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Prepare to get the body rows from the database
	bodySql := `select id,hid,rownumber,recipientid,positionname,
		deptname,ppeid,quantity,description,status,
		createtime,creatorid,confirmtime,confirmerid,modifytime,
		modifierid,dr,ts from ppeissuanceform_b
		where hid=$1 and dr=0 order by rownumber asc`
	bodyRows, err := db.Query(bodySql, pif.HID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.FillBody db.Query(bodySql) failed:", zap.Error(err))
		return
	}
	defer bodyRows.Close()
	var bodyRowNumber int32
	// Extract data row by row
	for bodyRows.Next() {
		bodyRowNumber++
		var pifr PPEIssuanceFormRow
		err = bodyRows.Scan(&pifr.BID, &pifr.HID, &pifr.RowNumber, &pifr.Recipient.ID, &pifr.PositionName,
			&pifr.DeptName, &pifr.PPE.ID, &pifr.Quantity, &pifr.Description, &pifr.Status,
			&pifr.CreateDate, &pifr.Creator.ID, &pifr.ConfirmDate, &pifr.Confirmer.ID, &pifr.ModifyDate,
			&pifr.Modifier.ID, &pifr.Dr, &pifr.Ts)
		if err != nil {
			zap.L().Error("PPEIssuanceForm.FillBody bodyRows.scan failed:", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		// Get Recipient Info
		if pifr.Recipient.ID > 0 {
			resStatus, err = pifr.Recipient.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get PPE Info
		if pifr.PPE.ID > 0 {
			resStatus, err = pifr.PPE.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		pifr.PPECode = pifr.PPE.Code
		pifr.PPEModel = pifr.PPE.Model
		pifr.PPEUnit = pifr.PPE.Unit

		// Get Creator Info
		if pifr.Creator.ID > 0 {
			resStatus, err = pifr.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// GEt Confirmer Info
		if pifr.Confirmer.ID > 0 {
			resStatus, err = pifr.Confirmer.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier Info
		if pifr.Modifier.ID > 0 {
			resStatus, err = pifr.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Row Attachments
		pifr.BFiles, resStatus, err = GetPPEIFRowFiles(pifr.BID)
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		pif.Body = append(pif.Body, pifr)
	}
	return
}

// Get PPE Issuance Form body Row Files
func GetPPEIFRowFiles(bid int32) (voucherFiles []VoucherFile, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	voucherFiles = make([]VoucherFile, 0)
	attachSql := `select id,billbid,billhid,fileid,createtime,
	creatorid,modifytime,modifierid,dr,ts 
	from ppeissuanceform_file where billbid=$1 and dr=0`
	// Fill in the attachments
	fileRows, err := db.Query(attachSql, bid)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetPPEIFRowFiles db.query(attachsql) failed:", zap.Error(err))
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
			zap.L().Error("GetPPEIFRowFiles fileRows.Scan failed:", zap.Error(fileErr))
			return
		}
		// Get File Info
		if f.File.ID > 0 {
			resStatus, err = f.File.GetFileInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Creator Info
		if f.Creator.ID > 0 {
			resStatus, err = f.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier Info
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

// Get PPE Issuance Form Header Files
func GetPPEIFHeaderFiles(hid int32) (voucherFiles []VoucherFile, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	voucherFiles = make([]VoucherFile, 0)
	// Prepare to get the attachments
	attachSql := `select id,billbid,billhid,fileid,createtime,
	creatorid,modifytime,modifierid,dr,ts 
	from ppeissuanceform_file where billhid=$1 and dr=0 and billbid=0`
	fileRows, err := db.Query(attachSql, hid)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetPPEIFHeaderFiles db.query(attachsql) failed:", zap.Error(err))
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
			zap.L().Error("GetPPEIFHeaderFiles fileRows.Scan failed:", zap.Error(fileErr))
			return
		}
		// Get file info
		if f.File.ID > 0 {
			resStatus, err = f.File.GetFileInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Creator info
		if f.Creator.ID > 0 {
			resStatus, err = f.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier info
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

// Fill PPE Issuance Form Header Information
func (pif *PPEIssuanceForm) FillHead() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get Department Info
	if pif.Department.ID > 0 {
		resStatus, err = pif.Department.GetSimpDeptInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Creator Info
	if pif.Creator.ID > 0 {
		resStatus, err = pif.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Confirmer Info
	if pif.Confirmer.ID > 0 {
		resStatus, err = pif.Confirmer.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Modifier Info
	if pif.Modifier.ID > 0 {
		resStatus, err = pif.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Header Attachments
	pif.HFiles, resStatus, err = GetPPEIFHeaderFiles(pif.HID)
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	return
}

// Get PPE Issuance Form List
func GetPPEIFList(queryString string) (pifs []PPEIssuanceForm, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	pifs = make([]PPEIssuanceForm, 0)
	var build strings.Builder
	// Concatenate SQL String for checking
	build.WriteString(`select count(h.id) as rownumber
	from ppeissuanceform_h as h
	left join department on h.deptid = department.id
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
		zap.L().Error("GetPPEIFList db.QueryRow(checkSql) failed:", zap.Error(err))
		return
	}
	if rowNumber == 0 {
		resStatus = i18n.StatusResNoData
		return
	}
	if rowNumber > setting.Conf.PqConfig.MaxRecord {
		return
	}
	build.Reset()

	// Concatenate SQL String for getting data
	build.WriteString(`select h.id,h.billnumber,h.billdate,h.deptid,h.description,
	h.period,h.startdate,h.enddate,h.sourcetype,h.status,
	h.createtime,h.creatorid,h.confirmtime,h.confirmerid,h.modifytime,
	h.modifierid,h.dr,h.ts 
	from ppeissuanceform_h as h
	left join department on h.deptid = department.id
	left join sysuser as creator on h.creatorid = creator.id
	left join sysuser as modifier on h.modifierid = modifier.id
	where (h.dr = 0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	headSql := build.String()

	headRows, err := db.Query(headSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetPPEIFList db.Query failed:", zap.Error(err))
		return
	}
	defer headRows.Close()
	// Extract data row by row
	for headRows.Next() {
		var pif PPEIssuanceForm
		err = headRows.Scan(&pif.HID, &pif.BillNumber, &pif.BillDate, &pif.Department.ID, &pif.Description,
			&pif.Period, &pif.StartDate, &pif.EndDate, &pif.SourceType, &pif.Status,
			&pif.CreateDate, &pif.Creator.ID, &pif.ConfirmDate, &pif.Confirmer.ID, &pif.ModifyDate,
			&pif.Modifier.ID, &pif.Dr, &pif.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetPPEIFList headRows.Next failed:", zap.Error(err))
			return
		}
		// Fill in the header items
		resStatus, err = pif.FillHead()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		pifs = append(pifs, pif)
	}

	return
}

// Edit PPE Issuance Form
func (pif *PPEIssuanceForm) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check the number of rows in the body, it cannot be zero
	if len(pif.Body) == 0 {
		resStatus = i18n.StatusVoucherNoBody
		return
	}
	// Check if the modifier is the creator
	if pif.Creator.ID != pif.Modifier.ID {
		resStatus = i18n.StatusVoucherOnlyCreateEdit
		return
	}

	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Edit db.Begin() failed:", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Modify the header content in the ppeissuanceform_h table
	ldHeadSql := `update ppeissuanceform_h set billdate=$1,deptid=$2,description=$3, period=$4,startdate=$5,
	enddate=$6,sourcetype=$7,modifytime=current_timestamp,modifierid=$8,ts=current_timestamp  
	where id=$9 and dr=0 and status=0 and ts=$10`
	ldHeadRes, err := tx.Exec(ldHeadSql, &pif.BillDate, &pif.Department.ID, &pif.Description, &pif.Period, &pif.StartDate,
		&pif.EndDate, &pif.SourceType, &pif.Modifier.ID,
		&pif.HID, &pif.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Edit tx.Exec(tritHeadSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of affected rows
	headUpdateNumber, err := ldHeadRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Edit EditHeadRes.RowsAffected failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	if headUpdateNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Prepare to modify header attachments
	updateHFileSql := `update ppeissuanceform_file set modifytime=current_timestamp,modifierid=$1,dr=$2,ts=current_timestamp
	where id=$3 and dr=0 and ts=$4 and billbid=0`
	updateHFileStmt, err := tx.Prepare(updateHFileSql)
	if err != nil {
		zap.L().Error("PPEIssuanceForm.Edit tx.Prepare(updateHFileSql) failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer updateHFileStmt.Close()
	// Prepare to add header attachments
	addHFileSql := `insert into ppeissuanceform_file(billhid,fileid,creatorid) 
	values($1,$2,$3) returning id`
	addHFileStmt, err := tx.Prepare(addHFileSql)
	if err != nil {
		zap.L().Error("PPEIssuanceForm.Edit tx.Prepare(addHFileSql) failed::", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer addHFileStmt.Close()
	// Write header attachments
	if len(pif.HFiles) > 0 {
		for _, hFile := range pif.HFiles {
			if hFile.ID == 0 { // If ID value is 0, it means a new attachment
				addHFileErr := addHFileStmt.QueryRow(pif.HID, hFile.File.ID, pif.Modifier.ID).Scan(&hFile.ID)
				if addHFileErr != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("PPEIssuanceForm.Edit old row addHFileStmt.QueryRow failed:", zap.Error(addHFileErr))
					tx.Rollback()
					return resStatus, addHFileErr
				}
			} else { // If ID value is not 0, it means modifying an existing attachment
				updateHFileRes, updateHFileErr := updateHFileStmt.Exec(pif.Modifier.ID, hFile.Dr, hFile.ID, hFile.Ts)
				if updateHFileErr != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("PPEIssuanceForm.Edit old row updateHFileStmt.Exec() failed:", zap.Error(updateHFileErr))
					tx.Rollback()
					return resStatus, updateHFileErr
				}
				updateHFileNumber, updateHFileEffectErr := updateHFileRes.RowsAffected()
				if updateHFileEffectErr != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("PPEIssuanceForm.Edit old row updateHFileRes.RowsAffected() failed:", zap.Error(updateHFileEffectErr))
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

	// Prepare to modify body rows
	updateRowSql := `update ppeissuanceform_b set rownumber=$1,recipientid=$2,positionname=$3,deptname=$4,ppeid=$5,
	quantity=$6,description=$7,modifytime=current_timestamp,modifierid=$8,ts=current_timestamp,dr=$9 
	where id=$10 and ts=$11 and status=0 and dr=0`
	updateRowStmt, err := tx.Prepare(updateRowSql)
	if err != nil {
		zap.L().Error("PPEIssuanceForm.Edit tx.Prepare(updateRowSql) failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer updateRowStmt.Close()
	// Prepare to add body rows
	addRowSql := `insert into ppeissuanceform_b(hid,rownumber,recipientid,positionname,deptname,
		ppeid,quantity,description,creatorid)
		values($1,$2,$3,$4,$5,$6,$7,$8,$9) 
		returning id`
	addRowStmt, err := tx.Prepare(addRowSql)
	if err != nil {
		zap.L().Error("PPEIssuanceForm.Edit tx.Prepare(addRowSql) failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer addRowStmt.Close()
	// Prepare to modify file attachments
	updateFileSql := `update ppeissuanceform_file set modifytime=current_timestamp,modifierid=$1,dr=$2,ts=current_timestamp 
		where id=$3 and dr=0 and ts=$4`
	updateFileStmt, err := tx.Prepare(updateFileSql)
	if err != nil {
		zap.L().Error("PPEIssuanceForm.Edit tx.Prepare(updateFileSql) failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer updateFileStmt.Close()
	// Prepare to add file attachments
	addFileSql := `insert into ppeissuanceform_file(billbid,billhid,fileid,creatorid) 
		values($1,$2,$3,$4) returning id`
	addFileStmt, err := tx.Prepare(addFileSql)
	if err != nil {
		zap.L().Error("PPEIssuanceForm.Edit tx.Prepare(addFileSql) failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer addFileStmt.Close()
	// Write body rows one by one
	for _, row := range pif.Body {
		// Check if the row status is free
		if row.Status != 0 {
			resStatus = i18n.StatusVoucherNoFree
			tx.Rollback()
			return
		}
		if row.BID == 0 { // If BID is 0, it means a new row
			addRowErr := addRowStmt.QueryRow(pif.HID, row.RowNumber, row.Recipient.ID, row.PositionName, row.DeptName,
				row.PPE.ID, row.Quantity, row.Description, pif.Modifier.ID).Scan(&row.BID)
			if addRowErr != nil {
				zap.L().Error("PPEIssuanceForm.Edit addRowStmt.QueryRow() failed:", zap.Error(addRowErr))
				resStatus = i18n.StatusInternalError
				tx.Rollback()
				return resStatus, addRowErr
			}

			// Handle attachments
			if len(row.BFiles) > 0 {
				for _, file := range row.BFiles {
					addFileErr := addFileStmt.QueryRow(row.BID, pif.HID, file.File.ID, pif.Creator.ID).Scan(&file.ID)
					if addFileErr != nil {
						resStatus = i18n.StatusInternalError
						zap.L().Error("PPEIssuanceForm.Edit new row addFileStmt.QueryRow failed:", zap.Error(err))
						tx.Rollback()
						return resStatus, addFileErr
					}
				}
			}

		} else { // If BID is not 0, it means modifying an existing row
			// Modify the row content
			updateRowRes, updateRowErr := updateRowStmt.Exec(row.RowNumber, row.Recipient.ID, row.PositionName, row.DeptName, row.PPE.ID,
				row.Quantity, row.Description, pif.Modifier.ID, row.Dr,
				row.BID, row.Ts)
			if updateRowErr != nil {
				zap.L().Error("PPEIssuanceForm.Edit updateRowStmt.Exec() failed:", zap.Error(updateRowErr))
				resStatus = i18n.StatusInternalError
				tx.Rollback()
				return
			}
			// Check the number of affected rows
			updateRowNumber, errUpdateEffect := updateRowRes.RowsAffected()
			if errUpdateEffect != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("PPEIssuanceForm.Edit updateRowRes.RowsAffected failed:", zap.Error(errUpdateEffect))
				tx.Rollback()
				return resStatus, errUpdateEffect
			}
			if updateRowNumber < 1 {
				resStatus = i18n.StatusOtherEdit
				tx.Rollback()
				return
			}

			// Handle attachments
			if len(row.BFiles) > 0 {
				for _, file := range row.BFiles {
					if file.ID == 0 { // If ID is 0, it means a new attachment
						addFileErr := addFileStmt.QueryRow(row.BID, pif.HID, file.File.ID, pif.Modifier.ID).Scan(&file.ID)
						if addFileErr != nil {
							resStatus = i18n.StatusInternalError
							zap.L().Error("PPEIssuanceForm.Edit old row addFileStmt.QueryRow failed:", zap.Error(addFileErr))
							tx.Rollback()
							return resStatus, addFileErr
						}
					} else { // If ID is not 0, it means modifying an existing attachment
						updateFileRes, updateFileErr := updateFileStmt.Exec(pif.Modifier.ID, file.Dr, file.ID, file.Ts)
						if updateFileErr != nil {
							resStatus = i18n.StatusInternalError
							zap.L().Error("PPEIssuanceForm.Edit old row updateFileRes.Exec() failed:", zap.Error(updateFileErr))
							tx.Rollback()
							return resStatus, updateFileErr
						}
						updateFileNumber, updateFileEffectErr := updateFileRes.RowsAffected()
						if updateFileEffectErr != nil {
							resStatus = i18n.StatusInternalError
							zap.L().Error("PPEIssuanceForm.Edit old row updateFileRes.RowsAffected() failed:", zap.Error(updateFileEffectErr))
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

// Delete PPE Issuance Form
func (pif *PPEIssuanceForm) Delete(operatorID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get PPE Issuance Form Details
	resStatus, err = pif.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check the status of the PPE Issuance Form
	if pif.Status != 0 { // If the status is not free, it cannot be deleted
		resStatus = i18n.StatusVoucherNoFree
		return
	}
	// Check if the modifier is the creator
	if pif.Creator.ID != operatorID {
		resStatus = i18n.StatusVoucherOnlyCreateEdit
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Delete db.Begin failed:", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Update delete flag in the ppeissuanceform_h table
	delHeadSql := `update ppeissuanceform_h set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	delHeadRes, err := tx.Exec(delHeadSql, operatorID, pif.HID, pif.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Delete tx.Exec(delHeadSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of affected rows
	delHeadNumber, err := delHeadRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Delete delHeadRes.RowsAffected() failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	if delHeadNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Delete header attachments
	delHFileSql := `update ppeissuanceform_file set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and billhid=$3 and ts=$4 and billbid=0`
	delHFileStmt, err := tx.Prepare(delHFileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Delete tx.Prepare(delFileSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delHFileStmt.Close()

	if len(pif.HFiles) > 0 {
		for _, hFile := range pif.HFiles {
			delHFileRes, delHFileErr := delHFileStmt.Exec(operatorID, hFile.ID, pif.HID, hFile.Ts)
			if delHFileErr != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("PPEIssuanceForm.Delete delHFileStmt.Exec() failed:", zap.Error(delHFileErr))
				tx.Rollback()
				return resStatus, delHFileErr
			}
			// Check the number of affected rows
			delHFileNumber, delHFileEffectedErr := delHFileRes.RowsAffected()
			if delHFileEffectedErr != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("PPEIssuanceForm.Delete delHFileRes.RowsAffected() failed:", zap.Error(delHFileEffectedErr))
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

	// Prepare to delete body rows and attachments
	delRowSql := `update ppeissuanceform_b set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	delRowStmt, err := tx.Prepare(delRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Delete tx.Prepare(delRowSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delRowStmt.Close()
	// Prepare to delete attachments
	delFileSql := `update ppeissuanceform_file set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and billbid=$3 and ts=$4`
	delFileStmt, err := tx.Prepare(delFileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Delete tx.Prepare(delFileSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delFileStmt.Close()
	// Update delete flag for each body row
	for _, row := range pif.Body {
		// Check the status of the body row
		if row.Status != 0 {
			resStatus = i18n.StatusVoucherNoFree
			tx.Rollback()
			return
		}
		delRowRes, errDelRow := delRowStmt.Exec(operatorID, row.BID, row.Ts)
		if errDelRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEIssuanceForm.Delete delRowStmt.Exec() failed:", zap.Error(errDelRow))
			tx.Rollback()
			return resStatus, errDelRow
		}
		// Check the number of affected rows
		delRowNumber, errDelRow := delRowRes.RowsAffected()
		if errDelRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEIssuanceForm.Delete delRowRes.RowsAffected() failed:", zap.Error(errDelRow))
			tx.Rollback()
			return resStatus, errDelRow
		}
		if delRowNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}

		if len(row.BFiles) > 0 {
			for _, file := range row.BFiles {
				delFileRes, delFileErr := delFileStmt.Exec(operatorID, file.ID, row.BID, file.Ts)
				if delFileErr != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("PPEIssuanceForm.Delete delFileStmt.Exec() failed:", zap.Error(delFileErr))
					tx.Rollback()
					return resStatus, delFileErr
				}
				// Check the number of affected rows
				delFileNumber, delFileErr := delFileRes.RowsAffected()
				if delFileErr != nil {
					resStatus = i18n.StatusInternalError
					zap.L().Error("PPEIssuanceForm.Delete delFileRes.RowsAffected() failed:", zap.Error(delFileErr))
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

// Confirm PPE Issuance Form
func (pif *PPEIssuanceForm) Confirm(operatorID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get PPE Issuance Form Details
	resStatus, err = pif.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check the status of the PPE Issuance Form
	if pif.Status != 0 {
		resStatus = i18n.StatusVoucherNoFree
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Confirm db.Begin failed:", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Update header to confirmed status
	confirmHeadSql := `update ppeissuanceform_h set status=1,confirmtime=current_timestamp,confirmerid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and status=0 and ts=$3`
	headRes, err := tx.Exec(confirmHeadSql, operatorID, pif.HID, pif.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Confirm tx.Exec(confirmHeadSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of affected rows
	confirmHeadNumber, err := headRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Confirm headRes.RowsAffected failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	if confirmHeadNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}

	// Prepare to confirm body rows
	confirmRowSql := `update ppeissuanceform_b set status=1,confirmtime=current_timestamp,confirmerid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and status=0 and ts=$3`
	rowStmt, err := tx.Prepare(confirmRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Confirm tx.Prepare(confirmRowSql) failed:", zap.Error(err))
		tx.Rollback()
		return
	}
	defer rowStmt.Close()
	// Write body row confirmations one by one
	for _, row := range pif.Body {
		// Check the status of the body row
		if row.Status != 0 {
			resStatus = i18n.StatusVoucherNoFree
			tx.Rollback()
			return
		}
		confirmRowRes, errConfirmRow := rowStmt.Exec(operatorID, row.BID, row.Ts)
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEIssuanceForm.Confirm rowStmt.Exec failed:", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}

		confirmRowNumber, errConfirmRow := confirmRowRes.RowsAffected()
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEIssuanceForm.Confirm confirmRowRes.RowsAffected failed:", zap.Error(errConfirmRow))
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

// Unconfirm PPE Issuance Form
func (pif *PPEIssuanceForm) Unconfirm(operatorID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get PPE Issuance Form Details
	resStatus, err = pif.GetDetailByHID()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check the status of the PPE Issuance Form
	if pif.Status != 1 { // if the status is not confirmed, it cannot be unconfirmed
		resStatus = i18n.StatusVoucherNoConfirm
		return
	}
	// Check if the operator is the confirmer
	if pif.Confirmer.ID != operatorID {
		resStatus = i18n.StatusVoucherCancelConfirmSelf
		return
	}
	// Check if all body rows are confirmed
	var noConfirmRowNumber int32
	for _, row := range pif.Body {
		if row.Status > 1 {
			noConfirmRowNumber++
		}
	}
	if noConfirmRowNumber > 0 {
		resStatus = i18n.StatusPPEIFBodyNoConfirm
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Unconfirm db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Update header to unconfirmed status
	confirmHeadSql := `update ppeissuanceform_h set status=0,confirmerid=0,confirmtime=to_timestamp(0),ts=current_timestamp 
	where id=$1 and dr=0 and status=1 and ts=$2`
	headRes, err := tx.Exec(confirmHeadSql, pif.HID, pif.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.PPEIssuanceForm.Unconfirm tx.Exec(confirmHeadSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of affected rows
	confirmHeadNumber, err := headRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Unconfirm headRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if confirmHeadNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Prepare to unconfirm body rows
	confirmRowSql := `update ppeissuanceform_b set status=0,confirmerid=0,confirmtime=to_timestamp(0),ts=current_timestamp 
	where id=$1 and dr=0 and status=1 and ts=$2`
	rowStmt, err := tx.Prepare(confirmRowSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PPEIssuanceForm.Unconfirm tx.Prepare(confirmRowSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer rowStmt.Close()
	// Write body row unconfirmations one by one
	for _, row := range pif.Body {
		// Check the status of the body row
		if row.Status != 1 {
			resStatus = i18n.StatusVoucherNoConfirm
			tx.Rollback()
			return
		}
		confirmRowRes, errConfirmRow := rowStmt.Exec(row.BID, row.Ts)
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEIssuanceForm.Unconfirm rowStmt.Exec failed", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}

		confirmRowNumber, errConfirmRow := confirmRowRes.RowsAffected()
		if errConfirmRow != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PPEIssuanceForm.Unconfirm confirmRowRes.RowsAffected failed", zap.Error(errConfirmRow))
			tx.Rollback()
			return resStatus, errConfirmRow
		}
		if confirmRowNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}
	}
	return i18n.StatusOK, nil
}

// Get PPE Issuance Form Report
func GetPPEIFReport(queryString string) (pifrs []PPEIssuanceFormReport, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	pifrs = make([]PPEIssuanceFormReport, 0)
	var build strings.Builder
	// Concatenate SQL to check the number of records
	build.WriteString(`select count(b.hid) as rowcount 
	from ppeissuanceform_b as b
	left join ppeissuanceform_h as h on b.hid = h.id
	left join ppe on b.ppeid = ppe.id
	left join sysuser as recipient on b.recipientid = recipient.id
	left join department as issuedept on h.deptid = issuedept.id
	left join sysuser as creator on h.creatorid=creator.id
	where (b.dr=0 and h.dr=0) `)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	checkSql := build.String()
	// Check the number of records
	var rowNumber int32
	err = db.QueryRow(checkSql).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetPPEIFReport db.QueryRow(checkSql) failed", zap.Error(err))
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

	// Concatenate SQL to get report data
	build.WriteString(`select b.hid as hid,
	b.id as bid,
	b.rownumber as rownumber,
	b.recipientid as recipientid,
	coalesce(recipient.code,'') as recipientcode,
	coalesce(recipient.name,'') as recipientname,
	b.positionname as recipientpoisitionname,
	b.deptname as recipientdeptname,
	b.ppeid as ppeid,
	coalesce(ppe.code,'') as ppecode,
	coalesce(ppe.name,'') as ppename,
	coalesce(ppe.model,'') as ppemodel,
	coalesce(ppe.unit,'') as ppeunit,
	b.quantity as quantity,
	b.description as bdescription,
	b.status as bstatus,
	h.billnumber as billnumber,
	h.billdate as billdate,
	h.deptid as issuingdeptid,
	coalesce(issuedept.code,'') as issuingdeptcode,
	coalesce(issuedept.name,'') as issuingdeptname,
	h.description as hdescription,
	h.period as period,
	h.startdate as startdate,
	h.enddate as enddate,
	h.sourcetype as sourcetype,
	h.status as hstatus,
	h.creatorid as creatorid,
	coalesce(creator.code,'') as creatorcode,
	coalesce(creator.name,'') as creatorname
	from ppeissuanceform_b as b
	left join ppeissuanceform_h as h on b.hid = h.id
	left join ppe on b.ppeid = ppe.id
	left join sysuser as recipient on b.recipientid = recipient.id
	left join department as issuedept on h.deptid = issuedept.id
	left join sysuser as creator on h.creatorid=creator.id
	where (b.dr=0 and h.dr=0)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	repSql := build.String()
	// Get report data
	ldRep, err := db.Query(repSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetQueryDocumentReport db.Query failed", zap.Error(err))
		return
	}
	defer ldRep.Close()
	// Extract data row by row
	for ldRep.Next() {
		var pifr PPEIssuanceFormReport
		err = ldRep.Scan(&pifr.HID, &pifr.BID, &pifr.RowNumber, &pifr.RecipientID, &pifr.RecipientCode,
			&pifr.RecipientName, &pifr.RecipientPositionName, &pifr.RecipientDeptName, &pifr.PPEID, &pifr.PPECode,
			&pifr.PPEName, &pifr.PPEModel, &pifr.PPEUnit, &pifr.Quantity, &pifr.BDescription,
			&pifr.BStatus, &pifr.Billnumber, &pifr.BillDate, &pifr.IssuingDeptID, &pifr.IssuingDeptCode,
			&pifr.IssuingDeptName, &pifr.HDescription, &pifr.Period, &pifr.StartDate, &pifr.EndDate,
			&pifr.SourceType, &pifr.Hstatus, &pifr.CreateUserID, &pifr.CreateUserCode, &pifr.CreateUserName)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetQueryDocumentReport ldRep.Next() ldRep.Scan failed", zap.Error(err))
			return
		}

		pifrs = append(pifrs, pifr)
	}

	return
}
