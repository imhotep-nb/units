# units
Unit system library for Go.

__Please don't use in production. I'm still working on this.__

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

```
There are more functions and methods. See `measurement.go` and `unit.go`.

The units are defined in the file `data.go`. I will extend this file with more units. 

The `Measurement` structs consist of a `float64` value and `*unit`; the unit may or may not be shared with other 
measurements, but is not public, and from the point of view of the client code, is immutable. The measurement remembers
its original unit. An `Add`, `Subtract`, `Mult` or `Div` will always return an SI unit though, but this can be converted to another compatible unit with `In(string)` or `ConvertTo(string)`, the latter doing compatibility checking. `In` will produce garbage if the unit is not compatible and won't warn you.

The internal storage of a unit consists of a struct with a symbol (e.g. "km/h", a conversion factor (1 for SI units) and a slice of 10 exponents ([]int8) for the SI units, and a few more handy ones. E.g. there is a exponent for currency, to allow 
for currency conversions, but I still have to do some work on handling dynamic factors for these (exchange rates). 


##todo
 * more tests
 * more units, the data.go only has a small set of units for testing
 * currently no support for prefixes (k for kilo, M for mega etc.) consider adding
 * exchange rate handling
 * more control over formatting and decimal places; add 'contexts' e.g. height vs length
 * conversion of unit.M(x, "s") to existing time.Duration
 * add degrees/minutes/seconds parsing
 * add degrees C and F and special conversions for these (formulas are not captured by simple factor)
 * godoc documentation
 * make Measurement sortable
 * parseing/printing of unitless


