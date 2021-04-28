package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Tests on Store
// GET Method
func TestStoreGet(t *testing.T) {
	// Passing case: Get all products at a valid store
	req, err := http.NewRequest(http.MethodGet, "/storefront/v0/testingstore/", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusOK, nil, req)

	// Failing case: Get all products at an invalid store
	req, err = http.NewRequest(http.MethodGet, "/storefront/v0/blah/", nil)
	if err != nil {
		t.Errorf(fmt.Sprintf("Request generation error: %s", err))
	}
	checkResponse(t, http.StatusNotFound, nil, req)

	// Passing case: Get some valid product at a valid store
	// Failing case: Get some invalid product at a valid store
	// Failing case: Get some valid product at an invalid store

}

func checkResponse(t *testing.T, targetStatus int, targetPayload interface{}, req *http.Request) {
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != targetStatus {
		t.Errorf(fmt.Sprintf("Expected response code: %s; Got : %s", targetStatus, responseRecorder.Code))
	}
	// TODO: Revisit with exact type of unmarshaled
	var unmarshaled interface{}
	unmarshaled := json.Unmarshal(responseRecorder.Body)
	if unmarshaled != targetPayload {
		t.Errorf(fmt.Sprintf("Expected content: %s; Got : %s", targetPayload, unmarshaled))
	}
}

// DB Queries for merchantID/productID can be joint queries
