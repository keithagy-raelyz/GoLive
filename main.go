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
	Name        string
	Description string
	Thumbnail   string
	Price       float64
	Quantity    int
}

type Merchant struct {
	Id          int
	Name        string
	Description string
	Products    []Product
}

const (
	homePath  = "/"
	merchPath = "/merchants/{merchantid}"
	prodPath  = "/products/{productid}"
)

var (
	db     *sql.DB
	router = mux.NewRouter()
)

func main() {
	router.HandleFunc(homePath, home)
	router.HandleFunc(merchPath, merch)
	router.HandleFunc(prodPath, prod)

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

// Handler Functions
func merch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// GET method - List products
	// If no merchantid, show all merchants
	// If have merchantid, show all products under merchantid (incl. nothing)
	// If invalid merchantid, show error and redirect back to home page
	// Callable by merchants and users
	if r.Method == "GET" {
		// Verify valid merchant ID
		merchID, ok := params["merchantid"]
		if !ok {
			// No merchant ID supplied. Show all merchants
			merchantRows, err := db.Query("Select username, description FROM storefront_db.merchants")
			if err != nil {
				panic(err.Error())
			}
			defer merchantRows.Close()

			var merchants = make([]Merchant, 0)
			for merchantRows.Next() {
				var newMerchant Merchant
				err = merchantRows.Scan(&newMerchant.Id, &newMerchant.Name, &newMerchant.Description, &newMerchant.Products)
				if err != nil {
					panic(err.Error())
				}
				merchants = append(merchants, newMerchant)
			}

			// TODO: Execute some template passing in merchants slice
			return
		}

		// Merchant ID supplied
		// Show all products under merchID; if invalid merchID handle error
		// TODO change query string
		merchProdsRows, err := db.Query("Select * FROM storefront_db.products WHERE merchid = ? AND quantity != 0", merchID)
		if err != nil {
			panic(err.Error())
		}
		defer merchProdsRows.Close()

		var merchProds = make([]Product, 0)
		for merchProdsRows.Next() {
			var p Product
			err = merchProdsRows.Scan(&p.Name, &p.Description, &p.Thumbnail, &p.Price, &p.Quantity)
			if err != nil {
				// Valid merchant ID but no products under merchant ID
				w.WriteHeader(http.StatusNoContent)
				w.Write([]byte("204 - Valid merchant ID, but store is empty"))
				return
			}
			merchProds = append(merchProds, p)
		}

		if len(merchProds) == 0 {
			// Invalid merchant ID inputted
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No merchant for inputted merchant ID"))
			return
		}
	}
	// POST method - Add a Merchant (ADMIN ONLY)
	// PUT method - Edit / Remove existing Merchant (ADMIN ONLY)
	// DELETE method - Delete merchant (ADMIN ONLY)
}

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
