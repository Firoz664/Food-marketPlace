package utils

import (
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var regionName = "ap-south-1"
var s3BucketName = "golang-ucket"

func UploadFileToS3(file *multipart.FileHeader, userID, key string) (string, error) {
	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(regionName), // Set your AWS region
	})
	if err != nil {
		return "", err
	}

	// Create S3 client
	svc := s3.New(sess)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Specify S3 bucket name
	bucket := s3BucketName // Set your S3 bucket name

	// Generate key using the provided structure (upload/profile/userid/key)
	fileKey := fmt.Sprintf("upload/profile/%s/%s", userID, key)

	// Upload file to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
		Body:   src,
	})
	if err != nil {
		return "", err
	}

	// Generate file URL
	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, fileKey)

	return fileURL, nil
}

func DeleteFileFromS3(userID, fileName string) error {
	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(regionName), // Set your AWS region
	})
	if err != nil {
		return err
	}

	// Create S3 client
	svc := s3.New(sess)

	// Specify S3 bucket name
	bucket := s3BucketName // Set your S3 bucket name

	// Generate key for the file to be deleted
	fileKey := fmt.Sprintf("upload/profile/%s/%s", userID, fileName)

	// Delete the file from S3
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return err
	}

	return nil
}
