package restapi

import (
	"net/http"
)

//RestClient A stuct to hold persistent data that is used between rest calls
type RestClient struct {
	URL          string
	APIKEY       string
	RefreshToken string
	Client       http.Client
}
