package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/foodmngtapp/food-management-apps/src/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

// Dummy user struct for demonstration purposes

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		fmt.Println("User", user)
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Assume validate is a valid validator instance of github.com/go-playground/validator/v10
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Printf("Error occurred while checking for the email: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the email"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this email already exists"})
			return
		}

		// Insert the user into the database, omitted for brevity

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "User signed up successfully",
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate fetching a user
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		fmt.Println("User", user)
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "false",
				"error":  err.Error(),
			})
			return
		}
		if user.Email == nil || *user.Email == "" || user.Password == nil || *user.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "false",
				"error":  "Email or Password missing",
			})
			return
		}

		// Assume validate is a valid validator instance of github.com/go-playground/validator/v10

		// Log fetching user operation
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"status":  http.StatusOK,
				"message": "User login successfully",
				"result": gin.H{
					"name":  "Shams Firoz",
					"email": "email@test.com",
				},
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

func GetUserDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate fetching a user

		// Log fetching user operation

		// Return the dummy user object along with status code and message
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "User fetched successfully",
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
