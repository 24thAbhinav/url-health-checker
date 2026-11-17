package main

import (
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	URL        string
	Healthy    bool
	StatusCode string
	Latency    time.Duration
	Err        error
}

func checkStatus(id int, jobs <-chan string, results chan<- Response, client *http.Client) {
	for url := range jobs {
		fmt.Println("Worker", id, "checking", url)

		start := time.Now()
		resp, err := client.Get(url)
		duration := time.Since(start)

		if err != nil {
			results <- Response{
				URL:     url,
				Healthy: false,
				Latency: duration,
				Err:     err,
			}
			continue
		}

		resp.Body.Close()

		healthy := resp.StatusCode >= 200 && resp.StatusCode < 300

		results <- Response{
			URL:        url,
			Healthy:    healthy,
			StatusCode: resp.Status,
			Latency:    duration,
		}
	}
}

func main() {
	urls := []string{
		"https://www.google.com",
		"https://github.com",
		"https://go.dev",
		"https://example.com",
		"https://invalid.invalid",
	}

	jobs := make(chan string, len(urls))
	results := make(chan Response, len(urls))

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for i := 1; i <= 3; i++ {
		go checkStatus(i, jobs, results, client)
	}

	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	for range urls {
		response := <-results

		fmt.Println("-----------------------------------")
		fmt.Printf("URL:         %s\n", response.URL)

		if response.Err != nil {
			fmt.Printf("Health:      UNHEALTHY\n")
			fmt.Printf("Error:       %v\n", response.Err)
		} else {
			status := "UNHEALTHY"
			if response.Healthy {
				status = "HEALTHY"
			}

			fmt.Printf("Health:      %s\n", status)
			fmt.Printf("HTTP Status: %s\n", response.StatusCode)
		}

		fmt.Printf("Latency:     %.2f ms\n",
			float64(response.Latency)/float64(time.Millisecond))
		fmt.Println("-----------------------------------")
	}
}
