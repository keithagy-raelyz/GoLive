package main

import "net/http"

// GET method: Show a particular product given productid, show all products otherwise (USER / MERCHANT)
func getProd(w http.ResponseWriter, r *http.Request) {
}

// POST method: Add a Product to their own store (MERCHANT ONLY)
func postProd(w http.ResponseWriter, r *http.Request) {
}

// PUT method: Update existing Product on their own store (MERCHANT ONLY)
func putProd(w http.ResponseWriter, r *http.Request) {
}

// DELETE method: Delete existing Product from their own store (MERCHANT ONLY)
func delProd(w http.ResponseWriter, r *http.Request) {

}
