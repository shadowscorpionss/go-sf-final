package main

import (
	commentsapi "ApiGate/package/comments_api"
	"ApiGate/package/storage"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// initialization from environment to keep safe your secrets =)
	connstr := os.Getenv("commentsdb")
	port:=os.Getenv("commentsport")

	// --db
	db, err := storage.New(connstr)
	if err != nil {
		log.Fatal(err)
	}

	// --api
	api := commentsapi.New(db)

	// start http server
	err = http.ListenAndServe(fmt.Sprintf(":%s",port), api.Router())
	if err != nil {
		log.Fatal(err)
	}

}
