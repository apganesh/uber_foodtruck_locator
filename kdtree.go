package main

import (
	"fmt"
	"go-priority-queue/prio"
	"math"
)

const kmtomiles = float64(0.621371192)
const earthRadius = float64(6371)

//use the haversine distance
//http://www.movable-type.co.uk/scripts/latlong.html

type KDNode struct {
	val         *Record
	left, right *KDNode
	pos         [2]float64
	level       int
}

type KDTree struct {
	root *KDNode
}

type PQNode struct {
	value *KDNode
	dist  float64
	index int // index in heap
}

var (
	tree *KDTree
)

//func (x *PQNode) Less(y prio.Interface) bool { return x.dist > y.(*PQNode).dist }
func (x *PQNode) Less(y prio.Interface) bool { return x.dist < y.(*PQNode).dist }
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

func (root *KDNode) findNearestNeighborsRadius(ll LatLng, pq *prio.Queue, radius float64, res *TrucksRes) {
	if root == nil {
		return
	}

	lvl := root.level
	rec := root.val

	dist := HaversineDistance(ll.Lng, ll.Lat, (*root.val).p[1], (*root.val).p[0])

	if dist <= radius {
		*res = append(*res, JsonRes{rec.Name, rec.Address, rec.Fooditems, LatLng{rec.p[0], rec.p[1]}, dist, rec.Foodtypes})
		pq.Push(&PQNode{root, dist, 0})
	}
	var diff float64
	if lvl%2 == 0 {
		diff = ll.Lat - rec.p[0]
	} else {
		diff = ll.Lng - rec.p[1]
	}

	if diff < 0.0 {
		if root.left != nil {
			root.left.findNearestNeighborsRadius(ll, pq, radius, res)
		}
		if math.Abs(diff) < radius {
			if root.right != nil {
				root.right.findNearestNeighborsRadius(ll, pq, radius, res)
			}
		}
	} else {
		if root.right != nil {
			root.right.findNearestNeighborsRadius(ll, pq, radius, res)
		}
		if math.Abs(diff) < radius {
			if root.left != nil {
				root.left.findNearestNeighborsRadius(ll, pq, radius, res)
			}
		}
	}
}

// Level order printing the kdtree
func printTree(root *KDNode) {
	queue := make([]*KDNode, 0)
	queue = append(queue, root)
	index := 0

	for {
		if len(queue) == 0 {
			break
		}
		qlen := len(queue)
		fmt.Printf("%d -- %d", index, qlen)
		index++
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
			//fmt.Printf("(%f %f)   ", x.val.p[0], x.val.p[1])
		}
		fmt.Printf("\n")
	}

}

func HaversineDistance(lonFrom float64, latFrom float64, lonTo float64, latTo float64) (distance float64) {
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
