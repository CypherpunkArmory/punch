package tunnel

import (
	"fmt"
	"net/url"

	"github.com/cypherpunkarmory/punch/restapi"
)

//Config Object to make passing config eaiser
type Config struct {
	ConnectionEndpoint url.URL
	RestAPI            restapi.RestClient
	TunnelEndpoint     restapi.Tunnel
	PrivateKeyPath     string
	LocalPort          string
	Subdomain          string
	EndpointType       string
	EndpointURL        url.URL
	LogLevel           string
}

type Endpoint struct {
	Host string
	Port string
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%s", e.Host, e.Port)
}
