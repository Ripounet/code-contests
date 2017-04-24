package main

// Google Code Jam 2017 Round 1B

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	pathIn = "/home/valou/Téléchargements/"
	// pathIn   = "./"
	pathOut = "./"
	letter  = "B"
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
	N, R, O, Y, G, B, V int
}

func (z *Case) readSingle() {
	z.N = readInt()
	z.R = readInt()
	z.O = readInt()
	z.Y = readInt()
	z.G = readInt()
	z.B = readInt()
	z.V = readInt()
}

func (z *Case) solveSmall() interface{} {
	count := []int{z.R, z.Y, z.B}
	letter := []string{"R", "Y", "B"}
	// sort desc!!
	if count[0] < count[1] {
		count[0], count[1] = count[1], count[0]
		letter[0], letter[1] = letter[1], letter[0]
	}
	if count[0] < count[2] {
		count[0], count[2] = count[2], count[0]
		letter[0], letter[2] = letter[2], letter[0]
	}
	if count[1] < count[2] {
		count[1], count[2] = count[2], count[1]
		letter[1], letter[2] = letter[2], letter[1]
	}

	// Now it's like r>y>b
	r, y, b := count[0], count[1], count[2]
	log("count ", r, y, b)

	// Corner case??  Should not happen because N>=3
	if r == 1 && y == 0 && b == 0 {
		return letter[0]
	}

	if r > y+b {
		return "IMPOSSIBLE"
	}

	v := (y + b) - r
	w := r - y
	var result bytes.Buffer
	for i := 0; i < v; i++ {
		fmt.Fprint(&result, letter[0], letter[1], letter[2])
	}
	for i := 0; i < y-v; i++ {
		fmt.Fprint(&result, letter[0], letter[1])
	}
	for i := 0; i < w; i++ {
		fmt.Fprint(&result, letter[0], letter[2])
	}
	//check(result.Len() == z.N)
	return result.String()
}

// didn't work :(
func (z *Case) solveLarge() interface{} {
	// Eliminate bicolors, then solve R,Y,B a small input!

	if z.G+z.R == z.N {
		if z.G == z.R {
			return strings.Repeat("GR", z.G)
		}
		return "IMPOSSIBLE"
	}

	if z.O+z.B == z.N {
		if z.O == z.B {
			return strings.Repeat("OB", z.O)
		}
		return "IMPOSSIBLE"
	}

	if z.V+z.Y == z.N {
		if z.V == z.Y {
			return strings.Repeat("VY", z.V)
		}
		return "IMPOSSIBLE"
	}

	if z.O > z.B-1 {
		return "IMPOSSIBLE"
	}
	if z.V > z.Y-1 {
		return "IMPOSSIBLE"
	}
	if z.G > z.R-1 {
		return "IMPOSSIBLE"
	}

	z.R -= z.G
	z.Y -= z.V
	z.B -= z.O
	s := z.solveSmall().(string)
	z.R += z.G
	z.Y += z.V
	z.B += z.O
	if s == "IMPOSSIBLE" {
		return "IMPOSSIBLE"
	}

	s = strings.Replace(s, "B", "B"+strings.Repeat("OB", z.O), 1)
	s = strings.Replace(s, "Y", "Y"+strings.Repeat("VY", z.V), 1)
	s = strings.Replace(s, "R", "R"+strings.Repeat("GR", z.G), 1)

	// log("s =", s)
	// log("len(s) =", len(s))
	// log("N =", z.N)
	check(len(s) == z.N)
	check(!strings.Contains(s, "RR"))
	check(!strings.Contains(s, "RO"))
	check(!strings.Contains(s, "OR"))
	check(!strings.Contains(s, "RV"))
	check(!strings.Contains(s, "VR"))

	check(!strings.Contains(s, "BB"))
	check(!strings.Contains(s, "BV"))
	check(!strings.Contains(s, "VB"))
	check(!strings.Contains(s, "BG"))
	check(!strings.Contains(s, "GB"))

	check(!strings.Contains(s, "YY"))
	check(!strings.Contains(s, "YO"))
	check(!strings.Contains(s, "OY"))
	check(!strings.Contains(s, "YG"))
	check(!strings.Contains(s, "GY"))

	check(!strings.Contains(s, "OO"))
	check(!strings.Contains(s, "VV"))
	check(!strings.Contains(s, "GG"))

	check(strings.Count(s, "R") == z.R)
	check(strings.Count(s, "Y") == z.Y)
	check(strings.Count(s, "B") == z.B)
	check(strings.Count(s, "O") == z.O)
	check(strings.Count(s, "V") == z.V)
	check(strings.Count(s, "G") == z.G)

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
