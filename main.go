package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	// Define the URL you want to send the GET request to.
	url := "https://example.com"

	// Send an HTTP GET request.
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	// Create a buffer to store the response body.
	var responseBody io.Writer

	// Read the response body using io.Copy.
	_, err = io.Copy(responseBody, response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

}
