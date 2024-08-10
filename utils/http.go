package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type HttpMngmnt struct {
	client *http.Client
}

func NewHttpMngmnt() *HttpMngmnt {
	return &HttpMngmnt{
		client: &http.Client{},
	}
}

func (h *HttpMngmnt) Get(url string) (*http.Response, error) {
	log.Println("Getting... ", url)

	req, httpRequestError := http.NewRequestWithContext(context.TODO(), "GET", url, nil)
	if httpRequestError != nil {
		return nil, httpRequestError
	}

	res, httpGetError := h.client.Do(req)
	if httpGetError != nil {
		return nil, httpGetError
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return res, nil
}
