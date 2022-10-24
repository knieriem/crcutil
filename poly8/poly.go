package poly8

import (
	"github.com/knieriem/crcutil"
)

type Poly = crcutil.Poly[uint8]

const N = 8

var (
	// CRC-8-Dallas/Maxim: x^8 + x^5 + x^4 + 1
	DOW = New(0x31)

	// CRC-8-SAE-J1850: x^8 + x^4 + x^3 + x^2 + 1
	SAEJ1850 = New(0x1D)
)

func New(poly uint8) *Poly {
	return &Poly{Word: poly, Width: N}
}
