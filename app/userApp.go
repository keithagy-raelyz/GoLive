package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/keithagy-raelyz/GoLive/db"
	"log"
	"net/http"
	"strings"
)

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// Verify valid merchant ID
	userID, _ := params["userid"]

	// Merchant ID supplied
	// Show all products under merchID; if invalid merchID handle error
	// TODO GetInventory errors need to be multiplexed to differentiate between invalid merchant and empty store
	var u db.User
	u, err := a.db.GetUser(userID)
	if err != nil {
		// Valid merchant ID but no products under merchant ID
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 - Valid merchant ID, but store is empty"))
		// return

		// Invalid merchant ID inputted
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No merchant for inputted merchant ID"))
		return
	}

	fmt.Println(u)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Valid merchant ID, displaying store"))
}

func (a *App) postUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	email := strings.ToLower(r.FormValue("email"))
	pw1 := r.FormValue("pw1")
	pw2 := r.FormValue("pw2")

	var u db.User

	u.Name = username
	u.Email = email
	if u.Name == "" {
		//TODO proper error handling
		log.Fatal()
	}
	if u.Email == "" {
		//TODO proper error handling
		log.Fatal()
	}
	err := a.db.CheckUser(u)
	if err != nil {
		//send the err msg back (err = errmsg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}
	if pw1 != pw2 {
		//t.ParseFiles("./templates/errorRegister.html")
		//data := Data{nil, ErrorMsg{"color:red", "Password entered are different"}, ErrorMsg{"display:none", ""}, ErrorMsg{"display:block", ""}, 0, nil}
		//t.Execute(w, data)
		return
	}
	err = a.db.CreateUser(u, pw1)
	if err != nil {
		//send the err msg back (err = errmsg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("201 - User Creation Successful"))
}

func (a *App) putUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, ok := params["userid"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Status Not Found"))
		// TODO display error message serve a proper template redirecting to registry of all merchants
		return
	}
	password := r.URL.Query().Get("password")
	if password == "" {
		//TODO proper error handling
		log.Fatal("no description supplied")
	}
	username := r.URL.Query().Get("username")
	if username == "" {
		//TODO proper error handling
		log.Fatal("no description supplied")
	}
	email := r.URL.Query().Get("email")
	if email == "" {
		//TODO proper error handling
		log.Fatal("no description supplied")
	}
	var u db.User
	u.Name = username
	u.Email = email
	u.Id = userID
	err := a.db.UpdateUser(u)
	if err != nil {
		//TODO proper error handling in template
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Unprocessable Entity"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Updated Successfully"))
}

func (a *App) delUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, ok := params["userid"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Status not Found"))
		return
	}
	err := a.db.DeleteUser(userID)
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
