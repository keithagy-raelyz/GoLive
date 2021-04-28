package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Product struct {
	name        string
	description string
	thumbnail   string
	price       float64
	quantity    int
}

type Merchant struct {
	Id       int
	Products []Product
}

const (
	merchPath = "/storefront/v0/{merchantid}/{productid}"
)

var (
	db     *sql.DB
	router = mux.NewRouter()
)

func main() {
	router.HandleFunc(merchPath, store)

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

// Handler Functions
func store(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// GET method - List products
	// If no merchantid, list all merchants
	// If no productid, list all products under a particular merchantid
	// Callable by merchants and users
	if r.Method == "GET" {
		// Verify valid merchant ID
		merchID, ok := params["merchantid"]
		if !ok {
			// No merchant ID supplied
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No merchant ID supplied"))
			return
		}

		// Verify valid product ID
		prodID, ok := params["productid"]
		if !ok {
			// Show all products under merchID; if invalid merchID handle error
			merchProdsRows, err := db.Query("Select * FROM storefront_db.products WHERE merchid = ? AND quantity != 0", merchID)
			defer merchProdsRows.Close()
			if err != nil {
				panic(err.Error())
			}

			var merchProds = make([]Product, 0)
			for merchProdsRows.Next() {
				var p Product
				err = merchProdsRows.Scan(&p.name, &p.description, &p.thumbnail, &p.price, &p.quantity)
				if err != nil {
					panic(err.Error())
				}
				merchProds = append(merchProds, p)
			}

			if len(merchProds) == 0 {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 - No merchant for inputted merchant ID"))
				return
			}

		} else {
			// Check if prodID is valid
			// Show merchID -> prodID
		}
	}

	// POST method - Add a Product to Store
	// Valid input would REQUIRE a specific merchantid and productid
	// Can only be called by merchants, on their own store (validation)

	// PUT method - Edit / Remove existing Product
	// Valid input would REQUIRE a specific merchantid and productid
	// Can only be called by merchants, on their own store (validation)

	// DELETE method - Delete the store (less important)
}

func cart(w http.ResponseWriter, r *http.Request) {
}

func merchants(w http.ResponseWriter, r *http.Request) {
}

func users(w http.ResponseWriter, r *http.Request) {
}
