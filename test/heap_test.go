package t

import (
	//"fmt"
	"testing"
	"unit"
)

func TestNewHeap(t *testing.T) {
	h := unit.NewHeap(unit.M(1, "kg"), unit.M(100, "kg"))
	if h == nil {
		t.Error("failed heap creation")
	}
}

func TestDeposit(t *testing.T) {
	h := unit.NewHeap(unit.M(1, "kg"), unit.M(100, "kg"))
	ok := h.Deposit(unit.M(50, "m2"))
	if ok {
		t.Error("incompatibility ignored")
	}
	if ok = h.Deposit(unit.M(-20, "kg")); ok {
		t.Error("min ignored")
	}
	if ok = h.Deposit(unit.M(200, "kg")); ok {
		t.Error("max ignored")
	}
	if !unit.Equal(h.Balance(), unit.M(1, "kg"), unit.M(1, "g")) {
		t.Error("not equal", h.Balance(), "1 kg")
	}
	h.Deposit(unit.M(150, "g"))
	if !unit.Equal(h.Balance(), unit.M(1150, "g"), unit.M(1, "g")) {
		t.Error("balance wrong", h.Balance())
	}
}

func TestWithdraw(t *testing.T) {
	h := unit.NewHeap(unit.M(-1, "kWh"), unit.M(100, "kWh"))
	if h.Set(unit.M(150, "kWh")) {
		t.Error("ignored out of bounds")
	}
	h.Set(unit.M(50, "kWh"))
	h.Withdraw(unit.M(1e6, "J"))
	if !unit.Equal(h.Balance(), unit.M(1.79e8, "m2.kg/s2"), unit.M(1, "J")) {
		t.Error("balance wrong", h.Balance())
	}
	if h.Withdraw(unit.M(90, "kWh")) {
		t.Error("min ignored")
	}
}

func TestMinMax(t *testing.T) {
	h := unit.NewHeap(unit.M(0, "m"), unit.M(100, "m"))
	h.Set(unit.M(30, "m"))
	if h.Min(unit.M(31, "m")) {
		t.Error("min more than balance accepted")
	}
	if h.Max(unit.M(29, "m")) {
		t.Error("max less than balance accepted")
	}
	min, max := h.Limits()
	if !unit.Equal(min, unit.M(0, "m"), unit.M(1, "cm")) {
		t.Error("lower bound wrong")
	}
	if !unit.Equal(max, unit.M(100, "m"), unit.M(1, "cm")) {
		t.Error("upper bound wrong")
	}
}
