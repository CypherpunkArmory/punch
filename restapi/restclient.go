package restapi

import (
	"net/http"
)

//need to manually update this
var apiVersion = "2019.4.19.1"

//RestClient A stuct to hold persistent data that is used between rest calls
type RestClient struct {
	URL          string
	RefreshToken string
	Client       http.Client
}

//NewRestClient Use this method to create a new rest client so headers can be setup
func NewRestClient(apiEndpoint string, refreshToken string) RestClient {
	client := http.DefaultClient
	rt := WithHeader(client.Transport)
	rt.Set("Api-Version", apiVersion)
	client.Transport = rt
	restAPI := RestClient{
		URL:          apiEndpoint,
		RefreshToken: refreshToken,
		Client:       *client,
	}
	return restAPI
}

//SetAPIKey set api key header
func (restClient *RestClient) SetAPIKey(apiKey string) {
	rt := WithHeader(restClient.Client.Transport)
	rt.Set("Authorization", "Bearer "+apiKey)
	restClient.Client.Transport = rt
}

//RoundTrip middleware that sets the headers of each request
func (h ClientHandler) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.Header {
		req.Header[k] = v
	}

	return h.rt.RoundTrip(req)
}

//ClientHandler struct to store header and roundtrip info
type ClientHandler struct {
	http.Header
	rt http.RoundTripper
}

//WithHeader Adds ability to edit headers without remaking everything
func WithHeader(rt http.RoundTripper) ClientHandler {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return ClientHandler{Header: make(http.Header), rt: rt}
}
