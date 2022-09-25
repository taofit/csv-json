package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type hierarchyHandler struct{}
type Data struct {
	Item_id string `json:"item_id"`
	Level_1 string `json:"level_1"`
	Level_2 string `json:"level_2"`
}

func (f *hierarchyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	records, err := csv.NewReader(r.Body).ReadAll()
	if err != nil {
		log.Fatalf("something wrong %v", err)
	}
	var allData = []Data{}
	for _, v := range records {
		allData = append(allData, Data{v[0], v[1], v[2]})
	}
	fmt.Println(allData)
	json.NewEncoder(w).Encode(allData)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", &hierarchyHandler{})
	err := http.ListenAndServe(":8080", &hierarchyHandler{})
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
	}
}
