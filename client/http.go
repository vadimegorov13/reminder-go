// HTTP client that communicates with Backend API
package client

import "net/http"

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
