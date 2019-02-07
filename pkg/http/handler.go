package http

import (
	"fmt"
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
	if r.Method == "GET" {
		id := strings.TrimPrefix(r.URL.Path, "/v1/posts/")
		id2, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		post, err := h.postService.Post(id2)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		_, err = fmt.Fprintf(w, "%s", post.Body)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	} else if r.Method == "POST" {
		// TODO create post from body, return
		h.postService.CreatePost(post)
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
