package pg

import (
	"time"

	"go.uber.org/zap"
)

type LandingPageInfo struct {
	SysNameDisp string    `db:"sysnamedisp" json:"sysnamedisp"`
	IntroText   string    `db:"introtext" json:"introtext"`
	File        File      `db:"fileid" json:"file"`
	ModifyDate  time.Time `db:"modifytime" json:"modifyDate"`
	ModifyUser  Person    `db:"modifierid" json:"modifier"`
	Ts          time.Time `db:"ts" json:"ts"`
}

// 初始化首页内容定义表
func initLandingPage() (isFinish bool, err error) {
	//检查首页内容定义表是否存在记录
	sqlStr := "select count(fileid) as rownum from landingpage"
	hasRecord, isFinish, err := genericCheckRecord("landingpage", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//表中没有数据则插入预置数据
	sqlStr = `insert into landingpage(sysnamedisp,introtext,fileid,modifierid) 
		values('SeaCloud现场管理系统',
		'一套实用有效的企业安全生产信息化系统,包含现场管理、文档管理、培训管理、劳保用品管理四大模块，帮助企业有效落实安全生产措施.'
		,0,10000);`
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("landingpage insert initvalue failed", zap.Error(err))
		return isFinish, err
	}
	return
}
