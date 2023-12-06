package routes

import (
	"WatchHive/pkg/api/handler"
	"WatchHive/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(engine *gin.RouterGroup, adminHandler *handler.AdminHandler, categoryHandler *handler.CategoryHandler, productHandler *handler.ProductHandler, paymentHandler *handler.PaymentHandler, orderHandler *handler.OrderHandler) {

	engine.POST("", adminHandler.LoginHandler)

	engine.Use(middleware.AdminAuthMiddleware)
	{
		userManagement := engine.Group("/users")
		{
			userManagement.PATCH("/block", adminHandler.BlockUser)
			userManagement.PATCH("/unblock", adminHandler.UnBlockUser)
			userManagement.GET("", adminHandler.GetUsers)
		}
		categorymanagement := engine.Group("/category")
		{
			categorymanagement.POST("", categoryHandler.AddCategory)
			categorymanagement.GET("", categoryHandler.GetCategory)
			categorymanagement.PATCH("", categoryHandler.UpdateCategory)
			categorymanagement.DELETE("", categoryHandler.DeleteCategory)
		}

		productManagement := engine.Group("/product")
		{
			productManagement.GET("", productHandler.ListProducts)
			productManagement.POST("", productHandler.AddProduct)
			productManagement.PATCH("", productHandler.EditProduct)
			productManagement.DELETE("", productHandler.DeleteProduct)
			productManagement.PATCH("/stock", productHandler.UpdateProduct)
		}

		paymentManagement := engine.Group("/payment")
		{
			paymentManagement.POST("", paymentHandler.AddPaymentMethod)
		}
		orderManagement := engine.Group("/orders")
		{
			orderManagement.GET("", orderHandler.GetAllOrderDetailsForAdmin)
			orderManagement.PATCH("", orderHandler.ApproveOrder)
			orderManagement.DELETE("", orderHandler.CancelOrderFromAdmin)
		}
	}
}
