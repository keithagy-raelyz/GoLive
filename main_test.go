package main_test

import (
	// "database/sql"
	//"encoding/json"
	"fmt"
	"os"

	"GoLive/app"

	// "log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	a *app.App
)

func TestMain(m *testing.M) {
	a = &app.App{}
	a.StartApp()
	//fmt.Println("Listening at port 5000")
	//log.Fatal(http.ListenAndServe(":5000", router))
	code := m.Run()
	os.Exit(code)
}

// Tests on Merchant Methods
// GET Method
func TestAllMerch(t *testing.T) {
	// Passing case: Get all products at a valid store
	req, err := http.NewRequest(http.MethodGet, "/merchants", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	// Passing case: Get some valid product at a valid store
	// Failing case: Get some invalid product at a valid store
	// Failing case: Get some valid product at an invalid store

}

func TestGetMerch(t *testing.T) {
	//Passing test case: Merchant has products
	req, err := http.NewRequest(http.MethodGet, "/merchants/10", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	// Passing test case: Merchant has no products
	req, err = http.NewRequest(http.MethodGet, "/merchants/13", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	//Failing case: Merchant does not exist
	req, err = http.NewRequest(http.MethodGet, "/merchants/300", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusNotFound, nil, req)
}

func TestPostMerch(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/merchants", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	req.ParseForm()
	req.Form.Add("username", "abcd1234")
	req.Form.Add("MerchDesc", "selling drugs")
	req.Form.Add("email", "test@email.com")
	req.Form.Add("pw1", "a")
	req.Form.Add("pw2", "a")

	checkResponse(t, http.StatusCreated, nil, req)
}

func TestPutMerch(t *testing.T) {

	req, err := http.NewRequest(http.MethodPut, "/merchants/1?MerchDesc=hello", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusBadRequest, nil, req)
	req, err = http.NewRequest(http.MethodPut, "/merchants/1?username=abcd1234&email=test@email.com&MerchDesc=hello", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)
}

func TestDelMerch(t *testing.T) {
	// No merchant, No product
	req, err := http.NewRequest(http.MethodDelete, "/merchants/1", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	// Have merchant, Have product
	req, err = http.NewRequest(http.MethodDelete, "/merchants/2", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	// Have merchant, No product
	req, err = http.NewRequest(http.MethodDelete, "/merchants/3", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)
}

// Tests on Product Methods
// GET Method

func TestAllProd(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/products", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)
}

func TestGetProd(t *testing.T) {
	// Product Exists
	req, err := http.NewRequest(http.MethodGet, "/products/6", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	// Product Does Not Exist
	req, err = http.NewRequest(http.MethodGet, "/products/100", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusNotFound, nil, req)
}

func TestPostProd(t *testing.T) {
	// Post product that does not yet exist
	req, err := http.NewRequest(http.MethodPost, "/products", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	req.ParseForm()
	req.Form.Add("Name", "'Dog Biscuits'")
	req.Form.Add("Quantity", "100")
	req.Form.Add("Thumbnail", "https://picsum.photos/200")
	req.Form.Add("Price", "499.99")
	req.Form.Add("ProdDesc", "It's a dog biscuit, dude.")
	req.Form.Add("MerchID", "15")
	checkResponse(t, http.StatusCreated, nil, req)

	// Post product with an invalid MerchID
	// TODO When session management is up, can only post to own MerchID
	req, err = http.NewRequest(http.MethodPost, "/products", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	req.ParseForm()
	req.Form.Add("Name", "Lettuce")
	req.Form.Add("Quantity", "10")
	req.Form.Add("Thumbnail", "https://picsum.photos/200")
	req.Form.Add("Price", "0.05")
	req.Form.Add("ProdDesc", "Nutritionally like cardboard.")
	req.Form.Add("MerchID", "6")
	checkResponse(t, http.StatusUnprocessableEntity, nil, req)

	// Post product with negative price
	req, err = http.NewRequest(http.MethodPost, "/products", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	req.ParseForm()
	req.Form.Add("Name", "Ketchup")
	req.Form.Add("Quantity", "23")
	req.Form.Add("Thumbnail", "https://picsum.photos/200")
	req.Form.Add("Price", "-0.05")
	req.Form.Add("ProdDesc", "Tomatoes.")
	req.Form.Add("MerchID", "17")
	checkResponse(t, http.StatusUnprocessableEntity, nil, req)

	// Post product with negative quantity
	req, err = http.NewRequest(http.MethodPost, "/products", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	req.ParseForm()
	req.Form.Add("Name", "Watch")
	req.Form.Add("Quantity", "-10")
	req.Form.Add("Thumbnail", "https://picsum.photos/200")
	req.Form.Add("Price", "0.05")
	req.Form.Add("ProdDesc", "Tells the time.")
	req.Form.Add("MerchID", "18")
	checkResponse(t, http.StatusUnprocessableEntity, nil, req)

	// Post product with empty ProdDesc
	req, err = http.NewRequest(http.MethodPost, "/products", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	req.ParseForm()
	req.Form.Add("Name", "Nothing")
	req.Form.Add("Quantity", "56789")
	req.Form.Add("Thumbnail", "https://picsum.photos/200")
	req.Form.Add("Price", "0.01")
	req.Form.Add("ProdDesc", "")
	req.Form.Add("MerchID", "19")
	checkResponse(t, http.StatusUnprocessableEntity, nil, req)
}

// func TestPutProd (t *testing.T){
// 	req, err := http.NewRequest(http.MethodPost, "/products", nil)
// 	if err != nil {
// 		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
// 	}
// }

// func TestDelProd (t *testing.T){

// }

func checkResponse(t *testing.T, targetStatus int, targetPayload interface{}, req *http.Request) {
	responseRecorder := httptest.NewRecorder()

	a.TestRoute(responseRecorder, req)
	if responseRecorder.Code != targetStatus {
		t.Errorf(fmt.Sprintf("Expected response code: %d; Got : %d", targetStatus, responseRecorder.Code))
	}
	// TODO: Revisit with exact type of unmarshaled
	// Merchant{a,b,c,d}
	// Response result := Merchant{d,e,f,g}
	//var unmarshaled interface{}
	//unmarshaled := json.Unmarshal(responseRecorder.Body)
	//if unmarshaled != targetPayload {
	//	t.Errorf(fmt.Sprintf("Expected content: %s; Got : %s", targetPayload, unmarshaled))
	//}
}
