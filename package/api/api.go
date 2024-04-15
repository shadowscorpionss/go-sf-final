package api

import (
	"ApiGate/package/middleware"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/gorilla/mux"
)

type API struct {
	r         *mux.Router
	censContr string
	comBase   string
	comContr  string
	newsBase  string
	newsContr string
	newsSrch  string
}

const (
	NEWS_PATH        = "/news"
	NEWS_SEARCH_PATH = NEWS_PATH + "/search"
	NEWS_LATEST_PATH = NEWS_PATH + "/latest"
	COMMENTS_PATH    = "/comments"
	CENSOR_PATH      = "/censor"
)

// constructor of API
func New(cfg ApiGatewayConfig) *API {

	a := API{r: mux.NewRouter()}
	a.comBase = HttpBaseUrl(cfg.Comments)
	a.newsBase = HttpBaseUrl(cfg.News)
	a.censContr = ControllerUrl(HttpBaseUrl(cfg.Censor), CENSOR_PATH)
	a.comContr = ControllerUrl(a.comBase, COMMENTS_PATH)
	a.newsContr = ControllerUrl(a.newsBase, NEWS_PATH)
	a.newsSrch = ControllerUrl(a.newsBase, NEWS_SEARCH_PATH)

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
	api.r.HandleFunc(NEWS_PATH, api.news).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc(NEWS_LATEST_PATH, api.newsLatest).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc(NEWS_SEARCH_PATH, api.newsSearch).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc(COMMENTS_PATH, api.addComment).Methods(http.MethodPost, http.MethodOptions)
	api.r.HandleFunc(COMMENTS_PATH+"/{id}", api.deleteComment).Methods(http.MethodDelete, http.MethodOptions)

}

// GET /news
func (api *API) news(w http.ResponseWriter, r *http.Request) {
}

// GET /news/latest
func (api *API) newsLatest(w http.ResponseWriter, r *http.Request) {

}

// GET /news/search
func (api *API) newsSearch(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != NEWS_SEARCH_PATH {
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

	//calling news microservice
	go func() {
		defer wg.Done()
		resp1, err := http.Get(QueryUrl(api.newsSrch, map[string]string{"id": idParam}))
		chErr <- err
		chNews <- resp1
	}()

	//calling comments microservice
	go func() {
		defer wg.Done()
		resp2, err := http.Get(QueryUrl(api.comContr, map[string]string{"news_id": idParam}))
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
			body, err := io.ReadAll(respNews.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = json.Unmarshal(body, &response.NewsDetailed)
		case respComment := <-chComments:
			body, err := io.ReadAll(respComment.Body)
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
	if r.URL.Path != COMMENTS_PATH {
		http.NotFound(w, r)
	}
	//reading request body
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body.Close()

	//sending to check to censor service
	censorBody := io.NopCloser(bytes.NewBuffer(bodyBytes))
	respCensor, err := middleware.MakeRequest(r, http.MethodPost, api.censContr, censorBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if respCensor.StatusCode != 200 {
		http.Error(w, "inaccaptable comment", respCensor.StatusCode)
		return
	}

	commentBody := io.NopCloser(bytes.NewBuffer(bodyBytes))
	resp, err := middleware.MakeRequest(r, http.MethodPost, api.comContr, commentBody)
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
	if r.URL.Path != COMMENTS_PATH {
		http.NotFound(w, r)
	}

	// proxy to microservice
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   api.comBase,
	})

	// proxy
	proxy.ServeHTTP(w, r)
}
