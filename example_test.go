package crcutil_test

import (
	"fmt"
)

// This example calculates the normal representation of the x^16 + x^15 + x^2 + 1
// polynomial back from its reversed and reversed reciprocal representations.
func ExamplePoly_NormalForm() {
	r := &poly16{Word: 0x8408, Width: 16, Reversed: true}
	rr := &poly16{Word: 0x8810, Width: 16, Reversed: true, Reciprocal: true}

	fmt.Printf("%#04x %#04x", r.NormalForm().Word, rr.NormalForm().Word)
	// Output: 0x1021 0x1021
}
