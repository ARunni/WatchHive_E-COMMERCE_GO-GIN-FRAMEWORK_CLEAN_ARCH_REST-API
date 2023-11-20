package routes

import (
	"WatchHive/pkg/api/handler"
	"WatchHive/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.RouterGroup, userHandler *handler.UserHandler) {

	engine.POST("/signup", userHandler.UserSignUp)
	engine.POST("/login", userHandler.LoginHandler)

	engine.Use(middleware.UserAUthMiddleware)
	{

	}

}
