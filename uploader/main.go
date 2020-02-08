package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Request is of type APIGatewayProxyRequest since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Request events.APIGatewayProxyRequest

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is the lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
	type data struct {
		File string `json:"file"`
	}

	d := data{}

	json.Unmarshal([]byte(request.Body), &d)

	contents, err := base64.StdEncoding.DecodeString(d.File)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = Upload(contents)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "OK\n",
		Headers: map[string]string{
			"Content-Type":           "text/plain",
			"X-MyCompany-Func-Reply": "uploader-handler",
		},
	}

	return resp, nil
}

func main() {
	if os.Getenv("LOCAL") == "true" {
		local()
	} else {
		lambda.Start(Handler)
	}
}
