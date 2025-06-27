package pg

import (
	"time"

	"go.uber.org/zap"
)

// Localization Struct
type Localization struct {
	ID              int32     `db:"id" json:"ID"`
	Language        string    `db:"language" json:"Language"`
	Name            string    `db:"name" json:"Name"`
	WeekFirstDay    string    `db:"weekfirstday" json:"WeekFirstDay"`
	ShortDateFormat string    `db:"shortdateformat" json:"ShortDateFormat"`
	LongDateFormat  string    `db:"longdateformat" json:"LongDateFormat"`
	ShortTimeFormat string    `db:"shorttimeformat" json:"ShortTimeFormat"`
	LongTimeFormat  string    `db:"longtimeformat" json:"LongTimeFormat"`
	TimeZone        string    `db:"timezone" json:"TimeZone"`
	Description     string    `db:"description" json:"Description"`
	SystemFlag      int16     `db:"systemflag" json:"SystemFlag"`
	CreateTime      time.Time `db:"createtime" json:"CreateTime"`
	Creator         Person    `db:"creatorid" json:"Creator"`
	ModifyTime      time.Time `db:"modifytime" json:"ModifyTime"`
	Modifier        Person    `db:"modifierid" json:"Modifier"`
	Ts              time.Time `db:"ts" json:"Ts"`
	Dr              int16     `db:"dr" json:"Dr"`
}

// Default locale
var DefaultLocale = Localization{
	ID:              1,
	Language:        "en_us",
	Name:            "English United States",
	WeekFirstDay:    "Sunday",
	ShortDateFormat: "MM/DD/YY",
	LongDateFormat:  "MM/DD/YYYY",
	ShortTimeFormat: "HH:MM AM/PM",
	LongTimeFormat:  "HH:MM:SS AM/PM",
	TimeZone:        "UTC-5",
	Description:     "System default locale",
	SystemFlag:      1,
}

// system locale list
var SysLocaleList []Localization

// Initialize the i18n table
func initI18n() (isFinish bool, err error) {
	isFinish = true
	// Step 1: Check if data already exists in the i18n table
	sqlStr := `select count(id) as rownum from i18n where dr=0`
	hasRecord, isFinish, err := genericCheckRecord("sysmsg", sqlStr)
	if hasRecord || !isFinish || err != nil {
		return
	}

	// Step 2: Insert Default localization into the i18n table
	insertSql := `insert into i18n(language,name,weekfirstday,shortdateformat,longdateformat,
		shorttimeformat,longtimeformat,timezone,description,systemflag) 
		values($1,$2,$3,$4,$5,$6,$7,$8,$9,1) returning id`
	err = db.QueryRow(insertSql,
		DefaultLocale.Language, DefaultLocale.Name, DefaultLocale.WeekFirstDay, DefaultLocale.ShortDateFormat, DefaultLocale.LongDateFormat,
		DefaultLocale.ShortTimeFormat, DefaultLocale.LongTimeFormat, DefaultLocale.TimeZone, DefaultLocale.Description).Scan(&DefaultLocale.ID)
	if err != nil {
		isFinish = false
		zap.L().Error("initI18n db.Exec(insertSql) failed:", zap.Error(err))
		return
	}

	return
}

// Initialize system locale list
func initSysLocalList() (err error) {
	// Retrieve system locale list from i18n table
	sqlStr := `select id,language,name,weekfirstday,shortdateformat,
	longdateformat,shorttimeformat,longtimeformat,timezone,description,
	systemflag,createtime,creatorid,modifytime,modifierid,
	ts from i18n where dr=0`
	rows, err := db.Query(sqlStr)
	if err != nil {
		zap.L().Error("initSysLocalList db.Query failed:", zap.Error(err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var l Localization
		err = rows.Scan(&l.ID, &l.Language, &l.Name, &l.WeekFirstDay, &l.ShortDateFormat,
			&l.LongDateFormat, &l.ShortDateFormat, &l.LongTimeFormat, &l.TimeZone, &l.Description,
			&l.SystemFlag, &l.CreateTime, &l.Creator.ID, &l.ModifyTime, &l.Modifier.ID,
			&l.Ts)
		if err != nil {
			zap.L().Error("initSysLocalList rows.Next failed:", zap.Error(err))
			return
		}
		// Get creator deatils
		// Get modifier details
		SysLocaleList = append(SysLocaleList, l)
	}
	zap.L().Info("Server locale list initilized successfully.")
	return
}
