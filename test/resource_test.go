package t

import (
	//"fmt"
	"testing"
	"unit"
	"unit/resource"
)

func TestNewHeap(t *testing.T) {
	rsc := resource.New(unit.M(1, "kg"), unit.M(100, "kg"))
	if rsc == nil {
		t.Error("failed heap creation")
	}
}

func TestDeposit(t *testing.T) {
	rsc := resource.New(unit.M(1, "kg"), unit.M(100, "kg"))
	ok := rsc.Deposit(unit.M(50, "m2"))
	if ok {
		t.Error("incompatibility ignored")
	}
	if ok = rsc.Deposit(unit.M(-20, "kg")); ok {
		t.Error("min ignored")
	}
	if ok = rsc.Deposit(unit.M(200, "kg")); ok {
		t.Error("max ignored")
	}
	if !unit.Equal(rsc.Balance(), unit.M(1, "kg"), unit.M(1, "g")) {
		t.Error("not equal", rsc.Balance(), "1 kg")
	}
	rsc.Deposit(unit.M(150, "g"))
	if !unit.Equal(rsc.Balance(), unit.M(1150, "g"), unit.M(1, "g")) {
		t.Error("balance wrong", rsc.Balance())
	}
}

func TestWithdraw(t *testing.T) {
	rsc := resource.New(unit.M(-1, "kWh"), unit.M(100, "kWh"))
	if rsc.Set(unit.M(150, "kWh")) {
		t.Error("ignored out of bounds")
	}
	rsc.Set(unit.M(50, "kWh"))
	rsc.Withdraw(unit.M(1e6, "J"))
	if !unit.Equal(rsc.Balance(), unit.M(1.79e8, "m2.kg/s2"), unit.M(1, "J")) {
		t.Error("balance wrong", rsc.Balance())
	}
	if rsc.Withdraw(unit.M(90, "kWh")) {
		t.Error("min ignored")
	}
}

func TestMinMax(t *testing.T) {
	rsc := resource.New(unit.M(0, "m"), unit.M(100, "m"))
	rsc.Set(unit.M(30, "m"))
	if rsc.Min(unit.M(31, "m")) {
		t.Error("min more than balance accepted")
	}
	if rsc.Max(unit.M(29, "m")) {
		t.Error("max less than balance accepted")
	}
	min, max := rsc.Limits()
	if !unit.Equal(min, unit.M(0, "m"), unit.M(1, "cm")) {
		t.Error("lower bound wrong")
	}
	if !unit.Equal(max, unit.M(100, "m"), unit.M(1, "cm")) {
		t.Error("upper bound wrong")
	}
}
