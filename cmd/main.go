package main

import (
	"fmt"

	calculate "github.com/ArteShow/Calculation_Service/pkg/Calculation"
)

func main(){
	result, err := calculate.CalcBasic("3 + 2.2")
	if err != nil{
		panic(err)
	}
	fmt.Println(result)
}