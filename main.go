package main

import (
	"net/http"
	"time"
	"fmt"
	"math/rand"
)

const (
	maxLength  = 8
	base62Char = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type URLStore struct {
	urls     map[string]urlData
	queue    []string
	capacity int
}

type urlData struct {
	url       string
	createdAt time.Time
}

func main() {
	store := URLStore{
		urls:     make(map[string]urlData),
		capacity: 20000,
	}
	http.HandleFunc("/", store.shortenURL)
	http.HandleFunc("/shortened-url/", store.redirectURL)
	http.ListenAndServe(":8080", nil)
}

func (s *URLStore) shortenURL(w http.ResponseWriter, request *http.Request) {
	fmt.Println("Shortening the URL")

	originalURL := request.URL.Query().Get("url")
	if originalURL == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if len(s.queue) >= s.capacity {
		http.Error(w, "Capacity Expired", http.StatusBadRequest)
		return
	}

	code := s.generateUniqueCode()
	s.urls[code] = urlData{url: originalURL, createdAt: time.Now()}
	s.queue = append(s.queue, code)

	shortenedURL := fmt.Sprintf("http://localhost:8080/shortened-url/%s", code)
	w.Write([]byte(shortenedURL))
}

func (s *URLStore) redirectURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Checking for long URL")

	code := r.URL.Path[len("/shortened-url/"):]

	if data, ok := s.urls[code]; ok {

		fmt.Println(data.url, "debug")

		if time.Now().Sub(data.createdAt) > 24*time.Hour {
			delete(s.urls, code)
			http.Error(w, "URL has expired", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, data.url, http.StatusFound)
	} else {
		http.NotFound(w, r)
	}
}

func (s *URLStore) generateUniqueCode() string {
	code := ""

	for {
		for i := 0; i < maxLength; i++ {
			code += string(base62Char[rand.Intn(len(base62Char))])
		}

		_, ok := s.urls[code]

		if !ok {
			break;
		}
	}

	return code
}