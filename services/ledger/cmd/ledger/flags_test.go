package main

import "testing"

func TestParseFlagsVersion(t *testing.T) {
	got, err := parseFlags([]string{"-version"})
	if err != nil {
		t.Fatal(err)
	}
	if !got {
		t.Fatal("expected version flag to be reported")
	}
}

func TestParseFlagsUnknown(t *testing.T) {
	if _, err := parseFlags([]string{"-nope"}); err == nil {
		t.Fatal("expected error for unknown flag")
	}
}
