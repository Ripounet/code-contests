package main

// Google Code Jam 2016 Round 1B

import (
	"fmt"
	"os"
	"runtime"
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

// Case data
type Case struct {
	caseNumber int
	// put fields here
	S string
}

func (z *Case) readSingle() {
	z.S = readString()
}

func (z *Case) solveSmall() interface{} {
	var digits [10]int
	count := map[rune]int{}

	for _, c := range z.S {
		count[c]++
	}

	words := []string{"ZERO", "ONE", "TWO", "THREE", "FOUR", "FIVE", "SIX", "SEVEN", "EIGHT", "NINE"}

	var n int

	// 0
	n = count['Z']
	digits[0] = n
	for _, c := range words[0] {
		count[c] -= n
	}

	// 2
	n = count['W']
	digits[2] = n
	for _, c := range words[2] {
		count[c] -= n
	}

	// 8
	n = count['G']
	digits[8] = n
	for _, c := range words[8] {
		count[c] -= n
	}

	// 6
	n = count['X']
	digits[6] = n
	for _, c := range words[6] {
		count[c] -= n
	}

	// 4
	n = count['U']
	digits[4] = n
	for _, c := range words[4] {
		count[c] -= n
	}

	// 7
	n = count['S']
	digits[7] = n
	for _, c := range words[7] {
		count[c] -= n
	}

	// 5
	n = count['V']
	digits[5] = n
	for _, c := range words[5] {
		count[c] -= n
	}

	// 3
	n = count['H']
	digits[3] = n
	for _, c := range words[3] {
		count[c] -= n
	}

	// 9
	n = count['I']
	digits[9] = n
	for _, c := range words[9] {
		count[c] -= n
	}

	// 1
	n = count['O']
	digits[1] = n
	for _, c := range words[1] {
		count[c] -= n
	}

	for c := rune('A'); c <= 'Z'; c++ {
		check(count[c] == 0)
	}

	s := ""
	for i, n := range digits {
		for j := 0; j < n; j++ {
			s += fmt.Sprintf("%d", i)
		}
	}
	return s
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
