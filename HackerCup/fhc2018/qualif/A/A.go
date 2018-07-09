package main

// Facebook Hacker Cup 2018 qualif

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	// pathIn = "/home/valou/Téléchargements/"
	pathIn     = "./"
	pathOut    = "./"
	letter     = "A"
	concurrent = false
)

// Put inputs here as global vars
var (
	N, K, V int
	A       []string
	times   []int
)

// Case data
type Case struct {
	caseNumber int
}

func (z *Case) readSingle() {
	check(!concurrent) // So we can use global vars
	N = readInt()
	K = readInt()
	V = readInt()
	A = make([]string, N)
	for i := range A {
		A[i] = readString()
	}
	times = make([]int, N)
}

func (z *Case) solve() interface{} {
	return z.large()
}

func (z *Case) naive() interface{} {
	return 0
}

func (z *Case) large() interface{} {
	d := ((V - 1) * K) % N
	var b bytes.Buffer
	for i := N; i < d+K; i++ {
		fmt.Fprintf(&b, "%s ", A[i%N])
	}
	for i := d; i < d+K && i < N; i++ {
		fmt.Fprintf(&b, "%s ", A[i])
	}
	s := b.String()
	s = strings.TrimSpace(s)
	return s
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
	logf("Usage: %v [sample] \n", os.Args[0])
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

func readInt64() int64 {
	var i int64
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
