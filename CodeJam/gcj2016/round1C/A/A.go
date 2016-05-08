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
	pathIn = "/home/val/Téléchargements/"
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
	N     int
	P     []int
	total int
)

// Case data
type Case struct {
	caseNumber int
}

func (z *Case) readSingle() {
	check(!concurrent) // So we can use global vars
	N = readInt()
	P = make([]int, N)
	total = 0
	for i := range P {
		P[i] = readInt()
		total += P[i]
	}
}

func (z *Case) solveSmall() interface{} {
	next := func() (int, string) {
		i, i2 := -1, -1
		high := 0
		present := 0
		for p, n := range P {
			present += n
			if n > high {
				high = n
				i = p
				i2 = -1
				continue
			}
			if i != -1 && n > 0 && n == high && i2 == -1 {
				i2 = p
			}
		}

		check(i >= 0)
		P[i]--
		if i2 == -1 || (high == 1 && present == 3) {
			return 1, fmt.Sprintf("%c", 'A'+byte(i))
		}
		P[i2]--
		return 2, fmt.Sprintf("%c%c", 'A'+byte(i), 'A'+byte(i2))
	}

	var plan bytes.Buffer
	evac := 0
	for evac < total {
		//log(evac, "/", total)
		m, str := next()
		evac += m
		if evac == total {
			plan.WriteString(str)
		} else {
			plan.WriteString(str + " ")
		}
	}

	return plan.String()
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
