package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func main() {
	for i := 0; i < 100000; i++ {
		fmt.Println("sending request to server")
		payload := []byte("request payload")
		resp, err := http.Post("http://localhost:8080/abc/123", "application/octet-stream", bytes.NewBuffer(payload))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		fmt.Printf("response status code: %d\n", resp.StatusCode)			
	}
}