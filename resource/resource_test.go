package resource

import (
	"testing"
	. "github.com/zn8nz/units/quantity"
	. "github.com/zn8nz/units/context"
)

func TestNewHeap(t *testing.T) {
	rsc := New(Q(1, "kg"), Q(100, "kg"), "")
	if rsc == nil {
		t.Error("failed heap creation")
	}
}

func TestDeposit(t *testing.T) {
	rsc := New(Q(1, "kg"), Q(100, "kg"), "")
	ok := rsc.Deposit(Q(50, "m2"))
	if ok {
		t.Error("incompatibility ignored")
	}
	if ok = rsc.Deposit(Q(-20, "kg")); ok {
		t.Error("min ignored")
	}
	if ok = rsc.Deposit(Q(200, "kg")); ok {
		t.Error("max ignored")
	}
	if !Equal(rsc.Balance(), Q(1, "kg"), Q(1, "g")) {
		t.Error("not equal", rsc.Balance(), "1 kg")
	}
	rsc.Deposit(Q(150, "g"))
	if !Equal(rsc.Balance(), Q(1150, "g"), Q(1, "g")) {
		t.Error("balance wrong", rsc.Balance())
	}
}

func TestWithdraw(t *testing.T) {
	rsc := New(Q(-1, "kWh"), Q(100, "kWh"), "")
	if rsc.Set(Q(150, "kWh")) {
		t.Error("ignored out of bounds")
	}
	rsc.Set(Q(50, "kWh"))
	rsc.Withdraw(Q(1e6, "J"))
	if !Equal(rsc.Balance(), Q(1.79e8, "m2.kg/s2"), Q(1, "J")) {
		t.Error("balance wrong", rsc.Balance())
	}
	if rsc.Withdraw(Q(90, "kWh")) {
		t.Error("min ignored")
	}
}

func TestMinMax(t *testing.T) {
	rsc := New(Q(0, "m"), Q(100, "m"), "")
	rsc.Set(Q(30, "m"))
	if rsc.Min(Q(31, "m")) {
		t.Error("min more than balance accepted")
	}
	if rsc.Max(Q(29, "m")) {
		t.Error("max less than balance accepted")
	}
	min, max := rsc.Limits()
	if !Equal(min, Q(0, "m"), Q(1, "cm")) {
		t.Error("lower bound wrong")
	}
	if !Equal(max, Q(100, "m"), Q(1, "cm")) {
		t.Error("upper bound wrong")
	}
}

func TestWithdrawPctContext(t *testing.T) {
	DefineContext("tank", "L", "%.1[1]fℓ")
	rsc := New(Q(1, "L"), Q(50, "L"), "tank")
	rsc.Set(Q(20, "L"))
	m, err := rsc.WithdrawPct(25)
	if err != nil || m.String() != "5.0000 L" {
		t.Error("expected: 5.0000 L, actual:", m.String())
	}
	s := rsc.Balance().String()
	if s != "15.0000 L" {
		t.Error("balance not reduced to 15.0000 L:", s)
	}
	m, err = rsc.WithdrawPct(200)
	if err == nil {
		t.Error("should not be allowed to withdraw 200 percent")
	}
	s = rsc.String()
	if s != "15.0ℓ" {
		t.Error("value withdrawn despite being invalid")
	}
}
