package controllers

import "github.com/gin-gonic/gin"

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func SignIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func hashPassword(password string) string {
	return ""
}

func verifyPassword(userPassword string, providePassword string) (bool, string) {
	return false, ""
}
