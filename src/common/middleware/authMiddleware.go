package middleware

import (
	"fmt"
	"net/http"

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
		ClientToken := c.Request.Header.Get("x-access-token")
		if ClientToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Authorization is Required"})
			c.Abort()
			return
		}

		claims, err := token.ValidateToken(ClientToken)
		fmt.Println("Claims----->>>", claims)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
			c.Abort()
			return
		}

		// Set both "Email" and "UserID" in the context
		c.Set("Email", claims.Email)
		c.Set("UserID", claims.Uid)

		c.Next()
	}
}
