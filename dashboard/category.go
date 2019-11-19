package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func Categories(w http.ResponseWriter, r *http.Request) {
	Categories := []Category{}
	C := Category{}
	rows, err := database.Query("SELECT id, category, date FROM category")

	if err != nil {
		log.Println("Gabim me databazen!")
	} else {
		for rows.Next() {
			rows.Scan(&C.Id, &C.Category, &C.Date)
			Categories = append(Categories, C)
		}
		fmt.Println(Categories)

		var templates = template.Must(template.ParseFiles("views/categories.gohtml", "views/base.gohtml"))
		varmap := map[string]interface{}{
			"Auth":       IsAuth(r),
			"Admin":      IsAdmin(r),
			"Categories": Categories,
		}

		err = templates.ExecuteTemplate(w, "base", varmap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func NewCategory(w http.ResponseWriter, r *http.Request) {

	var templates = template.Must(template.ParseFiles("views/edit_category.gohtml", "views/base.gohtml"))
	varmap := map[string]interface{}{
		"Auth":      IsAuth(r),
		"Admin":     IsAdmin(r),
		"Categorie": Categories,
	}

	err := templates.ExecuteTemplate(w, "base", varmap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func EditCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["id"]
	Category := Category{}

	err := database.QueryRow("SELECT id, category, date FROM category WHERE id=?", Id).Scan(&Category.Id, &Category.Category, &Category.Date)
	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("Couldn't get page!")
	} else {
		var templates = template.Must(template.ParseFiles("views/edit_category.gohtml", "views/base.gohtml"))
		varmap := map[string]interface{}{
			"Auth":     IsAuth(r),
			"Admin":    IsAdmin(r),
			"Category": Category,
		}

		err := templates.ExecuteTemplate(w, "base", varmap)
		checkErr(err)
	}
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var Id string
	Id = r.FormValue("id")
	category := r.FormValue("category")

	t := time.Now()
	date := t.Format("2006-01-02 15:04:05")
	id, _ := strconv.Atoi(Id)

	if id == 0 {
		stmt, err := database.Prepare("INSERT INTO category SET category=?, date=?")
		checkErr(err)
		_, err = stmt.Exec(category, date)
		checkErr(err)
		rp = "Category inserted successfully."
	} else {
		stmt, err := database.Prepare("UPDATE category SET category=?, date=? WHERE id=?")
		checkErr(err)
		_, err = stmt.Exec(category, date, id)
		checkErr(err)
		rp = "Category updated successfully."
	}

	var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/base.gohtml"))
	varmap := map[string]interface{}{
		"Auth":   IsAuth(r),
		"Admin":  IsAdmin(r),
		"Report": rp,
	}
	_ = templates.ExecuteTemplate(w, "base", varmap)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	Id := mux.Vars(r)["id"]

	stmt, err := database.Prepare("DELETE FROM category WHERE id=?")
	checkErr(err)

	_, err = stmt.Exec(Id)
	checkErr(err)
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}
