package http

import (
	"WatchHive/pkg/api/handler"
	"WatchHive/pkg/routes"
	"log"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(adminHandler *handler.AdminHandler,
	userHandler *handler.UserHandler,
	otpHandler *handler.OtpHandler,
	catgoryHandler *handler.CategoryHandler,
	productHandler *handler.ProductHandler,
	cartHandler *handler.CartHandler,
	orderHandler *handler.OrderHandler,
	paymentHandler *handler.PaymentHandler,
	walletHandler *handler.WalletHandler,
	offerHandler *handler.OfferHandler) *ServerHTTP {
	engine := gin.New()

	engine.Use(gin.Logger())
	engine.LoadHTMLGlob("pkg/templates/index.html")
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.GET("/validate_token", adminHandler.ValidateRefreshTokenAndCreateNewAccess)

	routes.UserRoutes(engine.Group("/user"), userHandler, otpHandler, productHandler, cartHandler, orderHandler, paymentHandler, walletHandler)
	routes.AdminRoutes(engine.Group("/admin"), adminHandler, catgoryHandler, productHandler, paymentHandler, orderHandler, offerHandler)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	err := sh.engine.Run(":7000")
	if err != nil {
		log.Fatal("gin engine couldn't start")
	}

}
