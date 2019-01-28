package restapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/google/jsonapi"
)

type Subdomain struct {
	ID       string `jsonapi:"primary,subdomain"`
	Name     string `jsonapi:"attr,name"`
	InUse    bool   `jsonapi:"attr,inUse"`
	Reserved bool   `jsonapi:"attr,reserved"`
}

//SubdomainListAPI get list of subdomains reserved
func (restClient *RestClient) ListSubdomainAPI() ([]Subdomain, error) {
	subdomainList := []Subdomain{}
	url := restClient.URL + "/subdomains"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+restClient.APIKEY)
	resp, err := restClient.Client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 399 {
		//errorBody := ErrorResponse{}
		errObject := new(jsonapi.ErrorObject)
		err = jsonapi.UnmarshalPayload(resp.Body, errObject)
		return subdomainList, err
	}
	responseBody, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Subdomain)))
	if err != nil {
		fmt.Println(err)
		responseBody := []Subdomain{}
		return responseBody, http.ErrAbortHandler
	}
	for _, subdomain := range responseBody {
		s, _ := subdomain.(*Subdomain)
		subdomainList = append(subdomainList, *s)
	}

	return subdomainList, nil
}

//ReserveSubdomainAPI calls holepunch web api to get reserve a subdomain
func (restClient *RestClient) ReserveSubdomainAPI(subdomainName string) (Subdomain, error) {
	subdomainReturn := Subdomain{}
	url := restClient.URL + "/subdomains"

	request := Subdomain{
		Name: subdomainName,
	}
	var outputBuffer bytes.Buffer
	_ = bufio.NewWriter(&outputBuffer)
	err := jsonapi.MarshalPayload(&outputBuffer, &request)

	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}

	req, err := http.NewRequest("POST", url, &outputBuffer)
	req.Header.Add("Content-Type", "application/vnd.api+json")
	req.Header.Add("Authorization", "Bearer "+restClient.APIKEY)
	resp, err := restClient.Client.Do(req)

	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 399 {
		//errorBody := ErrorResponse{}
		buf, _ := ioutil.ReadAll(resp.Body)
		errObject := ResponseError{}
		err = json.Unmarshal(buf, &errObject)
		return subdomainReturn, errObject
	}
	err = jsonapi.UnmarshalPayload(resp.Body, &subdomainReturn)
	if err != nil {
		fmt.Println(err)
		responseBody := Subdomain{}
		return responseBody, http.ErrAbortHandler
	}
	return subdomainReturn, nil
}

//ReleaseSubodmainAPI deletes tunnel
func (restClient *RestClient) ReleaseSubdomainAPI(subdomainName string) error {
	id, err := restClient.getSubdomainID(subdomainName)

	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}

	if id == "" {
		return errors.New("You do not own this subdomain")
	}

	url := restClient.URL + "/subdomains/" + id
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", "Bearer "+restClient.APIKEY)
	resp, err := restClient.Client.Do(req)

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

func (restClient *RestClient) getSubdomainID(subdomainName string) (string, error) {
	url := restClient.URL + "/subdomains?filter[name]=" + subdomainName
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+restClient.APIKEY)
	resp, err := restClient.Client.Do(req)

	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 399 {
		//errorBody := ErrorResponse{}
		buf, _ := ioutil.ReadAll(resp.Body)
		errObject := ResponseError{}
		err = json.Unmarshal(buf, &errObject)
		return "-1", errObject
	}

	subdomains, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Subdomain)))
	if err != nil {
		fmt.Println(err)
		return "-1", http.ErrAbortHandler
	}
	for _, subdomain := range subdomains {
		s, _ := subdomain.(*Subdomain)
		if s.ID != "" {
			return s.ID, nil
		} else {
			return "", nil
		}
	}

	return "", nil
}
func (restClient *RestClient) GetSubdomainName(subdomainID string) (string, error) {
	url := restClient.URL + "/subdomains/" + subdomainID
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+restClient.APIKEY)
	resp, err := restClient.Client.Do(req)

	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 399 {
		//errorBody := ErrorResponse{}
		buf, _ := ioutil.ReadAll(resp.Body)
		errObject := ResponseError{}
		err = json.Unmarshal(buf, &errObject)
		return "", errObject
	}
	subdomain := new(Subdomain)
	err = jsonapi.UnmarshalPayload(resp.Body, subdomain)
	if err != nil {
		fmt.Println(err)
		return "", http.ErrAbortHandler
	}

	if subdomain.Name != "" {
		return subdomain.Name, nil
	} else {
		return "", nil
	}

	return "", nil
}
