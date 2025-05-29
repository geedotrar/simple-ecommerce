package routes

import (
	"product-service/internal/handlers"
	"product-service/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ProductRouter interface {
	Mount()
}

type productRouterImpl struct {
	v       *gin.RouterGroup
	handler handlers.ProductHandler
}

func NewProductRouter(v *gin.RouterGroup, handler handlers.ProductHandler) ProductRouter {
	return &productRouterImpl{v: v, handler: handler}
}

func (r *productRouterImpl) Mount() {
	r.v.Use(cors.Default())
	r.v.Use(middleware.AuthMiddleware())
	r.v.GET("", middleware.RequireAnyPermission("view_all_products", "view_active_products"), r.handler.GetAllProducts)

	r.v.POST("/create", middleware.RequirePermission("create_products"), r.handler.CreateProduct)
	r.v.PUT("/update/:id", middleware.RequirePermission("update_products"), r.handler.UpdateProduct)
	r.v.PUT("/update-status/:id", middleware.RequirePermission("update_products"), r.handler.UpdateProductStatus)
	r.v.DELETE("/delete/:id", middleware.RequirePermission("delete_products"), r.handler.DeleteProduct)
}
