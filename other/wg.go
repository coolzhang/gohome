package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var urls = []string{
		"http://www.baidu.com/",
		"http://www.bing.com/",
		"http://www.sogou.com/",
	}
	for _, url := range urls {
		// Increment the WaitGroup counter.
		wg.Add(1)
		// Launch a goroutine to fetch the URL.
		go func(url string) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			// Fetch the URL.
			fmt.Println(url)
			http.Get(url)
		}(url)
	}
	// Wait for all HTTP fetches to complete.
	wg.Wait()
}
