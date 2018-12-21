package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TunnelAttributes struct {
	SSHPort     int    `json:"ssh_port"`
	SubdomainID string `json:"subdomain"`
	Public_Key  string `json:"ssh_key"`
}
type TunnelJsonData struct {
	Type       string           `json:"type"`
	Attributes TunnelAttributes `json:"attributes"`
	ID         string           `json:"id"`
}
type OpenTunnelRequest struct {
	Data TunnelJsonData `json:"data"`
}
type OpenTunnelResponse struct {
	Data TunnelJsonData `json:"data"`
}
type TunnelsListResponse struct {
	Data []TunnelJsonData `json:"data"`
}

//TunnelListAPI get list of subdomains reserved
func (restClient *RestClient) TunnelListAPI() (TunnelsListResponse, error) {
	responseBody := TunnelsListResponse{}
	url := restClient.URL + "/tunnel"
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

//CreateTunnelAPI calls holepunch web api to get tunnel details
func (restClient *RestClient) CreateTunnelAPI(requestBody OpenTunnelRequest) (OpenTunnelResponse, error) {
	responseBody := OpenTunnelResponse{}
	url := restClient.URL + "/tunnel"
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

	fmt.Println(string(body))
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		fmt.Println("error:", err)
		return responseBody, nil
	}
	return responseBody, nil
}

//DeleteTunnelAPI deletes tunnel
func (restClient *RestClient) DeleteTunnelAPI(subdomainName string) error {
	url := restClient.URL + "/tunnel/" + subdomainName
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
	fmt.Println(string(body))
	responseBody := OpenTunnelResponse{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		fmt.Println("error:", err)
		return http.ErrAbortHandler
	}
	return nil
}
