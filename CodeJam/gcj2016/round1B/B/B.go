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
	pathOut = "./"
	letter  = "B"
	//strategy = (*Case).solveAfterMatch
	//strategy = (*Case).solveSmall
	strategy = (*Case).solveLarge
	//strategy   = (*Case).solveLargeAndCheck
	concurrent = false
	maxProc    = 4
)

// Case data
type Case struct {
	caseNumber int
	// put fields here
	C, J string
}

func (z *Case) readSingle() {
	z.C = readString()
	z.J = readString()
	if len(z.C) != len(z.J) {
		panic("gder")
	}
}

func (z *Case) solveSmall() interface{} {
	var ac, aj [1000]bool
	var rec func(string, int, *[1000]bool)
	rec = func(s string, accu int, a *[1000]bool) {
		if s == "" {
			a[accu] = true
			return
		}
		if s[0] == '?' {
			for k := 0; k <= 9; k++ {
				rec(s[1:], 10*accu+k, a)
			}
		} else {
			k := int(s[0] - '0')
			rec(s[1:], 10*accu+k, a)
		}
	}
	rec(z.C, 0, &ac)
	rec(z.J, 0, &aj)

	bestDiff, bestC, bestJ := 9999, -1, -1

	for x := 0; x < 1000; x++ {
		if !ac[x] {
			continue
		}
		for y := 0; y < 1000; y++ {
			if !aj[y] {
				continue
			}
			if abs(x-y) < bestDiff {
				bestDiff, bestC, bestJ = abs(x-y), x, y
				continue
			}
			if abs(x-y) == bestDiff {
				if x < bestC {
					bestDiff, bestC, bestJ = abs(x-y), x, y
					continue
				}
				if x == bestC && y < bestJ {
					bestDiff, bestC, bestJ = abs(x-y), x, y
					continue
				}
			}
		}
	}

	N := fmt.Sprintf("%d", len(z.C))
	return fmt.Sprintf("%0"+N+"d %0"+N+"d", bestC, bestJ)
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (z *Case) solveLarge() interface{} {
	N, C, J := len(z.C), z.C, z.J

	bestDiff, bestC, bestJ := 999999999999999999, -1, -1

	var rec func(int, int, int)
	rec = func(i, c, j int) {
		if c < 0 {
			panic(c)
		}
		if j < 0 {
			panic(j)
		}
		if i == N {
			x, y := c, j
			if abs(x-y) < bestDiff {
				bestDiff, bestC, bestJ = abs(x-y), x, y
				return
			}
			if abs(x-y) == bestDiff {
				if x < bestC {
					bestDiff, bestC, bestJ = abs(x-y), x, y
					return
				}
				if x == bestC && y < bestJ {
					bestDiff, bestC, bestJ = abs(x-y), x, y
					return
				}
			}
			return
		}

		// ---

		ci, ji := C[i], J[i]

		if ci == '?' && ji == '?' {
			if c < j {
				rec(i+1, 10*c+9, 10*j+0)
			}
			if c > j {
				rec(i+1, 10*c+0, 10*j+9)
			}
			if c == j {
				rec(i+1, 10*c+0, 10*j+0)
				rec(i+1, 10*c+0, 10*j+1)
				rec(i+1, 10*c+1, 10*j+0)
			}
			return
		}

		if ci == '?' {
			J := int(ji - '0')
			if c < j {
				rec(i+1, 10*c+9, 10*j+J)
			}
			if c > j {
				rec(i+1, 10*c+0, 10*j+J)
			}
			if c == j {
				if J > 0 {
					rec(i+1, 10*c+J-1, 10*j+J)
				}
				rec(i+1, 10*c+J, 10*j+J)
				if J < 9 {
					rec(i+1, 10*c+J+1, 10*j+J)
				}
			}
			return
		}

		if ji == '?' {
			C := int(ci - '0')
			if c < j {
				rec(i+1, 10*c+C, 10*j+0)
			}
			if c > j {
				rec(i+1, 10*c+C, 10*j+9)
			}
			if c == j {
				if C > 0 {
					rec(i+1, 10*c+C, 10*j+C-1)
				}
				rec(i+1, 10*c+C, 10*j+C)
				if C < 9 {
					rec(i+1, 10*c+C, 10*j+C+1)
				}
			}
			return
		}

		C := int(ci - '0')
		J := int(ji - '0')
		rec(i+1, 10*c+C, 10*j+J)
	}
	rec(0, 0, 0)

	Ns := fmt.Sprintf("%d", N)
	return fmt.Sprintf("%0"+Ns+"d %0"+Ns+"d", bestC, bestJ)
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
