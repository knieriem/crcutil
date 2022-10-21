package crcutil

var (
	// CRC-3-GSM: x^3 + x^1 + 1
	GSM3 = &Poly[uint8]{Word: 3, Width: 3}
)
