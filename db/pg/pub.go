package pg

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"sccsmsserver/i18n"
	"sccsmsserver/pkg/mysf"
	"sccsmsserver/pkg/security"
	"time"

	"go.uber.org/zap"
)

// SeaCloud Date Type struct
type ScDataType struct {
	ID        int32  `json:"id"`
	TypeCode  string `json:"code"`
	TypeName  string `json:"name"`
	DataType  string `json:"dataType"`
	FrontDb   string `json:"frontDb"`
	InputMode string `json:"inputMode"`
}

// Data references check struct
type DataReferenceCheck struct {
	Description    string
	SqlStr         string
	UsedReturnCode i18n.ResKey
}

// Data query conditions
type QueryParams struct {
	QueryString string `json:"queryString"`
}

// Data query pagination
type PagingQueryParams struct {
	QueryString string `json:"queryString"`
	Page        int32  `json:"page"`
	PerPage     int32  `json:"perPage"`
}

// the struct Check the archive is refreneced
type ArchiveCheckUsed struct {
	Description    string
	SqlStr         string
	UsedReturnCode i18n.ResKey
}

// Frontend Database information struct
type FrontDBInfo struct {
	ID         int32     `db:"id" json:"id"`
	DbID       int64     `db:"dbid" json:"dbID,string"`
	FrontDbID  int64     `db:"frontdbid" json:"frontDbID,string"`
	CryptoKey  string    `db:"cryptokey" json:"cryptoKey"`
	CreateDate time.Time `db:"createtime" json:"createDate"`
	Creator    Person    `db:"creatorid" json:"creator"`
	Dr         int16     `db:"dr" json:"dr"`
	Ts         time.Time `db:"ts" json:"ts"`
}

// Generate Frontend DBID
func (f *FrontDBInfo) Generate() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check wether DBID is empty
	if f.DbID == 0 {
		resStatus = i18n.StatusDBIDEmpty
		return
	}
	// Check wether DBID is Match
	if f.DbID != ServerPubInfo.DbID {
		resStatus = i18n.StatusDBIDMissMatch
		return
	}
	// Generate Front-end DBID
	f.FrontDbID = mysf.GenID()
	// Generate CryptoKey
	byteKey, err := security.GenerateAESKey(32)
	if err != nil {
		zap.L().Error("FrontDbInfo.Generate security.GenerateAESKey failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	f.CryptoKey = base64.StdEncoding.EncodeToString(byteKey)

	// Write record into frontdb table
	sqlStr := `insert into frontdb(dbid,frontdbid,cryptokey,creatorid) 
	values($1,$2,$3,$4) 
	returning id,createtime,dr,ts`
	err = db.QueryRow(sqlStr, f.DbID, f.FrontDbID, f.CryptoKey, f.Creator.ID).Scan(&f.ID, &f.CreateDate, &f.Dr, &f.Ts)
	if err != nil {
		zap.L().Error("FrontDbInfo.Generate db.QueryRow failed:", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	return
}

// Get Frontend DBID info
func (f *FrontDBInfo) GetInfo() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check wether DBID is empty
	if f.DbID == 0 {
		resStatus = i18n.StatusDBIDEmpty
		return
	}
	// Check wether DBID is Match
	if f.DbID != ServerPubInfo.DbID {
		resStatus = i18n.StatusDBIDMissMatch
		return
	}
	// Retrieve the frontdb information from the frontdb table
	sqlStr := `select id,cryptokey from frontdb where dbid=$1 and frontdbid=$2 and dr=0 order by ts desc limit(1)`
	err = db.QueryRow(sqlStr, f.DbID, f.FrontDbID).Scan(&f.ID, &f.CryptoKey)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Info("FrontDBInfo.GetInfo db.QueryRow No match data")
			resStatus = i18n.StatusFrontDBIDNoRow
			return
		}
		zap.L().Error("FrontDBInfo.GetInfo db.QueryRow failed:", zap.Error(err))
		return
	}
	return
}

// Get the latest voucher sequence number
func GetLatestSerialNo(tx *sql.Tx, voucherType string, voucherDate string) (serialno string, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Query the latest sequence for the current voucher type
	var sn int32
	sqlString := "select serialno from serialno where vouchertype=$1 and datestring=$2"
	err = tx.QueryRow(sqlString, voucherType, voucherDate).Scan(&sn)
	if err != nil && err != sql.ErrNoRows {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetLatestSerialNo db.QueryRow failed", zap.Error(err))
		return
	}
	// If no data is queried, then create a new record for the current voucher type
	// and insert it into the database
	if err == sql.ErrNoRows {
		sqlString = `insert into serialno(datestring,vouchertype) values($1,$2)`
		_, err = tx.Exec(sqlString, voucherDate, voucherType)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetLatestSerial db.exec Insert failed", zap.Error(err))
			return
		}
	}
	// Update the serial number for the current voucher type in the database
	sqlString = `update serialno set serialno = serialno + 1 where datestring = $1 and vouchertype = $2`
	_, err = tx.Exec(sqlString, voucherDate, voucherType)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetLatestSerial db.Exec Update failed", zap.Error(err))
		return
	}
	serialno = fmt.Sprintf("%v%v%04d", voucherType, voucherDate, sn+1)

	return
}

/* // Cancel the voucher serial number
func CancelSerialNo(voucherType string, voucherDate string) {
	sqlStr := `update serialno set serialno = serialno-1 where dateString= $1 and vouchertype= $2`
	_, err := db.Exec(sqlStr, voucherType, voucherDate)
	if err != nil {
		zap.L().Error("UpdateSerialNo db.Exec failed", zap.Error(err))
	}
} */
