package kdtree

import (
	"encoding/json"
	"fmt"
	"go-priority-queue/prio"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strings"
)

const kmtomiles = float64(0.621371192)
const earthRadius = float64(6371)

//use the haversine distance
//http://www.movable-type.co.uk/scripts/latlong.html

type KDNode struct {
	val         *Record
	left, right *KDNode
	level       int
}
type Latlng struct {
	Lat, Lng float64
}

type JsonRes struct {
	Name      string
	Address   string
	Distance  float64
	Fooditems string
	Lat       float64
	Lng       float64
}
type LatlngRes []JsonRes

type Record struct {
	app       string
	address   string
	fooditems string
	p         [2]float64
	id        int
}

type Records []*Record

type PQNode struct {
	value *KDNode
	dist  float64
	index int // index in heap
}

var (
	tree *KDNode
	ftdb Records
)

func (x *PQNode) Less(y prio.Interface) bool { return x.dist > y.(*PQNode).dist }
func (x *PQNode) Index(i int)                { x.index = i }

type Comparator struct {
	r Records
	l int
}

// Comparators for sorting w.r.t X ()
func (c Comparator) Len() int {
	return len(c.r)
}

func (c Comparator) Swap(i, j int) {
	c.r[i], c.r[j] = c.r[j], c.r[i]
}

func (c Comparator) Less(i, j int) bool {
	if c.r[i].p[c.l] == c.r[j].p[c.l] {
		return c.r[i].p[(c.l+1)%2] < c.r[j].p[(c.l+1)%2]
	}
	return c.r[i].p[c.l] < c.r[j].p[c.l]
}

func readJsonFile() Records {
	var recs Records
	fitems := map[string]int{}

	type ObjectType struct {
		Status    string  `json:"status"`
		Address   string  `json:"address"`
		Applicant string  `json:"applicant"`
		Fooditems string  `json:"fooditems"`
		Lon       float64 `json:"longitude,string"`
		Lat       float64 `json:"latitude,string"`
	}

	url := "https://data.sfgov.org/resource/rqzj-sfat.json"

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	raw, e := ioutil.ReadAll(response.Body)
	if e != nil {
		fmt.Printf("Error occurred reading the file\n")
	}

	var elems []ObjectType

	json.Unmarshal(raw, &elems)
	index := 0
	for _, p := range elems {
		if p.Status == "EXPIRED" {
			continue
		}

		tokens := strings.Split(p.Fooditems, ":")
		for _, v := range tokens {
			fitems[v] = 1
		}

		newRec := &Record{id: index, address: p.Address, app: p.Applicant, fooditems: p.Fooditems, p: [2]float64{p.Lat, p.Lon}}
		ftdb = append(ftdb, newRec)
		index = index + 1
	}

	return recs
}

func BuildTree() *KDNode {
	fmt.Println("Starting to read the json file")
	readJsonFile()
	fmt.Println("Finished reading json file")
	fmt.Println("Starting to build the kdtree")
	tree = buildKDTree(ftdb, 0)
	fmt.Println("Finished to build the kdtree")
	fmt.Println("Waiting for clients")
	//printNode(tree)
	return tree
}

func buildKDTree(recs Records, lvl int) (root *KDNode) {
	if len(recs) == 0 {
		return nil
	}

	if len(recs) == 1 {
		root = &KDNode{val: recs[0], left: nil, right: nil, level: lvl % 2}
	} else {
		c := Comparator{r: recs, l: lvl % 2}
		sort.Sort(c)
		mid := len(recs) / 2

		root = &KDNode{val: recs[mid], left: nil, right: nil, level: lvl % 2}

		root.left = buildKDTree(recs[:mid], lvl+1)
		root.right = buildKDTree(recs[mid+1:], lvl+1)
	}

	return root
}

func printNode(root *KDNode) {
	if root == nil {
		return
	}

	fmt.Printf("%f -- %f\n", root.val.p[0], root.val.p[1])
}

// Level order printing the kdtree
func printTree(root *KDNode) {
	queue := make([]*KDNode, 0)
	queue = append(queue, root)

	for {
		if len(queue) == 0 {
			break
		}
		qlen := len(queue)
		for i := 0; i < qlen; i++ {
			x := queue[0]
			queue = queue[1:]
			if x == nil {
				continue
			}
			if x.left != nil {
				queue = append(queue, x.left)
			}
			if x.right != nil {
				queue = append(queue, x.right)
			}
			fmt.Printf("(%f %f)   ", x.val.p[0], x.val.p[1])
		}
		fmt.Printf("\n")
	}

}

func distance(ll Latlng, rr Record) float64 {
	// Use the hervesian or goog distance
	return math.Sqrt(math.Pow((ll.Lat-rr.p[0]), 2.0) + math.Pow((ll.Lng-rr.p[1]), 2.0))
}

func Haversine(lonFrom float64, latFrom float64, lonTo float64, latTo float64) (distance float64) {
	// From golang playground
	var deltaLat = (latTo - latFrom) * (math.Pi / 180)
	var deltaLon = (lonTo - lonFrom) * (math.Pi / 180)

	var a = math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(latFrom*(math.Pi/180))*math.Cos(latTo*(math.Pi/180))*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance = earthRadius * c
	distance = distance * kmtomiles
	return
}

func FindNearestRadius(ll Latlng, r float64) LatlngRes {
	knear := &LatlngRes{}
	pq := &prio.Queue{}
	findNearestRadius(tree, ll, pq, r, knear)
	pqlen := pq.Len()

	var knearres LatlngRes

	for i := 0; i < pqlen; i++ {
		x := pq.Remove(0)
		r := x.(*PQNode).value.val
		knearres = append(knearres, JsonRes{r.app, r.address, x.(*PQNode).dist, r.fooditems, r.p[0], r.p[1]})
		//fmt.Printf("%.2f %f %f\n", x.(*PQNode).dist, r.p[0], r.p[1])
	}

	// return *knear
	return knearres
}

func findNearestRadius(root *KDNode, ll Latlng, pq *prio.Queue, radius float64, res *LatlngRes) {
	if root == nil {
		return
	}
	lvl := root.level
	rec := root.val
	dist := Haversine(ll.Lng, ll.Lat, (*root.val).p[1], (*root.val).p[0])
	if dist <= radius {
		*res = append(*res, JsonRes{rec.app, rec.address, dist, rec.fooditems, rec.p[0], rec.p[1]})
		pq.Push(&PQNode{root, dist, 0})
	}
	var diff float64
	if lvl%2 == 0 {
		diff = ll.Lat - rec.p[0]
		//diff = Haversine(0.0, ll.Lat, 0.0, (*root.val).p[0])
	} else {
		diff = ll.Lng - rec.p[1]
		//diff = Haversine(ll.Lng, 0.0, (*root.val).p[1], 0.0)
	}

	if diff < 0.0 {
		findNearestRadius(root.left, ll, pq, radius, res)
		if math.Abs(diff) < radius {
			findNearestRadius(root.right, ll, pq, radius, res)
		}
	} else {
		findNearestRadius(root.right, ll, pq, radius, res)
		if math.Abs(diff) < radius {
			findNearestRadius(root.left, ll, pq, radius, res)
		}
	}
}

func FindKNearest(ll Latlng, k int) LatlngRes {
	pq := &prio.Queue{}
	kNN(tree, ll, pq, 0, k)
	pqlen := pq.Len()

	var knear LatlngRes

	for i := 0; i < pqlen; i++ {
		x := pq.Remove(0)
		r := x.(*PQNode).value.val
		knear = append(knear, JsonRes{r.app, r.address, x.(*PQNode).dist, r.fooditems, r.p[0], r.p[1]})
		fmt.Printf("%f %f %f\n", x.(*PQNode).dist, r.p[0], r.p[1])
	}
	return knear
}

func kNN(root *KDNode, ll Latlng, pq *prio.Queue, lvl int, k int) {
	if root == nil {
		return
	}

	//dd := distance(ll, *root.val)
	dist := Haversine(ll.Lng, ll.Lat, (*root.val).p[1], (*root.val).p[0])

	if pq.Len() < k {
		pq.Push(&PQNode{root, dist, 0})
	} else {
		if dist < pq.Peek().(*PQNode).dist {
			pq.Pop()
			pq.Push(&PQNode{root, dist, 0})
		}
	}

	lvl = lvl % 2
	kNN(root.left, ll, pq, lvl+1, k)
	kNN(root.right, ll, pq, lvl+1, k)

}
