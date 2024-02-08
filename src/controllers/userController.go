package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	tokenGenerate "github.com/foodmngtapp/food-management-apps/src/common/helpers/tokenHandler"
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/foodmngtapp/food-management-apps/src/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

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
		// Simulate fetching a user

		// Log fetching user operation
		// logger.Log.Info("Fetched user successfully", zap.String("userID", dummyUser.ID))
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "User fetched successfully",
		})
	}
}

func GetUserDetails(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve userID from query parameters
		userID := c.Query("userId") // Ensure this query parameter name matches what's sent by the client
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in query parameters"})
			return
		}

		var userDetail models.User                                                                  // Assuming models.User is your user model
		err := userCollection.FindOne(context.TODO(), bson.M{"userid": userID}).Decode(&userDetail) // Make sure the field name matches your MongoDB document
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
		filter := bson.M{"userId": userID}
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

func UpdateUserDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Dummy user before update

		// Attempt to bind the incoming JSON payload to the updatedUser struct

		// Log the successful update operation

		// Simulate updating the user's name and return the updated user object

		// Return the updated user object along with status code and message
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "User details updated successfully",
		})
	}
}
func ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate fetching a user

		// Log fetching user operation
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "User fetched successfully",
		})
	}
}
func ForgetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate fetching a user

		// Log fetching user operation
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "User fetched successfully",
		})
	}
}
func DeleteAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate fetching a user

		// Log fetching user operation
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "User fetched successfully",
		})
	}
}
