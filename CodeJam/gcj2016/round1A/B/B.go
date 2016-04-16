package main

// Google Code Jam 2016 Round 1A

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"
)

var (
	pathIn = "/home/val/Téléchargements/"
	// pathIn   = "./"
	pathOut  = "./"
	letter   = "B"
	strategy = (*Case).solveAfterMatch
	// strategy = (*Case).solveSmall
	// strategy = (*Case).solveLarge
	// strategy   = (*Case).solveLargeAndCheck
	concurrent = false
	maxProc    = 4
)

// Case data
type Case struct {
	caseNumber int
	// put fields here
	N     int
	lines [][]int
	grid  [][]int
}

func (z *Case) readSingle() {
	z.N = readInt()
	z.lines = make([][]int, 2*z.N-1)
	for i := range z.lines {
		z.lines[i] = make([]int, z.N)
		for j := range z.lines[i] {
			z.lines[i][j] = readInt()
		}
	}

	z.grid = make([][]int, z.N)
	for i := range z.grid {
		z.grid[i] = make([]int, z.N)
	}
}

func (z *Case) solveAfterMatch() interface{} {
	occ := map[int]int{}
	for _, line := range z.lines {
		for _, x := range line {
			occ[x]++
		}
	}

	odds := []int{}
	for v, count := range occ {
		if count%2 == 1 {
			odds = append(odds, v)
		}
	}
	sort.Ints(odds)
	s := fmt.Sprintf("%v", odds)
	s = strings.Replace(s, "[", "", 1)
	s = strings.Replace(s, "]", "", 1)
	return s
}

func (z *Case) solveSmall() interface{} {
	sort.Sort(List(z.lines))
	log(z.lines)
	missing := -1
	horiz := false

	p := 0
	for j, x := range z.lines[0] {
		z.grid[0][j] = x
	}
	p++
	if z.lines[0][0] == z.lines[1][0] {
		for i, x := range z.lines[1] {
			z.grid[i][0] = x
		}
		p++
	} else {
		missing = 0
		horiz = false
	}

	log(z.grid)
	// fmt.Sprintf(os.Stderr"%v", z.grid)

	for k := 1; k < z.N; k++ {
		okH := func(list []int) bool {
			for m := 0; m < k; m++ {
				if z.grid[k][m] != 0 && list[m] != z.grid[k][m] {
					return false
				}
			}
			return true
		}

		okV := func(list []int) bool {
			for m := 0; m < k; m++ {
				if z.grid[m][k] != 0 && list[m] != z.grid[m][k] {
					return false
				}
			}
			return true
		}

		if p+1 < len(z.lines) && z.lines[p][0] == z.lines[p+1][0] {
			if !okH(z.lines[p]) {
				z.lines[p], z.lines[p+1] = z.lines[p+1], z.lines[p]
			}

			if !okH(z.lines[p]) {
				panic(1)
			}

			for j, x := range z.lines[p] {
				z.grid[k][j] = x
			}
			p++

			if !okV(z.lines[p]) {
				panic(fmt.Sprintf("%v", z.grid))
			}

			for i, x := range z.lines[p] {
				z.grid[i][k] = x
			}
			p++
		} else {
			missing = k

			if okH(z.lines[p]) {
				for j, x := range z.lines[p] {
					z.grid[k][j] = x
				}
				horiz = false
				p++
			} else {
				if !okV(z.lines[p]) {
					panic(1)
				}
				horiz = true
				for i, x := range z.lines[p] {
					z.grid[i][k] = x
				}
				p++
			}
		}
	}
	answer := ""
	if horiz {
		for m := 0; m < z.N; m++ {
			answer += fmt.Sprintf(" %d", z.grid[missing][m])
		}
	} else {
		for m := 0; m < z.N; m++ {
			answer += fmt.Sprintf(" %d", z.grid[m][missing])
		}
	}
	return answer[1:]
}

type List [][]int

func (this List) Len() int           { return len(this) }
func (this List) Less(i, j int) bool { return this[i][0] < this[j][0] }
func (this List) Swap(i, j int)      { this[i], this[j] = this[j], this[i] }

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
