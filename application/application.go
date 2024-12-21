package application

import (
	"encoding/json"
	"fmt"
	"net/http"

	calculate "github.com/ArteShow/Calculation_Service/pkg/Calculation"
)

type InputDate struct{
	Expression string `json:"expression"`
}

type SuccessfulAnwser struct{
	Result float64 `json:"result"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request){
	request := new(InputDate)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil{
		http.Error(w, "International server error", 500)
	}
	result, err, status := calculate.Calc(request.Expression)
	if err != nil{
		NewError := ErrorResponse{
			Error: err.Error(),
		}
		jsonData, err := json.Marshal(NewError)
		if err != nil{
			http.Error(w, "International server error", 500)
		}
		fmt.Fprint(w, jsonData)
		fmt.Fprintln(w, jsonData)
	}else{
		ResultNotJson := SuccessfulAnwser{
			Result: result,
		}
		ResultJson, err := json.Marshal(ResultNotJson)
		if err != nil{
			http.Error(w, "International server error", 500)
		}
		fmt.Fprintf(w, string(ResultJson), status)
	}
}

func RunServer(){
	http.HandleFunc("/", CalcHandler)
	http.ListenAndServe(":8080", nil)
}