package restapi

import "fmt"

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

func (e ResponseError) Error() string {
	return fmt.Sprintf("%s", e.Data.Attributes.Detail)
}
