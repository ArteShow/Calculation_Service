package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	//calculate "github.com/ArteShow/Calculation_Service/pkg/Calculation"
)

type AgentAnwser struct{
	result string `json:"result"`
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
	resp, err := http.Post("http://localhost:8080/internal", "application/json", bytes.NewBuffer([]byte(Data)))
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

	w.Write(responseBody)
}

func FromServerToClient(w http.ResponseWriter, resp *http.Response){
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		//---------------------------------
		http.Error(w, "Empty", http.StatusInternalServerError)
		return
	}
	Result := AgentAnwser{}
	err = json.Unmarshal(responseBody, &Result)
	if err != nil{
		//----------------------------
		http.Error(w, "Empty", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "Result: %s", Result.result)
}

//Start the web-service
func RunServer(){
	http.HandleFunc("/api/v1/calculate", SendExpression)
	http.ListenAndServe(":8080", nil)
}
