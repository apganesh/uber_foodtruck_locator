package main

import (
	"encoding/json"
	"fmt"
	"go-priority-queue/prio"
	"net/http"
	"sort"
	"strings"
)

//use the haversine distance
//http://www.movable-type.co.uk/scripts/latlong.html

type LatLng struct {
	Lat, Lng float64
}

type JsonRes struct {
	Name      string
	Address   string
	Fooditems string
	LL        LatLng
	Distance  float64
	Foodtypes []string
}

type TrucksRes []JsonRes

type Record struct {
	Name      string
	Address   string
	Fooditems string
	p         [2]float64
	Foodtypes []string
}

type Records []*Record

type TruckDB struct {
	root    *KDNode
	records Records
}

func NewTruckDB() *TruckDB {
	tdb := &TruckDB{}
	tdb.root = nil
	return tdb
}

func (tdb *TruckDB) readJsonFile() error {

	type ObjectType struct {
		Status    string  `json:"status"`
		Address   string  `json:"address"`
		Applicant string  `json:"applicant"`
		Fooditems string  `json:"fooditems"`
		Lat       float64 `json:"latitude,string"`
		Lng       float64 `json:"longitude,string"`
	}

	var elems []ObjectType

	url := "https://data.sfgov.org/resource/rqzj-sfat.json"

	response, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return err
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&elems)

	fmt.Printf("Total records read: %d\n", len(elems))
	index := 1
	for _, p := range elems {
		if p.Status == "EXPIRED" {
			continue
		}
		fcats := getFoodCategories(p.Fooditems)
		newRec := &Record{Address: p.Address, Name: p.Applicant, Fooditems: p.Fooditems, Foodtypes: fcats, p: [2]float64{p.Lat, p.Lng}}
		tdb.records = append(tdb.records, newRec)
		index++
	}

	return nil
}

func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (tdb *TruckDB) categorizeFoodItems() {
	// Go thru each record and get the category and add the record to it
	cats := make(map[string]int)
	for _, r := range tdb.records {
		s := strings.ToLower(r.Fooditems)
		ws := strings.Split(s, ":")

		for _, w := range ws {
			w = strings.Trim(w, " ")
			cats[w]++

		}
	}

	var x PairList = rankByWordCount(cats)
	for _, v := range x {
		fmt.Printf("%s -- %d\n", v.Key, v.Value)
	}

}

func (tdb *TruckDB) updateKDTree() {

	if len(tdb.records) > 0 {
		tdb.root = tdb.buildKDTree(0, len(tdb.records), 0)
	}
}

func (tdb *TruckDB) buildKDTree(s int, e int, lvl int) (root *KDNode) {

	if e == s || e < s || s == len(tdb.records) {
		return nil
	}

	c := Comparator{r: tdb.records[s:e], l: lvl % 2}
	sort.Sort(c)
	mid := s + (e-s)/2

	root = &KDNode{val: tdb.records[mid], left: nil, right: nil, level: lvl % 2,
		pos: [2]float64{tdb.records[mid].p[0], tdb.records[mid].p[1]}}

	root.left = tdb.buildKDTree(s, mid, lvl+1)
	root.right = tdb.buildKDTree(mid+1, e, lvl+1)

	return root
}
func hasCommonTypes(s1 []string, s2 []string) bool {
	if len(s1) == 0 {
		return true
	}
	s1m := make(map[string]bool)
	for _, c := range s1 {
		s1m[c] = true
	}

	for _, c := range s2 {
		_, b := s1m[c]
		if b {
			return true
		}

	}
	return false
}

func (tdb *TruckDB) findNearestTrucks(ll LatLng, r float64, count int64, cats []string) TrucksRes {
	knear := &TrucksRes{}
	pq := &prio.Queue{}
	tdb.root.findNearestNeighborsRadius(ll, pq, r, knear)
	pqlen := pq.Len()

	var res TrucksRes

	for i := 0; i < pqlen && count > 0; i++ {
		x := pq.Remove(0)
		r := x.(*PQNode).value.val
		dist := x.(*PQNode).dist
		//fmt.Println("--------------------------")
		//fmt.Println("Distance: ", dist)
		ftypes := r.Foodtypes
		common := hasCommonTypes(cats, ftypes)
		if common {
			count--
			//fmt.Println(dist, " -> ", ftypes)
			res = append(res, JsonRes{r.Name, r.Address, r.Fooditems, LatLng{r.p[0], r.p[1]}, dist, r.Foodtypes})
		}
	}

	return res
}
