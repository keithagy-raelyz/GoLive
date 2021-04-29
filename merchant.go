package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type MerchantUser struct {
	User
	Description string
}

type Merchant struct {
	MerchantUser
	Products []Product
}

// GET method - List products
//Callable by merchants and users
func getMerch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// Verify valid merchant ID
	merchID, ok := params["merchantid"]
	if !ok {
		// No merchant ID supplied. Show all merchants
		merchantRows, err := db.Query("SELECT merchantID, Username, description FROM merchants")
		if err != nil {
			panic(err.Error())
		}
		defer merchantRows.Close()

		var merchants = make([]Merchant, 0)
		for merchantRows.Next() {
			var newMerchant Merchant
			err = merchantRows.Scan(&newMerchant.Id, &newMerchant.Name, &newMerchant.Description)
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
	merchProdsRows, err := db.Query("SELECT username, merchants.merchantid, merchants.description, products.ProductID, products.Product_Name, products.Quantity, products.Image, products.price,products.Description from merchants LEFT JOIN products on products.merchantid = merchants.merchantid where merchants.merchantid = ?;", merchID)
	if err != nil {
		panic(err.Error())
	}
	defer merchProdsRows.Close()

	var merchProds = make([]Product, 0)
	var merch = &Merchant{}
	for merchProdsRows.Next() {
		var p Product
		err = merchProdsRows.Scan(&merch.Name, &merch.Id, &merch.Description, &p.Id, &p.Name, &p.Quantity, &p.Thumbnail, &p.Price, &p.Description)
		if err != nil {
			// Valid merchant ID but no products under merchant ID
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("200 - Valid merchant ID, but store is empty"))
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
	// If invalid merchantid, show error and redirect back to home page
}

// POST method - Add a Merchant (ADMIN ONLY)
func postMerch(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	name := r.FormValue("username")
	email := strings.ToLower(r.FormValue("email"))
	description := r.FormValue("description")
	pw1 := r.FormValue("pw1")
	pw2 := r.FormValue("pw2")

	var m = &MerchantUser{}

	err := db.QueryRow("SELECT username,email FROM users where name=? OR email=?", name, email).Scan(m.Name, m.Email)
	if err != nil {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("201 - User Creation Successful"))
	}
	if m.Name == name || name == "" {
		//t.ParseFiles("./templates/errorRegister.html")
		//data := Data{nil, ErrorMsg{"color:red", "Email is already in use"}, ErrorMsg{"display:none", ""}, ErrorMsg{"display:block", ""}, 0, nil}
		//t.Execute(w, data)
		return
	}
	if m.Email == email {
		//t.ParseFiles("./templates/errorRegister.html")
		//data := Data{nil, ErrorMsg{"color:red", "Username is already in use"}, ErrorMsg{"display:none", ""}, ErrorMsg{"display:block", ""}, 0, nil}
		//t.Execute(w, data)
		return
	}
	if pw1 != pw2 {
		//t.ParseFiles("./templates/errorRegister.html")
		//data := Data{nil, ErrorMsg{"color:red", "Password entered are different"}, ErrorMsg{"display:none", ""}, ErrorMsg{"display:block", ""}, 0, nil}
		//t.Execute(w, data)
		return
	}
	if description == "" {
		//t.ParseFiles("./templates/errorRegister.html")
		//data := Data{nil, ErrorMsg{"color:red", "Email is already in use"}, ErrorMsg{"display:none", ""}, ErrorMsg{"display:block", ""}, 0, nil}
		//t.Execute(w, data)
		return
	}
}

// PUT method - Edit / Remove existing Merchant (ADMIN ONLY)
func putMerch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	merchID, ok := params["merchantid"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Status Not Found"))
		return
	}
	description := r.URL.Query().Get("description")
	if description == "" {
		log.Fatal("no description supplied")
	}
	res, err := db.Exec("UPDATE merchants set description =? where merchantID = ?", description, merchID)
	if err != nil {
		log.Fatal(err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Unprocessable Entity"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Updated Successfully"))
}

// DELETE method - Delete merchant (ADMIN ONLY)
func delMerch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	merchID, ok := params["merchantid"]
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("404 - Status not Found"))
		return
	}
	res, err := db.Exec("DELETE FROM merchants where merchantID =?", merchID)
	if err != nil {
		log.Fatal(err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Unprocessable Entity"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Deleted Successfully"))
}
