package routes

import (
	"life-signal/handlers"
	"life-signal/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Routes(engine *gin.Engine, db *mongo.Client) {
	dev := engine.Group("/dev")
	{
		dev.GET("/generate-doc", func(c *gin.Context) {
			handlers.GenerateRandomDoctor(c, db)
		})
		dev.GET("/medical-history/generate/:userid", func(c *gin.Context) {
			handlers.GenerateUserMedicalHistory(c, db)
		})
	}
	auth := engine.Group("/auth")
	{
		auth.POST("/login", func(c *gin.Context) { handlers.Login(c, db) })
		auth.POST("/signup", func(c *gin.Context) { handlers.Register(c, db) })
		auth.POST("/getOtp", func(c *gin.Context) { handlers.GetOtpHandler(c, db) })
		auth.POST("/verifyOtp", func(c *gin.Context) { handlers.VerifyOtpHandler(c, db) })
	}

	protected := engine.Group("/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/get-doctors", func(c *gin.Context) {
			handlers.GetAllDoctors(c, db)
		})
		protected.GET("/get-medical-history/:userid", func(c *gin.Context) {
			handlers.GetUserMedicalHistory(c, db)
		})
		protected.GET("/get-user/:userid", func(c *gin.Context) {
			handlers.GetUserDetails(c, db)
		})
		protected.POST("/set-medical-history/:userid", func(c *gin.Context) {
			handlers.SetUserMedicalHistory(c, db)
		})
	}
}
