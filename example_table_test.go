package crcutil_test

import (
	"fmt"

	"github.com/knieriem/crcutil"
)

// This example shows how to calculate a 3-bit CRC over
// 13 bits of data from two bytes of a frame of a serial
// communication protocol:
//
//	byte 0: crc (3 bits) | device address (5 bits)
//	byte 1: function code + length (8 bits)
//
// A specific algorithm will be used: 3-bit polynomial 0x3 (GSM),
// in reversed form. With an initial value of 5, the crc is calculated first
// over the five bits of the first byte, then the 8 bits of the second byte.
// After adding the 8 bit value, the crc result shall be mirrored (reflected).
// We will create
// -	one table for the 5 bit lookup, which already
// 	regards for the initial value, and
// -	one table for the 8 bit lookup, which produces a reflected result.

type crcTabs struct {
	t1 []uint8
	t2 []uint8
}

var crc3 crcTabs

func init() {
	poly := crcutil.GSM3.ReversedForm()

	initialValue := crcutil.WithInitialValue(5)
	fiveBits := crcutil.WithDataWidth(5)

	crc3.t1 = poly.MakeTable(initialValue, fiveBits)
	crc3.t2 = poly.MakeTable(crcutil.WithReversedBits())
}

// Then we can calculate the checksum by performing
// a table lookup first for the 5-bit data part, XOR-ing the result
// with the 8-bit data part, and doing a lookup of this value with
// the other table.
func (tabs *crcTabs) checksum(v1, v2 uint8) uint8 {
	crc := tabs.t1[v1]
	return tabs.t2[crc^v2]
}

// With the 5-bit device address being 10, and a value of the second
// byte of 0b0111_0101, the checksum should be 4.

func ExamplePoly_MakeTable() {
	sum := crc3.checksum(10, 0b0111_0101)
	fmt.Println(sum)
	// Output: 4
}
