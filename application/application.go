package application

import (
	"encoding/json"
	"fmt"
	"net/http"

	calculate "github.com/ArteShow/Calculation_Service/pkg/Calculation"
)

type InputDate struct {
	Expression string `json:"expression"`
}

type SuccessfulAnswer struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	request := new(InputDate)
	defer r.Body.Close()


	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "International Server Error", 500)
		return
	}
	result, err, status := calculate.Calc(request.Expression)
	if err != nil {
		
		NewError := ErrorResponse{
			Error: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(NewError)
		return
	}
	ResultNotJson := SuccessfulAnswer{
		Result: fmt.Sprintf("%f", result), 
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResultNotJson)
}

func RunServer(){
	http.HandleFunc("/", CalcHandler)
	http.ListenAndServe(":8082", nil)
}