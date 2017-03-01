package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"sort"
)

var V, E, R, C, X int

var S []int
var endpoints []Endpoint
var requests []Request
var requestsFrom = map[VideoEndpoint]int{}
var caches []Cache

// This is only an approximation (true gain depends on rest of world)
var gainPerRequest = map[VideoCache]int{}
var gain = map[VideoCache]int{}

var usedCaches = 0
var videoPenalty []int

func main() {
	read()
	precompute()
	solve()
	printSolution()
	scoring()
}

func read() {
	fmt.Scan(&V)
	fmt.Scan(&E)
	fmt.Scan(&R)
	fmt.Scan(&C)
	fmt.Scan(&X)

	S = make([]int, V)
	for i := range S {
		fmt.Scan(&S[i])
	}

	caches = make([]Cache, C)
	for c := range caches {
		caches[c].timeToEndpoint = make(map[int]int)
		caches[c].hasVideo = make(map[int]bool)
	}

	endpoints = make([]Endpoint, E)
	for e := range endpoints {
		fmt.Scan(&endpoints[e].LD)
		fmt.Scan(&endpoints[e].K)
		endpoints[e].con = make([]Connection, endpoints[e].K)
		for j := range endpoints[e].con {
			con := endpoints[e].con
			fmt.Scan(&con[j].c)
			fmt.Scan(&con[j].Lc)
			check(con[j].Lc <= endpoints[e].LD)

			caches[con[j].c].timeToEndpoint[e] = con[j].Lc
		}
	}

	requests = make([]Request, R)
	for i := range requests {
		fmt.Scan(&requests[i].v)
		fmt.Scan(&requests[i].e)
		fmt.Scan(&requests[i].n)
		ve := VideoEndpoint{requests[i].v, requests[i].e}
		requestsFrom[ve] += requests[i].n
	}

	videoPenalty = make([]int, V)

	// log("VERCX =", V, E, R, C, X)
	// log("S =", S)
	// log("endpoints =", endpoints)
	// log("requests =", requests)
}

func precompute() {
	// Do the heavy work ... only once, and save data for next time
	filename := fmt.Sprintf("%d_%d_%d_%d_%d.gob", V, E, R, C, X)
	if len(os.Args) >= 2 {
		filename = os.Args[1] + "_" + filename
	}
	_, err := os.Stat(filename)
	exists := !os.IsNotExist(err)

	if exists {
		// Deserialize already computed data
		// log("Reading", filename)
		f, err := os.Open(filename)
		check(err == nil)
		dec := gob.NewDecoder(f)
		err = dec.Decode(&gain)
		check(err == nil)
		return
	}

	// Compute expected gain of any placement
	for v := 0; v < V; v++ {
		for c, cache := range caches {
			vc := VideoCache{v, c}

			totalGainPerRequest := 0
			for e, t := range cache.timeToEndpoint {
				datacenterTime := endpoints[e].LD
				gain := datacenterTime - t
				check(gain >= 0)
				totalGainPerRequest += gain
			}

			// Let's factor the size!  Heavy videos are less interessant to cache.
			vsize := S[v]
			totalGainPerRequest = (totalGainPerRequest * 1000) / vsize

			gainPerRequest[vc] = totalGainPerRequest
		}
	}

	// Ouch!!  1,000,000 req * 1,000 caches == .....
	for _, req := range requests {
		for c := range caches {
			vc := VideoCache{req.v, c}
			gain[vc] += req.n * gainPerRequest[vc]
		}
	}

	// Serialize already computed data
	log("Writing", filename)
	f, err := os.Create(filename)
	check(err == nil)
	enc := gob.NewEncoder(f)
	err = enc.Encode(gain)
	if err != nil {
		panic(err)
	}
}

func solve() {

	// for vc := range gain {
	// 	vsize := S[vc.V]
	// 	gain[vc] *= (1000 - vsize)
	// }

	for c, cache := range caches {
		gainsPerVideo := make([]VideoGain, 0, V)
		for v := 0; v < V; v++ {
			vc := VideoCache{v, c}
			if g := gain[vc]; g > 0 {
				gainsPerVideo = append(gainsPerVideo, VideoGain{v, g})
			}
		}
		sort.Slice(gainsPerVideo, func(i, j int) bool {
			// DESCENDING order of gain!!
			gpvi := gainsPerVideo[i]
			gpvj := gainsPerVideo[j]
			worthi := gpvi.g / (10 + videoPenalty[gpvi.v])
			worthj := gpvj.g / (10 + videoPenalty[gpvj.v])
			return worthi > worthj
		})

		space := X
		for _, gv := range gainsPerVideo {
			vsize := S[gv.v]
			if vsize > space {
				continue
			}
			space -= vsize
			cache.hasVideo[gv.v] = true
			videoPenalty[gv.v] += 50
		}

		usedCaches++
	}

}

func printSolution() {
	fmt.Println(usedCaches)

	for c, cache := range caches {
		if len(cache.hasVideo) > 0 {
			fmt.Printf("%d", c)
			space := X
			for v := range cache.hasVideo {
				vsize := S[v]
				check(vsize <= space)
				fmt.Printf(" %d", v)
				space -= vsize
			}
		}
		fmt.Println()
	}
}

func scoring() {
	scoreSum := 0
	nbReq := 0
	for _, req := range requests {
		datacenterTime := endpoints[req.e].LD
		bestCacheTime := datacenterTime
		for _, cache := range caches {
			if t, connected := cache.timeToEndpoint[req.e]; connected && cache.hasVideo[req.v] && t < bestCacheTime {
				bestCacheTime = t
			}
		}
		saved := req.n * (datacenterTime - bestCacheTime)
		scoreSum += saved
		nbReq += req.n
	}
	score := (1000 * scoreSum) / nbReq
	log("Score =", score)
}

type Endpoint struct {
	LD  int
	K   int
	con []Connection
}

type Connection struct {
	c  int
	Lc int
}

type Request struct {
	v int
	e int
	n int
}

type Cache struct {
	// input (given)
	timeToEndpoint map[int]int

	// output (to be produced)
	hasVideo map[int]bool
}

type VideoEndpoint struct {
	v int
	e int
}

type VideoCache struct {
	V int
	C int
}

type VideoGain struct {
	v int
	g int
}

func check(condition bool) {
	if !condition {
		panic("Failed!")
	}
}

func logf(str string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, str, values...)
	fmt.Fprint(os.Stderr, "\n")
}

func log(values ...interface{}) {
	fmt.Fprintln(os.Stderr, values...)
}
