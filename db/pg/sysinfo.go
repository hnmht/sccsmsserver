package pg

import (
	"sccsmsserver/pkg/environment"
	"sccsmsserver/pkg/mysf"
	"sccsmsserver/pkg/security"
	"sccsmsserver/pub"

	"go.uber.org/zap"
)

// Initialize the sysinfo table
func initSysInfo() (isFinish bool, err error) {
	var rowNum int
	var sqlStr string
	isFinish = true
	// Step 1: Query the row count from sysinfo table.
	// There should be exactly one row in this table
	sqlStr = "select count(isfinish) from sysinfo"
	err = db.QueryRow(sqlStr).Scan(&rowNum)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysInfo db.QueryRow failed:", zap.Error(err))
		return
	}

	// Step 2: If the sysinfo table contains more than one row, Exit the database table creation operation
	if rowNum > 1 {
		isFinish = false
		zap.L().Error("initSysInfo query the sysinfo table records failed:  contains more than one row.")
		return
	}

	// Step 3: If there is a record in the sysinfo table, delete it.
	if rowNum == 1 {
		sqlStr = "delete from sysinfo"
		_, err = db.Exec(sqlStr)
		if err != nil {
			isFinish = false
			zap.L().Error("initSysInfo db.Exec delete old data failed:", zap.Error(err))
			return
		}
	}

	// Step 4: Insert data into the sysinfo table
	// Step 4.1： Get the machines's network adapter MAC address
	macArray, err := environment.GetMacArray()
	if err != nil {
		isFinish = false
		zap.L().Error("initSysInfo environment.GetMacArray failed:", zap.Error(err))
		return
	}

	// Step 4.2: Get the machine's motherboard serial number
	serialNumber, err := environment.GetSerialNumber()
	if err != nil {
		isFinish = false
		zap.L().Error("initSysInfo environment.GetSerialNumber failed:", zap.Error(err))
		return
	}

	// Step 4.3: Get the machine's hash value
	machineHash, err := environment.GetMachineHash(macArray, serialNumber)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysInfo environment.GetMachineHash failed:", zap.Error(err))
		return
	}

	// Step 4.4: Generate RSA public and private keys
	privateKey, publicKey, err := security.GenRsaKey(2048)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysInfo security.GenRsaKey failed:", zap.Error(err))
		return
	}

	// Step 4.5: Generate a unique database ID.
	dbID := mysf.GenID()

	// Step 4.6: Insert data into the sysinfo table.
	sqlInsert := `insert into sysinfo(dbid,serialnumber,macarray,machinehash,privatekey,publickey,dbversion,starttime) values($1,$2,$3,$4,$5,$6,$7,now())`
	_, err = db.Exec(sqlInsert, dbID, serialNumber, macArray, machineHash, privateKey, publicKey, pub.DbVersion)
	if err != nil {
		isFinish = false
		zap.L().Error("initSysInfo db.Exec insert data failed:", zap.Error(err))
		return
	}

	return
}
