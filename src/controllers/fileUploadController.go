package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	s3uploader "github.com/foodmngtapp/food-management-apps/src/common/helpers/utils"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func S3FileUploadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve user ID from query parameter
		userID := c.Query("userId")
		fmt.Println("userId------>>>,", userID)
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not provided"})
			return
		}

		// Retrieve file from form-data
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file from form-data"})
			return
		}

		// Upload file to Amazon S3
		fileURL, err := s3uploader.UploadFileToS3(file, userID, file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to S3"})
			return
		}

		// Save file URL to user collection in MongoDB
		err = saveFileURLToUserCollection(userID, fileURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file URL to user collection"})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "File uploaded successfully",
			"url":     fileURL,
		})
	}
}

// saveFileURLToUserCollection saves the file URL to user collection in MongoDB
func saveFileURLToUserCollection(userID, fileURL string) error {
	// Convert userID string to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	// Construct filter to find user by ID
	filter := bson.M{"_id": userObjectID}
	// Construct update to set file URL and update updatedAt field
	update := bson.M{
		"$set": bson.M{
			"profileimage": fileURL,
			"updatedAt":    time.Now(),
		},
	}
	// Perform update operation in MongoDB
	result, err := userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		// Log error details
		fmt.Println("Error updating user document:", err)
		return err
	}

	// Check if any document was modified
	if result.ModifiedCount == 0 {
		return fmt.Errorf("no document was updated")
	}

	return nil
}

func DeleteFileFromS3Bucket() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve user ID from query parameter
		userID := c.Query("userId")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not provided"})
			return
		}

		// Retrieve file name from query parameter
		fileName := c.Query("fileName")
		if fileName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File name not provided"})
			return
		}

		// Delete file from S3 bucket
		err := s3uploader.DeleteFileFromS3(userID, fileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from S3 bucket", "details": err.Error()})
			return
		}
		err = deleteFileUrlDb(userID, fileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file URL to user collection"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "File deleted successfully from S3 bucket",
		})
	}
}
func deleteFileUrlDb(userID, fileURL string) error {
	// Convert userID string to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	// Construct filter to find user by ID
	filter := bson.M{"_id": userObjectID}
	// Construct update to set file URL and update updatedAt field
	update := bson.M{
		"$set": bson.M{
			"profileimage": nil,
			"updatedAt":    time.Now(),
		},
	}
	// Perform update operation in MongoDB
	result, err := userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		// Log error details
		fmt.Println("Error updating user document:", err)
		return err
	}

	// Check if any document was modified
	if result.ModifiedCount == 0 {
		return fmt.Errorf("no document was updated")
	}

	return nil
}

func MultipleUploadFiles() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate fetching a user

		// Log fetching user operation
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "User fetched successfully",
		})
	}
}
