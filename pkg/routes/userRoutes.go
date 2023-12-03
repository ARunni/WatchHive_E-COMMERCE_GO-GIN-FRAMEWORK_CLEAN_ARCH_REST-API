package routes

import (
	"WatchHive/pkg/api/handler"
	"WatchHive/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.RouterGroup, userHandler *handler.UserHandler, otpHandler *handler.OtpHandler, productHandler *handler.ProductHandler, cartHandler *handler.CartHandler, orderHandler *handler.OrderHandler) {

	engine.POST("/signup", userHandler.UserSignUp)
	engine.POST("/login", userHandler.LoginHandler)

	engine.POST("/otplogin", otpHandler.SendOTP)
	engine.POST("/verifyotp", otpHandler.VerifyOTP)
	engine.GET("", productHandler.ListProducts)

	engine.Use(middleware.UserAuthMiddleware)
	{
		profile := engine.Group("/profile")
		{
			profile.POST("/address", userHandler.AddAddress)
			profile.GET("", userHandler.ShowUserDetails)
			profile.GET("/alladdress", userHandler.GetAllAddress)
			profile.PUT("", userHandler.EditProfile)
			profile.PATCH("", userHandler.ChangePassword)
		}
		cart := engine.Group("/cart")
		{
			cart.POST("", cartHandler.AddToCart)
			cart.GET("", cartHandler.ListCartItems)
			cart.PATCH("", cartHandler.UpdateProductQuantityCart)
			cart.PUT("", cartHandler.RemoveFromCart)
		}
		checkout := engine.Group("/orders")
		{
			checkout.GET("/checkout", orderHandler.CheckOut)
			checkout.POST("", orderHandler.OrderItemsFromCart)
			// checkout.POST("/placeorder", orderHandler.PlaceOrderCOD)
			checkout.GET("", orderHandler.GetOrderDetails)
			checkout.PUT("", orderHandler.CancelOrder)
			checkout.PATCH("",orderHandler.ReturnOrderCod)
		}
	}

}
