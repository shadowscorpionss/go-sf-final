package api

import (
	"fmt"
	"strings"
)

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

func HttpBaseUrl(cfg HostConfig) string {
	return fmt.Sprintf("http://%s:%d/", cfg.Host, cfg.Port)
}

func ControllerUrl(hostBaseUrl, controller string) string {
	return fmt.Sprintf("%s/%s", hostBaseUrl, strings.TrimLeft(controller, "/"))
}

func MethodUrl(controllerUrl, method string) string {
	return ControllerUrl(controllerUrl, method)
}

func QueryUrl(controlerUrl string, queryParams map[string]string) string {
	queryStrs := []string{}
	for k, v := range queryParams {
		queryStrs = append(queryStrs, fmt.Sprintf("%s=%s", k, v))
	}
	query := strings.Join(queryStrs, "&")
	if len(query) == 0 {
		return controlerUrl
	}
	return fmt.Sprintf("%s?%s", controlerUrl, query)
}
