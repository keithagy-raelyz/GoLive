package app

import (
	"log"
	"net/http"

	// "strings"

	"github.com/gorilla/mux"
)

// // MerchantUser has User's account details, with Description of storefront.
// type MerchantUser struct {
// 	User
// 	Description string
// }

// // Merchant contains MerchantUser (storefront details) and inventory.
// type Merchant struct {
// 	MerchantUser
// 	Products []Product
// }

func (a *App) allMerch(w http.ResponseWriter, r *http.Request) {

	// // No merchant ID supplied. Show all merchants
	// merchants, err := a.db.GetAllMerchants()
	// if err != nil {
	// 	panic(err.Error())
	// }

	// // TODO: Execute some template passing in merchants slice
	// return

}

// GET method - List products
// Callable by merchants and users
func (a *App) getMerch(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	// // Verify valid merchant ID
	// merchID, _ := params["merchantid"]

	// // Merchant ID supplied
	// // Show all products under merchID; if invalid merchID handle error
	// // TODO GetInventory errors need to be multiplexed to differentiate between invalid merchant and empty store
	// inventory, err := a.db.GetInventory(merchID)
	// if err != nil {
	// 	// Valid merchant ID but no products under merchant ID
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("200 - Valid merchant ID, but store is empty"))
	// 	return

	// 	// // Invalid merchant ID inputted
	// 	// w.WriteHeader(http.StatusNotFound)
	// 	// w.Write([]byte("404 - No merchant for inputted merchant ID"))
	// 	// return
	// }

	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("200 - Valid merchant ID, displaying store"))
}

// POST method - Add a Merchant (ADMIN ONLY)
func (a *App) postMerch(w http.ResponseWriter, r *http.Request) {

	// r.ParseForm()
	// username := r.FormValue("username")
	// email := strings.ToLower(r.FormValue("email"))
	// description := r.FormValue("description")
	// pw1 := r.FormValue("pw1")
	// pw2 := r.FormValue("pw2")

	// var m = &MerchantUser{}

	// m.Name = username
	// m.Email = email
	// if m.Name == "" {
	// 	//TODO proper error handling
	// 	log.Fatal()
	// }
	// if m.Email == "" {
	// 	//TODO proper error handling
	// 	log.Fatal()
	// }
	// err := a.db.CheckMerchant(m)
	// if err != nil {
	// 	//send the err msg back (err = errmsg)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte("400 - Bad Request"))
	// 	return
	// }
	// m.Description = description
	// if pw1 != pw2 {
	// 	//t.ParseFiles("./templates/errorRegister.html")
	// 	//data := Data{nil, ErrorMsg{"color:red", "Password entered are different"}, ErrorMsg{"display:none", ""}, ErrorMsg{"display:block", ""}, 0, nil}
	// 	//t.Execute(w, data)
	// 	return
	// }
	// err = a.db.CreateMerchant(m, pw1)
	// if err != nil {
	// 	//send the err msg back (err = errmsg)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte("400 - Bad Request"))
	// 	return
	// }

	// w.WriteHeader(http.StatusCreated)
	// w.Write([]byte("201 - User Creation Successful"))
}

// PUT method - Edit / Remove existing Merchant (ADMIN ONLY)
func (a *App) putMerch(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	// merchID, ok := params["merchantid"]
	// if !ok {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	w.Write([]byte("404 - Status Not Found"))
	// 	return
	// }
	// description := r.URL.Query().Get("description")
	// if description == "" {
	// 	//TODO proper error handling
	// 	log.Fatal("no description supplied")
	// }
	// username := r.URL.Query().Get("username")
	// if username == "" {
	// 	//TODO proper error handling
	// 	log.Fatal("no description supplied")
	// }
	// email := r.URL.Query().Get("email")
	// if email == "" {
	// 	//TODO proper error handling
	// 	log.Fatal("no description supplied")
	// }
	// var m MerchantUser
	// m.Name = username
	// m.Email = email
	// m.Id = merchID
	// err := a.db.UpdateMerchant(m)
	// if err != nil {
	// 	//TODO proper error handling in template
	// 	w.WriteHeader(http.StatusUnprocessableEntity)
	// 	w.Write([]byte("422 - Unprocessable Entity"))
	// 	return
	// }

	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("200 - Updated Successfully"))
}

// DELETE method - Delete merchant (ADMIN ONLY)
func (a *App) delMerch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	merchID, ok := params["merchantid"]
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("404 - Status not Found"))
		return
	}
	err := a.db.DeleteMerchant(merchID)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		//TODO proper error handling in template
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Unprocessable Entity"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Deleted Successfully"))
}
