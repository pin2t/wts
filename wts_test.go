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
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

// @todo #1:30min Add more tests.
//  Untested methods are balance and rate,
//  add tests using `gock` similar to id testing.
//  Also, add coverage checks.

func TestCreateNewClient(t *testing.T) {
	tkn := "test-token"
	c, err := Create(tkn)
	a := assert.New(t)
	a.Nil(err)
	tr, ok := c.cli.Transport.(*authTransport)
	a.True(ok)
	a.Equal(tr.token, tkn)
}

func TestErrorCreateClientWithInvalidToken(t *testing.T) {
	c, err := Create("")
	a := assert.New(t)
	a.Equal(err.Error(), "Invalid token")
	a.Nil(c)
}

func TestHttpGetId(t *testing.T) {
	id := "1234"
	token := "test-token"
	gock.DisableNetworking()
	defer gock.Off()
	gock.New("http://localhost").
		Get("/id").
		MatchHeader(hdrToken, token).
		Reply(200).
		BodyString(id)
	a := assert.New(t)
	w, err := Create(token)
	a.Nil(err)
	w.host = "http://localhost"
	resp, err := w.ID()
	a.Nil(err)
	a.Equal(resp, id)
}
