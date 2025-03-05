package main

import (
	"log"
	"time"

	"github.com/ArteShow/Calculation_Service/application"
	"github.com/ArteShow/Calculation_Service/internal"

)

func main() {
	log.Println("🚀 Starte beide Server...")

	
	go func() {
		internal.RunServerAgent()
	}()
	
	time.Sleep(2 * time.Second)

	
	application.RunServer()
}
