package main

import "testing"

func TestFibOne(t *testing.T) {
	actual := fibOne(6)
	if actual != 8 {
		t.Error("expected value of 8 at position 6.")
	}
}

func TestFibTwo(t *testing.T) {
	actual := fibTwo(6, nil)
	if actual != 8 {
		t.Error("expected value of 8 at position 6.")
	}
}
