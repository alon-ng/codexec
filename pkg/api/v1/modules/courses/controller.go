package courses

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

// Create godoc
// @Summary      Create a new course
// @Description  Create a new course
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        course  body      CreateCourseRequest  true  "Course creation data"
// @Success      201     {object}  db.Course
// @Failure      400     {object}  errors.ErrorResponse
// @Failure      401     {object}  errors.ErrorResponse
// @Failure      409     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /courses/create [post]
func (c *Controller) Create(ctx *gin.Context) {
	var req CreateCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid request data"), http.StatusBadRequest)
		return
	}

	course, err := c.svc.Create(ctx.Request.Context(), req)
	if err != nil {
		if err.ErrorMessage == ErrCourseNameAlreadyExists {
			e.HandleError(ctx, c.log, err, http.StatusConflict)
			return
		}
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(201, course)
}

// Update godoc
// @Summary      Update a course
// @Description  Update an existing course by UUID
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        course  body      UpdateRequest  true  "Course update data"
// @Success      200     {object}  db.Course
// @Failure      400     {object}  errors.ErrorResponse
// @Failure      401     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /courses/update [put]
func (c *Controller) Update(ctx *gin.Context) {
	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid request data"), http.StatusBadRequest)
		return
	}

	course, err := c.svc.Update(ctx.Request.Context(), req.Uuid, req.UpdateCourseRequest)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, course)
}

type IDRequest struct {
	Uuid uuid.UUID `json:"uuid" binding:"required"`
}

// Delete godoc
// @Summary      Delete a course
// @Description  Soft delete a course by UUID (sets deleted_at timestamp)
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        uuid  path      string  true  "Course UUID"
// @Success      200   {string}  string  "OK"
// @Failure      400   {object}  errors.ErrorResponse
// @Failure      401   {object}  errors.ErrorResponse
// @Failure      500   {object}  errors.ErrorResponse
// @Router       /courses/delete/{uuid} [delete]
func (c *Controller) Delete(ctx *gin.Context) {
	idStr := ctx.Param("uuid")
	if idStr == "" {
		e.HandleError(ctx, c.log, e.NewAPIError(nil, "UUID is required"), http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid UUID format"), http.StatusBadRequest)
		return
	}

	if err := c.svc.Delete(ctx.Request.Context(), id); err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.Status(200)
}

// Restore godoc
// @Summary      Restore a deleted course
// @Description  Restore a soft-deleted course by UUID (sets deleted_at to NULL)
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        course  body      IDRequest  true  "Course UUID"
// @Success      200     {string}  string     "OK"
// @Failure      400     {object}  errors.ErrorResponse
// @Failure      401     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /courses/restore [post]
func (c *Controller) Restore(ctx *gin.Context) {
	var req IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid request data"), http.StatusBadRequest)
		return
	}

	if err := c.svc.Restore(ctx.Request.Context(), req.Uuid); err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.Status(200)
}

// List godoc
// @Summary      List courses
// @Description  Get a paginated list of courses, optionally filtered by subject
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        limit   query     int     false  "Limit (default: 10)"  default(10)
// @Param        offset  query     int     false  "Offset (default: 0)"  default(0)
// @Param        subject query     string  false  "Filter by subject"
// @Success      200     {array}   db.Course
// @Failure      400     {object}  errors.ErrorResponse
// @Failure      401     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /courses [get]
func (c *Controller) List(ctx *gin.Context) {
	var req ListCoursesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid query parameters"), http.StatusBadRequest)
		return
	}

	courses, err := c.svc.List(ctx.Request.Context(), req)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, courses)
}

// Get godoc
// @Summary      Get a course by UUID
// @Description  Get a single course by UUID
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        uuid  path      string  true  "Course UUID"
// @Success      200   {object}  db.CourseFull
// @Failure      400   {object}  errors.ErrorResponse
// @Failure      401   {object}  errors.ErrorResponse
// @Failure      404   {object}  errors.ErrorResponse
// @Failure      500   {object}  errors.ErrorResponse
// @Router       /courses/{uuid} [get]
func (c *Controller) Get(ctx *gin.Context) {
	idStr := ctx.Param("uuid")
	id, parseErr := uuid.Parse(idStr)
	if parseErr != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(parseErr, "Invalid UUID format"), http.StatusBadRequest)
		return
	}

	course, err := c.svc.Get(ctx.Request.Context(), id)
	if err != nil {
		if err.ErrorMessage == ErrCourseNotFound {
			e.HandleError(ctx, c.log, err, http.StatusNotFound)
			return
		}

		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, course)
}
