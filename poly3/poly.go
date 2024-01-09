package poly3

import (
	"github.com/knieriem/crcutil/poly8"
)

const N = 3

var (
	// CRC-3-GSM: x³ + x + 1
	GSM = New(3)
)

func New(poly uint8) *poly8.Poly {
	return &poly8.Poly{Word: poly, Width: N}
}
