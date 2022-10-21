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
func (poly *Poly[T]) MakeTable(opts ...TableOption) []T {
	var conf tableConf
	conf.dataWidth = 8
	for _, o := range opts {
		o(&conf)
	}
	N := 1 << conf.dataWidth

	updateBitwise := BitwiseUpdateFn[T, T](poly)

	tab := make([]T, N)
	for i := range tab {
		tab[i] = updateBitwise(poly, T(conf.initial), T(i), conf.dataWidth)
		if conf.reverseBits {
			tab[i] = reverseBits(tab[i], poly.Width)
		}
	}
	return tab
}

type TableOption func(*tableConf)

type tableConf struct {
	initial     uint32
	dataWidth   int
	reverseBits bool
}

// WithDataWidthByte sets the data width to 8 bits,
// which will result in a table of 256 entries for byte-wise processing.
func WithDataWidth(w int) TableOption {
	return func(c *tableConf) {
		c.dataWidth = w
	}
}

// WithInitialValue ensures that when a table is created
// the specified initial value will be applied to the
// calculation of each table entry.
// In cases where a table later is used manually,
// like when a CRC is calulated over some bits only,
// this saves one XOR operation.
func WithInitialValue(initial uint32) TableOption {
	return func(c *tableConf) {
		c.initial = initial
	}
}

// WithReversedBits table option mirrors the bits of each table entry.
func WithReversedBits() TableOption {
	return func(c *tableConf) {
		c.reverseBits = true
	}
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
