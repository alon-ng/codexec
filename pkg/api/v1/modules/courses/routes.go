package courses

import (
	"codim/pkg/db"
	"codim/pkg/utils/logger"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(router *gin.RouterGroup, q *db.Queries, log *logger.Logger) {
	svc := NewService(q)
	ctrl := NewController(svc, log)

	coursesGroup := router.Group("/courses")
	{
		coursesGroup.GET("", ctrl.List)
		coursesGroup.GET("/:uuid", ctrl.Get)
	}
}

func RegisterAdminRoutes(router *gin.RouterGroup, q *db.Queries, log *logger.Logger) {
	svc := NewService(q)
	ctrl := NewController(svc, log)

	coursesGroup := router.Group("/courses")
	{
		coursesGroup.POST("/create", ctrl.Create)
		coursesGroup.PUT("/update", ctrl.Update)
		coursesGroup.DELETE("/delete/:uuid", ctrl.Delete)
		coursesGroup.POST("/restore", ctrl.Restore)
	}
}
