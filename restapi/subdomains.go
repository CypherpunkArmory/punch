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
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/google/jsonapi"
)

//Subdomain Subdomain object that holds all needed info
type Subdomain struct {
	ID       string `jsonapi:"primary,subdomain"`
	Name     string `jsonapi:"attr,name"`
	InUse    bool   `jsonapi:"attr,inUse"`
	Reserved bool   `jsonapi:"attr,reserved"`
}

//ListSubdomainAPI get list of subdomains reserved
func (restClient *RestClient) ListSubdomainAPI() ([]Subdomain, error) {
	subdomainList := []Subdomain{}
	url := restClient.URL + "/subdomains"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errorCantConnectRestCall
	}
	resp, err := restClient.Client.Do(req)
	if err != nil {
		return nil, errorCantConnectRestCall
	}
	defer resp.Body.Close()
	if resp.StatusCode > 399 {
		errObject := new(jsonapi.ErrorObject)
		err = jsonapi.UnmarshalPayload(resp.Body, errObject)
		if err != nil {
			return subdomainList, err
		}
		return subdomainList, errObject
	}
	responseBody, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Subdomain)))
	if err != nil {
		return subdomainList, errorUnableToParse
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
		return subdomainReturn, errorUnableToParse
	}

	req, err := http.NewRequest("POST", url, &outputBuffer)
	if err != nil {
		return subdomainReturn, errorCantConnectRestCall
	}
	req.Header.Add("Content-Type", "application/vnd.api+json")
	resp, err := restClient.Client.Do(req)

	if err != nil {
		return subdomainReturn, errorCantConnectRestCall
	}
	defer resp.Body.Close()
	if resp.StatusCode > 399 {
		buf, _ := ioutil.ReadAll(resp.Body)
		errObject := ResponseError{}
		err = json.Unmarshal(buf, &errObject)
		if err != nil {
			return subdomainReturn, err
		}
		return subdomainReturn, &errObject
	}
	err = jsonapi.UnmarshalPayload(resp.Body, &subdomainReturn)
	if err != nil {
		return subdomainReturn, errorUnableToParse
	}
	return subdomainReturn, nil
}

//ReleaseSubdomainAPI deletes subdomain
func (restClient *RestClient) ReleaseSubdomainAPI(subdomainName string) error {
	id, err := restClient.getSubdomainID(subdomainName)

	if err != nil {
		return err
	}

	url := restClient.URL + "/subdomains/" + id
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errorCantConnectRestCall
	}
	resp, err := restClient.Client.Do(req)
	if err != nil {
		return errorCantConnectRestCall
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		return nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode > 399 {
		errorBody := ResponseError{}
		err = json.Unmarshal(body, &errorBody)
		if err != nil {
			return err
		}
		return &errorBody
	}

	return errors.New("failed to delete")
}

//GetSubdomainName Returns subdomain name of a given subdomain id
func (restClient *RestClient) GetSubdomainName(subdomainID string) (string, error) {
	url := restClient.URL + "/subdomains/" + subdomainID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errorCantConnectRestCall
	}
	resp, err := restClient.Client.Do(req)
	if err != nil {
		return "", errorCantConnectRestCall
	}
	defer resp.Body.Close()
	if resp.StatusCode > 399 {
		buf, _ := ioutil.ReadAll(resp.Body)
		errObject := ResponseError{}
		err = json.Unmarshal(buf, &errObject)
		if err != nil {
			return "", err
		}
		return "", &errObject
	}
	subdomain := new(Subdomain)
	err = jsonapi.UnmarshalPayload(resp.Body, subdomain)
	if err != nil {
		return "", http.ErrAbortHandler
	}

	if subdomain.Name != "" {
		return subdomain.Name, nil
	}
	return "", errorUnownedSubdomain
}

func (restClient *RestClient) getSubdomainID(subdomainName string) (string, error) {
	url := restClient.URL + "/subdomains?filter[name]=" + subdomainName
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errorCantConnectRestCall
	}
	resp, err := restClient.Client.Do(req)
	if err != nil {
		return "", errorCantConnectRestCall
	}
	defer resp.Body.Close()
	if resp.StatusCode > 399 {
		buf, _ := ioutil.ReadAll(resp.Body)
		errObject := ResponseError{}
		err = json.Unmarshal(buf, &errObject)
		return "", err
	}

	subdomains, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Subdomain)))
	if err != nil {
		return "", errorUnableToParse
	}
	if len(subdomains) == 0 {
		return "", errorUnownedSubdomain
	}
	s, _ := subdomains[0].(*Subdomain)
	if s.ID == "" {
		return "", errorUnownedSubdomain
	}
	return s.ID, nil
}
