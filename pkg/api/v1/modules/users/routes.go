package users

import (
	"codim/pkg/db"
	"codim/pkg/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *gin.RouterGroup, q *db.Queries, p *pgxpool.Pool, log *logger.Logger) {
	svc := NewService(q, p)
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
