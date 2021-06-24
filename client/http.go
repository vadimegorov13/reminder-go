// HTTP client that communicates with Backend API
package client

import (
	"net/http"
	"time"
)

type HTTPClient struct {
	client     *http.Client
	BackendURI string
}

func NewHTTPClient(uri string) HTTPClient {
	return HTTPClient{
		client:     &http.Client{},
		BackendURI: uri,
	}
}

func (c HTTPClient) Create(title, message string, duration time.Duration) ([]byte, error) {
	res := []byte("Response for create reminder")
	return res, nil
}

func (c HTTPClient) Edit(id, title, message string, duration time.Duration) ([]byte, error) {
	res := []byte("Response for edit reminder")
	return res, nil
}

func (c HTTPClient) Fetch(id []string) ([]byte, error) {
	res := []byte("Response for fetch reminder")
	return res, nil
}

func (c HTTPClient) Delete(id []string) error {
	return nil
}

func (c HTTPClient) Health(host string) bool {
	return true
}
