package ui

import(
	"net/http"
	"log"
)

func StartUi(){
	http.Handle("/", http.FileServer(http.Dir("./")))

	
	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}