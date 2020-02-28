package filehandler

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	fileDir      = "/tmp"
	localFileDir = "./tmp"
)

// HandleFileRequest handle the rss file creating a temporary physical one
// and uploads it to S3
func HandleFileRequest(body string) error {
	contents, err := getFileContents(body)
	if err != nil {
		return err
	}

	filePath, err := createFile(contents)
	if err != nil {
		return err
	}

	u := s3Uploader{}

	err = u.createSession()
	if err != nil {
		return err
	}

	err = u.addFileToS3(filePath)
	if err != nil {
		return err
	}

	return nil
}

func getFileContents(body string) ([]byte, error) {
	type data struct {
		File string `json:"file"`
	}

	d := data{}

	json.Unmarshal([]byte(body), &d)

	return base64.StdEncoding.DecodeString(d.File)
}

func createFile(contents []byte) (string, error) {
	var filePath string

	if os.Getenv("LOCAL") == "true" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		filePath = filepath.Join(dir, localFileDir, os.Getenv("RSS_FILENAME"))
	} else {
		filePath = filepath.Join(fileDir, os.Getenv("RSS_FILENAME"))
	}

	err := makeEmptyFile(filePath)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(filePath, contents, 0644)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func makeEmptyFile(filePath string) error {
	if os.Getenv("LOCAL") == "true" {
		os.Mkdir(localFileDir, os.FileMode(0522))
	} else {
		os.Mkdir(fileDir, os.FileMode(0522))
	}

	file, err := os.Create(filePath)
	if err != nil {
		file.Close()
		return err
	}

	return nil
}
