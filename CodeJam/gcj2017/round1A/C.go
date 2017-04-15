package main

// Google Code Jam 2017 Round 1A

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
	letter   = "C"
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
	Hd, Ad, Hk, Ak, B, D int
}

func (z *Case) readSingle() {
	z.Hd = readInt()
	z.Ad = readInt()
	z.Hk = readInt()
	z.Ak = readInt()
	z.B = readInt()
	z.D = readInt()
}

type state struct {
	Hd, Ad, Hk, Ak int
}

type stateT struct {
	state
	turns int
}

// func (s state) worst(other state) bool {
// 	return
// }

func (z *Case) solveSmall() interface{} {
	initial := stateT{state{z.Hd, z.Ad, z.Hk, z.Ak}, 0}
	seen := map[state]bool{
		initial.state: true,
	}
	q := []stateT{initial}
	for len(q) > 0 {
		st := q[0]
		q = q[1:]
		// log(st)

		// Attack
		{
			st := st
			st.Hk -= st.Ad
			if st.Hk <= 0 {
				return st.turns + 1
			}
			st.Hd -= st.Ak
			if st.Hd > 0 {
				st.turns++
				if !seen[st.state] {
					seen[st.state] = true
					q = append(q, st)
				}
			}
		}

		// Buff
		{
			st := st
			st.Ad += z.B
			st.Hd -= st.Ak
			if st.Hd > 0 {
				st.turns++
				if !seen[st.state] {
					seen[st.state] = true
					q = append(q, st)
				}
			}
		}

		// Cure
		{
			st := st
			st.Hd = initial.Hd
			st.Hd -= st.Ak
			if st.Hd > 0 {
				st.turns++
				if !seen[st.state] {
					seen[st.state] = true
					q = append(q, st)
				}
			}
		}

		// Debuff
		{
			st := st
			st.Ak -= z.D
			if st.Ak < 0 {
				st.Ak = 0
			}
			st.Hd -= st.Ak
			if st.Hd > 0 {
				st.turns++
				if !seen[st.state] {
					seen[st.state] = true
					q = append(q, st)
				}
			}
		}
	}
	return "IMPOSSIBLE"
}

func (z *Case) solveLarge() interface{} {
	return "IMPOSSIBLE"
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
