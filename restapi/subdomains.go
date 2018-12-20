package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ReserveSubdomainRequest struct {
	Subdomain string `json:"subdomain"`
}
type ReserveSubdomainResponse struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Reserved bool   `json:"reserved"`
			In_use   bool   `json:"in_use"`
			Name     string `json:"name"`
		} `json:"attributes"`
		ID string `json:"id"`
	} `json:"data"`
}
type SubdomainListResponse struct {
	Data []struct {
		Type       string `json:"type"`
		Attributes struct {
			Reserved bool   `json:"reserved"`
			Name     string `json:"name"`
			InUse    bool   `json:"in_use"`
		} `json:"attributes"`
		ID string `json:"id"`
	} `json:"data"`
}

//SubdomainListAPI get list of subdomains reserved
func (restClient *RestClient) SubdomainListAPI() (SubdomainListResponse, error) {
	responseBody := SubdomainListResponse{}
	url := restClient.URL + "/subdomain"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+restClient.APIKEY)
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

//ReserveSubdomainAPI calls holepunch web api to get reserve a subdomain
func (restClient *RestClient) ReserveSubdomainAPI(requestBody ReserveSubdomainRequest) (ReserveSubdomainResponse, error) {
	responseBody := ReserveSubdomainResponse{}
	url := restClient.URL + "/subdomain"
	jsonStr, err := json.Marshal(&requestBody)
	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+restClient.APIKEY)
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

//ReleaseSubodmainAPI deletes tunnel
func (restClient *RestClient) ReleaseSubodmainAPI(subdomainName string) error {

	url := restClient.URL + "/subdomain/" + subdomainName
	client := &http.Client{}
	req, err := http.NewRequest("Delete", url, nil)
	req.Header.Add("Authorization", "Bearer "+restClient.APIKEY)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		return nil
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode > 399 {
		//errorBody := ErrorResponse{}
		errorBody := ResponseError{}
		err = json.Unmarshal(body, &errorBody)
		return errorBody
	}
	return nil
}
