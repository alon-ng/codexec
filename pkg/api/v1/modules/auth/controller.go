package auth

import (
	"codim/pkg/api/auth"
	e "codim/pkg/api/v1/errors"
	"codim/pkg/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	svc          *Service
	log          *logger.Logger
	authProvider *auth.Provider
}

func NewController(svc *Service, log *logger.Logger, authProvider *auth.Provider) *Controller {
	return &Controller{svc: svc, log: log, authProvider: authProvider}
}

// Signup godoc
// @Summary      Signup a new user
// @Description  Signup a new user. The authentication cookie (auth_token) is automatically set in the response.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body      SignupRequest  true  "Signup request"
// @Success      201
// @Failure      400     {object}  errors.ErrorResponse
// @Failure      409     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /auth/signup [post]
func (c *Controller) Signup(ctx *gin.Context) {
	var req SignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid request data"), http.StatusBadRequest)
		return
	}

	resp, err := c.svc.Signup(ctx.Request.Context(), req)
	if err != nil {
		if err.ErrorMessage == ErrEmailAlreadyExists {
			e.HandleError(ctx, c.log, err, http.StatusConflict)
			return
		}
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	c.authProvider.SetTokenCookie(ctx, resp.Token)

	ctx.Status(http.StatusCreated)
}

// Login godoc
// @Summary      Login a user
// @Description  Login a user. The authentication cookie (auth_token) is automatically set in the response.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body      LoginRequest  true  "Login request"
// @Success      200
// @Failure      400     {object}  errors.ErrorResponse
// @Failure      401     {object}  errors.ErrorResponse
// @Failure      500     {object}  errors.ErrorResponse
// @Router       /auth/login [post]
func (c *Controller) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		e.HandleError(ctx, c.log, e.NewAPIError(err, "Invalid request data"), http.StatusBadRequest)
		return
	}

	resp, err := c.svc.Login(ctx.Request.Context(), req)
	if err != nil {
		if err.ErrorMessage == ErrInvalidCredentials {
			e.HandleError(ctx, c.log, err, http.StatusUnauthorized)
			return
		}
		e.HandleError(ctx, c.log, err, http.StatusInternalServerError)
		return
	}

	c.authProvider.SetTokenCookie(ctx, resp.Token)
	ctx.Status(http.StatusOK)
}

// Logout godoc
// @Summary      Logout a user
// @Description  Logout a user by clearing the authentication cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200     {object}  map[string]string
// @Router       /auth/logout [post]
func (c *Controller) Logout(ctx *gin.Context) {
	c.authProvider.SetTokenCookie(ctx, "")

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
