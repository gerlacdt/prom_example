package main

import (
	"fmt"
	"log"

	"github.com/gerlacdt/prom_example/pkg/http"
	"github.com/gerlacdt/prom_example/pkg/inmemory"
)

func main() {
	fmt.Println("Start server...")

	postService := &inmemory.PostService{}
	h := http.New(postService)
	http.Handle("/v1/posts", h)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
