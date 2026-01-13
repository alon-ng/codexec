package me

import (
	"codim/pkg/db"
	"codim/pkg/utils/logger"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, q *db.Queries, log *logger.Logger) {
	svc := NewService(q)
	ctrl := NewController(svc, log)

	meGroup := router.Group("/me")
	{
		meGroup.GET("/", ctrl.Me)
	}
}
