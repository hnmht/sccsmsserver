package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func PPEIFRoute(g *gin.RouterGroup) {
	PPEIFGroup := g.Group("/ppeif", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Add PPE Issuance Form
		PPEIFGroup.POST("/add", handlers.AddPPEIFHandler)
		// Use the wizard to generate PPE Issuance Form
		PPEIFGroup.POST("/wizard", handlers.WiardAddPPEIFHandler)
		// Get PPE Issuance Form List
		PPEIFGroup.POST("/list", handlers.GetPPEIFListHandler)
		// Get PPE Issuance Form detail by HID
		PPEIFGroup.POST("/detail", handlers.GetPPEIFInfoByHIDHandler)
		// Modify PPE Issuance Form
		PPEIFGroup.POST("/edit", handlers.EditPPEIFHandler)
		// Delete PPE Issuance Form
		PPEIFGroup.POST("/del", handlers.DeletePPEIFHandler)
		// Confirm PPE Issuance Form
		PPEIFGroup.POST("/confirm", handlers.ConfirmPPEIFHandler)
		// Unconfirm PPE Issuance Form
		PPEIFGroup.POST("/unconfirm", handlers.UnconfirmPPEIFHandler)
		// Get PPE Issuance Form Report
		PPEIFGroup.POST("/rep", handlers.GetPPEIFReportHandler)
	}
}
