package restapi_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cypherpunkarmory/punch/restapi"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
)

var accessToken string = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE1NDY0NzI5NDYsIm5iZiI6MTU0NjQ3Mjk0NiwianRpIjoiZDBlNDQzZjYtY2Q4YS00YzIzLThhZTAtMjM0MWVkZmRmMWVjIiwiZXhwIjoxNTQ5MDY0OTQ2LCJpZGVudGl0eSI6InRlc3RAbG9uZG9udHJ1c3RtZWRpYS5jb20iLCJ0eXBlIjoicmVmcmVzaCJ9.MCCBlKs6hpPHFZ_citooalnQA2hoq7MBm2dSZnjdu5k"

func TestTunnelCreateAPI(t *testing.T) {
	r, err := recorder.New("./fixtures/tunnel_create")

	if err != nil {
		t.Fatal(err)
	}

	defer r.Stop()

	recording_client := http.Client{
		Transport: r,
	}

	client := restapi.RestClient{
		URL:    "http://127.0.0.1:5000",
		Client: recording_client,
	}

	client.StartSession(accessToken)

	tunnel, _ := client.CreateTunnelAPI("")

	schemaLoader := gojsonschema.NewReferenceLoader("file://./schemas/tunnel.json")
	jsonified := gojsonschema.NewGoLoader(tunnel)

	result, err := gojsonschema.Validate(schemaLoader, jsonified)

	if err != nil {
		fmt.Println(err)
	}

	require.True(t, result.Valid())
	require.NotNil(t, subdomain.Data.Attributes.Subdomain)
}

func TestTunnelDeleteAPI(t *testing.T) {
	r, err := recorder.New("./fixtures/tunnel_delete")

	if err != nil {
		t.Fatal(err)
	}

	defer r.Stop()

	recording_client := http.Client{
		Transport: r,
	}

	client := restapi.RestClient{
		URL:    "http://127.0.0.1:5000",
		Client: recording_client,
	}

	client.StartSession(accessToken)
	tunnel, _ := client.CreateTunnelAPI("")
	err = client.DeleteTunnelAPI(tunnel.Data.Attributes.Subdomain)

	require.Nil(t, err)
}
