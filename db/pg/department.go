package pg

import (
	"database/sql"
	"encoding/json"
	"sccsmsserver/cache"
	"sccsmsserver/i18n"
	"sccsmsserver/pub"
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

// Simplify Department Struct.
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

// Initialize Department table.
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

// Get Simplify department information by ID.
func (d *SimpDept) GetSimpDeptInfoByID() (resStatus i18n.ResKey, err error) {
	// Get simplify information from cache
	number, sdb, _ := cache.Get(pub.SimpDept, d.ID)
	if number > 0 {
		json.Unmarshal(sdb, &d)
		resStatus = i18n.StatusOK
		return
	}
	// If Simplify information isn't in cache, retrieve it from database
	sqlStr := `select code,name,leader,description,status,ts 
	from department where id=$1`
	err = db.QueryRow(sqlStr, d.ID).Scan(&d.Code, &d.Name, &d.Leader.ID, &d.Description, &d.Status, &d.Ts)
	if err != nil {
		if err == sql.ErrNoRows {
			resStatus = i18n.StatusDeptNotExist
			return
		}
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetSimpDeptInfoById failed", zap.Error(err))
		return
	}
	// Get Department Leader's information
	if d.Leader.ID > 0 {
		resStatus, err = d.Leader.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}

	// Write into cache
	jsonB, _ := json.Marshal(d)
	cache.Set(pub.SimpDept, d.ID, jsonB)

	return i18n.StatusOK, nil
}
