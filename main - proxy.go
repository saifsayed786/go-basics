package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func main() {
	// Define the target URL of the backend server
	targetURL := "http://localhost:8000"

	// Parse the target URL
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatal("Error parsing target URL:", err)
	}

	// Create a reverse proxy with a custom Director function
	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = target.Scheme
			r.URL.Host = target.Host
			r.Host = target.Host
			// r.URL.Path = "/addAlbums" + r.URL.Path

			// Modify the request here as needed
			modifyRequest(r)
		},
		ModifyResponse: func(r *http.Response) error {
			// Modify the response here as needed
			modifyResponse(r)
			return nil
		},
	}

	// Start the server
	log.Println("Reverse Proxy Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", proxy))
}

func modifyRequest(r *http.Request) {

}

func modifyResponse(r *http.Response) error {
	originalBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return err
	}
	defer r.Body.Close()

	// Modify the response body by appending a custom message
	modifiedBody := append(originalBody, []byte("\nModified by Reverse Proxy")...)

	// Set the modified body in the response
	r.Body = ioutil.NopCloser(bytes.NewReader(modifiedBody))
	r.ContentLength = int64(len(modifiedBody))

	// Update the Content-Length header
	r.Header.Set("Content-Length", strconv.Itoa(len(modifiedBody)))

	return nil
}
