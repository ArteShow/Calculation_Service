package application

import (
	"encoding/json"
	"fmt"
	"net/http"

	calculate "github.com/ArteShow/Calculation_Service/pkg/Calculation"
)

//Create some structions for possible answers
type InputDate struct {
	Expression string `json:"expression"`
}

type SuccessfulAnswer struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

//Handler for calculations
func CalcHandler(w http.ResponseWriter, r *http.Request) {
	//Get the expression
	request := new(InputDate)
	defer r.Body.Close()


	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	//Calculate the expression
	result, err, status := calculate.Calc(request.Expression)
	if err != nil {
		//Handel the error possibility
		NewError := ErrorResponse{
			Error: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(NewError)
		return
	}
	//Send the right anwser back
	ResultNotJson := SuccessfulAnswer{
		Result: fmt.Sprintf("%f", result), 
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResultNotJson)
}
//Start the web-service
func RunServer(){
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	http.ListenAndServe(":8082", nil)
}
