package main

import (
	"log"
	"time"
	"github.com/ArteShow/Calculation_Service/application"
	//"github.com/ArteShow/Calculation_Service/internal"

)

func main() {
	log.Println("🚀 Starte beide Server...")

	// Internal-Server parallel starten
	//go func() {
		//internal.RunServerAgent()
	//}()

	// Warten, damit internal zuerst startet
	time.Sleep(2 * time.Second)

	// Application-Server starten
	application.RunServer()
	//internal.RunServerAgent()
}