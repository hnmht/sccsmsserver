package pg

import (
	"sccsmsserver/pub"

	"go.uber.org/zap"
)

// Multilingual Logic Messages
var MultilingLogicMsg map[string]map[pub.ResStatus]string

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
		_, err = stmt.Exec(code, content, DefaultLocale.Language, content)
		if err != nil {
			isFinish = false
			zap.L().Error("initLogicMsgTranslate stmt.Exec failed:", zap.Error(err))
			return
		}
	}
	return
}

// Load Multilingual Logic Messages
func LoadMultilingualLogicMsg() (isFinish bool, err error) {
	isFinish = true
	// Step 1: Initialize map values
	MultilingLogicMsg = make(map[string]map[pub.ResStatus]string, 0)
	// Step 2: Check if the SysLocaleList variable content is empty.
	if len(SysLocaleList) < 1 {
		isFinish = false
		zap.L().Info("LoadMultilingualLogicMsg SysLocaleList is empty.")
		return
	}
	// Step 3: Preparing to read data from the logicmsg_t table in the database.
	sqlStr := `select code,content from logicmsg_t where language=$1 and dr=0`
	// Step 4: Get language from SysLocaleList
	for _, v := range SysLocaleList {
		lang := v.Language
		// Get data from logicmsg_t table
		rows, err := db.Query(sqlStr, lang)
		if err != nil {
			isFinish = false
			zap.L().Error("LoadMultilingualLogicMsg db.Query failed:", zap.Error(err))
			return isFinish, err
		}
		langMsgMap := make(map[pub.ResStatus]string)
		for rows.Next() {
			var code pub.ResStatus
			var content string
			err = rows.Scan(&code, &content)
			if err != nil {
				isFinish = false
				zap.L().Error("LoadMultilingualLogicMsg rows.Next() failed:", zap.Error(err))
				return isFinish, err
			}
			langMsgMap[code] = content
		}
		// Assign to the MultilingLogicMsg variable
		MultilingLogicMsg[lang] = langMsgMap
		// Close Rows
		rows.Close()
	}
	zap.L().Info("Multilingual logic messages loaded successfully.")
	return
}
