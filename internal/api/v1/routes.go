package api

import (
	apiDocs "codim/api"
	"codim/internal/api/v1/modules/users"
	"codim/internal/db"
	"codim/internal/utils/logger"

	"github.com/gin-gonic/gin"
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

func NewRouter(q *db.Queries, log *logger.Logger) *gin.Engine {
	r := gin.Default()

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		users.RegisterRoutes(v1, q, log)
	}

	return r
}
