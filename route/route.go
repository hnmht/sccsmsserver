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
	r.Use(logger.GinLogger(), logger.GinRecovery(true)) //Global middleware
	r.Use(Cors())                                       //Allow the browser to make cross-origin requests
	// r.Use(IpBlackListMiddleWare()) //Ip黑名单
	// Globle path
	superGroup := r.Group(pub.APIPath)
	{
		AuthRoute(superGroup) // Auth
		CSARoute(superGroup)  // Construction Site Archive
		CSCRoute(superGroup)  // Construction Site Category
		CSORoute(superGroup)  // Construction Site Options
		// DashboardRoute(superGroup) //看板数据
		// DCRoute(superGroup)        //文档类别
		// DocRoute(superGroup)       //文档
		// DDRoute(superGroup)        //dd问题处理单
		DeptRoute(superGroup) // Department
		// EDRoute(superGroup)        //ed执行单
		EPCRoute(superGroup) // Execution Project Category
		// EIDRoute(superGroup)       //eid执行项目档案
		// EITRoute(superGroup)       //eit执行模板
		// EventRoute(superGroup)     //Event
		// FileRoute(superGroup)      //file文件
		// LandPageRoute(superGroup)  //landingPage首页定义
		// LDRoute(superGroup)        //劳保用品发放单
		// LPRoute(superGroup)        //劳保用品档案
		// LQRoute(superGroup)        //劳保用品岗位定额
		// MsgRoute(superGroup)       //消息
		OuRoute(superGroup)       // Online user
		PersonRoute(superGroup)   // Person
		PositionRoute(superGroup) // Position
		PubRoute(superGroup)      // System public information
		// RepRoute(superGroup)       //报表
		// RLRoute(superGroup)        //风险等级
		RoleRoute(superGroup) // Role
		// TCRoute(superGroup)        //tc培训课程
		// TRRoute(superGroup)        //tr培训记录
		UDARoute(superGroup)  // User-defined Archive
		UDCRoute(superGroup)  // User-defined Category
		UserRoute(superGroup) // User
		// WORoute(superGroup)        //wo指令单
	}
	//ping
	// r.POST("/ping", control.PubServerPing)
	// Monolithic application
	ui.AddRoutes(r)

	return r
}
