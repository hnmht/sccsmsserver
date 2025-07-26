package pg

import (
	"time"

	"go.uber.org/zap"
)

// Department Master Data
type Department struct {
	ID          int32     `db:"id" json:"id"`
	Code        string    `db:"code" json:"code"`
	Name        string    `db:"name" json:"name"`
	FatherID    SimpDept  `db:"deptparent" json:"fatherid"` //上级部门
	Leader      Person    `db:"leader" json:"leader"`
	Description string    `db:"description" json:"description"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createdate"`
	Creator     Person    `db:"creatorid" json:"creator"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	ModifyDate  time.Time `db:"modify_time" json:"modifydate"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"` //删除标志
}

// Simplify Department Struct
type SimpDept struct {
	ID          int32     `db:"id" json:"id"`
	Code        string    `db:"deptcode" json:"code"`
	Name        string    `db:"deptname" json:"name"`
	FatherID    int32     `db:"fatherid" json:"fatherid"`
	Leader      Person    `db:"leader" json:"leader"`
	Description string    `db:"description" json:"description"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"create_time" json:"createdate"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"` //删除标志
}

// Initialize Department table
func initDepartment() (isFinish bool, err error) {
	// Step 1: Check if a record exists for the default department "Default Department" in the department table
	sqlStr := "select count(id) as rownum from department where id=10000"
	hasRecord, isFinish, err := genericCheckRecord("department", sqlStr)
	// Step 2: Exit if the record exists or an error occurs
	if hasRecord || !isFinish || err != nil {
		return
	}
	// Step 3: Insert a record for the system default department "Default Department" into the department table.
	sqlStr = `insert into department(id,code,name,description,creatorid) 
	values(10000,'default','Default Department','System pre-set departments',10000)`
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initDepartment insert initvalue failed", zap.Error(err))
		return isFinish, err
	}
	return
}
