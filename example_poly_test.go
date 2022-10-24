package crcutil_test

import (
	"fmt"

	"github.com/knieriem/crcutil"
	"github.com/knieriem/crcutil/poly16"
)

var ccitt16 = poly16.CCITT

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

func formatPoly[T crcutil.Word](variant string, poly *crcutil.Poly[T]) {
	fmt.Printf("%#04x\t%s\n", poly.Word, variant)
}
