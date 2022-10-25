package crcutil_test

import (
	_ "embed"
	"testing"

	"github.com/knieriem/crcutil"
	"github.com/knieriem/crcutil/poly32"
	"hash/crc32"
)

var (
	ieee = &crcutil.Model[uint32]{
		Poly:          poly32.IEEE.ReversedForm(),
		InitialInvert: true,
		FinalInvert:   true,
	}
)

//go:embed model.go
var ieeeData []byte

func TestComparePolyIEEEAgainstStd(t *testing.T) {
	inst := ieee.New()
	inst.Update(ieeeData)
	sum := inst.Sum()

	h := crc32.NewIEEE()
	h.Write(ieeeData)
	stdSum := h.Sum32()

	if sum != stdSum {
		t.Fatalf("ieee checksum does not match checksum from stdlib: %08x vs. %08x", sum, stdSum)
	}
}
