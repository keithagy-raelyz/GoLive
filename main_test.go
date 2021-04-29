package main

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

//
func TestMain(t *testing.T) {
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3305)/%s", "user", "password", "store_DB")
	var err error
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	InitializeDB()

	router.HandleFunc("/", home).Methods("GET")

	router.HandleFunc("/merchants/{merchantid}", getMerch).Methods("GET")
	router.HandleFunc("/merchants", postMerch).Methods("POST")
	router.HandleFunc("/merchants/{merchantid}", putMerch).Methods("PUT")
	router.HandleFunc("/merchants/{merchantid}", delMerch).Methods("DELETE")

	router.HandleFunc("/product/{productid}", getProd).Methods("GET")
	router.HandleFunc("/product", postProd).Methods("POST")
	router.HandleFunc("/product/{productid}", putProd).Methods("PUT")
	router.HandleFunc("/product/{productid}", delProd).Methods("DELETE")

	//fmt.Println("Listening at port 5000")
	//log.Fatal(http.ListenAndServe(":5000", router))
}

// Tests on Store
// GET Method
func TestMerchGet(t *testing.T) {
	// Passing case: Get all products at a valid store
	req, err := http.NewRequest(http.MethodGet, "/merchants/1", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	// Failing case: Get all products at an invalid store
	req, err = http.NewRequest(http.MethodGet, "/merchants/300", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusNotFound, nil, req)

	// Passing case: Get some valid product at a valid store
	// Failing case: Get some invalid product at a valid store
	// Failing case: Get some valid product at an invalid store

}

func TestMerchPost(t *testing.T) {
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
func checkResponse(t *testing.T, targetStatus int, targetPayload interface{}, req *http.Request) {
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != targetStatus {
		t.Errorf(fmt.Sprintf("Expected response code: %d; Got : %d", targetStatus, responseRecorder.Code))
	}
	// TODO: Revisit with exact type of unmarshaled
	//var unmarshaled interface{}
	//unmarshaled := json.Unmarshal(responseRecorder.Body)
	//if unmarshaled != targetPayload {
	//	t.Errorf(fmt.Sprintf("Expected content: %s; Got : %s", targetPayload, unmarshaled))
	//}
}

func TestMerchPut(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "/merchants/1?description=hello", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}

	checkResponse(t, http.StatusOK, nil, req)
}

// DB Queries for merchantID/productID can be joint queries

func TestMerchDel(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/merchants/1", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}

	checkResponse(t, http.StatusOK, nil, req)
}
