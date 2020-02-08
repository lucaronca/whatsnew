package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/lucaronca/whatsnew/s3urlconstructor"
)

const (
	s3Bucket = "whatsnew-bucket"
)

// S3Uploader make the upload call given
// a session
type s3Uploader struct {
	session *session.Session
}

// addFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func (u *s3Uploader) addFileToS3(uploadFileDir string) error {
	upFile, err := os.Open(uploadFileDir)

	if err != nil {
		return err
	}
	defer upFile.Close()

	// Get file size and read the file content into a buffer
	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	result, err := s3.New(u.session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(s3Bucket),
		Key:                  aws.String(filepath.Base(uploadFileDir)),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String("application/rss+xml"),
		ContentDisposition:   aws.String("inline"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("rss file updated on S3, file url: %v, ETag: %v\n", s3urlconstructor.MakeURL(), *result.ETag)
	return nil
}
