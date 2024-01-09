package crcutil

import (
	"math/bits"
)

// A Word holds the word representation of a polynomial.
type Word interface {
	uint8 | uint16 | uint32
}

// Poly defines a polynomial in a specific representation.
type Poly[T Word] struct {
	Word       T
	Width      int
	Reversed   bool
	Reciprocal bool
}

// NormalForm returns the normal form of the polynomial.
func (p *Poly[T]) NormalForm() *Poly[T] {
	if p.Reversed {
		p = p.reverse()
	}
	if p.Reciprocal {
		p = p.makeReciprocal()
	}
	return p
}

// ReversedForm returns the reversed, lsbit-first form of the polynomial.
// It returns the unchanged polynomial if it is already in its reversed form.
func (p *Poly[T]) ReversedForm() *Poly[T] {
	if p.Reversed {
		return p
	}
	p = p.reverse()
	return p
}

func (p *Poly[T]) reverse() *Poly[T] {
	r := new(Poly[T])
	r.Word = reverseBits(p.Word, p.Width)
	r.Width = p.Width
	r.Reversed = !p.Reversed
	r.Reciprocal = p.Reciprocal
	return r
}

// ReverseBits mirrors the lower n bits of the data value
// within the boundaries of those lower n bits.
func reverseBits[T Word](data T, n int) T {
	// support polynomials up to uint32
	// by allowing shift into bit32
	mask := uint32((uint64(1) << n) - 1)

	u := bits.Reverse32(uint32(data))
	u >>= 32 - n
	u &= mask
	return T(u)
}

// ReciprocalForm returns the reciprocal form of the polynomial
// that can be obtained by mirroring its coefficients.
// It returns the unchanged polynomial if it is already in its
// reciprocal form.
func (p *Poly[T]) ReciprocalForm() *Poly[T] {
	if p.Reciprocal {
		return p
	}
	if p.Reversed {
		p = p.reverse()
		p = p.makeReciprocal()
		return p.reverse()
	}
	return p.makeReciprocal()
}

func (p *Poly[T]) makeReciprocal() *Poly[T] {
	u := bits.Reverse32(uint32(p.Word))
	if shift := 32 - p.Width - 1; shift > 0 {
		u >>= shift
	} else if shift < 0 {
		u <<= -shift
	}

	u |= 1

	r := new(Poly[T])
	r.Word = T(u) & p.mask()
	r.Width = p.Width
	r.Reciprocal = !p.Reciprocal
	return r
}

func (p *Poly[T]) mask() T {
	return T((uint32(1) << p.Width) - 1)
}

// LSBitFirst reports whether the polynomial's representation is lsbit-first.
func (p *Poly[T]) LSBitFirst() bool {
	return p.Reversed
}
