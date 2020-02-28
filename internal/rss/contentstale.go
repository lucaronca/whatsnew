package rss

import (
	"log"
	"net/http"
	"os"
	"time"
)

const (
	// cronTimeSpan is the cron interval in hours
	cronTimeSpan = 4
	extraSpan    = 2
)

// IsContentStale verifies if content has changed
// counting from the last cron
func IsContentStale(resp *http.Response) bool {
	lastModified := resp.Header["Last-Modified"][0]

	lastModifiedTime, err := http.ParseTime(lastModified)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	now := time.Now().UTC()
	lastCheckedTime := now.Add((-cronTimeSpan * time.Hour) + (-extraSpan * time.Hour))

	isFresh := lastModifiedTime.After(lastCheckedTime)

	return !isFresh
}
