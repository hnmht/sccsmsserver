package pg

import (
	"encoding/json"
	"sccsmsserver/cache"
	"sccsmsserver/i18n"
	"sccsmsserver/pub"
	"time"

	"go.uber.org/zap"
)

// Position Master Data
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

// Get Position information by ID
func (p *Position) GetInfoByID() (resStatus i18n.ResKey, err error) {
	// Get Postion information from cache
	number, b, _ := cache.Get(pub.Position, p.ID)
	if number > 0 {
		json.Unmarshal(b, &p)
		resStatus = i18n.StatusOK
		return
	}
	// If Position information isn't in cache, retrieve it from database.
	sqlStr := `select name,description,status,createtime,creatorid,
	modifytime,modifierid,ts,dr
	from position 
	where id = $1`
	err = db.QueryRow(sqlStr, p.ID).Scan(&p.Name, &p.Description, &p.Status, &p.CreateDate, &p.Creator.ID,
		&p.ModifyDate, &p.Modifier.ID, &p.Ts, &p.Dr)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("Position.GetInfoByID db.QueryRow failed", zap.Error(err))
		return
	}
	// Get creator information.
	if p.Creator.ID > 0 {
		resStatus, err = p.Creator.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Get Modifier information.
	if p.Modifier.ID > 0 {
		resStatus, err = p.Modifier.GetPersonInfoByID()
		if resStatus != i18n.StatusOK || err != nil {
			return
		}
	}
	// Write into cache
	pB, _ := json.Marshal(p)
	cache.Set(pub.Position, p.ID, pB)

	return i18n.StatusOK, nil
}
