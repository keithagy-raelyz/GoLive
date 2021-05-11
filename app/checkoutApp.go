package app

import (
	"GoLive/cache"
	"encoding/json"
	"fmt"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/checkout/session"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

func (a *App) checkOutPage(w http.ResponseWriter, r *http.Request) {
	activeSession, ok := a.HaveValidSessionCookie(r)
	if !ok {
		fmt.Println("session is not valid")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/checkout.html", "templates/error.html")
	if err != nil {
		log.Fatal(err)
	}
	_, c := activeSession.(*cache.UserSession).GetSessionOwner()

	data := Data{Cart: c}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

type CartItem struct {
	Product products `json:"product"`
	Count   int      `json:"count,omitempty"`
}

type products struct {
	Id        string  `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Quantity  int     `json:"quantity,omitempty"`
	Thumbnail string  `json:"thumbnail,omitempty"`
	Price     float64 `json:"price,omitempty"`
	ProdDesc  string  `json:"prod_desc,omitempty"`
	MerchID   string  `json:"merch_id,omitempty"`
	Sales     int     `json:"sales,omitempty"`
}

func (a *App) payment(w http.ResponseWriter, r *http.Request) {

	rr, _ := ioutil.ReadAll(r.Body)
	var cart []CartItem
	json.Unmarshal(rr, &cart)
	for _, item := range cart {
		fmt.Println(item)
	}

	domain := "http://localhost:5000"

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencySGD)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("T-shirt"),
					},
					UnitAmount: stripe.Int64(2000),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "/success.html"),
		CancelURL:  stripe.String(domain + "/cancel.html"),
	}

	session, err := session.New(params)

	if err != nil {
		log.Printf("session.New: %v", err)
	}

	data := createCheckoutSessionResponse{
		SessionID: session.ID,
	}

	js, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")

	w.Write(js)
}

type createCheckoutSessionResponse struct {
	SessionID string `json:"id"`
}
