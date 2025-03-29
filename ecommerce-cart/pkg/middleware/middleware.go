package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	token "github.com/thutasann/ecommerce-cart/pkg/tokens"
)

// Authentication Middlware
func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ClientToken := ctx.Request.Header.Get("token")
		if ClientToken == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No authorization header provided"})
			ctx.Abort()
			return
		}
		claims, err := token.ValidateToken(ClientToken)
		if err != "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			ctx.Abort()
			return
		}
		ctx.Set("email", claims.Email)
		ctx.Set("uid", claims.Uid)
		ctx.Next()
	}
}
