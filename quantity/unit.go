// Package quantity provides a way to express and work with physical quantities, or measurements.
// A Quantity consists of a value and a unit.
package quantity

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	meter = iota
	kilogram
	kelvin
	ampere
	candela
	mole
	radian
	steradian
	currency
	byte
	second
	// when inserting a new base unit, then also update baseSymbols below
)

const (
	nBaseUnits = 11
)

const (
	yocto float64 = 1e-24
	zepto         = 1e-21
	atto          = 1e-18
	femto         = 1e-15
	pico          = 1e-12
	nano          = 1e-9
	micro         = 1e-6
	milli         = 1e-3
	centi         = 0.01
	deci          = 0.1
	deca          = 10
	hecto         = 100
	kilo          = 1e3
	mega          = 1e6
	giga          = 1e9
	tera          = 1e12
	peta          = 1e15
	exa           = 1e18
	zetta         = 1e21
	yotta         = 1e24
)

var (
	// DefaultFormat is the default formatstring for Quantities
	DefaultFormat = "%.4f %s"
	// UndefinedUnit represents a unit that is unknown to the system
	UndefinedUnit = Unit{"?", 0, emptyExponents()}
	// PanicOnIncompatibleUnits panic if operation with incompatible units happens
	PanicOnIncompatibleUnits = os.Getenv("GOUNITSPANIC") == "1"

	baseSymbols    = [nBaseUnits]string{"m", "kg", "K", "A", "cd", "mol", "rad", "sr", "¤", "byte", "s"}
	prefixValues   = [...]float64{deci, centi, hecto, milli, kilo, micro, mega, nano, giga, pico, tera, femto, peta, atto, exa, zepto, zetta, yotta, yocto}
	prefixSymbols  = "dchmkuMnGpTfPaEzZyY"
	symbolRx, muRx *regexp.Regexp
)

// Unit represents a unit of measure.
type Unit struct {
	symbol    string
	factor    float64
	exponents []int8
}

func def(dim *[nBaseUnits]int8) func(string, float64) *Unit {
	return func(symbol string, factor float64) *Unit {
		return &Unit{symbol, factor, dim[:]}
	}
}

func mapexp(e []int8, f func(int8) int8) []int8 {
	e1 := [nBaseUnits]int8{}
	for i := 0; i < nBaseUnits; i++ {
		e1[i] = f(e[i])
	}
	return e1[:]
}

// Symbol gets the string that represents the unit
func (u *Unit) Symbol() string {
	return u.symbol
}

func addu(a, b *Unit) *Unit {
	u := &Unit{"", a.factor * b.factor, addx(a.exponents, b.exponents)}
	u.symbol = makeSymbol(u.exponents)
	return u
}

func subu(a, b *Unit) *Unit {
	u := &Unit{"", a.factor / b.factor, addx(a.exponents, negx(b.exponents))}
	u.symbol = makeSymbol(u.exponents)
	return u
}

func addx(a, b []int8) []int8 {
	r := [nBaseUnits]int8{}
	for i := 0; i < nBaseUnits; i++ {
		r[i] = a[i] + b[i]
	}
	return r[:]
}

func negx(a []int8) []int8 {
	return mapexp(a, func(e int8) int8 { return -e })
}

func (u Unit) rcp() Unit {
	u.exponents = negx(u.exponents)
	return u
}

func (u *Unit) setSymbol() {
	u.symbol = makeSymbol(u.exponents)
}

func makeSymbol(expon []int8) string {
	var a []string
	for i := 0; i < nBaseUnits; i++ {
		e := expon[i]
		if e != 0 {
			a = append(a, "."+baseSymbols[i])
			if e != 1 {
				a = append(a, strconv.Itoa(int(e)))
			}
		}
	}
	if len(a) == 0 {
		return "?"
	}
	return strings.Join(a, "")[1:]
}

var units = make(map[string]*Unit)

// UnitFor looks up or construct a unit ref from a given symbol
func UnitFor(symbol string) *Unit {
	u := units[symbol]
	//fmt.Println("found in cache [", symbol, "] -> ", u)
	if u == nil {
		q, err := ParseSymbol(symbol)
		if err != nil {
			u = &UndefinedUnit
		} else {
			u = q.Unit
			units[u.symbol] = u // cache it
		}
	}
	return u
}

func prefix(symbol string) (f float64, base string, ok bool) {
	if len(symbol) < 2 {
		return 0, "", false
	}

	if len(symbol) > 2 && symbol[:2] == "da" {
		f = deca
		base = symbol[2:]
		ok = true
	} else {
		i := strings.IndexByte(prefixSymbols, symbol[0])
		if i != -1 {
			f = prefixValues[i]
			base = symbol[1:]
			ok = true
		}
	}
	if ok {
		u, found := units[base]
		if found {
			switch {
			case u.symbol == "g":
				f /= 1000
				base = "kg"
			case u.factor != 1 || strings.Contains(u.symbol, " "):
				ok = false
			}
		} else {
			ok = false
		}
	}
	return
}

func haveSameExponents(x, y []int8) bool {
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func emptyExponents() []int8 {
	x := [nBaseUnits]int8{}
	return x[:]
}

func (u Unit) toSI() (factor float64, si Unit) {
	si = Unit{"", 1, u.exponents}
	si.setSymbol()
	return u.factor, si
}

// ParseSymbol parses the given unit and returns a Quantity with the value set to 1.
func ParseSymbol(s string) (Quantity, error) {
	s = strings.ReplaceAll(s, "*", ".")
	s = strings.ReplaceAll(s, "^", "")
	resultSI := Quantity{1.0, units[""]}
	parts := strings.Split(s, "/")
	if len(parts) > 2 {
		return resultSI, errors.New("more than one '/' in unit")
	}

	for i, part := range parts {
		for _, symbol := range strings.Split(part, ".") {
			match := symbolRx.FindStringSubmatch(symbol)
			//fmt.Println("match", match)
			if len(match) != 3 {
				return resultSI, errors.New("cannot parse unit [" + s + "]")
			}
			u := units[match[1]]
			var pf float64 = 1
			if u == nil {
				p, baseUnit, ok := prefix(match[1])
				if !ok {
					return resultSI, errors.New("unknown symbol [" + match[1] + "]")
				}
				u = units[baseUnit]
				pf = p
			}
			factor, uSI := u.toSI()
			var x int
			mSI := Quantity{pf * factor, &uSI}
			if match[2] != "" {
				x, _ = strconv.Atoi(match[2])
				if i == 1 && x < 0 {
					return resultSI, errors.New("invalid format: negative exponent after the '/'")
				}
				mSI = Power(mSI, int8(x))
				//fmt.Println("x", x, "q^x", mSI.Format("%f %s"))
			}
			if i == 0 {
				resultSI = Mult(resultSI, mSI)
			} else {
				resultSI = Div(resultSI, mSI)
			}
			//fmt.Println("result so far", resultSI.value, resultSI.factor, resultSI.symbol, resultSI.exponents)
		}
	}
	resultSI.factor, resultSI.value = resultSI.value, resultSI.factor
	resultSI.symbol = s
	//fmt.Println("final result", resultSI.value, resultSI.factor, resultSI.symbol, resultSI.exponents)
	return resultSI, nil
}

// Define can be used to add a new unit to the unit table.
// The new unit symbol must be unique, the base symbol must either exist or be a calculation
// based on other units, e.g. "kg.q/s2", but not necessarily SI. 1 new unit = factor * base unit.
func Define(symbol string, factor float64, base string) (float64, error) {
	if _, found := units[symbol]; found {
		return 0, errors.New("duplicate symbol [" + symbol + "]")
	}
	mBase, err := ParseSymbol(base)
	if err != nil {
		return 0, err
	}
	siFactor := factor * mBase.factor
	units[symbol] = &Unit{symbol, siFactor, mBase.exponents}
	return siFactor, nil
}

func init() {
	fmt.Print("")
	symbolRx = regexp.MustCompile(`^([^\d-]+)(-?\d+)?$`)
	muRx = regexp.MustCompile(`^\s*(-?[\d.,]+)\s*(.*)$`)

	data := setup()
	for _, value := range data {
		if units[value.symbol] != nil {
			panic("duplicate unit symbol")
		}
		units[value.symbol] = value
	}
}
