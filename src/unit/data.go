package unit

import (
	"math"
)

func setup() []*unit {
	// keep alphabetic order!
	unitless := def(expMap{})
	acceleration := def(expMap{metre: 1, second: -2})
	angle := def(expMap{radian: 1})
	angularVelocity := def(expMap{radian: 1, second: -1})
	area := def(expMap{metre: 2})
	density := def(expMap{kilogram: 1, metre: -3})
	duration := def(expMap{second: 1})
	electricCurrent := def(expMap{ampere: 1})
	electricCharge := def(expMap{ampere: 1, second: 1})
	energy := def(expMap{kilogram: 1, metre: 2, second: -2})
	force := def(expMap{kilogram: 1, metre: 1, second: -2})
	frequency := def(expMap{second: -1})
	fuelEfficiency := def(expMap{metre: 2})
	illuminance := def(expMap{candela: 1, steradian: 1, metre: -2})
	information := def(expMap{byte: 1})
	length := def(expMap{metre: 1})
	luminousFlux := def(expMap{candela: 1, steradian: 1})
	luminousIntensity := def(expMap{candela: 1})
	mass := def(expMap{kilogram: 1})
	matter := def(expMap{mole: 1})
	money := def(expMap{currency: 1})
	power := def(expMap{kilogram: 1, metre: 2, second: -3})
	pressure := def(expMap{kilogram: 1, metre: -1, second: -2})
	solidAngle := def(expMap{steradian: 1})
	speed := def(expMap{metre: 1, second: -1})
	temperature := def(expMap{kelvin: 1})
	voltage := def(expMap{metre: 2, kilogram: 1, second: -3, ampere: -1})
	volume := def(expMap{metre: 3})

	return []*unit{
		unitless("", 1),

		acceleration("m/s2", 1),
		acceleration("G", 9.80665), //Earth's gravity

		angle("rad", 1),
		angle("deg", math.Pi/180),
		angle("cycles", math.Pi*2),

		angularVelocity("rpm", math.Pi*2/60),

		area("sqm", 1),
		area("ha", 1e4),
		area("acre", 4046.8564224),
		area("sq mi", 2589988.110336),
		area("sq in", 0.00064516),
		area("sq ft", 0.09290304),

		density("lb/cu ft", 0.0624279606),

		duration("s", 1),
		duration("min", 60),
		duration("h", 3600),
		duration("d", 24*3600),

		electricCharge("C", 1),

		electricCurrent("A", 1),

		energy("J", 1),
		energy("kWh", 3.6e6),

		force("N", 1),
		force("lbf", 4.4482216152605),

		frequency("Hz", 1),

		fuelEfficiency("m2", 1),
		fuelEfficiency("L/100km", 1e-8),

		illuminance("lx", 1),

		information("bit", 0.125),
		information("byte", 1),
		information("KiB", 1024),
		information("MiB", 1048576),
		information("GiB", 1073741824),
		information("TiB", 1099511627776),
		information("PiB", 1125899906842624),

		length("m", 1),
		length("cm", 0.01),
		length("mm", 0.001),
		length("km", 1000),
		length("mi", 1609.344), // mile
		length("in", 0.0254),
		length("ft", 0.3048),
		length("yd", 0.9144),
		length("M", 1852), // nautical mile

		luminousFlux("lm", 1),      // lumen
		luminousIntensity("cd", 1), // candela

		mass("kg", 1),
		mass("g", 0.001),
		mass("t", 1000),
		mass("lb", 0.45359237),
		mass("short ton", 907.18474),
		mass("long ton", 1016.04691),

		matter("mol", 1),

		money("¤", 1),
		money("$", 1),
		money("USD", 1),
		money("NZD", 1.57),

		power("W", 1), // watt
		power("kW", 1000),
		power("hp", 745.699872), // horsepower

		pressure("Pa", 1),
		pressure("psi", 6894.75729),
		pressure("bar", 1e5),
		pressure("mmHg", 133.322387415),

		solidAngle("sr", 1),

		speed("m/s", 1),
		speed("kph", 1000.0/3600.0),
		speed("mph", 1609.344/3600.0),
		speed("kn", 1852/3600.0),

		temperature("K", 1),

		voltage("V", 1),
		voltage("kV", 1000),

		volume("cu ft", 35.3146665722),
		volume("L", 1e-3),
		volume("us gal", 0.003785411784),
		volume("imp gal", 0.00454609188),
	}
}
