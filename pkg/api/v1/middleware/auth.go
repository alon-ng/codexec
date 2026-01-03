package middleware

import (
	"codim/pkg/api/auth"
	"codim/pkg/api/v1/cache"
	e "codim/pkg/api/v1/errors"
	"codim/pkg/db"
	"codim/pkg/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authProvider *auth.Provider, userCache *cache.UserCache, log *logger.Logger) gin.HandlerFunc {
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

		// Get user from cache (with fallback to DB and locking)
		user, err := userCache.GetUser(c.Request.Context(), userUUID)
		if err != nil {
			e.HandleError(c, log, e.NewAPIError(err, "Failed to get user"), http.StatusInternalServerError)
			c.Abort()
			return
		}

		// Store user UUID and user object in context
		c.Set("user_uuid", userUUID)
		c.Set("user", user)
		c.Next()
	}
}

// AdminMiddleware ensures the authenticated user is an admin
func AdminMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			e.HandleError(c, log, e.NewAPIError(nil, "User not found in context"), http.StatusInternalServerError)
			c.Abort()
			return
		}

		user, ok := userInterface.(db.User)
		if !ok {
			e.HandleError(c, log, e.NewAPIError(nil, "Invalid user type in context"), http.StatusInternalServerError)
			c.Abort()
			return
		}

		if !user.IsAdmin {
			e.HandleError(c, log, e.NewAPIError(nil, "Admin access required"), http.StatusForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
