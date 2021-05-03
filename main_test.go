package main_test

import (
	// "database/sql"
	//"encoding/json"
	"fmt"
	"os"

	"github.com/keithagy-raelyz/GoLive/app"

	// "log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	a *app.App
)

//
func TestMain(m *testing.M) {
	a = &app.App{}
	a.StartApp()
	//fmt.Println("Listening at port 5000")
	//log.Fatal(http.ListenAndServe(":5000", router))
	code := m.Run()
	os.Exit(code)
}

// Tests on Store
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
	req, err := http.NewRequest(http.MethodGet, "/merchants/2", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	// Passing test case: Merchant has no products
	req, err = http.NewRequest(http.MethodGet, "/merchants/1", nil)
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
	req.Form.Add("description", "selling drugs")
	req.Form.Add("email", "test@email.com")
	req.Form.Add("pw1", "a")
	req.Form.Add("pw2", "a")

	checkResponse(t, http.StatusCreated, nil, req)
}

func TestPutMerch(t *testing.T) {

	req, err := http.NewRequest(http.MethodPut, "/merchants/1?description=hello", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusBadRequest, nil, req)
	req, err = http.NewRequest(http.MethodPut, "/merchants/1?username=user?email=xxx@hotmail.com?description=hello", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)
}

// DB Queries for merchantID/productID can be joint queries

func TestDelMerch(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/merchants/1", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)
}

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
