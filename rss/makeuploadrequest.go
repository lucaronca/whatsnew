package rss

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// MakeUploadRequest encodes the xml data and sends its contents
// to the uploader lambda
func MakeUploadRequest(rg *Generator) error {
	buf := new(bytes.Buffer)
	rg.WriteRss(buf)

	file := buf.Bytes()
	contents := base64.StdEncoding.EncodeToString(file)

	body, err := json.Marshal(map[string]interface{}{
		"file": contents,
	})
	if err != nil {
		return err
	}

	// Payload should be a json containing a string serialized "body" property,
	// "body" contents will be then deserialized and passed to the lambda
	// in the request by aws
	type Payload struct {
		Body string `json:"body"`
	}
	p := Payload{
		Body: string(body),
	}
	payload, _ := json.Marshal(p)

	client := createClient()

	_, err = client.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String(os.Getenv("UPLOADER_LAMBDA_NAME")),
		Payload:      payload,
	})
	if err != nil {
		return err
	}
	return nil
}

func createClient() *lambda.Lambda {
	sess := session.Must(session.NewSession())
	region := os.Getenv("AWS_REGION")

	client := lambda.New(sess, &aws.Config{
		Region: aws.String(region),
	})

	return client
}
