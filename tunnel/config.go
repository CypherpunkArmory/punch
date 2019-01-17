package tunnel

import (
	"fmt"

	"github.com/cypherpunkarmory/punch/restapi"
)

type TunnelConfig struct {
	RestApi        restapi.RestClient
	TunnelEndpoint restapi.Tunnel
	PrivateKeyPath string
	LocalPort      int
	Subdomain      string
	EndpointType   string
	EndpointUrl    string
}
type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}
