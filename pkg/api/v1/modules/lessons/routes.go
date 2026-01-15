package lessons

import (
	"codim/pkg/db"
	"codim/pkg/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *gin.RouterGroup, q *db.Queries, p *pgxpool.Pool, log *logger.Logger) {
	svc := NewService(q, p)
	ctrl := NewController(svc, log)

	lessonsGroup := router.Group("/lessons")
	{
		lessonsGroup.POST("/create", ctrl.Create)
		lessonsGroup.PUT("/update", ctrl.Update)
		lessonsGroup.DELETE("/delete/:uuid", ctrl.Delete)
		lessonsGroup.POST("/restore", ctrl.Restore)
		lessonsGroup.POST("/add-translation", ctrl.AddTranslation)
		lessonsGroup.GET("", ctrl.List)
		lessonsGroup.GET("/:uuid", ctrl.Get)
	}
}
