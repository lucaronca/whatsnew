package url

import "os"

// MakeS3ObjectURL returns the S3 object url
func MakeS3ObjectURL() string {
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("s3RssBucketName")
	fileName := os.Getenv("s3RssFileName")

	return "https://" + bucket + ".s3." + region + ".amazonaws.com/" + fileName
}