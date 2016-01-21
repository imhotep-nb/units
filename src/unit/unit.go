package unit

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	nBaseUnits = 11
)

const (
	metre = iota
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

var (
	DefaultFormat            = "%.4f %s"
	UndefinedUnit            = unit{"?", 1, emptyExponents()}
	baseSymbols              = []string{"m", "kg", "K", "A", "cd", "mol", "rad", "sr", "¤", "byte", "s"}
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
	d := [nBaseUnits]int8{}
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
		m, err := ParseSymbol(symbol)
		if err != nil {
			u = &UndefinedUnit
		} else {
			u = m.unit
			units[u.symbol] = u // cache it
		}
	}
	return u
}

func isCompatible(x, y []int8) bool {
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

func ParseSymbol(s string) (Measurement, error) {
	resultSI := Measurement{1.0, units[""]}
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
			mSI := Measurement{factor, &uSI}
			if match[2] != "" {
				x, _ = strconv.Atoi(match[2])
				if i == 1 && x < 0 {
					return resultSI, errors.New("invalid format: negative exponent after the '/'")
				}
				mSI = Power(mSI, int8(x))
				//fmt.Println("x", x, "m^x", mSI.Format("%f %s"))
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
	fmt.Print()
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
