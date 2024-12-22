package handlers

import (
	"fmt"
	"life-signal/database"
	"life-signal/helpers"
	"life-signal/models"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func GenerateRandomDoctor(c *gin.Context, db *mongo.Client) {
	rand.Seed(time.Now().UnixNano())

	firstNames := []string{"John", "Jane", "Chris", "Pat", "Alex"}
	lastNames := []string{"Smith", "Doe", "Taylor", "Brown", "Johnson"}
	specialities := []string{"Cardiologist", "Dermatologist", "Neurologist", "Pediatrician", "General Practitioner"}
	clinicNames := []string{"City Health Clinic", "Wellness Center", "Prime Care", "Health Plus"}
	clinicAddresses := []string{"123 Main St", "456 Oak Ave", "789 Pine Dr", "321 Maple Rd"}
	availability := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}
	languages := []string{"English", "Spanish", "French", "German", "Chinese"}
	qualifications := []string{"MBBS", "MD", "PhD", "DO", "FACS"}
	services := []string{"Checkups", "Surgeries", "Consultations", "Therapy"}
	profilePictures := []string{
		"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Ftse1.mm.bing.net%2Fth%3Fid%3DOIP.JW_4m4RVV4ywf0aiB6TWrgHaLH%26pid%3DApi&f=1&ipt=dd58503b66ddab2a90e68b839f95d9de46e096c5f8f37e0d884c2fbdbab00a48&ipo=images",
		"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fimages.pexels.com%2Fphotos%2F433635%2Fpexels-photo-433635.jpeg%3Fcs%3Dsrgb%26dl%3Dman-person-portrait-433635.jpg%26fm%3Djpg&f=1&nofb=1&ipt=17b17bb9cfb36962e00a22945b4ced1193a1b4776389faaa756a1066f992384f&ipo=images",
		"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fwallpapers.com%2Fimages%2Fhd%2Fdoctor-pictures-l5y1qs2998u7rf0x.jpg&f=1&nofb=1&ipt=c31e94270e3f49f2ea16e9b7eb7e2478b74a5e742c0c17cb11ca301b9306a923&ipo=images",
		"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fleman-clinic.ch%2Fwp-content%2Fuploads%2F2018%2F11%2F02.jpg&f=1&nofb=1&ipt=21cbc1542e3deb2204bb116537903fc06bbf4766151a1867997821c8830f202b&ipo=images",
		"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fsegurancadotrabalhosempre.com%2Fwp-content%2Fuploads%2F2016%2F08%2FO6T8LS01.jpg&f=1&nofb=1&ipt=3fc23700e0428c6dcce21bd76eacca7034eabafd0c2e884ae2b79eab90a513e4&ipo=images",
		"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Ftse1.mm.bing.net%2Fth%3Fid%3DOIP.IVwf85npYYUcwRp4EIhqDgHaJm%26pid%3DApi&f=1&ipt=8b1382d10fd0e96c1a527ffc203d4bc2eec4bfcb3040a3a9138e18006e2e7504&ipo=images",
		"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fwww.advinohealthcare.com%2Fwp-content%2Fuploads%2F2020%2F08%2Fshutterstock_155685458.jpg&f=1&nofb=1&ipt=201a1977f129ca7d7068384bd60ec215ed83763aee70f43d2ba94df5fac5b45b&ipo=images",
		"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fpainlesshire.com%2Fwp-content%2Fuploads%2F2017%2F07%2Fdoctor.jpg&f=1&nofb=1&ipt=afd5f0341958babab3f3eb6bc3c3269d93db194b20fa990ab370c8ec424b422e&ipo=images",
		"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fimg.freepik.com%2Fpremium-photo%2Fmedical-concept-indian-beautiful-female-doctor-white-coat-with-stethoscope-waist-up-medical-student-woman-hospital-worker-looking-camera-smiling-studio-blue-background_185696-621.jpg%3Fw%3D2000&f=1&nofb=1&ipt=a134720cc6cf23db983914561e17706de1289e0b86f66022634cc92200f4f457&ipo=images",
		"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Ftse1.mm.bing.net%2Fth%3Fid%3DOIP.P1IfJNdtz7GmKkfPqR2yNAHaIO%26pid%3DApi&f=1&ipt=a362572be5e9b50b4e5eb90e9645e43c502cb622daaccd5a34a30b17bfc81c4f&ipo=images",
	}

	randString := func(slice []string) string {
		return slice[rand.Intn(len(slice))]
	}

	randInt := func(min, max int) int {
		return rand.Intn(max-min+1) + min
	}

	randFloat := func(min, max float64) float64 {
		return min + rand.Float64()*(max-min)
	}

	randSubset := func(slice []string, count int) []string {
		if count > len(slice) {
			count = len(slice)
		}
		result := make([]string, 0, count)
		perm := rand.Perm(len(slice))
		for i := 0; i < count; i++ {
			result = append(result, slice[perm[i]])
		}
		return result
	}

	doctor := models.Doctor{
		ID:             primitive.NewObjectID().Hex(),
		FirstName:      randString(firstNames),
		LastName:       randString(lastNames),
		Speciality:     randString(specialities),
		Phone:          fmt.Sprintf("+1-%03d-%03d-%04d", randInt(100, 999), randInt(100, 999), randInt(1000, 9999)),
		Email:          fmt.Sprintf("%s.%s@example.com", strings.ToLower(randString(firstNames)), strings.ToLower(randString(lastNames))),
		ClinicName:     randString(clinicNames),
		ClinicAddress:  randString(clinicAddresses),
		ProfilePicture: randString(profilePictures),
		Rating:         float32(randFloat(3.5, 5.0)),
		Experience:     randInt(1, 30),
		Availability:   randSubset(availability, randInt(2, 5)),
		Fee:            randFloat(50, 500),
		Languages:      randSubset(languages, randInt(1, 3)),
		Qualifications: randSubset(qualifications, randInt(1, 3)),
		Services:       randSubset(services, randInt(1, 4)),
		About:          "Highly experienced doctor dedicated to patient care.",
		SocialLinks: models.Social{
			LinkedIn:  "",
			Twitter:   "",
			Facebook:  "",
			Instagram: "",
			Website:   "",
		},
		CreatedAt: time.Now(),
	}

	doctorCollection := database.GetCollection(db, "life-signal", "doctors")
	if _, err := doctorCollection.InsertOne(c, doctor); err != nil {
		slog.Error("Failed to add doctor", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add doctor"})
		return
	}

	slog.Info("Doctor added successfully", "doctorID", doctor.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Doctor added successfully", "doctor": doctor})
}

func GenerateUserMedicalHistory(c *gin.Context, db *mongo.Client) {
	userID := c.Param("userid")
	if userID == "" {
		slog.Error("UserID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserID is required"})
		return
	}
	rand.Seed(time.Now().UnixNano())
	conditions := []string{"Hypertension", "Diabetes", "Asthma", "Arthritis", "Migraine"}
	severities := []string{"Mild", "Moderate", "Severe"}
	medications := []string{"Paracetamol", "Ibuprofen", "Metformin", "Insulin", "Aspirin"}
	doctors := []struct {
		ID   string
		Name string
	}{
		{"doc1", "Dr. John Smith"},
		{"doc2", "Dr. Jane Doe"},
		{"doc3", "Dr. Alice Brown"},
		{"doc4", "Dr. Bob Taylor"},
	}
	randString := func(slice []string) string {
		return slice[rand.Intn(len(slice))]
	}
	randInt := func(min, max int) int {
		return rand.Intn(max-min+1) + min
	}
	randTime := func() time.Time {
		return time.Now().AddDate(0, -randInt(1, 12), -randInt(0, 30))
	}
	randPtrTime := func() *time.Time {
		if rand.Intn(2) == 0 {
			return nil
		}
		endTime := randTime()
		return &endTime
	}
	issues := make([]models.Issue, randInt(1, 5))
	for i := range issues {
		issues[i] = models.Issue{
			Condition: randString(conditions),
			Severity:  randString(severities),
			Notes:     "Generated medical history notes.",
			StartDate: randTime(),
			EndDate:   randPtrTime(),
		}
	}
	prescriptions := make([]models.Prescription, randInt(1, 5))
	for i := range prescriptions {
		prescriptions[i] = models.Prescription{
			MedicationName: randString(medications),
			Dosage:         fmt.Sprintf("%d mg", randInt(100, 500)),
			StartDate:      randTime(),
			EndDate:        randPtrTime(),
		}
	}
	appointments := make([]models.Appointment, randInt(1, 3))
	for i := range appointments {
		doctor := doctors[rand.Intn(len(doctors))]
		appointments[i] = models.Appointment{
			DoctorID:        doctor.ID,
			DoctorName:      doctor.Name,
			AppointmentDate: randTime(),
			Notes:           "Routine checkup.",
		}
	}
	medicalHistory := models.MedicalHistory{
		ID:            primitive.NewObjectID().Hex(),
		UserID:        userID,
		MedicalIssues: issues,
		Prescriptions: prescriptions,
		Appointments:  appointments,
		CreatedAt:     time.Now(),
	}

	historyCollection := database.GetCollection(db, "life-signal", "user-medical-history")
	if _, err := historyCollection.InsertOne(c, medicalHistory); err != nil {
		slog.Error("Failed to add medical history", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add medical history"})
		return
	}

	slog.Info("Medical history added successfully", "userID", userID)
	c.JSON(http.StatusOK, gin.H{"message": "Medical history added successfully", "medical_history": medicalHistory})
}
func GetUserMedicalHistory(c *gin.Context, db *mongo.Client) {
	userID := c.Param("userid")
	if userID == "" {
		slog.Error("UserID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserID is required"})
		return
	}

	historyCollection := database.GetCollection(db, "life-signal", "user-medical-history")
	var medicalHistory models.MedicalHistory

	err := historyCollection.FindOne(c, bson.M{"user_id": userID}).Decode(&medicalHistory)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			slog.Warn("No medical history found", "userID", userID)
			c.JSON(http.StatusNotFound, gin.H{"error": "No medical history found"})
		} else {
			slog.Error("Error fetching medical history", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	slog.Info("Medical history retrieved successfully", "userID", userID)
	c.JSON(http.StatusOK, gin.H{"medical_history": medicalHistory})
}

func GetAllDoctors(c *gin.Context, db *mongo.Client) {
	doctorCollection := database.GetCollection(db, "life-signal", "doctors")
	cursor, err := doctorCollection.Find(c, bson.M{})
	if err != nil {
		slog.Error("Error fetching doctors", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer cursor.Close(c)

	var doctors []models.Doctor
	if err = cursor.All(c, &doctors); err != nil {
		slog.Error("Error decoding doctors", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	slog.Info("Doctors retrieved successfully", "count", len(doctors))
	c.JSON(http.StatusOK, gin.H{"doctors": doctors})
}
