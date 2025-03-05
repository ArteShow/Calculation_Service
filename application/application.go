package application

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetExpressionById(w http.ResponseWriter, r *http.Request){
	log.Println("📩 Anfrage auf:", r.URL.Path)

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/expression/"), "/")
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
		http.Error(w, "Empty", http.StatusBadRequest)
		return
	}

	url := "http://localhost:8083/internal/expression/" + strconv.Itoa(id)
	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Fehler beim Abrufen der Expression:", err)
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Fehler beim Lesen der Antwort:", err)
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Antwort von internal:", string(body))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func GetExpressionsList(w http.ResponseWriter, r *http.Request){
	resp, err := http.Get("http://localhost:8083/internal/expression/list")
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Fehler beim Lesen der Antwort:", err)
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Fehler beim Abrufen der Expression:", err)
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Antwort von internal:", string(body))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// Sendet eine Expression an den Internal-Server (8083) und fordert danach eine ID an
func SendExpression(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Fehler beim Lesen des Bodys:", err)
		http.Error(w, "1", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	log.Println(body)
	log.Println("📡 Sende Expression an internal:", string(body))
	_, err2 := http.Post("http://localhost:8083/internal/task", "application/json", bytes.NewBuffer(body))

	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Fehler beim Senden:", err)
		http.Error(w, "2", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Expression gespeichert, fordere jetzt eine ID an...")

	// Fordere eine neue ID an
	idResp, err := http.Post("http://localhost:8083/internal/expression", "application/json", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Fehler beim Anfordern der ID:", err)
		http.Error(w, "4", http.StatusInternalServerError)
		return
	}
	defer idResp.Body.Close()

	idBody, err := io.ReadAll(idResp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Fehler beim Lesen der ID-Antwort:", err)
		http.Error(w, "5", http.StatusInternalServerError)
		return
	}

	// Debugging - Zeige den Inhalt der ID-Antwort an
	log.Printf("📜 Erhaltene ID-Antwort: %s\n", string(idBody))

	if len(idBody) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Fehler: Keine ID erhalten")
		http.Error(w, "6", http.StatusInternalServerError)
		return
	}

	log.Println("✅ ID erhalten:", string(idBody))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(idBody)
}

func RunServer() {
	log.Println("🌍 API-Server gestartet auf Port 8082")
	http.HandleFunc("/api/v1/calculate", SendExpression)
	http.HandleFunc("/api/v1/expression/", GetExpressionById)
	http.HandleFunc("/api/v1/expressions", GetExpressionsList)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
