package api

import (
	"ApiGate/package/middleware"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/gorilla/mux"
)

type API struct {
	r       *mux.Router
	cens    string
	coms    string
	nws     string
	comBase string
}

// constructor of API
func New(cfg ApiGatewayConfig) *API {
	a := API{r: mux.NewRouter(),
		cens:    NewEndpointConfig(cfg.Censor, "censor", ""),
		coms:    NewEndpointConfig(cfg.Comments, "comments", ""),
		nws:     NewEndpointConfig(cfg.News, "news", "search"),
		comBase: fmt.Sprintf("%s:%d", cfg.Comments.Host, cfg.Comments.Port),
	}

	a.r.StrictSlash(true)
	a.addRoutes()
	return &a
}

// returns router for HTTP Server
func (api *API) Router() *mux.Router {
	return api.r
}

// register addRoutes
func (api *API) addRoutes() {
	api.r.HandleFunc("/news", api.news).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/news/latest", api.newsLatest).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/news/search", api.newsSearch).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/comments", api.addComment).Methods(http.MethodPost, http.MethodOptions)
	api.r.HandleFunc("/comments/{id}", api.deleteComment).Methods(http.MethodDelete, http.MethodOptions)

}

// GET /news
func (api *API) news(w http.ResponseWriter, r *http.Request) {
}

// GET /news/latest
func (api *API) newsLatest(w http.ResponseWriter, r *http.Request) {
}

// GET /news/search
func (api *API) newsSearch(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/news/search" {
		http.NotFound(w, r)
	}

	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "id not set", http.StatusBadRequest)
		return
	}

	chNews := make(chan *http.Response, 2)
	chComments := make(chan *http.Response, 2)
	chErr := make(chan error, 2)
	var response ResponseDetailed
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// Отправляем запрос на порт 8081
		resp1, err := http.Get("http://localhost" + portNews + "/news/search" + "?id=" + idParam)
		chErr <- err
		chNews <- resp1
	}()

	go func() {
		defer wg.Done()
		// Отправляем запрос на порт 8082
		resp2, err := http.Get("http://localhost" + portComment + "/comments" + "?news_id=" + idParam)
		chErr <- err
		chComments <- resp2
	}()

	wg.Wait()
	close(chErr)

	for err := range chErr {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

block:
	for {
		select {
		case respNews := <-chNews:
			body, err := ioutil.ReadAll(respNews.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = json.Unmarshal(body, &response.NewsDetailed)
		case respComment := <-chComments:
			body, err := ioutil.ReadAll(respComment.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = json.Unmarshal(body, &response.Comments)
		default:
			break block
		}

	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// POST /comments
func (api *API) addComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comments" {
		http.NotFound(w, r)
	}
	//reading request body
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body.Close()

	//sending to check to censor service
	censorBody := io.NopCloser(bytes.NewBuffer(bodyBytes))
	respCensor, err := middleware.MakeRequest(r, http.MethodPost, api.cens, censorBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if respCensor.StatusCode != 200 {
		http.Error(w, "inaccaptable comment", respCensor.StatusCode)
		return
	}

	commentBody := io.NopCloser(bytes.NewBuffer(bodyBytes))
	resp, err := middleware.MakeRequest(r, http.MethodPost, api.coms, commentBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for name, values := range resp.Header {
		w.Header()[name] = values
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// DELETE /comments
func (api *API) deleteComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comments" {
		http.NotFound(w, r)
	}

	// proxy to microservice
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   api.comBase, // microservice
	})

	// proxy
	proxy.ServeHTTP(w, r)
}
