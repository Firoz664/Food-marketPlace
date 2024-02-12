package controllers

import (
	"context"

	"fmt"
	"math/rand"
	"net/http"
	"time"

	otpUtils "github.com/foodmngtapp/food-management-apps/src/common/helpers/utils"
	"github.com/foodmngtapp/food-management-apps/src/config/database"

	"github.com/foodmngtapp/food-management-apps/src/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var otpCollection *mongo.Collection = database.OpenCollection(database.Client, "otp")
var validateOtp = validator.New()

// GetInvoice returns a Gin handler function

func CreateOTP(otpCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var mongoCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.OTP

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"error":  err.Error(),
			})
			return
		}

		// Generating a 6-digit OTP
		rand.Seed(time.Now().UnixNano())
		otpCode := rand.Intn(900000) + 100000 // Ensures OTP is always 6 digits

		// Set expiration time to 5 minutes from now
		expirationTime := time.Now().Add(5 * time.Minute)

		// Update user struct with OTP details
		user.OtpCode = int64(otpCode)
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.ExpiredAt = expirationTime

		// Save OTP to the database
		result, err := otpCollection.InsertOne(mongoCtx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": false,
				"error":  "Failed to save OTP",
			})
			return
		}

		// Send OTP based on whether email or mobile number is provided
		if user.Email != "" {
			// Assuming you have sendOtpOnEmail function
			if err := otpUtils.SendOtpViaEmail(user.Email, otpCode); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": false,
					"error":  "Failed to send OTP via email",
				})
				return
			}
		} else if user.MobileNumber != "" {
			// Assuming you have sendOtpviaMobile function
			if err := otpUtils.SendOtpViaMobile(user.MobileNumber, otpCode); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": false,
					"error":  "Failed to send OTP via mobile",
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"error":  "No email or mobile number provided",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "OTP generated and sent successfully",
			"otpID":   result.InsertedID,
		})
	}
}

type VerifyInputOtp struct {
	Email        string `json:"email"`
	MobileNumber string `json:"mobileNumber"`
	OtpCode      int64  `json:"otpCode" bson:"otpCode" validate:"required"`
}

func VerifyOtp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input VerifyInputOtp
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "Invalid request data"})
			return
		}

		// Check if either email or mobile number is provided
		if input.Email == "" && input.MobileNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "Either email or mobile number is required"})
			return
		}

		// Build search criteria based on provided email or mobile number
		searchCriteria := bson.M{}
		if input.Email != "" {
			searchCriteria["email"] = input.Email
		}
		if input.MobileNumber != "" {
			searchCriteria["mobileNumber"] = input.MobileNumber
		}

		var foundOTP models.OTP
		err := otpCollection.FindOne(context.TODO(), searchCriteria).Decode(&foundOTP)
		if err != nil {
			errorMsg := "Invalid User!"
			if err == mongo.ErrNoDocuments {
				// If email or mobile number is provided, dynamically include it in the error message
				if input.Email != "" {
					errorMsg = fmt.Sprintf("%s %s not found", errorMsg, input.Email)
				} else if input.MobileNumber != "" {
					errorMsg = fmt.Sprintf("%s %s not found", errorMsg, input.MobileNumber)
				}
				c.JSON(http.StatusNotFound, gin.H{"status": false, "msg": errorMsg})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "Failed to query database"})
			}
			return
		}

		// Check if OTP has expired
		if time.Now().After(foundOTP.ExpiredAt) {
			c.JSON(http.StatusOK, gin.H{"status": false, "msg": "OTP has expired"})
			return
		}

		// Verify OTP code
		if foundOTP.OtpCode == input.OtpCode {
			// Update IsVerified to true
			update := bson.M{"$set": bson.M{"isVerify": true}}
			_, err := otpCollection.UpdateOne(context.TODO(), bson.M{"_id": foundOTP.ID}, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "Failed to update verification status"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": true, "msg": "OTP verified successfully"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": false, "msg": "Verification failed due to incorrect OTP code"})
		}
	}
}
