package app

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"GoLive/db"

	"github.com/gorilla/mux"
)

func (a *App) allProd(w http.ResponseWriter, r *http.Request) {

	// No product ID supplied. Show all products
	products, err := a.db.GetAllProducts()
	if err != nil {
		panic(err.Error())
	}

	// TODO: Execute some template passing in products slice
	fmt.Println(products)
}

func (a *App) getProd(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// Verify valid product ID
	prodID := params["productid"]

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
	name := r.FormValue("Name")
	ProdDesc := r.FormValue("ProdDesc")
	thumbnail := r.FormValue("Thumbnail")
	price, err := strconv.ParseFloat(r.FormValue("Price"), 64)
	if err != nil {
		// TODO error handling
		log.Fatal(err)
	}
	quantity, err := strconv.Atoi(r.FormValue("Quantity"))
	if err != nil {
		// TODO error handling
		log.Fatal(err)
	}

	merchID := r.FormValue("MerchID")
	if err != nil {
		// TODO error handling
		log.Fatal(err)
	}

	if price <= 0 || quantity < 0 || ProdDesc == "" || name == "" {

		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Invalid Price and/or Quantity submitted"))
		return
	}

	// TODO Session handling to get MerchID
	p := db.Product{Name: name, ProdDesc: ProdDesc, Thumbnail: thumbnail, Price: price, Quantity: quantity, MerchID: merchID, Sales: 0}
	err = a.db.CreateProduct(p)
	if err != nil {
		// TODO error handling
		fmt.Println(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Merchant ID provided is Invalid"))
		return
	}
	fmt.Println(p, "TESTING PRODUCT")
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

	name := r.URL.Query().Get("Name")
	if name == "" {
		//TODO proper error handling
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Invalid Name submitted"))
		return
	}

	quantity, err := strconv.Atoi(r.URL.Query().Get("Quantity"))
	if r.URL.Query().Get("Quantity") == "" || err != nil {
		//TODO proper error handling
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Invalid Quantity submitted"))
		return
	}

	price, err := strconv.ParseFloat(r.URL.Query().Get("Price"), 64)
	if r.URL.Query().Get("Price") == "" || err != nil {
		//TODO proper error handling
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Invalid Price submitted"))
		return
	}

	ProdDesc := r.URL.Query().Get("ProdDesc")
	if ProdDesc == "" {
		//TODO proper error handling
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Invalid ProdDesc submitted"))
		return
	}

	merchID := r.URL.Query().Get("MerchID")
	if merchID == "" {
		//TODO proper error handling
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Invalid MerchID submitted"))
		return
	}

	thumbnail := r.URL.Query().Get("Thumbnail")
	if thumbnail == "" {
		//TODO proper error handling
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Invalid Thumbnail submitted"))
		return
	}

	// TODO session handling to supply correct MerchID
	p := db.Product{
		Id:        prodID,
		Name:      name,
		Quantity:  quantity,
		Thumbnail: thumbnail,
		Price:     price,
		ProdDesc:  ProdDesc,
		MerchID:   merchID,
	}
	err = a.db.UpdateProduct(p)
	if err != nil {
		fmt.Println(err)
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
		//TODO proper error handling in template
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Unprocessable Entity"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Deleted Successfully"))
}
