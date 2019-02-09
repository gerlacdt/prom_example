package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gerlacdt/prom_example/pkg/domain"
)

// PostHandler type for all http routes
type PostHandler struct {
	postService domain.PostService
}

// New creates a new http PostHandler
func New(postService domain.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// ServeHTTP ...
func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		h.getPost(w, r)
	} else if r.Method == "POST" {
		h.createPost(w, r)
	} else if r.Method == "DELETE" {
		h.deletePost(w, r)
	} else {
		handleError(w, fmt.Errorf("normal error"))
	}
}

// Handle from http package
func Handle(pattern string, handler http.Handler) {
	http.DefaultServeMux.Handle(pattern, handler)
}

// ListenAndServe http server
func ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, nil)
}

type myError struct {
	err        error
	statusCode int
}

func (e *myError) Error() string {
	return fmt.Sprintf("myError string, %v", e.err)
}

func handleError(w http.ResponseWriter, err error) {
	switch v := err.(type) {
	case *myError:
		w.WriteHeader(v.statusCode)
		_, err := w.Write([]byte(fmt.Sprintf("{\"message\": %s}", err)))
		if err != nil {
			w.WriteHeader(500)
		}
	default:
		w.WriteHeader(500)
		_, err := w.Write([]byte(fmt.Sprintf("{\"message\": %s}", err)))
		if err != nil {
			// ignore
		}
	}
}

func getID(r *http.Request) (int, error) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/posts/")
	return strconv.Atoi(id)

}

func (h *PostHandler) getPost(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		handleError(w, &myError{err: fmt.Errorf("ID could not be read"),
			statusCode: 500})
		return
	}
	post, err := h.postService.Post(id)
	if err != nil {
		handleError(w, &myError{err: err, statusCode: 404})
		return
	}
	data, err := json.Marshal(post)
	if err != nil {
		handleError(w, &myError{err: err, statusCode: 500})
		return
	}
	_, err = w.Write(data)
	if err != nil {
		handleError(w, &myError{err: fmt.Errorf("Body could not be read"),
			statusCode: 500})
		return
	}
}

func (h *PostHandler) createPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, &myError{err: fmt.Errorf("Body could not be read, %s", err),
			statusCode: 400})
		return
	}
	var post domain.Post
	err = json.Unmarshal(body, &post)
	if err != nil {
		handleError(w, &myError{err: fmt.Errorf("Body could not be json parsed: %s", err),
			statusCode: 400})
	}
	err = h.postService.CreatePost(&post)
	if err != nil {
		handleError(w, &myError{err: fmt.Errorf("Could not store post"),
			statusCode: 500})
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *PostHandler) deletePost(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		handleError(w, &myError{err: fmt.Errorf("Error parsing id, %s", err)})
		return
	}
	err = h.postService.DeletePost(id)
	if err != nil {
		handleError(w, &myError{err: fmt.Errorf("Error during deletion, %s", err)})
		return
	}
	w.WriteHeader(http.StatusOK)
}
