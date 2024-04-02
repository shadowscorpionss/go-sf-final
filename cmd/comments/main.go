package main

import (
	"ApiGate/package/api"
	"ApiGate/package/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	// initialization from environment to keep safe your secrets =)
	connstr := os.Getenv("commentsdb")

	// --db
	db, err := storage.New(connstr)
	if err != nil {
		log.Fatal(err)
	}

	// --api
	api := api.NewCommentsApi(db)

	// start http server
	err = http.ListenAndServe(":8181", api.Router())
	if err != nil {
		log.Fatal(err)
	}

}
