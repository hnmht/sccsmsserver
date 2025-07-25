package pg

// 部门档案表初始化
func initDepartment() (isFinish bool, err error) {
	//检查部门档案表中是否存在记录
	/* sqlStr := "select count(id) as rownum from department where id=10000"
	hasRecord, isFinish, err := checkRecord("department", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//如果表中没有数据则插入预置数据
	sqlStr = "insert into department(id,deptcode,deptname,description,createuserid) values(10000,'default','预置部门','系统预置部门',10000)"
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("initDepartment insert initvalue failed", zap.Error(err))
		return isFinish, err
	} */
	return
}
