package main

import (
	"flag"
	"fmt"
	"github.com/g4s8/wts"
	"os"
)

func main() {
	var token string
	flag.StringVar(&token, "token", "", "API token")
	flag.Parse()
	if token == "" {
		flag.Usage()
	}
	w, err := wts.Create(token)
	if err != nil {
		panic(err)
	}
	args := flag.Args()
	if len(args) == 0 {
		fail("action required: id|balance|txns|pull")
	}
	switch action := args[0]; action {
	case "id":
		printID(w)
	case "balance":
		printBalance(w)
	default:
		// @todo #1:30min Implement other actions, such as
		//  pull and txns, see WTS readme file for
		//  more details about API methods.
		fail(action + " - not implemented")
	}
}

func fail(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func printID(w *wts.WTS) {
	id, err := w.ID()
	if err != nil {
		fail(err.Error())
	}
	fmt.Printf("ID: %s\n", id)
}

func printBalance(w *wts.WTS) {
	// @todo #1:30min Use math/big & big.Float
	//  to operate with arbitrary-precision
	//  numbers. It should be used to calculate
	//  ZLD amount from zents and USD from ZLD.
	zents, err := w.Balance()
	if err != nil {
		fail(err.Error())
	}
	zld := float64(zents) / float64(2<<31)
	rate, err := w.UsdRate()
	if err != nil {
		fail(err.Error())
	}
	usd := rate * zld
	fmt.Printf("Balance: %f ZLD (%f USD)\n", zld, usd)
}
