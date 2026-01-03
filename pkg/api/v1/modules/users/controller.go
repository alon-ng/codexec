package users

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
// @Summary      Create a new user
// @Description  Create a new user (admin function, can set is_admin)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        user  body      CreateUserRequest  true  "User creation data"
// @Success      201   {object}  db.User
// @Failure      400   {object}  errors.ErrorResponse
// @Failure      401   {object}  errors.ErrorResponse
// @Failure      500   {object}  errors.ErrorResponse
// @Router       /users/create [post]
func (c *Controller) Create(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid request data"), http.StatusBadRequest)
		return
	}

	user, err := c.svc.Create(ctx.Request.Context(), req)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(201, user)
}

// Update godoc
// @Summary      Update a user
// @Description  Update an existing user by UUID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        user  body      UpdateRequest  true  "User update data"
// @Success      200   {object}  db.User
// @Failure      400   {object}  errors.ErrorResponse
// @Failure      401   {object}  errors.ErrorResponse
// @Failure      500   {object}  errors.ErrorResponse
// @Router       /users/update [put]
func (c *Controller) Update(ctx *gin.Context) {
	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid request data"), http.StatusBadRequest)
		return
	}

	user, err := c.svc.Update(ctx.Request.Context(), req.Uuid, req.UpdateUserRequest)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, user)
}

type IDRequest struct {
	Uuid uuid.UUID `json:"uuid" binding:"required"`
}

// Delete godoc
// @Summary      Delete a user
// @Description  Soft delete a user by UUID (sets deleted_at timestamp)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        uuid  path      string  true  "User UUID"
// @Success      200   {string}  string  "OK"
// @Failure      400   {object}  errors.ErrorResponse
// @Failure      401   {object}  errors.ErrorResponse
// @Failure      500   {object}  errors.ErrorResponse
// @Router       /users/delete/{uuid} [delete]
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
// @Summary      Restore a deleted user
// @Description  Restore a soft-deleted user by UUID (sets deleted_at to NULL)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        user  body      IDRequest  true  "User UUID"
// @Success      200   {string}  string     "OK"
// @Failure      400   {object}  errors.ErrorResponse
// @Failure      401   {object}  errors.ErrorResponse
// @Failure      500   {object}  errors.ErrorResponse
// @Router       /users/restore [post]
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
// @Summary      List users
// @Description  Get a paginated list of users
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        limit   query     int     false  "Limit (default: 10)"  default(10)
// @Param        offset  query     int     false  "Offset (default: 0)"  default(0)
// @Success      200     {array}   db.User
// @Failure      400     {object}  errors.ErrorResponse
// @Failure      401     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /users [get]
func (c *Controller) List(ctx *gin.Context) {
	var req ListUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid query parameters"), http.StatusBadRequest)
		return
	}

	users, err := c.svc.List(ctx.Request.Context(), req)
	if err != nil {
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, users)
}

// Get godoc
// @Summary      Get a user by UUID
// @Description  Get a single user by UUID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        uuid  path      string  true  "User UUID"
// @Success      200   {object}  db.User
// @Failure      400   {object}  errors.ErrorResponse
// @Failure      401   {object}  errors.ErrorResponse
// @Failure      404   {object}  errors.ErrorResponse
// @Failure      500   {object}  errors.ErrorResponse
// @Router       /users/{uuid} [get]
func (c *Controller) Get(ctx *gin.Context) {
	idStr := ctx.Param("uuid")
	id, parseErr := uuid.Parse(idStr)
	if parseErr != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(parseErr, "Invalid UUID format"), http.StatusBadRequest)
		return
	}

	user, err := c.svc.Get(ctx.Request.Context(), id)
	if err != nil {
		if err.ErrorMessage == ErrUserNotFound {
			e.HandleError(ctx, c.log, err, http.StatusNotFound)
			return
		}

		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, user)
}
