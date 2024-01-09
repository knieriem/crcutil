package poly16

import (
	"github.com/knieriem/crcutil"
)

type Poly = crcutil.Poly[uint16]

const N = 16

var (
	// CRC-16-CCITT: x¹⁶ + x¹² + x⁵ + 1
	CCITT = New(0x1021)

	// CRC-16-IBM: x¹⁶ + x¹⁵ + x² + 1
	IBM = New(0x8005)
)

func New(poly uint16) *Poly {
	return &Poly{Word: poly, Width: N}
}
