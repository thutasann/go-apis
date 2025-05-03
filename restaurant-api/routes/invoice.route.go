package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thutasann/restaurant-api/controllers"
)

func InvoiceRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/invoices", controllers.GetInvoices())
	incomingRoutes.GET("/invoice/:invoice_id", controllers.GetInvoice())
	incomingRoutes.POST("/invoices", controllers.CreateInvoice())
	incomingRoutes.PATCH("/invoice/:invoice_id", controllers.UpdateInvoice())
}
