// Punch CLI used for interacting with holepunch.io
// Copyright (C) 2018-2019  Orb.House, LLC
// 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package restapi

import (
	"net/http"
)

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
