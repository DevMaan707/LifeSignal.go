package routes

import (
	"life-signal/handlers"
	"life-signal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Routes(engine *gin.Engine, db *mongo.Client) {
	auth := engine.Group("/auth")
	{
		auth.POST("/login", func(c *gin.Context) { handlers.Login(c, db) })
		auth.POST("/signup", func(c *gin.Context) { handlers.Register(c, db) })
		auth.POST("/getOtp", func(c *gin.Context) { handlers.GetOtpHandler(c, db) })
		auth.POST("/verifyOtp", func(c *gin.Context) { handlers.VerifyOtpHandler(c, db) })
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
