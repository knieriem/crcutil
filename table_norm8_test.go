package crcutil_test

import (
	"fmt"
	"testing"

	"github.com/knieriem/crcutil/poly8"
)

// The humidity sensor [Si7021] that is, for example,
// part of the PG22-DK2503A development kit,
// includes a CRC with the x^8 + x^5 + x^4 + 1 polynomial (Dallas/Maxim)
// in normal form in some of its I²C messages.
// This test verifies the checksum over the electronic serial number
// response frames (id1, id2) of a real device.
// In the first frame each second byte is a crc result of the
// preceding data byte(s), in the second frame each third byte.
// TestDowCRCNorm verifies all these idFrameParts.
//
// [Si7021]: https://www.silabs.com/documents/public/data-sheets/Si7021-A20.pdf
func TestDowCRCNorm(t *testing.T) {
	tab := poly8.DOW.MakeTable()

	for i, part := range idFrameParts {
		t.Run(fmt.Sprintf("validate idFrame part %d", i), func(t *testing.T) {
			// crc(data + sum) must result in 0
			if sum := update8(0, tab, part); sum != 0 {
				t.Errorf("invalid crc8 for frame % x: %#02x", part, sum)
			}
		})
	}
}

var idFrameParts = [][]byte{
	// id1
	{0xa1, 0xcd},
	{0xa1, 0x68, 0x09},
	{0xa1, 0x68, 0xa5, 0x81},
	{0xa1, 0x68, 0xa5, 0xdc, 0x32},

	// id2
	{0x15, 0xff, 0xb5},
	{0x15, 0xff, 0xff, 0xff, 0xcb},
}

func update8(crc uint8, tab []uint8, p []byte) uint8 {
	for _, v := range p {
		crc = tab[byte(crc)^v]
	}
	return crc
}
