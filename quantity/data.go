package quantity

import (
	"math"
)

func setup() []*Unit {
	// keep alphabetic order!
	// only define quantities here that have a unit symbol that is not a combination of existing unit symbols
	acceleration := def(&[nBaseUnits]int8{meter: 1, second: -2})
	angle := def(&[nBaseUnits]int8{radian: 1})
	angularVelocity := def(&[nBaseUnits]int8{radian: 1, second: -1})
	area := def(&[nBaseUnits]int8{meter: 2})
	capacitance := def(&[nBaseUnits]int8{ampere: 2, second: 4, kilogram: -1, meter: -2})
	duration := def(&[nBaseUnits]int8{second: 1})
	electricCharge := def(&[nBaseUnits]int8{ampere: 1, second: 1})
	electricCurrent := def(&[nBaseUnits]int8{ampere: 1})
	electricResistance := def(&[nBaseUnits]int8{kilogram: 1, meter: 2, ampere: -2, second: -3})
	energy := def(&[nBaseUnits]int8{kilogram: 1, meter: 2, second: -2})
	force := def(&[nBaseUnits]int8{kilogram: 1, meter: 1, second: -2})
	frequency := def(&[nBaseUnits]int8{second: -1})
	fuelEfficiency := def(&[nBaseUnits]int8{meter: 2})
	illuminance := def(&[nBaseUnits]int8{candela: 1, steradian: 1, meter: -2})
	information := def(&[nBaseUnits]int8{byte: 1})
	length := def(&[nBaseUnits]int8{meter: 1})
	luminousFlux := def(&[nBaseUnits]int8{candela: 1, steradian: 1})
	luminousIntensity := def(&[nBaseUnits]int8{candela: 1})
	mass := def(&[nBaseUnits]int8{kilogram: 1})
	matter := def(&[nBaseUnits]int8{mole: 1})
	money := def(&[nBaseUnits]int8{currency: 1})
	power := def(&[nBaseUnits]int8{kilogram: 1, meter: 2, second: -3})
	pressure := def(&[nBaseUnits]int8{kilogram: 1, meter: -1, second: -2})
	solidAngle := def(&[nBaseUnits]int8{steradian: 1})
	speed := def(&[nBaseUnits]int8{meter: 1, second: -1})
	temperature := def(&[nBaseUnits]int8{kelvin: 1})
	unitless := def(&[nBaseUnits]int8{})
	voltage := def(&[nBaseUnits]int8{meter: 2, kilogram: 1, second: -3, ampere: -1})
	volume := def(&[nBaseUnits]int8{meter: 3})

	return []*Unit{
		// define only basic unit symbols here, no derived symbols like m/s2, lb/cu ft

		unitless("", 1),

		acceleration("G", 9.80665), //Earth's gravity constant

		angle("rad", 1),           // radians
		angle("deg", math.Pi/180), // degrees (360deg per full circle)
		angle("cycles", math.Pi*2),

		angularVelocity("rpm", math.Pi*2/60), // rounds per minute

		area("sqm", 1),  // square meter, alt unit
		area("ha", 1e4), // hectare
		area("acre", 4046.8564224),
		area("sq mi", 2589988.110336), // square mile
		area("sq in", 0.00064516),     // square inch
		area("sq ft", 0.09290304),     // square feet

		capacitance("F", 1), // farad

		duration("s", 1),
		duration("min", 60),
		duration("h", 3600),
		duration("d", 24*3600),

		electricCharge("C", 1),

		electricCurrent("A", 1),

		electricResistance("Ω", 1),

		energy("J", 1), // joule
		energy("kWh", 3.6e6),

		force("N", 1),                 // newton
		force("lbf", 4.4482216152605), // pound force

		frequency("Hz", 1), // hertz

		fuelEfficiency("L/100km", 1e-8), // Liter per 100km = 1e-3 m3 / 1e5 m = 1e-8 m2

		illuminance("lx", 1),

		information("bit", 0.125),
		information("byte", 1),
		information("KiB", 1024),    // note: KB is 1000
		information("MiB", 1048576), // note: MB is 1e6
		information("GiB", 1073741824),
		information("TiB", 1099511627776),
		information("PiB", 1125899906842624),

		length("m", 1), // meter, metre
		length("mi", 1609.344), // mile
		length("in", 0.0254),   // inch
		length("ft", 0.3048),   // foot
		length("yd", 0.9144),   // yard
		length("M", 1852),      // nautical mile

		luminousFlux("lm", 1),      // lumen
		luminousIntensity("cd", 1), // candela

		mass("kg", 1),              // kilogram
		mass("g", 0.001),           // gram
		mass("t", 1000),            // tonne, metric ton
		mass("lb", 0.45359237),     // pound
		mass("lbs", 0.45359237),     // pound
		mass("oz", 0.028349523125), // ounce avdp
		mass("short ton", 907.18474),
		mass("long ton", 1016.04691),
		mass("st", 6.35029318), // stone

		matter("mol", 1),

		money("¤", 1),      // generic currency symbol
		money("$", 1),      // dollar
		money("USD", 1),    // US dollar
		money("NZD", 1.57), // todo: use conversion table updated by function

		power("W", 1), // watts
		power("hp", 745.699872), // horsepower

		pressure("Pa", 1),           // pascal
		pressure("psi", 6894.75729), // pounds per square inch
		pressure("bar", 1e5),
		pressure("mbar", 100),  // millibar, bar is not SI unit cannot use just any prefix
		pressure("kbar", 1e8), // kilobar
		pressure("mmHg", 133.322387415), // millimeter mercury
		pressure("cmHg", 1333.22387415), // centimeter mercury

		solidAngle("sr", 1), // steradian

		speed("kph", 1000.0/3600.0),   // kilometer per hour, alt unit
		speed("mph", 1609.344/3600.0), // mile per hour
		speed("kn", 1852/3600.0),      // knots

		temperature("K", 1), // kelvin
		temperature("degC", 1), // degree celsius, relative temperature
		temperature("degF", 5.0/9), // degree fahrenheit, relative temperature

		voltage("V", 1), // volt

		volume("cu ft", 35.3146665722),           // cubic foot
		volume("L", 1e-3),                        // liter
		volume("us gal", 0.003785411784),         // US gallon
		volume("imp gal", 0.00454609188),         // Imperial gallon
		volume("us fl oz", 0.0000295735295625),   // US fluid ounce
		volume("imp fl oz", 0.00002841307424375), // Imperial fluid ounce
	}
}
