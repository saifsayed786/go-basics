package routes

import (
	"go-basics/controllers"
	_ "go-basics/docs"
	"go-basics/middlewares"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func PublicEndpoints(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	// Generate public endpoints - [signup] - api/v1/signup

	r.POST("/signup", controllers.Signup)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/signin", authMiddleware.LoginHandler)
}

func AuthenticatedEndpoints(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	// Generate Authenticated endpoints - [] - api/v1/auth/
	r.Use(authMiddleware.MiddlewareFunc())

	// User related endpoints
	r.POST("supervisor/create", controllers.CreateSupervisor)
	r.POST("admin/create", controllers.CreateAdmin)

	// Category related endpoints
	r.GET("category/", controllers.ListAllCategories)
	r.POST("category/create", controllers.CreateCategory)
	r.GET("category/:id", controllers.GetCategory)
	r.PATCH("category/:id", controllers.UpdateCatagory)

	// Product related endpoints.
	r.GET("product/", controllers.ListAllProduct)
	r.POST("product/create", controllers.CreateProduct)
	// r.POST("product/:id/image/upload", controllers.UploadProductImages)
	r.GET("product/:id", controllers.GetProduct)
	r.PATCH("product/:id", controllers.UpdateProduct)
}

func GetRouter(router chan *gin.Engine) {
	gin.ForceConsoleColor()
	r := gin.Default()

	// r.Use(middlewares.RequestLogger)
	// r.Use(gin.CustomRecovery(middlewares.LogFailedRequests))
	r.Use(cors.Default())
	authMiddleware, _ := middlewares.GetAuthMiddleware()

	// Create a BASE_URL - /api/v1
	v1 := r.Group("/api/v1/")
	PublicEndpoints(v1, authMiddleware)
	AuthenticatedEndpoints(v1.Group("auth"), authMiddleware)
	router <- r
}
