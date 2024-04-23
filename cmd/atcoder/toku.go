package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	defer wr.Flush()

	// 標準入力を受け取る

	res := solve()
	out(res)
}

func solve() int { // 受け取った入力を全て受け取り答えを出力
	return 0
}

const INF int = math.MaxInt64

var (
	sc *bufio.Scanner = initSc()
	wr *bufio.Writer  = initWr()
)

// the followings are the template IO functions
func initSc() *bufio.Scanner {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer([]byte{}, INF) // どでか入力対策
	return sc
}

func initWr() *bufio.Writer {
	return bufio.NewWriter(os.Stdout)
}

func ri() int {
	sc.Scan()
	return atoi(sc.Text())
}

func ri2() (int, int) {
	sc.Scan()
	ss := strings.Split(sc.Text(), " ")
	return atoi(ss[0]), atoi(ss[1])
}

func ri3() (int, int, int) {
	sc.Scan()
	ss := strings.Split(sc.Text(), " ")
	return atoi(ss[0]), atoi(ss[1]), atoi(ss[2])
}

func ris() []int {
	sc.Scan()
	ss := strings.Split(sc.Text(), " ")
	is := make([]int, len(ss))
	for i, s := range ss {
		is[i] = atoi(s)
	}
	return is
}

func risn(n int) [][]int { // nは行数
	sl := make([][]int, n)
	for i := 0; i < n; i++ {
		sl[i] = ris()
	}
	return sl
}

func rs() string {
	sc.Scan()
	return sc.Text()
}

func rss() []string {
	sc.Scan()
	return strings.Split(sc.Text(), " ")
}

func rssn(n int) [][]string {
	sl := make([][]string, n)
	for i := 0; i < n; i++ {
		sl[i] = rss()
	}
	return sl
}

func rcs() []string { // 文字列を1文字ずつ受け取る
	sc.Scan()
	return strings.Split(sc.Text(), "")
}

func rcsn(n int) [][]string {
	sl := make([][]string, n)
	for i := 0; i < n; i++ {
		sl[i] = rcs()
	}
	return sl
}

func out(v ...interface{}) {
	_, err := fmt.Fprintln(wr, v...)
	if err != nil {
		panic(err)
	}
}

func outIS(is []int) {
	ss := make([]string, len(is))
	for i, v := range is {
		ss[i] = itoa(v)
	}
	out(strings.Join(ss, " "))
}

func outSS(ss []string) {
	out(strings.Join(ss, " "))
}

// the followings are the slice functions
func del[T interface{}](a []T, i int) []T {
	return a[:i+copy(a[i:], a[i+1:])]
}

func delR[T interface{}](a []T, i int) []T { // 順番を考慮せずに削除(遥かに高速)
	a[i] = a[len(a)-1]
	return a[:len(a)-1]
}

func prepend[T interface{}](h T, sl []T) []T {
	return append([]T{h}, sl...)
}

// the followings are the generic template functions
func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

func min(a int, b int) int {
	if a >= b {
		return b
	} else {
		return a
	}
}

func max(a int, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func chMin(a *int, b int) bool {
	tmp := *a
	*a = min(*a, b)
	if tmp == *a {
		return false
	} else {
		return true
	}
}

func chMax(a *int, b int) bool {
	tmp := *a
	*a = max(*a, b)
	if tmp == *a {
		return false
	} else {
		return true
	}
}

func nmin(a []int) (int, int) {
	if len(a) == 0 {
		return -1, -1
	}
	idx, mini := 0, a[0]
	for i, v := range a {
		if chMin(&mini, v) {
			idx = i
		}
	}
	return idx, mini
}

func nmax(a []int) (int, int) {
	if len(a) == 0 {
		return -1, -1
	}
	idx, maxi := 0, a[0]
	for i, v := range a {
		if chMax(&maxi, v) {
			idx = i
		}
	}
	return idx, maxi
}

func nsum(a []int) int {
	sum := 0
	for _, v := range a {
		sum += v
	}
	return sum
}

func ntime(a []int) int {
	prod := 1
	for _, v := range a {
		prod *= v
	}
	return prod
}

func pow(a int, b int) int {
	exp := 1
	for i := 0; i < b; i++ {
		exp *= a
	}
	return exp
}

func fact(n int) int {
	res := 1
	for i := 1; i <= n; i++ {
		res *= i
	}
	return res
}

func abs(a int) int {
	if a >= 0 {
		return a
	} else {
		return -a
	}
}

func nroot(n int, a int) int { // n乗根を二分探索
	l, r := 0, a
	var m int
	for l+1 < r {
		m = (l + r) / 2
		if pow(m, n) > a {
			r = m
		} else {
			l = m
		}
	}
	return l
}
