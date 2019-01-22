package restapi

import "net/http"

type RestClient struct {
	URL           string
	APIKEY        string
	ResfreshToken string
	Client        http.Client
}
