package main

import (
	"BeamGo/store"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientDB *mongo.Client

const (
	mongoConnectionString = "BEAM_MONGO_CONNECTION_STRING"
	port                  = "BEAM_SERVER_PORT"
)

func getEnvWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func GetScooterList(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("content-type", "application/json")

	entries, _ := store.GetScooterListCommand(clientDB)

	_ = json.NewEncoder(response).Encode(entries)
}

func getPaginationParameterValue(request *http.Request, key string, defaultValue int64) int64 {
	values := request.URL.Query()[key]
	if len(values) > 0 {
		limitParam := values[0]
		i, err := strconv.Atoi(limitParam)
		if err == nil && i < 100 {
			defaultValue = int64(i)
		}
	}

	return defaultValue
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Print("Error loading .env file")
	}

	log.Print("Starting the application...")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/get-scooters", GetScooterList).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./dist/")))

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// initialize database connection
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	clientOptions := options.Client().ApplyURI(getEnvWithDefault(mongoConnectionString, "mongodb://localhost:27017"))
	clientDB, _ = mongo.Connect(ctx, clientOptions)

	port := getEnvWithDefault(port, "9090")
	fmt.Println("Server listening on port " + port)
	log.Fatal(http.ListenAndServe(":" + port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
