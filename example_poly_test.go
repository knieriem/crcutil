package crcutil_test

import (
	"fmt"

	"github.com/knieriem/crcutil"
)

type poly16 = crcutil.Poly[uint16]

var ccitt16 = &poly16{Word: 0x1021, Width: 16}

// This example calculates representations of the CCITT-16
// polynomial as shown in
// https://en.wikipedia.org/wiki/Mathematics_of_cyclic_redundancy_checks#Polynomial_representations
func ExamplePoly() {
	reversed := ccitt16.ReversedForm()

	formatPoly("reversed", reversed)
	formatPoly("reciprocal", ccitt16.ReciprocalForm())
	formatPoly("reversed reciprocal", reversed.ReciprocalForm())

	// Output:
	// 0x8408	reversed
	// 0x0811	reciprocal
	// 0x8810	reversed reciprocal
}

func formatPoly(variant string, poly *poly16) {
	fmt.Printf("%#04x\t%s\n", poly.Word, variant)
}
