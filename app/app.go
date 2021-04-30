package app

import (
	"database/sql"
	"net/http/httptest"
	"os"

	//"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/keithagy-raelyz/GoLive/db"
)

type App struct {
	router *mux.Router
	db     *db.Database
}

// StartApp initializes the application (called by main).
func (a *App) StartApp() {
	a.connectDB()
	a.setRoutes()
	a.startRouter()
}

// Helpers for starting application.
func (a *App) connectDB() {
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	var err error
	var db = &db.Database{}
	sqlDB, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	db.InitializeDB(sqlDB)
}

func (a *App) setRoutes() {
	a.router.HandleFunc("/", home).Methods("GET")

	//Get all Merchants
	a.router.HandleFunc("/merchants", a.allMerch).Methods("GET")

	//Restful Route for Merchants
	a.router.HandleFunc("/merchants/{merchantid}", a.getMerch).Methods("GET")
	a.router.HandleFunc("/merchants", a.postMerch).Methods("POST")
	a.router.HandleFunc("/merchants/{merchantid}", a.putMerch).Methods("PUT")
	a.router.HandleFunc("/merchants/{merchantid}", a.delMerch).Methods("DELETE")

	//Restful Route for Products
	a.router.HandleFunc("/products/{productid}", a.getProd).Methods("GET")
	a.router.HandleFunc("/products", a.postProd).Methods("POST")
	a.router.HandleFunc("/products/{productid}", a.putProd).Methods("PUT")
	a.router.HandleFunc("/products/{productid}", a.delProd).Methods("DELETE")
}

func (a *App) startRouter() {
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", a.router))
}

func (a *App) TestRoute(recorder *httptest.ResponseRecorder, request *http.Request) {
	a.router.ServeHTTP(recorder, request)
}

// Product Handler Functions
func (a *App) prod(w http.ResponseWriter, r *http.Request) {
	// GET method: Show a particular product given productid, show all products otherwise (USER / MERCHANT)
	// POST method: Add a Product (MERCHANT ONLY)
	// PUT method: Update existing Product (MERCHANT ONLY)
	// DELETE method: Delete existing Product (MERCHANT ONLY)

	// (a.)SomeCRUD(arg1 arg2)
}

func cart(w http.ResponseWriter, r *http.Request) {
	// GET method: Show a user's entire cart (USER ONLY)
	// POST method: Add a Product to a user's cart (USER ONLY)
	// PUT method: Edit item in user's cart (USER ONLY)
	// DELETE method: Delete existing Product (MERCHANT ONLY)
}

func users(w http.ResponseWriter, r *http.Request) {
	// GET method: Username registry (ADMIN ONLY)
	// POST method: Add new user (ADMIN ONLY)
	// PUT method: Edit new user (ADMIN ONLY)
	// DELETE method: Delete user (ADMIN ONLY)
}

func home(w http.ResponseWriter, r *http.Request) {

}
