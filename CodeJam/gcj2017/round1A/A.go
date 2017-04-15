package main

// Google Code Jam 2017 Round 1A

import (
	"bytes"
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
	R, C    int
	cake    [][]byte
	letters map[byte]bool
}

func (z *Case) readSingle() {
	z.letters = map[byte]bool{}
	z.R = readInt()
	z.C = readInt()
	z.cake = make([][]byte, z.R)
	for i := range z.cake {
		line := readString()
		z.cake[i] = []byte(line)
		for _, letter := range line {
			z.letters[byte(letter)] = true
		}
	}
}

func (z *Case) solveSmall() interface{} {
	// Fill non-empty rows
	for i := 0; i < z.R; i++ {
		var last byte = '?'
		for j := 0; j < z.C; j++ {
			if z.cake[i][j] == '?' {
				z.cake[i][j] = last
			}
			last = z.cake[i][j]
		}
		last = z.cake[i][z.C-1]
		for j := z.C - 2; j >= 0; j-- {
			if z.cake[i][j] == '?' {
				z.cake[i][j] = last
			}
			last = z.cake[i][j]
		}
	}

	// Fill empty rows
	for i := 1; i < z.R; i++ {
		if z.cake[i][0] == '?' && z.cake[i-1][0] != '?' {
			for j := 0; j < z.C; j++ {
				z.cake[i][j] = z.cake[i-1][j]
			}
		}
	}
	for i := z.R - 2; i >= 0; i-- {
		if z.cake[i][0] == '?' && z.cake[i+1][0] != '?' {
			for j := 0; j < z.C; j++ {
				z.cake[i][j] = z.cake[i+1][j]
			}
		}
	}

	sb := new(bytes.Buffer)
	for i := 0; i < z.R; i++ {
		fmt.Fprintln(sb)
		fmt.Fprint(sb, string(z.cake[i]))
	}
	return sb.String()
}

func (z *Case) solveLarge() interface{} {
	return ""
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
