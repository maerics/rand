package main

import (
	"bytes"
	"io"
	"regexp"
	"testing"
)

func TestMainHelp(t *testing.T) {
	helpArgs := [][]string{{}, {"-h"}, {"--help"}}

	helpMessages := []*regexp.Regexp{
		regexp.MustCompile(`^Print random bytes from a secure source to stdout.\n`),
		regexp.MustCompile(`Usage:\n  rand LENGTH_BYTES \[flags\]`),
		regexp.MustCompile(`Flags:\n  -`),
	}

	for _, args := range helpArgs {
		stdout := &bytes.Buffer{}
		randCmd = newRandCmd()
		randCmd.SetArgs(args)
		randCmd.SetOut(stdout)
		main()

		output := stdout.String()

		for i, helpMessage := range helpMessages {
			if !helpMessage.MatchString(output) {
				t.Fatalf(
					"unexpected help message example %#v:\n  wanted: %q\n     got: %q",
					i+1, helpMessage.String(), output)
			}
		}
	}
}

func TestAll(t *testing.T) {
	// args, stdin, output
	for i, example := range []struct {
		args      []string
		input     io.Reader
		predicate func([]byte) bool
	}{
		// Encodings
		{[]string{"2"}, nil, func(bs []byte) bool { return regexp.MustCompile(`^[A-Fa-f0-9]{4}\n$`).Match(bs) }},
		{[]string{"9", "-a"}, nil, func(bs []byte) bool { return regexp.MustCompile(`^[A-Za-z0-9+/]{12}\n$`).Match(bs) }},
		{[]string{"10", "--base64"}, nil, func(bs []byte) bool { return regexp.MustCompile(`^[A-Za-z0-9+/]{14}={2}\n$`).Match(bs) }},
		{[]string{"11", "--base64"}, nil, func(bs []byte) bool { return regexp.MustCompile(`^[A-Za-z0-9+/]{15}={1}\n$`).Match(bs) }},

		// UUID
		{[]string{"--uuid"}, nil, func(bs []byte) bool {
			return regexp.MustCompile(`^[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}\n$`).Match(bs)
		}},
		{[]string{"-u"}, nil, func(bs []byte) bool {
			return regexp.MustCompile(`^[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}\n$`).Match(bs)
		}},
	} {
		randCmd = newRandCmd()
		randCmd.SetArgs(example.args)
		randCmd.SetIn(example.input)
		output := new(bytes.Buffer)
		randCmd.SetOut(output)
		main()
		actual := output.Bytes()
		if !example.predicate(actual) {
			t.Fatalf("unexpected output for example #%v (args=%#v), got: %#v (%s)", i+1, example.args, actual, actual)
		}
	}
}
