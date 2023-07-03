package crcutil_test

import (
	"fmt"
	"testing"

	"github.com/knieriem/crcutil"
)

func ExampleFromImplicit1Notation() {
	p := crcutil.FromImplicit1Notation(uint32(0x19d17))
	fmt.Printf("%#05x %d", p.Word, p.Width)
	// Output: 0x13a2f 17
}

func ExampleFromImplicit1NotationReciprocal() {
	p := crcutil.FromImplicit1NotationReciprocal(uint32(0x1e8b9)).NormalForm()
	fmt.Printf("%#05x %d", p.Word, p.Width)
	// Output: 0x13a2f 17
}

// This test ensures that the implicit +1 notation is converted
// to a normal form with the same word as when converting a Poly
// in reverse reciprocal form with a visually identical word
// to normal form. Same for the reciprocal variant.
//
// The polynomial used for the test is CCITT-16, see:
// https://users.ece.cmu.edu/~koopman/crc/c16/0x8810.txt
func Test_Implicit1NotationVsAppearance(t *testing.T) {
	k := uint16(0x8810)
	p := crcutil.FromImplicit1Notation(k)
	pa := &crcutil.Poly[uint16]{
		Word:       k,
		Width:      p.Width,
		Reversed:   true,
		Reciprocal: true,
	}
	runTest := func(expectedWord uint16) {
		pn := pa.NormalForm()
		if p.Word != pn.Word {
			t.Errorf("Words in normal form derived from %#04x do not match: %#04x != %#04x", k, p.Word, pn.Word)
		}
		if p.Word != expectedWord {
			t.Errorf("Expected word in normal form: %#04x, got: %#04x", expectedWord, p.Word)
		}
	}
	runTest(0x1021)

	// test reciprocal variant
	k = 0x8408
	p = crcutil.FromImplicit1NotationReciprocal(k)
	pa.Reciprocal = false
	runTest(0x0811)
}
