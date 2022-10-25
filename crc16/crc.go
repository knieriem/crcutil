package crc16

import (
	"github.com/knieriem/crcutil"
	"github.com/knieriem/crcutil/poly16"
)

type Model = crcutil.Model[uint16]
type Inst = crcutil.Inst[uint16]

var (
	Modbus = &Model{
		Poly:          poly16.IBM.ReversedForm(),
		InitialInvert: true,
	}
)
