package routes

import (
	controller "github.com/Adnanoff029/go_jwt/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("users/signup", controller.Signup())
	router.POST("users/login", controller.Login())
}
