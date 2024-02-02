package routes

import(
	"github.com/gin-gonic/gin"
controllers "golang-Restaurant-Management/controllers"
)

func OrderItemsRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/orderItems", controllers.GetOrderItems())
	incomingRoutes.GET("/orderItems/:orderItems_id", controllers.GetOrderItem())
	incomingRoutes.GET("/orderItems-order/:order_id", controllers.GetOrderItemsByOrder())
	incomingRoutes.POST("/orderItems", controllers.CreateOrderItem())
	incomingRoutes.PATCH("/orderItems/:orderItems_id", controllers.UpdateOrderItem())
}