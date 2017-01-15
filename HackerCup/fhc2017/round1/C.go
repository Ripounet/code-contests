package main

// Facebook Hacker Cup 2017 Round 1

// 3rd problem: nice attempt, but failed!!
// Answers are wrong for some reason.
// The answers to sample cases were correct, though.

import (
	"container/heap"
	"fmt"
	"os"
	"runtime"
	"time"
)

var (
	pathIn = "/home/valou/Téléchargements/"
	// pathIn   = "./"
	pathOut  = "./"
	letter   = "C"
	strategy = (*Case).solveSmall
	// strategy = (*Case).solveLarge
	// strategy   = (*Case).solveLargeAndCheck
	concurrent = false
	maxProc    = 4
)

type Town struct {
	n int

	// out[t2] (if exists) is "cheapest gas single road to reach t2"
	out map[int]int

	// dijsktra
	dist  int
	prev  int
	index int // The index of the item in the heap.
}

// Put inputs here as global vars
var (
	N, M, K int
	A, B, G []int
	S, D    []int

	towns []*Town

	// How much gas to drive from i to j
	cost [][]int
)

// Case data
type Case struct {
	caseNumber int
}

func (z *Case) readSingle() {
	check(!concurrent) // So we can use global vars
	N = readInt()
	M = readInt()
	K = readInt()

	towns = make([]*Town, N)
	for i := range towns {
		towns[i] = new(Town)
		towns[i].n = i
		towns[i].out = make(map[int]int)
	}
	cost = make([][]int, N)
	for i := range cost {
		cost[i] = make([]int, N)
	}

	A = make([]int, M)
	B = make([]int, M)
	G = make([]int, M)
	for i := 0; i < M; i++ {
		A[i] = readInt() - 1
		B[i] = readInt() - 1
		G[i] = readInt()
		if gas, ok := towns[A[i]].out[B[i]]; !ok || gas > G[i] {
			towns[A[i]].out[B[i]] = G[i]
			towns[B[i]].out[A[i]] = G[i]
		}
	}

	S = make([]int, K)
	D = make([]int, K)
	for i := 0; i < K; i++ {
		S[i] = readInt() - 1
		D[i] = readInt() - 1
	}
}

func (z *Case) solveSmall() interface{} {
	// shortest path from any to any
	dijkstra := func(source *Town) {

		source.dist = 0

		// Create a priority queue, put the items in it, and
		// establish the priority queue (heap) invariants.
		q := make(PriorityQueue, N)
		for i, v := range towns {
			q[i] = v
			v.index = i
			if v != source {
				v.dist = 999999999
				v.prev = -1
			}
		}
		heap.Init(&q)

		for q.Len() > 0 {
			u := heap.Pop(&q).(*Town)
			cost[source.n][u.n] = u.dist

			for v, gasUV := range towns[u.n].out {
				alt := u.dist + gasUV
				if alt < towns[v].dist {
					towns[v].prev = u.n
					q.update(towns[v], alt)
				}
			}
		}
	}
	for i := 0; i < N; i++ {
		dijkstra(towns[i])
	}

	// log(A)
	// log(B)
	// log(G)
	// log("--")
	// log(S)
	// log(D)
	// log("--")
	// log("cost", cost)

	bestEmpty := make([]int, K)
	bestLoaded := make([]int, K)

	bestEmpty[0] = cost[0][S[0]] + cost[S[0]][D[0]]
	if bestEmpty[0] >= 999999999 {
		return -1
	}
	if K >= 2 {
		bestLoaded[0] = cost[0][S[0]] + cost[S[0]][S[1]] + cost[S[1]][D[0]]
	}

	for i := 1; i < K; i++ {
		loc := D[i-1]
		e1 := bestEmpty[i-1] + cost[loc][S[i]] + cost[S[i]][D[i]]
		e2 := bestLoaded[i-1] + cost[loc][D[i]]
		bestEmpty[i] = min(e1, e2)

		if i+1 < K {
			bestLoaded[i] = bestEmpty[i-1] + cost[loc][S[i]] + cost[S[i]][S[i+1]] + cost[S[i+1]][D[i]]
		}

		if bestEmpty[i] >= 999999999 || bestLoaded[i] >= 999999999 {
			return -1
		}
	}
	// log("bestEmpty", bestEmpty)
	// log("bestLoaded", bestLoaded)

	return bestEmpty[K-1]
}

func (z *Case) solveLarge() interface{} {
	return nil
}

// Global precomputed data (if needed)

func precompute() {
}

//
//
//
// NO NEED TO EDIT BELOW
//
//
//

func (z *Case) solveSingle() (answer interface{}) {
	defer func() {
		logf("------------ Case #%v: %v --------", z.caseNumber, answer)
	}()
	return strategy(z)
}

func (z *Case) solveLargeAndCheck() interface{} {
	// But make sure solveLarge is non-destructive!!
	sol1 := z.solveLarge()
	sol2 := z.solveSmall()
	if sol1 != sol2 {
		log("small = ", sol2, "; large = ", sol1)
		panic("Small and Large strategies must be equivalent")
	}
	return sol1
}

func solve() {
	if concurrent {
		runtime.GOMAXPROCS(maxProc)
	}

	top1 := time.Now() //.UnixNano()
	precompute()

	T := readInt()
	solutions := make([]chan interface{}, T)
	for i := range solutions {
		solutions[i] = make(chan interface{})
	}
	go func() {
		for i := range solutions {
			currentCase := &Case{caseNumber: 1 + i}
			currentCase.readSingle()
			if concurrent {
				go func(ch chan interface{}) {
					ch <- currentCase.solveSingle()
					close(ch)
				}(solutions[i])
			} else {
				solutions[i] <- currentCase.solveSingle()
				close(solutions[i])
			}
		}
	}()
	for i, ch := range solutions {
		solution := <-ch
		outf("Case #%v: %v\n", 1+i, solution)
	}
	duration := time.Since(top1)
	seconds := float64(duration) / 1000000000.0
	logf("Took %6.1fs \n", seconds)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	sample := os.Args[1]

	var fileIn = pathIn + letter + "-" + sample + ".in"
	var fileOut = pathOut + letter + "-" + sample + ".out"

	var err error
	input, err = os.Open(fileIn)
	if err != nil {
		panic(fmt.Sprintf("open %s: %v", fileIn, err))
	}
	output, err = os.Create(fileOut)
	if err != nil {
		panic(fmt.Sprintf("creating %s: %v", fileOut, err))
	}
	defer input.Close()
	defer output.Close()

	solve()
}

func usage() {
	logf("Usage: %v <sample> \n", os.Args[0])
	os.Exit(1)
}

var input *os.File
var output *os.File

func check(condition bool) {
	if !condition {
		panic("Failed!")
	}
}

func out(str string) {
	//	fmt.Print(str)
	fmt.Fprint(output, str)
}

func outf(pattern string, values ...interface{}) {
	str := fmt.Sprintf(pattern, values...)
	out(str)
}

func logf(str string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, str, values...)
	fmt.Fprint(os.Stderr, "\n")
}

func log(values ...interface{}) {
	fmt.Fprintln(os.Stderr, values...)
}

func readInt() int {
	var i int
	fmt.Fscanf(input, "%d", &i)
	return i
}

func readString() string {
	var str string
	fmt.Fscanf(input, "%s", &str)
	return str
}

func readFloat() float64 {
	var x float64
	fmt.Fscanf(input, "%f", &x)
	return x
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

// PQ

// An Item is something we manage in a priority queue.
/*
type Item struct {
	iTown int // The value of the item; arbitrary.
	dist  int // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.

	prev int
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
*/

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Town

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Town)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(town *Town, dist int) {
	town.dist = dist
	heap.Fix(pq, town.index)
}
