package main

import (
	"log"
	"net/http"
	"time"

	middl "github.com/nmrshll/go-httpclient-middl"
	"github.com/nmrshll/go-httpclient-middl/middleware/logger"
	"github.com/nmrshll/go-httpclient-middl/middleware/statusvalidator"
)

func main() {
	httpClient := http.Client{Timeout: 30 * time.Second}
	client, err := middl.NewClient(&httpClient)
	if err != nil {
		log.Fatal(err)
	}

	// add middleware to you client (classic examples provided in this library or custom)
	client.UseMiddleware(logger.New())
	client.UseMiddleware(statusvalidator.New())

	// then do your requests as usual
	resp, err := client.Get("https://google.com")
	if err != nil {
		log.Fatal(err)
	}
	if resp == nil {
		log.Fatalf("no response from server")
	} // else
	defer resp.Body.Close()

	// do something with the response here
}
