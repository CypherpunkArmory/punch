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

func TestSubdomainListAPI(t *testing.T) {
	r, err := recorder.New("./fixtures/subdomain")

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

	subList, _ := client.ListSubdomainAPI()

	schemaLoader := gojsonschema.NewReferenceLoader("file://./schemas/subdomains.json")
	jsonified := gojsonschema.NewGoLoader(subList)

	result, err := gojsonschema.Validate(schemaLoader, jsonified)

	if err != nil {
		fmt.Println(err)
	}

	require.True(t, result.Valid())
}

func TestSubdomainReserveAPI(t *testing.T) {
	r, err := recorder.New("./fixtures/reserve")

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

	subdomain, _ := client.ReserveSubdomainAPI("iamatestdomain")

	schemaLoader := gojsonschema.NewReferenceLoader("file://./schemas/subdomain.json")
	jsonified := gojsonschema.NewGoLoader(subdomain)

	result, err := gojsonschema.Validate(schemaLoader, jsonified)

	if err != nil {
		fmt.Println(err)
	}

	require.True(t, result.Valid())
	require.True(t, subdomain.Data.Attributes.Reserved)
	require.False(t, subdomain.Data.Attributes.InUse)
	require.Equal(t, subdomain.Data.Attributes.Name, "iamatestdomain")
}

func TestSubdomainReleaseAPI(t *testing.T) {
	r, err := recorder.New("./fixtures/release")

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

	err = client.ReleaseSubdomainAPI("sillymonkey")

	require.Nil(t, err)
}
