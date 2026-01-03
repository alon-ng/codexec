package middleware

import (
	"codim/internal/api/auth"
	e "codim/internal/api/v1/errors"
	"codim/internal/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authProvider *auth.Provider, log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(auth.AuthCookieName)
		if err != nil {
			e.HandleError(c, log, e.NewAPIError(err, "Authentication cookie required"), http.StatusUnauthorized)
			c.Abort()
			return
		}

		if tokenString == "" {
			e.HandleError(c, log, e.NewAPIError(nil, "Authentication cookie required"), http.StatusUnauthorized)
			c.Abort()
			return
		}

		userUUID, renewalRequired, err := authProvider.VerifyToken(tokenString)
		if err != nil {
			e.HandleError(c, log, e.NewAPIError(err, "Invalid or expired token"), http.StatusUnauthorized)
			c.Abort()
			return
		}

		if renewalRequired {
			newToken, err := authProvider.GenerateToken(userUUID)
			if err != nil {
				e.HandleError(c, log, e.NewAPIError(err, "Failed to generate new token"), http.StatusInternalServerError)
				c.Abort()
				return
			}
			authProvider.SetTokenCookie(c, newToken)
		}

		// Store user UUID in context
		c.Set("user_uuid", userUUID)
		c.Next()
	}
}
