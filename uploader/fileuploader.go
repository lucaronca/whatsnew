package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	fileDir      = "/tmp"
	localFileDir = "./tmp"
)

// Upload creates the rss file and uploads
// it to S3
func Upload(contents []byte) error {
	var filePath string

	if os.Getenv("LOCAL") == "true" {
		filePath = filepath.Join(localFileDir, os.Getenv("s3RssFileName"))
	} else {
		filePath = filepath.Join(fileDir, os.Getenv("s3RssFileName"))
	}

	err := createFile(filePath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, contents, 0644)
	if err != nil {
		return err
	}

	u := s3Uploader{}

	err = createSession(&u)
	if err != nil {
		return err
	}

	err = u.addFileToS3(filePath)
	if err != nil {
		return err
	}

	return nil
}

func createFile(filePath string) error {
	os.Mkdir(fileDir, os.FileMode(0522))

	file, err := os.Create(filePath)
	if err != nil {
		file.Close()
		return err
	}

	return nil
}

func createSession(u *s3Uploader) error {
	// Create a single AWS session (we can re use this if we're uploading many files)
	s3Region := os.Getenv("AWS_REGION")

	s, err := session.NewSession(&aws.Config{Region: aws.String(s3Region)})
	if err != nil {
		return err
	}

	u.session = s

	return nil
}
