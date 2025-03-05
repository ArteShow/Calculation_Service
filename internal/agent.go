package internal

import (
	"encoding/json"
	"io"
	"context"
	"net/http"
	"log"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"sync"

	calculate "github.com/ArteShow/Calculation_Service/pkg/Calculation"
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



func GenerateID(w http.ResponseWriter, r *http.Request) {
	if GenerateIdBool {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Println("⚠️ Anfrage ignoriert: ID-Generierung läuft bereits")
		http.Error(w, "ID-Generierung läuft bereits", http.StatusConflict)
		return
	}

	GenerateIdBool = true
	defer func() { GenerateIdBool = false }()

	log.Println("🔢 Generiere neue ID...")

	newID := atomic.AddInt64(&id, 1)
	Calculations["expression"][expression] = Expressions{
		ID:     int(newID),
		Status: 0,
		Result: 0,
		Error: nil,
	}
	ExpressionByID[int(newID)] = expression
	log.Println("✅ Neue ID generiert:", newID)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func SendExpressionById(w http.ResponseWriter, r *http.Request){
	 // ID aus der URL extrahieren
    parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/internal/expression/"), "/")
    if len(parts) == 0 || parts[0] == "" {
		w.WriteHeader(http.StatusNotFound)
        log.Println("❌ Fehler: Ungültige URL")
        http.Error(w, "Invalid URL format", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(parts[0])
    if err != nil {
		w.WriteHeader(http.StatusNotFound)
        log.Println("❌ Fehler: ID ist keine Zahl:", parts[0])
        http.Error(w, "Invalid ID format", http.StatusBadRequest)
        return
    }

    // Suche die Expression basierend auf der ID
    exprStr, found := ExpressionByID[id]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		log.Println("❌ Fehler: Expression nicht gefunden")
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	expression, found := Calculations["expression"][exprStr]

    // JSON Antwort senden
    response, err := json.Marshal(expression)
    if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
        log.Println("❌ Fehler beim Erstellen der Antwort:", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
	
    w.Header().Set("Content-Type", "application/json")
    w.Write(response)
	w.WriteHeader(http.StatusOK)
}

func GetExpression(w http.ResponseWriter, r *http.Request) Expression {
	response, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Fehler beim Lesen der Expression:", err)
		http.Error(w, "3", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var returning Expression
	body, err := io.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Fehler beim Unmarshalen:", err)
		log.Println("📜 Erhaltene Expression:", string(body))
		http.Error(w, "4", http.StatusBadRequest)
		return
	}

	log.Println("✅ Expression erhalten:", expr.Expression)

	w.WriteHeader(http.StatusCreated)
	expression = expr.Expression
}


func Calculate(expression string) {
	log.Println("I am in the Calculate function! Juhu! 😄")
	//////////////////////
	TIME_ADDITION_MS = 5 * time.Millisecond
	TIME_SUBTRACTION_MS = 1 * time.Millisecond
	TIME_MULTIPLICATIONS_MS = 5 * time.Millisecond
	TIME_MULTIPLICATIONS_MS = 5 * time.Millisecond

	//////////////////////
	var wg sync.WaitGroup
	var mu sync.Mutex
	var statusCode int
	var finalResult float64
	var finalError error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parts := strings.Fields(expression) 

	wg.Add(len(parts))
	for _, expr := range parts {
		go func(expr string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				log.Println("🛑 Timeout: Berechnung abgebrochen")
				return
			default:
				result, err, code := calculate.Calc(expr)
				mu.Lock()
				if err != nil {
					log.Println("❌ Fehler bei der Berechnung:", err)
					finalError = err
				} else {
					log.Printf("✅ Berechnung: %s = %f, StatusCode: %d\n", expr, result, code)
					finalResult += result
					statusCode = code
				}
				mu.Unlock()
			}
		}(expr)
	}

	wg.Wait()

	log.Println("✅ Alle Berechnungen abgeschlossen!")
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
