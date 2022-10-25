// Package impl contains implementations of the crcutil.Impl interface.
package impl

type Word interface {
	uint8 | uint16 | uint32
}

type Impl8[T Word] struct{}

func (impl Impl8[T]) Update(crc T, tab []T, p []byte) T {
	for _, v := range p {
		crc = tab[byte(crc)^v]
	}
	return crc
}

func (impl Impl8[T]) Append(in []byte, crc T) []byte {
	return append(in, byte(crc))
}

type Impl16[T Word] struct{}

func (impl Impl16[T]) Update(crc T, tab []T, p []byte) T {
	for _, v := range p {
		crc = tab[byte(crc>>8)^v] ^ (crc << 8)
	}
	return crc
}

func (impl Impl16[T]) Append(in []byte, crc T) []byte {
	return append(in, byte(crc>>8), byte(crc))
}

type Impl32[T Word] struct{}

func (impl Impl32[T]) Update(crc T, tab []T, p []byte) T {
	for _, v := range p {
		crc = tab[byte(crc>>24)^v] ^ (crc << 8)
	}
	return crc
}

func (impl Impl32[T]) Append(in []byte, crc T) []byte {
	return append(in, byte(crc>>24), byte(crc>>16), byte(crc>>8), byte(crc))
}

type Impl16LSBitFirst[T Word] struct{}

func (impl Impl16LSBitFirst[T]) Update(crc T, tab []T, p []byte) T {
	for _, v := range p {
		crc = tab[byte(crc)^v] ^ (crc >> 8)
	}
	return crc
}

func (impl Impl16LSBitFirst[T]) Append(in []byte, crc T) []byte {
	return append(in, byte(crc), byte(crc>>8))
}

type Impl32LSBitFirst[T Word] struct{ Impl16LSBitFirst[T] }

func (impl Impl32LSBitFirst[T]) Append(in []byte, crc T) []byte {
	return append(in, byte(crc), byte(crc>>8), byte(crc>>16), byte(crc>>24))
}
