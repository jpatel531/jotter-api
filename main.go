package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	// Allow a 10KB body size
	maxBodySize = 10000

	// This formats the date based on constants defined in
	// https://golang.org/src/time/format.go
	// timeReferenceLayout = "15:04:05"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := time.Now().Format(time.RFC850)
	log.Println("Timestamp", t, "Body", r.Body)
	fmt.Fprint(w, "Hey baaaaaby\n")
}

type updateRequest struct {
	Type string `json:"type"`

	Timestamp string `json:"timestamp,omitempty"`

	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Altitude  float64 `json:"altitude,omitempty"`
	Speed     float64 `json:"speed,omitempty"`

	Distance          float64 `json:"distance,omitempty"`
	NumberOfSteps     float64 `json:"numberOfSteps,omitempty"`
	AverageActivePace float64 `json:"averageActivePace,omitempty"`
	FloorsAscended    float64 `json:"floorsAscended,omitempty"`
	FloorsDescended   float64 `json:"floorsDescended,omitempty"`
}

func ReadJSONBody(body io.Reader, dest interface{}) (failure bool) {
	decoder := json.NewDecoder(body)
	if decodeErr := decoder.Decode(&dest); decodeErr != nil {
		log.Println("Unable to JSON parse the request body", decodeErr)
		failure = true
	}
	return
}

func update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var (
		updateReq updateRequest
		failure   bool
	)

	t := time.Now().Format(time.RFC850)
	log.Println("Timestamp", t)

	// Use a max bytes reader to limit the size of the body
	maxReader := http.MaxBytesReader(w, r.Body, maxBodySize)
	maxBody, readErr := ioutil.ReadAll(maxReader)
	if readErr != nil {
		log.Println("Request body size > 10KB")
		return
	}

	if failure = ReadJSONBody(bytes.NewReader(maxBody), &updateReq); failure {
		log.Println("body", string(maxBody))
		return
	}

	// log.Println("updateReq: ", updateReq)
	fmt.Printf("%+v\n", updateReq)
}

func main() {
	router := httprouter.New()
	router.GET("/", index)
	router.POST("/", index)
	router.POST("/update", update)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
