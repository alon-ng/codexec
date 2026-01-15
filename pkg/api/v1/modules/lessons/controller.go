package lessons

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
// @Summary      Create a new lesson
// @Description  Create a new lesson
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        lesson  body      CreateLessonRequest  true  "Lesson creation data"
// @Success      201     {object}  db.LessonWithTranslation
// @Failure      400     {object}  errors.ErrorResponse
// @Failure      401     {object}  errors.ErrorResponse
// @Failure      409     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /lessons/create [post]
func (c *Controller) Create(ctx *gin.Context) {
	var req CreateLessonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid request data"), http.StatusBadRequest)
		return
	}

	lesson, err := c.svc.Create(ctx.Request.Context(), req)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(201, lesson)
}

// Update godoc
// @Summary      Update a lesson
// @Description  Update an existing lesson by UUID
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        lesson  body      UpdateRequest  true  "Lesson update data"
// @Success      200     {object}  db.LessonWithTranslation
// @Failure      400     {object}  errors.ErrorResponse
// @Failure      401     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /lessons/update [put]
func (c *Controller) Update(ctx *gin.Context) {
	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid request data"), http.StatusBadRequest)
		return
	}

	lesson, err := c.svc.Update(ctx.Request.Context(), req.Uuid, req.UpdateLessonRequest)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, lesson)
}

type IDRequest struct {
	Uuid uuid.UUID `json:"uuid" binding:"required"`
}

// Delete godoc
// @Summary      Delete a lesson
// @Description  Soft delete a lesson by UUID (sets deleted_at timestamp)
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        uuid  path      string  true  "Lesson UUID"
// @Success      200   {string}  string  "OK"
// @Failure      400   {object}  errors.ErrorResponse
// @Failure      401   {object}  errors.ErrorResponse
// @Failure      500   {object}  errors.ErrorResponse
// @Router       /lessons/delete/{uuid} [delete]
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
// @Summary      Restore a deleted lesson
// @Description  Restore a soft-deleted lesson by UUID (sets deleted_at to NULL)
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        lesson  body      IDRequest  true  "Lesson UUID"
// @Success      200     {string}  string     "OK"
// @Failure      400     {object}  errors.ErrorResponse
// @Failure      401     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /lessons/restore [post]
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
// @Summary      List lessons
// @Description  Get a paginated list of lessons, optionally filtered by course UUID
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        limit       query     int     false  "Limit (default: 10)"  default(10)
// @Param        offset      query     int     false  "Offset (default: 0)"  default(0)
// @Param        course_uuid query     string  false  "Filter by course UUID"
// @Success      200         {array}   db.LessonWithTranslation
// @Failure      400         {object}  errors.ErrorResponse
// @Failure      401         {object}  errors.ErrorResponse
// @Failure      500         {object}  errors.ErrorResponse
// @Router       /lessons [get]
func (c *Controller) List(ctx *gin.Context) {
	var req ListLessonsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid query parameters"), http.StatusBadRequest)
		return
	}

	lessons, err := c.svc.List(ctx.Request.Context(), req)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, lessons)
}

// Get godoc
// @Summary      Get a lesson by UUID
// @Description  Get a single lesson by UUID
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        uuid       path      string  true  	"Lesson UUID"
// @Param        language   query     string  false  	"Language"  	default(en) example(en)
// @Success      200   {object}  db.LessonWithTranslation
// @Failure      400   {object}  errors.ErrorResponse
// @Failure      401   {object}  errors.ErrorResponse
// @Failure      404   {object}  errors.ErrorResponse
// @Failure      500   {object}  errors.ErrorResponse
// @Router       /lessons/{uuid} [get]
func (c *Controller) Get(ctx *gin.Context) {
	idStr := ctx.Param("uuid")
	id, parseErr := uuid.Parse(idStr)
	if parseErr != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(parseErr, "Invalid UUID format"), http.StatusBadRequest)
		return
	}

	language := ctx.Query("language")
	if language == "" {
		language = "en"
	}

	lesson, err := c.svc.Get(ctx.Request.Context(), id, language)
	if err != nil {
		if err.ErrorMessage == ErrLessonNotFound {
			e.HandleError(ctx, c.log, err, http.StatusNotFound)
			return
		}

		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, lesson)
}

// AddTranslation godoc
// @Summary      Add a translation to an existing lesson
// @Description  Add a new translation for an existing lesson
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        translation  body      AddLessonTranslationRequest  true  "Translation data"
// @Success      201          {object}  db.LessonTranslation
// @Failure      400          {object}  errors.ErrorResponse
// @Failure      401          {object}  errors.ErrorResponse
// @Failure      500          {object}  errors.ErrorResponse
// @Router       /lessons/add-translation [post]
func (c *Controller) AddTranslation(ctx *gin.Context) {
	var req AddLessonTranslationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid request data"), http.StatusBadRequest)
		return
	}

	translation, err := c.svc.AddTranslation(ctx.Request.Context(), req)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(201, translation)
}
