package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type FoodCategory struct {
	Name     string `json:"name"`
	Keywords string `json:"keywords"`
}

var (
	foodCategoryMap = make(map[int]string)
	catTrie         *Trie
)

func initializeFoodCategories() {
	// Read json file
	catTrie = NewTrie()
	raw, err := ioutil.ReadFile("./html/food_categories.json")

	if err != nil {
		fmt.Println(err)
	}

	var elems []FoodCategory

	json.Unmarshal(raw, &elems)

	index := 1

	for _, e := range elems {
		foodCategoryMap[index] = e.Name
		catTrie.addString(e.Keywords, index)
		index++
	}

}

func getFoodCategories(ss string) []string {

	s := strings.ToLower(ss)
	ws := strings.Split(s, ":")

	var str string
	if len(ws) == 0 {
		str = s
	} else {
		str = strings.Join(ws, " ")
	}

	res := catTrie.getIndices(str)

	var cats []string
	for _, r := range res {
		cats = append(cats, foodCategoryMap[r])
	}
	return cats
}
