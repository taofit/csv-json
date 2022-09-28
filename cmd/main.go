package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"strings"
)

type hierarchyHandler struct{}

type element struct {
	Value      string `json:"value"`
	ParentPath string `json:"parentPath"`
}
type node struct {
	Item     bool             `json:"item,omitempty"`
	Children map[string]*node `json:"children,omitempty"`
}

func (f *hierarchyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	records, err := csv.NewReader(r.Body).ReadAll()
	if err != nil {
		handleBadRequest(w, err)
		return
	}
	allElements, err := getAllElements(records)
	if err != nil {
		handleBadRequest(w, err)
		return
	}
	var node = node{}
	var currentPath = ""
	generateNode(&node, currentPath, allElements)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(node)
}

func handleBadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	resp := make(map[string]string)
	resp["message"] = err.Error()
	json.NewEncoder(w).Encode(resp)
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

		if !isFirstElementLevel {
			row = append(row[1:], row[0])
		}
		if !validateRecord(row) {
			return []element{}, errors.New("invalid structure")
		}
		row = getRowRemoveEptElements(row)
		for j, ele := range row {
			parentPath := strings.Join(row[:j], ",")
			elements = append(elements, element{ele, parentPath})
		}

	}

	return elements, nil
}

func validateRecord(row []string) bool {
	if len(row) < 2 {
		return false
	}
	if row[0] == "" || row[len(row)-1] == "" {
		return false
	}
	var hasEptElement = false
	for i, element := range row {
		if element == "" {
			hasEptElement = true
		}
		if element != "" && hasEptElement && i != len(row)-1 {
			return false
		}
	}

	return true
}

func getRowRemoveEptElements(row []string) []string {
	indiceOfEptElements := getIndiceOfEptElements(row)
	if len(indiceOfEptElements) > 0 {
		row = removeEptElements(row, indiceOfEptElements)
	}
	return row
}

func getIndiceOfEptElements(row []string) (indiceOfEptElements []int) {
	indiceOfEptElements = []int{}
	for i, element := range row {
		if element == "" {
			indiceOfEptElements = append(indiceOfEptElements, i)
		}
	}
	return
}

func removeEptElements(row []string, indiceOfEptElements []int) []string {
	firstIndexOfEptElement := indiceOfEptElements[0]
	lastIndexOfEptElement := indiceOfEptElements[len(indiceOfEptElements)-1]
	return append(row[:firstIndexOfEptElement], row[lastIndexOfEptElement+1:]...)
}

func generateNode(aNode *node, currentPath string, elements []element) {
	children := getChildren(currentPath, elements)
	aNode.Children = map[string]*node{}

	for key := range children {
		aNode.Children[key] = &node{}
		var childCurrentPath string
		if currentPath == "" {
			childCurrentPath = key
		} else {
			childCurrentPath = currentPath + "," + key
		}
		generateNode(aNode.Children[key], childCurrentPath, elements)
	}
	if len(children) == 0 {
		var pathArr = strings.Split(currentPath, ",")
		var currentValue = pathArr[len(pathArr)-1]
		if currentValue != "" {
			aNode.Item = true
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
