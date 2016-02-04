package unit

import (
	"errors"
)


// -- temperature ------------------------------

const abszero = 273.15

func KtoC(q Quantity) (float64, error) {
	if !q.HasCompatibleUnit("K") {
		return 0, errors.New("not a temperature:" + q.String())
	}
	return q.value - abszero, nil
}

func KtoF(q Quantity) (float64, error) {
	if !q.HasCompatibleUnit("K") {
		return 0, errors.New("not a temperature:" + q.String())
	}
	return (q.value - abszero) * 1.8 + 32, nil
}

func CtoF(c float64) float64 {
	return c * 1.8 + 32
}

func FtoC(f float64) float64 {
	return (f - 32) / 1.8
}

func CtoK(c float64) Quantity {
	return Q(c + abszero, "K")
}

func FtoK(f float64) Quantity {
	return Q((f - 32) / 1.8 + abszero, "K")
}


