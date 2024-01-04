package routes

import (
	"WatchHive/pkg/api/handler"
	"WatchHive/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.RouterGroup,
	userHandler *handler.UserHandler,
	otpHandler *handler.OtpHandler,
	productHandler *handler.ProductHandler,
	cartHandler *handler.CartHandler,
	orderHandler *handler.OrderHandler,
	paymentHandler *handler.PaymentHandler,
	walletHandler *handler.WalletHandler,
	couponHandler *handler.CouponHandler) {

	engine.POST("/signup", userHandler.UserSignUp)
	engine.POST("/login", userHandler.LoginHandler)

	engine.POST("/otplogin", otpHandler.SendOTP)
	engine.POST("/verifyotp", otpHandler.VerifyOTP)
	engine.GET("/products", productHandler.ListProductsUser)
	engine.GET("/payment", paymentHandler.MakePaymentRazorpay)
	engine.GET("/verifypayment", paymentHandler.VerifyPayment)

	engine.Use(middleware.UserAuthMiddleware)
	{
		profile := engine.Group("/profile")
		{
			profile.POST("/address", userHandler.AddAddress)
			profile.GET("", userHandler.ShowUserDetails)
			profile.GET("/alladdress", userHandler.GetAllAddress)
			profile.PATCH("", userHandler.EditProfile)
			profile.PATCH("/password", userHandler.ChangePassword)
		}
		cart := engine.Group("/cart")
		{
			cart.POST("", cartHandler.AddToCart)
			cart.GET("", cartHandler.ListCartItems)
			cart.PATCH("", cartHandler.UpdateProductQuantityCart)
			cart.DELETE("", cartHandler.RemoveFromCart)
		}
		checkout := engine.Group("/orders")
		{
			checkout.GET("/checkout", orderHandler.CheckOut)
			checkout.POST("", orderHandler.OrderItemsFromCart)
			checkout.GET("", orderHandler.GetOrderDetails)
			checkout.DELETE("", orderHandler.CancelOrder)
			checkout.PATCH("", orderHandler.ReturnOrder)
			checkout.GET("/print", orderHandler.PrintInvoice)

		}
		wallet := engine.Group("/wallet")

		{
			wallet.GET("", walletHandler.GetWallet)
			wallet.GET("/history",walletHandler.GetWalletHistory)
		}
		coupon := engine.Group("/coupon")
		{
			coupon.GET("", couponHandler.GetCouponUser)
		}

	}

}
