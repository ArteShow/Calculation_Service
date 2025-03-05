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
	log.Println("📩 Request to:", r.URL.Path)

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/expression/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		w.WriteHeader(http.StatusNotFound)
		log.Println("❌ Error: Invalid URL")
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("❌ Error: ID is not a number:", parts[0])
		http.Error(w, "Empty", http.StatusBadRequest)
		return
	}

	url := "http://localhost:8083/internal/expression/" + strconv.Itoa(id)
	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error while retrieving expression:", err)
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error while reading response:", err)
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Response from internal:", string(body))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func GetExpressionsList(w http.ResponseWriter, r *http.Request){
	resp, err := http.Get("http://localhost:8083/internal/expression/list")
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error while reading response:", err)
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error while retrieving expressions:", err)
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Response from internal:", string(body))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// Sends an expression to the internal server (8083) and requests an ID afterward
func SendExpression(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error while reading the body:", err)
		http.Error(w, "1", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	log.Println(body)
	log.Println("📡 Sending expression to internal:", string(body))
	_, err2 := http.Post("http://localhost:8083/internal/task", "application/json", bytes.NewBuffer(body))

	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error while sending:", err)
		http.Error(w, "2", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Expression saved, requesting an ID now...")

	// Request a new ID
	idResp, err := http.Post("http://localhost:8083/internal/expression", "application/json", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error while requesting ID:", err)
		http.Error(w, "4", http.StatusInternalServerError)
		return
	}
	defer idResp.Body.Close()

	idBody, err := io.ReadAll(idResp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error while reading ID response:", err)
		http.Error(w, "5", http.StatusInternalServerError)
		return
	}

	// Debugging - Show the content of the ID response
	log.Printf("📜 Received ID response: %s\n", string(idBody))

	if len(idBody) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error: No ID received")
		http.Error(w, "6", http.StatusInternalServerError)
		return
	}

	log.Println("✅ ID received:", string(idBody))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(idBody)
}

func RunServer() {
	log.Println("🌍 API server started on port 8082")
	http.HandleFunc("/api/v1/calculate", SendExpression)
	http.HandleFunc("/api/v1/expression/", GetExpressionById)
	http.HandleFunc("/api/v1/expressions", GetExpressionsList)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
