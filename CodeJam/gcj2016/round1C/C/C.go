package main

// Google Code Jam 2016 Round 1C

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	pathIn = "/home/val/Téléchargements/"
	// pathIn   = "./"
	pathOut  = "./"
	letter   = "C"
	strategy = (*Case).solveSmall
	// strategy = (*Case).solveLarge
	// strategy   = (*Case).solveLargeAndCheck
	concurrent = false
	maxProc    = 4
)

// Put inputs here as global vars
var (
	J, P, S, K int
)

// Case data
type Case struct {
	caseNumber int
}

func (z *Case) readSingle() {
	check(!concurrent) // So we can use global vars
	J = readInt()
	P = readInt()
	S = readInt()
	K = readInt()
}

func (z *Case) solveSmall() interface{} {
	bestY := 0
	bestList := make([]string, 27)
	//list := make([]string, 27)
	var answer string = "???"

	type worn [27]bool
	type pairWorn [9]int

	expand := func(worn worn) string {
		str := ""
		for combi, w := range worn {
			if !w {
				continue
			}
			j, p, s := (combi / 9), ((combi / 3) % 3), (combi % 3)
			str = str + fmt.Sprintf("[%d %d %d]", 1+j, 1+p, 1+s)
		}
		return str
	}
	_ = expand

	var rec func(worn, pairWorn, pairWorn, pairWorn, int, int)
	rec = func(worn worn, jp, ps, js pairWorn, days int, last int) {
		//log(expand(worn))
		check(days < 28)
		if days > bestY {
			bestY = days
			//copy(bestList, list[:days])
			d := 0
			for combi, w := range worn {
				if !w {
					continue
				}
				j, p, s := (combi / 9), ((combi / 3) % 3), (combi % 3)
				bestList[d] = fmt.Sprintf("%d %d %d", 1+j, 1+p, 1+s)
				d++
			}

			answer = fmt.Sprintf("%d\n", bestY)
			answer += strings.Join(bestList[:bestY], "\n")
		}

		for combi := last + 1; combi < 27; combi++ {
			//log(combi)
			if worn[combi] {
				//log(combi, "already worn")
				continue // can't...
			}
			j, p, s := (combi / 9), ((combi / 3) % 3), (combi % 3)
			if j > J-1 || p > P-1 || s > S-1 {
				//log(j, p, s, "invalid")
				continue // invalid
			}
			ijp := 3*j + p
			ips := 3*p + s
			ijs := 3*j + s
			if jp[ijp] == K {
				//log(combi, jp, "jp over K")
				continue // can't...
			}
			if ps[ips] == K {
				//log(combi, ps, "ps over K")
				continue // can't...
			}
			if js[ijs] == K {
				//log(combi, js, "js over K")
				continue // can't...
			}
			nextWorn := worn
			nextWorn[combi] = true
			nextjp := jp
			nextjp[ijp]++
			nextps := ps
			nextps[ips]++
			nextjs := js
			nextjs[ijs]++
			rec(nextWorn, nextjp, nextps, nextjs, days+1, combi)
		}
	}
	rec(worn{}, pairWorn{}, pairWorn{}, pairWorn{}, 0, -1)

	check(answer != "???")
	return answer
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
