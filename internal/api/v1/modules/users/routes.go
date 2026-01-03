package users

import (
	"codim/internal/db"
	"codim/internal/utils/logger"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, q *db.Queries, log *logger.Logger) {
	svc := NewService(q)
	ctrl := NewController(svc, log)

	usersGroup := router.Group("/users")
	{
		usersGroup.POST("/create", ctrl.Create)
		usersGroup.PUT("/update", ctrl.Update)
		usersGroup.DELETE("/delete/:uuid", ctrl.Delete)
		usersGroup.POST("/restore", ctrl.Restore)
		usersGroup.GET("", ctrl.List)
		usersGroup.GET("/:uuid", ctrl.Get)
	}
}
