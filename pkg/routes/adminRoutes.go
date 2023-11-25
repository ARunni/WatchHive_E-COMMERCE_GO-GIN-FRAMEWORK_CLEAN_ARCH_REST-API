package routes

import (
	"WatchHive/pkg/api/handler"
	"WatchHive/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(engine *gin.RouterGroup, adminHandler *handler.AdminHandler, categoryHandler *handler.CategoryHandler, productHandler *handler.ProductHandler) {

	engine.POST("", adminHandler.LoginHandler)

	engine.Use(middleware.AdminAuthMiddleware)
	{
		userManagement := engine.Group("/users")
		{
			userManagement.PUT("/block", adminHandler.BlockUser)
			userManagement.PUT("/unblock", adminHandler.UnBlockUser)
			userManagement.GET("", adminHandler.GetUsers)
		}
		categorymanagement := engine.Group("/category")
		{
			categorymanagement.POST("", categoryHandler.AddCategory)
			categorymanagement.GET("", categoryHandler.GetCategory)
			categorymanagement.PUT("", categoryHandler.UpdateCategory)
			categorymanagement.DELETE("", categoryHandler.DeleteCategory)
		}

		productManagement := engine.Group("/product")
		{
			productManagement.GET("", productHandler.ListProducts)
			productManagement.POST("", productHandler.AddProduct)
			productManagement.PUT("", productHandler.EditProduct)
			productManagement.DELETE("", productHandler.DeleteProduct)
			productManagement.PATCH("", productHandler.UpdateProduct)
		}

	}
}
