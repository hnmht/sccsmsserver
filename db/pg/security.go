package pg

import (
	"sccsmsserver/pkg/security"

	"go.uber.org/zap"
)

// Get the RSA keys and Initialize RSA
func initRsa() (isFinish bool, err error) {
	var sqlStr string
	isFinish = true
	// Step 1: Get the public key and private key from the sysinfo table
	sqlStr = "select publickey,privatekey from sysinfo"
	var publickey, privatekey string
	err = db.QueryRow(sqlStr).Scan(&publickey, &privatekey)
	if err != nil {
		isFinish = false
		zap.L().Error("initRsa db.QueryRow failed:", zap.Error(err))
		return
	}
	if publickey == "" {
		isFinish = false
		zap.L().Error("initRsa public key is null.")
		return
	}
	if privatekey == "" {
		isFinish = false
		zap.L().Error("initRsa private key is null.")
		return
	}
	// step 2: Initialize RSA
	security.NewRsa(publickey, privatekey)
	return
}
