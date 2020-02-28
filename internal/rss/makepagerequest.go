package rss

import (
	"log"
	"net/http"
	"os"
)

// MakePageRequest does the request to main page
func MakePageRequest() *http.Response {
	client := &http.Client{}

	req, err := http.NewRequest("GET", os.Getenv("TARGET_PAGE_URL"), nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return resp
}
