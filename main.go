package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"golang-Restaurant-Management/database"
	"golang-Restaurant-Management/routes"
	middleware "golang-Restaurant-Management/middleware"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	port := os.Getenv("PORT")

	if port==""{
		port = "8000"
	}

	router := gin.New()

	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemsRoutes(router)
	routes.InvoiceRoutes(router)
	
	router.Run(":"+port)
}