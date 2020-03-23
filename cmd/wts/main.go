package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/caarlos0/spin"
	"github.com/g4s8/wts"
)

var (
	cfg    = new(Config)
	filter string
	limit  int
	period string
	pind   bool
)

func main() {
	var debug bool
	var pull bool
	var config string
	var token string
	var templ string
	flag.StringVar(&token, "token", "", "API token")
	flag.StringVar(&filter, "filter", ".*", "Transactions regexp filter")
	flag.IntVar(&limit, "limit", -1, "Transactions limit")
	flag.BoolVar(&pull, "pull", false, "Pull wallet first")
	flag.BoolVar(&debug, "debug", false, "Debug output")
	flag.StringVar(&period, "period", "", "Days for statistic")
	flag.StringVar(&config, "config", "$HOME/.config/wts/config.yml",
		"Config file location")
	flag.StringVar(&templ, "fmt", "", "Output format (advanced)")
	flag.BoolVar(&pind, "progress", true, "Show progress indicator")
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
		fail("action required: id|balance|txns|pull|rate|pay")
	}
	switch action := args[0]; action {
	case "id":
		printID(w)
	case "balance":
		if templ == "" {
			templ = "{{printf \"%.2f\" .Zld}} ZLD (${{ printf \"%.2f\" .Usd}})\n"
		}
		printBalance(w, templ)
	case "txns":
		printTransactions(w)
	case "pull":
		pullWallet(w)
	case "rate":
		printRate(w)
	case "stats":
		printStats(w)
	case "pay":
		if len(args) != 5 {
			fail("Invalid payment command format. The correct one is " +
				"wts -token ... pay <destination> <amount> <keygap> <description>")
		}
		pay(w, args[1], args[2], args[3], args[4])
	default:
		fail(action + " - not implemented")
	}
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func failErr(err error) {
	fail(err.Error())
}

func printStats(w *wts.WTS) {
	if period == "" {
		fail("-period was not specified")
	}
	days, err := strconv.ParseInt(period, 10, 64)
	if err != nil {
		failErr(err)
	}
	dur := time.Duration(days) * 24 * time.Hour
	t := time.Now().Add(-dur)
	filter := wts.TxFilterSince(&t)
	pullIfNeeded(w)
	s := spinner(" Loading %s")
	txns, err := w.Transactions(filter, limit)
	s.Stop()
	if err != nil {
		failErr(err)
	}
	s = spinner(" Calculating %s")
	in := new(big.Float)
	out := new(big.Float)
	amt := new(big.Float)
	znts := new(big.Float).SetInt64(wts.ZldZents)
	zero := new(big.Float).SetInt64(0)
	for _, t := range txns {
		amt.SetInt64(t.Amount)
		amt.Quo(amt, znts)
		if amt.Cmp(zero) > 0 {
			in.Add(in, amt)
		} else if amt.Cmp(zero) < 0 {
			out.Add(out, amt)
		}
	}
	s.Stop()
	fmt.Printf("Stats for period: %s\n", period)
	fmt.Printf("\tIN:  %s\n", in.Text('f', 6))
	fmt.Printf("\tOUT: %s\n", out.Text('f', 6))
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

func printBalance(w *wts.WTS, t string) {
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
	var b struct {
		Zld float64
		Usd float64
	}
	b.Zld = float64(zents) / float64(wts.ZldZents)
	rate, err := w.UsdRate()
	s.Stop()
	if err != nil {
		failErr(err)
	}

	b.Usd = rate * b.Zld
	tpl := template.Must(template.New("balance").Parse(t))
	tpl.Execute(os.Stdout, b)
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
	var f wts.TxFilter
	if filter == "" {
		f = wts.TxFilterNone
	} else {
		ft, err := wts.TxFilterRegex(filter)
		if err != nil {
			failErr(err)
		}
		f = ft
	}
	txns, err := w.Transactions(f, limit)
	s.Stop()
	if err != nil {
		failErr(err)
	}
	for _, t := range txns {
		fmt.Println(t.String())
	}
}

func pay(w *wts.WTS, to string, amount string, keygap string, desc string) {
	pullIfNeeded(w)
	defer spinner(fmt.Sprintf(" Sending %s ZLD to %s", amount, to)).Stop()
	amt, _ := strconv.ParseFloat(amount, 64)
	err := w.Pay(to, uint64(amt*wts.ZldZents), keygap, desc)
	if err != nil {
		failErr(err)
	}
}

type progress interface {
	Stop()
}

type stubProg struct{}

func (s *stubProg) Stop() {}

type spinnerProg struct {
	spinner *spin.Spinner
}

func (s *spinnerProg) Stop() {
	s.spinner.Stop()
}

func spinner(lbl string) progress {
	if pind {
		s := spin.New(lbl)
		s.Set(spin.Spin1)
		s.Start()
		return &spinnerProg{s}
	}
	return new(stubProg)
}
