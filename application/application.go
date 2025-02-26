package application

import (
	"bytes"
	"io"
	"net/http"
	//calculate "github.com/ArteShow/Calculation_Service/pkg/Calculation"
)


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


//Start the web-service
func RunServer(){
	http.HandleFunc("/api/v1/calculate", SendExpression)
	http.ListenAndServe(":8080", nil)
}
