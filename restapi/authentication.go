package restapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SessionResponse struct {
	Access_Token string `json:"access_token"`
	Token_Type   string `json:"token_type"`
	Expires_In   int    `json:"expires-in"`
}

func (restClient *RestClient) StartSession(refresh_token string) (SessionResponse, error) {
	responseBody := SessionResponse{}
	url := restClient.URL + "/session"
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, nil)
	req.Header.Add("Authorization", "Bearer "+refresh_token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode > 399 {
		//errorBody := ErrorResponse{}
		errorBody := ResponseError{}
		err = json.Unmarshal(body, &errorBody)
		return responseBody, errorBody
	}

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		fmt.Println("error:", err)
		return responseBody, http.ErrAbortHandler
	}
	return responseBody, nil
}
