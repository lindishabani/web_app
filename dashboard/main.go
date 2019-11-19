package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	DBHost     = "127.0.0.1"
	DBPort     = ":3306"
	DBUser     = "root"
	DBPass     = ""
	DBDbase    = "blog"
	ASSETS_DIR = "/assets/"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type Post struct {
	Id            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Category_id   string `json:"category_id"`
	Date          string `json:"date"`
	Author        string `json:"author_id"`
	Id_category   string `json:"id_category"`
	Category      string `json:"category"`
	Date_category string `json:"date_category"`
}

type Category struct {
	Id       string `json:"id"`
	Category string `json:"category"`
	Date     string `json:"date"`
}

type Comment struct {
	Id              string `json:"id"`
	Comment         string `json:"comment"`
	Date            string `json:"date"`
	Post_related    string `json:"post_related"`
	Comment_author  string `json:"comment_author"`
	Post_related_id string `json:"post_related_id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Category_id     string `json:"category_id"`
	Post_date       string `json:"post_date"`
	Author_id       string `json:"author_id"`
	Author_c_id     string `json:"author_c_id"`
	Author_name     string `json:"author_name"`
	Author_surname  string `json:"author_surname"`
	Author_email    string `json:"author_email"`
	Author_password string `json:"author_password"`
	Author_role     string `json:"author_role"`
}

var (
	database *sql.DB
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("ahfy7634*&^%23DFERko923456df0&")
	store = sessions.NewCookieStore(key)
	auth  bool
	rp    string
)

func main() {
	dbConn := fmt.Sprintf("%s:%s@tcp(%s)/%s", DBUser, DBPass, DBHost, DBDbase)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		panic("Couldn't connect!")
	}

	database = db
	routes := NewRouter()

	// Public
	routes.HandleFunc("/post/{id:[0-9]+}", ShowPost).Methods("GET")
	routes.HandleFunc("/", ShowPosts).Methods("GET")

	// Posts
	routes.Handle("/dashboard", adminMiddleware(http.HandlerFunc(Dashboard))).Methods("GET")
	routes.Handle("/add_post", adminMiddleware(http.HandlerFunc(NewPost))).Methods("GET")
	routes.Handle("/edit/{id:[0-9]+}", adminMiddleware(http.HandlerFunc(EditPost))).Methods("GET")
	routes.Handle("/update", adminMiddleware(http.HandlerFunc(UpdatePost))).Methods("POST")
	routes.Handle("/delete/{id:[0-9]+}", adminMiddleware(http.HandlerFunc(DeletePost))).Methods("GET")

	// Categories
	routes.Handle("/categories", adminMiddleware(http.HandlerFunc(Categories))).Methods("GET")
	routes.Handle("/add_category", adminMiddleware(http.HandlerFunc(NewCategory))).Methods("GET")
	routes.Handle("/edit_category/{id:[0-9]+}", adminMiddleware(http.HandlerFunc(EditCategory))).Methods("GET")
	routes.Handle("/update_category", adminMiddleware(http.HandlerFunc(UpdateCategory))).Methods("POST")
	routes.Handle("/delete_category/{id:[0-9]+}", adminMiddleware(http.HandlerFunc(DeleteCategory))).Methods("GET")

	// Users
	routes.Handle("/users", adminMiddleware(http.HandlerFunc(Users))).Methods("GET")
	routes.Handle("/edit_user/{id:[0-9]+}", adminMiddleware(http.HandlerFunc(EditUser))).Methods("GET")
	routes.Handle("/update_user", adminMiddleware(http.HandlerFunc(UpdateUser))).Methods("POST")
	routes.Handle("/delete_user/{id:[0-9]+}", adminMiddleware(http.HandlerFunc(DeleteUser))).Methods("GET")
	routes.Handle("/comment", authMiddleware(http.HandlerFunc(CommentsList))).Methods("POST")

	// Profile
	routes.Handle("/profile", authMiddleware(http.HandlerFunc(Profile))).Methods("GET")
	routes.Handle("/update_profile", authMiddleware(http.HandlerFunc(UpdateProfile))).Methods("POST")
	routes.Handle("/delete_profile_form", authMiddleware(http.HandlerFunc(DeleteProfileForm))).Methods("GET")
	routes.Handle("/delete_profile", authMiddleware(http.HandlerFunc(DeleteProfile))).Methods("GET")
	routes.Handle("/change_password", authMiddleware(http.HandlerFunc(ChangePassword))).Methods("GET")
	routes.Handle("/update_password", authMiddleware(http.HandlerFunc(UpdatePassword))).Methods("POST")

	// Login
	routes.HandleFunc("/login", LoginForm).Methods("GET")
	routes.HandleFunc("/logout", Logout).Methods("GET")
	routes.HandleFunc("/logincheck", LoginCheck).Methods("POST")

	// Register
	routes.HandleFunc("/register", RegisterForm).Methods("GET")
	routes.HandleFunc("/registrationcheck", RegistrationCheck).Methods("POST")

	log.Println(http.ListenAndServe(":8080", routes))
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Serve CSS, JS & Images Statically.
	router.PathPrefix(ASSETS_DIR).Handler(http.StripPrefix(ASSETS_DIR, http.FileServer(http.Dir("."+ASSETS_DIR))))

	return router
}
