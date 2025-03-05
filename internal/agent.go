package internal

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"
)

var (
	id           int
	Calculations = map[string]map[string]Expressions{
		"expression": {},
	}
	mu sync.Mutex
)

type Expression struct {
	Expression string `json:"expression"`
}

type Id struct {
	Id int `json:"id"`
}

type Expressions struct {
	ID     int     `json:"id"`
	Status int     `json:"status"`
	Result float64 `json:"result"`
}

func SendExpressionById(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:8080/internal/expression")
	if err != nil {
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}

	var data Id
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}

	idStr := strconv.Itoa(data.Id)

	mu.Lock()
	expression, exists := Calculations["expression"][idStr]
	mu.Unlock()

	if !exists {
		http.Error(w, "Empty", http.StatusNotFound)
		return
	}

	jsonData, _ := json.Marshal(expression)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func SendToClientExpressions(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	jsonData, err := json.Marshal(Calculations["expression"])
	mu.Unlock()

	if err != nil {
		http.Error(w, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GetExpression(w http.ResponseWriter, r *http.Request) Expression {
	response, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		http.Error(w, "Empty", http.StatusInternalServerError)
		return Expression{}
	}
	defer response.Body.Close()

	var returning Expression
	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Empty", http.StatusInternalServerError)
		return Expression{}
	}
	err = json.Unmarshal(body, &returning)
	if err != nil {
		http.Error(w, "Empty", http.StatusInternalServerError)
		return Expression{}
	}
	return returning
}

func GenerateID(w http.ResponseWriter, r *http.Request) {
	expression := GetExpression(w, r)

	mu.Lock()
	id++
	newID := id
	idStr := strconv.Itoa(newID)

	Calculations["expression"][idStr] = Expressions{
		ID:     newID,
		Status: 0,
		Result: 0,
	}
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Id{Id: newID})

	Calculate(expression)
}

func Calculate(expression Expression) {
	// Berechnung hinzufügen
}

func RunServerAgent() {
	http.HandleFunc("/internal/task", GenerateID)
	http.HandleFunc("/internal/task/expression", SendExpressionById)
	http.HandleFunc("/internal/task/expressions", SendToClientExpressions)
	http.ListenAndServe(":8080", nil)
}
