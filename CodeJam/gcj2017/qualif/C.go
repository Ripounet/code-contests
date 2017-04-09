package main

// Google Code Jam 2017 Qualification round

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
	pathOut = "./"
	letter  = "C"
	// strategy = (*Case).solveSmall
	// strategy = (*Case).solveSmall2
	strategy = (*Case).solveLarge
	// strategy   = (*Case).solveLargeAndCheck
	concurrent = false
	maxProc    = 4
)

// Case data
type Case struct {
	caseNumber int
	// put fields here
	K, N int
}

func (z *Case) readSingle() {
	z.N = readInt()
	z.K = readInt()
}

func (z *Case) solveSmall() interface{} {
	stalls := make([]bool, z.N+2)
	stalls[0], stalls[z.N+1] = true, true
	var best int
	for i := 0; i < z.K; i++ {
		best = -1
		mark := 0
		for mark < z.N+2 && stalls[mark] {
			mark++
		}
		check(mark < z.N+2)
		best = mark
		check(best >= 1)
		check(best <= z.N)
		longestEmpty := 1
		for j := mark + 1; j < z.N+2; j++ {
			if stalls[j] && mark <= j {
				check(mark <= j)
				empty := j - mark
				middle := (mark + j - 1) / 2
				if empty > longestEmpty {
					longestEmpty = empty
					best = middle
				}
				mark = j + 1
				for mark < z.N+2 && stalls[mark] {
					mark++
				}
				if mark >= z.N+1 {
					break
				}
			}
		}
		check(!stalls[best])
		// log("   (Choosing", best, ")")
		stalls[best] = true
		// log(stalls)
	}
	// log(stalls)
	L, R := 0, 0
	for best-L > 0 && !stalls[best-L-1] {
		L++
	}
	for best+R < z.N+1 && !stalls[best+R+1] {
		R++
	}
	Y, Z := L, R
	if Y < Z {
		Y, Z = Z, Y
	}
	return fmt.Sprintf("%d %d", Y, Z)
}

func (z *Case) solveSmall2() interface{} {
	h := &IntHeap{z.N}
	heap.Init(h)
	var left, right int
	for i := 0; i < z.K; i++ {
		space := heap.Pop(h).(int)
		space-- // Consume for myself
		left = space / 2
		right = space - left
		if left > 0 {
			heap.Push(h, left)
		}
		if right > 0 {
			heap.Push(h, right)
		}
	}
	return fmt.Sprintf("%d %d", right, left)
}

/*
func (z *Case) solveLarge() interface{} {
	left, right := 0, 0
	k := z.K
	h := &PopuHeap{Popu{stalls: z.N, occurrences: 1}}
	heap.Init(h)
	for k > 0 {
		popu := heap.Pop(h).(Popu)
		k -= popu.occurrences
		left = (popu.stalls - 1) / 2
		right = popu.stalls - 1 - left
		if left == right {
			heap.Push(h, Popu{left, 2 * popu.occurrences})
		} else {
			heap.Push(h, Popu{left, popu.occurrences})
			heap.Push(h, Popu{right, popu.occurrences})
		}
	}
	return fmt.Sprintf("%d %d", right, left)
}
*/

func (z *Case) solveLarge() interface{} {
	left, right := 0, 0
	k := z.K
	h := &IntHeap{z.N}
	heap.Init(h)
	popu := map[int]int{z.N: 1}
	for k > 0 {
		stalls := heap.Pop(h).(int)
		occurrences := popu[stalls]
		k -= occurrences
		left = (stalls - 1) / 2
		right = stalls - 1 - left
		if _, ok := popu[left]; ok {
			popu[left] += occurrences
		} else {
			popu[left] = occurrences
			heap.Push(h, left)
		}
		if _, ok := popu[right]; ok {
			popu[right] += occurrences
		} else {
			popu[right] = occurrences
			heap.Push(h, right)
		}
	}
	return fmt.Sprintf("%d %d", right, left)
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

// An IntHeap is a max-heap of ints.
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// An IntHeap is a max-heap of ints.
type Popu struct {
	stalls      int
	occurrences int
}
type PopuHeap []Popu

func (h PopuHeap) Len() int           { return len(h) }
func (h PopuHeap) Less(i, j int) bool { return h[i].stalls > h[j].stalls }
func (h PopuHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *PopuHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(Popu))
}

func (h *PopuHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
