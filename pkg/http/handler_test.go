package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gerlacdt/prom_example/pkg/domain"
	"github.com/gerlacdt/prom_example/pkg/inmemory"
)

func TestCreatePost(t *testing.T) {
	postService := inmemory.New()
	p := domain.Post{ID: 100, Body: "foobar42"}
	err := postService.CreatePost(&p)
	if err != nil {
		t.Fatalf("ERROR creating post, %s", err)
	}
	h := New(postService)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/v1/posts/100", nil)
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
	if !reflect.DeepEqual(p, rp) {
		t.Errorf("Got %v, expected: %v", rp, p)
	}
}
