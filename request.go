package letterboxdgo

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"
)

var REQ_DELAY = time.Second * 10 // Delay between requests to avoid rate limiting
var MAX_RETRIES = 7              // Maximum number of retries for a request

func Get(url string) (io.ReadCloser, error) {

	for i := 0; i < MAX_RETRIES; i++ {

		// Send request
		res, err := http.Get(url)
		if err != nil {
			slog.Error("Error sending GET request", "error", err)
			return nil, err
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
