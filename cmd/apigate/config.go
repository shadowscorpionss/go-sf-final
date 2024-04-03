package main

import (
	"ApiGate/package/api"
	"os"
	"strconv"
)

type Config struct {
	Censor   Censor
	Comments Comments
	News     News
	Gateway  Gateway
}

type Censor struct {
	AdrPort string
	URLdb   string
}

type Comments struct {
	AdrPort string
	URLdb   string
}

type News struct {
	AdrPort string
	URLdb   string
}

type Gateway struct {
	AdrPort string
}

// New возвращает новую Config структуру
func NewConfig() *Config {
	return &Config{

		Censor: Censor{
			AdrPort: getEnv("CENSOR_PORT", "8083"),
			URLdb:   getEnv("CENSOR_DB", ""),
		},
		Comments: Comments{
			AdrPort: getEnv("COMMENTS_PORT", "8082"),
			URLdb:   getEnv("COMMENTS_DB", ""),
		},
		News: News{
			AdrPort: getEnv("NEWS_PORT", "8081"),
			URLdb:   getEnv("NEWS_DB", ""),
		},
		Gateway: Gateway{
			AdrPort: getEnv("GATEWAY_PORT", "8080"),
		},
	}
}

func CtoApiConfig(c Config) *api.ApiGatewayConfig {
	apicfg := api.ApiGatewayConfig{
		GatewayPort:  stoid(c.Gateway.AdrPort, 8080),
		CensorPort:   stoid(c.Censor.AdrPort, 8083),
		CommentsPort: stoid(c.Comments.AdrPort, 8082),
		NewsPort:     stoid(c.News.AdrPort, 8081),
	}
	return &apicfg
}
func stoi(port string) (int, bool) {
	if len(port) == 0 {
		return 0, false
	}

	if a, err := strconv.Atoi(port); err == nil && a > 0 {
		return a, true
	}
	return 0, false
}

func stoid(port string, def int) int {
	a, b := stoi(port)
	if !b {
		return def
	}
	return a
}

// Простая вспомогательная функция для считывания окружения или возврата значения по умолчанию
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
