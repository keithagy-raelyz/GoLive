package main_test

import (
	// "database/sql"
	//"encoding/json"
	"fmt"
	"io"
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
	// Create a session
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
	req, err := http.NewRequest(http.MethodGet, "/products/5", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	// Product Does Not Exist
	req, err = http.NewRequest(http.MethodGet, "/products/9999999", nil)
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
	req.Form.Add("MerchID", "600")
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

func TestPutProd(t *testing.T) {

	//Updating with all parameters provided
	req, err := http.NewRequest(http.MethodPut, "/products/1?Name=something&Quantity=5&Price=4%2E5&ProdDesc=definitely%20something&MerchID=15&Thumbnail=thumbnail?", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	//Updating with wrong productID
	req, err = http.NewRequest(http.MethodPut, "/products/9000?Name=something&Quantity=5&Price=4%2E0&ProdDesc=definitely%20something&MerchID=15&Thumbnail=thumbnail?", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	//Updating with wrong merchantID
	req, err = http.NewRequest(http.MethodPut, "/products/1?Name=something&Quantity=5&Price=40&ProdDesc=definitely%20something&MerchID=5000&Thumbnail=thumbnail?", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	//Updating with Thumbnail missing
	req, err = http.NewRequest(http.MethodPut, "/products/1?Name=something&Quantity=5&Price=40&ProdDesc=definitely%20something&MerchID=15", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusUnprocessableEntity, nil, req)

	//Updating with MerchID missing
	req, err = http.NewRequest(http.MethodPut, "/products/1?Name=something&Quantity=5&Price=40&ProdDesc=definitely%20something&Thumbnail=thumbnail?", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusUnprocessableEntity, nil, req)

	//Updating with ProdDesc missing
	req, err = http.NewRequest(http.MethodPut, "/products/1?Name=something&Quantity=5&Price=40&MerchID=15&Thumbnail=thumbnail?", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusUnprocessableEntity, nil, req)

	//Updating with price missing
	req, err = http.NewRequest(http.MethodPut, "/products/1?Name=something&Quantity=5&ProdDesc=definitely%20something&MerchID=15&Thumbnail=thumbnail?", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusUnprocessableEntity, nil, req)

	//Updating with quantity missing
	req, err = http.NewRequest(http.MethodPut, "/products/1?Name=something&Price=40&ProdDesc=definitely%20something&MerchID=15&Thumbnail=thumbnail?", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusUnprocessableEntity, nil, req)

	//Updating with name missing
	req, err = http.NewRequest(http.MethodPut, "/products/1?Quantity=5&Price=40&ProdDesc=definitely%20something&MerchID=15&Thumbnail=thumbnail?", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusUnprocessableEntity, nil, req)
}

func TestDelProd(t *testing.T) {

	//Deleting with valid product ID
	req, err := http.NewRequest(http.MethodDelete, "/products/46", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	//Deleting with invalid product iD
	req, err = http.NewRequest(http.MethodDelete, "/products/999999", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)
}

// User Tests
func TestAllUsers(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)
}

func TestGetUser(t *testing.T) {
	//Get valid UserID
	req, err := http.NewRequest(http.MethodGet, "/users/1", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	//Get invalid UserID
	req, err = http.NewRequest(http.MethodGet, "/users/999999", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusNotFound, nil, req)
}

func TestPostUser(t *testing.T) {
	// Post new User
	req, err := http.NewRequest(http.MethodPost, "/users", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	req.ParseForm()
	req.Form.Add("Username", "DogShiet")
	req.Form.Add("Email", "feecalmatter@hotmail.com")
	req.Form.Add("Pw1", "abc")
	req.Form.Add("Pw2", "abc")
	checkResponse(t, http.StatusCreated, nil, req)

	// Post repeated Username
	req, err = http.NewRequest(http.MethodPost, "/users", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	req.ParseForm()
	req.Form.Add("Username", "DogShit")
	req.Form.Add("Email", "fecalmatters@hotmail.com")
	req.Form.Add("Pw1", "abc")
	req.Form.Add("Pw2", "abc")
	checkResponse(t, http.StatusBadRequest, nil, req)

	// Post repeated Email
	req, err = http.NewRequest(http.MethodPost, "/users", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	req.ParseForm()
	req.Form.Add("Username", "DogShitz")
	req.Form.Add("Email", "fecalmatter@hotmail.com")
	req.Form.Add("Pw1", "abc")
	req.Form.Add("Pw2", "abc")
	checkResponse(t, http.StatusBadRequest, nil, req)

	// Post mismatch pw
	req, err = http.NewRequest(http.MethodPost, "/users", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	req.ParseForm()
	req.Form.Add("Username", "DipShitz")
	req.Form.Add("Email", "fecalymatter@hotmail.com")
	req.Form.Add("Pw1", "abcd")
	req.Form.Add("Pw2", "abc")
	checkResponse(t, http.StatusBadRequest, nil, req)
}

func TestPutUser(t *testing.T) {

	//Updating with all parameters provided
	req, err := http.NewRequest(http.MethodPut, "/users/1?Username=testing&Email=testing90316%40hotmail%2Ecom", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	//Updating with no username
	req, err = http.NewRequest(http.MethodPut, "/users/1?Email=testing90316%40hotmail%2Ecom", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusBadRequest, nil, req)

	//Updating with no Email
	req, err = http.NewRequest(http.MethodPut, "/users/1?Username=what", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusBadRequest, nil, req)
}

func TestDelUser(t *testing.T) {

	//Deleting with valid product ID
	req, err := http.NewRequest(http.MethodDelete, "/users/1", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	//Deleting with invalid product iD
	req, err = http.NewRequest(http.MethodDelete, "/users/999999", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)
}

// Session Tests

func checkValidation(t *testing.T)

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

func NewRequestWithCookie(method string, url string, body io.Reader, cookieValid bool) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if cookieValid {
		sessionCookie := &http.Cookie{
			Name:  "sessionCookie",
			Value: "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b", // userID 1
		}
		secondTest := &http.Cookie{
			Name:  "i184m",
			Value: "74dcfbc208bb6aa08c90fb05bda0f2bc53285713e89611dfdd97ae129b5f6195", // dogshit
		}
		req.AddCookie(sessionCookie)
		req.AddCookie(secondTest)
		return req, nil
	}
	//error url
	return req, nil
}
