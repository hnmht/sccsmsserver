package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func PPEQuotaRoute(g *gin.RouterGroup) {
	LQGroup := g.Group("/ppeq", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Get Personal Protective Equipment Quota List
		LQGroup.POST("/list", handlers.GetPPEQuotaListHandler)
		// Get The Personal Protective Equipment Quota detail by HID
		LQGroup.POST("/detail", handlers.GetPPEQuotaInfoByHIDHandler)
		// Add Personal Protective Equipment Quota
		LQGroup.POST("/add", handlers.AddPPEQuotaHandler)
		// Modify Personal Protective Equipment Quota
		LQGroup.POST("/edit", handlers.EditPPEQuotaHandler)
		// Delete Personal Protective Equipment Quota
		LQGroup.POST("/del", handlers.DeletePPEQuotaHandler)
		// Confirm PPE Quota
		LQGroup.POST("/confirm", handlers.ConfirmPPEQuotaHandler)
		// Unconfirm PPE Quota
		LQGroup.POST("/unconfirm", handlers.UnconfirmPPEQuotaHandler)
		// Check if a PPE Position Quota for the same period
		LQGroup.POST("/check", handlers.CheckPPEQuotaExistHandler)
		// Get the list of all position that have PPE Quotas within the same period
		LQGroup.POST("/positions", handlers.GetPPEPositionsPeriodHandler)
	}
}
