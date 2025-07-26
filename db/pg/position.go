package pg

import (
	"time"

	"go.uber.org/zap"
)

// Position 岗位
type Position struct {
	ID          int32     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Creator     Person    `db:"creatorid" json:"creator"`
	ModifyDate  time.Time `db:"modify_time" json:"modifyDate"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	Dr          int16     `db:"dr" json:"dr"`
	Ts          time.Time `db:"ts" json:"ts"`
}

// Initialize postion table
func initPosition() (isFinish bool, err error) {
	// Step 1: Check if a record exists for the default position
	sqlStr := "select count(id) as rownum from position where id=10000"
	hasRecord, isFinish, err := genericCheckRecord("position", sqlStr)
	// Step 2: Exit if the record exists or an error occurs
	if hasRecord || !isFinish || err != nil {
		return
	}
	// Step 3: Insert a record for the system default positon "Default position" into the position table.
	sqlStr = `insert into position(id,name,description,creatorid) 
	values(10000,'Default position','System pre-set position',10000)`
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initPosition insert initvalue failed", zap.Error(err))
		return isFinish, err
	}
	return
}
