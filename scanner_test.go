// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json5

import (
	"math"
	"math/rand"
	"testing"
)

// Tests of simple examples.

type example struct {
	compact string
	indent  string
}

var examples = []example{
	{`1`, `1`},
	{`{}`, `{}`},
	{`[]`, `[]`},
	{`{"":2}`, "{\n\t\"\": 2\n}"},
	{`[3]`, "[\n\t3\n]"},
	{`[1,2,3]`, "[\n\t1,\n\t2,\n\t3\n]"},
	{`{"x":1}`, "{\n\t\"x\": 1\n}"},
	{ex1, ex1i},
}

var ex1 = `[true,false,null,"x",1,1.5,0,-5e+2]`

var ex1i = `[
	true,
	false,
	null,
	"x",
	1,
	1.5,
	0,
	-5e+2
]`

func diff(t *testing.T, a, b []byte) {
	for i := 0; ; i++ {
		if i >= len(a) || i >= len(b) || a[i] != b[i] {
			j := i - 10
			if j < 0 {
				j = 0
			}
			t.Errorf("diverge at %d: «%s» vs «%s»", i, trim(a[j:]), trim(b[j:]))
			return
		}
	}
}

func trim(b []byte) []byte {
	if len(b) > 20 {
		return b[0:20]
	}
	return b
}

func genValue(n int) interface{} {
	if n > 1 {
		switch rand.Intn(2) {
		case 0:
			return genArray(n)
		case 1:
			return genMap(n)
		}
	}
	switch rand.Intn(3) {
	case 0:
		return rand.Intn(2) == 0
	case 1:
		return rand.NormFloat64()
	case 2:
		return genString(30)
	}
	panic("unreachable")
}

func genString(stddev float64) string {
	n := int(math.Abs(rand.NormFloat64()*stddev + stddev/2))
	c := make([]rune, n)
	for i := range c {
		f := math.Abs(rand.NormFloat64()*64 + 32)
		if f > 0x10ffff {
			f = 0x10ffff
		}
		c[i] = rune(f)
	}
	return string(c)
}

func genArray(n int) []interface{} {
	f := int(math.Abs(rand.NormFloat64()) * math.Min(10, float64(n/2)))
	if f > n {
		f = n
	}
	if f < 1 {
		f = 1
	}
	x := make([]interface{}, f)
	for i := range x {
		x[i] = genValue(((i+1)*n)/f - (i*n)/f)
	}
	return x
}

func genMap(n int) map[string]interface{} {
	f := int(math.Abs(rand.NormFloat64()) * math.Min(10, float64(n/2)))
	if f > n {
		f = n
	}
	if n > 0 && f == 0 {
		f = 1
	}
	x := make(map[string]interface{})
	for i := 0; i < f; i++ {
		x[genString(10)] = genValue(((i+1)*n)/f - (i*n)/f)
	}
	return x
}
