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
func worker(id int, jobs <-chan string, results chan<- FetchResult) {
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
	for i := 1; i <= numWorkers ; i++ {
		go worker(i, jobs, results)
	}

	// TODO: Send jobs 
	for n := 0; n <= len(urls)-1; n++ {
		jobs <- urls[n]
	}
	close (jobs)
	// TODO: Collect results
	for result := range results{
		fmt.Println(result)
	} 
	go func(){
		wg.Wait()
		close(results)	
	}()
	fmt.Println("\nScraping complete!")	
}