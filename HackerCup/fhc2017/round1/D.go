package main

// Facebook Hacker Cup 2017 Round 1

// Correct answers but took 1000s to compute -> expired!
// Should have gone multicore for this one.

import (
	"fmt"
	"math/big"
	"os"
	"runtime"
	"time"
)

var (
	pathIn = "/home/valou/Téléchargements/"
	// pathIn   = "./"
	pathOut  = "./"
	letter   = "D"
	strategy = (*Case).solveSmall
	// strategy = (*Case).solveLarge
	// strategy   = (*Case).solveLargeAndCheck
	concurrent = false
	maxProc    = 4
)

// Put inputs here as global vars
var (
	N, M int
	R    []int

	facto = map[int]int{
		0: 1,
	}

	billionSeven = big.NewInt(1000000007)
)

// Case data
type Case struct {
	caseNumber int
}

func (z *Case) readSingle() {
	check(!concurrent) // So we can use global vars
	N = readInt()
	M = readInt()
	R = make([]int, N)
	for i := range R {
		R[i] = readInt()
	}
}

func (z *Case) solveSmall() interface{} {
	if N == 1 {
		return M
	}

	width := 0
	for _, r := range R {
		width += 2 * r
	}

	binom := func(n, k int) int {

		// Too slow!!
		// z := new(big.Int)
		// z.Binomial(int64(n), int64(k))
		// z.Mod(z, billionSeven)
		// return int(z.Int64())

		return choose_mod(n, k, 1000000007)
	}

	x := 0

	// Pick 2 extremities, and compute for each pair
	for i := 0; i < N-1; i++ {
		for j := i + 1; j < N; j++ {
			minM := width - R[i] - R[j] + 1
			if minM > M {
				continue
			}
			groupOrder := facto[N-2]

			holes := M - minM
			combiHoles := binom(N+1+holes-1, N+1-1) // stars and bars

			symmetry := 2

			y := (groupOrder * combiHoles * symmetry) % 1000000007
			x = (x + y) % 1000000007
		}
	}

	return x
}

func (z *Case) solveLarge() interface{} {
	return nil
}

// Global precomputed data (if needed)

func precompute() {
	for i := 1; i < 3000; i++ {
		facto[i] = (facto[i-1] * i) % 1000000007
	}
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

// From http://stackoverflow.com/a/10862881/871134 :

func factorial_exponent(n, p int) int {
	ex := 0
	for {
		n /= p
		ex += n

		if n == 0 {
			break
		}
	}
	return ex
}

func choose_mod(n, k, p int) int {
	// We deal with the trivial cases first
	if k < 0 || n < k {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}
	// Now check whether choose(n,k) is divisible by p
	if factorial_exponent(n, p) > factorial_exponent(k, p)+factorial_exponent(n-k, p) {
		return 0
	}
	// If it's not divisible, do the generic work
	return choose_mod_one(n, k, p)
}

// Preconditions: 0 <= k <= n; p > 1 prime
func choose_mod_one(n, k, p int) int {
	// For small k, no recursion is necessary
	if k < p {
		return choose_mod_two(n, k, p)
	}
	q_n := n / p
	r_n := n % p
	q_k := k / p
	r_k := k % p
	choose := choose_mod_two(r_n, r_k, p)
	// If the exponent of p in choose(n,k) isn't determined to be 0
	// before the calculation gets serious, short-cut here:
	/* if (choose == 0) return 0; */
	choose *= choose_mod_one(q_n, q_k, p)
	return choose % p
}

// Preconditions: 0 <= k <= min(n,p-1); p > 1 prime
func choose_mod_two(n, k, p int) int {
	// reduce n modulo p
	n %= p
	// Trivial checks
	if n < k {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}
	// Now 0 < k < n, save a bit of work if k > n/2
	if k > n/2 {
		k = n - k
	}
	// calculate numerator and denominator modulo p
	num := n
	den := 1
	for n = n - 1; k > 1; n, k = n-1, k-1 {
		num = (num * n) % p
		den = (den * k) % p
	}
	// Invert denominator modulo p
	den = invert_mod(den, p)
	return (num * den) % p
}

func invert_mod(k, m int) int {
	if m == 0 {
		if k == 1 || k == -1 {
			return k
		} else {
			return 0
		}
	}
	if m < 0 {
		m = -m
	}
	k %= m
	if k < 0 {
		k += m
	}
	neg := true
	p1 := 1
	p2 := 0
	k1 := k
	m1 := m
	for k1 > 0 {
		q := m1 / k1
		r := m1 % k1
		temp := q*p1 + p2
		p2 = p1
		p1 = temp
		m1 = k1
		k1 = r
		neg = !neg
	}
	if neg {
		return m - p2
	} else {
		return p2
	}
}
