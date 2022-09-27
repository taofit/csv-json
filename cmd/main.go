package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	// "reflect"
	"strings"
	// "github.com/formulatehq/data-engineer"
)

type hierarchyHandler struct{}
type Data struct {
	Item_id string `json:"item_id"`
	Level_1 string `json:"level_1"`
	Level_2 string `json:"level_2"`
}

type element struct {
	Value      string `json:"value"`
	ParentPath string `json:"parentPath"`
	// Children map[string]element
}

type Node struct {
	Item     bool
	Children map[string]*Node
}

func (f *hierarchyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	records, err := csv.NewReader(r.Body).ReadAll()
	if err != nil {
		log.Fatalf("something wrong %v", err.Error())
	}
	allData, err := getAllElements(records)
	if err != nil {
		log.Fatalf("somtthing wrong %v", err.Error())
	}
	var initialNode = Node{}
	var parentPath = ""
	generateNode(&initialNode, parentPath, allData)
	json.NewEncoder(w).Encode(initialNode)
}

func getAllElements(records [][]string) ([]element, error) {
	var elements = []element{}
	var isFirstElementLevel = true
	for i, row := range records {
		if i == 0 {
			if row[0] != "level_1" && row[0] != "item_id" {
				return []element{}, errors.New("wrong header structure")
			}
			if row[0] == "item_id" {
				isFirstElementLevel = false
			}
			continue
		}
		if isFirstElementLevel {
			for j, ele := range row {
				parentPath := strings.Join(row[:j], ",")
				elements = append(elements, element{ele, parentPath})
			}
		}
	}
	return elements, nil
}

func validateRecord() {
	
}

func generateNode(node *Node, currentPath string, elements []element) {
	children := getChildren(currentPath, elements)
	fmt.Println("children------->", children)
	node.Children = make(map[string]*Node)

	for key := range children {
		node.Children[key] = &Node{}
		var tempCurrentPath string
		if currentPath == "" {
			tempCurrentPath = key
		} else {
			tempCurrentPath = currentPath + "," + key
		}
		fmt.Println("key----->", "\""+key+"\"", tempCurrentPath)
		generateNode(node.Children[key], tempCurrentPath, elements)
	}
	if len(children) == 0 {
		var pathArr = strings.Split(currentPath, ",")
		var currentValue = pathArr[len(pathArr)-1]
		if currentValue != "" {
			node.Item = true
		}
	}
}

func getChildren(parentPath string, elements []element) (childrenSet map[string]bool) {
	childrenSet = make(map[string]bool)

	for _, element := range elements {
		if element.ParentPath == parentPath {
			key := element.Value
			if _, found := childrenSet[key]; !found {
				childrenSet[key] = true
			}
		}
	}

	return
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", &hierarchyHandler{})
	err := http.ListenAndServe(":8080", &hierarchyHandler{})
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
	}
}
