package pg

import (
	"sccsmsserver/pub"

	"go.uber.org/zap"
)

// Upgrade database schema version
func upgradeDb() (isFinish bool, err error) {
	isFinish = true
	// Retrieve the current data schema version from the sysinfo table.
	var currentDbVer string
	sqlStr := `select dbversion from sysinfo limit 1`
	err = db.QueryRow(sqlStr).Scan(&currentDbVer)
	if err != nil {
		zap.L().Error("upgradeDb db.QueryRow failed:", zap.Error(err))
		isFinish = false
		return
	}
	// If the database schema version and the application's database version are the same, no upgrade is needed.
	if pub.DbVersion == currentDbVer {
		return
	}

	// If the database schema version is greater than the application's database version,
	// output an error log asking the user to upgrade the application.
	if pub.DbVersion < currentDbVer {
		isFinish = false
		zap.L().Error("The current database version is newer than the version used by the application. The application cannot start. Please upgrade the application.")
		return
	}

	// if the database schema version is less than the application's database version,
	// call the relevant function to upgrade the database schema.
	if pub.DbVersion > currentDbVer {
		// call the relevant function
	}

	return
}
