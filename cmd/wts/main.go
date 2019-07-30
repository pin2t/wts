package main

import (
	"flag"
	"fmt"
	"github.com/g4s8/wts"
	"os"
)

var (
	token  string
	filter string
	limit  int
	debug  bool
	pull   bool
)

func main() {
	flag.StringVar(&token, "token", "", "API token")
	flag.StringVar(&filter, "filter", ".*", "Transactions regexp filter")
	flag.IntVar(&limit, "limit", -1, "Transactions limit")
	flag.BoolVar(&pull, "pull", false, "Pull wallet first")
	flag.BoolVar(&debug, "debug", false, "Debug output")
	flag.Parse()
	if token == "" {
		flag.Usage()
		os.Exit(1)
	}
	w, err := wts.Create(token)
	if err != nil {
		failErr(err)
	}
	w.Debug = debug
	args := flag.Args()
	if len(args) == 0 {
		fail("action required: id|balance|txns|pull")
	}
	switch action := args[0]; action {
	case "id":
		printID(w)
	case "balance":
		printBalance(w)
	case "txns":
		printTransactions(w)
	case "pull":
		pullWallet(w)
	default:
		// @todo #3:30min Implement other actions, such as
		//  pull and others, see WTS readme file for
		//  more details about API methods.
		fail(action + " - not implemented")
	}
}

func fail(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func failErr(err error) {
	fail(err.Error())
}

func printID(w *wts.WTS) {
	pullIfNeeded(w)
	id, err := w.ID()
	if err != nil {
		failErr(err)
	}
	fmt.Printf("ID: %s\n", id)
}

func printBalance(w *wts.WTS) {
	// @todo #1:30min Use math/big & big.Float
	//  to operate with arbitrary-precision
	//  numbers. It should be used to calculate
	//  ZLD amount from zents and USD from ZLD.
	pullIfNeeded(w)
	zents, err := w.Balance()
	if err != nil {
		fail(err.Error())
	}
	zld := float64(zents) / float64(wts.ZldZents)
	rate, err := w.UsdRate()
	if err != nil {
		failErr(err)
	}
	usd := rate * zld
	fmt.Printf("Balance: %f ZLD (%f USD)\n", zld, usd)
}

func pullIfNeeded(w *wts.WTS) {
	if pull {
		pullWallet(w)
	}
}

func pullWallet(w *wts.WTS) {
	if err := w.Pull(); err != nil {
		failErr(err)
	}
}

func printTransactions(w *wts.WTS) {
	pullIfNeeded(w)
	txns, err := w.Transactions(filter, limit)
	if err != nil {
		failErr(err)
	}
	for _, t := range txns {
		fmt.Println(t.String())
	}
}
