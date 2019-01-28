package restapi

import "net/http"

type RestClient struct {
	URL          string
	APIKEY       string
	RefreshToken string
	Client       http.Client
}
