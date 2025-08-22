package controllers

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"errors" // ✅ <-- Add this line
	"strings"
	"time"

	tokenGenerate "github.com/foodmngtapp/food-management-apps/src/common/helpers/tokenHandler"
	"github.com/foodmngtapp/food-management-apps/src/common/logger"
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/foodmngtapp/food-management-apps/src/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection
var validate = validator.New()

// getUserCollection returns the user collection with lazy initialization
func getUserCollection() *mongo.Collection {
	if userCollection == nil {
		client := database.DBInstance()
		userCollection = database.OpenCollection(client, "users")
	}
	return userCollection
}

// Centralized Response Structures
type APIResponse struct {
	Status         bool            `json:"status"`
	StatusCode     int             `json:"statusCode"`
	Message        string          `json:"message"`
	Data           interface{}     `json:"data,omitempty"`
	PaginationData *PaginationInfo `json:"paginationData,omitempty"`
	Error          string          `json:"error,omitempty"`
}

type PaginationInfo struct {
	Total           int64 `json:"total"`
	Page            int   `json:"page"`
	TotalPage       int   `json:"totalPage"`
	Limit           int   `json:"limit"`
	HasNextPage     bool  `json:"hasNextPage"`
	HasPreviousPage bool  `json:"hasPreviousPage"`
}

// Helper Response Functions
func sendSuccessResponse(c *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    message,
		Data:       data,
	}
	c.JSON(http.StatusOK, response)
}

func sendSuccessWithPagination(c *gin.Context, message string, data interface{}, total int64, page, limit int) {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	
	paginationData := &PaginationInfo{
		Total:           total,
		Page:            page,
		TotalPage:       totalPages,
		Limit:           limit,
		HasNextPage:     page < totalPages,
		HasPreviousPage: page > 1,
	}

	response := APIResponse{
		Status:         true,
		StatusCode:     http.StatusOK,
		Message:        message,
		Data:           data,
		PaginationData: paginationData,
	}
	c.JSON(http.StatusOK, response)
}

func sendErrorResponse(c *gin.Context, statusCode int, message string, errorDetail ...string) {
	errorMsg := ""
	if len(errorDetail) > 0 {
		errorMsg = errorDetail[0]
	}

	response := APIResponse{
		Status:     false,
		StatusCode: statusCode,
		Message:    message,
		Error:      errorMsg,
	}
	c.JSON(statusCode, response)
}

func sendCreatedResponse(c *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Status:     true,
		StatusCode: http.StatusCreated,
		Message:    message,
		Data:       data,
	}
	c.JSON(http.StatusCreated, response)
}

// Helper Functions
func HashPassword(password string) string {
	// Reduced bcrypt cost from 14 to 10 for better performance
	// Cost 10: ~50ms vs Cost 14: ~800ms while maintaining good security
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	checkValidPassword := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Invalid login credentials!")
		checkValidPassword = false
	}
	return checkValidPassword, msg
}

type Credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserUpdate , definds user update model
type UserUpdates struct {
	FirstName    string
	LastName     string
	LastActive   time.Time
	ProfileImage *string
}

type PasswordReset struct {
	Email           string `json:"email" binding:"required,email"`
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var mongoCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			sendErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
			return
		}

		// Validate the user input
		validationErr := validate.Struct(user)
		if validationErr != nil {
			logger.Log.Error("Error binding JSON", zap.Error(validationErr))
			sendErrorResponse(c, http.StatusBadRequest, "Validation failed", validationErr.Error())
			return
		}

		// Normalize email before querying
		*user.Email = strings.ToLower(*user.Email)
		// Check if the user already exists
		var result bson.M
		collection := getUserCollection()
		err := collection.FindOne(mongoCtx, bson.M{
			"$or": []bson.M{
				{"email": user.Email},
				{"mobileNumber": user.MobileNumber},
			},
		}).Decode(&result)

		if err == nil {
			sendErrorResponse(c, http.StatusConflict, "User already exists", "User already exists with this email or mobile number")
			return
		} else if err != mongo.ErrNoDocuments {
			log.Println(err)
			logger.Log.Error("Mongo error", zap.Error(err))
			sendErrorResponse(c, http.StatusInternalServerError, "Database error", "internal server error")
			return
		}

		// Proceed with user registration
		userID := primitive.NewObjectID()
		password := HashPassword(*user.Password)
		user.Password = &password
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.LastActive = time.Now()
		user.ID = userID
		user.UserID = userID.Hex()

		// Generate tokens
		token, refreshToken, err := tokenGenerate.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, user.UserID)
		if err != nil {
			log.Println(err)
			sendErrorResponse(c, http.StatusInternalServerError, "Authentication error", "Error generating tokens")
			return
		}

		user.Token = &token
		user.RefreshToken = &refreshToken

		// Insert user into the database
		collection = getUserCollection()
		_, insertErr := collection.InsertOne(mongoCtx, user)
		if insertErr != nil {
			log.Println(insertErr)
			logger.Log.Error("Getting error while creating user", zap.Error(err))
			sendErrorResponse(c, http.StatusInternalServerError, "Registration failed", "Error creating user, please try again")
			return
		}

		// Respond with user details
		userData := user
		userData.Password = nil
		sendCreatedResponse(c, "User signup successful", userData)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		
		// Optimized timeout from 100s to 5s for faster response
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			logger.Log.Error("JSON binding error", zap.Error(err), zap.Duration("duration", time.Since(startTime)))
			sendErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
			return
		}

		if user.Email == nil || user.Password == nil {
			logger.Log.Warn("Missing credentials", zap.Duration("duration", time.Since(startTime)))
			sendErrorResponse(c, http.StatusBadRequest, "Missing credentials", "Email and password are required")
			return
		}

		*user.Email = strings.ToLower(*user.Email)
		
		// Optimized database query with projection to fetch only needed fields
		var foundUser models.User
		collection := getUserCollection()
		
		projection := bson.M{
			"_id":         1,
			"firstname":   1,
			"lastname":    1,
			"email":       1,
			"password":    1,
			"profileType": 1,
			"userId":      1,
			"isVerify":    1,
			"isComplete":  1,
			"mobileNumber": 1,
			"location":    1,
			"createdAt":   1,
			"updatedAt":   1,
			"lastActive":  1,
		}
		
		queryStart := time.Now()
		err := collection.FindOne(ctx, 
			bson.M{"email": user.Email},
			options.FindOne().SetProjection(projection),
		).Decode(&foundUser)
		queryDuration := time.Since(queryStart)
		
		if err != nil {
			if err == mongo.ErrNoDocuments {
				logger.Log.Warn("User not found", 
					zap.String("email", *user.Email),
					zap.Duration("total_duration", time.Since(startTime)))
				sendErrorResponse(c, http.StatusUnauthorized, "Authentication failed", "Invalid credentials")
			} else {
				logger.Log.Error("Database query error", zap.Error(err))
				sendErrorResponse(c, http.StatusInternalServerError, "Service unavailable", "Authentication service temporarily unavailable")
			}
			return
		}

		// Optimized password verification
		passwordStart := time.Now()
		checkValidPassword, _ := VerifyPassword(*user.Password, *foundUser.Password)
		passwordDuration := time.Since(passwordStart)
		
		if !checkValidPassword {
			logger.Log.Warn("Invalid password attempt", 
				zap.String("email", *user.Email),
				zap.Duration("password_check_duration", passwordDuration),
				zap.Duration("total_duration", time.Since(startTime)))
			sendErrorResponse(c, http.StatusUnauthorized, "Authentication failed", "Invalid credentials")
			return
		}

		// Token generation
		tokenStart := time.Now()
		token, refreshToken, err := tokenGenerate.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.UserID)
		if err != nil {
			logger.Log.Error("Token generation failed", zap.Error(err))
			sendErrorResponse(c, http.StatusInternalServerError, "Authentication error", "Authentication service error")
			return
		}
		tokenDuration := time.Since(tokenStart)

		// Optimized token update - single atomic operation instead of separate UpdateAllTokens call
		updateStart := time.Now()
		updateResult := collection.FindOneAndUpdate(
			ctx,
			bson.M{"userId": foundUser.UserID},
			bson.M{
				"$set": bson.M{
					"token":        token,
					"refreshToken": refreshToken,
					"lastActive":   time.Now(),
					"updatedAt":    time.Now(),
				},
			},
			options.FindOneAndUpdate().SetProjection(bson.M{"_id": 1}),
		)
		if updateResult.Err() != nil {
			logger.Log.Error("Token update failed", zap.Error(updateResult.Err()))
		}
		updateDuration := time.Since(updateStart)

		userResponse := foundUser
		userResponse.Password = nil
		
		// Include fresh tokens in response
		userResponse.Token = &token
		userResponse.RefreshToken = &refreshToken

		totalDuration := time.Since(startTime)
		
		// Performance logging for monitoring
		logger.Log.Info("Login completed",
			zap.String("email", *foundUser.Email),
			zap.Duration("total_duration", totalDuration),
			zap.Duration("query_duration", queryDuration),
			zap.Duration("password_duration", passwordDuration),
			zap.Duration("token_duration", tokenDuration),
			zap.Duration("update_duration", updateDuration))

		responseData := gin.H{
			"user": userResponse,
			"tokens": gin.H{
				"accessToken":  token,
				"refreshToken": refreshToken,
			},
		}

		sendSuccessResponse(c, "User login successful", responseData)
	}
}

func GetAllUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters for pagination and search
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			sendErrorResponse(c, http.StatusBadRequest, "Invalid pagination", "Invalid page number")
			return
		}

		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil || limit < 1 {
			sendErrorResponse(c, http.StatusBadRequest, "Invalid pagination", "Invalid limit")
			return
		}

		// Fetch search query parameters (firstName and lastName)
		firstName := c.Query("firstName")
		lastName := c.Query("lastName")

		// Construct MongoDB query filter based on search parameters
		filter := bson.M{}
		if firstName != "" {
			filter["firstname"] = bson.M{"$regex": primitive.Regex{Pattern: firstName, Options: "i"}}
		}
		if lastName != "" {
			filter["lastname"] = bson.M{"$regex": primitive.Regex{Pattern: lastName, Options: "i"}}
		}

		// Count total number of documents (users)
		collection := getUserCollection()
		total, err := collection.CountDocuments(context.Background(), filter)
		if err != nil {
			sendErrorResponse(c, http.StatusInternalServerError, "Query failed", "Failed to count users")
			return
		}

		// Calculate skip count based on pagination parameters
		skip := (page - 1) * limit

		// Fetch users with pagination and search filters applied
		cur, err := collection.Find(context.Background(), filter, options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)))
		if err != nil {
			sendErrorResponse(c, http.StatusInternalServerError, "Query failed", "Failed to fetch users")
			return
		}
		defer cur.Close(context.Background())

		var users []models.User
		if err := cur.All(context.Background(), &users); err != nil {
			sendErrorResponse(c, http.StatusInternalServerError, "Decode failed", "Failed to decode users")
			return
		}

		// Hide password field in each user
		for i := range users {
			users[i].Password = nil
		}

		sendSuccessWithPagination(c, "Users retrieved successfully", users, total, page, limit)
	}
}

func GetUserDetails(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		
		// Step 1: Extract userID from context (injected by JWT middleware)
		userIDVal, exists := c.Get("userID")
		if !exists {
			logger.Log.Error("UserID not found in token context")
			sendErrorResponse(c, http.StatusUnauthorized, "Authentication required", "userID not found in token")
			return
		}

		userIDStr, ok := userIDVal.(string)
		if !ok {
			logger.Log.Error("Invalid userID format in context", zap.Any("userID", userIDVal))
			sendErrorResponse(c, http.StatusBadRequest, "Invalid token", "Invalid userID format in token")
			return
		}

		// Step 2: Optimized user fetch with projection (only get needed fields)
		var user models.User
		collection := getUserCollection()
		
		// Use timeout context for better performance
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		// Projection to fetch only required fields
		projection := bson.M{
			"_id":           1,
			"firstname":     1,
			"lastname":      1,
			"email":         1,
			"profileType":   1,
			"userId":        1,
			"isVerify":      1,
			"isComplete":    1,
			"mobileNumber":  1,
			"location":      1,
			"createdAt":     1,
			"updatedAt":     1,
			"lastActive":    1,
			"profileImage":  1,
		}
		
		queryStart := time.Now()
		err := collection.FindOne(ctx, 
			bson.M{"userId": userIDStr}, 
			options.FindOne().SetProjection(projection),
		).Decode(&user)
		queryDuration := time.Since(queryStart)
		
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				logger.Log.Warn("User not found", zap.String("userID", userIDStr))
				sendErrorResponse(c, http.StatusNotFound, "User not found", "User not found")
			} else {
				logger.Log.Error("Database error fetching user details", zap.Error(err), zap.String("userID", userIDStr))
				sendErrorResponse(c, http.StatusInternalServerError, "Database error", "Error fetching user details")
			}
			return
		}

		// Step 3: Optimized last active update (atomic operation)
		updateStart := time.Now()
		updateResult := collection.FindOneAndUpdate(
			ctx,
			bson.M{"userId": userIDStr},
			bson.M{
				"$set": bson.M{
					"lastActive": time.Now(),
					"updatedAt":  time.Now(),
				},
			},
			options.FindOneAndUpdate().SetProjection(bson.M{"_id": 1}),
		)
		updateDuration := time.Since(updateStart)
		
		if updateResult.Err() != nil {
			logger.Log.Error("Failed to update last active timestamp", zap.Error(updateResult.Err()))
		}

		// Step 4: Ensure sensitive fields are not included (already excluded by projection)
		user.Password = nil

		totalDuration := time.Since(startTime)
		
		// Performance logging
		logger.Log.Info("User details fetched",
			zap.String("userID", userIDStr),
			zap.String("email", *user.Email),
			zap.Duration("total_duration", totalDuration),
			zap.Duration("query_duration", queryDuration),
			zap.Duration("update_duration", updateDuration))

		sendSuccessResponse(c, "User details fetched successfully", user)
	}
}

func UpdateUserDetails(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("userId")
		if userID == "" {
			sendErrorResponse(c, http.StatusBadRequest, "Missing parameter", "user ID not found in query parameters")
			return
		}

		var userDetail models.User // Assuming models.User is your user model
		collection := getUserCollection()
		err := collection.FindOne(context.TODO(), bson.M{"userId": userID}).Decode(&userDetail)
		if err != nil {
			sendErrorResponse(c, http.StatusNotFound, "User not found", "User not found!")
			return
		}

		var updates UserUpdates
		if err := c.BindJSON(&updates); err != nil {
			sendErrorResponse(c, http.StatusBadRequest, "Invalid request", "invalid request body")
			return
		}

		updateDoc := bson.M{"$set": bson.M{}}

		// Set firstName if provided
		if updates.FirstName != "" {
			updateDoc["$set"].(bson.M)["firstname"] = updates.FirstName
		}

		// Set lastName if provided
		if updates.LastName != "" {
			updateDoc["$set"].(bson.M)["lastname"] = updates.LastName
		}

		// Set lastActive if provided
		if !updates.LastActive.IsZero() {
			updateDoc["$set"].(bson.M)["lastactive"] = time.Now()
		}

		// Set new field to true
		updateDoc["$set"].(bson.M)["new"] = true

		if len(updateDoc["$set"].(bson.M)) == 0 {
			sendErrorResponse(c, http.StatusBadRequest, "No updates", "no valid fields to update")
			return
		}

		filter := bson.M{"userid": userID}
		result, err := collection.UpdateOne(context.TODO(), filter, updateDoc)
		if err != nil {
			sendErrorResponse(c, http.StatusInternalServerError, "Update failed", "error updating user details")
			return
		}

		if result.ModifiedCount == 0 {
			sendSuccessResponse(c, "No changes made to the user", nil)
			return
		}

		sendSuccessResponse(c, "User updated successfully", nil)
	}
}

func ResetPassword(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request parameters
		var payload PasswordReset
		if err := c.ShouldBindJSON(&payload); err != nil {
			sendErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
			return
		}

		// Check if the user exists and credentials are valid
		var user models.User
		collection := getUserCollection()
		err := collection.FindOne(context.TODO(), bson.M{"email": payload.Email}).Decode(&user)
		if err != nil {
			sendErrorResponse(c, http.StatusNotFound, "User not found", "User not found!")
			return
		}

		// Verify current password
		if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(payload.CurrentPassword)); err != nil {
			logger.Log.Error("Invalid email or password", zap.Error(err))
			sendErrorResponse(c, http.StatusUnauthorized, "Authentication failed", "Invalid email or password")
			return
		}

		// Check if the new password is different from the previous password
		if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(payload.NewPassword)); err == nil {
			sendErrorResponse(c, http.StatusBadRequest, "Invalid password", "New password must be different from the previous password")
			return
		}

		// Update the password in the database
		hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			logger.Log.Error("Failed to generate password hash", zap.Error(err))
			sendErrorResponse(c, http.StatusInternalServerError, "Encryption failed", "Failed to generate password hash")
			return
		}

		_, err = collection.UpdateOne(
			context.TODO(),
			bson.M{"email": payload.Email},
			bson.M{"$set": bson.M{"password": hashedNewPassword}},
		)
		if err != nil {
			logger.Log.Error("Failed to update password", zap.Error(err))
			sendErrorResponse(c, http.StatusInternalServerError, "Update failed", "Failed to update password")
			return
		}

		sendSuccessResponse(c, "Password reset successfully", nil)
	}
}

func ForgetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user input (email or mobile number) from the request
		var req struct {
			Email        string `json:"email" binding:"omitempty,email"`
			MobileNumber string `json:"mobileNumber" binding:"omitempty"`
		}
		if err := c.BindJSON(&req); err != nil {
			sendErrorResponse(c, http.StatusBadRequest, "Invalid request", "Invalid request data")
			return
		}

		if req.Email == "" && req.MobileNumber == "" {
			sendErrorResponse(c, http.StatusBadRequest, "Missing data", "Either email or mobile number must be provided")
			return
		}

		// Find user by email or mobile number
		filter := bson.M{}
		if req.Email != "" {
			filter["email"] = req.Email
		}
		if req.MobileNumber != "" {
			filter["mobilenumber"] = req.MobileNumber
		}

		var user models.User
		collection := getUserCollection()
		err := collection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			sendErrorResponse(c, http.StatusNotFound, "User not found", "User not found")
			return
		}

		resetToken, err := generateResetToken(user.ID)
		if err != nil {
			sendErrorResponse(c, http.StatusInternalServerError, "Token generation failed", "Failed to generate reset token")
			return
		}
		// Send password reset email or SMS to the user with the reset token
		if req.Email != "" {
			err = sendPasswordResetEmail(req.Email, resetToken)
			if err != nil {
				sendErrorResponse(c, http.StatusInternalServerError, "Email failed", "Failed to send password reset email")
				return
			}
		} else if req.MobileNumber != "" {
			err = sendPasswordResetSMS(req.MobileNumber, resetToken)
			if err != nil {
				sendErrorResponse(c, http.StatusInternalServerError, "SMS failed", "Failed to send password reset SMS")
				return
			}
		}

		sendSuccessResponse(c, "Password reset instructions sent", nil)
	}
}

// generateResetToken generates a secure token for password reset
func generateResetToken(userID primitive.ObjectID) (string, error) {
	// Implement your logic to generate a secure token (e.g., using JWT)
	// Here, you can create a JWT token containing the user ID and a short expiration time
	// Ensure the token is securely generated with appropriate expiration time
	return "", nil
}

// sendPasswordResetEmail sends an email to the user with the password reset instructions and the reset token
func sendPasswordResetEmail(email string, resetToken string) error {
	// Implement your logic to send an email containing the password reset instructions and the reset token
	// Ensure you securely send the email and handle any errors gracefully
	return nil
}

func sendPasswordResetSMS(email string, resetToken string) error {
	// Implement your logic to send an email containing the password reset instructions and the reset token
	// Ensure you securely send the email and handle any errors gracefully
	return nil
}

func DeleteAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("userId")
		if userID == "" {
			sendErrorResponse(c, http.StatusBadRequest, "Missing parameter", "user ID not found in query parameters")
			return
		}

		// Delete the user account
		collection := getUserCollection()
		result, err := collection.DeleteOne(context.Background(), bson.M{"userid": userID})
		if err != nil {
			sendErrorResponse(c, http.StatusInternalServerError, "Delete failed", "Failed to delete user account")
			return
		}

		if result.DeletedCount == 0 {
			sendErrorResponse(c, http.StatusNotFound, "User not found", "User not found")
			return
		}

		sendSuccessResponse(c, "User deleted successfully", nil)
	}
}