package pg

import (
	"sccsmsserver/pub"

	"go.uber.org/zap"
)

// Initialize the logicmsg table
func initLogicMsg() (isFinish bool, err error) {
	isFinish = true
	//Step 1: Check if data already exists in the logicmsg table
	sqlStr := "select count(id) as rownum from logicmsg where dr=0"
	hasRecord, isFinish, err := genericCheckRecord("logicmsg", sqlStr)
	if hasRecord || !isFinish || err != nil {
		return
	}
	//Step 2: Insert system preset data into the logicmsg table
	insertSql := `insert into logicmsg(code,content) values($1,$2)`
	stmt, err := db.Prepare(insertSql)
	if err != nil {
		isFinish = false
		zap.L().Error("initLogicMsg db.Prepare failed:", zap.Error(err))
		return
	}
	defer stmt.Close()
	// Extract system prompt messages from system constants
	for code, content := range pub.ResStatusCodeMsg {
		_, err = stmt.Exec(code, content)
		if err != nil {
			isFinish = false
			zap.L().Error("initLogicMsg stmt.Exec failed:", zap.Error(err))
			return
		}
	}
	return
}

// Initialize the logicmsg_t table
func initLogicMsgTranslate() (isFinish bool, err error) {
	isFinish = true
	// Step 1: Check if data already exists in the logicmsg_t table
	sqlStr := "select count(id) as rownum from logicmsg_t where dr=0"
	hasRecord, isFinish, err := genericCheckRecord("logicmsg_t", sqlStr)
	if hasRecord || !isFinish || err != nil {
		return
	}
	// Step 2: Insert system preset data into the logicmsg_t table
	insertSql := `insert into logicmsg_t(code,defaultcontent,language,content) 
	values($1,$2,$3,$4)`
	stmt, err := db.Prepare(insertSql)
	if err != nil {
		isFinish = false
		zap.L().Error("initLogicMsgTranslate db.Prepare failed:", zap.Error(err))
		return
	}
	defer stmt.Close()
	// Extract system prompt messages from system constants
	for code, content := range pub.ResStatusCodeMsg {
		_, err = stmt.Exec(code, content, pub.DefaultLocale.Language, content)
		if err != nil {
			isFinish = false
			zap.L().Error("initLogicMsgTranslate stmt.Exec failed:", zap.Error(err))
			return
		}
	}
	return
}
