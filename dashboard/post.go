package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	Posts := []Post{}
	P := Post{}

	rows, err := database.Query("SELECT * FROM posts INNER JOIN category ON category_id=category.id")

	if err != nil {
		log.Println("Gabim me databazen!")
	} else {
		for rows.Next() {
			rows.Scan(&P.Id, &P.Title, &P.Description, &P.Category_id, &P.Date, &P.Author, &P.Id_category, &P.Category, &P.Date_category)
			Posts = append(Posts, P)
		}

		var templates = template.Must(template.ParseFiles("views/index.gohtml", "views/base.gohtml"))
		varmap := map[string]interface{}{
			"Auth":  IsAuth(r),
			"Admin": IsAdmin(r),
			"Posts": Posts,
		}

		err = templates.ExecuteTemplate(w, "base", varmap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func NewPost(w http.ResponseWriter, r *http.Request) {
	Categories := []Category{}
	C := Category{}

	rows, err := database.Query("SELECT id, category, date FROM Category")
	checkErr(err)
	for rows.Next() {
		rows.Scan(&C.Id, &C.Category, &C.Date)
		Categories = append(Categories, C)
	}

	var templates = template.Must(template.ParseFiles("views/edit.gohtml", "views/base.gohtml"))
	varmap := map[string]interface{}{
		"Auth":      IsAuth(r),
		"Admin":     IsAdmin(r),
		"Categorie": Categories,
	}

	err = templates.ExecuteTemplate(w, "base", varmap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageID := vars["id"]
	P := Post{}
	Categories := []Category{}
	C := Category{}

	rows, err := database.Query("SELECT id, category, date FROM Category")
	checkErr(err)
	for rows.Next() {
		rows.Scan(&C.Id, &C.Category, &C.Date)
		Categories = append(Categories, C)
	}

	err = database.QueryRow("SELECT id, title, description, category_id, date, author_id FROM posts WHERE id=?", pageID).Scan(&P.Id, &P.Title, &P.Description, &P.Category_id, &P.Date, &P.Author)
	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("Couldn't get page!")
	} else {
		var templates = template.Must(template.ParseFiles("views/edit.gohtml", "views/base.gohtml"))
		varmap := map[string]interface{}{
			"Auth":      IsAuth(r),
			"Admin":     IsAdmin(r),
			"Post":      P,
			"Categorie": Categories,
		}

		err := templates.ExecuteTemplate(w, "base", varmap)
		checkErr(err)
	}
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "gosession")
	var post_id string
	post_id = r.FormValue("id")

	title := r.FormValue("title")
	description := r.FormValue("description")
	category_id := r.FormValue("category")

	t := time.Now()
	date := t.Format("2006-01-02 15:04:05")
	id, _ := strconv.Atoi(post_id)

	if id == 0 {
		stmt, err := database.Prepare("INSERT INTO posts SET title=?, description=?, category_id=?, date=?, author_id=?")
		checkErr(err)
		_, err = stmt.Exec(title, description, category_id, date, session.Values["user_id"])
		checkErr(err)
		rp = "Record inserted successfully."
	} else {
		stmt, err := database.Prepare("UPDATE posts SET title=?, description=?, category_id=?, date=? WHERE id=?")
		checkErr(err)
		_, err = stmt.Exec(title, description, category_id, date, id)
		checkErr(err)
		rp = "Record updated successfully."
	}

	if uploadFile(w, r, post_id) {
		rp = rp + " Image uploaded successfully."
	} else {
		rp = rp + " Image upload failed."
	}

	var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/base.gohtml"))
	varmap := map[string]interface{}{
		"Auth":   IsAuth(r),
		"Admin":  IsAdmin(r),
		"Report": rp,
	}
	_ = templates.ExecuteTemplate(w, "base", varmap)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	Id := mux.Vars(r)["id"]

	stmt, err := database.Prepare("DELETE FROM posts WHERE id=?")
	checkErr(err)

	_, err = stmt.Exec(Id)
	checkErr(err)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func uploadFile(w http.ResponseWriter, r *http.Request, post_id string) bool {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("imageFile")
	if err != nil {
		fmt.Println("Gabim gjate leximit te fajllit")
		fmt.Println(err)
		return false
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
	header := handler.Header
	var ext string

	switch header["Content-Type"][0] {
	case "image/jpeg":
		ext = "jpg"
	case "image/png":
		ext = "png"
	default:
		return false
	}

	tempFile, err := ioutil.TempFile("uploads", "upload-*."+ext)
	if err != nil {
		fmt.Println(err)
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)
	tempFile.Close()
	err = os.Rename(tempFile.Name(), "assets/images/image-"+post_id+".jpg")

	if err != nil {
		log.Fatal(err)
	}

	return true
}

func ShowPosts(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	fmt.Println(url)

	Posts := []Post{}

	faqja, err := http.Get("http://localhost:8081/posts")
	if err != nil {
		fmt.Println("Nuk munda ta hap faqen")
	} else {
		defer faqja.Body.Close()

		body, _ := ioutil.ReadAll(faqja.Body)
		json.Unmarshal(body, &Posts)

		var templates = template.Must(template.ParseFiles("views/blogen.gohtml", "views/base.gohtml"))
		varmap := map[string]interface{}{
			"Auth":  IsAuth(r),
			"Admin": IsAdmin(r),
			"Posts": Posts,
			"Url":   url,
		}
		_ = templates.ExecuteTemplate(w, "base", varmap)
	}

}

func ShowPost(w http.ResponseWriter, r *http.Request) {
	Post := Post{}
	Comments := []Comment{}
	C := Comment{}
	vars := mux.Vars(r)
	post_id := vars["id"]
	url := "http://localhost:8081/post/" + post_id

	faqja, err := http.Get(url)
	if err != nil {
		fmt.Println("Nuk munda ta hap faqen")
	} else {
		defer faqja.Body.Close()

		body, err := ioutil.ReadAll(faqja.Body)
		json.Unmarshal(body, &Post)

		rows, err := database.Query("SELECT * FROM comments INNER JOIN posts ON post_related=posts.id INNER JOIN users ON comment_author=users.id WHERE post_related=?", post_id)
		if err != nil {
			log.Println("Gabim me databazen!")
		}

		for rows.Next() {
			rows.Scan(&C.Id, &C.Comment, &C.Date, &C.Post_related, &C.Comment_author, &C.Post_related_id, &C.Title, &C.Description, &C.Category_id, &C.Post_date, &C.Author_id, &C.Author_c_id, &C.Author_name, &C.Author_surname, &C.Author_email, &C.Author_password, &C.Author_role)
			Comments = append(Comments, C)
		}

		fmt.Println(Comments)

		var templates = template.Must(template.ParseFiles("views/details.gohtml", "views/base.gohtml"))
		varmap := map[string]interface{}{
			"Auth":     IsAuth(r),
			"Admin":    IsAdmin(r),
			"Post":     Post,
			"Comments": Comments,
		}
		err = templates.ExecuteTemplate(w, "base", varmap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

}
