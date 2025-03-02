package application

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetExpressionById(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	idStr := parts[4]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		//--------------------------
		http.Error(w, "Empty", http.StatusBadRequest)
		return
	}

	jsonData, _ := json.Marshal(map[string]int{"id": id})
	resp, err := http.Post("http://localhost:8080/internal/expression", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		//-------------------------
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}


func SendExpression(w http.ResponseWriter, r *http.Request){
	body, err := io.ReadAll(r.Body)
	if err != nil{
		//---------------------------------------
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	Data := string(body)
	resp, err := http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer([]byte(Data)))
	if err != nil{
		//-----------------------------------
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		//-------------------------------------
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
	log.Print("A")
}

func GetExpressions(w http.ResponseWriter, r *http.Request){
	resp, err := http.Get("http://localhost:8080/internal/expressions")
	if err != nil{
		//-----------------------------------
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		//----------------------------------
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

//Start the web-service
func RunServer(){
	http.HandleFunc("/api/v1/calculate", SendExpression)
	http.HandleFunc("/api/v1/expressions", GetExpressions)
	http.HandleFunc("/api/v1/expressions/", GetExpressionById)
	http.ListenAndServe(":8082", nil)
}
