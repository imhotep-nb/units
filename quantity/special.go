package quantity

import (
	"errors"
)

// -- temperature ------------------------------

const abszero = 273.15

// KtoC converts Kelvin to Celsius
func KtoC(q Quantity) (float64, error) {
	if !q.HasCompatibleUnit("K") {
		return 0, errors.New("not a temperature:" + q.String())
	}
	return q.value - abszero, nil
}

// KtoF converts Kelvin to Fahrenheit
func KtoF(q Quantity) (float64, error) {
	if !q.HasCompatibleUnit("K") {
		return 0, errors.New("not a temperature:" + q.String())
	}
	return (q.value-abszero)*1.8 + 32, nil
}

// CtoF converts Celsius to Fahrenheit
func CtoF(c float64) float64 {
	return c*1.8 + 32
}

// FtoC converts Fahrenheit to Celsius
func FtoC(f float64) float64 {
	return (f - 32) / 1.8
}

// CtoK converts Celsius to Kelvin
func CtoK(c float64) Quantity {
	return Q(c+abszero, "K")
}

// FtoK converts Fahrenheit to Kelvin
func FtoK(f float64) Quantity {
	return Q((f-32)/1.8+abszero, "K")
}
