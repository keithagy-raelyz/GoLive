package app

import (
	"database/sql"
	"html/template"
	"net/http/httptest"
	"os"

	//"encoding/json"
	"fmt"
	"log"
	"net/http"

	"GoLive/db"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type App struct {
	router *mux.Router
	db     *db.Database
}

type Database interface {
}

// StartApp initializes the application (called by main).
func (a *App) StartApp() {
	godotenv.Load()
	a.connectDB()
	a.setRoutes()
	a.startRouter()
}

// Helpers for starting application.
func (a *App) connectDB() {
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	var err error
	a.db = &db.Database{}
	sqlDB, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.db.InitializeDB(sqlDB)
}

func (a *App) setRoutes() {

	a.router = mux.NewRouter()

	auth := a.db.InitializeAndGetAuth()
	a.router.Use(auth.Middleware)

	a.router.HandleFunc("/", home).Methods("GET")

	// Login Page
	a.router.HandleFunc("/login", a.allMerch).Methods("GET")    // Display Login Page
	a.router.HandleFunc("/login", a.allMerch).Methods("POST")   // Begin Validation
	a.router.HandleFunc("/login", a.allMerch).Methods("DELETE") // Delete session i.e logout

	// //Authentication
	// a.router.HandleFunc("/sessions", a.allMerch).Methods("GET") // ADMIN ONLY view/delete active sessions
	// a.router.HandleFunc("/sessions/{username}", a.allMerch).Methods("GET") // Validate session
	// a.router.HandleFunc("/sessions/{username}", a.allMerch).Methods("POST") // Validate login (user input)
	// a.router.HandleFunc("/sessions/{username}", a.allMerch).Methods("PUT") // Extend session
	// a.router.HandleFunc("/sessions/{username}", a.allMerch).Methods("GET") // Delete session (logout)

	// Get all Merchants and Products
	a.router.HandleFunc("/merchants", a.allMerch).Methods("GET")
	a.router.HandleFunc("/products", a.allProd).Methods("GET")
	a.router.HandleFunc("/users", a.allUser).Methods("GET")

	// Restful Route for Merchants
	a.router.HandleFunc("/merchants/{merchantid}", a.getMerch).Methods("GET")
	a.router.HandleFunc("/merchants", a.postMerch).Methods("POST")
	a.router.HandleFunc("/merchants/{merchantid}", a.putMerch).Methods("PUT")
	a.router.HandleFunc("/merchants/{merchantid}", a.delMerch).Methods("DELETE")

	// Restful Route for Products
	a.router.HandleFunc("/products/{productid}", a.getProd).Methods("GET")
	a.router.HandleFunc("/products", a.postProd).Methods("POST")
	a.router.HandleFunc("/products/{productid}", a.putProd).Methods("PUT")
	a.router.HandleFunc("/products/{productid}", a.delProd).Methods("DELETE")

	// Restful Route for Users
	a.router.HandleFunc("/users/{userid}", a.getUser).Methods("GET")
	a.router.HandleFunc("/users", a.postUser).Methods("POST")
	a.router.HandleFunc("/users/{userid}", a.putUser).Methods("PUT")
	a.router.HandleFunc("/users/{userid}", a.delUser).Methods("DELETE")

	// Restful Route for Cart
	a.router.HandleFunc("/cart", a.getUser).Methods("GET")
	a.router.HandleFunc("/cart/{productid}", a.getUser).Methods("POST")
	a.router.HandleFunc("/cart/{productid}", a.getUser).Methods("PUT")
	a.router.HandleFunc("/cart/{productid}", a.getUser).Methods("DELETE")

	// Checkout
	a.router.HandleFunc("/checkout", a.getUser).Methods("GET")
	a.router.HandleFunc("/checkout", a.getUser).Methods("POST")

}

func (a *App) startRouter() {
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", a.router))
}

func (a *App) TestRoute(recorder *httptest.ResponseRecorder, request *http.Request) {
	a.router.ServeHTTP(recorder, request)
}

func home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}
