package routes

import (
	"WatchHive/pkg/api/handler"
	"WatchHive/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(engine *gin.RouterGroup, adminHandler *handler.AdminHandler, categoryHandler *handler.CategoryHandler, inventoryHandler *handler.InventoryHandler) {

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

		inventorymanagement := engine.Group("/inventory")
		{
			inventorymanagement.POST("", inventoryHandler.AddInventory)
			inventorymanagement.GET("", inventoryHandler.ListProducts)
			inventorymanagement.PUT("", inventoryHandler.EditInventory)
			inventorymanagement.DELETE("", inventoryHandler.DeleteInventory)
			inventorymanagement.PATCH("", inventoryHandler.UpdateInventory)
		}

	}
}
