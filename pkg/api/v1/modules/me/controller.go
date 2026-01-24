package me

import (
	e "codim/pkg/api/v1/errors"
	"codim/pkg/api/v1/modules/progress"
	"codim/pkg/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
	svc *Service
	log *logger.Logger
}

func NewController(svc *Service, progressSvc *progress.Service, log *logger.Logger) *Controller {
	return &Controller{svc: svc, log: log}
}

// Me godoc
// @Summary      Get the current user from the JWT token
// @Description  Get the current user from the JWT token
// @Tags         me
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Success      200   {object}  users.User
// @Failure      400   {object}  errors.ErrorResponse
// @Failure      401   {object}  errors.ErrorResponse
// @Failure      500   {object}  errors.ErrorResponse
// @Router       /me [get]
func (c *Controller) Me(ctx *gin.Context) {
	meUUID := uuid.MustParse(ctx.GetString("user_uuid"))
	user, err := c.svc.Me(ctx.Request.Context(), meUUID)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, user)
}

// ListUserCoursesWithProgress godoc
// @Summary      List the user courses with progress
// @Description  List the user courses with progress
// @Tags         me
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        limit   	query    int     false  "Limit (default: 10)" default(10)
// @Param        offset  	query    int     false  "Offset (default: 0)" default(0)
// @Param        subject 	query    string  false  "Filter by subject"
// @Param        language 	query    string  false  "Filter by language"  default(en) 	example(en)
// @Param        is_active 	query    bool    false  "Filter by is_active" default(true) example(true)
// @Success      200     {array}   progress.UserCourseWithProgress
// @Failure      400     {object}  errors.ErrorResponse
// @Failure      401     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /me/courses [get]
func (c *Controller) ListUserCoursesWithProgress(ctx *gin.Context) {
	meUUID := uuid.MustParse(ctx.GetString("user_uuid"))
	var req ListUserCoursesWithProgressRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid request data"), http.StatusBadRequest)
		return
	}

	userCourses, err := c.svc.ListUserCoursesWithProgress(ctx.Request.Context(), meUUID, req)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, userCourses)
}

// GetUserCourseFull godoc
// @Summary      Get the user course full
// @Description  Get the user course full
// @Tags         me
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        course_uuid path string true "Course UUID"
// @Success      200     {object}  progress.UserCourseFull
// @Failure      400     {object}  errors.ErrorResponse
// @Failure      401     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /me/courses/{course_uuid} [get]
func (c *Controller) GetUserCourseFull(ctx *gin.Context) {
	meUUID := uuid.MustParse(ctx.GetString("user_uuid"))
	courseUUID := uuid.MustParse(ctx.Param("course_uuid"))
	userCourse, err := c.svc.GetUserCourseFull(ctx.Request.Context(), meUUID, courseUUID)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, userCourse)
}

// GetUserExercise godoc
// @Summary      Get the user exercise
// @Description  Get the user exercise
// @Tags         me
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        exercise_uuid path string true "Exercise UUID"
// @Success      200     {object}  UserExercise
// @Failure      400     {object}  errors.ErrorResponse
// @Failure      401     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /me/exercises/{exercise_uuid} [get]
func (c *Controller) GetUserExercise(ctx *gin.Context) {
	meUUID := uuid.MustParse(ctx.GetString("user_uuid"))
	exerciseUUID := uuid.MustParse(ctx.Param("exercise_uuid"))
	userExercise, err := c.svc.GetUserExercise(ctx.Request.Context(), meUUID, exerciseUUID)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, userExercise)
}

// SaveUserExerciseSubmission godoc
// @Summary      Save the user exercise submission
// @Description  Save the user exercise submission
// @Tags         me
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        exercise_uuid path string true "Exercise UUID"
// @Param        submission body SaveUserExerciseSubmissionRequest true "Submission"
// @Success      200     {string}  string  "OK"
// @Failure      400     {object}  errors.ErrorResponse
// @Failure      401     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /me/exercises/{exercise_uuid} [put]
func (c *Controller) SaveUserExerciseSubmission(ctx *gin.Context) {
	meUUID := uuid.MustParse(ctx.GetString("user_uuid"))
	exerciseUUID := uuid.MustParse(ctx.Param("exercise_uuid"))
	var req SaveUserExerciseSubmissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid request data"), http.StatusBadRequest)
		return
	}

	err := c.svc.SaveUserExerciseSubmission(ctx.Request.Context(), meUUID, exerciseUUID, req)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
