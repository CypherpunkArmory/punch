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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//SessionResponse Json response for login/session refresh
type SessionResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires-in"`
	RefreshToken string `json:"refresh_token"`
}

type loginRequest struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

//StartSession Start a session and set the restClient to the current access token
func (restClient *RestClient) StartSession(refreshToken string) error {
	responseBody := SessionResponse{}
	url := restClient.URL + "/session"
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return errorCantConnectRestCall
	}
	req.Header.Add("Authorization", "Bearer "+refreshToken)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errorCantConnectRestCall
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode > 399 {
		errorBody := ResponseError{}
		err = json.Unmarshal(body, &errorBody)
		if err != nil {
			return err
		}
		return &errorBody
	}

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return errorUnableToParse
	}
	restClient.SetAPIKey(responseBody.AccessToken)
	return nil
}

//Login Login user with given username and password. Returns sessionresponse so cmd can set viper configs
func (restClient *RestClient) Login(username string, password string) (SessionResponse, error) {
	responseBody := SessionResponse{}
	url := restClient.URL + "/login"

	reqBody := loginRequest{
		Username: username,
		Password: password,
	}
	jsonStr, _ := json.Marshal(&reqBody)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")

	resp, err := restClient.Client.Do(req)

	if err != nil {
		return responseBody, errorCantConnectRestCall
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode > 399 {
		errorBody := ResponseError{}
		err = json.Unmarshal(body, &errorBody)
		if err != nil {
			return responseBody, err
		}
		return responseBody, &errorBody
	}

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return responseBody, errorUnableToParse
	}
	restClient.RefreshToken = responseBody.RefreshToken
	restClient.SetAPIKey(responseBody.AccessToken)
	return responseBody, nil
}
