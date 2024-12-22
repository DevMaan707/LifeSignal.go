package helpers

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func GetSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("JWT_SECRET_KEY environment variable is not set")
	}
	return secretKey
}
func GenerateJWT(userID string, expiresAt time.Time) (string, error) {
	claims := &jwt.StandardClaims{
		Issuer:    userID,
		ExpiresAt: expiresAt.Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := GetSecretKey()

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return tokenString, nil
}
func ParseJWT(tokenString string) (*jwt.StandardClaims, error) {
	secretKey := GetSecretKey()

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("incorrect password: %w", err)
	}
	return nil
}

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(1000000)
	return fmt.Sprintf("%06d", otp)
}

func VerifyOTP(storedOtp, enteredOtp string) bool {
	return storedOtp == enteredOtp
}

func SaveOTP(phone string, otp string) {
	log.Printf("Simulating save: OTP %s saved for phone %s", otp, phone)
}

func RetrieveOTP(phone string) (string, error) {
	return "123456", nil
}
func ValidateOTPExpiry(generatedAt time.Time, expiryDuration time.Duration) bool {
	return time.Since(generatedAt) <= expiryDuration
}
