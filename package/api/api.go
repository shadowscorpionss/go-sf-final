package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	r  *mux.Router
	hp string //home path
}

// constructor of API
func New(homepath string) *API {
	a := API{r: mux.NewRouter(), hp: homepath}
	a.endpoints()
	return &a
}

// returns router for HTTP Server
func (api *API) Router() *mux.Router {
	return api.r
}

// register endpoints
func (api *API) endpoints() {
	// get list of news
	api.r.HandleFunc("/news/list", api.posts).Methods(http.MethodGet, http.MethodOptions)
	// get filtered news
	api.r.HandleFunc("/news/list/{filter}", api.filterposts).Methods(http.MethodGet, http.MethodOptions)

	// get post by id
	api.r.HandleFunc("/news/{id}", api.post).Methods(http.MethodGet, http.MethodOptions)

	//add comment
	api.r.HandleFunc("/comments/", api.addcomment).Methods(http.MethodPost, http.MethodOptions)

	// html web server
	//api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(api.hp+"/webapp"))))

}

// GET /news/list - returns posts
func (api *API) posts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	json.NewEncoder(w).Encode("")
}

// GET /news/list/{filter} - returns filtered posts
func (api *API) filterposts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	filter := mux.Vars(r)["filter"]
	_ = filter
	json.NewEncoder(w).Encode("")
}

// GET /news/{id} - returns post by id
func (api *API) post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	id := mux.Vars(r)["id"]
	_ = id
	json.NewEncoder(w).Encode("")
}

// POST /comments/ - adds comment
func (api *API) addcomment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	//json.NewEncoder(w).Encode(&post)
}
