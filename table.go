package crcutil

import (
	"fmt"
	"sync"

	"github.com/knieriem/crcutil/internal/impl"
)

// MakeTable creates a lookup table for the specified polynomial.
// The size of the table returned equals the length of
// the value range determined by the number of data bits.
func (poly *Poly[T]) MakeTable(opts ...TableOption) []T {
	var conf tableConf
	conf.dataWidth = 8
	for _, o := range opts {
		o(&conf)
	}

	k := tableCacheKey(poly, &conf)
	tableCacheMu.RLock()
	t := tableCache[k]
	tableCacheMu.RUnlock()
	if t != nil {
		return t.([]T)
	}

	N := 1 << conf.dataWidth

	updateBitwise := BitwiseUpdateFn[T, T](poly)

	tab := make([]T, N)
	for i := range tab {
		crc := updateBitwise(poly, T(conf.initial), T(i), conf.dataWidth)
		if conf.reverseBits {
			crc = reverseBits(crc, poly.Width)
		}
		tab[i] = crc
	}

	tableCacheMu.Lock()
	defer tableCacheMu.Unlock()
	tableCache[k] = tab
	return tab
}

type TableOption func(*tableConf)

type tableConf struct {
	initial     uint32
	dataWidth   int
	reverseBits bool
}

// WithDataWidth sets the data width to n bits,
// which will result in a table of 2^n entries for n-bit-wise processing.
// The default width is 8 bits for byte-wise processing.
func WithDataWidth(n int) TableOption {
	return func(c *tableConf) {
		c.dataWidth = n
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

var tableCacheMu sync.RWMutex
var tableCache = map[string]any{}

func tableCacheKey[T Word](p *Poly[T], c *tableConf) string {
	rep := ""
	if p.Reversed {
		rep += ".r"
	}
	if p.Reciprocal {
		rep += ".R"
	}
	tabMod := ""
	if c.reverseBits {
		tabMod = ".r"
	}
	return fmt.Sprintf("%x%s-%x.%d%s",
		p.Word, rep,
		c.initial, c.dataWidth, tabMod)
}

type Impl[T Word] interface {
	Update(crc T, tab []T, p []byte) T
	Append(b []byte, crc T) []byte
}

func (p *Poly[T]) Impl() Impl[T] {
	var x T

	switch any(x).(type) {
	case uint8:
		return impl.Impl8[T]{}
	case uint16:
		if p.LSBitFirst() {
			return impl.Impl16LSBitFirst[T]{}
		}
		return impl.Impl16[T]{}
	case uint32:
		if p.LSBitFirst() {
			return impl.Impl32LSBitFirst[T]{}
		}
		return impl.Impl32[T]{}
	}
	panic("type not supported")
}
