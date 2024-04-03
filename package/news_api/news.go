package newsapi

import (
	"ApiGate/package/middleware"
	"ApiGate/package/models"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type NewsApi struct {
}

func (c *NewsApi) addRoutes(r *mux.Router) {
	s := r.PathPrefix("/news").Subrouter()
	//get posts list
	s.HandleFunc("/", c.posts).Methods(http.MethodGet, http.MethodOptions)
	//create post
	s.HandleFunc("/", c.addPost).Methods(http.MethodPost, http.MethodOptions)
	//get post by id
	s.HandleFunc("/{id}", c.post).Methods(http.MethodGet, http.MethodOptions)
	//delete post
	s.HandleFunc("/{id}", c.deletePost).Methods(http.MethodDelete, http.MethodOptions)
	//update post
	s.HandleFunc("/{id}", c.updatePost).Methods(http.MethodPut, http.MethodOptions)
}

// GET /posts/
func (c *NewsApi) posts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET posts")
	res := []models.NewsShortDetailed{
		{
			Id:      0,
			Title:   "Test",
			PubTime: 0,
			Link:    "http://localhost:8080/",
		},
	}
	middleware.RespondWithJSON(w, http.StatusAccepted, res)
}

// GET /posts/{id}
func (c *NewsApi) post(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET post")

	strId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(strId)

	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	res := models.NewsShortDetailed{
		Id:      id,
		Title:   "Test response",
		PubTime: 0,
		Link:    "http://localhost:8080/",
	}

	middleware.RespondWithJSON(w, http.StatusAccepted, res)
}

// DELETE /posts/{id}
func (c *NewsApi) deletePost(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (c *NewsApi) updatePost(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// POST /posts/ - adds post
func (api *NewsApi) addPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("add post")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var c models.NewsShortDetailed
	if err := decoder.Decode(&c); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	c.Id = rand.Intn(99) + 1

	middleware.RespondWithJSON(w, http.StatusAccepted, c)

}
