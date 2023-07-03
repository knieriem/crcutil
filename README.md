[![Go Reference](https://pkg.go.dev/badge/github.com/knieriem/crcutil.svg)](https://pkg.go.dev/github.com/knieriem/crcutil)

# crcutil

When communicating with hardware devices like integrated circuits,
often CRC sums need to be calculated with polynomials of orders lower than 32, like 3, 4, 6, 8 or 16.
Sometimes the so-called reflected form of algorithms has to be used on the wire,
sometimes the normal form — or even both forms in the same protocol;
in some cases obscure additional reflections need to be performed.

As the Go standard library offers no support for these smaller widths of CRCs —
also its crc32 implementation implements solely the commonly used reflected form of the algorithm,
and does not provide a generic way to specify initial values or final XORing —,
this module tries to provide that functionality in a generic way.

Using the type `Poly`, a polynomial may be specified by its word value and bit width,
and whether the word refers to the [normal, reflected, or a reciprocal representation][polyreprs].
A value of type `Poly` may be converted to other representations, e.g. from normal to reflected form, using the corresponding method.
From a `Poly`, a lookup table may be created that may be used directly, if the amount of data is very small, like e.g. in case of a CRC-3.

[polyreprs]: https://en.wikipedia.org/wiki/Mathematics_of_cyclic_redundancy_checks#Polynomial_representations


A `Model` specifies a CRC algorithm: the polynomial,
the initial value (resp. optional initial inversion) to be used,
and whether a final operation like inversion (XORing) shall be applied. From a static `Model`,
an `Inst` may be created that allows to calculate a checksum over an amount of data.

There exist sub-packages `poly{n}` and `crc{n}`
that provide concrete versions of the generic `Poly` and `Model` types for some common values of bit widths and word types.
One example for a predefined `Model` is [`crc8.SAEJ1850`], which uses `poly8.SAEJ1850` with initial and final bitwise inversion. Another example is [`crc16.Modbus`], based on the reversed form of `poly16.IBM`.

[`crc8.SAEJ1850`]: https://pkg.go.dev/github.com/knieriem/crcutil/crc8#SAEJ1850
[`crc16.Modbus`]: https://pkg.go.dev/github.com/knieriem/crcutil/crc16#Modbus


## Example

As a real example, the LED driver circuit [Infineon TLD7002-16ES] uses both a CRC-3 in reflected form and a CRC-8 (SAEJ1850) in normal form in the same serial protocol.

### CRC-3

The CRC-3 polynomial x³ + x¹ + 1 (GSM, value = `0b11` = 3) can be created in the following equivalent ways:

```Go
p := poly3.GSM                          // using predefined polynomial
p := poly3.New(0b11)                    // 0b11 = 2¹ + 2°
p := &poly8.Poly{Word: 0b11, Width: 3}
p := &Poly[uint8]{Word: 0b11, Width: 3} // using the generic type
```

Since the 3-bit polynomial can be represented by an 8-bit value,
the `poly3` internally uses `poly8.Poly`,
like shown in the third line above.

If we wanted to create a lookup table for five bits of data,
with an initial value of 5 — as defined by the specification —
encoded into the table,
we could write

```Go
p := poly3.GSM.ReversedForm()
tab := p.MakeTable(crcutil.WithInitialValue(5), crcutil.WithDataWidth(5))
```

To get the CRC-sum of a 5-bit value of 13 it would be sufficient to
evaluate `tab[13]`.

For more details see the example for `func (*Poly[T]) MakeTable()` in the documentation. 

[Infineon TLD7002-16ES]: https://www.mouser.de/pdfDocs/Infineon-TLD7002-16ES-DataSheet-v01_00-EN.pdf#page=67


### CRC-8-SAE-J1850

To calculate a checksum over some bytes
using the SAE-J1850 polynomial x⁸ + x⁴ + x³ + x² + 1,
follow this example (compare the result the one listed at
[AUTOSAR Specification of CRC Routines, p.24]):

```Go
// create an instance
crc := crc8.SAEJ1850.New()

// insert some bytes
crc.Write([]byte{0xF2, 0x01, 0x83})

fmt.Printf("0x02x\n", crc.Sum())
// => 0x37
```

[AUTOSAR Specification of CRC Routines, p.24]: https://www.autosar.org/fileadmin/standards/R22-11/CP/AUTOSAR_SWS_CRCLibrary.pdf#page=24


## hash.Hash interface

For an implementation aligned with Go's `hash.Hash` interface, see [github.com/knieriem/hash], which is a thin wrapper
around `crcutil`.

[github.com/knieriem/hash]: https://pkg.go.dev/github.com/knieriem/hash
