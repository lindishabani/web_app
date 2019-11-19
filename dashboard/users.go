package main

import (
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

func Users(w http.ResponseWriter, r *http.Request) {
	Users := []User{}
	U := User{}
	rows, err := database.Query("SELECT id, name, surname, email, role FROM users")

	if err != nil {
		log.Println("Gabim me databazen!")
	} else {
		for rows.Next() {
			rows.Scan(&U.Id, &U.Name, &U.Surname, &U.Email, &U.Role)
			Users = append(Users, U)
		}
		fmt.Println(Users)

		var templates = template.Must(template.ParseFiles("views/users.gohtml", "views/base.gohtml"))
		varmap := map[string]interface{}{
			"Auth":  IsAuth(r),
			"Admin": IsAdmin(r),
			"Users": Users,
		}

		err = templates.ExecuteTemplate(w, "base", varmap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["id"]
	User := User{}

	err := database.QueryRow("SELECT id, name, surname, email, role FROM users WHERE id=?", Id).Scan(&User.Id, &User.Name, &User.Surname, &User.Email, &User.Role)
	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("Couldn't get page!")
	} else {
		var templates = template.Must(template.ParseFiles("views/edit_user.gohtml", "views/base.gohtml"))
		varmap := map[string]interface{}{
			"Auth":  IsAuth(r),
			"Admin": IsAdmin(r),
			"User":  User,
		}

		err := templates.ExecuteTemplate(w, "base", varmap)
		checkErr(err)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var Id string
	Id = r.FormValue("id")
	name := r.FormValue("name")
	surname := r.FormValue("surname")
	email := r.FormValue("email")
	role := r.FormValue("role")

	id, _ := strconv.Atoi(Id)

	if id == 0 {
		stmt, err := database.Prepare("INSERT INTO users SET name=?, surname=?, email=?, role=?")
		checkErr(err)
		_, err = stmt.Exec(name, surname, email, role)
		checkErr(err)
		rp = "User inserted successfully."
	} else {
		stmt, err := database.Prepare("UPDATE users SET name=?, surname=?, email=?, role=? WHERE id=?")
		checkErr(err)
		_, err = stmt.Exec(name, surname, email, role, id)
		checkErr(err)
		rp = "User updated successfully."
	}

	var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/base.gohtml"))
	varmap := map[string]interface{}{
		"Auth":   IsAuth(r),
		"Admin":  IsAdmin(r),
		"Report": rp,
	}
	_ = templates.ExecuteTemplate(w, "base", varmap)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	Id := mux.Vars(r)["id"]

	stmt, err := database.Prepare("DELETE FROM users WHERE id=?")
	checkErr(err)

	_, err = stmt.Exec(Id)
	checkErr(err)
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func Profile(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "gosession")

	id := session.Values["user_id"]
	name := session.Values["user_name"]
	surname := session.Values["user_surname"]
	email := session.Values["user_email"]

	U := User{
		Id:      fmt.Sprintf("%v", id),
		Name:    fmt.Sprintf("%v", name),
		Surname: fmt.Sprintf("%v", surname),
		Email:   fmt.Sprintf("%v", email),
	}

	var templates = template.Must(template.ParseFiles("views/profile.gohtml", "views/base.gohtml"))
	varmap := map[string]interface{}{
		"Auth":    IsAuth(r),
		"Admin":   IsAdmin(r),
		"Report":  rp,
		"Profile": U,
	}
	_ = templates.ExecuteTemplate(w, "base", varmap)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "gosession")

	id := session.Values["user_id"]
	name := r.FormValue("name")
	surname := r.FormValue("surname")
	email := r.FormValue("email")

	stmt, err := database.Prepare("UPDATE users SET name=?, surname=?, email=? WHERE id=?")
	checkErr(err)
	_, err = stmt.Exec(name, surname, email, id)
	checkErr(err)
	rp = "Profile updated successfully."

	session.Values["user_name"] = name
	session.Values["user_surname"] = surname
	session.Values["user_email"] = email

	if uploadProfilePhoto(w, r, fmt.Sprintf("%v", id)) {
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

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func DeleteProfileForm(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.ParseFiles("views/delete_profil.gohtml", "views/base.gohtml"))
	varmap := map[string]interface{}{
		"Auth":  IsAuth(r),
		"Admin": IsAdmin(r),
	}

	err := templates.ExecuteTemplate(w, "base", varmap)
	checkErr(err)
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "gosession")
	id := session.Values["user_id"]
	U := User{}
	password := r.FormValue("password")

	err := database.QueryRow("SELECT password FROM users WHERE id=?", id).Scan(&U.Password)
	if err != nil {
		fmt.Println("Query failed")
	}

	if CheckPasswordHash(password, U.Password) {

		stmt, err := database.Prepare("DELETE FROM users WHERE id=?")
		if err != nil {
			fmt.Println("Query failed")
		}

		_, err = stmt.Exec(id)
		checkErr(err)

		session.Values["authenticated"] = false
		session.Values["user_id"] = ""
		session.Values["user_name"] = ""
		session.Values["user_surname"] = ""
		session.Values["user_email"] = ""
		session.Save(r, w)

		http.Redirect(w, r, "/register", http.StatusSeeOther)
	} else {
		rp := "Wrong Password"
		var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/base.gohtml"))
		varmap := map[string]interface{}{
			"Auth":   IsAuth(r),
			"Admin":  IsAdmin(r),
			"Report": rp,
		}
		_ = templates.ExecuteTemplate(w, "base", varmap)

	}
}

func uploadProfilePhoto(w http.ResponseWriter, r *http.Request, id string) bool {
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
	err = os.Rename(tempFile.Name(), "assets/profile/image-"+id+".jpg")

	if err != nil {
		log.Fatal(err)
	}

	return true
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.ParseFiles("views/change_password.gohtml", "views/base.gohtml"))
	varmap := map[string]interface{}{
		"Auth":  IsAuth(r),
		"Admin": IsAdmin(r),
	}

	err := templates.ExecuteTemplate(w, "base", varmap)
	checkErr(err)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "gosession")
	id := session.Values["user_id"]
	U := User{}

	err := database.QueryRow("SELECT password FROM users WHERE id=?", id).Scan(&U.Password)
	checkErr(err)

	current := r.FormValue("current")
	newpass := r.FormValue("newpass")
	newpass2 := r.FormValue("newpass2")

	p, _ := HashPassword(newpass)

	if newpass == newpass2 {
		if CheckPasswordHash(current, U.Password) {
			stmt, err := database.Prepare("UPDATE users SET password=? WHERE id=?")
			if err != nil {
				fmt.Println("Query failed")
			}
			_, err = stmt.Exec(p, id)
			checkErr(err)
			rp = "Password updated successfully."

			var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/base.gohtml"))
			varmap := map[string]interface{}{
				"Auth":   IsAuth(r),
				"Admin":  IsAdmin(r),
				"Report": rp,
			}
			_ = templates.ExecuteTemplate(w, "base", varmap)

			http.Redirect(w, r, "/profile", http.StatusSeeOther)
		} else {
			rp = "Password update failed"

			var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/base.gohtml"))
			varmap := map[string]interface{}{
				"Auth":   IsAuth(r),
				"Admin":  IsAdmin(r),
				"Report": rp,
			}
			_ = templates.ExecuteTemplate(w, "base", varmap)

			http.Redirect(w, r, "/profile", http.StatusSeeOther)
		}
	} else {
		rp = "Passwords did not match"
		var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/base.gohtml"))
		varmap := map[string]interface{}{
			"Auth":   IsAuth(r),
			"Admin":  IsAdmin(r),
			"Report": rp,
		}
		_ = templates.ExecuteTemplate(w, "base", varmap)
	}
}

func CommentsList(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "gosession")

	comment_author := session.Values["user_id"]

	post_id := r.FormValue("id")
	comment := r.FormValue("comment")

	t := time.Now()
	date := t.Format("2006-01-02 15:04:05")

	stmt, err := database.Prepare("INSERT INTO comments SET comment=?, date=?, post_related=?, comment_author=?")
	checkErr(err)
	_, err = stmt.Exec(comment, date, post_id, comment_author)
	checkErr(err)

	http.Redirect(w, r, "/post/"+post_id, http.StatusSeeOther)

}
