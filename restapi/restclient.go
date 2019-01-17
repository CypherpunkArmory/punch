package restapi

import "net/http"

type RestClient struct {
	URL    string
	APIKEY string
	Client http.Client
}
