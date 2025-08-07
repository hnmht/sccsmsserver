package pg

import (
	"time"

	"go.uber.org/zap"
)

// Document Category Master Data
type DC struct {
	ID          int32     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	HigherClass SimpDC    `db:"fatherid" json:"fatherID"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Creator     Person    `db:"creatorid" json:"creator"`
	ModifyDate  time.Time `db:"modifytime" json:"modifyDate"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

// Simplify Document category
type SimpDC struct {
	ID          int32     `db:"id" json:"id"`
	ClassName   string    `db:"classname" json:"name"`
	Description string    `db:"description" json:"description"`
	FatherID    int32     `db:"fatherid" json:"fatherID"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Creator     Person    `db:"creatorid" json:"creator"`
	ModifyDate  time.Time `db:"modifytime" json:"modifyDate"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

// Initialize Document Category table
func initDocumentCategory() (isFinish bool, err error) {
	// Step 1: Check if a record exists for the defualt Document Category.
	sqlStr := "select count(id) as rownum from dc where id=1"
	// Step 2: Exit if the record exists or an error occurs.
	hasRecord, isFinish, err := genericCheckRecord("dc", sqlStr)
	if hasRecord || !isFinish || err != nil {
		return
	}
	// Step 3: Insert a record for the system default DC "Default Category" into the
	sqlStr = `insert into dc(id,name,description,creatorid) values(10000,'Default Category','System Pre-Set',10000)`
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initDocumentCategory insert initvalue failed", zap.Error(err))
		return isFinish, err
	}
	return
}
