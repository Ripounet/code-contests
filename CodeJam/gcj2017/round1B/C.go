package main

// Google Code Jam 2017 Round 1B

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"runtime"
	"time"
)

var (
	pathIn = "/home/valou/Téléchargements/"
	// pathIn   = "./"
	pathOut = "./"
	letter  = "C"
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
	N, Q int
	E, S []int
	D    [][]int
	U, V []int
}

func (z *Case) readSingle() {
	z.N = readInt()
	z.Q = readInt()

	z.E = make([]int, z.N)
	z.S = make([]int, z.N)
	for i := 0; i < z.N; i++ {
		z.E[i] = readInt()
		z.S[i] = readInt()
	}

	z.D = make([][]int, z.N)
	for i := 0; i < z.N; i++ {
		z.D[i] = make([]int, z.N)
		for j := 0; j < z.N; j++ {
			z.D[i][j] = readInt()
		}
	}

	z.U = make([]int, z.Q)
	z.V = make([]int, z.Q)
	for i := 0; i < z.Q; i++ {
		z.U[i] = readInt()
		z.V[i] = readInt()
	}
}

func (z *Case) solveSmall() interface{} {

	var rec func(i, e, s int, t float64) (bool, float64)
	rec = func(i, e, s int, t float64) (bool, float64) {
		if i == z.N-1 {
			return true, t
		}
		j := i + 1
		d := z.D[i][j]
		e -= d
		if e < 0 {
			return false, -1
		}
		t += float64(d) / float64(s)

		// No switch
		oka, ta := rec(j, e, s, t)

		// Switch
		okb, tb := rec(j, z.E[j], z.S[j], t)

		switch {
		case !oka && !okb:
			return false, -1
		case !oka:
			return true, tb
		case !okb:
			return true, ta
		case ta <= tb:
			return true, ta
		case tb < ta:
			return true, tb
		}
		panic("ouch")
	}

	ok, t := rec(0, z.E[0], z.S[0], 0.0)
	check(ok)

	return fmt.Sprintf("%.7f", t)
}

// didn't work :(
func (z *Case) solveLarge() interface{} {
	m := make([][]float64, z.N)
	for i := 0; i < z.N; i++ {
		m[i] = make([]float64, z.N)
		for j := 0; j < z.N; j++ {
			if i != j {
				m[i][j] = math.Inf(1)
			}
		}
	}

	for i := 0; i < z.N; i++ {
		z.Bfs(i, m)
	}
	log(matstr(m))

	total := make([][]float64, z.N)
	for i := 0; i < z.N; i++ {
		total[i] = make([]float64, z.N)
		for j := 0; j < z.N; j++ {
			if i != j {
				total[i][j] = math.Inf(1)
			}
		}
	}

	for i := 0; i < z.N; i++ {
		z.Bfs2(i, m, total)
	}
	log(matstr(total))

	var result bytes.Buffer

	for k := 0; k < z.Q; k++ {
		u, v := z.U[k]-1, z.V[k]-1
		fmt.Fprint(&result, fmt.Sprintf("%.7f", total[u][v]), " ")
	}

	return result.String()
}

func (z *Case) Bfs(start int, m [][]float64) {
	queue := []int{start}
	for len(queue) > 0 {
		i := queue[0]
		queue = queue[1:]

		t := m[start][i]

		for j := 0; j < z.N; j++ {
			if j == i {
				continue
			}
			nexttime := float64(z.D[i][j]) / float64(z.S[start])
			if z.D[i][j] != -1 && t+nexttime < m[start][j] {
				m[start][j] = t + nexttime
				queue = append(queue, j)
			}
		}
	}
}

func (z *Case) Bfs2(start int, m, total [][]float64) {
	queue := []int{start}
	for len(queue) > 0 {
		i := queue[0]
		queue = queue[1:]

		t := m[start][i]

		for j := 0; j < z.N; j++ {
			if j == i {
				continue
			}
			nexttime := m[i][j]
			if z.D[i][j] != -1 && t+nexttime < total[start][j] {
				log(start, i, j, ":", t+nexttime, "<", total[start][j])
				total[start][j] = t + nexttime
				queue = append(queue, j)
			}
		}
	}
}

func matstr(mat [][]float64) string {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, "[")
	for _, line := range mat {
		fmt.Fprintln(&buf, " ", line)
	}
	fmt.Fprintln(&buf, "]")
	return buf.String()
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
