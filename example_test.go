package crcutil_test

import (
	"fmt"

	"github.com/knieriem/crcutil/poly16"
)

// This example calculates the normal representation of the x^16 + x^15 + x^2 + 1
// polynomial back from its reversed and reversed reciprocal representations.
func ExamplePoly_NormalForm() {
	r := &poly16.Poly{Word: 0x8408, Width: poly16.N, Reversed: true}
	rr := &poly16.Poly{Word: 0x8810, Width: poly16.N, Reversed: true, Reciprocal: true}

	fmt.Printf("%#04x %#04x", r.NormalForm().Word, rr.NormalForm().Word)
	// Output: 0x1021 0x1021
}
