package crc8

import (
	"github.com/knieriem/crcutil"
	"github.com/knieriem/crcutil/poly8"
)

type Model = crcutil.Model[uint8]
type Inst = crcutil.Inst[uint8]

var (
	DOW = &Model{
		Poly: poly8.DOW.ReversedForm(),
	}

	SAEJ1850 = &Model{
		Poly:          poly8.SAEJ1850,
		InitialInvert: true,
		FinalInvert:   true,
	}
)
