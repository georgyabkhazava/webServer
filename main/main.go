package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found"))
		return
	}
	var data Post
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if data.LenTitle() == false {
		http.Error(w, "Invalid Title", http.StatusBadRequest)
		return
	}
	if data.LenAuthor() == false {
		http.Error(w, "Invalid Author", http.StatusBadRequest)
		return
	}
	if data.LenText() == false {
		http.Error(w, "Invalid Text", http.StatusBadRequest)
		return
	}

	fmt.Println(data.Title, data.Author, data.Text)

	myInformation = append(myInformation, data)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
}
func firstPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if len(myInformation) == 0 {
		http.Error(w, "No posts available", http.StatusNotFound)
		return
	}
	firstPost := myInformation[0]
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"title": "%s", "author": "%s", "text": "%s"}`, firstPost.Title, firstPost.Author, firstPost.Text)
}

var myInformation = make([]Post, 0, 0)

type Post struct {
	Title  string
	Author string
	Text   string
}

func (p *Post) LenTitle() bool {
	if len(p.Title) >= 30 {
		return false
	}
	return true
}

func (p *Post) LenAuthor() bool {
	if len(p.Author) >= 20 {
		return false
	}
	return true
}

func (p *Post) LenText() bool {
	if len(p.Text) >= 500 {
		return false
	}
	return true
}

// сделать получение первого поста из слайса, проверять что это метод GET
func main() {
	// Регистрируем обработчик для всех запросов
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/post/create", postHandler)
	http.HandleFunc("/first-post", firstPostHandler)
	// Запускаем сервер на порту 8080
	fmt.Println("Starting server at port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
