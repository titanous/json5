package json5

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/robertkrimen/otto"
)

type ErrorSpec struct {
	At           int
	LineNumber   int
	ColumnNumber int
	Message      string
}

func TestJSON5Decode(t *testing.T) {
	filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			t.Errorf("error reading file: %s", err)
			return nil
		}

		parseJSON5 := func() (interface{}, error) {
			var res interface{}
			return res, Unmarshal(data, &res)
		}
		parseJSON := func() (interface{}, error) {
			var res interface{}
			return res, json.Unmarshal(data, &res)
		}
		parseES5 := func() (interface{}, error) {
			vm := otto.New()
			_, err := vm.Run("x=" + string(data))
			if err != nil {
				return nil, err
			}
			v, err := vm.Get("x")
			if err != nil {
				return nil, err
			}
			return v.Export()
		}

		t.Logf("file: %s", path)
		switch filepath.Ext(path) {
		case ".json":
			jd, err := parseJSON()
			if err != nil {
				t.Errorf("unexpected error from json decoder: %s", err)
				return nil
			}
			j5d, err := parseJSON5()
			if err != nil {
				t.Errorf("unexpected error from json5 decoder: %s", err)
				return nil
			}
			if diff := pretty.Compare(jd, j5d); diff != "" {
				t.Errorf("data is not equal\n%s", diff)
				return nil
			}
		case ".json5":
			if _, err := parseJSON(); err == nil {
				t.Errorf("expected JSON parsing to fail")
				return nil
			}
			es5d, err := parseES5()
			if err != nil {
				t.Errorf("unexpected error from ES5 decoder: %s", err)
				return nil
			}
			j5d, err := parseJSON5()
			if err != nil {
				t.Errorf("unexpected error from json5 decoder: %s", err)
				return nil
			}
			if diff := pretty.Compare(j5d, es5d); diff != "" {
				t.Errorf("data is not equal\n%s", diff)
				return nil
			}
		case ".js":
			if _, err := parseJSON(); err == nil {
				t.Errorf("expected JSON parsing to fail")
				return nil
			}
			if _, err := parseES5(); err != nil {
				t.Errorf("unexected error from ES5 decoder: %s", err)
				return nil
			}
			if _, err := parseJSON5(); err == nil {
				t.Errorf("expected JSON5 parsing to fail")
				return nil
			}
		case ".txt":
			var expectedErr *ErrorSpec
			specName := strings.TrimRight(path, filepath.Ext(path)) + ".errorSpec"
			specFile, err := os.Open(specName)
			if err != nil && !os.IsNotExist(err) {
				t.Errorf("error trying to open errorSpec file %s: %s", specName, err)
				return nil
			}
			if specFile != nil {
				defer specFile.Close()
				expectedErr = &ErrorSpec{}
				if err := NewDecoder(specFile).Decode(expectedErr); err != nil {
					t.Errorf("error decoding %s: %s", specName, err)
					return nil
				}
			}
			_, err = parseJSON5()
			if err == nil {
				t.Errorf("expected JSON5 parsing to fail")
				return nil
			}
		}

		return nil
	})
}

// The tests below this comment were found with go-fuzz

func TestQuotedQuote(t *testing.T) {
	var v struct {
		E string
	}
	if err := Unmarshal([]byte(`{e:"'"}`), &v); err != nil {
		t.Error(err)
	}
	if v.E != "'" {
		t.Errorf(`expected "'", got %q`, v.E)
	}
}

func TestInvalidNewline(t *testing.T) {
	expected := "invalid character '\\n' in string literal"
	var v interface{}
	if err := Unmarshal([]byte("{a:'\\\r0\n'}"), &v); err == nil || err.Error() != expected {
		t.Errorf("expected error %q, got %s", expected, err)
	}
}
