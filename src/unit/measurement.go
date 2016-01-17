package unit

import (
	"fmt"
	"math"
)

type Measurement struct {
	value float64
	*unit
}

func (m Measurement) String() string {
	return m.Format(DefaultFormat)
}

func (m Measurement) Inspect() string {
	return fmt.Sprintf("%f %s -> %f %s %v", m.value, m.symbol, m.factor, makeSymbol(m.exponents), m.exponents)
}

func (m Measurement) Format(format string) string {
	var a, b interface{}
	if m.unit == nil {
		a, b = m.value, "?"
	} else {
		a, b = m.value, m.symbol
	}
	return fmt.Sprintf(format, a, b)
}

func (m Measurement) Split() (float64, string) {
	return m.value, m.symbol
}

func (m Measurement) Value() float64 {
	return m.value
}

// Unit conversion
func (m Measurement) ConvertTo(u string) (Measurement, bool) {
	target := get(u)
	compatible := isCompatible(m.exponents, target.exponents)
	if target == nil || !compatible {
		return Measurement{}, false
	}
	f := target.factor / m.factor
	return Measurement{m.value / f, target}, true
}

func (m Measurement) In(u string) Measurement {
	target := get(u)
	return Measurement{m.value * m.factor / target.factor, target}
}

func M(value float64, symbol string) Measurement {
	return Measurement{value, get(symbol)}
}

func (m Measurement) Invalid() bool {
	return m.unit == nil
}

func SameUnit(a, b Measurement) bool {
	return isCompatible(a.exponents, b.exponents)
}

func check(a, b Measurement) {
	if PanicOnIncompatibleUnits && !isCompatible(a.exponents, b.exponents) {
		panic(fmt.Sprintf("units not compatible: %q <> %q", a, b))
	}
}

func Add(a, b Measurement) Measurement {
	check(a, b)
	u := &unit{"", 1, a.exponents}
	u.setSymbol()
	return Measurement{a.value*a.factor + b.value*b.factor, u}
}

func AddN(a Measurement, more ...Measurement) Measurement {
	return multi(a, func(m *float64, b Measurement) { *m += b.value * b.factor }, more)
}

func Subtract(a, b Measurement) Measurement {
	return Add(a, Neg(b))
}

func SubtractN(a Measurement, more ...Measurement) Measurement {
	return multi(a, func(m *float64, b Measurement) { *m -= b.value * b.factor }, more)
}

func multi(
	a Measurement,
	op func(*float64, Measurement),
	more []Measurement) Measurement {

	result := a.value * a.factor
	for _, b := range more {
		check(a, b)
		op(&result, b)
	}
	u := &unit{"", 1, a.exponents}
	u.setSymbol()
	return Measurement{result, u}
}

func Neg(a Measurement) Measurement {
	return Measurement{-a.value, a.unit}
}

func Mult(a, b Measurement) Measurement {
	return Measurement{a.value * a.factor * b.value * b.factor, addu(a.unit, b.unit)}
}

func Div(a, b Measurement) Measurement {
	return Measurement{(a.value * a.factor) / (b.value * b.factor), subu(a.unit, b.unit)}
}

func Reciprocal(a Measurement) Measurement {
	u := &unit{"", 1, negx(a.exponents)}
	u.setSymbol()
	return Measurement{1 / (a.value * a.factor), u}
}

func MultF(m Measurement, f float64) Measurement {
	return Measurement{m.value * f, m.unit}
}

func DivF(m Measurement, f float64) Measurement {
	return Measurement{m.value / f, m.unit}
}

func Power(a Measurement, n int8) Measurement {
	calc := func(e int8) int8 { return e * n }
	u := &unit{"", 1, mapexp(a.exponents, calc)}
	u.setSymbol()
	return Measurement{math.Pow(a.value*a.factor, float64(n)), u}
}

func Abs(a Measurement) Measurement {
	if a.value < 0 {
		return Neg(a)
	}
	return a
}

func Equal(a, b, epsilon Measurement) bool {
	check(a, b)
	check(a, epsilon)
	return Abs(Subtract(a, b)).value < epsilon.value*epsilon.factor
}

func (m Measurement) ToSI() Measurement {
	factor, u := m.toSI()
	return Measurement{m.value * factor, &u}
}


func (m *Measurement) Normalize() {
	m.value *= m.factor
	m.unit = &unit{makeSymbol(m.exponents), 1, m.exponents}
}
