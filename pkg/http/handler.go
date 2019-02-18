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
	random      *Random
}

// New creates a new http PostHandler
func New(postService domain.PostService) *PostHandler {
	random := NewRandom()
	return &PostHandler{postService: postService, random: random}
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

// Handler ...
type Handler http.Handler

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
		msg := domain.HTTPError{Message: err.Error()}
		err := json.NewEncoder(w).Encode(msg)
		if err != nil {
			w.WriteHeader(500)
		}
	default:
		w.WriteHeader(500)
		msg := domain.HTTPError{Message: err.Error()}
		err := json.NewEncoder(w).Encode(msg)
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
	if h.random.randomError(5) {
		handleError(w, &myError{err: fmt.Errorf("GET /v1/posts/:id Random injected Error"),
			statusCode: 500})
		return
	}
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
	// in order to have some fake stats for prometheues
	h.random.randomSleep(30, 110)
	_, err = w.Write(data)
	if err != nil {
		handleError(w, &myError{err: fmt.Errorf("Body could not be read"),
			statusCode: 500})
		return
	}
}

func (h *PostHandler) createPost(w http.ResponseWriter, r *http.Request) {
	if h.random.randomError(10) {
		handleError(w, &myError{err: fmt.Errorf("POST /v1/posts Random injected Error"),
			statusCode: 500})
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, &myError{err: fmt.Errorf("Body could not be read, %s", err),
			statusCode: 400})
		return
	}
	var message domain.Message
	err = json.Unmarshal(body, &message)
	if err != nil {
		handleError(w, &myError{err: fmt.Errorf("Body could not be json parsed: %s", err),
			statusCode: 400})
		return
	}
	if message.Body == "" {
		handleError(w, &myError{err: fmt.Errorf("Message.Body must be set"),
			statusCode: 400})
		return
	}
	post, err := h.postService.CreatePost(&message)
	if err != nil {
		handleError(w, &myError{err: fmt.Errorf("Could not store post"),
			statusCode: 500})
		return
	}
	// in order to have some fake stats for prometheues
	h.random.randomSleep(100, 510)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		handleError(w, &myError{err: fmt.Errorf("Could not json-encode created post"),
			statusCode: 500})
	}
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
