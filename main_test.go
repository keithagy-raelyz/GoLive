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

//Seed Data and INitialize the DB values
func InitializeDB() {
	Query1 := `CREATE TABLE IF NOT EXISTS Users (
    UserID int NOT NULL AUTO_INCREMENT,
    Username VARCHAR(255) NOT NULL,
    Password VARCHAR(255) NOT NULL,
    Email varchar(255) NOT NULL,
    PRIMARY KEY (UserID)
	)`
	_, err := db.Exec(Query1)
	if err != nil {
		log.Fatal(err)
	}
	Query2 := `
	CREATE TABLE IF NOT EXISTS Merchants (
		MerchantID int NOT NULL AUTO_INCREMENT,
		Username VARCHAR(255) NOT NULL,
		Password VARCHAR(255) NOT NULL,
		Email varchar(255) NOT NULL,
		Description VARCHAR(255),
		PRIMARY KEY (MerchantID)
	);`
	_, err = db.Exec(Query2)
	if err != nil {
		log.Fatal(err)
	}
	Query3 := `CREATE TABLE IF NOT EXISTS Products (
    ProductID int NOT NULL AUTO_INCREMENT,
    Product_Name VARCHAR(255) NOT NULL,
    Quantity int NOT NULL,
    Image varchar(255) NOT NULL,
    Price int not null,
    Description VARCHAR(255),
    MerchantID int NOT NULL,
    Foreign Key (MerchantID) REFERENCES Merchants (MerchantID),
    PRIMARY KEY (ProductID)
	);`
	_, err = db.Exec(Query3)
	if err != nil {
		log.Fatal(err)
	}
	Query4 := `	INSERT INTO USERS (username,password,email) VALUES("A","B","C");`
	_, err = db.Exec(Query4)
	if err != nil {
		log.Fatal(err)
	}
	Query5 := `	INSERT INTO MERCHANTS (username,password,email) VALUES("A","B","C");`
	_, err = db.Exec(Query5)
	if err != nil {
		log.Fatal(err)
	}
	Query6 := `	INSERT INTO MERCHANTS (username,password,email,description) VALUES("A","B","C","D");`
	_, err = db.Exec(Query6)
	if err != nil {
		log.Fatal(err)
	}
	Query7 := `	INSERT INTO PRODUCTS (product_name,quantity,image,price,description,merchantid) VALUES("A",2,"B",20.5,"C",1);`
	_, err = db.Exec(Query7)
	if err != nil {
		log.Fatal(err)
	}
	Query8 := `	INSERT INTO PRODUCTS (product_name,quantity,image,price,description,merchantid) VALUES("D",0,"B",30.1,"C",8);`
	_, err = db.Exec(Query8)
	if err != nil {
		log.Fatal(err)
	}
}
