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

// Department Master Data.
type Department struct {
	ID          int32     `db:"id" json:"id"`
	Code        string    `db:"code" json:"code"`
	Name        string    `db:"name" json:"name"`
	FatherID    SimpDept  `db:"deptparent" json:"fatherID"`
	Leader      Person    `db:"leader" json:"leader"`
	Description string    `db:"description" json:"description"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Creator     Person    `db:"creatorid" json:"creator"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	ModifyDate  time.Time `db:"modifytime" json:"modifyDate"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

// Simplify Department Struct.
type SimpDept struct {
	ID          int32     `db:"id" json:"id"`
	Code        string    `db:"code" json:"code"`
	Name        string    `db:"name" json:"name"`
	FatherID    int32     `db:"fatherid" json:"fatherID"`
	Leader      Person    `db:"leader" json:"leader"`
	Description string    `db:"description" json:"description"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createDate"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

// Simplify Department front end Cache struct
type SimpDeptCache struct {
	QueryTs      time.Time  `json:"queryts"`
	ResultNumber int32      `json:"resultnum"`
	DelDepts     []SimpDept `json:"delitems"`
	UpdateDepts  []SimpDept `json:"updateitems"`
	NewDepts     []SimpDept `json:"newitems"`
	ResultTs     time.Time  `json:"resultts"`
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

// Get department list
func GetDepts() (depts []Department, resStatus i18n.ResKey, err error) {
	depts = make([]Department, 0)
	resStatus = i18n.StatusOK
	// Get Data from department table
	sqlStr := `select a.id,a.code,a.name,a.fatherid,a.leader,
	a.description,a.status,	a.createtime,a.creatorid,COALESCE(a.modifytime,createtime) as modifytime,
	COALESCE(a.modifierid,0) as modifierid,a.ts 
	from department a
	where a.dr=0 
	order by a.ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetDepts db.Query failed", zap.Error(err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var dept Department
		err = rows.Scan(&dept.ID, &dept.Code, &dept.Name, &dept.FatherID.ID, &dept.Leader.ID,
			&dept.Description, &dept.Status, &dept.CreateDate, &dept.Creator.ID, &dept.ModifyDate,
			&dept.Modifier.ID, &dept.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetDepts rows.next failed", zap.Error(err))
			return
		}
		// Get Parent department detail
		if dept.FatherID.ID > 0 {
			resStatus, err = dept.FatherID.GetSimpDeptInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Leader detail
		if dept.Leader.ID > 0 {
			resStatus, err = dept.Leader.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Creator deatail
		if dept.Creator.ID > 0 {
			resStatus, err = dept.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get modifier detail
		if dept.Modifier.ID > 0 {
			resStatus, err = dept.Modifier.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		depts = append(depts, dept)
	}

	resStatus = i18n.StatusOK
	return
}

// Get simplify department list
func GetSimpDepts() (depts []SimpDept, resStatus i18n.ResKey, err error) {
	depts = make([]SimpDept, 0)
	resStatus = i18n.StatusOK
	// Get data from department list
	sqlStr := `select id,code,name,fatherid,leader,
	description,status,createtime,ts,dr 
	from department 
	where dr=0 
	order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetSimpDepts db.Query failed", zap.Error(err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var dept SimpDept
		err = rows.Scan(&dept.ID, &dept.Code, &dept.Name, &dept.FatherID, &dept.Leader.ID,
			&dept.Description, &dept.Status, &dept.CreateDate, &dept.Ts, &dept.Dr)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetSimpDepts rows.next failed", zap.Error(err))
			return
		}
		// Get department leader details
		if dept.Leader.ID > 0 {
			resStatus, err = dept.Leader.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		depts = append(depts, dept)
	}
	return
}

// Get latest Department Master Data
func (dc *SimpDeptCache) GetLatestSimpDepts() (resStatus i18n.ResKey, err error) {
	dc.DelDepts = make([]SimpDept, 0)
	dc.NewDepts = make([]SimpDept, 0)
	dc.UpdateDepts = make([]SimpDept, 0)
	resStatus = i18n.StatusOK
	// Get the lastest timestamp from department table
	sqlStr := "select ts from department where ts > $1 order by ts desc limit(1)"
	err = db.QueryRow(sqlStr, dc.QueryTs).Scan(&dc.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			dc.ResultNumber = 0
			dc.ResultTs = dc.QueryTs
			resStatus = i18n.StatusOK
			return
		}
		resStatus = i18n.StatusInternalError
		zap.L().Error("SimpDeptCache.GetLatestSimpDepts get latest ts db.QueryRow failed", zap.Error(err))
		return
	}

	// Retrieve all data greater than the  latest timestamp.
	sqlStr = `select id,code,name,fatherid,leader,
	description,status,createtime,ts,dr 
	from department 
	where ts > $1 order by ts desc`
	rows, err := db.Query(sqlStr, dc.QueryTs)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("SimpDeptCache.GetLatestSimpDepts db.Qeury failed", zap.Error(err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var dept SimpDept
		err = rows.Scan(&dept.ID, &dept.Code, &dept.Name, &dept.FatherID, &dept.Leader.ID,
			&dept.Description, &dept.Status, &dept.CreateDate, &dept.Ts, &dept.Dr)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("SimpDeptCache.GetLatestSimpDepts rows.next failed", zap.Error(err))
			return
		}
		// Get department leader details.
		if dept.Leader.ID > 0 {
			resStatus, err = dept.Leader.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		if dept.Dr == 0 {
			if dept.CreateDate.Before(dc.QueryTs) || dept.CreateDate.Equal(dc.QueryTs) {
				dc.ResultNumber++
				dc.UpdateDepts = append(dc.UpdateDepts, dept)
			} else {
				dc.ResultNumber++
				dc.NewDepts = append(dc.NewDepts, dept)
			}
		} else {
			if dept.CreateDate.Before(dc.QueryTs) || dept.CreateDate.Equal(dc.QueryTs) {
				dc.ResultNumber++
				dc.DelDepts = append(dc.DelDepts, dept)
			}
		}
	}
	return
}

// Get Department information by ID
func (d *Department) GetDeptInfoByID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	sqlStr := `select code,name,leader,description,status,
	ts from department where id=$1`
	err = db.QueryRow(sqlStr, d.ID).Scan(&d.Code, &d.Name, &d.Leader.ID, &d.Description, &d.Status,
		&d.Ts)
	if err != nil && err != sql.ErrNoRows {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Department.GetDeptInfoById db.QueryRow failed", zap.Error(err))
		return
	}

	// Get department leader details.
	if d.Leader.ID > 0 {
		resStatus, err = d.Leader.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	return
}

// Check if the department code exists.
func (d *Department) CheckDeptCodeExist() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	var count int32
	sqlStr := "select count(id) from department where dr=0 and code = $1 and id <> $2"
	err = db.QueryRow(sqlStr, d.Code, d.ID).Scan(&count)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("CheckDeptCodeExist query failed", zap.Error(err))
		return
	}
	if count > 0 {
		resStatus = i18n.StatusDeptCodeExist
		return
	}
	return
}

// Add department
func (d *Department) AddDept() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the department code exists.
	resStatus, err = d.CheckDeptCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Write a record to the department table
	sqlStr := `insert into department(code,name,fatherid,leader,description,
		status,createtime,creatorid) 
		values($1,$2,$3,$4,$5,$6,now(),$7) 
		returning id`
	err = db.QueryRow(sqlStr, d.Code, d.Name, d.FatherID.ID, d.Leader.ID, d.Description, d.Status, d.Creator.ID).Scan(&d.ID)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("department.AddDept stmt.QueryRow failed", zap.Error(err))
		return
	}

	return
}

// Edit Department
func (dept *Department) Edit() (resStatus i18n.ResKey, err error) {
	// Check if the parent department is compliant.
	resStatus, err = dept.CheckFather()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Check if the department code exists.
	resStatus, err = dept.CheckDeptCodeExist()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}

	// Update the record in the department table
	sqlStr := `update department set code=$1,name=$2,fatherid=$3,leader=$4,description=$5,
	status=$6,modifierid=$7,modifytime=now(),ts=current_timestamp where id = $8 and ts = $9`
	res, err := db.Exec(sqlStr, dept.Code, dept.Name, dept.FatherID.ID, dept.Leader.ID, dept.Description,
		dept.Status, dept.Modifier.ID, dept.ID, dept.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Department.Edit stmt.exec failed", zap.Error(err))
		return
	}
	updateNumber, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Department.Edit res.RowsAffected falied", zap.Error(err))
		return
	}

	if updateNumber < 1 {
		resStatus = i18n.StatusOtherEdit
		return
	}

	// Delete from local cache
	dept.DelFromLocalCache()

	return
}

// Check if the parent department is compliant.
func (dept *Department) CheckFather() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK

	if dept.FatherID.ID > 0 {
		// The parent department cannot be itself
		if dept.ID == dept.FatherID.ID {
			resStatus = i18n.StatusDeptFatherSelf
			return
		}

		// The parent department cannot be in a circular depandency.
		depts, resStatus, err := GetSimpDepts()
		if resStatus != i18n.StatusOK || err != nil {
			return resStatus, err
		}
		// Get all subordinate departments of this department,
		// The parent department cannot be any of them.
		childrens := FindSimpDeptChildrens(depts, dept.ID)
		var rowNum int32
		for _, child := range childrens {
			if child.ID == dept.FatherID.ID {
				rowNum++
			}
		}
		if rowNum > 0 {
			resStatus = i18n.StatusDeptFatherCircle
			return resStatus, nil
		}
	}

	return
}

// Find all child departments by ID.
func FindSimpDeptChildrens(arr []SimpDept, id int32) []SimpDept {
	childrens := make([]SimpDept, 0)
	for _, dept := range arr {
		if dept.FatherID == id {
			childrens = append(childrens, dept)
			child := FindSimpDeptChildrens(arr, dept.ID)
			childrens = append(childrens, child...)
		}
	}
	return childrens
}

// Delete department
func (dept *Department) Delete() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Check if the department ID redrenced.
	resStatus, err = dept.CheckIsUsed()
	if resStatus != i18n.StatusOK || err != nil {
		return
	}
	// Update the department table with  a deletion flag.
	sqlStr := `update department set dr=1,modifytime=now(),modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	res, err := db.Exec(sqlStr, dept.Modifier.ID, dept.ID, dept.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Department.DelDept stmt.exec failed", zap.Error(err))
		return
	}
	// Check the number of rows affected by the SQL update operation.
	affected, err := res.RowsAffected()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Department.DelDept check RowsAffected failed", zap.Error(err))
		return
	}
	if affected < 1 {
		resStatus = i18n.StatusOtherEdit
		zap.L().Info("Department.DelDept other edit")
		return
	}

	// Delete from local cache
	dept.DelFromLocalCache()

	return
}

// Batch delete department
func DeleteDepts(depts *[]Department, modifyUserid int32) (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Begin a database trnasaction.
	tx, err := db.Begin()
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteDepts db.begin failed", zap.Error(err))
		return
	}
	defer tx.Commit()
	// Pre-processing for update the department table deletion flag.
	delSqlStr := `update department set dr=1,modifytime=now(),modifierid=$1,ts=current_timestamp 
	where id=$2 and dr=0 and ts=$3`
	stmt, err := tx.Prepare(delSqlStr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("DeleteDepts Delete prepare failed", zap.Error(err))
		_ = tx.Rollback()
		return
	}
	defer stmt.Close()

	for _, dept := range *depts {
		// Check if the department ID refrenced.
		resStatus, err = dept.CheckIsUsed()
		if resStatus != i18n.StatusOK || err != nil {
			_ = tx.Rollback()
			return
		}

		// Update the department table with a deletction flag.
		res, err := stmt.Exec(modifyUserid, dept.ID, dept.Ts)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("DeleteDepts stmt.exec failed", zap.Error(err))
			_ = tx.Rollback()
			return resStatus, err
		}
		// Check the row number of effected by the SQL operation.
		affected, err := res.RowsAffected()
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("DeleteDepts check RowsAffected failed", zap.Error(err))
			_ = tx.Rollback()
			return resStatus, err
		}
		if affected < 1 {
			resStatus = i18n.StatusOtherEdit
			zap.L().Info("DeleteDepts other edit")
			_ = tx.Rollback()
			return resStatus, err
		}
		// Delete from local cache
		dept.DelFromLocalCache()
	}
	return
}

// Delete department from local cache
func (d *Department) DelFromLocalCache() {
	number, _, _ := cache.Get(pub.SimpDept, d.ID)
	if number > 0 {
		cache.Del(pub.SimpDept, d.ID)
	}
}

// Check if the department ID is refrenced.
func (dept *Department) CheckIsUsed() (resStatus i18n.ResKey, err error) {
	// Define the items to be checked.
	checkItems := []ArchiveCheckUsed{
		{
			Description:    "There are subordinate departments",
			SqlStr:         `select count(id) as usednum from department where fatherid=$1 and dr=0`,
			UsedReturnCode: i18n.StatusDeptLowLevelExist,
		},
		{
			Description:    "Refrenced by the user master data",
			SqlStr:         `select count(id) as usernum from sysuser where deptid = $1 and dr=0`,
			UsedReturnCode: i18n.StatusUserUsed,
		},
		/* {
			Description:    "被执行项目默认值引用",
			SqlStr:         `select count(id) as usednum from exectiveitem where resulttypeid = '520' and dr=0 and defaultvalue=cast($1 as varchar)`,
			UsedReturnCode: i18n.StatusEPDefaultUsed,
		},
		{
			Description:    "被执行项目错误值引用",
			SqlStr:         `select count(id) as usednum from exectiveitem where resulttypeid = '520' and dr=0 and errorvalue=cast($1 as varchar)`,
			UsedReturnCode: i18n.StatusEPErrorUsed,
		},
		{
			Description:    "被执行模板默认值引用",
			SqlStr:         `select count(id) from exectivetemplate_b where eid_id in (select id from exectiveitem where resulttypeid='520' and dr=0) and dr=0 and defaultvalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEPTDefaultUsed,
		},
		{
			Description:    "被执行模板错误值引用",
			SqlStr:         `select count(id) from exectivetemplate_b where eid_id in (select id from exectiveitem where resulttypeid='520' and dr=0) and dr=0 and errorvalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEPTErrorUsed,
		},
		{
			Description:    "被现场档案所属部门引用",
			SqlStr:         `select count(id) from sceneitem where subdept_id=$1 and dr=0`,
			UsedReturnCode: i18n.StatusCSUsed,
		},
		{
			Description:    "被现场档案责任部门引用",
			SqlStr:         `select count(id) from sceneitem where respdept_id=$1 and dr=0`,
			UsedReturnCode: i18n.StatusCSUsed,
		},
		{
			Description:    "被指令单引用",
			SqlStr:         `select count(id) from workorder_h where dept_id=$1 and dr=0`,
			UsedReturnCode: i18n.StatusWOUsed,
		},
		{
			Description:    "被执行单引用",
			SqlStr:         `select count(id) from executedoc_h where dept_id=$1 and dr=0`,
			UsedReturnCode: i18n.StatusEOUsed,
		},
		{
			Description:    "被执行单执行值引用",
			SqlStr:         `select count(id) from executedoc_b where eid_id in (select id from exectiveitem where resulttypeid='520' and dr=0) and dr=0 and exectivevalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEOValueUsed,
		},
		{
			Description:    "被执行单错误值引用",
			SqlStr:         `select count(id) from executedoc_b where eid_id in (select id from exectiveitem where resulttypeid='520' and dr=0) and dr=0 and errorvalue=CAST($1 as varchar)`,
			UsedReturnCode: i18n.StatusEOErrorUsed,
		},
		{
			Description:    "被问题处理单引用",
			SqlStr:         `select count(id) from disposedoc where dept_id=$1 and dr=0`,
			UsedReturnCode: i18n.StatusIRFUsed,
		},

		{
			Description:    "被培训记录表头部门引用",
			SqlStr:         `select count(id) from trainrecord_h where dept_id=$1 and dr=0`,
			UsedReturnCode: i18n.StatusTRDeptUsed,
		},
		{
			Description:    "被劳保用品发放单发放部门引用",
			SqlStr:         `select count(id) from lpaissuedoc_h where dept_id=$1 and dr=0`,
			UsedReturnCode: i18n.StatusPPEIFDeptUsed,
		}, */
	}

	// Check item by item
	var usedNum int32
	for _, item := range checkItems {
		err = db.QueryRow(item.SqlStr, dept.ID).Scan(&usedNum)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("UserDefineDoc.CheckIsUsed "+item.Description+"failed", zap.Error(err))
			return
		}
		if usedNum > 0 {
			resStatus = item.UsedReturnCode
			return
		}
	}
	return
}
