package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

// var db *sql.DB

type Article struct {
	Id      int
	Title   string
	Preview string
	Content string
}

type Log struct {
	Login        string
	PasswordHash string
}

func mainPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	Data := showDB(db)
	templ := template.Must(template.ParseFiles("html/homepage.html"))
	templ.Execute(w, Data)
}

func loginPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	templ := template.Must(template.ParseFiles("html/login.html"))
	templ.Execute(w, nil)
	CheckPass("lmk", "6656")
}

func HashPassword(password, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password + salt))
	return hex.EncodeToString(hasher.Sum(nil))
}

func CheckPass(login string, pass string) {
	data, err := ioutil.ReadFile("pass/pass.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	lines := strings.Split(string(data), "\n")
	log := Log{}
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {
			case "Login":
				log.Login = value
			case "PasswordHash":
				log.PasswordHash = value
			}
		}
	}

	fmt.Printf("Login: %s\nPassword Hash: %s\n", log.Login, log.PasswordHash)
}

func initQuery(db *sql.DB) {
	var title string
	var preview string
	var content string
	fmt.Scan(&title, &preview, content)
	article := Article{
		Title:   title,
		Preview: preview,
		Content: content,
	}
	pk := insertArtice(db, article)
	fmt.Println(pk)
}

func showDB(db *sql.DB) []Article {
	var id int
	var title string
	var preview string
	var content string
	data := []Article{}
	rows, err := db.Query("SELECT id, title, preview, content FROM article")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&id, &title, &preview, &content)
		data = append(data, Article{
			Title:   title,
			Preview: preview,
			Content: content,
		})
	}
	return data
}

func connectToDB() *sql.DB {

	connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return nil
	}
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		return nil
	}
	createArticleTable(db)
	return db
}

func insertArtice(db *sql.DB, article Article) int {
	query := `INSERT INTO article (title, preview, content)
    VALUES ($1, $2, $3) RETURNING id`
	var pk int
	err := db.QueryRow(query, article.Title, article.Preview, article.Content).Scan(&pk) //returns at most one value
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return pk
}

func createArticleTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS article (
		id SERIAL PRIMARY KEY, 
		title VARCHAR(50) NOT NULL,
		preview TEXT NOT NULL, 
		content TEXT NOT NULL,
		created timestamp DEFAULT NOW()
	);`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {

	db := connectToDB()
	initQuery(db)
	Data := showDB(db)
	fmt.Println(Data)
	htmlDir := http.FileServer(http.Dir("html"))
	http.Handle("/html/", http.StripPrefix("/html", htmlDir))

	cssDir := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css", cssDir))

	picDir := http.FileServer(http.Dir("images"))
	http.Handle("/images/", http.StripPrefix("/images", picDir))

	jsDir := http.FileServer(http.Dir("js"))
	http.Handle("/js/", http.StripPrefix("/js", jsDir))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mainPage(w, r, db)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginPage(w, r, db)
	})

	fmt.Println("Server is listening on port 8888")
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
}
