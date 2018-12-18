package main

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantProgram string
		wantArgs    []string
		wantOK      bool
	}{
		{"test 1", []string{"echo"}, "echo", nil, true},
		{"test 2", []string{"echo", "-n"}, "echo", []string{"-n"}, true},
		{"test 3", []string{}, "", nil, false},
		{"test 4", nil, "", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, args, ok := parseArgs(tt.args)

			if program != tt.wantProgram {
				t.Fatalf("Got program: %s, want: %s", program, tt.wantProgram)
			}
			if !reflect.DeepEqual(args, tt.wantArgs) {
				t.Fatalf("Got args: %#v, want: %#v", args, tt.wantArgs)
			}
			if ok != tt.wantOK {
				t.Fatalf("Got ok: %t, want: %t", ok, tt.wantOK)
			}
		})
	}
}

func TestLineReader(t *testing.T) {
	tests := []struct {
		name   string
		reader io.Reader
		want   []string
	}{
		{"test 1", bytes.NewBuffer([]byte("")), []string{""}},
		{"test 2", bytes.NewBuffer([]byte("1\n")), []string{"1"}},
		{"test 3", bytes.NewBuffer([]byte("1\n2\n")), []string{"1", "2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := lineReader(tt.reader)
			got := []string{}

			for str := range reader {
				got = append(got, str)
			}
			if !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("Got: %#v, want: %#v", got, tt.want)
			}
		})
	}
}
