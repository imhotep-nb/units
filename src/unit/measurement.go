package unit

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// Measurement represents a physical quantity: a value and a unit.
// The units have to be registered in the unit table with DefineUnit.
type Measurement struct {
	value float64
	*unit
}

// String returns a default string representation of the Measurement
func (m Measurement) String() string {
	return m.Format(DefaultFormat)
}

// Inspect returns a string representation of the Measurement for debugging
func (m Measurement) Inspect() string {
	return fmt.Sprintf("%f %s -> %f %s %v", m.value, m.symbol, m.factor, makeSymbol(m.exponents), m.exponents)
}

// Format returns a string representation of the Measurement according to the
// format string passed in. The first argument of the format string is the value,
// the second one is the unit. The unit and value can be swapped by using
// format string indexes such as in "%[2]s %.2[1]f". If only one argument is to be
// used, then an index must be used as well, e.g. "%[1]e radians".
// A better way to format measurements is by using a Context.
func (m Measurement) Format(format string) string {
	var a, b interface{}
	if m.unit == nil {
		a, b = m.value, "?"
	} else {
		a, b = m.value, m.symbol
	}
	return fmt.Sprintf(format, a, b)
}

// Split returns the value and the unit symbol of the Measurement
func (m Measurement) Split() (float64, string) {
	return m.value, m.symbol
}

// Value returns only the value part of the measurement. 
func (m Measurement) Value() float64 {
	return m.value
}

// ConvertTo creates and returns a new Measurement that has undergone conversion to the given unit.
// It also returns true/false to indicate success/failure. The conversion fails if the given unit
// cannot be found or calculated, or if that unit is not compatible.
func (m Measurement) ConvertTo(u string) (Measurement, bool) {
	target := get(u)
	compatible := haveSameExponents(m.exponents, target.exponents)
	if target == nil || !compatible {
		return Measurement{}, false
	}
	f := target.factor / m.factor
	return Measurement{m.value / f, target}, true
}

// In returns a Measurement converted to the given unit. No unit compatibility check is 
// performed. If the target unit is not compatible the function will return garbage.
func (m Measurement) In(u string) Measurement {
	target := get(u)
	return Measurement{m.value * m.factor / target.factor, target}
}

// M returns a Measurement with the given value and unit.
func M(value float64, symbol string) Measurement {
	u := get(symbol)
	if u == &UndefinedUnit {
		panic(fmt.Sprintf("undefined unit: %s", symbol))
	}
	return Measurement{value, u}
}

// Parse can be used to parse text input. The input is expected to contain a number
// followed by a unit string. Whitespace between number and unit string is optional.
// The number can have a negative sign and optional group separators (,). 
// The unit string has to be a registered unit symbol using the dot and slash to connect 
// factors, numbers for exponents and optional minus signs, e.g. "-1,500 N.m/s2" =
// -1500 newton meter per square second. This function returns the Measurement and an 
// error which is nil in case the string has been correctly parsed into a Measurement.
func Parse(s string) (Measurement, error) {
	undef := Measurement{0, &UndefinedUnit}
	match := muRx.FindStringSubmatch(s)
	if len(match) != 3 {
		return undef, errors.New("invalid measurement format [" + s + "]")
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
	return Measurement{value, mu.unit}, nil
}

// Invalid checks if the Measurement is valid, i.e. if it has a unit.
func (m Measurement) Invalid() bool {
	return m.unit == nil
}

// AreCompatible checks if two measurements are compatible. Compatibility means the exponents
// of the SI base units are the same. A return value of true means the measurements
// have compatible units.
func AreCompatible(a, b Measurement) bool {
	return haveSameExponents(a.exponents, b.exponents)
}

// HasCompatibleUnit check whether the Measurment can be converted to the given unit.
func (m Measurement) HasCompatibleUnit(symbol string) bool {
	return haveSameExponents(m.exponents, get(symbol).exponents)
}

func check(a, b Measurement) {
	if PanicOnIncompatibleUnits && !haveSameExponents(a.exponents, b.exponents) {
		panic(fmt.Sprintf("units not compatible: %q <> %q", a, b))
	}
}

// Add adds 2 Measurements that should have compatible units. If not compatible
// a panic happens or garbage is returned, depending on the setting of GOUNITSPANIC environment
// variable: 1 = panic, else no panic.
// The returned Measurement will be represented in SI units. This can be converted
// to the desired units with methods In or ConvertTo.
func Add(a, b Measurement) Measurement {
	check(a, b)
	u := &unit{"", 1, a.exponents}
	u.setSymbol()
	return Measurement{a.value*a.factor + b.value*b.factor, u}
}

// Sum adds one or more Measurements. The Measurements should have compatible units.
// If not compatible a panic happens or garbage is returned, depending on the setting 
// of GOUNITSPANIC environment variable: 1 = panic, else no panic.
func Sum(a Measurement, more ...Measurement) Measurement {
	return multi(a, func(m *float64, b Measurement) { *m += b.value * b.factor }, more)
}

// Subtract subtracts the second argument from the first one. Compatible units are required.
func Subtract(a, b Measurement) Measurement {
	return Add(a, Neg(b))
}

// Diff can be used to do multiple subtractions from the first argument. Compatible units are
// required.
func Diff(a Measurement, more ...Measurement) Measurement {
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

// Neg negates a Measurement value. The unit does not change.
func Neg(a Measurement) Measurement {
	return Measurement{-a.value, a.unit}
}

// Mult multiplies 2 Measurements. A new unit will be calculated. The returned Measurement will
// have SI units. Use In or ConvertTo to convert it to the desired unit.
func Mult(a, b Measurement) Measurement {
	return Measurement{a.value * a.factor * b.value * b.factor, addu(a.unit, b.unit)}
}

// Div divides the first argument by the second. A new unit will be calculated. 
// The returned Measurement will have SI units. Use In or ConvertTo to convert it to the desired unit.
func Div(a, b Measurement) Measurement {
	return Measurement{(a.value * a.factor) / (b.value * b.factor), subu(a.unit, b.unit)}
}

// Reciprocal calculates 1 divided by the given Measurement. The unit changes accordingly but
// will be represented in SI units. 
func Reciprocal(a Measurement) Measurement {
	u := &unit{"", 1, negx(a.exponents)}
	u.setSymbol()
	return Measurement{1 / (a.value * a.factor), u}
}

// MultFac multiplies a Measurement with a factor and returns the new Measurement. The unit
// does not change.
func MultFac(m Measurement, f float64) Measurement {
	return Measurement{m.value * f, m.unit}
}

// DivFac divides a Measurement by a factor and returns the new Measurement. The unit does not
// change.
func DivFac(m Measurement, f float64) Measurement {
	return Measurement{m.value / f, m.unit}
}

// Power raises the Measurement to the given power n. The exponents of the resulting unit must
// be in the range -128..127.
func Power(a Measurement, n int8) Measurement {
	calc := func(e int8) int8 { return e * n }
	u := &unit{"", 1, mapexp(a.exponents, calc)}
	u.setSymbol()
	return Measurement{math.Pow(a.value*a.factor, float64(n)), u}
}

// Abs returns the absolute of Measurement: the result is always >= 0.
func Abs(a Measurement) Measurement {
	if a.value < 0 {
		return Neg(a)
	}
	return a
}

// Equal checks if two Measurements are equal. A tolerance epsilon is allowed, this value should
// be much smaller compared to the two Measurements being compared. All arguments must have 
// compatible units.
func Equal(a, b, epsilon Measurement) bool {
	check(a, b)
	check(a, epsilon)
	return Abs(Subtract(a, b)).value < epsilon.value*epsilon.factor
}

// More checks if the first argument is greater than the second.
func More(a, b Measurement) bool {
	check(a, b)
	return a.ToSI().Value() > b.ToSI().Value()
}

// Less checks if the first argument is less than the second.
func Less(a, b Measurement) bool {
	check(a, b)
	return a.ToSI().Value() < b.ToSI().Value()
}

// ToSI returns a converted Measurement represented in SI units.
func (m Measurement) ToSI() Measurement {
	factor, u := m.toSI()
	return Measurement{m.value * factor, &u}
}

// Normalize changes the Measurement to SI units.
func (m *Measurement) Normalize() {
	m.value *= m.factor
	m.unit = &unit{makeSymbol(m.exponents), 1, m.exponents}
}

// Duration converts a Measurement with a duration unit to a time.Duration.
// An error or nil is provided as second return value.
func Duration(m Measurement) (time.Duration, error) {
	if si, ok := m.ConvertTo("s"); ok {
		return time.Duration(si.Value()) * time.Second, nil
	}
	return time.Duration(0), errors.New("not a Duration: " + m.String())
}

// Slice of Measurements. Useful for sorting.
type MeasurementSlice []Measurement

// Len is used by Sort
func (a MeasurementSlice) Len() int {
	return len(a)
}

// Swap is used by Sort
func (a MeasurementSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less is used by Sort
func (a MeasurementSlice) Less(i, j int) bool {
	return Less(a[i], a[j])
}
