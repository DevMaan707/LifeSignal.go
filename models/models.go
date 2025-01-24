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
	PhoneNumber string `json:"phone_number" validate:"required,phone"`
	Otp         string `json:"otp"`
}

type OTPRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
}

type VerifyOTPRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

type UserDetails struct {
	ID           string    `json:"id" bson:"_id"`
	Username     string    `json:"username" bson:"username"`
	Email        string    `json:"email" bson:"email"`
	Phone        string    `json:"phone" bson:"phone"`
	FirstName    string    `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName     string    `json:"last_name,omitempty" bson:"last_name,omitempty"`
	PasswordHash string    `json:"password_hash" bson:"password_hash"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}
type MedicalHistory struct {
	ID            string         `json:"id" bson:"_id"`
	UserID        string         `json:"user_id" bson:"user_id"`
	MedicalIssues []Issue        `json:"medical_issues" bson:"medical_issues"`
	Prescriptions []Prescription `json:"prescriptions" bson:"prescriptions"`
	Appointments  []Appointment  `json:"appointments" bson:"appointments"`
	CreatedAt     time.Time      `json:"created_at" bson:"created_at"`
}

type Issue struct {
	Condition string     `json:"condition" bson:"condition"`
	Severity  string     `json:"severity" bson:"severity"`
	Notes     string     `json:"notes" bson:"notes"`
	StartDate time.Time  `json:"start_date" bson:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty" bson:"end_date,omitempty"`
}

type Prescription struct {
	MedicationName string     `json:"medication_name" bson:"medication_name"`
	Dosage         string     `json:"dosage" bson:"dosage"`
	StartDate      time.Time  `json:"start_date" bson:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty" bson:"end_date,omitempty"`
}

type Appointment struct {
	DoctorID        string    `json:"doctor_id" bson:"doctor_id"`
	DoctorName      string    `json:"doctor_name" bson:"doctor_name"`
	AppointmentDate time.Time `json:"appointment_date" bson:"appointment_date"`
	Notes           string    `json:"notes" bson:"notes"`
}
type Doctor struct {
	ID             string    `json:"id" bson:"_id"`
	FirstName      string    `json:"first_name" bson:"first_name"`
	LastName       string    `json:"last_name" bson:"last_name"`
	Speciality     string    `json:"speciality" bson:"speciality"`
	Phone          string    `json:"phone" bson:"phone"`
	Email          string    `json:"email" bson:"email"`
	ClinicName     string    `json:"clinic_name" bson:"clinic_name"`
	ClinicAddress  string    `json:"clinic_address" bson:"clinic_address"`
	ProfilePicture string    `json:"profile_picture" bson:"profile_picture"`
	Rating         float32   `json:"rating" bson:"rating"`
	Experience     int       `json:"experience" bson:"experience"`
	Availability   []string  `json:"availability" bson:"availability"`
	Fee            float64   `json:"fee" bson:"fee"`
	Languages      []string  `json:"languages" bson:"languages"`
	Qualifications []string  `json:"qualifications" bson:"qualifications"`
	Services       []string  `json:"services" bson:"services"`
	About          string    `json:"about" bson:"about"`
	SocialLinks    Social    `json:"social_links" bson:"social_links"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
}

type Social struct {
	LinkedIn  string `json:"linkedin,omitempty" bson:"linkedin,omitempty"`
	Twitter   string `json:"twitter,omitempty" bson:"twitter,omitempty"`
	Facebook  string `json:"facebook,omitempty" bson:"facebook,omitempty"`
	Instagram string `json:"instagram,omitempty" bson:"instagram,omitempty"`
	Website   string `json:"website,omitempty" bson:"website,omitempty"`
}
