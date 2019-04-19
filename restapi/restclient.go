package restapi

import (
	"net/http"
)

//RestClient A stuct to hold persistent data that is used between rest calls
type RestClient struct {
	URL          string
	RefreshToken string
	Client       http.Client
}

func NewRestClient(apiEndpoint string, refreshToken string) RestClient {
	client := http.DefaultClient
	rt := WithHeader(client.Transport)
	rt.Set("Api-Version", "2019.4.19.1")
	client.Transport = rt
	restAPI := RestClient{
		URL:          apiEndpoint,
		RefreshToken: refreshToken,
		Client:       *client,
	}
	return restAPI
}

func (restClient *RestClient) SetAPIKey(apiKey string) {
	rt := WithHeader(restClient.Client.Transport)
	rt.Set("Authorization", "Bearer "+apiKey)
	restClient.Client.Transport = rt
}

func (h withHeader) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.Header {
		req.Header[k] = v
	}

	return h.rt.RoundTrip(req)
}

type withHeader struct {
	http.Header
	rt http.RoundTripper
}

func WithHeader(rt http.RoundTripper) withHeader {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return withHeader{Header: make(http.Header), rt: rt}
}
