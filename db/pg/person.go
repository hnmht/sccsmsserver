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
	SystemFlag   int16     `db:"systemflag" json:"systemflag"`
	Status       int16     `db:"status" json:"status"`
	CreateDate   time.Time `db:"createtime" json:"createDate"`
	Ts           time.Time `db:"ts" json:"ts"`
	Dr           int16     `db:"dr" json:"dr"`
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
