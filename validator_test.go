package main

import (
	"testing"
)

func TestCleanBody(t *testing.T) {
	got := cleanBody("This is a kerfuffle opinion I need to share with the world")
	want := "This is a **** opinion I need to share with the world"

	if got != want {
		t.Errorf("\n got: %q\nwant: %q", got, want)
	}

	got = cleanBody("This is kerfuffle, but not an opinion I need to share with the world")
	want = "This is kerfuffle, but not an opinion I need to share with the world"

	if got != want {
		t.Errorf("\n got: %q\nwant: %q", got, want)
	}

	got = cleanBody("Sharbert! This is a kerfuffle opinion I need to share with the world")
	want = "Sharbert! This is a **** opinion I need to share with the world"

	if got != want {
		t.Errorf("\n got: %q\nwant: %q", got, want)
	}
}
