package application

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetExpressionById(t *testing.T) {
	// Setup
	handler := http.NewServeMux()
	handler.HandleFunc("/api/v1/expression/", GetExpressionById)

	// Test 1: Successfully retrieve expression by ID
	// First, generate an ID
	// GenerateIDTest() // Uncomment if you have a function to generate IDs

	// Test retrieving the expression by ID
	req, err := http.NewRequest("GET", "/api/v1/expression/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Simulate the request
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Failed to retrieve expression, StatusCode: %v", status)
	}
}

func TestGetExpressionsList(t *testing.T) {
	// Setup
	handler := http.NewServeMux()
	handler.HandleFunc("/api/v1/expressions", GetExpressionsList)

	// Test 1: Successfully retrieve the expressions list
	req, err := http.NewRequest("GET", "/api/v1/expressions", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Simulate the request
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Failed to retrieve expressions list, StatusCode: %v", status)
	}
}

type Expression struct {
	Expression string `json:"expression"`
}

func TestSendExpression(t *testing.T) {
	// Setup
	handler := http.NewServeMux()
	handler.HandleFunc("/api/v1/calculate", SendExpression)

	// Test 1: Successfully send an expression
	expr := Expression{Expression: "2 + 2"}
	body, _ := json.Marshal(expr)
	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Simulate the request
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Failed to send expression, StatusCode: %v", status)
	}
}
