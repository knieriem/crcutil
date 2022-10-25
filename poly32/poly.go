package poly32

import (
	"github.com/knieriem/crcutil"
)

type Poly = crcutil.Poly[uint32]

const N = 32

var (
	IEEE = New(0x04C11DB7)
)

func New(poly uint32) *Poly {
	return &Poly{Word: poly, Width: N}
}
