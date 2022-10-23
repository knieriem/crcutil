package crcutil

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
