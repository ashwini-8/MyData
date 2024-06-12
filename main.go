package main

import (
	"fmt"
	"log"
	"net/http"

	"MyData/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/fetch-data", handlers.FetchDataHandler).Methods("GET") // it will get all data curl cmd for it is

	// curl -X GET http://localhost:8080//api/v1/fetch-data
	r.HandleFunc("/api/v1/data", handlers.CreateDataHandler).Methods("POST") /* curl -X POST http://localhost:8080//api/v1/data \
	-H "Content-Type: application/json" \
	-d '{
		  "Make_ID": 123,
		  "Make_Name": "TestMake"
		}'
	*/
	r.HandleFunc("/api/v1/data", handlers.ReadAllDataHandler).Methods("GET")     // curl -X GET http://localhost:8080//api/v1/data
	r.HandleFunc("/api/v1/data/{id}", handlers.ReadDataHandler).Methods("GET")   // curl -X GET http://localhost:8080/api/v1/data/<id>
	r.HandleFunc("/api/v1/data/{id}", handlers.UpdateDataHandler).Methods("PUT") /* curl -X PUT http://localhost:8080/api/v1/data/<id> \
	-H "Content-Type: application/json" \
	-d '{
		  "Make_ID": 123,
		  "Make_Name": "UpdatedMake"
		}'
	*/
	r.HandleFunc("/api/v1/data/{id}", handlers.DeleteDataHandler).Methods("DELETE") // curl -X DELETE http://localhost:8080/api/v1/data/<id>
	fmt.Println("Starting server on :8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
