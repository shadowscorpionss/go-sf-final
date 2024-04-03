package api

import (
	"ApiGate/package/middleware"
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

type API struct {
	r           *mux.Router
	portCensor  string
	portComment string
}

// constructor of API
func New(cfg ApiGatewayConfig) *API {
	a := API{r: mux.NewRouter()}

	a.portCensor = strconv.Itoa(cfg.CensorPort)
	a.portComment = strconv.Itoa(cfg.CommentsPort)

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
	api.r.HandleFunc("/news", api.news).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/news/latest", api.newsLatest).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/news/search", api.newsSearch).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/comments", api.addComment).Methods(http.MethodPost, http.MethodOptions)
	api.r.HandleFunc("/comments", api.deleteComment).Methods(http.MethodDelete, http.MethodOptions)

}

// GET /news
func (api *API) news(w http.ResponseWriter, r *http.Request) {
}

// GET /news/latest
func (api *API) newsLatest(w http.ResponseWriter, r *http.Request) {
}

// GET /news/search
func (api *API) newsSearch(w http.ResponseWriter, r *http.Request) {
}

// POST /comments
func (api *API) addComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comments/add" {
		http.NotFound(w, r)
	}
	portCensor := api.portCensor
	portComment := api.portComment

	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body.Close()
	Body1 := io.NopCloser(bytes.NewBuffer(bodyBytes))
	Body := io.NopCloser(bytes.NewBuffer(bodyBytes))

	respCensor, err := middleware.MakeRequest(r, http.MethodPost, "http://localhost:"+portCensor+"/comments/check", Body1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if respCensor.StatusCode != 200 {
		http.Error(w, "неправильное содержание комментария", respCensor.StatusCode)
		return
	}

	resp, err := middleware.MakeRequest(r, http.MethodPost, "http://localhost:"+portComment+"/comments/add", Body)
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
	portComment := ":9090" //api.portComment

	// Создаем прокси-сервер для первого микросервиса
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost" + portComment, // адрес микросервиса
	})

	// Проксируем запрос к первому микросервису
	proxy.ServeHTTP(w, r)
}
