package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)
var wg sync.WaitGroup

//Structure to store results
type FetchResult struct {
	URL string
	StatusCode int
	Size int
	Error error
}

//Worker function
func worker(id int, jobs <-chan string, results chan<- FetchResult, wg *sync.WaitGroup) {
	// TODO: fetch the url
	// TODO: send Result struct to results channel
	// hint: use resp, err := http.Get(url)

	defer wg.Done()
	for url := range jobs {
		resp, err := http.Get(url)
		status := resp.StatusCode
		body, err := ioutil.ReadAll(resp.Body)
		size := len(body)
		results <- FetchResult{URL: url, StatusCode: status, Size: size, Error: err }
	}
}

func main() {
	urls := []string{
		"https://example.com",
		"https://golang.org",
		"https://uottawa.ca",
		"https://github.com",
		"https://httpbin.org/get",
	}

	numWorkers := 3

	jobs := make(chan string, len(urls))
	results := make(chan FetchResult, len(urls))

	// TODO: Start workers
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers ; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	// TODO: Send jobs 
	for n := 0; n <= len(urls)-1; n++ {
		jobs <- urls[n]
	}
	
	close (jobs)
	wg.Wait()
	close(results)
	// TODO: Collect results
	for result := range results {
		if result.Error != nil {
			fmt.Printf("Error fetching %s: %v\n", result.URL, result.Error)
		} else {
			fmt.Printf("%s | Status Code: %d | Size: %d bytes\n",
				result.URL, result.StatusCode, result.Size)
		}
	}
}