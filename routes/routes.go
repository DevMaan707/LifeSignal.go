package routes

import (
	"life-signal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine) {

	auth := engine.Group("/auth")
	{
		auth.POST("/login", loginHandler)
		auth.POST("/signup", signupHandler)
		auth.POST("/getOtp", getOtpHandler)
		auth.POST("/verifyOtp", verifyOtpHandler)
	}

	protected := engine.Group("/protected")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/dashboard", dashboardHandler)
	}
}

func loginHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func signupHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Signup successful"})
}

func getOtpHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent"})
}

func verifyOtpHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OTP verified"})
}

func dashboardHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the dashboard"})
}
