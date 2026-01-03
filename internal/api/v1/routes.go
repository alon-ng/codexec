package api

import (
	apiDocs "codim/api"
	authProvider "codim/internal/api/auth"
	"codim/internal/api/v1/cache"
	"codim/internal/api/v1/middleware"
	"codim/internal/api/v1/modules/auth"
	"codim/internal/api/v1/modules/users"
	"codim/internal/db"
	"codim/internal/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	apiDocs.SwaggerInfo.Title = "Codim API"
	apiDocs.SwaggerInfo.Description = "This is the Codim API server."
	apiDocs.SwaggerInfo.Version = "1.0"
	apiDocs.SwaggerInfo.Host = "localhost:8080"
	apiDocs.SwaggerInfo.BasePath = "/api/v1"
	apiDocs.SwaggerInfo.Schemes = []string{"http", "https"}
}

func NewRouter(q *db.Queries, log *logger.Logger, authProvider *authProvider.Provider, redisClient *redis.Client) *gin.Engine {
	r := gin.Default()

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userCache := cache.NewUserCache(redisClient, q, log)

	v1 := r.Group("/api/v1")
	{
		auth.RegisterRoutes(v1, q, log, authProvider)

		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(authProvider, userCache, log))
		{
			admin := protected.Group("/")
			admin.Use(middleware.AdminMiddleware(log))
			{
				users.RegisterRoutes(admin, q, log)
			}
		}
	}

	log.Info("Router initialized successfully")

	return r
}
