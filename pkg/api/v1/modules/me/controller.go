package me

import (
	e "codim/pkg/api/v1/errors"
	_ "codim/pkg/db"
	"codim/pkg/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
	svc *Service
	log *logger.Logger
}

func NewController(svc *Service, log *logger.Logger) *Controller {
	return &Controller{svc: svc, log: log}
}

// Me godoc
// @Summary      Get the current user from the JWT token
// @Description  Get the current user from the JWT token
// @Tags         me
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Success      200   {object}  db.User
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
// @Success      200     {array}   db.UserCourseWithProgress
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
