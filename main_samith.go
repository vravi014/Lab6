package main

import (
	"fmt"
	"net/http"
	"sync"
)

type FetchResult struct {
	URL        string
	StatusCode int
	Size       int64
	Error      error
}

func worker(id int, jobs <-chan string, results chan<- FetchResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range jobs {
		resp, err := http.Get(url)
		if err != nil {
			results <- FetchResult{URL: url, Error: err}
			continue
		}
		defer resp.Body.Close()
		results <- FetchResult{
			URL:        url,
			StatusCode: resp.StatusCode,
			Size:       resp.ContentLength,
			Error:      nil,
		}
	}
}

func main() {z
	urls := []string{
		"https://www.example.com",
		"https://www.golang.org",
		"https://www.uottawa.ca",
		"https://www.github.com",
		"https://www.httpbin.org/get",
	}

	const numWorkers = 3
	jobs := make(chan string, len(urls))
	results := make(chan FetchResult, len(urls))

	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	wg.Wait()
	close(results)

	for result := range results {
		if result.Error != nil {
			fmt.Printf("Error fetching %s: %v\n", result.URL, result.Error)
		} else {
			fmt.Printf("%s | Status Code: %d | Size: %d bytes\n",
				result.URL, result.StatusCode, result.Size)
		}
	}
}