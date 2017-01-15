package main

// Facebook Hacker Cup 2017 Round 1

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"time"
)

var (
	pathIn = "/home/valou/Téléchargements/"
	// pathIn   = "./"
	pathOut  = "./"
	letter   = "A"
	strategy = (*Case).solveSmall
	// strategy = (*Case).solveLarge
	// strategy   = (*Case).solveLargeAndCheck
	concurrent = false
	maxProc    = 4
)

// Put inputs here as global vars
var (
	N, M int
	C    [][]int
)

// Case data
type Case struct {
	caseNumber int
}

func (z *Case) readSingle() {
	check(!concurrent) // So we can use global vars
	N = readInt()
	M = readInt()
	C = make([][]int, N)
	for i := range C {
		C[i] = make([]int, M)
		for j := range C[i] {
			C[i][j] = readInt()
		}
		// Cheapest first please!
		sort.Ints(C[i])
	}
}

func (z *Case) solveSmall() interface{} {
	// best[i][p] is "Cheapest way at end of day (after dinner) i to have p pies"
	var best [300][301]int
	for i := range best {
		for p := range best[i] {
			best[i][p] = -1
		}
	}

	// First day
	piescost := 0
	for p := 1; p <= M; p++ {
		piescost += C[0][p-1]
		tax := p * p
		best[0][p-1] = piescost + tax
	}

	for i := 1; i < N; i++ {
		// p==0
		for k := 1; k < 301 && best[i-1][k] != -1; k++ {
			best[i][k-1] = best[i-1][k]
		}

		piescost := 0
		for p := 1; p <= M; p++ {
			piescost += C[i][p-1]
			tax := p * p

			for k := 0; k+p-1 < 301 && best[i-1][k] != -1; k++ {
				candidate := best[i-1][k] + piescost + tax
				if best[i][k+p-1] == -1 || candidate < best[i][k+p-1] {
					best[i][k+p-1] = candidate
				}
			}
			// No point in having more than 300 pies.
		}
	}

	// Cheapest path ends with exactly 0 pie after last supper.
	return best[N-1][0]
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
