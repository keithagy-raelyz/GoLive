package app

import (
	"GoLive/cache"
	"GoLive/db"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func (a *App) getCart(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/viewCart.html", "templates/error.html")
	if err != nil {
		log.Fatal(err)
	}
	cartItem := cache.CartItem{db.Product{
		Id:        "5",
		Name:      "Test",
		Quantity:  5,
		Thumbnail: "test",
		Price:     30,
		ProdDesc:  "fucker",
		MerchID:   "10",
		Sales:     0,
	}, 5}
	cart := cache.Cart{cartItem}
	data := Data{Cart: cart}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
	return

	//TODO

	//code requesting for userdata from cookies

	//getcartdata from cache using the userID from cookie
	if userSession, ok := a.cacheManager.GetFromCache("sessionid", "activeUsers"); !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		user, cart := userSession.(*cache.UserSession).Data()
		data := Data{User: user, Cart: cart}

		t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/viewCart.html", "templates/error.html")
		if err != nil {
			log.Fatal(err)
		}
		err = t.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
}

func (a *App) postCart(w http.ResponseWriter, r *http.Request) {

}

func (a *App) updateCart(w http.ResponseWriter, r *http.Request) {
	jData, err := json.Marshal(Response{true})
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func (a *App) deleteCart(w http.ResponseWriter, r *http.Request) {

}

type Response struct {
	Redirect bool `json:"redirect"`
}
