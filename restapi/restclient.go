package restapi

import "net/http"

type RestClient struct {
	URL    string
	APIKEY string
	client http.Client
}
