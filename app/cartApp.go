package app

import (
	"GoLive/cache"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *App) getCart(w http.ResponseWriter, r *http.Request) {
	//TODO
	activeSession, ok := a.HaveValidSessionCookie(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u, c := activeSession.(*cache.UserSession).GetSessionOwner()
	data := Data{
		User: u,
		Cart: c,
	}
	parseCartPage(&w, data)
}

func (a *App) postCart(w http.ResponseWriter, r *http.Request) {
	//Obtain user session Data, redirect if invalid
	activeSession, ok := a.HaveValidSessionCookie(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//Obtain item Data
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	id := r.FormValue("Id")

	a.cacheManager.UpdateCart(activeSession, id, "append")
	a.cacheManager.BlockStock(id)

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func (a *App) updateCart(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	productID := params["productid"]

	activeSession, ok := a.HaveValidSessionCookie(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	a.cacheManager.UpdateCart(activeSession, productID, "+")
	a.cacheManager.BlockStock(productID)

	jData, _ := json.Marshal(Response{true})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func (a *App) deleteCart(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productID := params["productid"]
	activeSession, ok := a.HaveValidSessionCookie(r)
	if !ok {
		fmt.Println("session is not valid")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//Delete one from cart and update in the cache

	a.cacheManager.UpdateCart(activeSession, productID, "-")
	a.cacheManager.ReleaseStock(productID)

	jData, err := json.Marshal(Response{true})
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

type Response struct {
	Redirect bool `json:"redirect"`
}

func parseCartPage(w *http.ResponseWriter, data Data) {
	t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/viewCart.html", "templates/error.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(*w, data)
	if err != nil {
		log.Fatal(err)
	}
}
