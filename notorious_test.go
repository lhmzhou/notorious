package main_test

import (
	"bytes"
	main "local/notorious"
	"local/notorious/opts"
	"os"
	"strings"
	"testing"
)

func mustNew(pattern string) opts.Opts {
	o, err := opts.New(pattern)
	if err != nil {
		panic(err)
	}
	return o
}

func Test_Notorious(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		input                    []string
		wantErrContains          string
		outContains, outExcludes []string
		opts                     opts.Opts
	}{
		"base case": {
			input:       []string{"foo", "bar"},
			opts:        mustNew("foo"),
			outContains: []string{"foo"},
		},
		"context": {
			input:       []string{"foo", "bar", "baz", "boo"},
			opts:        mustNew("bar").WithContext(1, 1),
			outContains: []string{"foo", "bar", "baz"},
			outExcludes: []string{"boo", "0", "1"},
		},
		"line numbers": {
			input:       []string{"foo", "bar", "baz", "boo"},
			opts:        mustNew("bar").WithContext(1, 1).WithLineNumbers(true),
			outContains: []string{"0", "1", "2"},
			outExcludes: []string{"3"},
		},
		"no match": {
			input:       []string{"ahsldjkashjkldas"},
			opts:        mustNew("foo"),
			outExcludes: []string{"a", "b", "c", "d"},
		},
	}
	for name, tt := range tests {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := new(bytes.Buffer)
			in := strings.NewReader(strings.Join(tt.input, "\n"))
			opts := tt.opts.WithVerbose(testing.Verbose())
			err := main.Notorious(opts, in, out, os.Stderr)
			if tt.wantErrContains != "" && err == nil || err != nil && !strings.Contains(err.Error(), tt.wantErrContains) {
				t.Fatalf("expected an error containing %q, but got %v", tt.wantErrContains, err)
			}
			got := out.String()
			for _, want := range tt.outContains {
				if !strings.Contains(got, want) {
					t.Logf("expected the output to contain %q: output: %q", want, got)
					t.Fail()
				}
			}
			for _, forbidden := range tt.outExcludes {
				if strings.Contains(got, forbidden) {
					t.Logf("expected the output not to contain %q: output: %q", forbidden, got)
					t.Fail()
				}
			}
		})
	}
}
