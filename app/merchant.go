package app

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// MerchantUser has User's account details, with Description of storefront.
type MerchantUser struct {
	User
	Description string
}

// Merchant contains MerchantUser (storefront details) and inventory.
type Merchant struct {
	MerchantUser
	Products []Product
}

// GET method - List products
// Callable by merchants and users
func (a *App) getMerch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// Verify valid merchant ID
	merchID, ok := params["merchantid"]
	if !ok {
		// No merchant ID supplied. Show all merchants
		merchants, err := a.db.GetAllMerchants
		if err != nil {
			panic(err.Error())
		}

		// TODO: Execute some template passing in merchants slice
		return
	}

	// Merchant ID supplied
	// Show all products under merchID; if invalid merchID handle error
	// TODO GetInventory errors need to be multiplexed to differentiate between invalid merchant and empty store
	inventory, err := a.db.GetInventory(merchID)
	if err != nil {
		// Valid merchant ID but no products under merchant ID
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 - Valid merchant ID, but store is empty"))
		return

		// // Invalid merchant ID inputted
		// w.WriteHeader(http.StatusNotFound)
		// w.Write([]byte("404 - No merchant for inputted merchant ID"))
		// return
	}
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

	err := db.QueryRow("SELECT username,email FROM users where Username=? OR email=?", name, email).Scan(m.Name, m.Email)
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

	res, err := db.Exec("DELETE FROM products where merchantID =?", merchID)
	if err != nil {
		log.Fatal(err)
	}
	_, err = res.RowsAffected()

	res, err = db.Exec("DELETE FROM merchants where merchantID =?", merchID)
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
