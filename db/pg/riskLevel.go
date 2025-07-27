package pg

import (
	"time"

	"go.uber.org/zap"
)

type RiskLevel struct {
	ID          int32     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Color       string    `db:"color" json:"color"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Createor    Person    `db:"creatorid" json:"creator"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	ModifyDate  time.Time `db:"modifytime" json:"modifyDate"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

func initRiskLevel() (isFinish bool, err error) {
	// Step 1: Check if a record exists for the default Risk Level.
	sqlStr := "select count(id) as rownum from risklevel where dr=0"
	// Step 2: Exit if the record exists or an error occurs.
	hasRecord, isFinish, err := genericCheckRecord("risklevel", sqlStr)
	if hasRecord || !isFinish || err != nil {
		return
	}
	// Step 3: Insert default Risk Level records.
	sqlStrs := []string{
		"insert into risklevel(id,name,description,color,creatorid) values(1,'Major Risk','System pre-set','red',10000)",
		"insert into risklevel(id,name,description,color,creatorid) values(2,'Significant Risk','System pre-set','orange',10000)",
		"insert into risklevel(id,name,description,color,creatorid) values(3,'General Risk','System pre-set','yellow',10000)",
		"insert into risklevel(id,name,description,color,creatorid) values(4,'Low Risk','System pre-set','blue',10000)",
		"insert into risklevel(id,name,description,color,creatorid) values(5,'No Risk','System pre-set','white',10000)",
	}

	for _, t := range sqlStrs {
		_, err = db.Exec(t)
		if err != nil {
			isFinish = false
			zap.L().Error("initRiskLevel insert default data:"+t+" failed.", zap.Error(err))
			return
		}
	}
	return
}
