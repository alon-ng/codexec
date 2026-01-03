package auth

import (
	authProvider "codim/pkg/api/auth"
	"codim/pkg/db"
	"codim/pkg/utils/logger"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, q *db.Queries, log *logger.Logger, authProvider *authProvider.Provider) {
	svc := NewService(q, authProvider)
	controller := NewController(svc, log, authProvider)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", controller.Signup)
		authGroup.POST("/login", controller.Login)
		authGroup.POST("/logout", controller.Logout)
	}
}
