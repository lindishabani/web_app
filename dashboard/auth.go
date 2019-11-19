package main

import (
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func LoginForm(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.ParseFiles("views/login.gohtml", "views/base.gohtml"))
	varmap := map[string]interface{}{
		"Auth":  IsAuth(r),
		"Admin": IsAdmin(r),
	}
	_ = templates.ExecuteTemplate(w, "base", varmap)
}

func LoginCheck(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "gosession")
	email := r.FormValue("email")
	password := r.FormValue("password")

	fmt.Println(email, password)

	U := User{}
	err := database.QueryRow("SELECT id, name, surname, email, password, role FROM users WHERE email=?", email).Scan(&U.Id, &U.Name, &U.Surname, &U.Email, &U.Password, &U.Role)
	fmt.Println(U.Password)
	if err != nil {
		var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/base.gohtml"))
		varmap := map[string]interface{}{
			"Auth":   IsAuth(r),
			"Admin":  IsAdmin(r),
			"Report": "Authentication failed",
		}
		_ = templates.ExecuteTemplate(w, "base", varmap)
	} else {
		if CheckPasswordHash(password, U.Password) {
			session.Values["authenticated"] = true
			session.Values["user_id"] = U.Id
			session.Values["user_name"] = U.Name
			session.Values["user_surname"] = U.Surname
			session.Values["user_email"] = U.Email
			session.Values["user_role"] = U.Role

			session.Save(r, w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			varmap := map[string]interface{}{
				"Auth":   IsAuth(r),
				"Admin":  IsAdmin(r),
				"Report": "Authentication failed",
			}
			var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/base.gohtml"))
			_ = templates.ExecuteTemplate(w, "base", varmap)
		}
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "gosession")

	session.Values["authenticated"] = false
	session.Values["user_id"] = ""
	session.Values["user_name"] = ""
	session.Values["user_surname"] = ""
	session.Values["user_email"] = ""
	session.Values["user_role"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func RegisterForm(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.ParseFiles("views/register.gohtml", "views/base.gohtml"))
	varmap := map[string]interface{}{
		"Auth":  IsAuth(r),
		"Admin": IsAdmin(r),
	}
	_ = templates.ExecuteTemplate(w, "base", varmap)
}

func RegistrationCheck(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	surname := r.FormValue("surname")
	email := r.FormValue("email")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")

	if password1 == password2 {
		sqlquery := "INSERT INTO users SET name=?, surname=?, email=?, password=?"

		stmt, err := database.Prepare(sqlquery)
		checkErr(err)
		p, _ := HashPassword(password1)
		_, err = stmt.Exec(name, surname, email, p)
		checkErr(err)

		rp = "Regjistrimi u krye me sukses."
		var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/base.gohtml"))
		varmap := map[string]interface{}{
			"Auth":   IsAuth(r),
			"Admin":  IsAdmin(r),
			"Report": rp,
		}
		_ = templates.ExecuteTemplate(w, "base", varmap)
	} else {
		rp = "Fjalekalimet nuk perputhen"
		var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/base.gohtml"))
		varmap := map[string]interface{}{
			"Auth":   IsAuth(r),
			"Admin":  IsAdmin(r),
			"Report": rp,
		}
		_ = templates.ExecuteTemplate(w, "base", varmap)
	}

}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	fmt.Println(string(bytes))
	return string(bytes), err
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func IsAuth(r *http.Request) bool {
	session, _ := store.Get(r, "gosession")
	var ok bool
	auth, ok = session.Values["authenticated"].(bool)

	if ok {
		return auth
	}
	return false
}

func IsAdmin(r *http.Request) bool {
	session, _ := store.Get(r, "gosession")

	if session.Values["user_role"] == 1 {
		return true
	}
	return false
}

func authMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if IsAuth(r) {
			nextHandler.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func adminMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if IsAuth(r) {
			session, _ := store.Get(r, "gosession")
			id := session.Values["user_id"]
			U := User{}
			err := database.QueryRow("SELECT role FROM users WHERE id=?", id).Scan(&U.Role)
			if err != nil {
				fmt.Println("Gabim ne SQL")
			}

			fmt.Println(database)

			if U.Role == 1 {
				nextHandler.ServeHTTP(w, r)
			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
