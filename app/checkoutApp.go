package app

import (
	"GoLive/cache"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/checkout/session"
)

type createCheckoutSessionResponse struct {
	SessionID string `json:"id"`
}

//type CartItem struct {
//	Product products `json:"product"`
//	Count   int      `json:"count,omitempty"`
//}
//
//type products struct {
//	Id        string  `json:"id,omitempty"`
//	Name      string  `json:"name,omitempty"`
//	Quantity  int     `json:"quantity,omitempty"`
//	Thumbnail string  `json:"thumbnail,omitempty"`
//	Price     float64 `json:"price,omitempty"`
//	ProdDesc  string  `json:"prod_desc,omitempty"`
//	MerchID   string  `json:"merch_id,omitempty"`
//	Sales     int     `json:"sales,omitempty"`
//}

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
	jsonCart, err := json.Marshal(c.Contents())
	if err != nil {
		fmt.Println("cart error")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	fmt.Println("Cart:", c)
	fmt.Println("jsonCart:", string(jsonCart))

	data := Data{Cart: c,
		JSON: string(jsonCart)}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

type CheckOutCart []cache.CartItem

func (a *App) payment(w http.ResponseWriter, r *http.Request) {

	rawData, _ := ioutil.ReadAll(r.Body)
	var cart CheckOutCart
	json.Unmarshal(rawData, &cart)
	fmt.Println(cart, "cart at payment")
	sessionCookie, err := r.Cookie("sessionCookie")
	if err != nil {
		fmt.Println("session is not valid")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userID := sessionCookie.Value
	a.cacheManager.AddCartProcessing(userID)

	domain := "http://localhost:5000"

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems:  []*stripe.CheckoutSessionLineItemParams{},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "/success" + "?userID=" + userID),
		CancelURL:  stripe.String(domain + "/cancel" + "?userID=" + userID),
	}

	for _, item := range cart {
		params.LineItems = append(params.LineItems,
			&stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyUSD)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(item.Product.Name),
					},
					UnitAmount: stripe.Int64(int64(item.Product.Price * 100)),
				},
				Quantity: stripe.Int64(int64(item.Count)),
			},
		)
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

func (a *App) paymentSuccess(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	fmt.Println(userID)
	a.cacheManager.ClearActiveUserCart(userID)

	activeSession, ok := a.HaveValidSessionCookie(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u, _ := activeSession.(*cache.UserSession).GetSessionOwner()
	data := Data{
		User: u.User,
	}
	//fmt.Println("payment success!")
	//http.Redirect(w, r, "/", http.StatusSeeOther)
	// Execute Template
	t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/success.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (a *App) paymentCancelled(w http.ResponseWriter, r *http.Request) {
	cartID := r.URL.Query().Get("cartID")
	a.cacheManager.CartFailure(cartID)

	fmt.Println("payment cancelled!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
	// Execute Template
}
