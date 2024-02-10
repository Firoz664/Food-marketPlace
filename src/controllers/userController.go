package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

// var logger *zap.Logger

// Helper Fucntions
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
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
			c.JSON(http.StatusBadRequest, gin.H{
				"status":     false,
				"statusCode": http.StatusBadRequest,
				"error":      err.Error(),
			})
			return
		}

		// validate date type and key
		validationErr := validate.Struct(user)
		if validationErr != nil {
			logger.Log.Error("Error binding JSON", zap.Error(validationErr))

			c.JSON(http.StatusBadRequest, gin.H{
				"status":     false,
				"statusCode": http.StatusBadRequest,
				"error":      validationErr.Error(),
			})
			return
		}

		var result bson.M
		err := userCollection.FindOne(mongoCtx, bson.M{
			"$or": []bson.M{
				{"email": user.Email},
				{"mobilenumber": user.MobileNumber},
			},
		}).Decode(&result)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				// No user found with the same email or mobile, safe to proceed with user registration
			} else {
				// An error occurred during the query execution
				log.Println(err)
				logger.Log.Error("Mongo error", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":     false,
					"statusCode": http.StatusInternalServerError,
					"error":      "internal server error",
				})
				return
			}
		} else {
			// A user was found with the same email or mobile
			c.JSON(http.StatusConflict, gin.H{
				"status":     false,
				"statusCode": http.StatusConflict,
				"error":      "User already exits! with email or mobile",
			})
			return
		}

		userID := primitive.NewObjectID()
		// Assuming HashPassword and generate.TokenGenerator are implemented elsewhere and working correctly.
		password := HashPassword(*user.Password)
		user.Password = &password
		user.CreatedAt = time.Now() // Direct assignment, no need for parsing
		user.UpdatedAt = time.Now()
		user.LastActive = time.Now()
		user.ID = userID // Assuming ID is of type primitive.ObjectID
		user.UserID = userID.Hex()
		user.ID = primitive.NewObjectID()
		fmt.Println("password:", password)
		*user.Email = strings.ToLower(*user.Email)

		// // Handling errors from TokenGenerator
		token, refreshToken, err := tokenGenerate.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, user.UserID)
		if err != nil {
			log.Println(err) // Use log or c.JSON to return an error response
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":     false,
				"statusCode": http.StatusInternalServerError,
				"error":      "Error generating tokens",
			})
			return
		}

		user.Token = &token
		user.RefreshToken = &refreshToken
		_, insertErr := userCollection.InsertOne(c, user)
		if insertErr != nil {
			log.Println(insertErr) // It's a good practice to log the actual error too
			logger.Log.Error("Getting error while creating user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "false",
				"error":  "Getting error while creating user, please try again"})
			return
		}
		defer cancel()

		userData := user
		userData.Password = nil
		// Respond with user details, message, and custom status code
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"status":  http.StatusOK,
				"message": "User signup successfully",
				"result":  userData,
			},
		})

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate fetching a user
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var updates Credential
		if err := c.BindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		var user models.User
		fmt.Println("User", user)
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":     false,
				"statusCode": http.StatusBadRequest,
				"error":      err.Error(),
			})
			return
		}

		if user.Email == nil || *user.Email == "" || user.Password == nil || *user.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":     false,
				"statusCode": http.StatusBadRequest,
				"error":      "Email or Password missing",
			})
			return
		}

		*user.Email = strings.ToLower(*user.Email)
		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":     false,
				"statusCode": http.StatusNotFound,
				"error":      "User not exist!"})
			return
		}
		checkValidPassword, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !checkValidPassword {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"error":  msg})
			return
		}

		token, refreshToken, err := tokenGenerate.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		tokenGenerate.UpdateAllTokens(token, refreshToken, foundUser.UserID)

		userResponse := foundUser
		userResponse.Password = nil

		// Log fetching user operation
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"status":     true,
				"statusCode": http.StatusOK,
				"message":    "User login successfully",
				"result":     userResponse,
			},
		})

	}
}

func GetAllUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters for pagination and search
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
			return
		}

		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil || limit < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
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
		total, err := userCollection.CountDocuments(context.Background(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count users"})
			return
		}

		// Calculate skip count based on pagination parameters
		skip := (page - 1) * limit

		// Fetch users with pagination and search filters applied
		cur, err := userCollection.Find(context.Background(), filter, options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		defer cur.Close(context.Background())

		var users []models.User
		if err := cur.All(context.Background(), &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
			return
		}

		// Hide password field in each user
		for i := range users {
			users[i].Password = nil
		}
		// Construct response payload
		response := gin.H{
			"status": http.StatusOK,
			"data": gin.H{
				"users":   users,
				"total":   total,
				"page":    page,
				"perPage": limit,
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetUserDetails(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve userID from context
		// userID := c.Query("userId")
		userID, exists := c.Get("UserID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
			return
		}

		// Convert userID to string
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID format"})
			return
		}

		var userDetail models.User                                                                     // Assuming models.User is your user model
		err := userCollection.FindOne(context.TODO(), bson.M{"userid": userIDStr}).Decode(&userDetail) // Make sure the field name matches your MongoDB document
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching user details"})
			}
			return
		}

		update := bson.M{
			"$set": bson.M{
				"lastactive": time.Now(),
			},
		}

		// Perform the update operation
		filter := bson.M{"userid": userIDStr} // Correct the key to match your document
		result, err := userCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating user lastActive"})
			return
		}

		// Assuming you want to exclude certain sensitive details from the response
		// If Password is a string, you should set it to an empty string instead of nil
		userDetail.Password = nil // Adjust according to your model's field type

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"status":     true,
				"statusCode": http.StatusOK,
				"message":    "User fetch successfully",
				"result":     userDetail,
				"updated":    result,
			},
		})
	}
}

func UpdateUserDetails(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("userId")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in query parameters"})
			return
		}

		var updates UserUpdates
		if err := c.BindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		// Convert UserUpdates struct to a bson.M for MongoDB update, omitting zero values
		updateDoc := bson.M{}
		if !updates.LastActive.IsZero() {
			updateDoc["lastactive"] = updates.LastActive
		}

		if len(updateDoc) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no valid fields to update"})
			return
		}

		filter := bson.M{"userid": userID}
		update := bson.M{"$set": updateDoc}
		result, err := userCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating user details"})
			return
		}
		fmt.Println("result,", result)

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "no changes made to the user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":     true,
			"statusCode": http.StatusOK,
			"message":    "User updated successfully",
		})
	}
}
func ResetPassword(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request parameters
		var payload PasswordReset
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the user exists and credentials are valid
		var user models.User
		err := userCollection.FindOne(context.TODO(), bson.M{"email": payload.Email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found!"})
			return
		}

		// Verify current password
		if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(payload.CurrentPassword)); err != nil {
			logger.Log.Error("Invalid email or password", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Check if the new password is different from the previous password
		if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(payload.NewPassword)); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "New password must be different from the previous password"})
			return
		}

		// Update the password in the database
		hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			logger.Log.Error("Failed to generate password hash", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate password hash"})
			return
		}

		_, err = userCollection.UpdateOne(
			context.TODO(),
			bson.M{"email": payload.Email},
			bson.M{"$set": bson.M{"password": hashedNewPassword}},
		)
		if err != nil {
			logger.Log.Error("Failed to update password", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":     true,
			"statusCode": http.StatusOK,
			"message":    "Password reset successfully",
		})
	}
}

func ForgetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user input (e.g., email or username) from the request
		var req struct {
			Email string `json:"email" binding:"required,email"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request data",
			})
			return
		}

		// Check if the email exists in the database
		var user models.User
		err := userCollection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// User with the provided email does not exist
				c.JSON(http.StatusNotFound, gin.H{
					"error": "User not found",
				})
				return
			}
			// Error occurred while fetching user from the database
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}

		// Generate a secure token for password reset
		resetToken, err := generateResetToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate reset token",
			})
			return
		}

		// Send password reset email to the user with the reset token
		err = sendPasswordResetEmail(*user.Email, resetToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to send password reset email",
			})
			return
		}

		// Password reset email sent successfully
		c.JSON(http.StatusOK, gin.H{
			"message": "Password reset instructions sent to your email",
		})
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

func DeleteAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID string `json:"userId" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request data",
			})
			return
		}

		// Convert userID string to ObjectID
		userObjectID, err := primitive.ObjectIDFromHex(req.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid userID format",
			})
			return
		}

		// Check if the user exists
		var user models.User
		err = userCollection.FindOne(context.Background(), bson.M{"_id": userObjectID}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// User with the provided userID does not exist
				c.JSON(http.StatusNotFound, gin.H{
					"error": "User not found",
				})
				return
			}
			// Error occurred while fetching user from the database
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}

		// Delete the user account
		_, err = userCollection.DeleteOne(context.Background(), bson.M{"_id": userObjectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete user account",
			})
			return
		}

		// User account deleted successfully
		c.JSON(http.StatusOK, gin.H{
			"message": "User account deleted successfully",
		})
	}
}
