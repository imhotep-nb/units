package resource

import (
	"errors"
	"fmt"
	"unit"
)

// Resource is similar to an account, but can handle other values than money.
// For example use for inventory, limited resources. A Resource has a min
// and max value and guarantees the balance is between these two at all times.
// Initially a Resource has a balance equal to the min value.
type Resource struct {
	min, max, balance unit.Quantity
	*unit.Context
}

// New creates a new Resource with the given minimum and maximum values.
// min should be less than max and the units should be compatible.
// The initial balance value is set to min. A Context name can be provided, or ""
// if no Context is required.
func New(min unit.Quantity, max unit.Quantity, context string) *Resource {
	var ctx *unit.Context
	if context != "" {
		ctx = unit.Ctx(context)
	} else {
		ctx, _ = unit.DefineContext("", min.Symbol(), unit.DefaultFormat)
	}
	if unit.AreCompatible(min, max) && unit.Less(min, max) {
		return &Resource{ctx.Convert(min), ctx.Convert(max), min, ctx}
	}
	return nil
}

// Set the Resource to the given value. The value should be between the min
// and max of the Resource. Return true for success, false for incompatible unit
// or out of bounds.
func (h *Resource) Set(q unit.Quantity) bool {
	if !unit.AreCompatible(h.balance, q) || h.outOfBounds(q) {
		return false
	}
	h.balance = q
	return true
}

// Deposit adds the Measurement to the Resource. Return true for success, false for
// incompatible unit or out of bounds.
func (h *Resource) Deposit(q unit.Quantity) bool {
	if !unit.AreCompatible(h.balance, q) {
		return false
	}
	n := unit.Add(h.balance, q)
	if h.outOfBounds(n) {
		return false
	}
	h.balance = n
	return true
}

// Withdraw subtracts the given amount from the Resource.
// Return true for success, false for incompatible unit or out of bounds
func (h *Resource) Withdraw(q unit.Quantity) bool {
	if !unit.AreCompatible(h.balance, q) {
		return false
	}
	n := unit.Subtract(h.balance, q)
	if h.outOfBounds(n) {
		return false
	}
	h.balance = n
	return true
}

// WithdrawPct subtracts a percentage of the balance. It returns the
// quantity that has been deducted and an error or nil if the percentage
// is not in the range 0..100.
func (h *Resource) WithdrawPct(percentage float64) (unit.Quantity, error) {
	if percentage < 0 || percentage > 100 {
		msg := fmt.Sprintf("percentage not in range 0..1")
		return unit.Quantity{}, errors.New(msg)
	}
	taken := unit.MultFac(h.balance, percentage/100.0)
	h.balance = unit.Subtract(h.balance, taken)
	return h.Convert(taken), nil
}

func (h *Resource) outOfBounds(q unit.Quantity) bool {
	return unit.Less(q, h.min) || unit.More(q, h.max)
}

// Balance returns the current balance.
func (h *Resource) Balance() unit.Quantity {
	return h.Convert(h.balance)
}

// Min sets a new minimum balance. Returns true for success, false for incompatible unit
// or in case the min is more than the current balance.
func (h *Resource) Min(min unit.Quantity) bool {
	if !unit.AreCompatible(h.balance, min) || unit.More(min, h.balance) {
		return false
	}
	h.min = min
	return true
}

// Min sets a new minimum balance. Returns true for success, false for incompatible unit
// or in case the max is less than the current balance.
func (h *Resource) Max(max unit.Quantity) bool {
	if !unit.AreCompatible(h.balance, max) || unit.Less(max, h.balance) {
		return false
	}
	h.max = max
	return true
}

// Limits returns the min and max Measurements of the resource.
func (h *Resource) Limits() (min unit.Quantity, max unit.Quantity) {
	min, max = h.min, h.max
	return
}

// String returns a string value formatted according to the Context.
func (h Resource) String() string {
	return h.Context.String(h.balance)
}
