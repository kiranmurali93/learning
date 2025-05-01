package main

import (
	"job-claim-concurrency-demo/db"
	"job-claim-concurrency-demo/handlers"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// Ignore this error
func main() {
	godotenv.Load()
	db.InitDB() 
	http.HandleFunc("/claim-job", handlers.ClaimJob)
	log.Println("Server Running at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}