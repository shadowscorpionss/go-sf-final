package main

import (
	"ApiGate/package/api"
	"ApiGate/package/storage"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	// initialization from environment to keep safe your secrets =)
	connstr := os.Getenv("apigatedb")

	_, exeFilename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("unable to get the current filename")
	}

	homePath := filepath.Dir(exeFilename)

	// --db
	db, err := storage.New(connstr)
	if err != nil {
		log.Fatal(err)
	}

	// --api
	api := api.New(db, homePath)

	// start http server
	err = http.ListenAndServe(":8080", api.Router())
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe("localhost:8080", api.Router())

}
