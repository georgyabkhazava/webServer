package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

func savePost(title string, author string, text string) error {
	_, err := db.Exec(
		"INSERT INTO posts (title, author, text) VALUES ($1, $2, $3)",
		title,
		author,
		text,
	)
	return err
}

func getPosts() ([]Post, error) {
	rows, err := db.Query("SELECT id, title, author, text FROM posts")
	if err != nil {
		log.Fatal("Ошибка запроса:", err)
	}
	defer rows.Close()

	fmt.Println("\nСохраненные посты:")
	var posts = make([]Post, 0, 0)
	for rows.Next() {
		var id int
		var title, author, text string
		err = rows.Scan(&id, &title, &author, &text)
		if err != nil {
			log.Fatal("Ошибка чтения:", err)
		}
		post := Post{
			Title:  title,
			Author: author,
			Text:   text,
		}
		posts = append(posts, post)

		fmt.Printf("%d. %s: %s: %s\n", id, title, author, text)
	}
	return posts, nil
}

func getPostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	posts, err := getPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
	}
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

	//myInformation = append(myInformation, data)

	savePost(data.Title, data.Author, data.Text)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
}

func getFirstPost(db *sql.DB) (*Post, error) {
	query := `
        SELECT id, title, author, text 
        FROM posts 
        ORDER BY id ASC 
        LIMIT 1
    `

	var post Post
	err := db.QueryRow(query).Scan(
		&post.ID,
		&post.Title,
		&post.Author,
		&post.Text,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("посты не найдены")
		}
		return nil, fmt.Errorf("ошибка получения поста: %v", err)
	}

	return &post, nil
}

func firstPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Получаем первый пост
	post, err := getFirstPost(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Отправляем JSON
	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
	}
}

//var myInformation = make([]Post, 0, 0)

type Post struct {
	ID     int
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

var db *sql.DB

// сделать получение первого поста из слайса, проверять что это метод GET
// создать пару табличек в базе данных(1. Posts  с данными самого поста и id, 2. сделать табличку пользователей, написать в консоле код для добавления 5 постов, затем получить все посты определенного автора, получить первые 2 поста
// запрос для удаления поста, у которого id = 4, удалить таблицу books
// попробовать подключиться к базе данных и создать в ней запись
func main() {
	connStr := "user=postgres dbname=postgres password=1234 host=localhost port=5433 sslmode=disable"

	// Подключение к базе
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		log.Fatal("Ошибка подключения:", err)
	}
	defer db.Close()

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		log.Fatal("Не удалось проверить подключение:", err)
	}
	fmt.Println("Успешное подключение к базе!")

	// Регистрируем обработчик для всех запросов
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/post/create", postHandler)
	http.HandleFunc("/first-post", firstPostHandler)
	http.HandleFunc("/posts", getPostsHandler)

	// Запускаем сервер на порту 8080
	fmt.Println("Starting server at port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}

}
