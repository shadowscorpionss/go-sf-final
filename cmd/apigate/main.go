package main

import (
	"ApiGate/package/api"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// init вызывается перед main()
func init() {

	// загружает значения из файла .env в систему
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	var err error

	// initialization from environment to keep safe your secrets =)

	cfg := api.ApiConfig{
		PortGate:    8080,
		PortCensor:  8083,
		PortComment: 8082,
		PortNews:    8081,
	}

	// --api
	api := api.New(cfg)

	// start http server
	err = http.ListenAndServe(":8080", api.Router())
	if err != nil {
		log.Fatal(err)
	}

}
