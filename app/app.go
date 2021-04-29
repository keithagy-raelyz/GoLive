package app

import (
	"database/sql"
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

// StartApp initializes the appliation (called by main).
func (a *App) StartApp() {
	a.connectDB()
	a.setRoutes()
	a.startRouter()
}

// Helpers for starting application.
func (a *App) connectDB() {
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	var err error
	a.db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	initializeDB(a.db)
}

func (a *App) setRoutes() {
	a.router.HandleFunc("/", home).Methods("GET")

	a.router.HandleFunc("/merchants/{merchantid}", getMerch).Methods("GET")
	a.router.HandleFunc("/merchants", postMerch).Methods("POST")
	a.router.HandleFunc("/merchants/{merchantid}", putMerch).Methods("PUT")
	a.router.HandleFunc("/merchants/{merchantid}", delMerch).Methods("DELETE")

	a.router.HandleFunc("/product/{productid}", getProd).Methods("GET")
	a.router.HandleFunc("/product", postProd).Methods("POST")
	a.router.HandleFunc("/product/{productid}", putProd).Methods("PUT")
	a.router.HandleFunc("/product/{productid}", delProd).Methods("DELETE")
}

func (a *App) startRouter() {
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", a.router))
}

// Initialize the DB schema
func initializeDB(db *sql.DB) {
	Query1 := `CREATE TABLE IF NOT EXISTS Users (
    UserID int NOT NULL AUTO_INCREMENT,
    Username VARCHAR(255) NOT NULL,
    Password VARCHAR(255) NOT NULL,
    Email varchar(255) NOT NULL,
    PRIMARY KEY (UserID)
	)`
	_, err := db.Exec(Query1)
	if err != nil {
		log.Fatal(err)
	}
	Query2 := `
	CREATE TABLE IF NOT EXISTS Merchants (
		MerchantID int NOT NULL AUTO_INCREMENT,
		Username VARCHAR(255) NOT NULL,
		Password VARCHAR(255) NOT NULL,
		Email varchar(255) NOT NULL,
		Description VARCHAR(255) NOT NULL,
		PRIMARY KEY (MerchantID)
	);`
	_, err = db.Exec(Query2)
	if err != nil {
		log.Fatal(err)
	}
	Query3 := `CREATE TABLE IF NOT EXISTS Products (
    ProductID int NOT NULL AUTO_INCREMENT,
    Product_Name VARCHAR(255) NOT NULL,
    Quantity int NOT NULL,
    Image varchar(255) NOT NULL,
    Price float not null,
    Description VARCHAR(255),
    MerchantID int NOT NULL,
    Foreign Key (MerchantID) REFERENCES Merchants (MerchantID),
    PRIMARY KEY (ProductID)
	);`
	_, err = db.Exec(Query3)
	if err != nil {
		log.Fatal(err)
	}
}

// Product Handler Functions
func (a *App) prod(w http.ResponseWriter, r *http.Request) {
	// GET method: Show a particular product given productid, show all products otherwise (USER / MERCHANT)
	// POST method: Add a Product (MERCHANT ONLY)
	// PUT method: Update existing Product (MERCHANT ONLY)
	// DELETE method: Delete existing Product (MERCHANT ONLY)

	// (a.db)SomeCRUD(arg1 arg2)
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
