package crcutil_test

import (
	"fmt"

	"github.com/knieriem/crcutil"
	"github.com/knieriem/crcutil/crc16"
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

	fmt.Printf("%#04x\n", inst.Sum())

	// Alternatively, use a predefined model:
	sum := crc16.Modbus.Checksum([]byte{2, 7})
	fmt.Printf("%#04x", sum)
	// Output:
	// 0x1241
	// 0x1241
}
