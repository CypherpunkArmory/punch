package restapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SubdomainAttributes struct {
	InUse    bool   `json:"in_use"`
	Name     string `json:"name"`
	Reserved bool   `json:"reserved"`
}
type SubdomainJsonData struct {
	Type       string              `json:"type"`
	Attributes SubdomainAttributes `json:"attributes"`
	ID         string              `json:"id"`
}
type ReserveSubdomainRequest struct {
	Data SubdomainJsonData `json:"data"`
}
type ReserveSubdomainResponse struct {
	Data SubdomainJsonData `json:"data"`
}
type SubdomainListResponse struct {
	Data []SubdomainJsonData `json:"data"`
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
	id, err := restClient.GetSubdomainID(subdomainName)
	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
	if id == "" {
		return errors.New("You do not own this subdomain")
	}
	url := restClient.URL + "/subdomain/" + id
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
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
	return errors.New("Failed to delete")
}

func (restClient *RestClient) GetSubdomainID(subdomainName string) (string, error) {
	responseBody, err := restClient.SubdomainListAPI()
	if err != nil {
		panic(err)
	}
	SubdomainList := responseBody.Data
	for _, domain := range SubdomainList {
		if domain.Attributes.Name == subdomainName {

			return domain.ID, nil
		}
	}
	return "", nil
}
