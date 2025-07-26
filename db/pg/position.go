package pg

import "time"

// Position 岗位
type Position struct {
	ID          int32     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"createtime" json:"createdate"`
	Creator     Person    `db:"creatorid" json:"createuser"`
	ModifyDate  time.Time `db:"modify_time" json:"modifydate"`
	Modifier    Person    `db:"modifierid" json:"modifyuser"`
	Dr          int16     `db:"dr" json:"dr"`
	Ts          time.Time `db:"ts" json:"ts"`
}

// 初始化岗位档案
func initPosition() (isFinish bool, err error) {
	/* //检查岗位档案表中是否存在记录
	sqlStr := "select count(id) as rownum from operatingpost where id=10000"
	hasRecord, isFinish, err := checkRecord("operatingpost", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//如果表中没有数据则插入预置数据
	sqlStr = "insert into operatingpost(id,name,description,createuserid) values(10000,'预置岗位','系统预置岗位',10000)"
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initOperatingPost insert initvalue failed", zap.Error(err))
		return isFinish, err
	} */
	return
}
