package main

import (
	"reflect"
	"testing"
)

func TestGetSpiralArray(t *testing.T) {
	assertEqual := func(t *testing.T, got [63]int, want [63]int) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, wanted %q", got, want)
		}
	}

	got := get_sprial_array([63]int{1})
	want := [63]int{}
	assertEqual(t, got, want)
}
