package application

import (
	"bytes"
	"net/http/httptest"
	"testing"
)
func TestCalcHandler(t *testing.T){
	requestBody := `{"expression": "-5"}`
	expected := `{"error":"There is a letter"}`
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(requestBody)))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	CalcHandler(rr, req)

	if rr.Body.String() != expected {
		t.Errorf("Wrong Anwser. Wanted %v, got %v . Status: %v.", expected, rr.Body.String(), rr.Code)
	}
}