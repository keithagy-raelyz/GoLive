package app

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/keithagy-raelyz/GoLive/db"
)

func (a *App) allProd(w http.ResponseWriter, r *http.Request) {

	// No product ID supplied. Show all products
	products, err := a.db.GetAllProducts()
	if err != nil {
		panic(err.Error())
	}

	// TODO: Execute some template passing in products slice
	fmt.Println(products)
	return
}

func (a *App) getProd(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// Verify valid product ID
	prodID, _ := params["productid"]

	// Product ID supplied
	// Show all products under prodID; if invalid prodID handle error
	product, err := a.db.GetProduct(prodID)
	if err != nil {
		// Product ID not registered
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Product ID not registered"))
		return
	}

	fmt.Println(product)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Valid merchant ID, displaying store"))
}

func (a *App) postProd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// Checking for non-empty inputs to be handled by HTML form
	name := r.FormValue("name")
	description := r.FormValue("description")
	thumbnail := r.FormValue("thumbnail")
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		// TODO error handling
		log.Fatal(err)
	}
	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil {
		// TODO error handling
		log.Fatal(err)
	}

	// TODO Session handling to get MerchID
	p := db.Product{name, 0, description, thumbnail, price, quantity, 0}
	err = a.db.CreateProduct(p)
	if err != nil {
		// TODO error handling
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("201 - Product Creation Successful"))
}

func (a *App) putProd(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	prodID, ok := params["productid"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Status Not Found"))
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		//TODO proper error handling
		log.Fatal("no description supplied")
	}

	description := r.URL.Query().Get("description")
	if description == "" {
		//TODO proper error handling
		log.Fatal("no description supplied")
	}

	thumbnail := r.URL.Query().Get("thumbnail")
	if thumbnail == "" {
		//TODO proper error handling
		log.Fatal("no description supplied")
	}

	price := r.URL.Query().Get("price")
	if thumbnail == "" {
		//TODO proper error handling
		log.Fatal("no description supplied")
	}

	quantity := r.URL.Query().Get("quantity")
	if thumbnail == "" {
		//TODO proper error handling
		log.Fatal("no description supplied")
	}

	// TODO session handling to supply correct MerchID
	p := db.Product{name, prodID, description, thumbnail, price, quantity, 0}
	err := a.db.UpdateProduct(p)
	if err != nil {
		//TODO proper error handling in template
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Unprocessable Entity"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Updated Successfully"))
}

func (a *App) delProd(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	prodID, ok := params["productid"]
	if !ok {
		//TODO proper error handling in template
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Status not Found"))
		return
	}

	// TODO Session handling to pass in correct merchID
	err := a.db.DeleteProduct(prodID, "0")
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
