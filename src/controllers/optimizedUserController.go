package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	tokenGenerate "github.com/foodmngtapp/food-management-apps/src/common/helpers/tokenHandler"
	"github.com/foodmngtapp/food-management-apps/src/common/logger"
	"github.com/foodmngtapp/food-management-apps/src/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// OptimizedLogin provides a performance-optimized login endpoint
func OptimizedLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		
		// Reduced timeout for faster response
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			logger.Log.Error("JSON binding error", zap.Error(err), zap.Duration("duration", time.Since(startTime)))
			c.JSON(http.StatusBadRequest, gin.H{
				"status":     false,
				"statusCode": http.StatusBadRequest,
				"error":      "Invalid request format",
			})
			return
		}

		if user.Email == nil || user.Password == nil {
			logger.Log.Warn("Missing credentials", zap.Duration("duration", time.Since(startTime)))
			c.JSON(http.StatusBadRequest, gin.H{
				"status":     false,
				"statusCode": http.StatusBadRequest,
				"error":      "Email and password are required",
			})
			return
		}

		*user.Email = strings.ToLower(*user.Email)
		
		// Optimized database query with projection
		var foundUser models.User
		collection := getUserCollection()
		
		// Only fetch required fields to reduce network I/O
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
		logger.Log.Debug("Database query completed", 
			zap.Duration("query_duration", queryDuration),
			zap.String("email", *user.Email))

		if err != nil {
			if err == mongo.ErrNoDocuments {
				logger.Log.Warn("User not found", 
					zap.String("email", *user.Email),
					zap.Duration("total_duration", time.Since(startTime)))
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":     false,
					"statusCode": http.StatusUnauthorized,
					"error":      "Invalid credentials",
				})
			} else {
				logger.Log.Error("Database query error", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":     false,
					"statusCode": http.StatusInternalServerError,
					"error":      "Authentication service temporarily unavailable",
				})
			}
			return
		}

		// Password verification
		passwordStart := time.Now()
		if err := bcrypt.CompareHashAndPassword([]byte(*foundUser.Password), []byte(*user.Password)); err != nil {
			passwordDuration := time.Since(passwordStart)
			logger.Log.Warn("Invalid password attempt", 
				zap.String("email", *user.Email),
				zap.Duration("password_check_duration", passwordDuration),
				zap.Duration("total_duration", time.Since(startTime)))
			
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":     false,
				"statusCode": http.StatusUnauthorized,
				"error":      "Invalid credentials",
			})
			return
		}
		passwordDuration := time.Since(passwordStart)

		// Token generation
		tokenStart := time.Now()
		token, refreshToken, err := tokenGenerate.GenerateAllTokens(
			*foundUser.Email, 
			*foundUser.FirstName, 
			*foundUser.LastName, 
			foundUser.UserID)
		if err != nil {
			logger.Log.Error("Token generation failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":     false,
				"statusCode": http.StatusInternalServerError,
				"error":      "Authentication service error",
			})
			return
		}
		tokenDuration := time.Since(tokenStart)

		// Optimized token update - single atomic operation
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

		// Prepare response (exclude sensitive data)
		userResponse := foundUser
		userResponse.Password = nil
		
		// Include fresh tokens in response
		userResponse.Token = &token
		userResponse.RefreshToken = &refreshToken

		totalDuration := time.Since(startTime)
		
		// Performance logging
		logger.Log.Info("Login completed successfully",
			zap.String("email", *foundUser.Email),
			zap.Duration("total_duration", totalDuration),
			zap.Duration("query_duration", queryDuration),
			zap.Duration("password_duration", passwordDuration),
			zap.Duration("token_duration", tokenDuration),
			zap.Duration("update_duration", updateDuration))

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"status":     true,
				"statusCode": http.StatusOK,
				"message":    "Login successful",
				"result":     userResponse,
				"tokens": gin.H{
					"accessToken":  token,
					"refreshToken": refreshToken,
					"expiresIn":    "24h",
					"tokenType":    "Bearer",
				},
				"performance": gin.H{
					"duration_ms": totalDuration.Milliseconds(),
				},
			},
		})
	}
}

// OptimizedHashPassword uses bcrypt cost of 10 for better performance
func OptimizedHashPassword(password string) (string, error) {
	// Cost 10 provides good security while being much faster than cost 14
	// Cost 10: ~50ms, Cost 14: ~800ms
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FastSignup optimized version of user signup
func FastSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		
		// Reduced timeout
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			logger.Log.Error("JSON binding error", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"status":     false,
				"statusCode": http.StatusBadRequest,
				"error":      "Invalid request format",
			})
			return
		}

		// Validate user input
		if err := validate.Struct(user); err != nil {
			logger.Log.Error("Validation error", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"status":     false,
				"statusCode": http.StatusBadRequest,
				"error":      "Validation failed: " + err.Error(),
			})
			return
		}

		*user.Email = strings.ToLower(*user.Email)

		// Check for existing user with optimized query
		collection := getUserCollection()
		existingUser := collection.FindOne(ctx, bson.M{
			"$or": []bson.M{
				{"email": user.Email},
				{"mobileNumber": user.MobileNumber},
			},
		}, options.FindOne().SetProjection(bson.M{"_id": 1}))

		if existingUser.Err() == nil {
			logger.Log.Warn("User already exists", zap.String("email", *user.Email))
			c.JSON(http.StatusConflict, gin.H{
				"status":     false,
				"statusCode": http.StatusConflict,
				"error":      "User already exists with this email or mobile number",
			})
			return
		} else if existingUser.Err() != mongo.ErrNoDocuments {
			logger.Log.Error("Database error during user check", zap.Error(existingUser.Err()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":     false,
				"statusCode": http.StatusInternalServerError,
				"error":      "Registration service temporarily unavailable",
			})
			return
		}

		// Hash password with optimized cost
		hashedPassword, err := OptimizedHashPassword(*user.Password)
		if err != nil {
			logger.Log.Error("Password hashing failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":     false,
				"statusCode": http.StatusInternalServerError,
				"error":      "Registration failed",
			})
			return
		}

		// Prepare user document
		user.Password = &hashedPassword
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.LastActive = time.Now()
		user.ID = user.ID // Keep existing ID if provided, otherwise MongoDB will generate
		if user.UserID == "" {
			user.UserID = user.ID.Hex()
		}

		// Generate tokens
		token, refreshToken, err := tokenGenerate.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, user.UserID)
		if err != nil {
			logger.Log.Error("Token generation failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":     false,
				"statusCode": http.StatusInternalServerError,
				"error":      "Registration failed",
			})
			return
		}

		user.Token = &token
		user.RefreshToken = &refreshToken

		// Insert user
		_, insertErr := collection.InsertOne(ctx, user)
		if insertErr != nil {
			logger.Log.Error("User insertion failed", zap.Error(insertErr))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":     false,
				"statusCode": http.StatusInternalServerError,
				"error":      "Registration failed, please try again",
			})
			return
		}

		// Prepare response
		user.Password = nil
		totalDuration := time.Since(startTime)

		logger.Log.Info("User registration completed",
			zap.String("email", *user.Email),
			zap.Duration("duration", totalDuration))

		c.JSON(http.StatusCreated, gin.H{
			"data": gin.H{
				"status":  http.StatusCreated,
				"message": "User registration successful",
				"result":  user,
				"performance": gin.H{
					"duration_ms": totalDuration.Milliseconds(),
				},
			},
		})
	}
}