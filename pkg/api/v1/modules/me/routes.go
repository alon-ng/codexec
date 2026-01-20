package me

import (
	"codim/pkg/api/v1/modules/progress"
	"codim/pkg/db"
	"codim/pkg/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *gin.RouterGroup, q *db.Queries, p *pgxpool.Pool, log *logger.Logger) {
	progressSvc := progress.NewService(q, p)
	svc := NewService(q, p, progressSvc)
	ctrl := NewController(svc, progressSvc, log)

	meGroup := router.Group("/me")
	{
		meGroup.GET("/", ctrl.Me)
		meGroup.GET("/courses", ctrl.ListUserCoursesWithProgress)
		meGroup.GET("/courses/:course_uuid", ctrl.GetUserCourseFull)
		meGroup.GET("/exercises/:exercise_uuid", ctrl.GetUserExercise)
		meGroup.PUT("/exercises/:exercise_uuid", ctrl.SaveUserExerciseSubmission)
	}

}
