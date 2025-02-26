package internal

import (
	"encoding/json"
	"io"
	"net/http"
	//"time"
)

var(
	id int
)

type Expression struct{
	Expression string `json:"expression"`
}

type Id struct{
	id int `json:id`
}

func GetExpression(w http.ResponseWriter, r *http.Request) Expression{
	response, err := http.Get("http://localhost:8080/internal")
	if err != nil{
		//-----------------
		http.Error(w, "Empty", http.StatusInternalServerError)
		return Expression{}
	}
	defer response.Body.Close()
	var returning Expression
	body, err := io.ReadAll(response.Body)
	if err != nil{
		//-----------------
		http.Error(w, "Emtpy", http.StatusInternalServerError)
		return Expression{}
	}
	err2 := json.Unmarshal(body, &returning)
	if err2 != nil{
		//--------------------
		http.Error(w, "Empty", http.StatusInternalServerError)
		return Expression{}
	}
	return returning
}

func GenerateID(w http.ResponseWriter, r *http.Request){
	expression := GetExpression(w, r)
	NewId := Id{id: id+1}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewId)
	Calculate(expression)

}

func Calculate(expression Expression){
	o := 1
}

func main(){
	http.HandleFunc("/internal", GenerateID)
	http.ListenAndServe(":8082", nil)
}