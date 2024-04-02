package api

import (
	"ApiGate/package/storage"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	r  *mux.Router
	hp string //home path
	//c  CommentsApi
	n  newsApi
	db *storage.DB
}

// constructor of API
func New(db *storage.DB, homepath string) *API {
	a := API{r: mux.NewRouter(), hp: homepath, db: db}
	//a.c = CommentsApi{}
	a.r.StrictSlash(true)
	a.endpoints()
	return &a
}

// returns router for HTTP Server
func (api *API) Router() *mux.Router {
	return api.r
}

// register endpoints
func (api *API) endpoints() {
	//fmt.Println("adding comments routes")
	//api.c.addRoutes(api.r, api.db)
	fmt.Println("adding news routes")
	api.n.addRoutes(api.r)
	// html web server
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(api.hp+"/webapp"))))

}
