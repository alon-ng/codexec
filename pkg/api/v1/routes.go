package api

import (
	apiDocs "codim/api"
	authProvider "codim/pkg/api/auth"
	"codim/pkg/api/v1/cache"
	"codim/pkg/api/v1/middleware"
	"codim/pkg/api/v1/modules/auth"
	"codim/pkg/api/v1/modules/courses"
	"codim/pkg/api/v1/modules/exercises"
	"codim/pkg/api/v1/modules/lessons"
	"codim/pkg/api/v1/modules/me"
	"codim/pkg/api/v1/modules/users"
	"codim/pkg/db"
	"codim/pkg/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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
	// Disable Gin's default logger output
	gin.DefaultWriter = log.Writer()
	gin.DefaultErrorWriter = log.WriterLevel(logrus.ErrorLevel)

	r := gin.New()
	r.Use(middleware.GinLogger(log))
	r.Use(gin.Recovery())

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userCache := cache.NewUserCache(redisClient, q, log)

	v1 := r.Group("/api/v1")
	{
		auth.RegisterRoutes(v1, q, log, authProvider)
		courses.RegisterPublicRoutes(v1, q, log)

		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(authProvider, userCache, log))
		{
			me.RegisterRoutes(protected, q, log)
			admin := protected.Group("/")
			admin.Use(middleware.AdminMiddleware(log))
			{
				users.RegisterRoutes(admin, q, log)
				courses.RegisterAdminRoutes(admin, q, log)
				lessons.RegisterRoutes(admin, q, log)
				exercises.RegisterRoutes(admin, q, log)
			}
		}
	}

	log.Info("Router initialized successfully")

	return r
}
