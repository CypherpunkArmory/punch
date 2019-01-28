package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type SessionResponse struct {
	Access_Token  string `json:"access_token"`
	Token_Type    string `json:"token_type"`
	Expires_In    int    `json:"expires-in"`
	Refresh_Token string `json:"refresh_token"`
}

type LoginRequest struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func (restClient *RestClient) StartSession(refresh_token string) (SessionResponse, error) {
	responseBody := SessionResponse{}
	url := restClient.URL + "/session"
	req, err := http.NewRequest("PUT", url, nil)
	req.Header.Add("Authorization", "Bearer "+refresh_token)

	resp, err := restClient.Client.Do(req)
	if err != nil {
		fmt.Println("Could not connect to the api server")
		os.Exit(1)
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
	restClient.APIKEY = responseBody.Access_Token
	return responseBody, nil
}

func (restClient *RestClient) Login(username string, password string) (SessionResponse, error) {
	responseBody := SessionResponse{}
	url := restClient.URL + "/login"

	reqBody := LoginRequest{
		Username: username,
		Password: password,
	}
	jsonStr, _ := json.Marshal(&reqBody)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")

	resp, err := restClient.Client.Do(req)

	if err != nil {
		return responseBody, err
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
		return responseBody, http.ErrAbortHandler
	}
	restClient.RefreshToken = responseBody.Refresh_Token
	restClient.APIKEY = responseBody.Access_Token
	return responseBody, nil
}
