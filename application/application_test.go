package application

import (
	"encoding/json"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)



func TestCalcHandler(t *testing.T){
	requestBody := `{"expression": "2/0"}`
	expected := `{"error":Division by zero}%!(EXTRA int=200)`
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(requestBody)))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	CalcHandler(rr, req)
	if status := rr.Code; status != http.StatusOK{
		t.Errorf("Wrong Code. Want: %v, got %v", http.StatusOK, status)
	}
	var actual ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&actual) // JSON-Body in Struct dekodieren
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	if actual.Error != expected{
		t.Errorf("Wrong Anwser. Wanted %v, got %v", expected, actual.Error)
	}
	t.Log(actual.Error)
}