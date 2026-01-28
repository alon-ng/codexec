package me

import (
	"codim/pkg/ai"
	"codim/pkg/db"
	"codim/pkg/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *gin.RouterGroup, q *db.Queries, p *pgxpool.Pool, log *logger.Logger, aiClient *ai.Client) {
	svc := NewService(q, p, aiClient)
	ctrl := NewController(svc, log)

	meGroup := router.Group("/me")
	{
		meGroup.GET("/", ctrl.Me)
		meGroup.GET("/courses", ctrl.ListUserCoursesWithProgress)
		meGroup.GET("/courses/:course_uuid", ctrl.GetUserCourseFull)
		meGroup.GET("/exercises/:exercise_uuid", ctrl.GetUserExercise)
		meGroup.PUT("/exercises/:exercise_uuid", ctrl.SaveUserExerciseSubmission)
		meGroup.GET("/exercises/:exercise_uuid/chat", ctrl.ListChatMessages)
		meGroup.POST("/exercises/:exercise_uuid/chat", ctrl.SendChatMessage)
	}

}
