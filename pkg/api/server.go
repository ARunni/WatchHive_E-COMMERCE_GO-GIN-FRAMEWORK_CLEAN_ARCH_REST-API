package http

import (
	"WatchHive/pkg/api/handler"
	"WatchHive/pkg/routes"
	"log"

	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(adminHandler *handler.AdminHandler, userHandler *handler.UserHandler, otpHandler *handler.OtpHandler, catgoryHandler *handler.CategoryHandler, productHandler *handler.ProductHandler) *ServerHTTP {
	engine := gin.New()

	engine.Use(gin.Logger())

	engine.GET("/validate_token", adminHandler.ValidateRefreshTokenAndCreateNewAccess)

	routes.UserRoutes(engine.Group("/user"), userHandler, otpHandler, productHandler)
	routes.AdminRoutes(engine.Group("/admin"), adminHandler, catgoryHandler, productHandler)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	err := sh.engine.Run(":7000")
	if err != nil {
		log.Fatal("gin engine couldn't start")
	}

}
