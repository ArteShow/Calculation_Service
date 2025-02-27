package main

import (
	"github.com/ArteShow/Calculation_Service/application"
	"github.com/ArteShow/Calculation_Service/internal"
)

//Starting the web-service
func main(){
	application.RunServer()
	internal.RunServerAgent()
}