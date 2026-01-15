package courses

import (
	"codim/pkg/db"
	"codim/pkg/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPublicRoutes(router *gin.RouterGroup, q *db.Queries, p *pgxpool.Pool, log *logger.Logger) {
	svc := NewService(q, p)
	ctrl := NewController(svc, log)

	coursesGroup := router.Group("/courses")
	{
		coursesGroup.GET("", ctrl.List)
		coursesGroup.GET("/:uuid", ctrl.Get)
	}
}

func RegisterAdminRoutes(router *gin.RouterGroup, q *db.Queries, p *pgxpool.Pool, log *logger.Logger) {
	svc := NewService(q, p)
	ctrl := NewController(svc, log)

	coursesGroup := router.Group("/courses")
	{
		coursesGroup.POST("/create", ctrl.Create)
		coursesGroup.PUT("/update", ctrl.Update)
		coursesGroup.DELETE("/delete/:uuid", ctrl.Delete)
		coursesGroup.POST("/restore", ctrl.Restore)
		coursesGroup.POST("/add-translation", ctrl.AddTranslation)
	}
}
