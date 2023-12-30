package routes

import (
	"WatchHive/pkg/api/handler"
	"WatchHive/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(engine *gin.RouterGroup, adminHandler *handler.AdminHandler, categoryHandler *handler.CategoryHandler, productHandler *handler.ProductHandler, paymentHandler *handler.PaymentHandler, orderHandler *handler.OrderHandler, offerHandler *handler.OfferHandler) {

	engine.POST("", adminHandler.LoginHandler)

	engine.Use(middleware.AdminAuthMiddleware)
	{
		engine.GET("/dashboard", adminHandler.AdminDashBoard)
		engine.GET("/printsales", adminHandler.PrintSalesByDate)
		engine.GET("/currentsalesreport", adminHandler.FilteredSalesReport)
		engine.GET("/salesreport", adminHandler.SalesReportByDate)

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
			productManagement.GET("", productHandler.ListProductsAdmin)
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

		offer := engine.Group("offer")
		{
			offer.POST("/product_offer", offerHandler.AddProductOffer)
			offer.GET("/product_offer", offerHandler.GetProductOffer)
			offer.DELETE("/product_offer", offerHandler.ExpireProductOffer)

			offer.POST("/category_offer", offerHandler.AddCategoryOffer)
			offer.GET("/category_offer", offerHandler.GetCategoryOffer)
			offer.DELETE("/category_offer", offerHandler.ExpireCategoryOffer)
		}

	}
}
