package route

import (
	"sccsmsserver/logger"
	"sccsmsserver/pub"
	"sccsmsserver/ui"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	gin.SetMode(mode)
	r := gin.New()
	//全局中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.Use(Cors()) //允许浏览器跨域
	// r.Use(IpBlackListMiddleWare()) //Ip黑名单
	//全局路径
	superGroup := r.Group(pub.APIPath)
	{
		AuthRoute(superGroup) //auth权限
		// DashboardRoute(superGroup) //看板数据
		// DCRoute(superGroup)        //文档类别
		// DocRoute(superGroup)       //文档
		// DDRoute(superGroup)        //dd问题处理单
		// DeptRoute(superGroup)      //dept部门
		// EDRoute(superGroup)        //ed执行单
		// EICRoute(superGroup)       //eic执行项目类别
		// EIDRoute(superGroup)       //eid执行项目档案
		// EITRoute(superGroup)       //eit执行模板
		// EventRoute(superGroup)     //Event
		// FileRoute(superGroup)      //file文件
		// LandPageRoute(superGroup)  //landingPage首页定义
		// LDRoute(superGroup)        //劳保用品发放单
		// LPRoute(superGroup)        //劳保用品档案
		// LQRoute(superGroup)        //劳保用品岗位定额
		// MsgRoute(superGroup)       //消息
		// OPRoute(superGroup)        //岗位档案
		// OuRoute(superGroup)        //在线用户
		// PersonRoute(superGroup)    //person人员
		// PubRoute(superGroup)       //pub公共信息
		// RepRoute(superGroup)       //报表
		// RLRoute(superGroup)        //风险等级
		// RoleRoute(superGroup)      //role角色
		// SIRoute(superGroup)        //si现场档案
		// SICRoute(superGroup)       //sic现场档案分类
		// TCRoute(superGroup)        //tc培训课程
		// TRRoute(superGroup)        //tr培训记录
		// UDCRoute(superGroup)       //udc用户自定义档案类别
		// UDDRoute(superGroup)       //udd用户自定义档案
		UserRoute(superGroup) //user用户
		// WORoute(superGroup)        //wo指令单

	}
	//ping
	// r.POST("/ping", control.PubServerPing)
	// Monolithic application
	ui.AddRoutes(r)

	return r
}
