package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"strings"
	// "github.com/formulatehq/data-engineer"
)

type hierarchyHandler struct{}

// type Data struct {
// 	Item_id string `json:"item_id"`
// 	Level_1 string `json:"level_1"`
// 	Level_2 string `json:"level_2"`
// }

type element struct {
	Value      string `json:"value"`
	ParentPath string `json:"parentPath"`
}
type Node struct {
	Item     bool
	Children map[string]*Node
}

func (f *hierarchyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	records, err := csv.NewReader(r.Body).ReadAll()
	if err != nil {
		handleBadRequest(w, err)
		return
	}
	allData, err := getAllElements(records)
	if err != nil {
		handleBadRequest(w, err)
		return
	}
	var node = Node{}
	var parentPath = ""
	generateNode(&node, parentPath, allData)
	w.WriteHeader(http.StatusCreated)
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
		if isFirstElementLevel {
			fmt.Println(row)
			if !validateRecord(row) {
				return []element{}, errors.New("invidate structure")
			}
			row = getRowRemoveEptElements(row)
			for j, ele := range row {
				parentPath := strings.Join(row[:j], ",")
				elements = append(elements, element{ele, parentPath})
			}
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

func generateNode(node *Node, currentPath string, elements []element) {
	children := getChildren(currentPath, elements)
	// fmt.Println("children------->", children)
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
