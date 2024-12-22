package handlers

import (
	"life-signal/database"
	"life-signal/helpers"
	"life-signal/models"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(c *gin.Context, db *mongo.Client) {
	var payload models.CreateAccountReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("Registration failed: Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storedOtp, err := helpers.RetrieveOTP(payload.Phone)
	if err != nil || !helpers.VerifyOTP(storedOtp, payload.OTP) {
		slog.Warn("Registration failed: Invalid or expired OTP", "phone", payload.Phone)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	if payload.Password != payload.ConfirmPassword {
		slog.Warn("Registration failed: Passwords do not match", "email", payload.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	userID := uuid.New().String()
	passwordHash, err := helpers.HashPassword(payload.Password)
	if err != nil {
		slog.Error("Registration failed: Error hashing password", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.UserDetails{
		ID:           userID,
		Username:     payload.Username,
		Email:        payload.Email,
		Phone:        payload.Phone,
		FirstName:    payload.FirstName,
		LastName:     payload.LastName,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	userCollection := database.GetCollection(db, "life-signal", "users")
	_, err = userCollection.InsertOne(c, user)
	if err != nil {
		slog.Error("Registration failed: Error inserting user into database", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user into database"})
		return
	}

	token, err := helpers.GenerateJWT(userID, time.Now().Add(24*time.Hour))
	if err != nil {
		slog.Error("Registration failed: Error generating JWT", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	slog.Info("Registration successful", "userID", userID)
	c.JSON(http.StatusOK, gin.H{"token": token, "userId": userID})
}

func Login(c *gin.Context, db *mongo.Client) {
	var login models.LoginReq
	if err := c.BindJSON(&login); err != nil {
		slog.Error("Login failed: Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userCollection := database.GetCollection(db, "life-signal", "users")
	var user models.UserDetails
	err := userCollection.FindOne(c, bson.M{"email": login.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			slog.Warn("Login failed: No user found", "email", login.Email)
			c.JSON(http.StatusBadRequest, gin.H{"error": "No user found"})
		} else {
			slog.Error("Login failed: Error fetching user", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}

	if err := helpers.VerifyPassword(user.PasswordHash, login.Password); err != nil {
		slog.Warn("Login failed: Invalid password", "email", login.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	token, err := helpers.GenerateJWT(user.ID, time.Now().Add(24*time.Hour))
	if err != nil {
		slog.Error("Login failed: Error generating JWT", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Login successful", "userID", user.ID)
	c.JSON(http.StatusOK, gin.H{"token": token, "userID": user.ID})

}

func GetOtpHandler(c *gin.Context, db *mongo.Client) {
	var request models.OTPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		slog.Error("GetOtp failed: Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	otp := helpers.GenerateOTP()
	helpers.SaveOTP(request.Phone, otp)

	slog.Info("OTP sent to user", "phone", request.Phone)
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func VerifyOtpHandler(c *gin.Context, db *mongo.Client) {
	var request models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		slog.Error("VerifyOtp failed: Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storedOtp, err := helpers.RetrieveOTP(request.Phone)
	if err != nil || !helpers.VerifyOTP(storedOtp, request.OTP) {
		slog.Warn("VerifyOtp failed: Invalid or expired OTP", "phone", request.Phone)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	slog.Info("OTP verified successfully", "phone", request.Phone)
	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}
