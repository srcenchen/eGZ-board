package utility

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var HTTPClient = &http.Client{Timeout: 10 * time.Second}

func HttpGet(url string) (string, error) {
	res, err := HTTPClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("GET %s: %w", url, err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(io.LimitReader(res.Body, 10<<20))
	if err != nil {
		return "", fmt.Errorf("read GET %s response: %w", url, err)
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("GET %s returned %s: %s", url, res.Status, string(body))
	}
	if len(body) == 0 {
		return "", fmt.Errorf("GET %s returned an empty body", url)
	}
	return string(body), nil
}
