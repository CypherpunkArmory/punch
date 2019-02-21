package tunnel

import (
	"fmt"

	"github.com/cypherpunkarmory/punch/restapi"
)

//Config Object to make passing config eaiser
type Config struct {
	RestAPI        restapi.RestClient
	TunnelEndpoint restapi.Tunnel
	PrivateKeyPath string
	LocalPort      int
	Subdomain      string
	EndpointType   string
	EndpointURL    string
}

type endpoint struct {
	Host string
	Port int
}

func (endpoint *endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}
