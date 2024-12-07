// main.go
package main

import (
	"log"

	"zocket-assignment/api"
	"zocket-assignment/db"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize MongoDB connection
	connectionString := "mongodb+srv://rahul:kushwah@cluster0.46gecm5.mongodb.net/product_management?retryWrites=true&w=majority"
	mongo, err := db.NewMongoDB(connectionString)
	if err != nil {
		panic(err)
	}

	// Pass the MongoDB connection to the api package
	api.Mongodb = mongo

	// Create a new Gin router
	router := gin.Default()

	// Define your API routes
	ProductRoutes := router.Group("/products")
	{
		ProductRoutes.POST("/add", api.CreateProductHandler)
		ProductRoutes.GET("/:id", api.FetchProductHandler)
		ProductRoutes.GET("/", api.FetchAllProductsHandler)
		ProductRoutes.PUT("/update/:id", api.UpdateProductHandler)
		ProductRoutes.DELETE("/delete/:id", api.DeleteProductHandler)
		// Add other routes as needed
	}

	userGroup := router.Group("/users")
	{
		userGroup.POST("/add", api.CreateUserHandler)
		userGroup.GET("/:id", api.FetchUserHandler)
		userGroup.GET("/", api.FetchAllUsersHandler)
		userGroup.PUT("/update/:id", api.UpdateUserHandler)
		userGroup.DELETE("/delete/:id", api.DeleteUserHandler)
		
		
		// Add other user routes as needed
	}
	// Start the Gin server
	router.Run(":8080") // You can specify the port of your choice
	if err != nil {
		log.Fatal("Failed to start Gin server:", err)
	}
}
