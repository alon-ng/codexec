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
