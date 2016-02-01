# units
Physical units system library for Go.

One of my first pieces of Go software. I developed this to get a feel for the language, but I will keep improving this
in 2016.

## purpose
With this library you can express measurements as value + physical units and you can do calculations with these
measurements. The units are automatically handled in these calculations. Unit conversions between metric / imperial and US
units are taken care of. 

More features will be added, see todo list below.

## examples
Here are some usage examples:

```go
// create a measurement with unit.M(value float64, unit string)
	m1 := unit.M(33, "N/m2")
	
// conversions to other compatible units can be done with ConvertTo(string)
// which returns a new measurement with the new unit, and 
	if m, ok := unit.M(6894.757, "N/m2").ConvertTo("lbf.in-2"); ok {
		fmt.Println(m)
	} else {
		fmt.Println("error")
	}
// add your own units with Define(newUnitSymbol string, factorForBaseUnit float64, baseUnit string)
	siFactor, err := unit.Define("foo", 7, "lbf/sq in")
	
// return an SI version of the unit
	m2 = m.ToSI()
	
// split into value and unit
	t := unit.M(12345, "psi")
	value, symbol := t.Split()
	
// conversion
	t1 := t.In("bar")
	value, ok := t.ConvertTo("bar")
	
// change unit of the measurement to SI
	t.Normalize()
	
// math operations are Add, AddN, Subtract, SubtractN, Mult, MultF, Div, DivF, Neg, Power
	s := unit.Add(unit.M(3.5, "km"), unit.M(1.2, "mi")).In("ft").String()
	
// return a measurement with the given unit; calculate new conversion factor.
	measurement, err := unit.ParseSymbol("psi.kg-1")
	
// parse user input
	measurement, err := unit.Parse(" -1,234,566.88 sq in/min  ")
	
// create a Context for maintaining a certain unit and output formatting
	unit.DefineContext(personHeight, "cm", "%.0[1]fcm")
	height := unit.Ctx(personHeight)
	m := height.M(1.75, "m")
	s := height.String(m) // -> "175cm"

// other Contexts e.g.
	unit.DefineContext(landArea, "acre", "%0.[1]f acres")
	unit.DefineContext(money, "$", "%[2]s%.2[1]f") // unit before value
	unit.DefineContext(rainIntensity, "mm/h", "%.1f %s")

//----------
	
// optionally create static types like this:
type Area struct {
	unit.Measurement
	*unit.Context
}

var area, _ = unit.DefineContext("landArea", "acre", "%.1f %s")

func NewArea(m unit.Measurement) (Area, error) {
	if !m.HasCompatibleUnit(area.Symbol()) {
		return Area{}, errors.New(fmt.Sprintf("%v is not a %s", m, area.Name))
	}
	return Area{m, area}, nil
}

func (a Area) String() string {
	return a.Context.String(a.Measurement)
}
...
	a, err := NewArea(unit.M(2.5, "km2"))
	fmt.Println(a.String()) // => 617.8 acre
//----------

// package units/resource
// create a Measurement based resource with a min and max value.
// you can Deposit and Withdraw from the resource.

	rsc := resource.New(unit.M(0, "kWh"), unit.M(500, "kWh"), "")
	rsc.Set(unit.M(25, "kWh"))
	rsc.Withdraw(unit.M(1e6, "J"))
	rsc.Deposit(unit.M(1.23 * unit.Kilo, "J"))	
	b := rsc.Balance() // => 24.7226 kWh



```
See also the 'test' folder for more examples of how to use the packages.

There are more functions and methods. See `measurement.go` and `unit.go`.

The units are defined in the file `data.go`. I will extend this file with more units. 

The `Measurement` structs consist of a `float64` value and `*unit`; the unit may or may not be shared with other 
measurements, but is not public, and from the point of view of the client code, is immutable. The measurement remembers
its original unit. An `Add`, `Subtract`, `Mult` or `Div` will always return an SI unit though, but this can be converted to another compatible unit with `In(string)` or `ConvertTo(string)`, the latter doing compatibility checking. `In` will produce garbage if the unit is not compatible and won't warn you.

The internal storage of a unit consists of a struct with a symbol (e.g. "km/h", a conversion factor (1 for SI units) and a slice of 10 exponents ([]int8) for the SI units, and a few more handy ones. E.g. there is a exponent for currency, to allow 
for currency conversions, but I still have to do some work on handling dynamic factors for these (exchange rates). 



##todo
 * a few more units, the data.go only has a small set of units for testing
 * change support for prefixes(k for kilo, M for mega etc.) by parsing unit
 * exchange rate handling
 * add degrees/minutes/seconds parsing
 * add degrees C and F and special conversions for these (formulas are not captured by simple factor)
 * parsing/printing of unitless
 * parsing of combined units such as "5ft 10in"
 * resource: make safe for concurrent access by goroutines?


