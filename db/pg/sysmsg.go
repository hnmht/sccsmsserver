package pg

import (
	"sccsmsserver/pub"

	"go.uber.org/zap"
)

// Initialize the sysmsg table
func initSysMsg() (isFinish bool, err error) {
	isFinish = true
	// Step 1: Check if data already exists in the sysmsg table
	sqlStr := "select count(id) as rownum from sysmsg where dr=0"
	hasRecord, isFinish, err := genericCheckRecord("sysmsg", sqlStr)
	if hasRecord || !isFinish || err != nil {
		return
	}

	// Setp 2: Insert system preset data into the sysmsg table
	insertSql := "insert into sysmsg(code,content) values($1,$2)"
	stmt, err := db.Prepare(insertSql)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysMsg db.Prepare failed:", zap.Error(err))
		return
	}
	defer stmt.Close()
	// Extract system prompt messages from system constants
	for code, content := range pub.ResCodeMsg {
		_, err = stmt.Exec(code, content)
		if err != nil {
			isFinish = false
			zap.L().Error("initSysMsg stmt.Exec failed:", zap.Error(err))
			return
		}
	}
	return
}

// Initialize the sysmsg_t Table
func initSysMsgTranslate() (isFinish bool, err error) {
	isFinish = true
	// Step 1: Check if data already exists in the sysmsg_t table
	sqlStr := "select count(id) as rownum from sysmsg_t where dr=0"
	hasRecord, isFinish, err := genericCheckRecord("sysmsg_t", sqlStr)
	if hasRecord || !isFinish || err != nil {
		return
	}

	// Step 2: Insert system preset data into the sysmsg_t table
	insertSql := `insert into sysmsg_t(code,defaultcontent,language,content) 
	values($1,$2,$3,$4)`
	stmt, err := db.Prepare(insertSql)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysMsgTranslate db.Prepare failed:", zap.Error(err))
		return
	}
	defer stmt.Close()
	// Extract system prompt messages from system constants
	for code, content := range pub.ResCodeMsg {
		_, err = stmt.Exec(code, content, pub.DefaultLocale.Language, content)
		if err != nil {
			isFinish = false
			zap.L().Error("initSysMsgTranslate stmt.Exec failed:", zap.Error(err))
			return
		}
	}
	return
}
