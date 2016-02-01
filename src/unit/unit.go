// Package unit provides a way to express and work with physical quantities, or measurements.
// A Quantity consists of a value and a unit.
package unit

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"math"
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
	Yocto = 1e-24
	Zepto = 1e-21
	Atto  = 1e-18
	Femto = 1e-15
	Pico  = 1e-12
	Nano  = 1e-9
	Micro = 1e-6
	Milli = 1e-3
	Centi = 0.01
	Deci  = 0.1
	Deca  = 10
	Hecto = 100
	Kilo  = 1e3
	Mega  = 1e6
	Giga  = 1e9
	Tera  = 1e12
	Peta  = 1e15
	Exa   = 1e18
	Zetta = 1e21
	Yotta = 1e24
)

// Square can be used to apply to a SI metric prefix, e.g. unit.Square(unit.Deci)
func Square(f float64) float64 {
	return f * f
}

// Cubic can be used to apply to a SI metric prefix, e.g. unit.Cubic(unit.Deci)
func Cubic(f float64) float64 {
	return f * f * f
}

// Pow can be used to apply to a SI metric prefix, e.g. unit.Pow(unit.Deci, 4)
func Pow(f float64, exp int8) float64 {
	return math.Pow(f, float64(exp))
}

var (
	DefaultFormat            = "%.4f %s"
	UndefinedUnit            = unit{"?", 0, emptyExponents()}
	baseSymbols              = [nBaseUnits]string{"m", "kg", "K", "A", "cd", "mol", "rad", "sr", "¤", "byte", "s"}
	PanicOnIncompatibleUnits = os.Getenv("GOUNITSPANIC") == "1"
	symbolRx, muRx           *regexp.Regexp
)

type unit struct {
	symbol    string
	factor    float64
	exponents []int8
}

type expMap map[int]int8

func exp(u expMap) []int8 {
	d := [len(baseSymbols)]int8{}
	for k, v := range u {
		d[k] = v
	}
	return d[:]
}

func def(exponents expMap) func(string, float64) *unit {
	return func(symbol string, factor float64) *unit {
		return &unit{symbol, factor, exp(exponents)}
	}
}

func mapexp(e []int8, f func(int8) int8) []int8 {
	e1 := [nBaseUnits]int8{}
	for i := 0; i < nBaseUnits; i++ {
		e1[i] = f(e[i])
	}
	return e1[:]
}

func (u *unit) Symbol() string {
	return u.symbol
}

func addu(a, b *unit) *unit {
	u := &unit{"", a.factor * b.factor, addx(a.exponents, b.exponents)}
	u.symbol = makeSymbol(u.exponents)
	return u
}

func subu(a, b *unit) *unit {
	u := &unit{"", a.factor / b.factor, addx(a.exponents, negx(b.exponents))}
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

func (u unit) rcp() unit {
	u.exponents = negx(u.exponents)
	return u
}

func (u *unit) setSymbol() {
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

var units = make(map[string]*unit)

// Look up or construct a unit ref from a given symbol
func get(symbol string) *unit {
	u := units[symbol]
	//fmt.Println("found in cache [", symbol, "] -> ", u)
	if u == nil {
		q, err := ParseSymbol(symbol)
		if err != nil {
			u = &UndefinedUnit
		} else {
			u = q.unit
			units[u.symbol] = u // cache it
		}
	}
	return u
}

func haveSameExponents(x, y []int8) bool {
	for i, _ := range x {
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

func (u unit) toSI() (factor float64, si unit) {
	si = unit{"", 1, u.exponents}
	si.setSymbol()
	return u.factor, si
}


func ParseSymbol(s string) (Quantity, error) {
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
			if u == nil {
				return resultSI, errors.New("unknown symbol [" + match[1] + "]")
			}
			factor, uSI := u.toSI()
			var x int
			mSI := Quantity{factor, &uSI}
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
	units[symbol] = &unit{symbol, siFactor, mBase.exponents}
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
