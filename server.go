package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

type dataSiswa struct {
	Nisn   string `json:"nisn,omitempty"`
	Nama   string `json:"nama,omitempty"`
	Umur   int    `json:"umur,omitempty"`
	Alamat string `json:"alamat,omitempty"`
}

// function to connect to db
func connect() (*mongo.Database, error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI("mongodb://localhost:27017")

	// try to establish the client
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	// try to connect the client to the context
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database("sekolah"), nil
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func insert_siswa(w http.ResponseWriter, r *http.Request) {

	// reqBody will be the JSON data
	reqBody, _ := ioutil.ReadAll(r.Body)
	var d dataSiswa

	// parse JSON data (reqBody) and store into d (use pointer)
	json.Unmarshal(reqBody, &d)

	// try connect to db
	db, err := connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	response, err := db.Collection("siswa").InsertOne(ctx, d)
	if err != nil {
		log.Fatal(err.Error())
	}

	json.NewEncoder(w).Encode(response)
	fmt.Println("Endpoint Hit: insert_siswa, Time: ", time.Now())
}

func get_siswa(w http.ResponseWriter, r *http.Request) {
	db, err := connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	data, err := db.Collection("siswa").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer data.Close(ctx)

	// convert data to json
	result := make([]dataSiswa, 0)
	for data.Next(ctx) {
		var row dataSiswa
		err := data.Decode(&row)
		if err != nil {
			log.Fatal(err.Error())
		}

		result = append(result, row)
	}

	json.NewEncoder(w).Encode(result)
	fmt.Println("Endpoint Hit: get_siswa, Time: ", time.Now())
}

func update_siswa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var d dataSiswa
	json.Unmarshal(reqBody, &d)

	db, err := connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	data_update := bson.M{
		"$set": d,
	}

	objID, _ := primitive.ObjectIDFromHex(id)
	response, err := db.Collection("siswa").UpdateOne(ctx, bson.M{"_id": objID}, data_update)
	if err != nil {
		log.Fatal(err.Error())
	}

	json.NewEncoder(w).Encode(response)

	fmt.Println("Endpoint Hit: update_siswa, Time: ", time.Now())
}

func delete_siswa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	db, err := connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	objID, _ := primitive.ObjectIDFromHex(id)
	response, err := db.Collection("siswa").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Fatal(err.Error())
	}

	json.NewEncoder(w).Encode(response)

	fmt.Println("Endpoint Hit: delete_siswa, Time: ", time.Now())
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/insert_siswa", insert_siswa).Methods("POST")
	myRouter.HandleFunc("/get_siswa", get_siswa).Methods("GET")
	myRouter.HandleFunc("update_siswa/{id}", update_siswa).Methods("PUT")
	myRouter.HandleFunc("/delete_siswa/{id}", delete_siswa).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Program Started!!")
	handleRequests()
}
