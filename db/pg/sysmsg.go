package pg

import (
	"sccsmsserver/pub"

	"go.uber.org/zap"
)

// multilingual system messages
var MultilingualSysMsg map[string]map[pub.ResCode]string

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
		_, err = stmt.Exec(code, content, DefaultLocale.Language, content)
		if err != nil {
			isFinish = false
			zap.L().Error("initSysMsgTranslate stmt.Exec failed:", zap.Error(err))
			return
		}
	}
	return
}

// Load multilingual system messages
func LoadMultilingualSysMsg() (isFinish bool, err error) {
	isFinish = true
	// Step 1: Initialize map values
	MultilingualSysMsg = make(map[string]map[pub.ResCode]string, 0)
	// Step 2: Check if the SysLocaleList variable content is empty.
	if len(SysLocaleList) < 1 {
		isFinish = false
		zap.L().Info("LoadMultilingualSysMsg SysLocaleList is empty.")
		return
	}
	// Step 3: Preparing to read data from the sysmsg_t table in the database.
	sqlStr := `select code,content from sysmsg_t where language=$1 and dr=0`
	// Step 4: Get language from SysLocaleList
	for _, v := range SysLocaleList {
		lang := v.Language
		// Get data from sysmsg_t table
		rows, err := db.Query(sqlStr, lang)
		if err != nil {
			isFinish = false
			zap.L().Error("LoadMultilingualSysMsg db.Query failed:", zap.Error(err))
			return isFinish, err
		}
		langMsgMap := make(map[pub.ResCode]string)
		for rows.Next() {
			var code pub.ResCode
			var content string
			err = rows.Scan(&code, &content)
			if err != nil {
				isFinish = false
				zap.L().Error("LoadMultilingualSysMsg rows.Next() failed:", zap.Error(err))
				return isFinish, err
			}
			langMsgMap[code] = content
		}
		// Assign to the MutilingualSysMsg variable
		MultilingualSysMsg[lang] = langMsgMap
		// Close Rows
		rows.Close()
	}
	zap.L().Info("Multilingual system messages loaded successfully.")
	return
}
