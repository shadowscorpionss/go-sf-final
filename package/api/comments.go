package api

import (
	"ApiGate/package/models"
	"ApiGate/package/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CommentsApi struct {
	db *storage.DB
	r  *mux.Router
}

func NewCommentsApi(db *storage.DB) *CommentsApi {
	// constructor of API
	a := CommentsApi{r: mux.NewRouter(), db: db}
	a.r.StrictSlash(true)
	a.addRoutes(a.r, db)
	return &a
}

// returns router for HTTP Server
func (c *CommentsApi) Router() *mux.Router {
	return c.r
}

func (c *CommentsApi) addRoutes(r *mux.Router, db *storage.DB) {
	s := r.PathPrefix("/comments").Subrouter()
	//get comments list
	s.HandleFunc("/{postId}", c.postComments).Methods(http.MethodGet, http.MethodOptions)
	//create comment
	s.HandleFunc("/", c.addComment).Methods(http.MethodPost, http.MethodOptions)
	//get comment by id
	s.HandleFunc("/by/{id}", c.comment).Methods(http.MethodGet, http.MethodOptions)
	//delete comment
	s.HandleFunc("/{id}", c.deleteComment).Methods(http.MethodDelete, http.MethodOptions)
	//update comment
	s.HandleFunc("/{id}", c.updateComment).Methods(http.MethodPut, http.MethodOptions)
}

// GET /comments/{postId}
func (c *CommentsApi) postComments(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET comments by postId")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	sPostId := mux.Vars(r)["postId"]
	postId, err := strconv.Atoi(sPostId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	res, err := c.db.Comments(postId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error while reading from db: %v", err))
		return
	}
	respondWithJSON(w, http.StatusAccepted, res)
}

// GET /comments/by/{id}
func (c *CommentsApi) comment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET comment by Id")
	respondWithError(w, http.StatusNotImplemented, "not implemented")
}

// DELETE /comments/{id}
func (ca *CommentsApi) deleteComment(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusNotImplemented, "not implemented")
}

func (ca *CommentsApi) updateComment(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusNotImplemented, "not implemented")
}

// POST /comments/ - adds comment
func (ca *CommentsApi) addComment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("add comment")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var nc models.NewComment
	if err := decoder.Decode(&nc); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	res, err := ca.db.NewComment(nc)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	respondWithJSON(w, http.StatusAccepted, res)

}
