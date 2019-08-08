package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/spin"
	"github.com/g4s8/wts"
	"os"
)

var (
	cfg    = new(Config)
	filter string
	limit  int
)

func main() {
	var debug bool
	var pull bool
	var config string
	var token string
	flag.StringVar(&token, "token", "", "API token")
	flag.StringVar(&filter, "filter", ".*", "Transactions regexp filter")
	flag.IntVar(&limit, "limit", -1, "Transactions limit")
	flag.BoolVar(&pull, "pull", false, "Pull wallet first")
	flag.BoolVar(&debug, "debug", false, "Debug output")
	flag.StringVar(&config, "config", "$HOME/.config/wts/config.yml",
		"Config file location")
	flag.Parse()
	if err := cfg.ParseFile(config); err != nil {
		failErr(err)
	}
	if token != "" && token != cfg.Wts.Token {
		cfg.Wts.Token = token
	}
	if debug != cfg.Wts.Debug {
		cfg.Wts.Debug = debug
	}
	if pull != cfg.Wts.Pull {
		cfg.Wts.Pull = pull
	}
	if cfg.Wts.Token == "" {
		flag.Usage()
		fail("No token provided")
	}
	w, err := wts.Create(cfg.Wts.Token)
	if err != nil {
		failErr(err)
	}
	w.Debug = debug
	args := flag.Args()
	if len(args) == 0 {
		fail("action required: id|balance|txns|pull|rate")
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
	case "rate":
		printRate(w)
	default:
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
	s := spinner(" Loading %s")
	id, err := w.ID()
	s.Stop()
	if err != nil {
		failErr(err)
	}
	fmt.Printf("ID: %s\n", id)
}

func printRate(w *wts.WTS) {
	s := spinner(" Loading %s")
	r, err := w.UsdRate()
	s.Stop()
	if err != nil {
		failErr(err)
	}
	fmt.Printf("ZLD:USD = %f\n", r)
}

func printBalance(w *wts.WTS) {
	// @todo #1:30min Use math/big & big.Float
	//  to operate with arbitrary-precision
	//  numbers. It should be used to calculate
	//  ZLD amount from zents and USD from ZLD.
	pullIfNeeded(w)
	s := spinner(" Loading %s")
	zents, err := w.Balance()
	if err != nil {
		s.Stop()
		fail(err.Error())
	}
	zld := float64(zents) / float64(wts.ZldZents)
	rate, err := w.UsdRate()
	s.Stop()
	if err != nil {
		failErr(err)
	}
	usd := rate * zld
	fmt.Printf("Balance: %f ZLD (%f USD)\n", zld, usd)
}

func pullIfNeeded(w *wts.WTS) {
	if cfg.Wts.Pull {
		pullWallet(w)
	}
}

func pullWallet(w *wts.WTS) {
	defer spinner(" Pulling %s").Stop()
	if err := w.Pull(); err != nil {
		failErr(err)
	}
}

func printTransactions(w *wts.WTS) {
	pullIfNeeded(w)
	s := spinner(" Loading %s")
	txns, err := w.Transactions(filter, limit)
	s.Stop()
	if err != nil {
		failErr(err)
	}
	for _, t := range txns {
		fmt.Println(t.String())
	}
}

func spinner(lbl string) *spin.Spinner {
	s := spin.New(lbl)
	s.Set(spin.Spin1)
	return s.Start()
}
