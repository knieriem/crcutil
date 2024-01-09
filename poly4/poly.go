package poly4

import (
	"github.com/knieriem/crcutil/poly8"
)

const N = 4

var (
	// CRC-4-ITU: x⁴ + x + 1
	ITU = New(3)
)

func New(poly uint8) *poly8.Poly {
	return &poly8.Poly{Word: poly, Width: N}
}
