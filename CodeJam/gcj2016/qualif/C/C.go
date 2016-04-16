package main

// Google Code Jam 2016 Qualif

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

var (
	pathIn = "/home/valou/Téléchargements/"
	// pathIn   = "./"
	pathOut = "./"
	letter  = "C"
	// strategy = (*Case).solveSmall
	// strategy = (*Case).solveLarge
	strategy = (*Case).solveLargeGenerate
	// strategy   = (*Case).solveLargeAndCheck
	concurrent = false
	maxProc    = 4
)

// Case data
type Case struct {
	caseNumber int
	// put fields here
	N, J int
}

func (z *Case) readSingle() {
	z.N = readInt()
	z.J = readInt()
}

func (z *Case) solveSmall() interface{} {
	type coin [32]byte
	result := ""
	memo := map[string]bool{}
	var divisors [11]int
	for i := 0; i < z.J; i++ {
	search:
		for {
			var x coin
			x[0], x[z.N-1] = '1', '1'
			for k := 1; k < z.N-1; k++ {
				x[k] = '0' + byte(rand.Intn(2))
			}
			s := string(x[:z.N])
			if memo[s] {
				continue search
			}
			// log("s =", s)
		bases:
			for base := 2; base <= 10; base++ {
				t64, _ := strconv.ParseInt(s, base, 0)
				t := int(t64)
				// log("  t =", t)
				sq := int(math.Sqrt(float64(t)))
				for w := 0; w < 100000; w++ {
					divisor := 2 + rand.Intn(sq-1)
					if t%divisor == 0 {
						divisors[base] = divisor
						continue bases
					}
				}
				// divisor not found :(
				//log("Divisor not found, probable prime", s, "in base", base)
				continue search
			}
			// Compound Found :)
			log("Found ", s)
			memo[s] = true
			result += "\n" + s
			for _, divisor := range divisors[2:] {
				result += fmt.Sprintf(" %d", divisor)
			}
			break search
		}
	}
	return result
}

func (z *Case) solveLarge() interface{} {
	type coin [32]byte
	result := ""
	memo := map[string]bool{}
	var divisors [11]int
	for i := 0; i < z.J; i++ {
	search:
		for {
			var x coin
			bi := new(big.Int)
			rem := new(big.Int)
			two := big.NewInt(2)
			//sqr := big.NewInt(2)
			x[0], x[z.N-1] = '1', '1'
			for k := 1; k < z.N-1; k++ {
				x[k] = '0' + byte(rand.Intn(2))
			}
			s := string(x[:z.N])
			if memo[s] {
				continue search
			}
		bases:
			for base := 2; base <= 10; base++ {
				bi.SetString(s, base)
				if !bi.ProbablyPrime(200) {
					rem.Mod(bi, two)
					//				for divisor:=3;divisor*divisor<=
					//				divisors[base] = 9 //divisor
					continue bases
				}
				// divisor not found :(
				//log("Divisor not found, probable prime", s, "in base", base)
				continue search
			}
			// Compound Found :)
			log("Found ", s)
			memo[s] = true
			result += "\n" + s
			for _, divisor := range divisors[2:] {
				result += fmt.Sprintf(" %d", divisor)
			}
			break search
		}
	}
	return result
}

func (z *Case) solveLargeGenerate() interface{} {
	// Observe that :
	// ABCDABCD == 10001 * ABCD is never prime
	// 1XXXXXXXXXXXXXX11XXXXXXXXXXXXXX1 == 10000000000000001 * 1XXXXXXXXXXXXXX1 is never prime
	result := ""
	j := 0

	// Always the same "non-trivial" divisors!
	divisors := ""
	bid := new(big.Int)
	for base := 2; base <= 10; base++ {
		bid.SetString("10000000000000001", base)
		divisors = fmt.Sprintf("%v %d", divisors, bid)
	}

	for b9 := '0'; b9 <= '1'; b9++ {
		for b8 := '0'; b8 <= '1'; b8++ {
			for b7 := '0'; b7 <= '1'; b7++ {
				for b6 := '0'; b6 <= '1'; b6++ {
					for b5 := '0'; b5 <= '1'; b5++ {
						for b4 := '0'; b4 <= '1'; b4++ {
							for b3 := '0'; b3 <= '1'; b3++ {
								for b2 := '0'; b2 <= '1'; b2++ {
									for b1 := '0'; b1 <= '1'; b1++ {
										result += fmt.Sprintf("\n100000%c%c%c%c%c%c%c%c%c1100000%c%c%c%c%c%c%c%c%c1%v", b9, b8, b7, b6, b5, b4, b3, b2, b1, b9, b8, b7, b6, b5, b4, b3, b2, b1, divisors)
										//log("Found ", s)
										j++
										if j == z.J {
											return result
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	panic("nope")
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
