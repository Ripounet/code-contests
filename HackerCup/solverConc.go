package main

// Facebook Hacker Cup XXX

import (
	"fmt"
	"os"
	"time"
)

var (
	pathIn = "/home/valou/Téléchargements/"
	// pathIn   = "./"
	pathOut    = "./"
	letter     = "X"
	concurrent = true
)

// Global vars, for precomputed data only.
var (
	facto = map[int]int{
		0: 1,
	}
)

// Case data
type Case struct {
	caseNumber int
	solving    chan<- Ø
}

func (z *Case) solve() interface{} {
	//
	// Read case
	//
	N := readInt()

	//
	// Solve
	//
	z.solving <- ø

	// Implement closures that reference inputs, if aux
	// funcs are needed.

	return 0 * N
}

// Global precomputed data (if needed)

func precompute() {
	for i := 1; i < 3000; i++ {
		facto[i] = (facto[i-1] * i) % 1000000007
	}
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
	return z.solve()
}

func solve() {
	top1 := time.Now()
	precompute()

	T := readInt()
	solutions := make([]chan interface{}, T)
	for i := range solutions {
		solutions[i] = make(chan interface{})
	}
	go func() {
		for i := range solutions {
			// Chan just to be sure inputs are read in
			// original order.
			solvingch := make(chan Ø, 1)

			currentCase := &Case{
				caseNumber: 1 + i,
				solving:    solvingch,
			}
			if concurrent {
				go func(ch chan interface{}) {
					ch <- currentCase.solveSingle()
					close(ch)
				}(solutions[i])
				<-solvingch
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
	var fileIn, fileOut string
	if len(os.Args) < 2 {
		fileIn = pathIn + letter + ".in"
		fileOut = pathOut + letter + ".out"
	} else {
		sample := os.Args[1]
		fileIn = pathIn + letter + "-" + sample + ".in"
		fileOut = pathOut + letter + "-" + sample + ".out"
	}

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

type Ø struct{}

var ø = Ø{}
