package easyhttp

import (
	"net/http"
	"time"
)

var (
	httpClient *http.Client
)

func init() {
	httpClient = &http.Client{
		Timeout: time.Second * 5,
	}
}

func Do(request *http.Request) (*http.Response, error) {
	return httpClient.Do(request)
}
