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
	//TODO
	activeSession, ok := a.HaveValidSessionCookie(r)
	if !ok {
		fmt.Println("session is not valid")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u, c := activeSession.(*cache.UserSession).GetSessionOwner()
	data := Data{
		User: u,
		Cart: c,
	}
	t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/viewCart.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}

}

func (a *App) postCart(w http.ResponseWriter, r *http.Request) {
	//destructure the data from the form into a product
	//if we loaded the product page? did we cache it?
	//if we destructure, then where are we actually updating it? does this mean we keep a cached copy of the item before
	//updating the db?
	//every single user's cart qty vs DB amount
	//are we even caching pages at the moment? NO
	//we only cache when someone adds to cart
	//map[productID]product on expiry, update the DB IF changes have been made changes boolean

	//Obtain user session Data, redirect if invalid
	activeSession, ok := a.HaveValidSessionCookie(r)
	if !ok {
		fmt.Println("session is not valid")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//u, _ := activeSession.(*cache.UserSession).GetSessionOwner()

	//Obtain item Data
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	id := r.FormValue("Id")

	a.cacheManager.UpdateCart(activeSession, id, "append")

	http.Redirect(w, r, "/cart", http.StatusSeeOther)

	//SCENARIO A (EVERYONE CAN ADD TO CART, DB TRANSACTION WILL VERIFY IF IT CNA GO THROUGH)
	//eg 10 ppl have varying amounts of shit in their cart, lets say 30x shit
	// DB only has 10 shit, should this be allowed?
	//When does the verification / deletion etc should be carried out.

	//SCENARIO B (CACHED variable is king)
	//5 ppl have added 30 shit into their cart
	//db has 40 shit
	//6th person assessed the shit page, should the person see 10? or 40?
	//
	//case SEe 10
	//usercheckout experience is guaranteed.
	//
	//expired Cart / Session
	//needs to update the cached page.

	//Cache value as safeguard
	//

	//1)  Add to cache only after it has been added to a cart
	//2)  Add to cache the moment it has been read GET req

	//individual A just loaded apples qty 37
	//individual B checked out 30 apples
	//individual A  tries to add 30 apples to his cart.

	// in Scenario 2, in cache memory will reflect 7 and reject the request.
	// in scenario 1, in cache memory will reflect 7 too? (what if the cache memory expired already)

}

func (a *App) updateCart(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	productID := params["productid"]

	activeSession, ok := a.HaveValidSessionCookie(r)
	if !ok {
		fmt.Println("session is not valid")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//u, c := activeSession.(*cache.UserSession).GetSessionOwner()
	//for _,cartitem:= range c{
	//	caritem
	//}
	//Add to cart and udpate in the session

	a.cacheManager.UpdateCart(activeSession, productID, "+")
	jData, _ := json.Marshal(Response{true})

	//data := Data{
	//	User: u,
	//	Cart: c,
	//}
	//if err != nil {
	//	fmt.Println(err)
	//}
	//t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/viewCart.html")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = t.Execute(w, data)
	//if err != nil {
	//	fmt.Println(err)
	//}
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
	//u, c := activeSession.(*cache.UserSession).GetSessionOwner()
	//data := Data{
	//	User: u,
	//	Cart: c,
	//}
	////Delete one from carrt and update in the cache
	//
	a.cacheManager.UpdateCart(activeSession, productID, "-")
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
