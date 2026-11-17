# url-health-checker

A small Go program for checking the health of a list of URLs concurrently. This is a learning project focused on Go and concurrency patterns.

## What it does

- Sends an HTTP GET request to each URL in a list
- Marks a URL healthy if it returns a 2xx status code
- Reports status code, latency, and any request errors
- Runs checks in parallel using a worker pool

## Concurrency concepts used

- Goroutines for workers
- Buffered channels for jobs and results
- Worker pool pattern (3 workers)
- Closing a channel to signal no more work
- HTTP client timeout per request

## Run

```sh
go run main.go
```

## Notes

- URLs are defined in the `urls` slice in `main.go`
- The worker count is set in the loop that starts goroutines
- Each request has a 5 second timeout
