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
	r.Use(logger.GinLogger(), logger.GinRecovery(true)) // Global middleware
	// r.Use(Cors())                                       // Allow the browser to make cross-origin requests
	// r.Use(IpBlackListMiddleWare()) //Ip黑名单
	// Globle path
	superGroup := r.Group(pub.APIPath)
	{
		AuthRoute(superGroup) // Auth
		CSARoute(superGroup)  // Construction Site Archive
		CSCRoute(superGroup)  // Construction Site Category
		CSORoute(superGroup)  // Construction Site Options
		// DashboardRoute(superGroup) //看板数据
		DCRoute(superGroup) // Document Category
		// DocRoute(superGroup)       //文档
		DeptRoute(superGroup) // Department
		EPARoute(superGroup)  // Execution Project
		EPCRoute(superGroup)  // Execution Project Category
		EPTRoute(superGroup)  // Execution Project Template
		EORoute(superGroup)   // Execution Order
		// EventRoute(superGroup)     //Event
		FileRoute(superGroup) // File
		IRFRoute(superGroup)  // Issue Resolution Form
		// LandPageRoute(superGroup)  //landingPage首页定义
		// LDRoute(superGroup)        //劳保用品发放单
		// LQRoute(superGroup)        //劳保用品岗位定额
		// MsgRoute(superGroup)       //消息
		OuRoute(superGroup)       // Online user
		PersonRoute(superGroup)   // Person
		PositionRoute(superGroup) // Position
		PPERoute(superGroup)      //劳保用品档案
		PubRoute(superGroup)      // System public information
		// RepRoute(superGroup)       //报表
		RLRoute(superGroup)   // Risk Level
		RoleRoute(superGroup) // Role
		TCRoute(superGroup)   // Training Course
		// TRRoute(superGroup)        //tr培训记录
		UDARoute(superGroup)  // User-defined Archive
		UDCRoute(superGroup)  // User-defined Category
		UserRoute(superGroup) // User
		WORoute(superGroup)   // Work Order
	}
	//ping
	// r.POST("/ping", control.PubServerPing)
	// Monolithic application
	ui.AddRoutes(r)

	return r
}
