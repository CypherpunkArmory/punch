package restapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/google/jsonapi"
)

type Tunnel struct {
	ID        string     `jsonapi:"primary,tunnel"`
	Port      string     `jsonapi:"attr,port,omitempty"`
	PublicKey string     `jsonapi:"attr,sshKey,omitempty"`
	SSHPort   string     `jsonapi:"attr,sshPort,omitempty"`
	IPAddress string     `jsonapi:"attr,ipAddress,omitempty"`
	Subdomain *Subdomain `jsonapi:"relation,subdomain,omitempty"`
}

//SubdomainListAPI get list of subdomains reserved
func (restClient *RestClient) listTunnelAPI() ([]Tunnel, error) {
	tunnelList := []Tunnel{}
	url := restClient.URL + "/tunnels"
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
		return tunnelList, err
	}
	responseBody, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Tunnel)))
	if err != nil {
		fmt.Println(err.Error())
		responseBody := []Tunnel{}
		return responseBody, http.ErrAbortHandler
	}
	for _, tunnel := range responseBody {
		t, _ := tunnel.(*Tunnel)
		tunnelList = append(tunnelList, *t)
	}

	return tunnelList, nil
}

//CreateTunnelAPI calls holepunch web api to get tunnel details
func (restClient *RestClient) CreateTunnelAPI(subdomain string, publicKey string, protocol string) (Tunnel, error) {
	tunnelReturn := Tunnel{}
	var outputBuffer bytes.Buffer
	if subdomain != "" {
		subdomainID, err := restClient.getSubdomainID(subdomain)
		if err != nil {
			fmt.Println("Cant connect to holepunch api")
			os.Exit(0)
		}
		if subdomainID == "" {
			fmt.Println("You do not own this domain")
			os.Exit(0)
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
			fmt.Println("error:", err)
			panic(err)
		}
	} else {
		request := Tunnel{
			Port:      protocol,
			PublicKey: publicKey,
		}

		_ = bufio.NewWriter(&outputBuffer)
		err := jsonapi.MarshalPayload(&outputBuffer, &request)

		if err != nil {
			fmt.Println("error:", err)
			panic(err)
		}
	}
	url := restClient.URL + "/tunnels"
	req, err := http.NewRequest("POST", url, &outputBuffer)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+restClient.APIKEY)
	resp, err := restClient.Client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 399 {

		buf, _ := ioutil.ReadAll(resp.Body)
		errObject := ResponseError{}
		err = json.Unmarshal(buf, &errObject)
		return tunnelReturn, errObject
	}
	err = jsonapi.UnmarshalPayload(resp.Body, &tunnelReturn)
	if err != nil {
		fmt.Println(err)
		responseBody := Tunnel{}
		return responseBody, http.ErrAbortHandler
	}
	return tunnelReturn, nil
}

//DeleteTunnelAPI deletes tunnel
func (restClient *RestClient) DeleteTunnelAPI(subdomainName string) error {
	id, err := restClient.getTunnelID(subdomainName)

	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	if id == "" {
		return errors.New("You do not own this subdomain")
	}

	url := restClient.URL + "/tunnels/" + id
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", "Bearer "+restClient.APIKEY)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

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
func (restClient *RestClient) getTunnelID(subdomainName string) (string, error) {
	url := restClient.URL + "/tunnels?filter[subdomain][name]=" + subdomainName
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+restClient.APIKEY)
	resp, err := restClient.Client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 399 {

		buf, _ := ioutil.ReadAll(resp.Body)
		errObject := ResponseError{}
		err = json.Unmarshal(buf, &errObject)
		return "", errObject
	}

	tunnels, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Tunnel)))
	if err != nil {
		fmt.Println(err)
		return "", http.ErrAbortHandler
	}
	for _, tunnel := range tunnels {
		t, _ := tunnel.(*Tunnel)
		if t.ID != "" {
			return t.ID, nil
		} else {
			return "", nil
		}
	}

	return "", nil
}
