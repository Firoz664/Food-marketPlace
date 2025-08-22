package tokenhandler

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/foodmngtapp/food-management-apps/src/config"
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedTokendDetails struct {
	Email     string
	FirstName string
	LastName  string
	UserID    string
	Uid       string
	ID        string
	jwt.RegisteredClaims
}

var userCollection *mongo.Collection

// getUserCollection returns the user collection with lazy initialization
func getUserCollection() *mongo.Collection {
	if userCollection == nil {
		client := database.DBInstance()
		userCollection = database.OpenCollection(client, "users")
	}
	return userCollection
}

// GetJWTSecret returns the JWT secret key from config or environment
func GetJWTSecret() string {
	// Try to get from config first
	if config.AppConfig != nil && config.AppConfig.JWT.SecretKey != "" {
		return config.AppConfig.JWT.SecretKey
	}
	
	// Fallback to environment variable
	secret := os.Getenv("FOOD_APP_JWT_SECRET_KEY")
	if secret == "" {
		secret = os.Getenv("SIGNING_SECRET_KEY") // Legacy support
	}
	
	if secret == "" {
		log.Fatal("JWT secret key is required. Set FOOD_APP_JWT_SECRET_KEY environment variable or configure it in config file")
	}
	
	return secret
}

func GenerateAllTokens(email string, firstName string, lastName string, userId string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedTokendDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UserID:    userId,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(24))),
		},
	}
	fmt.Println("Claims", claims)

	refreshClaims := &SignedTokendDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(168))),
		},
	}

	secretKey := GetJWTSecret()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secretKey))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedTokendDetails, msg string) {
	secretKey := GetJWTSecret()
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedTokendDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedTokendDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt.Before(time.Now().Local()) {
		msg = fmt.Sprintf("token is expired")
		return
	}
	return claims, msg
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refreshToken", Value: signedRefreshToken})

	UpdatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: UpdatedAt})

	upsert := true
	filter := bson.M{"userid": userId}
	fmt.Println("filter", filter) //print filter apply
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	collection := getUserCollection()
	_, err := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)

	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
	return
}
