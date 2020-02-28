package url

import "os"

// MakeS3ObjectURL returns the S3 object url
func MakeS3ObjectURL() string {
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("S3_BUCKET")
	fileName := os.Getenv("RSS_FILENAME")

	return "https://" + bucket + ".s3." + region + ".amazonaws.com/" + fileName
}
