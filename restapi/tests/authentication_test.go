package restapi_test

import (
	"github.com/cypherpunkarmory/punch/restapi"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRestAPILogin(t *testing.T) {
	r, err := recorder.New("./fixtures/login")

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

	res, err := client.Login("me@stephenprater.com", "secret")

	if err != nil {
		t.Fatal("Login failed")
	}

	require.NotNil(t, res.Access_Token)
	require.NotEqual(t, client.APIKEY, "")
}

func TestRestAPISession(t *testing.T) {
	r, err := recorder.New("./fixtures/refresh")

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
		APIKEY: "yadayadayada",
	}

	// let's call this "expirementally verified"
	res, err := client.StartSession("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE1NDU5ODEzNjYsIm5iZiI6MTU0NTk4MTM2NiwianRpIjoiYzcxOTgxNWQtNWYwMi00MDhlLWI1MDctYTAyMDA5OTNhMzgwIiwiZXhwIjoxNTQ4NTczMzY2LCJpZGVudGl0eSI6Im1lQHN0ZXBoZW5wcmF0ZXIuY29tIiwidHlwZSI6InJlZnJlc2gifQ.Zl6sdAcXwdAxlMAUzQ_RpFgf0FX69mS5JRpUZir8bn8")

	if err != nil {
		t.Fatal("Refresh failed")
	}

	require.NotNil(t, res.Access_Token)
	require.NotEqual(t, client.APIKEY, "yadayadayada")
}

func TestRestAPISessionNoRefresh(t *testing.T) {
	client := restapi.RestClient{
		URL:    "http://127.0.0.1:5000",
		APIKEY: "",
	}

	_, err := client.StartSession("")

	require.NotNil(t, err)
}
