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
	// r.Use(IpBlackListMiddleWare()) // IP Black list
	// Globle path
	superGroup := r.Group(pub.APIPath)
	{
		AuthRoute(superGroup)      // Auth
		CSARoute(superGroup)       // Construction Site Archive
		CSCRoute(superGroup)       // Construction Site Category
		CSORoute(superGroup)       // Construction Site Options
		DashboardRoute(superGroup) // Dashboard
		DCRoute(superGroup)        // Document Category
		DeptRoute(superGroup)      // Department
		DocRoute(superGroup)       // Document
		EPARoute(superGroup)       // Execution Project
		EPCRoute(superGroup)       // Execution Project Category
		EPTRoute(superGroup)       // Execution Project Template
		EORoute(superGroup)        // Execution Order
		EventRoute(superGroup)     // User Events
		FileRoute(superGroup)      // File
		IRFRoute(superGroup)       // Issue Resolution Form
		LandPageRoute(superGroup)  // Landing Page define
		MsgRoute(superGroup)       // Message
		OuRoute(superGroup)        // Online user
		PersonRoute(superGroup)    // Person
		PositionRoute(superGroup)  // Position
		PPEIFRoute(superGroup)     // Personal Protective Equipment Issuance Form
		PPEQuotaRoute(superGroup)  // Personal Protective Equipment Quota
		PPERoute(superGroup)       // Personal Protective Equipment
		PubRoute(superGroup)       // System public information
		RepRoute(superGroup)       // Report
		RLRoute(superGroup)        // Risk Level
		RoleRoute(superGroup)      // Role
		TCRoute(superGroup)        // Training Course
		TRRoute(superGroup)        // Training Record
		UDARoute(superGroup)       // User-defined Archive
		UDCRoute(superGroup)       // User-defined Category
		UserRoute(superGroup)      // User
		WORoute(superGroup)        // Work Order
	}
	//ping
	// r.POST("/ping", control.PubServerPing)
	// Monolithic application
	ui.AddRoutes(r)

	return r
}
