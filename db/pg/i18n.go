package pg

import (
	"sccsmsserver/pub"

	"go.uber.org/zap"
)

// Initialize the i18n table
func initI18n() (isFinish bool, err error) {
	// Step 1: Check if data already exists in the i18n table
	sqlStr := `select count(id) as rownum from i18n where dr=0`
	hasRecord, isFinish, err := genericCheckRecord("sysmsg", sqlStr)
	if hasRecord || !isFinish || err != nil {
		return
	}

	// Step 2: Insert Default localization into the i18n table
	insertSql := `insert into i18n(language,name,weekfirstday,shortdateformat,longdateformat,
		shorttimeformat,longtimeformat,timezone) 
		values($1,$2,$3,$4,$5,$6,$7,$8) returning id`
	err = db.QueryRow(insertSql,
		pub.DefaultLocale.Language, pub.DefaultLocale.Name, pub.DefaultLocale.WeekFirstDay, pub.DefaultLocale.ShortDateFormat, pub.DefaultLocale.LongDateFormat,
		pub.DefaultLocale.ShortTimeFormat, pub.DefaultLocale.LongTimeFormat, pub.DefaultLocale.TimeZone).Scan(&pub.DefaultLocale.ID)
	if err != nil {
		isFinish = false
		zap.L().Error("initI18n db.Exec(insertSql) failed:", zap.Error(err))
		return
	}

	return
}
