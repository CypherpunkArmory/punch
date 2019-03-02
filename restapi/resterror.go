package restapi

import (
	"errors"
	"fmt"
)

var errorUnownedSubdomain = errors.New("you do not own this subdomain")
var errorCantConnectRestCall = errors.New("problem contacting the server")
var errorUnableToParse = errors.New("can't parse the json response")
var errorUnownedTunnel = errors.New("you do not own this subdomain")
var errorUnableToDelete = errors.New("failed to delete")

//ResponseError JSONapi response error
type ResponseError struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Title  string `json:"title"`
			Status string `json:"status"`
			Detail string `json:"detail"`
		} `json:"attributes"`
		ID string `json:"id"`
	} `json:"data"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf(e.Data.Attributes.Detail)
}
