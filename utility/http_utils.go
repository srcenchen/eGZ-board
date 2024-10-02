package utility

import (
	"io"
	"net/http"
)

func HttpGet(url string) string {
	res, _ := http.Get(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)
	body, _ := io.ReadAll(res.Body)
	return string(body)
}
