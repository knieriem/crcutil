package crcutil

import (
	"fmt"
)

var ccitt16 = &Poly[uint16]{Word: 0x1021, Width: 16}

// This example calculates representations of the CCITT-16
// polynomial as shown in
// https://en.wikipedia.org/wiki/Mathematics_of_cyclic_redundancy_checks#Polynomial_representations
func ExamplePoly() {
	reversed := ccitt16.ReversedForm()
	fmt.Printf("%#04x %#04x %#04x",
		reversed.Word,
		ccitt16.ReciprocalForm().Word,
		reversed.ReciprocalForm().Word)
	// Output: 0x8408 0x0811 0x8810
}
