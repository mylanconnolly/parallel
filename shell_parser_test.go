package main

import (
	"reflect"
	"strconv"
	"testing"
)

func TestShellParser(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"one two three", []string{"one", "two", "three"}},
		{`'"' "'"`, []string{`"`, "'"}},
		{`"\"something\"" here;`, []string{`"something"`, "here;"}},
		{"\tsomething\telse", []string{"something", "else"}},
		{"\nsomething\nelse", []string{"something", "else"}},
		{"\fsomething\felse", []string{"something", "else"}},
		{`'\' " "`, []string{`\`, " "}},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := shellParser(tt.input)

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Got: %#v Want: %#v", got, tt.want)
			}
		})
	}
}
