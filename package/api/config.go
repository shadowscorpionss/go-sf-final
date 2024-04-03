package api

import "fmt"

type ApiGatewayConfig struct {
	Gateway  HostConfig
	Censor   HostConfig
	Comments HostConfig
	News     HostConfig
}

type HostConfig struct {
	Host string
	Port int
}


func NewEndpointConfig(cfg HostConfig, controller, method string) string {
	return fmt.Sprintf("http://%s:%d/%s/%s", cfg.Host, cfg.Port, controller, method)
}
