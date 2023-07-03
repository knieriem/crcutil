package crcutil

// FromImplicit1Notation returns a [Poly] with the word k converted from
// Koopman's implicit +1 notation to the explicit +1 notation used within this module.
// An advantage of the implicit +1 notation is that the bit width
// of the polynomial can be derived from the word value.
//
// Notation examples:
//
//	implicit +1:  1*x^3  + 0*x^2 + 1*x^1 + (1*x^0) = 0b 1 01(1) => 0x5
//	explicit +1: (1*x^3) + 0*x^2 + 1*x^1 +  1*x^0  = 0b(1)01 1  => 0x3
//
// See https://users.ece.cmu.edu/~koopman/crc/notes.html for details.
//
// Note that the implicit +1 notation is only a different way to write
// the polynomial word, but there is no algorithmic difference calculating
// the CRC.
// A word in implicit +1 notation looks identical to the word value
// of a [Poly] refering to the same polynomial,
// but converted to reverse reciprocal form.
func FromImplicit1Notation[T Word](k T) *Poly[T] {
	n := 0
	for v := k; v != 0; v >>= 1 {
		n++
	}

	// Make room for the explicit +1 bit.
	k <<= 1

	// Clear the topmost bit. This is a no-op if n equals the word width.
	k ^= 1 << n

	// Add explicit +1 bit.
	w := k | 1

	return &Poly[T]{Word: w, Width: n}
}

// FromImplicit1NotationReciprocal is like [FromImplicit1Notation] but returns a [Poly]
// with the Reciprocal flag set.
//
// Note: A word in implicit +1 notation looks identical to the word value
// of a [Poly] refering to the same polynomial,
// but converted to reverse form.
func FromImplicit1NotationReciprocal[T Word](k T) *Poly[T] {
	p := FromImplicit1Notation(k)
	p.Reciprocal = true
	return p
}
