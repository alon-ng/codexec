package exercises

import (
	"codim/pkg/db"
	"codim/pkg/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *gin.RouterGroup, q *db.Queries, p *pgxpool.Pool, log *logger.Logger) {
	svc := NewService(q, p)
	ctrl := NewController(svc, log)

	exercisesGroup := router.Group("/exercises")
	{
		exercisesGroup.POST("/create", ctrl.Create)
		exercisesGroup.PUT("/update", ctrl.Update)
		exercisesGroup.DELETE("/delete/:uuid", ctrl.Delete)
		exercisesGroup.POST("/restore", ctrl.Restore)
		exercisesGroup.POST("/add-translation", ctrl.AddTranslation)
		exercisesGroup.GET("", ctrl.List)
		exercisesGroup.GET("/:uuid", ctrl.Get)
	}
}
