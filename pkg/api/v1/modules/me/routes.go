package me

import (
	"codim/pkg/db"
	"codim/pkg/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *gin.RouterGroup, q *db.Queries, p *pgxpool.Pool, log *logger.Logger) {
	svc := NewService(q, p)
	ctrl := NewController(svc, log)

	meGroup := router.Group("/me")
	{
		meGroup.GET("/", ctrl.Me)
		meGroup.GET("/courses", ctrl.ListUserCoursesWithProgress)
	}

}
