package crcutil

import (
	"unsafe"
)

// A Word holds the word representation of a polynomial.
type Word interface {
	~uint8 | ~uint16 | ~uint32
}

// MakeTable creates a lookup table for the specified polynomial.
// The size of the table returned equals the length of
// the value range determined by the number of data bits.
func (poly *Poly[T]) MakeTable(initial T, dataWidth int) []T {
	N := 1 << dataWidth

	updateBitwise := BitwiseUpdateFn[T, T](poly)

	tab := make([]T, N)
	for i := range tab {
		tab[i] = updateBitwise(poly, initial, T(i), dataWidth)
	}
	return tab
}

// ReverseTableBits mirrors the lower bits of each table entry within
// the bit range specified by width.
// The calculations are done in-place, so the original table
// is altered; the modified table is also returned for convenience.
func ReverseTableBits[T Word](tab []T, width int) []T {
	for i := range tab {
		tab[i] = reverseBits(tab[i], width)
	}
	return tab
}

// BitwiseUpdateFn returns a function that returns the result of
// adding the bits defined by data and dataWidth into the crc.
// The actual function returned depends on the polynomial's
// representation (normal or reversed form).
func BitwiseUpdateFn[T, D Word](poly *Poly[T]) func(poly *Poly[T], crc T, data D, dataWidth int) T {
	if !poly.Reversed {
		return updateBitwiseNormal[T, D]
	}
	return updateBitwiseReversed[T, D]
}

func UpdateBitwise[T, D Word](poly *Poly[T], crc T, data D, dataWidth int) T {
	return BitwiseUpdateFn[T, D](poly)(poly, crc, data, dataWidth)
}

func updateBitwiseNormal[T, D Word](poly *Poly[T], crc T, data D, dataWidth int) T {
	mask := (T(1) << poly.Width) - 1
	msb := T(1) << (poly.Width - 1)
	msbData := D(1) << (dataWidth - 1)

	for i := msbData; i > 0; i >>= 1 {
		bit := crc & msb
		if (data & i) != 0 {
			bit ^= msb
		}
		if bit != 0 {
			crc = (crc << 1) ^ poly.Word
		} else {
			crc <<= 1
		}
	}
	return crc & mask
}

func updateBitwiseReversed[T, D Word](poly *Poly[T], crcIn T, data D, dataWidth int) T {
	if unsafe.Sizeof(crcIn) > unsafe.Sizeof(data) {
		panic("D must be at least as wide as T")
	}
	mask := (D(1) << poly.Width) - 1
	crc := D(crcIn) ^ data
	for i := 0; i < dataWidth; i++ {
		if crc&1 == 1 {
			crc = (crc >> 1) ^ D(poly.Word)
		} else {
			crc >>= 1
		}
	}
	return T(crc & mask)
}
