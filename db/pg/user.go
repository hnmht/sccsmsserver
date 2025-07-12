package pg

import (
	"crypto/md5"
	"encoding/hex"
	"sccsmsserver/pub"

	"go.uber.org/zap"
)

// Initialize user table
func initSysUser() (isFinish bool, err error) {
	// Step 1: Check if a record exists for the system defaut user 'admin' in the sysuser table.
	sqlStr := "select count(id) as rownum from sysuser where id=10000 and dr=0"
	hasRecord, isFinish, err := genericCheckRecord("sysuser", sqlStr)
	// Step 2: Exit if the record exists or an error occurs
	if hasRecord || !isFinish || err != nil {
		return
	}
	// Step 3: Insert a record for the system default user 'admin' into the sysuser tabel.
	sqlStr = `insert into sysuser(id,username,password,createtime,description,
		systemflag,usercode,createuserid) 
		values(10000,'admin',$1,now(),'System default',
		1,'admin',10000)`
	_, err = db.Exec(sqlStr, EncryptPassword(pub.DefaultPassword))
	if err != nil {
		isFinish = false
		zap.L().Error("initSysUser db.Exec failed:", zap.Error(err))
		return isFinish, err
	}
	return
}

// encryptPassword 加密密码
func EncryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(pub.Secret))

	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
