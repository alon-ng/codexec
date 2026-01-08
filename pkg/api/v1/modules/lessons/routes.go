package lessons

import (
	"codim/pkg/db"
	"codim/pkg/utils/logger"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, q *db.Queries, log *logger.Logger) {
	svc := NewService(q)
	ctrl := NewController(svc, log)

	lessonsGroup := router.Group("/lessons")
	{
		lessonsGroup.POST("/create", ctrl.Create)
		lessonsGroup.PUT("/update", ctrl.Update)
		lessonsGroup.DELETE("/delete/:uuid", ctrl.Delete)
		lessonsGroup.POST("/restore", ctrl.Restore)
		lessonsGroup.GET("", ctrl.List)
		lessonsGroup.GET("/:uuid", ctrl.Get)
	}
}
