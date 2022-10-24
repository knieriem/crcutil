package crcutil

import (
	"fmt"
	"testing"
)

var ibmcrc16 = (&Poly[uint16]{Word: 0x8005, Width: 16}).ReversedForm()

// TestMakeTable is creating a lookup table multiple times,
// each time calculating the same example checksum.
// Between sub tests, the table cache state is verified.
func TestMakeTable(t *testing.T) {
	cacheKey := tableCacheKey(ibmcrc16, &tableConf{
		dataWidth: 8,
	})
	t.Run("verify cache key", func(t *testing.T) {
		if cacheKey != "a001.r-0.8" {
			t.Fatalf("cache key mismatch: %q", cacheKey)
		}
	})
	t.Run("verify nil cache state", func(t *testing.T) {
		if tableCache[cacheKey] != nil {
			t.Fatal("tableCache entry not nil")
		}
	})
	t.Run("with empty cache entry", checkMakeTable)

	var tab []uint16

	// now there must be an entry for ibmcrc16 in the cache
	t.Run("verify initialized cache state", func(t *testing.T) {
		u, ok := tableCache[cacheKey].([]uint16)
		if !ok {
			t.Fatal("tableCache entry missing")
		}
		tab = u
	})

	t.Run("with existing cache entry", checkMakeTable)

	if update16(0xFFFF, tab, modExFrame) != 0x1241 {
		t.Fatal("unexpected failure when using cached table directly")
	}
}

var modExFrame = []byte{2, 7}

// CheckMakeTable creates a table for the CRC-16-IBM polynomial,
// then calculates a checksum over an example frame {2, 7}
// from the Modbus over serial line specification, and compares
// the result with the expected value.
func checkMakeTable(t *testing.T) {
	tab := ibmcrc16.MakeTable()
	crc := update16(0xFFFF, tab, modExFrame)
	if crc != 0x1241 {
		t.Errorf("crc mismatch: want 0x1241, got %#04x", crc)
	}
}

func update16(crc uint16, tab []uint16, p []byte) uint16 {
	for _, v := range p {
		crc = tab[byte(crc)^v] ^ (crc >> 8)
	}
	return crc
}

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
	poly := &Poly[uint8]{Word: 0x31, Width: 8}
	tab := poly.MakeTable()

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
