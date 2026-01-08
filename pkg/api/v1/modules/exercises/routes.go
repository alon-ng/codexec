package exercises

import (
	"codim/pkg/db"
	"codim/pkg/utils/logger"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, q *db.Queries, log *logger.Logger) {
	svc := NewService(q)
	ctrl := NewController(svc, log)

	exercisesGroup := router.Group("/exercises")
	{
		exercisesGroup.POST("/create", ctrl.Create)
		exercisesGroup.PUT("/update", ctrl.Update)
		exercisesGroup.DELETE("/delete/:uuid", ctrl.Delete)
		exercisesGroup.POST("/restore", ctrl.Restore)
		exercisesGroup.GET("", ctrl.List)
		exercisesGroup.GET("/:uuid", ctrl.Get)
	}
}
