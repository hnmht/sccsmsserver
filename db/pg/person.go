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

// Person Master Data (simplify User)
type Person struct {
	ID           int32     `db:"id" json:"id"`
	Code         string    `db:"code" json:"code"`
	Name         string    `db:"name" json:"name"`
	Avatar       File      `db:"fileid" json:"avatar"`
	DeptID       int32     `db:"deptid" json:"deptID"`
	DeptCode     string    `json:"deptCode"`
	DeptName     string    `json:"deptName"`
	IsOperator   int16     `json:"isOperator"`
	PositionID   int32     `db:"positionid" json:"positionID"`
	PositionName string    `json:"positionName"`
	Description  string    `db:"description" json:"description"`
	Mobile       string    `db:"mobile" json:"mobile"`
	Email        string    `db:"email" json:"email"`
	Gender       int16     `db:"gender" json:"gender"`
	SystemFlag   int16     `db:"systemflag" json:"systemFlag"`
	Status       int16     `db:"status" json:"status"`
	CreateDate   time.Time `db:"createtime" json:"createDate"`
	Ts           time.Time `db:"ts" json:"ts"`
	Dr           int16     `db:"dr" json:"dr"`
}

// Latest Person Master Data
type PersonCache struct {
	QueryTs      time.Time `json:"queryTs"`
	ResultNumber int32     `json:"resultNumber"`
	DelItems     []Person  `json:"delItems"`
	UpdateItems  []Person  `json:"updateItems"`
	NewItems     []Person  `json:"newItems"`
	ResultTs     time.Time `json:"resultTs"`
}

// Get Person information by User ID.
func (p *Person) GetPersonInfoByID() (resStatus i18n.ResKey, err error) {
	// Get Person information from local cache
	e, pb, _ := cache.Get(pub.Person, p.ID)
	if e > 0 {
		err = json.Unmarshal(pb, &p)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetPersonInfoByID json.Unmarshal failed", zap.Error(err))
			return
		}
		resStatus = i18n.StatusOK
		return
	}
	// If Person information isn't in the cache, retrieve it from the database.
	sqlStr := `select u.code as usercode,
	u.name as username,
	COALESCE(u.fileid,'0') as fileid,
	u.deptid as dept_id,
	COALESCE(d.code,'') as deptcode,
	COALESCE(d.name,'') as deptname,
	u.isoperator as isoperator,
	u.positionid as positonid,
	COALESCE(p.name,'') as positionname,
	COALESCE(u.description,'') as description,
	u.mobile as mobile,
	u.email as email,
	COALESCE(u.gender,0) as gender,
	u.systemflag as systemflag,
	u.status as status,
	u.createtime as createtime,
	u.ts as ts,
	u.dr as dr 	
	from sysuser as u
	left join department as d on u.deptid=d.id
	left join position as p on u.positionid=p.id
	where  u.id=$1`
	err = db.QueryRow(sqlStr, p.ID).Scan(&p.Code, &p.Name, &p.Avatar.ID, &p.DeptID, &p.DeptCode,
		&p.DeptName, &p.IsOperator, &p.PositionID, &p.PositionName, &p.Description,
		&p.Mobile, &p.Email, &p.Gender, &p.SystemFlag, &p.Status,
		&p.CreateDate, &p.Ts, &p.Dr)
	if err != nil && err != sql.ErrNoRows {
		resStatus = i18n.StatusInternalError
		zap.L().Error("dap.GetPersonInfoByID failed", zap.Error(err))
		return
	}
	if err == sql.ErrNoRows {
		resStatus = i18n.StatusUserNotExist
		return
	}

	// Get User Avatar deatil
	if p.Avatar.ID > 0 {
		resStatus, err = p.Avatar.GetFileInfoByID()
		if err != nil || resStatus != i18n.StatusOK {
			return
		}
	}
	// Write into cache
	jsonP, _ := json.Marshal(p)
	cache.Set(pub.Person, p.ID, jsonP)

	return i18n.StatusOK, nil
}

// Get Person List
func GetPersons() (persons []Person, resStatus i18n.ResKey, err error) {
	persons = make([]Person, 0)
	resStatus = i18n.StatusOK
	// Retrieve from sysuser table
	sqlStr := `select id from sysuser where dr=0 order by ts desc`
	rows, err := db.Query(sqlStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		resStatus = i18n.StatusInternalError
		zap.L().Error("Get Persons from database failed", zap.Error(err))
		return
	}
	defer rows.Close()
	// Extract data from database query results
	for rows.Next() {
		var p Person
		err = rows.Scan(&p.ID)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetPerson from rows scan person info failed", zap.Error(err))
			return
		}
		resStatus, err = p.GetPersonInfoByID()
		if err != nil || resStatus != i18n.StatusOK {
			return
		}
		persons = append(persons, p)
	}

	return
}

// Get Latest Person Master data
func (pc *PersonCache) GetLatestPersons() (resStatus i18n.ResKey, err error) {
	pc.DelItems = make([]Person, 0)
	pc.NewItems = make([]Person, 0)
	pc.UpdateItems = make([]Person, 0)
	resStatus = i18n.StatusOK
	// Get the latest timestamp from sysuser table
	sqlStr := "select ts from sysuser where ts > $1 order by ts desc limit(1)"
	err = db.QueryRow(sqlStr, pc.QueryTs).Scan(&pc.ResultTs)
	if err != nil {
		if err == sql.ErrNoRows {
			pc.ResultNumber = 0
			pc.ResultTs = pc.QueryTs
			return
		}
		resStatus = i18n.StatusInternalError
		zap.L().Error("PersonCache.GetLatestPersons  db.QueryRow failed:", zap.Error(err))
		return
	}
	// Retrieve all data greater than the latest timestamp.
	sqlStr = `select a.id,
	a.code,
	a.name,
	a.fileid,
	a.deptid,
	COALESCE((select b.code from department b where b.id = a.deptid),'') as deptcode,
	COALESCE((select b.name from department b where b.id = a.deptid),'') as deptname,
	a.isoperator as isoperator,
	a.positionid as positionid,
	COALESCE((select o.name from position o where o.id = a.positionid),'') as positionname,
	COALESCE(a.description,'') as description,
	COALESCE(a.mobile,'') as mobile,
	COALESCE(a.email,'') as email,
	a.gender,
	a.systemflag,
	a.status,
	a.createtime,
	a.ts,
	a.dr 
	from sysuser a
	where a.ts > $1 and a.systemflag=0 order by a.ts desc`
	rows, err := db.Query(sqlStr, pc.QueryTs)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("PersonCache.GetLatestPersons()  db.Query() failed", zap.Error(err))
		return
	}
	defer rows.Close()

	// Extract data from database query results
	for rows.Next() {
		var p Person
		err = rows.Scan(&p.ID, &p.Code, &p.Name, &p.Avatar.ID, &p.DeptID,
			&p.DeptCode, &p.DeptName, &p.IsOperator, &p.PositionID, &p.PositionName,
			&p.Description, &p.Mobile, &p.Email, &p.Gender, &p.SystemFlag,
			&p.Status, &p.CreateDate, &p.Ts, &p.Dr)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("PersonCache.GetLatestPersons rows.scan failed", zap.Error(err))
			return
		}
		// Get Avatar detail
		if p.Avatar.ID > 0 {
			resStatus, err = p.Avatar.GetFileInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}

		// Determine the latest data category.
		if p.Dr == 0 { // The Dr field being equal 0 means the data has not been deleted
			// If the data's createDate field is less than or equal to the QueryTs.
			if p.CreateDate.Before(pc.QueryTs) || p.CreateDate.Equal(pc.QueryTs) { // It means the client needs to update the data
				pc.ResultNumber++
				pc.UpdateItems = append(pc.UpdateItems, p)
			} else { // It means the client needs to add the data
				pc.ResultNumber++
				pc.NewItems = append(pc.NewItems, p)
			}
		} else { // The Dr field being not equal 0 means the data has been deleted
			// If the data's createDate field is less than or equal to the QueryTs
			if p.CreateDate.Before(pc.QueryTs) || p.CreateDate.Equal(pc.QueryTs) { // It means the client needs to delete the data
				pc.ResultNumber++
				pc.DelItems = append(pc.DelItems, p)
			}
		}
	}
	return
}
