package app

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"strings"

	"GoLive/db"

	"github.com/gorilla/mux"
)

func (a *App) allMerch(w http.ResponseWriter, r *http.Request) {

	// No merchant ID supplied. Show all merchants
	merchants, err := a.db.GetAllMerchants()
	if err != nil {
		panic(err.Error())
	}

	// TODO: Execute some template passing in merchants slice
	fmt.Println("Merchants:", merchants)

	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("200 - Displaying all merchants"))
	t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/viewAllMerchants.html", "templates/error.html")
	if err != nil {
		log.Fatal(err)
	}
	data := Data{Merchants: merchants}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
	return

}

// GET method - List products
// Callable by merchants and users
func (a *App) getMerch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// Verify valid merchant ID
	merchID := params["merchantid"]

	// Merchant ID supplied
	// Show all products under merchID; if invalid merchID handle error
	// TODO inventory has to be renamed to accurately reflect no rows were found
	merchant, err := a.db.GetInventory(merchID)
	if err != nil {
		// Merchant exist but no product
		fmt.Println("Blank inventory for merchid:", merchID, merchant)
		w.WriteHeader(http.StatusOK)
		//w.Write([]byte("200 - Valid merchID but empty inventory"))
		// fmt.Println("200, empty inv, merchID:", merchID)
		t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/viewMerchantProducts.html", "templates/error.html")
		if err != nil {
			log.Fatal(err)
		}
		data := Data{Merchant: merchant}
		err = t.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	//TODO is this logic correct?
	if merchant.Products == nil {
		// Invalid Merchant ID
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Invalid merchID"))
		// fmt.Println("404, invalid merchID, merchID:", merchID)
		return
	}
	// fmt.Println("Inventory for merchID:", merchID, inventory)
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("200 - Valid merchant ID, displaying store"))
	t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/viewMerchantProducts.html", "templates/error.html")
	if err != nil {
		log.Fatal(err)
	}
	data := Data{Merchant: merchant}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// POST method - Add a Merchant (ADMIN ONLY)
func (a *App) postMerch(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	username := r.FormValue("username")
	email := strings.ToLower(r.FormValue("email"))
	MerchDesc := r.FormValue("MerchDesc")
	pw1 := r.FormValue("pw1")
	pw2 := r.FormValue("pw2")

	var m db.MerchantUser

	m.Name = username
	m.Email = email
	if m.Name == "" {
		//TODO proper error handling
		log.Fatal()
	}
	if m.Email == "" {
		//TODO proper error handling
		log.Fatal()
	}
	err := a.db.CheckMerchant(m)
	if err != nil {

		//check MerchDesc and pw
		m.MerchDesc = MerchDesc
		if pw1 != pw2 {
			//t.ParseFiles("./templates/errorRegister.html")
			//data := Data{nil, ErrorMsg{"color:red", "Password entered are different"}, ErrorMsg{"display:none", ""}, ErrorMsg{"display:block", ""}, 0, nil}
			//t.Execute(w, data)
			return
		}
		err = a.db.CreateMerchant(m, pw1)

		if err != nil {
			//send the err msg back (err = errmsg)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Bad Request"))
			return
		}

		//send the err msg back (err = errmsg)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("201 - User Creation Successful"))
		return

	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("400 - Bad Request"))

}

// PUT method - Edit existing Merchant (ADMIN ONLY)
func (a *App) putMerch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	merchID, ok := params["merchantid"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Status Not Found"))
		// TODO display error message serve a proper template redirecting to registry of all merchants
		return
	}
	MerchDesc := r.URL.Query().Get("MerchDesc")
	if MerchDesc == "" {
		//TODO proper error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - no MerchDesc supplied"))
		return
	}
	username := r.URL.Query().Get("username")
	if username == "" {
		//TODO proper error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - no username supplied"))
		return
	}
	email := r.URL.Query().Get("email")
	if email == "" {
		//TODO proper error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - no email supplied"))
		return
	}
	var m db.MerchantUser
	m.Name = username
	m.Email = email
	m.Id = merchID
	err := a.db.UpdateMerchant(m)
	if err != nil {
		//TODO proper error handling in template
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Unprocessable Entity"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Updated Successfully"))
}

// DELETE method - Delete merchant (ADMIN ONLY)
func (a *App) delMerch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	merchID, ok := params["merchantid"]
	// fmt.Println(merchID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Status not Found"))
		return
	}
	err := a.db.DeleteMerchant(merchID)
	if err != nil {
		// fmt.Println(err)
		//TODO proper error handling in template
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Unprocessable Entity"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Deleted Successfully"))
}
