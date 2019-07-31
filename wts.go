// MIT License

// Copyright (c) 2019 g4s8

// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files
// (the "Software"), to deal in the Software without restriction,
// including without limitation the rights * to use, copy, modify,
// merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:

// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package wts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var (
	errInvalidToken = errors.New("Invalid token")
)

const (
	// ZldZents - zents in one ZLD
	ZldZents = 2 << 31
	hdrToken = "X-Zold-Wts"
)

// WTS client
type WTS struct {
	cli   *http.Client
	host  string
	Debug bool
}

// Txn - WTS transaction
type Txn struct {
	ID      uint64 `json:"id"`
	Date    string `json:"date"`
	Amount  int64  `json:"amount"`
	Prefix  string `json:"prefix"`
	Details string `json:"details"`
}

// authTransport adds auth headers to HTTP requests
type authTransport struct {
	base  http.RoundTripper
	token string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(hdrToken, t.token)
	return t.base.RoundTrip(req)
}

type wtsErrTransport struct {
	base http.RoundTripper
}

func (t *wtsErrTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	url := req.URL
	q := url.Query()
	q.Add("noredirect", "1")
	url.RawQuery = q.Encode()
	req.URL = url
	rsp, err := t.base.RoundTrip(req)
	if err != nil {
		return rsp, err
	}
	if rsp.StatusCode != 200 {
		hdr := rsp.Header.Get("X-Zold-Error")
		if hdr != "" {
			return nil, errors.New(hdr)
		}
		return nil, fmt.Errorf("WTS returned %d without error description",
			rsp.StatusCode)
	}
	return rsp, nil
}

// Create new WTS client
func Create(token string) (*WTS, error) {
	if token == "" {
		return nil, errInvalidToken
	}
	t := &authTransport{
		&wtsErrTransport{http.DefaultTransport},
		token,
	}
	return &WTS{
		&http.Client{Transport: t},
		"https://wts.zold.io",
		false,
	}, nil
}

// ID of wallet
func (w *WTS) ID() (string, error) {
	return w.getText("/id")
}

// Balance of wallet
func (w *WTS) Balance() (uint64, error) {
	rsp, err := w.getText("/balance")
	if err != nil {
		return 0, err
	}
	balance, err := strconv.ParseUint(rsp, 10, 64)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

// UsdRate at WTS
func (w *WTS) UsdRate() (float64, error) {
	rsp, err := w.getText("/usd_rate")
	if err != nil {
		return 0, err
	}
	rate, err := strconv.ParseFloat(rsp, 64)
	if err != nil {
		return 0, err
	}
	return rate, nil
}

// Transactions of wallet
func (w *WTS) Transactions(filter string, limit int) ([]Txn, error) {
	var txns []Txn
	rsp, err := w.cli.Get(w.host + "/txns.json?sort=desc")
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	if err := json.NewDecoder(rsp.Body).Decode(&txns); err != nil {
		return nil, err
	}
	w.debug(fmt.Sprintf("Transactions(filter=%s, limit=%d): txns[%d]",
		filter, limit, len(txns)))
	if limit < 0 {
		limit = len(txns)
	}
	result := make([]Txn, 0, limit)
	ptn, err := regexp.Compile(filter)
	if err != nil {
		return nil, err
	}
	for pos, txn := range txns {
		if pos >= limit {
			break
		}
		if ptn.MatchString(txn.Details) {
			result = append(result, txn)
		}
	}
	return result, nil
}

// Pull wallet from the network
func (w *WTS) Pull() error {
	rsp, err := w.cli.Get(w.host + "/pull")
	if err != nil {
		return err
	}
	job := rsp.Header.Get("X-Zold-Job")
	w.debug(fmt.Sprintf("pulling wallet, job=%s", job))
	var done bool
	for {
		rsp, err := w.cli.Get(w.host + "/job?id=" + job)
		if err != nil {
			return err
		}
		defer rsp.Body.Close()
		b, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return err
		}
		w.debug(fmt.Sprintf("pulling, status=%s", string(b)))
		done = string(b) == "OK"
		if done {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
}

func (w *WTS) getText(path string) (string, error) {
	rsp, err := w.cli.Get(w.host + path)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (t *Txn) String() string {
	return fmt.Sprintf("[%d] (%s) %f ZLD: %s",
		t.ID, t.Date, float64(t.Amount)/float64(ZldZents), t.Details)
}

func (w *WTS) debug(msg string) {
	if w.Debug {
		fmt.Println("DEBUG: " + msg)
	}
}
