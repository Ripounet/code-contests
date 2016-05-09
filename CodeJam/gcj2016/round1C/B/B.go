package main

// Google Code Jam 2016 Round 1C

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"time"
)

var (
	pathIn = "/home/valentin/Downloads/"
	// pathIn   = "./"
	pathOut = "./"
	letter  = "B"
	//strategy = (*Case).solveSmall
	strategy = (*Case).solveLarge
	// strategy   = (*Case).solveLargeAndCheck
	concurrent = false
	maxProc    = 4
)

// Put inputs here as global vars
var (
	B, M int
	mat  Matrix
	seen [2000000]bool
)

// Case data
type Case struct {
	caseNumber int
}

func (z *Case) readSingle() {
	check(!concurrent) // So we can use global vars
	B = readInt()
	M = readInt()
	seen = [2000000]bool{}
	mat = Matrix{}
}

type Matrix [50][50]bool

func (m *Matrix) print() string {
	var buf bytes.Buffer
	for i := 0; i < B; i++ {
		for j := 0; j < B; j++ {
			c := "0"
			if m[i][j] {
				c = "1"
			}
			buf.WriteString(c)
		}
		if i < B-1 {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

func (m *Matrix) count() int {
	var c [50]int
	c[0] = 1
	for i := 0; i < B-1; i++ {
		for j := i; j < B; j++ {
			if m[i][j] {
				c[j] += c[i]
			}
		}
	}
	return c[B-1]
}

func (z *Case) solveSmall() interface{} {

	var rec func(int, int) string
	rec = func(i, j int) string {
		if i == B-1 {
			if mat.count() == M {
				return "POSSIBLE\n" + mat.print()
			}
			return ""
		}

		var ii, jj int
		if j == B-1 {
			ii, jj = i+1, i+2
		} else {
			ii, jj = i, j+1
		}
		a := rec(ii, jj)
		if a != "" {
			return a
		}

		mat[i][j] = true
		a = rec(ii, jj)
		if a != "" {
			return a
		}
		mat[i][j] = false
		return ""
	}
	answer := rec(0, 1)
	if answer == "" {
		return "IMPOSSIBLE"
	}
	return answer
}

func (z *Case) solveLarge() interface{} {
	if M > pow(2, B-2) {
		return "IMPOSSIBLE"
	}

	for k := 0; k < B-1; k++ {
		mat[k][k+1] = true
	}
	// Now the "diagonal tangent" contains 1 path 1->2->...->B
	left := M - 1
	// Build full triangle : contains 2^X paths to "next building"
	p := 1
	j := 2
	for ; j < B && left >= p; j++ {
		for i := 0; i < j-1; i++ {
			mat[i][j] = true
		}
		left -= p
		p *= 2
	}
	// Build last, partial column
	if left > 0 {
		for i := j - 2; i >= 0; i-- {
			if left >= pow(2, i-1) {
				mat[i][j] = true
				left -= pow(2, i-1)
			} else {
			}
		}
	}
	check(left == 0)

	return "POSSIBLE\n" + mat.print()
}

func pow(a, b int) int {
	n := 1
	for i := 0; i < b; i++ {
		n *= a
	}
	return n
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
