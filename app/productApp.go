package app

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"GoLive/cache"
	"GoLive/db"

	"github.com/gorilla/mux"
)

func (a *App) allProd(w http.ResponseWriter, r *http.Request) {
	p, _ := a.cacheManager.GetAllProducts()
	data := Data{
		Products: p,
	}
	activeSession, ok := a.HaveValidSessionCookie(r)
	fmt.Println(activeSession)
	if !ok {
		parseAllProducts(&w, data)
		return
	}
	switch v := activeSession.(type) {
	case *cache.UserSession:
		user, cart := v.GetSessionOwner()
		data.User, data.Cart = user.User, cart
	case *cache.MerchantSession:
		data.Merchant.MerchantUser, _ = v.GetSessionOwner()
	}
	parseAllProducts(&w, data)
}

func (a *App) getProd(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// Verify valid product ID
	prodID := params["productid"]

	// Product ID supplied
	// Show all products under prodID; if invalid prodID handle error
	product, err := a.cacheManager.GetFromCache(prodID, "cachedItems")
	var products = []db.Product{product.(*cache.CachedItem).Copy()}
	data := Data{Products: products}
	if err != nil {
		// Product ID not registered
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Product ID not registered"))
		return
	}
	activeSession, ok := a.HaveValidSessionCookie(r)
	if !ok {
		parseProductPage(&w, data)
		return
	}
	switch v := activeSession.(type) {
	case *cache.UserSession:
		user, cart := v.GetSessionOwner()
		data.User, data.Cart = user.User, cart
	case *cache.MerchantSession:
		data.Merchant.MerchantUser, _ = v.GetSessionOwner()
	}

	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("200 - Valid merchant ID, displaying store"))
	parseProductPage(&w, data)

}

func (a *App) postProd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// Checking for non-empty inputs to be handled by HTML form
	name := r.FormValue("Name")
	ProdDesc := r.FormValue("ProdDesc")
	thumbnail := r.FormValue("Thumbnail")
	price, err := strconv.ParseFloat(r.FormValue("Price"), 64)
	if err != nil {
		log.Fatal(err)
	}
	quantity, err := strconv.Atoi(r.FormValue("Quantity"))
	if err != nil {
		log.Fatal(err)
	}

	merchID := a.merchIDFromSession(w, r)

	if price <= 0 || quantity < 0 || ProdDesc == "" || name == "" {

		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Invalid Price and/or Quantity submitted"))
		return
	}

	p := db.Product{Name: name, ProdDesc: ProdDesc, Thumbnail: thumbnail, Price: price, Quantity: quantity, MerchID: merchID, Sales: 0}
	err = a.cacheManager.CreateProduct(p)
	if err != nil {
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
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Invalid Name submitted"))
		return
	}

	quantity, err := strconv.Atoi(r.URL.Query().Get("Quantity"))
	if r.URL.Query().Get("Quantity") == "" || err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Invalid Quantity submitted"))
		return
	}

	price, err := strconv.ParseFloat(r.URL.Query().Get("Price"), 64)
	if r.URL.Query().Get("Price") == "" || err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Invalid Price submitted"))
		return
	}

	ProdDesc := r.URL.Query().Get("ProdDesc")
	if ProdDesc == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Invalid ProdDesc submitted"))
		return
	}

	merchID := a.merchIDFromSession(w, r)

	thumbnail := r.URL.Query().Get("Thumbnail")
	if thumbnail == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Invalid Thumbnail submitted"))
		return
	}

	p := db.Product{
		Id:        prodID,
		Name:      name,
		Quantity:  quantity,
		Thumbnail: thumbnail,
		Price:     price,
		ProdDesc:  ProdDesc,
		MerchID:   merchID,
	}
	err = a.cacheManager.UpdateProduct(p)
	if err != nil {
		fmt.Println(err)
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
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Status not Found"))
		return
	}

	merchID := a.merchIDFromSession(w, r)

	err := a.cacheManager.DeleteProduct(prodID, merchID)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Unprocessable Entity"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Deleted Successfully"))
}

func (a *App) merchIDFromSession(w http.ResponseWriter, r *http.Request) string {
	sessionCookie, err := r.Cookie("sessionCookie")
	if err != nil {
		return ""
	}
	sessionValStr := sessionCookie.String()
	activeSession, err := a.cacheManager.GetFromCache(sessionValStr, "activeMerchants")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
	}
	accountOwner, _ := activeSession.(*cache.MerchantSession).GetSessionOwner()
	return accountOwner.Id
}

func parseProductPage(w *http.ResponseWriter, data Data) {
	t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/product.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(*w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func parseAllProducts(w *http.ResponseWriter, data Data) {
	t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/allProducts.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(*w, data)
	if err != nil {
		fmt.Println(err)
	}
}
