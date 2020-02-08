package rss

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"bou.ke/monkey"
)

// TestIsContentStale verifies that function returns true
// if content was last modified before the last cron check
// true otherwise
func TestIsContentStale(t *testing.T) {
	tables := []struct {
		now          string
		lastModified string
		isStale      bool
	}{
		// {"Mon, 20 Jan 2020 12:00:00 UTC", "Wed, 15 Jan 2020 21:22:59 GMT", true},
		// {"Mon, 20 Jan 2020 12:00:00 UTC", "Wed, 31 Dec 2019 13:20:15 GMT", true},
		// {"Mon, 20 Jan 2020 12:00:00 UTC", "Mon, 20 Jan 2020 09:00:00 GMT", false},
		// {"Mon, 20 Jan 2020 12:00:00 UTC", "Mon, 20 Jan 2020 00:00:00 GMT", false},
		// {"Mon, 20 Jan 2020 12:00:00 UTC", "Sun, 19 Jan 2020 23:00:00 GMT", false},
		// {"Fri, 24 Jan 2020 12:00:00 UTC", "Fri, 24 Jan 2020 11:20:15 GMT", false},
		// {"Fri, 24 Jan 2020 12:00:00 UTC", "Thu, 23 Jan 2020 21:59:59 GMT", true},
		// {"Fri, 24 Jan 2020 12:00:00 UTC", "Thu, 23 Jan 2020 22:00:01 GMT", false},
	}

	for _, table := range tables {
		monkey.Patch(time.Now, func() time.Time {
			time, _ := time.Parse(time.RFC1123, table.now)
			return time
		})

		h := make(http.Header)
		h["Last-Modified"] = []string{table.lastModified}

		r := http.Response{
			Body:   ioutil.NopCloser(bytes.NewBufferString("mock")),
			Header: h,
		}

		isStale := IsContentStale(&r)

		defer r.Body.Close()

		monkey.Unpatch(time.Now)

		if isStale != table.isStale {
			t.Errorf("isStale was incorrect for Last-Modified: %v at %v, got: %t, want: %t.", table.lastModified, table.now, isStale, table.isStale)
		}
	}
}
