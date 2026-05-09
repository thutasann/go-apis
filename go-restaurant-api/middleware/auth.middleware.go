package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/thutasann/restaurant-api/helpers"
)

// Auth Middleware
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			helpers.Error(c, "No authorization header provided", 403, c.Err())
			c.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			helpers.Error(c, "token validate error", 0, c.Err())
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)

		c.Next()
	}
}
