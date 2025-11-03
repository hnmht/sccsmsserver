package pg

import (
	"encoding/json"
	"math"
	"sccsmsserver/cache"
	"sccsmsserver/i18n"
	"sccsmsserver/pub"
	"sccsmsserver/setting"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Document Struct
type Document struct {
	ID          int32         `db:"id" json:"id"`
	DC          SimpDC        `db:"dcid" json:"dc"`
	Name        string        `db:"name" json:"name"`
	Edition     string        `db:"edition" json:"edition"`
	Author      string        `db:"author" json:"author"`
	UploadDate  time.Time     `db:"uploaddate" json:"uploadDate"`
	ReleaseDate time.Time     `db:"releasedate" json:"releaseDate"`
	Tags        string        `db:"tags" json:"tags"`
	Description string        `db:"description" json:"description"`
	Files       []VoucherFile `json:"files"`
	CreateDate  time.Time     `db:"createtime" json:"createDate"`
	Creator     Person        `db:"creatorid" json:"creator"`
	ModifyDate  time.Time     `db:"modifytime" json:"modifyDate"`
	Modifier    Person        `db:"modifierid" json:"modifier"`
	Ts          time.Time     `db:"ts" json:"ts"`
	Dr          int16         `db:"dr" json:"dr"`
}

// Document query pagination parameters by category
type DCPagingParams struct {
	DC      SimpDC     `json:"dc"`
	Count   int32      `json:"count"`
	Page    int32      `json:"page"`
	PerPage int32      `json:"perPage"`
	Docs    []Document `json:"docs"`
}

// Document Report struct
type QueryDocument struct {
	DocID       int32         `json:"docID"`
	DocName     string        `json:"docName"`
	DCID        int32         `json:"dcID"`
	DCName      string        `json:"dcName"`
	Edition     string        `db:"edition" json:"edition"`
	Author      string        `db:"author" json:"author"`
	UploadDate  time.Time     `db:"uploaddate" json:"uploadDate"`
	ReleaseDate time.Time     `db:"releasedate" json:"releaseDate"`
	Description string        `db:"description" json:"description"`
	Files       []VoucherFile `json:"files"`
	CreatorID   int32         `json:"creatorID"`
	CreatorCode string        `json:"creatorCode"`
	CreatorName string        `json:"creatorName"`
}

// Get Document pagination data list based on document category
func (dpp *DCPagingParams) Get() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	dpp.Docs = make([]Document, 0)
	// Check if the document exists
	checkSql := `select count(id) from document where dr=0 and dcid=$1`
	err = db.QueryRow(checkSql, dpp.DC.ID).Scan(&dpp.Count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DCPagingParams.Get db.QueryRow(checksql) failed:", zap.Error(err))
		return
	}
	if dpp.Count == 0 {
		resStatus = i18n.StatusResNoData
		return
	}
	// Recalculate pagination
	if dpp.PerPage > dpp.Count {
		dpp.Page = 0
	} else {
		var totalPage = int32(math.Ceil(float64(dpp.Count) / float64(dpp.PerPage)))
		if (dpp.Page + 1) > totalPage {
			dpp.Page = totalPage - 1
		}
	}
	// Retrieve document list from database
	querySql := `select id from document where dr=0 and dcid=$1 order by id limit $2 offset $3`
	rows, err := db.Query(querySql, dpp.DC.ID, dpp.PerPage, dpp.Page*dpp.PerPage)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DCPagingParams.Get db.Query(querySql) failed:", zap.Error(err))
		return
	}
	// Extract document data row by row
	for rows.Next() {
		var d Document
		err = rows.Scan(&d.ID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("DCPagingParams.Getrows.Scan failed:", zap.Error(err))
			return
		}
		resStatus, err = d.GetDetailByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		dpp.Docs = append(dpp.Docs, d)
	}

	return
}

// Get Document Attachments
func GetDocumentFiles(did int32) (voucherFiles []VoucherFile, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	voucherFiles = make([]VoucherFile, 0)
	// Get Attachments from database
	filesStr := `select id,billhid,fileid,createtime,creatorid,
		modifytime,modifierid,ts,dr 
		from document_file 
		where dr=0 and billhid=$1`
	rows, err := db.Query(filesStr, did)
	if err != nil {
		zap.L().Error("GetDocumentFiles db.Query(files)  failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Extract file row by row
	for rows.Next() {
		var file VoucherFile
		err = rows.Scan(&file.ID, &file.BillHID, &file.File.ID, &file.CreateDate, &file.Creator.ID,
			&file.ModifyDate, &file.Modifier.ID, &file.Ts, &file.Dr)
		if err != nil {
			zap.L().Error("Document.GetDetailByID db.Query(files)  rows.Scan  failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}

		// Get File details
		if file.File.ID > 0 {
			resStatus, err = file.File.GetFileInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Creator details
		if file.Creator.ID > 0 {
			resStatus, err = file.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Modifier details
		if file.Modifier.ID > 0 {
			resStatus, err = file.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		voucherFiles = append(voucherFiles, file)
	}
	return
}

// Get Document details by ID
func (d *Document) GetDetailByID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get document details from Cache
	number, b, _ := cache.Get(pub.Document, d.ID)
	if number > 0 {
		json.Unmarshal(b, &d)
		return
	}
	// If Document isn't in cache, retrieve it from database
	sqlStr := `select dcid,name,edition,author,uploaddate,releasedate,
		tags,description,createtime,creatorid,modifytime,
		modifierid,ts,dr 
		from document where id=$1`
	err = db.QueryRow(sqlStr, d.ID).Scan(&d.DC.ID, &d.Name, &d.Edition, &d.Author, &d.UploadDate, &d.ReleaseDate,
		&d.Tags, &d.Description, &d.CreateDate, &d.Creator.ID, &d.ModifyDate,
		&d.Modifier.ID, &d.Ts, &d.Dr)
	if err != nil {
		zap.L().Error("Document.GetDetailByID db.QueryRow  failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}

	// Get Document Category details
	if d.DC.ID > 0 {
		resStatus, err = d.DC.GetSDCInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Creator details
	if d.Creator.ID > 0 {
		resStatus, err = d.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Modifier details
	if d.Modifier.ID > 0 {
		resStatus, err = d.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Document attachments
	d.Files, resStatus, err = GetDocumentFiles(d.ID)
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Write into cache
	dB, _ := json.Marshal(d)
	cache.Set(pub.Document, d.ID, dB)

	return
}

// Add Document
func (d *Document) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check the number of files, zero is not allowed
	if len(d.Files) == 0 {
		resStatus = i18n.StatusDocumentNoFile
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Document.Add db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Insert document record into the document table
	headSql := `insert into document(dcid,name,edition,author,releasedate,
		tags,description,creatorid) 
		values($1,$2,$3,$4,$5,$6,$7,$8)
		returning id`
	err = tx.QueryRow(headSql, d.DC.ID, d.Name, d.Edition, d.Author, d.ReleaseDate,
		d.Tags, d.Description, d.Creator.ID).Scan(&d.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Document.Add tx.QueryRow failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Insert Document Attachments to the document_file table
	fileSql := `insert into document_file(billhid,fileid,creatorid) 
	values($1,$2,$3) 
	returning id`
	fileStmt, err := tx.Prepare(fileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Document.Add  tx.Prepare(fileSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer fileStmt.Close()
	for _, file := range d.Files {
		err = fileStmt.QueryRow(d.ID, file.File.ID, d.Creator.ID).Scan(&file.ID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("Document.Add  fileStmt.QueryRow failed", zap.Error(err))
			tx.Rollback()
			return
		}
	}
	return
}

// Edit Document
func (d *Document) Edit() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if  creator and modifier are the same person
	if d.Creator.ID != d.Modifier.ID {
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Check the number of files, zero is not allowed
	var fileNumber int32
	for _, file := range d.Files {
		if file.Dr == 0 {
			fileNumber++
		}
	}
	if fileNumber == 0 {
		resStatus = i18n.StatusDocumentNoFile
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Document.Edit db.Begin() failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Modify the Document info in the document table
	editDocSql := `update document set dcid=$1,name=$2,edition=$3,author=$4,releasedate=$5,
		tags=$6,description=$7,modifytime=current_timestamp,modifierid=$8,ts=current_timestamp 
		where id=$9 and dr=0 and ts=$10`
	editDocRes, err := tx.Exec(editDocSql, &d.DC.ID, &d.Name, &d.Edition, &d.Author, &d.ReleaseDate,
		&d.Tags, &d.Description, &d.Modifier.ID,
		&d.ID, &d.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Document.Edit tx.Exec(editDocSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of affected rows
	updateNumber, err := editDocRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Document.Edit editDocRes.RowsAffected failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if updateNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}

	// Update the Attachments info
	updateFileSql := `update document_file set modifierid=$1,modifytime=current_timestamp,dr=$2,ts=current_timestamp 
		where id=$3 and billhid=$4 and dr=0 and ts=$5`
	addFileSql := `insert into document_file(billhid,fileid,creatorid) values($1,$2,$3) returning id`
	// Prepare update attachments in the document table
	updateFileStmt, err := tx.Prepare(updateFileSql)
	if err != nil {
		zap.L().Error("Document.Edit tx.Prepare(updateFileSql) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer updateFileStmt.Close()
	// Prepare Add attachments in the document table
	addFileStmt, err := tx.Prepare(addFileSql)
	if err != nil {
		zap.L().Error("Document.Edit tx.Prepare(addFileStmt) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		tx.Rollback()
		return
	}
	defer addFileStmt.Close()

	for _, file := range d.Files {
		if file.ID != 0 { // If the file.ID is not 0, it means it is a row that needs to be modified
			updateFileRes, updateFileErr := updateFileStmt.Exec(d.Modifier.ID, file.Dr, file.ID, d.ID, file.Ts)
			if updateFileErr != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("Document.Edit  updateFileRes.Exec() failed", zap.Error(updateFileErr))
				tx.Rollback()
				return resStatus, updateFileErr
			}
			updateFileNumber, updateFileEffectErr := updateFileRes.RowsAffected()
			if updateFileEffectErr != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("Document.EditupdateFileRes.RowsAffected failed", zap.Error(updateFileEffectErr))
				tx.Rollback()
				return resStatus, updateFileEffectErr
			}
			if updateFileNumber < 1 {
				resStatus = i18n.StatusOtherEdit
				tx.Rollback()
				return
			}
		} else { // If the file.ID is 0, it means it is a new file
			addFileErr := addFileStmt.QueryRow(d.ID, file.File.ID, d.Modifier.ID).Scan(&file.ID)
			if addFileErr != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("Document.Edit addFileStmt.QueryRow failed", zap.Error(addFileErr))
				tx.Rollback()
				return resStatus, addFileErr
			}
		}
	}
	// Delete from cache
	d.DelFromLocalCache()
	return
}

// Delete Document from cache
func (d *Document) DelFromLocalCache() {
	number, _, _ := cache.Get(pub.Document, d.ID)
	if number > 0 {
		cache.Del(pub.Document, d.ID)
	}
}

// Delete Document
func (d *Document) Delete(modifyUserId int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the Document is referenced
	// Check the operator and creator are a same person
	if d.Creator.ID != modifyUserId {
		resStatus = i18n.StatusOtherEdit
		return
	}
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Document.Del db.Begin() failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Modify the delete flag for the document in the document table
	delDocSql := `update document set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
		where id=$2 and dr=0 and ts=$3`
	delDocRes, err := tx.Exec(delDocSql, modifyUserId, d.ID, d.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Document.Delete tx.Exec(delDocSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	// Check the number of rows affected by SQL statement
	delDocNumber, err := delDocRes.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Document.Delete delDocRes.RowsAffected() failed", zap.Error(err))
		tx.Rollback()
		return
	}
	if delDocNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		tx.Rollback()
		return
	}
	// Modify the delete flag for the document attachements in the document_file table
	delFileSql := `update document_file set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
		where id=$2 and dr=0 and billhid=$3 and ts=$4`
	// Prepare update the document file
	delFileStmt, err := tx.Prepare(delFileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Document.Delete tx.Prepare(delFileSql) failed", zap.Error(err))
		tx.Rollback()
		return
	}
	defer delFileStmt.Close()
	// update the file delete flag one by one
	for _, file := range d.Files {
		delFileRes, errDelFile := delFileStmt.Exec(modifyUserId, file.ID, file.BillHID, file.Ts)
		if errDelFile != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("Document.Delete delFileStmt.Exec failed", zap.Error(errDelFile))
			tx.Rollback()
			return
		}
		delFileNumber, errDelEff := delFileRes.RowsAffected()
		if errDelEff != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("Document.Delete delFileRes.RowsAffected failed", zap.Error(errDelEff))
			tx.Rollback()
			return
		}
		if delFileNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}
	}
	// Delete from cache
	d.DelFromLocalCache()

	return
}

// Batche Delete Documents
func DeleteDocuments(docs *[]Document, modifierID int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Begin a database transaction
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteDocuments db.Begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()

	// Prepare modify the delete flag in the document table
	delDocSql := `update document set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
		where id=$2 and dr=0 and ts=$3`
	docStmt, err := tx.Prepare(delDocSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteDocuments tx.Prepare(delDocSql) failed", zap.Error(err))
		tx.Rollback()
		return resStatus, err
	}
	defer docStmt.Close()
	// Prepare modify the delete flag in the document_file table
	delFileSql := `update document_file set dr=1,modifytime=current_timestamp,modifierid=$1,ts=current_timestamp 
		where id=$2 and dr=0 and billhid=$3 and ts=$4`
	fileStmt, err := tx.Prepare(delFileSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteDocuments tx.Prepare(delFileSql) failed", zap.Error(err))
		tx.Rollback()
		return resStatus, err
	}
	defer fileStmt.Close()

	// Modify Document one by one
	for _, d := range *docs {
		// Check the creator and modifier are same person
		if d.Creator.ID != modifierID {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}
		// Modify record in the document table
		delDocRes, errDelDoc := docStmt.Exec(modifierID, d.ID, d.Ts)
		if errDelDoc != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("DeleteDocuments docStmt.Exec failed", zap.Error(errDelDoc))
			tx.Rollback()
			return
		}
		// Check the number of rows effected by SQL statement
		delDocNumber, errDelDocEff := delDocRes.RowsAffected()
		if errDelDocEff != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("DeleteDocuments delDocRes.RowsAffected failed", zap.Error(errDelDocEff))
			tx.Rollback()
			return
		}
		if delDocNumber < 1 {
			resStatus = i18n.StatusOtherEdit
			tx.Rollback()
			return
		}

		// Modify attachment record in the document_file table
		for _, file := range d.Files {
			delFileRes, errDelFile := fileStmt.Exec(modifierID, file.ID, file.BillHID, file.Ts)
			if errDelFile != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("DeleteDocuments fileStmt.Exec failed", zap.Error(errDelFile))
				tx.Rollback()
				return
			}
			delFileNumber, errDelEff := delFileRes.RowsAffected()
			if errDelEff != nil {
				resStatus = i18n.StatusInternalError
				zap.L().Error("DeleteDocuments delFileRes.RowsAffected failed", zap.Error(errDelEff))
				tx.Rollback()
				return
			}
			if delFileNumber < 1 {
				resStatus = i18n.StatusOtherEdit
				tx.Rollback()
				return
			}
		}
		// Delte from cache
		d.DelFromLocalCache()
	}

	return
}

// Get the document Report
func GetQueryDocumentReport(queryString string) (qds []QueryDocument, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	qds = make([]QueryDocument, 0)
	var build strings.Builder
	// Concatenate the SQL string for check
	build.WriteString(`select count(d.id) as rowcount
	from document as d
	left join dc on d.dcid = dc.id
	left join sysuser as creator on d.creatorid = creator.id
	where d.dr=0 `)
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
		zap.L().Error("GetQueryDocumentReport db.QueryRow(checkSql) failed", zap.Error(err))
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

	// Concatenate the SQL string for data retrieve
	build.WriteString(`select d.id as docid,
	d.name as docname,
	d.dcid as dcid,
	coalesce(dc.name,'') as dcname,
	d.edition as edition,
	d.author as author,
	d.uploaddate as uploaddate,
	d.releasedate as releasedate,
	d.description as desecription,
	d.creatorid as creatorid,
	coalesce(creator.code,'') as creatorcode,
	coalesce(creator.name,'') as creatorname
	from document as d
	left join dc on d.dcid = dc.id
	left join sysuser as creator on d.creatorid = creator.id
	where d.dr=0`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	repSql := build.String()
	// Retrieve Document report from database
	qdRep, err := db.Query(repSql)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetQueryDocumentReport db.Query failed", zap.Error(err))
		return
	}
	defer qdRep.Close()

	// Extract data row by row
	for qdRep.Next() {
		var qd QueryDocument
		err = qdRep.Scan(&qd.DocID, &qd.DocName, &qd.DCID, &qd.DCName, &qd.Edition,
			&qd.Author, &qd.UploadDate, &qd.ReleaseDate, &qd.Description, &qd.CreatorID,
			&qd.CreatorCode, &qd.CreatorName)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetQueryDocumentReport qdRep.Next() qdRep.Scan failed", zap.Error(err))
			return
		}
		// Get Document attachments
		qd.Files, resStatus, err = GetDocumentFiles(qd.DocID)
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		qds = append(qds, qd)
	}
	return
}
