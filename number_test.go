// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json5

import (
	"regexp"
	"testing"
)

func TestNumberIsValid(t *testing.T) {
	validTests := []string{
		"0",
		"-0",
		"1",
		"+1",
		"-1",
		"0.1",
		"-0.1",
		"+0.1",
		"1234",
		"-1234",
		"+1234",
		"12.34",
		"-12.34",
		"+12.34",
		"12E0",
		"12E1",
		"12e34",
		"12E-0",
		"12e+1",
		"12e-34",
		"-12E0",
		"-12E1",
		"-12e34",
		"-12E-0",
		"-12e+1",
		"-12e-34",
		"+12E0",
		"+12E1",
		"+12e34",
		"+12E-0",
		"+12e+1",
		"+12e-34",
		"1.2E0",
		"1.2E1",
		"1.2e34",
		"1.2E-0",
		"1.2e+1",
		"1.2e-34",
		"-1.2E0",
		"-1.2E1",
		"-1.2e34",
		"-1.2E-0",
		"-1.2e+1",
		"-1.2e-34",
		"+1.2E0",
		"+1.2E1",
		"+1.2e34",
		"+1.2E-0",
		"+1.2e+1",
		"+1.2e-34",
		"0E0",
		"0E1",
		"0e34",
		"0E-0",
		"0e+1",
		"0e-34",
		"-0E0",
		"-0E1",
		"-0e34",
		"-0E-0",
		"-0e+1",
		"-0e-34",
		"+0E0",
		"+0E1",
		"+0e34",
		"+0E-0",
		"+0e+1",
		"+0e-34",
		"1.",
		"+1.",
		"-1.",
		"1.e1",
		"+1.e1",
		"-1.e1",
		".5",
		"-.5",
		"+.5",
		"0xa",
		"0xA",
		"0XA",
		"-0XA",
		"+0XA",
		"0x0",
		"0x0",
		"0X0",
		"-0X0",
		"+0X0",
		"0x2",
		"0x2",
		"0X2",
		"-0X2",
		"+0X2",
		"0xDEADBeef",
		"-0xDEADBeef",
		"+0xDEADBeef",
		"0XDEADBeef",
		"-0XDEADBeef",
		"+0XDEADBeef",
		"0xDEAD3eef",
		"-0xDEAD3eef",
		"+0xDEAD3eef",
		"0XDEAD3eef",
		"-0XDEAD3eef",
		"+0XDEAD3eef",
		"NaN",
		"+Infinity",
		"-Infinity",
		"Infinity",
	}

	for _, test := range validTests {
		if !isValidNumber(test) {
			t.Errorf("%s should be valid", test)
		}

		var f float64
		if err := Unmarshal([]byte(test), &f); err != nil {
			t.Errorf("%s should be valid but Unmarshal failed: %v", test, err)
		}
	}

	invalidTests := []string{
		"",
		"invalid",
		"1.0.1",
		"1..1",
		"-1-2",
		"012a42",
		"01.2",
		"012",
		"12E12.12",
		"1e2e3",
		"1e+-2",
		"1e--23",
		"1e",
		"e1",
		"1e+",
		"1ea",
		"1a",
		"1.a",
		"01",
		"0xDsADBeef",
		".0xDEADBeef",
		"0XDsADBeef",
		".0XDEADBeef",
		"+NaN",
		"-NaN",
		".NaN",
		".Infinity",
		"0xs",
	}

	for _, test := range invalidTests {
		if isValidNumber(test) {
			t.Errorf("%s should be invalid", test)
		}

		var f float64
		if err := Unmarshal([]byte(test), &f); err == nil {
			t.Errorf("%s should be invalid but unmarshal wrote %v", test, f)
		}
	}
}

func BenchmarkNumberIsValid(b *testing.B) {
	s := "-61657.61667E+61673"
	for i := 0; i < b.N; i++ {
		isValidNumber(s)
	}
}

func BenchmarkNumberIsValidRegexp(b *testing.B) {
	var jsonNumberRegexp = regexp.MustCompile(`^-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?$`)
	s := "-61657.61667E+61673"
	for i := 0; i < b.N; i++ {
		jsonNumberRegexp.MatchString(s)
	}
}
