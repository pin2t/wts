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
)

var (
	errInvalidToken = errors.New("Invalid token")
)

const (
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
	token string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(hdrToken, t.token)
	return http.DefaultTransport.RoundTrip(req)
}

// Create new WTS client
func Create(token string) (*WTS, error) {
	if token == "" {
		return nil, errInvalidToken
	}
	return &WTS{
		&http.Client{Transport: &authTransport{token}},
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
	w.debug(fmt.Sprintf("Transactions(filter=%s, limit=%d): txns[%s]",
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
	return fmt.Sprintf("%d %d %s",
		t.ID, t.Amount, t.Details)
}

func (w *WTS) debug(msg string) {
	if w.Debug {
		fmt.Println("DEBUG: " + msg)
	}
}
