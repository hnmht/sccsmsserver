package pg

import (
	"database/sql"
	"fmt"
	"sccsmsserver/i18n"
	"sccsmsserver/setting"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Execution Order Comment Message struct
type CommentMessage struct {
	ID                 int32         `db:"id" json:"id"`
	HID                int32         `db:"hid" json:"hid"`
	BID                int32         `db:"bid" json:"bid"`
	RowNUmber          int32         `db:"rownumber" json:"rowNumber"`
	BillNumber         string        `db:"billnumber" json:"billNumber"`
	CSAID              int32         `json:"csaID"`
	CSACode            string        `json:"csaCode"`
	CSAName            string        `json:"csaName"`
	EPAID              int32         `json:"epaID"`
	EPACode            string        `json:"epaCode"`
	EPAName            string        `json:"epaName"`
	ExecutionValueDisp string        `json:"executionValueDisp"`
	EOFiles            []VoucherFile `json:"eoFiles"`
	SendTo             Person        `db:"sendtoid" json:"sendTo"`
	IsRead             int16         `db:"isread" json:"isRead"`
	ReadTime           time.Time     `db:"readtime" json:"readtime"`
	Content            string        `db:"content" json:"content"`
	SendTime           time.Time     `db:"sendtime" json:"sendTime"`
	CreateDate         time.Time     `db:"createtime" json:"createDate"`
	Creator            Person        `db:"creatorid" json:"creator"`
	ModifyDate         time.Time     `db:"modifytime" json:"modifyDate"`
	Modifier           Person        `db:"modifierid" json:"modifier"`
	Ts                 time.Time     `db:"ts" json:"ts"`
	Dr                 int16         `db:"dr" json:"dr"`
}

// Get User UnRead Comments list
func GetUserUnReadComments(userID int32) (comments []CommentMessage, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	comments = make([]CommentMessage, 0)
	var build strings.Builder
	// Concatenate the SQL string for inspection
	build.WriteString(`select count(c.id) as rowcount 				
	from executionorder_comment as c
	left join executionorder_h as h on c.hid = h.id
	left join executionorder_b as b on c.bid = b.id
	left join epa on b.epaid = epa.id
	left join csa on h.csaid = csa.id
	where  (c.dr=0 and b.dr=0 and c.isread=0) `)
	build.WriteString(`and (c.sendtoid=`)
	build.WriteString(strconv.Itoa(int(userID)))
	build.WriteString(`)`)
	// CHeck
	checkSql := build.String()
	var rowNumber int32
	err = db.QueryRow(checkSql).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetUserUnReadComments db.QueryRow(checkSql) failed", zap.Error(err))
		return
	}
	if rowNumber == 0 {
		resStatus = i18n.StatusResNoData
		return
	}
	// Check MaxRecord
	if rowNumber > setting.Conf.PqConfig.MaxRecord {
		resStatus = i18n.StatusOverRecord
		return
	}
	build.Reset()

	// Concatenate the SQL for data retrieval.
	build.WriteString(`select c.id as id,
	c.bid as bid,
	c.hid as hid,
	c.rownumber as rownumber,
	c.billnumber as billnumber,
	h.csaid as csaid,
	csa.code as csacode,
	csa.name as csaname,
	b.epaid as epaid,
	epa.code as epacode,
	epa.name as epaname,
	b.executionvaluedisp as executionvaluedisp,
	c.sendtoid as sendtoid,
	c.isread as isread,
	c.readtime as readtime,
	c.content as content,
	c.sendtime as sendtime,
	c.createtime as createtime,
	c.creatorid as creatorid,
	c.dr as dr,
	c.ts as ts
	from executionorder_comment as c
	left join executionorder_h as h on c.hid = h.id
	left join executionorder_b as b on c.bid = b.id
	left join epa on b.epaid = epa.id
	left join csa on h.csaid = csa.id
	where  (c.dr=0 and b.dr=0 and c.isread=0)`)
	build.WriteString(`and (c.sendtoid=`)
	build.WriteString(strconv.Itoa(int(userID)))
	build.WriteString(`) order by c.id desc`)
	sqlStr := build.String()
	// Retrieve UnRead comments list from database
	res, err := db.Query(sqlStr)
	if err != nil {
		if err == sql.ErrNoRows {
			resStatus = i18n.StatusResNoData
			return
		}
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetUserUnReadComments db.Query failed", zap.Error(err))
		return
	}
	defer res.Close()

	// Extract data row by row
	for res.Next() {
		var cm CommentMessage
		err = res.Scan(&cm.ID, &cm.BID, &cm.HID, &cm.RowNUmber, &cm.BillNumber,
			&cm.CSAID, &cm.CSACode, &cm.CSAName, &cm.EPAID, &cm.EPACode,
			&cm.EPAName, &cm.ExecutionValueDisp, &cm.SendTo.ID, &cm.IsRead, &cm.ReadTime,
			&cm.Content, &cm.SendTime, &cm.CreateDate, &cm.Creator.ID, &cm.Dr,
			&cm.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetUserUnReadComments res.scan failed", zap.Error(err))
			return
		}
		// Get SendTo Person details
		if cm.SendTo.ID > 0 {
			resStatus, err = cm.SendTo.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Creator details
		if cm.Creator.ID > 0 {
			resStatus, err = cm.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Attachements
		cm.EOFiles, resStatus, err = GetEORowFiles(cm.BID)
		if resStatus != i18n.StatusOK || err != nil {
			return
		}

		comments = append(comments, cm)
	}

	return
}

// Get User Read Comments
func GetUserReadComments(userID int32, queryString string) (comments []CommentMessage, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	comments = make([]CommentMessage, 0)
	var build strings.Builder
	// Concatenate the SQL string for inspection
	build.WriteString(`select count(c.id) as rowcount 
	from executionorder_comment as c
	left join executionorder_h as h on c.hid = h.id
	left join executionorder_b as b on c.bid = b.id
	left join epa on b.epaid = epa.id
	left join csa on h.csaid = csa.id
	where  (c.dr=0 and b.dr=0 and c.isread=1 and c.sendtoid=`)
	build.WriteString(strconv.Itoa(int(userID)))
	build.WriteString(`)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	// Check
	checkSql := build.String()
	var rowNumber int32
	err = db.QueryRow(checkSql).Scan(&rowNumber)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetUserReadComments db.QueryRow(checkSql) failed", zap.Error(err))
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

	// Concatenate the SQL from data retrieve
	build.WriteString(`select c.id as id,
	c.bid as bid,
	c.hid as hid,
	c.rownumber as rownumber,
	c.billnumber as billnumber,
	h.csaid as csaid,
	csa.code as csacode,
	csa.name as csaname,
	b.epaid as epaid,
	epa.code as epacode,
	epa.name as epaname,
	b.executionvaluedisp as executioncaluedisp,
	c.sendtoid as sendtoid,
	c.isread as isread,
	c.readtime as readtime,
	c.content as content,
	c.sendtime as sendtime,
	c.createtime as createtime,
	c.creatorid as creatorid,
	c.dr as dr,
	c.ts as ts
	from executionorder_comment as c
	left join executionorder_h as h on c.hid = h.id
	left join executionorder_b as b on c.bid = b.id
	left join epa on b.epaid = epa.id
	left join csa on h.csaid = csa.id
	where  (c.dr=0 and b.dr=0 and c.isread=1 and c.sendtoid=`)
	build.WriteString(strconv.Itoa(int(userID)))
	build.WriteString(`)`)
	if queryString != "" {
		build.WriteString(" and (")
		build.WriteString(queryString)
		build.WriteString(")")
	}
	build.WriteString(` order by c.id desc`)
	sqlStr := build.String()
	// Retrieve comments from database
	res, err := db.Query(sqlStr)
	if err != nil {
		if err == sql.ErrNoRows {
			resStatus = i18n.StatusResNoData
			return
		}
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetUserReadComments db.Query(sqlStr) failed", zap.Error(err))
		return
	}
	defer res.Close()

	// Extract comment row by row
	for res.Next() {
		var cm CommentMessage
		err = res.Scan(&cm.ID, &cm.BID, &cm.HID, &cm.RowNUmber, &cm.BillNumber,
			&cm.CSAID, &cm.CSACode, &cm.CSAName, &cm.EPAID, &cm.EPACode,
			&cm.EPAName, &cm.ExecutionValueDisp, &cm.SendTo.ID, &cm.IsRead, &cm.ReadTime,
			&cm.Content, &cm.SendTime, &cm.CreateDate, &cm.Creator.ID, &cm.Dr,
			&cm.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetUserReadComments res.scan failed", zap.Error(err))
			return
		}
		// Get SendTo user details
		if cm.SendTo.ID > 0 {
			resStatus, err = cm.SendTo.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Creator details
		if cm.Creator.ID > 0 {
			resStatus, err = cm.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Attachements
		cm.EOFiles, resStatus, err = GetEORowFiles(cm.BID)
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
		comments = append(comments, cm)
	}

	return
}

// Get User To-Do Issues
func GetUserEORefs(userID int32) (reds []ReferExecutionOrder, resStatus i18n.ResKey, err error) {
	queryString := fmt.Sprintf("issueownerid=%d", userID)
	reds, resStatus, err = GetReferEOs(queryString)
	return
}

// Get User work Orders awaiting execution
func GetUserWORefs(userID int32) (wors []WorkOrderRow, resStauts i18n.ResKey, err error) {
	queryString := fmt.Sprintf("executorid=%d", userID)
	wors, resStauts, err = GetWORefer(queryString)
	return
}

// User Read Execution Order Comment
func (cm *CommentMessage) Read() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check the Send To User and current User are the same person
	if cm.SendTo.ID != cm.Modifier.ID {
		resStatus = i18n.StatusMsgOnlyReadSelf
		return
	}
	// Update the record in the executionorder_comment
	sqlStr := `update executionorder_comment 
	set isread=1,current_timestamp,modifierid=$1,modifytime=current_timestamp,ts=current_timestamp 
	where dr=0 and isread=0 and id=$2 and ts=$3`
	res, err := db.Exec(sqlStr, cm.Modifier.ID, cm.ID, cm.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("CommentMessage.Read db.exec failed", zap.Error(err))
		return
	}
	// Check the number of rows affected by SQL statement
	updateNumber, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("CommentMessage.Read res.RowsAffected falied", zap.Error(err))
		return
	}
	if updateNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		return
	}
	return
}
