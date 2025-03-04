package internal

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	calculate "github.com/ArteShow/Calculation_Service/pkg/Calculation"
)

var (
	id            int64
	Calculations  = map[string]map[string]Expressions{"expression": {}}
	GenerateIdBool bool
	expression string
	ExpressionByID     = map[int]string{}
	TIME_ADDITION_MS time.Duration
	TIME_SUBTRACTION_MS time.Duration
	TIME_MULTIPLICATIONS_MS time.Duration
	TIME_DIVISIONS_MS time.Duration
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
	Error error
}

func SendExpressionsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(Calculations)
	if err != nil {
		log.Println("❌ Fehler beim Erstellen der Antwort:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}



func GenerateID(w http.ResponseWriter, r *http.Request) {
	if GenerateIdBool {
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Id{Id: int(newID)})
	Calculate(expression)
}

func SendExpressionById(w http.ResponseWriter, r *http.Request){
	 // ID aus der URL extrahieren
    parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/internal/expression/"), "/")
    if len(parts) == 0 || parts[0] == "" {
        log.Println("❌ Fehler: Ungültige URL")
        http.Error(w, "Invalid URL format", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(parts[0])
    if err != nil {
        log.Println("❌ Fehler: ID ist keine Zahl:", parts[0])
        http.Error(w, "Invalid ID format", http.StatusBadRequest)
        return
    }

    // Suche die Expression basierend auf der ID
    exprStr, found := ExpressionByID[id]
	if !found {
		log.Println("❌ Fehler: Expression nicht gefunden")
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	expression, found := Calculations["expression"][exprStr]

    // JSON Antwort senden
    response, err := json.Marshal(expression)
    if err != nil {
        log.Println("❌ Fehler beim Erstellen der Antwort:", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
	
    w.Header().Set("Content-Type", "application/json")
    w.Write(response)
}

func StoreExpression(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("❌ Fehler beim Lesen der Expression:", err)
		http.Error(w, "3", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var expr Expression
	err = json.Unmarshal(body, &expr)
	if err != nil {
		log.Println(body)
		log.Println("❌ Fehler beim Unmarshalen:", err)
		log.Println("📜 Erhaltene Expression:", string(body))
		http.Error(w, "4", http.StatusBadRequest)
		return
	}

	log.Println("✅ Expression erhalten:", expr.Expression)

	w.WriteHeader(http.StatusOK)
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
	Calculations["expression"][expression] = Expressions{
		ID:     Calculations["expression"][expression].ID,
		Status: statusCode,
		Result: finalResult,
		Error:  finalError,
	}
	mu.Unlock()

	log.Println("Gesamtergebnis:", finalResult)
}

func RunServerAgent() {
	log.Println("🚀 Internal-Server gestartet auf Port 8083")
	http.HandleFunc("/internal/task", StoreExpression)
	http.HandleFunc("/internal/expression/", SendExpressionById)
	http.HandleFunc("/internal/expression", GenerateID)
	http.HandleFunc("/internal/expression/list", SendExpressionsList)  // Richtig: list-Route für Liste

	log.Fatal(http.ListenAndServe(":8083", nil))
}

