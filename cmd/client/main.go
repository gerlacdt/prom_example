package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gerlacdt/prom_example/pkg/domain"
)

func main() {
	endpoint := flag.String("endpoint", "http://localhost:8080", "the complete endpoint url for the service to call")

	flag.Parse()

	go workerPOST(*endpoint)
	go workerGET(*endpoint)
	go workerHealth(*endpoint)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	s := <-c
	fmt.Println("Finished, got signal:", s)
}

func workerPOST(endpoint string) {
	i := 0
	for {
		time.Sleep(300 * time.Millisecond)
		i++
		start := time.Now()
		msg := domain.Message{Body: fmt.Sprintf("foobar-%d", i)}
		_, err := createPost(endpoint, &msg)
		if err != nil {
			fmt.Printf("POST /v1/posts failed: %s duration: %v\n", err, time.Since(start))
			continue
		}
		fmt.Printf("POST /v1/posts duration: %v\n", time.Since(start))
	}
}

func workerGET(endpoint string) {
	for {
		time.Sleep(500 * time.Millisecond)
		start := time.Now()
		_, err := getPost(endpoint, 1)
		if err != nil {
			fmt.Printf("GET /v1/posts/:id failed: %s duration: %v\n", err, time.Since(start))
			continue
		}
		fmt.Printf("GET /v1/posts/:id duration: %v\n", time.Since(start))
	}
}

func workerHealth(endpoint string) {
	for {
		time.Sleep(1 * time.Second)
		start := time.Now()
		err := getHealth(endpoint)
		if err != nil {
			fmt.Printf("GET /health failed: %v duration: %v\n", err, time.Since(start))
			continue
		}
		fmt.Printf("GET /health duration: %v\n", time.Since(start))
	}
}

func handleHTTPError(method string, statusCode int, body []byte) error {
	var httpError domain.HTTPError
	err := json.Unmarshal(body, &httpError)
	if err != nil {
		return fmt.Errorf("Failed to parse http error: %s, method: %s, statusCode: %d", err, method, statusCode)
	}
	return fmt.Errorf("HTTP call failed, method: %s,statusCode: %d, error: %s",
		method, statusCode, httpError.Message)
}

func createPost(endpoint string, msg *domain.Message) (*domain.Post, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(msg)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/posts", endpoint), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("ERROR closing resp.Body, %v", err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 {
		return nil, handleHTTPError(req.Method, resp.StatusCode, body)
	}
	var p domain.Post
	err = json.Unmarshal(body, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func getPost(endpoint string, id int) (*domain.Post, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/posts/%d", endpoint, id), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("ERROR closing resp.Body, %v", err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, handleHTTPError(req.Method, resp.StatusCode, body)
	}
	var p domain.Post
	err = json.Unmarshal(body, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func getHealth(endpoint string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/health", endpoint), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("ERROR closing resp.Body, %v", err)
		}
	}()
	if resp.StatusCode != 200 {
		return fmt.Errorf("GET /health failed, statusCode: %d", resp.StatusCode)
	}
	return nil
}
