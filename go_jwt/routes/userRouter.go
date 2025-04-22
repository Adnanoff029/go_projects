package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Adnanoff029/go_jwt/controllers"
	"github.com/Adnanoff029/go_jwt/middleware"
)

func UserRoutes(router *gin.Engine){
	router.Use(middleware.Authenticate())
	router.GET("users", controllers.GetUsers())
	router.GET("users/:user_id", controllers.GetUser())
}