package letterboxdgo

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var REQ_DELAY = time.Second * 10 // Delay between requests to avoid rate limiting
var MAX_RETRIES = 7              // Maximum number of retries for a request

func Get(url string) (io.ReadCloser, error) {

	for i := 0; i < MAX_RETRIES; i++ {

		// Send request
		userAgent := "curl"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			slog.Error("Could not create GET request", "error", err)
			os.Exit(1)
		}
		req.Header.Set("User-Agent", userAgent)

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			slog.Error("Could not send GET request", "error", err)
			os.Exit(1)
		}

		if res.StatusCode == 200 {
			return res.Body, nil
		} else {
			res.Body.Close()
			if res.StatusCode == 429 {
				slog.Debug("Hit rate limit, waiting...", "code", res.StatusCode, "message", res.Status, "try", i+1)
				time.Sleep(REQ_DELAY)
			} else {
				slog.Debug("Non-200 response from GET request", "code", res.StatusCode, "message", res.Status, "try", i+1)
				return nil, errors.New("non-200 response from GET request")
			}
		}
	}

	return nil, errors.New("could not complete GET request")
}
