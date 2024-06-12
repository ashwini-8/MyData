package handlers

import (
	"MyData/database"
	"MyData/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const apiURL = "https://vpic.nhtsa.dot.gov/api/vehicles/getallmakes?format=json"

// Get the data from the URL
func fetchData(url string) (*models.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Insert data to mongodb
func insertData(client *mongo.Client, data []models.Data) error {
	collection := client.Database(database.DatabaseName).Collection("data")

	var documents []interface{}
	for _, datum := range data {
		documents = append(documents, datum)
	}

	_, err := collection.InsertMany(context.Background(), documents)
	return err
}

// Check connection with mongdb and insert data into it
func FetchDataHandler(w http.ResponseWriter, r *http.Request) {
	data, err := fetchData(apiURL)
	if err != nil {
		http.Error(w, "Error fetching data from API", http.StatusInternalServerError)
		return
	}
	client, err := database.GetMongoClient()
	if err != nil {
		http.Error(w, "Error connecting to MongoDB", http.StatusInternalServerError)
		return
	}

	err = insertData(client, data.Results)
	if err != nil {
		http.Error(w, "Error inserting data into MongoDB", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data successfully inserted into MongoDB"))
}

// Create handler and insert one data into it, used for POST
func CreateDataHandler(w http.ResponseWriter, r *http.Request) {
	var dat models.Data
	if err := json.NewDecoder(r.Body).Decode(&dat); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	client, err := database.GetMongoClient()
	if err != nil {
		http.Error(w, "Error connecting to MongoDB", http.StatusInternalServerError)
		return
	}

	collection := client.Database(database.DatabaseName).Collection("data")
	_, err = collection.InsertOne(context.Background(), dat)
	if err != nil {
		http.Error(w, "Error inserting data into MongoDB", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dat)
}

// Handler to get all the data from mongodb
func ReadAllDataHandler(w http.ResponseWriter, r *http.Request) {
	client, err := database.GetMongoClient()
	if err != nil {
		http.Error(w, "Error connecting to MongoDB", http.StatusInternalServerError)
		return
	}

	collection := client.Database(database.DatabaseName).Collection("data")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Error reading data from MongoDB", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var data []models.Data
	for cursor.Next(context.Background()) {
		var datum models.Data
		if err := cursor.Decode(&datum); err != nil {
			http.Error(w, "Error decoding data from MongoDB", http.StatusInternalServerError)
			return
		}
		data = append(data, datum)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// Handler to read specific data using id
func ReadDataHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	client, err := database.GetMongoClient()
	if err != nil {
		http.Error(w, "Error connecting to MongoDB", http.StatusInternalServerError)
		return
	}

	collection := client.Database(database.DatabaseName).Collection("data")
	var datum models.Data
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&datum)
	if err != nil {
		http.Error(w, "Error reading data from MongoDB, or the object is not present in the data base", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(datum)
}

// Handler used to update the data into mongodb
func UpdateDataHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("value for r is: ", r)
	params := mux.Vars(r)
	fmt.Println("params is: ", params)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var datum models.Data
	if err := json.NewDecoder(r.Body).Decode(&datum); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	client, err := database.GetMongoClient()
	if err != nil {
		http.Error(w, "Error connecting to MongoDB", http.StatusInternalServerError)
		return
	}

	collection := client.Database(database.DatabaseName).Collection("data")
	fmt.Println("collection is: ", collection)

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": datum})
	if err != nil {
		http.Error(w, "Error updating data in MongoDB", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(datum)
}

// Handler to Delete data
func DeleteDataHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	client, err := database.GetMongoClient()
	if err != nil {
		http.Error(w, "Error connecting to MongoDB", http.StatusInternalServerError)
		return
	}

	collection := client.Database(database.DatabaseName).Collection("data")
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Error deleting data from MongoDB", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data successfully deleted from MongoDB"))
}
