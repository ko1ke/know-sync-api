package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ko1ke/know-sync-api/utils/auth_utils"
	"github.com/ko1ke/know-sync-api/utils/res_utils"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth_utils.TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, &res_utils.ErrObj{Message: err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
