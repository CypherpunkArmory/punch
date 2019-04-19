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
	Port      []string   `jsonapi:"attr,port,omitempty"`
	PublicKey string     `jsonapi:"attr,sshKey,omitempty"`
	SSHPort   string     `jsonapi:"attr,sshPort,omitempty"`
	IPAddress string     `jsonapi:"attr,ipAddress,omitempty"`
	Subdomain *Subdomain `jsonapi:"relation,subdomain,omitempty"`
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
			Port:      protocol,
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
			Port:      protocol,
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
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errorCantConnectRestCall
	}
	resp, err := client.Do(req)
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
