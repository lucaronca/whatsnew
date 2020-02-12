package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gorilla/feeds"
	"github.com/lucaronca/whatsnew/rss"
	"github.com/lucaronca/whatsnew/url"
)

const (
	rssTitle = "JW.org cosa c'è di nuovo"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

func makeResponse() Response {
	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "OK\n",
		Headers: map[string]string{
			"Content-Type":           "text/plain",
			"X-MyCompany-Func-Reply": "generator-handler",
		},
	}
}

// Handler is the lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	start := time.Now()
	fmt.Printf("start fetching %v\n", rss.TargetPageURL)

	resp := rss.MakePageRequest()

	secs := time.Since(start).Seconds()
	fmt.Printf("%s request fulfilled, %.2fs elapsed\n", rss.TargetPageURL, secs)

	if rss.IsContentStale(resp) {
		defer resp.Body.Close()
		return makeResponse(), nil
	}

	rg := rss.Generator{
		URL: rss.TargetPageURL,
		Feed: feeds.Feed{
			Title:       rssTitle,
			Link:        &feeds.Link{Href: url.MakeS3ObjectURL()},
			Description: rssTitle,
		},
		Data: rss.Data{},
	}

	liveDone := make(chan bool)
	storedDone := make(chan bool)

	go rg.GetStored(storedDone)
	go func(done chan bool) {

		rg.GetDocument(resp)
		defer resp.Body.Close()
		rg.ParseDocument()
		rg.GetRssData()

		done <- true
	}(liveDone)

	_, _ = <-storedDone, <-liveDone

	hasChanged := rg.Compare()

	if hasChanged {
		err := rss.MakeUploadRequest(&rg)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}

	return makeResponse(), nil
}

func main() {
	if os.Getenv("LOCAL") == "true" {
		local()
	} else {
		lambda.Start(Handler)
	}
}
