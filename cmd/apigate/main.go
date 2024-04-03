package main

import (
	"ApiGate/package/api"
	"fmt"
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
	cfg:=NewConfig()
	apicfg:=CtoApiConfig(*cfg)

	// --api
	api := api.New(*apicfg)

	// start http server
	err = http.ListenAndServe(fmt.Sprintf(":%d", apicfg.GatewayPort), api.Router())
	if err != nil {
		log.Fatal(err)
	}

}
