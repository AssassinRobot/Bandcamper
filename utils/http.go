package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func (h *HttpMngmnt) Get(url string, headers map[string]string) (*http.Response, error) {
	log.Println("Getting... ", url)

	req, httpRequestError := http.NewRequestWithContext(context.TODO(), "GET", url, nil)
	if httpRequestError != nil {
		return nil, httpRequestError
	}

	for key, value := range headers {
		req.Header.Add(key, value)
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

func (h *HttpMngmnt) Post(url string, headers map[string]string, body map[string]any) (*http.Response, error) {
	log.Println("Posting... ", url)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader = bytes.NewReader(jsonBody)

	req, httpRequestError := http.NewRequestWithContext(context.TODO(), "POST", url, bodyReader)
	if httpRequestError != nil {
		return nil, httpRequestError
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", fmt.Sprintf("%d", len(jsonBody)))
	req.Header.Add("Accept", "application/json")

	res, httpPostError := h.client.Do(req)
	if httpPostError != nil {
		return nil, httpPostError
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return res, nil
}
