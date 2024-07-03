// About Locks: https://chatgpt.com/share/598d4cac-69bb-4e1b-873b-ea90ac9ddf60
package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/victoronofiok/gourlshortener/routes"
)
type URLStore struct {
	urls map[string]string
	mu sync.RWMutex
}

func (s *URLStore) Get(key string) string {
	// lock so no update (write) can happen on the store
	s.mu.RLock()
	defer s.mu.Unlock()

	url := s.urls[key]
	return url
}

func (s *URLStore) Set(key, url string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, present := s.urls[key]
	if !present {
		return false
	}
	s.urls[key] = url
	return true
}

// generate the key for the long url
func (s *URLStore) Put(url string) string {
	for {
		key := genKey(s.Count())
		if s.Set(key, url) {
			return key
		}
	}
}

func (s *URLStore) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.urls)
}

func genKey(n int) string {
	return ""
}


// Making the value of a struct is done in Go by defining a function with the prefix New, that returns an initialized value of the type
func NewURLStore() *URLStore {
	return &URLStore{ urls: make(map[string]string) }
}

func main() {
	// var store = NewURLStore()

	http.HandleFunc("/add", routes.Add)
	http.HandleFunc("/redirect", routes.Redirect)
	log.Fatal(http.ListenAndServe(":3000", nil))
}