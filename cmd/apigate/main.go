package main

import (
	"ApiGate/package/api"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
)

func main() {

	_, exeFilename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("unable to get the current filename")
	}

	homePath := filepath.Dir(exeFilename)
	// --api
	api := api.New(homePath)

	// start http server
	err := http.ListenAndServe(":8080", api.Router())
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe("localhost:8080", api.Router())

}
