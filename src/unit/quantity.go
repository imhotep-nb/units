package unit

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// Quantity represents a physical quantity: a value and a unit.
// The units have to be registered in the unit table with DefineUnit.
type Quantity struct {
	value float64
	*unit
}

// String returns a default string representation of the Quantity
func (m Quantity) String() string {
	return m.Format(DefaultFormat)
}

// Inspect returns a string representation of the Quantity for debugging
func (m Quantity) Inspect() string {
	return fmt.Sprintf("%f %s -> %f %s %v", m.value, m.symbol, m.factor, makeSymbol(m.exponents), m.exponents)
}

// Format returns a string representation of the Quantity according to the
// format string passed in. The first argument of the format string is the value,
// the second one is the unit. The unit and value can be swapped by using
// format string indexes such as in "%[2]s %.2[1]f". If only one argument is to be
// used, then an index must be used as well, e.g. "%[1]e radians".
// A better way to format quantities is by using a Context.
func (m Quantity) Format(format string) string {
	var a, b interface{}
	if m.unit == nil {
		a, b = m.value, "?"
	} else {
		a, b = m.value, m.symbol
	}
	return fmt.Sprintf(format, a, b)
}

// Split returns the value and the unit symbol of the Quantity
func (m Quantity) Split() (float64, string) {
	return m.value, m.symbol
}

// Value returns only the value part of the Quantity.
func (m Quantity) Value() float64 {
	return m.value
}

// ConvertTo creates and returns a new Quantity that has undergone conversion to the given unit.
// It also returns true/false to indicate success/failure. The conversion fails if the given unit
// cannot be found or calculated, or if that unit is not compatible.
func (m Quantity) ConvertTo(u string) (Quantity, bool) {
	target := get(u)
	compatible := haveSameExponents(m.exponents, target.exponents)
	if target == nil || !compatible {
		return Quantity{}, false
	}
	f := target.factor / m.factor
	return Quantity{m.value / f, target}, true
}

// In returns a Quantity converted to the given unit. No unit compatibility check is
// performed. If the target unit is not compatible the function will return garbage.
func (m Quantity) In(u string) Quantity {
	target := get(u)
	return Quantity{m.value * m.factor / target.factor, target}
}

// Q returns a Quantity with the given value and unit.
func Q(value float64, symbol string) Quantity {
	u := get(symbol)
	if u == &UndefinedUnit {
		panic(fmt.Sprintf("undefined unit: %s", symbol))
	}
	return Quantity{value, u}
}

// Parse can be used to parse text input. The input is expected to contain a number
// followed by a unit string. Whitespace between number and unit string is optional.
// The number can have a negative sign and optional group separators (,).
// The unit string has to be a registered unit symbol using the dot and slash to connect
// factors, numbers for exponents and optional minus signs, e.g. "-1,500 N.m/s2" =
// -1500 newton meter per square second. This function returns the Quantity and an
// error which is nil in case the string has been correctly parsed into a Quantity.
func Parse(s string) (Quantity, error) {
	undef := Quantity{0, &UndefinedUnit}
	match := muRx.FindStringSubmatch(s)
	if len(match) != 3 {
		return undef, errors.New("invalid quantity format [" + s + "]")
	}
	f := match[1]
	if strings.Count(f, ".") > 1 {
		return undef, errors.New("more than one decimal point in [" + s + "]")
	}
	f = strings.Replace(f, ",", "", -1)
	value, err := strconv.ParseFloat(f, 64)
	if err != nil {
		return undef, err
	}
	sym := strings.Trim(match[2], " \r\n\t")
	mu, err := ParseSymbol(sym)
	if err != nil {
		return undef, err
	}
	return Quantity{value, mu.unit}, nil
}

// Invalid checks if the Quantity is valid, i.e. if it has a unit.
func (m Quantity) Invalid() bool {
	return m.unit == nil
}

// AreCompatible checks if two quantities are compatible. Compatibility means the exponents
// of the SI base units are the same. A return value of true means the quantities
// have compatible units.
func AreCompatible(a, b Quantity) bool {
	return haveSameExponents(a.exponents, b.exponents)
}

// HasCompatibleUnit check whether the Measurment can be converted to the given unit.
func (m Quantity) HasCompatibleUnit(symbol string) bool {
	return haveSameExponents(m.exponents, get(symbol).exponents)
}

func check(a, b Quantity) {
	if PanicOnIncompatibleUnits && !haveSameExponents(a.exponents, b.exponents) {
		panic(fmt.Sprintf("units not compatible: %q <> %q", a, b))
	}
}

// Add adds 2 Quantities that should have compatible units. If not compatible
// a panic happens or garbage is returned, depending on the setting of GOUNITSPANIC environment
// variable: 1 = panic, else no panic.
// The returned Quantity will be represented in SI units. This can be converted
// to the desired units with methods In or ConvertTo.
func Add(a, b Quantity) Quantity {
	check(a, b)
	u := &unit{"", 1, a.exponents}
	u.setSymbol()
	return Quantity{a.value*a.factor + b.value*b.factor, u}
}

// Sum adds one or more Quantities. The Quantities should have compatible units.
// If not compatible a panic happens or garbage is returned, depending on the setting
// of GOUNITSPANIC environment variable: 1 = panic, else no panic.
func Sum(a Quantity, more ...Quantity) Quantity {
	return multi(a, func(m *float64, b Quantity) { *m += b.value * b.factor }, more)
}

// Subtract subtracts the second argument from the first one. Compatible units are required.
func Subtract(a, b Quantity) Quantity {
	return Add(a, Neg(b))
}

// Diff can be used to do multiple subtractions from the first argument. Compatible units are
// required.
func Diff(a Quantity, more ...Quantity) Quantity {
	return multi(a, func(m *float64, b Quantity) { *m -= b.value * b.factor }, more)
}

func multi(
	a Quantity,
	op func(*float64, Quantity),
	more []Quantity) Quantity {

	result := a.value * a.factor
	for _, b := range more {
		check(a, b)
		op(&result, b)
	}
	u := &unit{"", 1, a.exponents}
	u.setSymbol()
	return Quantity{result, u}
}

// Neg negates a Quantity value. The unit does not change.
func Neg(a Quantity) Quantity {
	return Quantity{-a.value, a.unit}
}

// Mult multiplies 2 Quantities. A new unit will be calculated. The returned Quantity will
// have SI units. Use In or ConvertTo to convert it to the desired unit.
func Mult(a, b Quantity) Quantity {
	return Quantity{a.value * a.factor * b.value * b.factor, addu(a.unit, b.unit)}
}

// Div divides the first argument by the second. A new unit will be calculated.
// The returned Quantity will have SI units. Use In or ConvertTo to convert it to the desired unit.
func Div(a, b Quantity) Quantity {
	return Quantity{(a.value * a.factor) / (b.value * b.factor), subu(a.unit, b.unit)}
}

// Reciprocal calculates 1 divided by the given Quantity. The unit changes accordingly but
// will be represented in SI units.
func Reciprocal(a Quantity) Quantity {
	u := &unit{"", 1, negx(a.exponents)}
	u.setSymbol()
	return Quantity{1 / (a.value * a.factor), u}
}

// MultFac multiplies a Quantity with a factor and returns the new Quantity. The unit
// does not change.
func MultFac(m Quantity, f float64) Quantity {
	return Quantity{m.value * f, m.unit}
}

// DivFac divides a Quantity by a factor and returns the new Quantity. The unit does not
// change.
func DivFac(m Quantity, f float64) Quantity {
	return Quantity{m.value / f, m.unit}
}

// Power raises the Quantity to the given power n. The exponents of the resulting unit must
// be in the range -128..127.
func Power(a Quantity, n int8) Quantity {
	calc := func(e int8) int8 { return e * n }
	u := &unit{"", 1, mapexp(a.exponents, calc)}
	u.setSymbol()
	return Quantity{math.Pow(a.value*a.factor, float64(n)), u}
}

// Abs returns the absolute of Quantity: the result is always >= 0.
func Abs(a Quantity) Quantity {
	if a.value < 0 {
		return Neg(a)
	}
	return a
}

// Equal checks if two Quantities are equal. A tolerance epsilon is allowed, this value should
// be much smaller compared to the two Quantities being compared. All arguments must have
// compatible units.
func Equal(a, b, epsilon Quantity) bool {
	check(a, b)
	check(a, epsilon)
	return Abs(Subtract(a, b)).value < epsilon.value*epsilon.factor
}

// More checks if the first argument is greater than the second.
func More(a, b Quantity) bool {
	check(a, b)
	return a.ToSI().Value() > b.ToSI().Value()
}

// Less checks if the first argument is less than the second.
func Less(a, b Quantity) bool {
	check(a, b)
	return a.ToSI().Value() < b.ToSI().Value()
}

// ToSI returns a converted Quantity represented in SI units.
func (m Quantity) ToSI() Quantity {
	factor, u := m.toSI()
	return Quantity{m.value * factor, &u}
}

// Normalize changes the Quantity to SI units.
func (m *Quantity) Normalize() {
	m.value *= m.factor
	m.unit = &unit{makeSymbol(m.exponents), 1, m.exponents}
}

// Duration converts a Quantity with a duration unit to a time.Duration.
// An error or nil is provided as second return value.
func Duration(m Quantity) (time.Duration, error) {
	if si, ok := m.ConvertTo("s"); ok {
		return time.Duration(si.Value()) * time.Second, nil
	}
	return time.Duration(0), errors.New("not a Duration: " + m.String())
}

// Slice of Quantity values. Useful for sorting.
type Quantities []Quantity

// Len is used by Sort
func (a Quantities) Len() int {
	return len(a)
}

// Swap is used by Sort
func (a Quantities) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less is used by Sort
func (a Quantities) Less(i, j int) bool {
	return Less(a[i], a[j])
}
