package main

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Product struct {
	Name        string
	Id          int
	Description string
	Thumbnail   string
	Price       float64
	Quantity    int
}

type User struct {
	Id    string
	Name  string
	Email string
}

var (
	db     *sql.DB
	router = mux.NewRouter()
)

func main() {
	router.HandleFunc("/", home).Methods("GET")

	router.HandleFunc("/merchants/{merchantid}", getMerch).Methods("GET")
	router.HandleFunc("/merchants", postMerch).Methods("POST")
	router.HandleFunc("/merchants/{merchantid}", putMerch).Methods("PUT")
	router.HandleFunc("/merchants/{merchantid}", delMerch).Methods("DELETE")

	router.HandleFunc("/product/{productid}", getProd).Methods("GET")
	router.HandleFunc("/product", postProd).Methods("POST")
	router.HandleFunc("/product/{productid}", putProd).Methods("PUT")
	router.HandleFunc("/product/{productid}", delProd).Methods("DELETE")

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

// Handler Functions

func prod(w http.ResponseWriter, r *http.Request) {
	// GET method: Show a particular product given productid, show all products otherwise (USER / MERCHANT)
	// POST method: Add a Product (MERCHANT ONLY)
	// PUT method: Update existing Product (MERCHANT ONLY)
	// DELETE method: Delete existing Product (MERCHANT ONLY)
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
