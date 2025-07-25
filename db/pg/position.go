package pg

// 初始化岗位档案
func initPosition() (isFinish bool, err error) {
	/* //检查岗位档案表中是否存在记录
	sqlStr := "select count(id) as rownum from operatingpost where id=10000"
	hasRecord, isFinish, err := checkRecord("operatingpost", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//如果表中没有数据则插入预置数据
	sqlStr = "insert into operatingpost(id,name,description,createuserid) values(10000,'预置岗位','系统预置岗位',10000)"
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initOperatingPost insert initvalue failed", zap.Error(err))
		return isFinish, err
	} */
	return
}
