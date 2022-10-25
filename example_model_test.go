package crcutil_test

import (
	"fmt"

	"github.com/knieriem/crcutil"
	"github.com/knieriem/crcutil/poly16"
)

var (
	ibmcrc = &crcutil.Model[uint16]{
		Poly:          poly16.IBM.ReversedForm(),
		InitialInvert: true,
	}
)

// This example calculates a Modbus crc over an example frame;
// see “MODBUS over serial line specification and implementation guide V1.02”, p. 41
func ExampleModel() {
	inst := ibmcrc.New()
	inst.Update([]byte{2, 7})

	fmt.Printf("%#04x", inst.Sum())
	// Output: 0x1241
}
