package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	DBHost  = "127.0.0.1"
	DBPort  = ":3306"
	DBUser  = "root"
	DBPass  = ""
	DBDbase = "blog"
)

type Post struct {
	ID            string `json:"ID"`
	Title         string `json:"Title"`
	Description   string `json:"Description"`
	Category_id   string `json:"category_id"`
	Date          string `json:"date"`
	Author_id     string `json:"author_id"`
	Id_category   string `json:"id_category"`
	Category      string `json:"category"`
	Date_category string `json:"date_category"`
}

type allPosts []Post

// var posts = allPosts{}

var database *sql.DB

func main() {
	dbConn := fmt.Sprintf("%s:%s@tcp(%s)/%s", DBUser, DBPass, DBHost, DBDbase)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		panic("Couldn't connect!")
	}

	database = db

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/posts", createPost).Methods("POST")
	router.HandleFunc("/posts", getAllPosts).Methods("GET")
	router.HandleFunc("/post/{id:[0-9]+}", getOnePost).Methods("GET")
	router.HandleFunc("/posts/{id:[0-9]+}", updatePost).Methods("PATCH")
	router.HandleFunc("/posts/{id:[0-9]+}", deletePost).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8081", router))
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var newPost Post
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newPost)
	stmt, err := database.Prepare("INSERT INTO posts SET title = ?, description = ?, category_id=?, date=?, author_id=?")
	checkErr(err)
	_, err = stmt.Exec(newPost.Title, newPost.Description, newPost.Category_id, newPost.Date, newPost.Author_id)
	checkErr(err)

	w.WriteHeader(http.StatusCreated)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newPost)

}

func getOnePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]

	P := Post{}

	// err := database.QueryRow("SELECT id, title, description, category_id, date, author_id FROM posts WHERE id=?", postID).Scan(&P.ID, &P.Title, &P.Description, &P.Category_id, &P.Date, &P.Author_id)
	err := database.QueryRow("SELECT * FROM posts INNER JOIN category ON category_id=category.id WHERE posts.id=?", postID).Scan(&P.ID, &P.Title, &P.Description, &P.Category_id, &P.Date, &P.Author_id, &P.Id_category, &P.Category, &P.Date_category)

	checkErr(err)

	b, err := json.Marshal(P)
	if err != nil {
		log.Println(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(b))
	}
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {
	posts := []Post{}
	P := Post{}
	// rows, err := database.Query("SELECT id, title, description, category_id, date, author_id FROM posts")
	rows, err := database.Query("SELECT * FROM posts INNER JOIN category ON category_id=category.id ORDER BY posts.id DESC")

	checkErr(err)

	for rows.Next() {
		rows.Scan(&P.ID, &P.Title, &P.Description, &P.Category_id, &P.Date, &P.Author_id, &P.Id_category, &P.Category, &P.Date_category)
		posts = append(posts, P)
	}

	b, err := json.Marshal(posts)
	if err != nil {
		log.Println(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(b))
	}
}

func updatePost(w http.ResponseWriter, r *http.Request) {

	postID := mux.Vars(r)["id"]
	P := Post{}

	reqBody, err := ioutil.ReadAll(r.Body)
	checkErr(err)

	json.Unmarshal(reqBody, &P)

	stmt, err := database.Prepare("UPDATE posts SET title = ?, description = ?, category_id=?, date=?, author_id=? WHERE id=?")
	checkErr(err)
	_, err = stmt.Exec(P.Title, P.Description, P.Category_id, P.Date, P.Author_id, postID)
	checkErr(err)

}

func deletePost(w http.ResponseWriter, r *http.Request) {

	eventID := mux.Vars(r)["id"]

	stmt, err := database.Prepare("DELETE FROM posts WHERE id=?")
	checkErr(err)

	_, err = stmt.Exec(eventID)
	checkErr(err)
}

func checkErr(err error) {

	if err != nil {
		panic(err)
	}
}
