package middlewares

import (
	"Training/internal/common_middlewares/common_middlewares_errors"
	"Training/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("Role")
		if !exists {
			logger.Logger.Error("Error getting Role: %v", common_middlewares_errors.RoleNotFoundInContext)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": common_middlewares_errors.RoleNotFoundInContext.Error()})
			return
		}

		if role != "admin" {
			logger.Logger.Error("Error current user not admin: %v", common_middlewares_errors.CurrentUserNotAdmin)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": common_middlewares_errors.CurrentUserNotAdmin.Error()})
			return
		}

		c.Next()
	}
}
