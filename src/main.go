package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

var _ = godotenv.Load()

var ctx = context.Background()
var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: os.Getenv("REDIS_PASSWORD"),
	DB:       0,
})

func ResponseWriter(w http.ResponseWriter, error bool, data string) {
	if error {
		w.WriteHeader(http.StatusBadRequest)
	}

	response, _ := json.Marshal(&JSONResponse{
		Error: error,
		Data:  data,
	})
	w.Write(response)
}

func Analyzer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	body := json.NewDecoder(r.Body)
	var data RequestData
	err := body.Decode(&data)

	if err != nil {
		ResponseWriter(w, true, "Invalid request: "+err.Error())
		return
	}

	key := "data"
	value, _ := json.Marshal(data)
	DataWriter(key, string(value))

	ResponseWriter(w, false, "ok")
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", Analyzer)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}

func DataWriter(key string, value string) {
	err := client.Set(key, value, 0).Err()

	if err != nil {
		log.Print(err)
	}
}
