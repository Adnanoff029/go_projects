package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Adnanoff029/go_jwt/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file in database package")
	}
}
 
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "Access granted for api-1"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": "Access granted to api-2",
		})
	})

	router.Run(":" + port)
}
