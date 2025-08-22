package middleware

import (
	"fmt"
	"net/http"

	"strings"

	token "github.com/foodmngtapp/food-management-apps/src/common/helpers/tokenHandler"

	"github.com/gin-gonic/gin"
)

// func Authentication() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ClientToken := c.Request.Header.Get("x-access-token")
// 		if ClientToken == "" {
// 			c.JSON(http.StatusBadRequest, gin.H{"Error": " Authorization is Required"})
// 			c.Abort()
// 		}
// 		claims, err := token.ValidateToken(ClientToken)
// 		fmt.Println("Claims", claims)
// 		if err != "" {
// 			c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
// 			c.Abort()
// 		}
// 		c.Set("Email", claims.Email)
// 		c.Set("Email", claims.Uid)
// 		c.Next()

// 	}
// }

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract token string from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}
		clientToken := tokenParts[1]
		// Validate the token using your token.ValidateToken method
		claims, err := token.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}
		fmt.Println(" claims Email userID from context:", claims)
		// Set values from token into context
		c.Set("email", claims.Email)
		c.Set("userID", claims.UserID) // Use claims.UserID instead of claims.Uid
		fmt.Println(" claims Email userID from context:", claims)
		c.Next()
	}
}

