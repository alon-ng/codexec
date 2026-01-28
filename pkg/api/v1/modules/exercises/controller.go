package exercises

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
// @Summary      Create a new exercise
// @Description  Create a new exercise
// @Tags         exercises
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        exercise  body      CreateExerciseRequest  true  "Exercise creation data"
// @Success      201       {object}  models.ExerciseWithTranslation
// @Failure      400       {object}  errors.ErrorResponse
// @Failure      401       {object}  errors.ErrorResponse
// @Failure      409       {object}  errors.ErrorResponse
// @Failure      500       {object}  errors.ErrorResponse
// @Router       /exercises/create [post]
func (c *Controller) Create(ctx *gin.Context) {
	var req CreateExerciseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid request data"), http.StatusBadRequest)
		return
	}

	exercise, err := c.svc.Create(ctx.Request.Context(), req)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(201, exercise)
}

// Update godoc
// @Summary      Update an exercise
// @Description  Update an existing exercise by UUID
// @Tags         exercises
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        exercise  body      UpdateExerciseRequest  true  "Exercise update data"
// @Success      200       {object}  models.ExerciseWithTranslation
// @Failure      400       {object}  errors.ErrorResponse
// @Failure      401       {object}  errors.ErrorResponse
// @Failure      500       {object}  errors.ErrorResponse
// @Router       /exercises/update [put]
func (c *Controller) Update(ctx *gin.Context) {
	var req UpdateExerciseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid request data"), http.StatusBadRequest)
		return
	}

	exercise, err := c.svc.Update(ctx.Request.Context(), req)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, exercise)
}

// Delete godoc
// @Summary      Delete an exercise
// @Description  Soft delete an exercise by UUID (sets deleted_at timestamp)
// @Tags         exercises
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        uuid  path      string  true  "Exercise UUID"
// @Success      200   {string}  string  "OK"
// @Failure      400   {object}  errors.ErrorResponse
// @Failure      401   {object}  errors.ErrorResponse
// @Failure      500   {object}  errors.ErrorResponse
// @Router       /exercises/delete/{uuid} [delete]
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
// @Summary      Restore a deleted exercise
// @Description  Restore a soft-deleted exercise by UUID (sets deleted_at to NULL)
// @Tags         exercises
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        exercise  body      IDRequest  true  "Exercise UUID"
// @Success      200       {string}  string     "OK"
// @Failure      400       {object}  errors.ErrorResponse
// @Failure      401       {object}  errors.ErrorResponse
// @Failure      500       {object}  errors.ErrorResponse
// @Router       /exercises/restore [post]
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
// @Summary      List exercises
// @Description  Get a paginated list of exercises, optionally filtered by lesson UUID
// @Tags         exercises
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        limit       query     int      false  "Limit (default: 10)"  default(10)
// @Param        offset      query     int      false  "Offset (default: 0)"  default(0)
// @Param        lesson_uuid query     string   false  "Filter by lesson UUID"
// @Param        language    query     string   false  "Filter by language"  default(en)
// @Success      200         {array}   models.ExerciseWithTranslation
// @Failure      400         {object}  errors.ErrorResponse
// @Failure      401         {object}  errors.ErrorResponse
// @Failure      500         {object}  errors.ErrorResponse
// @Router       /exercises [get]
func (c *Controller) List(ctx *gin.Context) {
	var req ListExercisesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid query parameters"), http.StatusBadRequest)
		return
	}

	exercises, err := c.svc.List(ctx.Request.Context(), req)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, exercises)
}

// Get godoc
// @Summary      Get an exercise by UUID
// @Description  Get a single exercise by UUID
// @Tags         exercises
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        uuid  		path      string  true  	"Exercise UUID"
// @Param        language   query     string  false  	"Language"  	default(en) example(en)
// @Success      200   {object}  models.ExerciseWithTranslation
// @Failure      400   {object}  errors.ErrorResponse
// @Failure      401   {object}  errors.ErrorResponse
// @Failure      404   {object}  errors.ErrorResponse
// @Failure      500   {object}  errors.ErrorResponse
// @Router       /exercises/{uuid} [get]
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

	exercise, err := c.svc.Get(ctx.Request.Context(), id, language)
	if err != nil {
		if err.ErrorMessage == ErrExerciseNotFound {
			e.HandleError(ctx, c.log, err, http.StatusNotFound)
			return
		}

		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, exercise)
}

// AddTranslation godoc
// @Summary      Add a translation to an existing exercise
// @Description  Add a new translation for an existing exercise
// @Tags         exercises
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        translation  body      AddExerciseTranslationRequest  true  "Translation data"
// @Success      201          {object}  models.ExerciseTranslation
// @Failure      400          {object}  errors.ErrorResponse
// @Failure      401          {object}  errors.ErrorResponse
// @Failure      500          {object}  errors.ErrorResponse
// @Router       /exercises/add-translation [post]
func (c *Controller) AddTranslation(ctx *gin.Context) {
	var req AddExerciseTranslationRequest
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
