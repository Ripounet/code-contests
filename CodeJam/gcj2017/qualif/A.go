package main

// Google Code Jam 2017 Qualification round

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

var (
	pathIn = "/home/valou/Téléchargements/"
	// pathIn   = "./"
	pathOut = "./"
	letter  = "A"
	// strategy = (*Case).solveSmall
	strategy = (*Case).solveLarge
	// strategy   = (*Case).solveLargeAndCheck
	concurrent = false
	maxProc    = 4
)

// Case data
type Case struct {
	caseNumber int
	// put fields here
	S string
	K int
}

func (z *Case) readSingle() {
	z.S = readString()
	z.K = readInt()
}

func allHappy(s string) bool {
	for _, c := range s {
		if c != '+' {
			return false
		}
	}
	return true
}

func flip(s string, i, j int) string {
	t := []byte(s)
	for k := i; k < j; k++ {
		t[k] ^= ('+' ^ '-')
	}
	return string(t)
}

func (z *Case) solveSmall() interface{} {
	// BFS
	seen := map[string]int{z.S: 0}
	q := []string{z.S}

	for len(q) > 0 {
		s := q[0]
		q = q[1:]
		uses := seen[s]
		if allHappy(s) {
			return uses
		}

		for i := 0; i+z.K <= len(s); i++ {
			next := flip(s, i, i+z.K)
			if _, already := seen[next]; !already {
				seen[next] = uses + 1
				q = append(q, next)
			}
		}
	}

	return "IMPOSSIBLE"
}

func (z *Case) solveLarge() interface{} {
	row := []byte(z.S)
	N := len(row)
	uses := 0
	for i := 0; i < N-z.K+1; i++ {
		if row[i] != '+' {
			uses++
			for j := i; j < i+z.K; j++ {
				row[j] ^= ('+' ^ '-')
			}
		}
	}
	log(string(row))
	for i := N - z.K + 1; i < len(z.S); i++ {
		if row[i] != '+' {
			return "IMPOSSIBLE"
		}
	}
	return uses
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
