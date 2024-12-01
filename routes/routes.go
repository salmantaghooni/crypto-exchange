// routes/routes.go
package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog"
    "crypto-exchange/controllers"
)

// SetupRoutes initializes all the routes for the application.
func SetupRoutes(router *gin.Engine, txController *controllers.TransactionController, logger zerolog.Logger) {
    // Define transaction routes
    router.POST("/transactions", txController.CreateTransaction)
    router.GET("/transactions/:id", txController.GetTransaction)

    // Add more routes as needed
}