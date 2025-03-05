package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"log"
)

func TestGenerateID(t *testing.T) {
	// Setup
	handler := http.NewServeMux()
	handler.HandleFunc("/internal/expression", GenerateID)

	// Test 1: Successfully generate ID
	req, err := http.NewRequest("POST", "/internal/expression", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Simulate the request
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Error generating ID, status code: %v", status)
	}

	// Test 2: ID generation is already in progress
	req, err = http.NewRequest("POST", "/internal/expression", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Simulate the request again when ID generation is already in progress
	GenerateIdBool = true
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("Error in ID generation conflict, status code: %v", status)
	}

	// Test 3: ID generation successful after reset
	GenerateIdBool = false
	req, err = http.NewRequest("POST", "/internal/expression", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Simulate the request
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Error generating ID, status code: %v", status)
	}
}

func TestStoreExpression(t *testing.T) {
	// Setup
	handler := http.NewServeMux()
	handler.HandleFunc("/internal/task", StoreExpression)

	// Test 1: Successfully store expression
	expr := Expression{Expression: "2 + 2"}
	body, _ := json.Marshal(expr)
	req, err := http.NewRequest("POST", "/internal/task", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Simulate the request
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Error storing expression, status code: %v", status)
	}
}

func TestSendExpressionById(t *testing.T) {
	// Setup
	handler := http.NewServeMux()
	handler.HandleFunc("/internal/expression/", SendExpressionById)

	// Test 1: Successfully retrieve expression
	// First generate an ID
	GenerateIDTest()

	// Test retrieving expression by ID
	req, err := http.NewRequest("GET", "/internal/expression/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Simulate the request
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error retrieving expression, status code: %v", status)
	}
}

func GenerateIDTest() {
	// Create a new ID
	handler := http.NewServeMux()
	handler.HandleFunc("/internal/expression", GenerateID)

	req, err := http.NewRequest("POST", "/internal/expression", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Simulate the request
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
}

func TestSendExpressionsList(t *testing.T) {
	// Setup
	handler := http.NewServeMux()
	handler.HandleFunc("/internal/expression/list", SendExpressionsList)

	// Test 1: Successfully retrieve list of expressions
	req, err := http.NewRequest("GET", "/internal/expression/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Simulate the request
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error retrieving expressions list, status code: %v", status)
	}
}
