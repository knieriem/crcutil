package crcutil

import (
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
