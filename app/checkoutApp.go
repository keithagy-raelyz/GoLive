package app

import (
	"GoLive/cache"
	"GoLive/db"
	"html/template"
	"log"
	"net/http"
)

func (a *App) checkOutPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/checkout.html", "templates/error.html")
	if err != nil {
		log.Fatal(err)
	}
	cartItem := cache.CartItem{
		Product: db.Product{
			Id:        "5",
			Name:      "Test",
			Quantity:  5,
			Thumbnail: "test",
			Price:     30,
			ProdDesc:  "fucker",
			MerchID:   "10",
			Sales:     0},
		Count: 5}
	cart := []cache.CartItem{cartItem}
	data := Data{Cart: cart}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) payment(w http.ResponseWriter, r *http.Request) {

}
