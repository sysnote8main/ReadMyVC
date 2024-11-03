package easyhttp

import (
	"io"
	"net/http"
)

func RequestGet(url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, url, body)
}

func RequestPost(url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(http.MethodPost, url, body)
}
