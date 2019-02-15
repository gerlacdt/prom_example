package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gerlacdt/prom_example/pkg/domain"
	"github.com/gerlacdt/prom_example/pkg/inmemory"
)

func TestGetPost(t *testing.T) {
	postService := inmemory.New()
	m := domain.Message{Body: "foobar42"}
	p, err := postService.CreatePost(&m)
	if err != nil {
		t.Fatalf("ERROR creating post, %s", err)
	}
	h := New(postService)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", fmt.Sprintf("/v1/posts/%d", p.ID), nil)
	h.ServeHTTP(w, r)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Fatalf("ERROR closing resp.Body, %v", err)
		}
	}()
	if err != nil {
		t.Fatalf("ERROR read http response body, %s", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode, got: %d, expected: %d", resp.StatusCode, 200)
	}
	var rp domain.Post
	err = json.Unmarshal(body, &rp)
	if err != nil {
		t.Fatalf("ERROR json unmarshalling: %s", err)
	}
	if !reflect.DeepEqual(*p, rp) {
		t.Errorf("Got %v, expected: %v", rp, p)
	}
}

func TestCreatePost(t *testing.T) {
	postService := inmemory.New()
	m := domain.Message{Body: "foobar42"}
	h := New(postService)
	w := httptest.NewRecorder()

	// create http request
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(m)
	if err != nil {
		t.Fatalf("ERROR in json body encoding, %v", err)
	}
	r, _ := http.NewRequest("POST", "/v1/posts", buf)
	r.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(w, r)
	resp := w.Result()
	if resp.StatusCode != 201 {
		t.Errorf("StatusCode, got: %d, expected: %d", resp.StatusCode, 201)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Fatalf("ERROR closing resp.Body, %v", err)
		}
	}()
	var p domain.Post
	err = json.Unmarshal(body, &p)
	if err != nil {
		t.Fatalf("ERROR json unmarshalling: %s", err)
	}
	rp, err := postService.Post(p.ID)
	if err != nil {
		t.Errorf("New post was not stored.")
	}
	if !reflect.DeepEqual(&p, rp) {
		t.Errorf("Got %v, expected: %v", rp, p)
	}
}

func TestDeletePost(t *testing.T) {
	postService := inmemory.New()
	m := domain.Message{Body: "foobar42"}
	p, err := postService.CreatePost(&m)
	if err != nil {
		t.Fatalf("ERROR creating a post, %v", err)
	}
	h := New(postService)
	// create http request
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/posts/%d", p.ID), nil)
	h.ServeHTTP(w, r)
	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode, got: %d, expected: %d", resp.StatusCode, 200)
	}
	_, err = postService.Post(p.ID)
	if err == nil {
		t.Errorf("Post should not be deleted.")
	}
}
