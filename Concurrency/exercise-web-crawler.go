package main

import (
	"fmt"
	"sync"
)

type SafeCache struct {
	v   map[string]string
	mux sync.Mutex
}

// Add puts new key value pair in.
func (c *SafeCache) Add(key string, value string) {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key] = value
	c.mux.Unlock()
}

// Value returns the current value of the counter for the given key.
func (c *SafeCache) GetValue(key string) string {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	value := c.v[key]
	defer c.mux.Unlock()
	return value
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
// messy but easiest way I found to close/manage channels was create one
// per iteration that collects values from all children and closes/returns
// at end of function
func Crawl(url string, depth int, fetcher Fetcher, ch chan string, sc SafeCache) {
	defer close(ch)

	if depth <= 0 {
		return
	}

	cacheval := sc.GetValue(url)

	if cacheval != "" {
		// Cache hit
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		// dont need to have separate cache entries for err and body but may as well
		sc.Add(url, err.Error())
		ch <- err.Error()
		return
	}

	sc.Add(url, body)
	ch <- fmt.Sprintf("found: %s %q", url, body)

	result := make([]chan string, len(urls))
	for i, u := range urls {
		result[i] = make(chan string)
		go Crawl(u, depth-1, fetcher, result[i], sc)
	}

	for i := range result {
		for s := range result[i] {
			ch <- s
		}
	}

	return
}

func main() {
	safec := SafeCache{v: make(map[string]string)}

	ch := make(chan string)
	go Crawl("http://golang.org/", 4, fetcher, ch, safec)
	for i := range ch {
		fmt.Println(i)
	}
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
