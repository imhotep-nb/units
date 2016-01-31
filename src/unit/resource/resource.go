package resource

import "unit"

// Resource is similar to an account, but can handle other values than money.
// For example use for inventory, limited resources. A Resource has a min
// and max value and guarantees the balance is between these two at all times.
// Initially a Resource has a balance equal to the min value.
type Resource interface {
	Set(unit.Measurement) bool
	Deposit(unit.Measurement) bool
	Withdraw(unit.Measurement) bool
	Balance() unit.Measurement
	Min(unit.Measurement) bool
	Max(unit.Measurement) bool
	Limits() (min unit.Measurement, max unit.Measurement)
}

type resource struct {
	min, max, balance unit.Measurement
}

func New(min unit.Measurement, max unit.Measurement) Resource {
	if unit.AreCompatible(min, max) && unit.Less(min, max) {
		return &resource{min, max, min}
	}
	return nil
}

func (h *resource) Set(m unit.Measurement) bool {
	if !unit.AreCompatible(h.balance, m) || h.outOfBounds(m) {
		return false
	}
	h.balance = m
	return true
}

func (h *resource) Deposit(m unit.Measurement) bool {
	if !unit.AreCompatible(h.balance, m) {
		return false
	}
	n := unit.Add(h.balance, m)
	if h.outOfBounds(n) {
		return false
	}
	h.balance = n
	return true
}

func (h *resource) Withdraw(m unit.Measurement) bool {
	if !unit.AreCompatible(h.balance, m) {
		return false
	}
	n := unit.Subtract(h.balance, m)
	if h.outOfBounds(n) {
		return false
	}
	h.balance = n
	return true
}

func (h *resource) outOfBounds(m unit.Measurement) bool {
	return unit.Less(m, h.min) || unit.More(m, h.max)
}

func (h *resource) Balance() unit.Measurement {
	return h.balance
}

func (h *resource) Min(min unit.Measurement) bool {
	if unit.More(min, h.max) || unit.More(min, h.balance) {
		return false
	}
	h.min = min
	return true
}

func (h *resource) Max(max unit.Measurement) bool {
	if unit.Less(max, h.min) || unit.Less(max, h.balance) {
		return false
	}
	h.max = max
	return true
}

func (h *resource) Limits() (min unit.Measurement, max unit.Measurement) {
	min, max = h.min, h.max
	return
}
