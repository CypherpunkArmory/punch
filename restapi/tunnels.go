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
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/google/jsonapi"
)

//Tunnel JSONAPI response of tunnel object
type Tunnel struct {
	ID        string     `jsonapi:"primary,tunnel"`
	PortTypes []string   `jsonapi:"attr,port,omitempty"`
	PublicKey string     `jsonapi:"attr,sshKey,omitempty"`
	SSHPort   string     `jsonapi:"attr,sshPort,omitempty"`
	TCPPorts  []string   `jsonapi:"attr,tcpPorts,omitempty"`
	IPAddress string     `jsonapi:"attr,ipAddress,omitempty"`
	Subdomain *Subdomain `jsonapi:"relation,subdomain,omitempty"`
}

//ListTunnelsAPI get list of tunnels
func (restClient *RestClient) ListTunnelsAPI() ([]Tunnel, error) {
	tunnelList := []Tunnel{}
	url := restClient.URL + "/tunnels"
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
			return tunnelList, err
		}
		return tunnelList, errObject
	}
	responseBody, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Tunnel)))
	if err != nil {
		return tunnelList, errorUnableToParse
	}
	for _, tunnel := range responseBody {
		s, _ := tunnel.(*Tunnel)
		tunnelList = append(tunnelList, *s)
	}

	return tunnelList, nil
}

//CreateTunnelAPI calls holepunch web api to get tunnel details
func (restClient *RestClient) CreateTunnelAPI(subdomain string, publicKey string, protocol []string) (Tunnel, error) {
	tunnelReturn := Tunnel{}
	var outputBuffer bytes.Buffer

	if subdomain != "" {
		subdomainID, err := restClient.getSubdomainID(subdomain)
		if err != nil {
			return tunnelReturn, errorUnownedSubdomain
		}
		request := Tunnel{
			PortTypes: protocol,
			PublicKey: publicKey,
			Subdomain: &Subdomain{
				ID: subdomainID,
			},
		}
		_ = bufio.NewWriter(&outputBuffer)
		err = jsonapi.MarshalPayload(&outputBuffer, &request)
		if err != nil {
			return tunnelReturn, errorUnableToParse
		}
	} else {
		request := Tunnel{
			PortTypes: protocol,
			PublicKey: publicKey,
		}

		_ = bufio.NewWriter(&outputBuffer)
		err := jsonapi.MarshalPayload(&outputBuffer, &request)
		if err != nil {
			return tunnelReturn, errorUnableToParse
		}
	}

	url := restClient.URL + "/tunnels"
	req, err := http.NewRequest("POST", url, &outputBuffer)
	if err != nil {
		return tunnelReturn, errorCantConnectRestCall
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := restClient.Client.Do(req)
	if err != nil {
		return tunnelReturn, errorCantConnectRestCall
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		buf, _ := ioutil.ReadAll(resp.Body)
		errObject := ResponseError{}
		err = json.Unmarshal(buf, &errObject)
		if err != nil {
			return tunnelReturn, err
		}
		return tunnelReturn, &errObject
	}

	err = jsonapi.UnmarshalPayload(resp.Body, &tunnelReturn)
	if err != nil {
		return tunnelReturn, errorUnableToParse
	}
	return tunnelReturn, nil
}

//DeleteTunnelAPI deletes tunnel
func (restClient *RestClient) DeleteTunnelAPI(subdomainName string) error {
	id, err := restClient.getTunnelID(subdomainName)
	if err != nil {
		return err
	}

	url := restClient.URL + "/tunnels/" + id
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

	return errorUnableToDelete
}
func (restClient *RestClient) getTunnelID(subdomainName string) (string, error) {
	url := restClient.URL + "/tunnels?filter[subdomain][name]=" + subdomainName
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

	tunnels, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Tunnel)))
	if err != nil {
		return "", errorUnableToParse
	}
	if len(tunnels) == 0 {
		return "", errorUnownedTunnel
	}
	t, _ := tunnels[0].(*Tunnel)
	if t.ID == "" {
		return "", errorUnownedTunnel
	}
	return t.ID, nil
}
