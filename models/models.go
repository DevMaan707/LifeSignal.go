package models

import (
	"time"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type CreateAccountReq struct {
	Username        string `json:"username" validate:"required,min=3,max=32"`
	Email           string `json:"email" validate:"required,email"`
	Phone           string `json:"phone" validate:"required,e164"`
	FirstName       string `json:"first_name" validate:"omitempty,min=1,max=32"`
	LastName        string `json:"last_name" validate:"omitempty,min=1,max=32"`
	Password        string `json:"password" validate:"required,min=8,max=128"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=128,eqfield=Password"`
	OTP             string `json:"otp" validate:"omitempty,len=6"`
}

type LoginReq struct {
	Phone    string `json:"phone" validate:"required,e164"`
	Password string `json:"password"`
	OTP      string `json:"otp,omitempty"` 
}
type OTPRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
	OTP   string `json:"otp" validate:"required,len=6"`  
}

type VerifyOTPRequest struct {
	Phone string `json:"phone" validate:"required,e164"` 
	OTP   string `json:"otp" validate:"required,len=6"` 
}

type CreateAccountDetails struct {
	ID        string    `json:"id"`
	Username  string    `json:"username" validate:"required,min=3,max=32"`
	Email     string    `json:"email" validate:"required,email"`
	Phone     string    `json:"phone" validate:"required,e164"`
	FirstName string    `json:"first_name" validate:"omitempty,min=1,max=32"`
	LastName  string    `json:"last_name" validate:"omitempty,min=1,max=32"`
	Password  string    `json:"password" validate:"required,min=8,max=128"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	OTP       string    `json:"otp,omitempty"`
}
