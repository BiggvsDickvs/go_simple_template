package routes

import (
	"go_proj_example/internal/handlers"
	"go_proj_example/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {
	router.GET("/login", handlers.LoginPage)
	router.POST("/login", handlers.Login)
	router.GET("/signup", handlers.SignupPage)
	router.POST("/signup", handlers.Register)
	router.GET("/logout", handlers.SignOut)
	router.GET("/404", handlers.NotFound)

	protected := router.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/", handlers.HomePage)
	}
}
