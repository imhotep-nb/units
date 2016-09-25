package resource

import (
	"errors"
	"fmt"
	us "github.com/zn8nz/units/quantity"
	"github.com/zn8nz/units/context"
)

// Resource is similar to an account, but can handle other values than money.
// For example use for inventory, limited resources. A Resource has a min
// and max value and guarantees the balance is between these two at all times.
// Initially a Resource has a balance equal to the min value.
type Resource struct {
	min, max, balance us.Quantity
	*context.Context
}

// New creates a new Resource with the given minimum and maximum values.
// min should be less than max and the units should be compatible.
// The initial balance value is set to min. A Context name can be provided, or ""
// if no Context is required.
func New(min us.Quantity, max us.Quantity, c string) *Resource {
	var ctx *context.Context
	if c != "" {
		ctx = context.Ctx(c)
	} else {
		ctx, _ = context.DefineContext("", min.Symbol(), us.DefaultFormat)
	}
	if us.AreCompatible(min, max) && us.Less(min, max) {
		return &Resource{ctx.Convert(min), ctx.Convert(max), min, ctx}
	}
	return nil
}

// Set the Resource to the given value. The value should be between the min
// and max of the Resource. Return true for success, false for incompatible unit
// or out of bounds.
func (h *Resource) Set(q us.Quantity) bool {
	if !us.AreCompatible(h.balance, q) || h.outOfBounds(q) {
		return false
	}
	h.balance = q
	return true
}

// Deposit adds the Measurement to the Resource. Return true for success, false for
// incompatible unit or out of bounds.
func (h *Resource) Deposit(q us.Quantity) bool {
	if !us.AreCompatible(h.balance, q) {
		return false
	}
	n := us.Add(h.balance, q)
	if h.outOfBounds(n) {
		return false
	}
	h.balance = n
	return true
}

// Withdraw subtracts the given amount from the Resource.
// Return true for success, false for incompatible unit or out of bounds
func (h *Resource) Withdraw(q us.Quantity) bool {
	if !us.AreCompatible(h.balance, q) {
		return false
	}
	n := us.Subtract(h.balance, q)
	if h.outOfBounds(n) {
		return false
	}
	h.balance = n
	return true
}

// WithdrawPct subtracts a percentage of the balance. It returns the
// quantity that has been deducted and an error or nil if the percentage
// is not in the range 0..100.
func (h *Resource) WithdrawPct(percentage float64) (us.Quantity, error) {
	if percentage < 0 || percentage > 100 {
		msg := fmt.Sprintf("percentage not in range 0..1")
		return us.Quantity{}, errors.New(msg)
	}
	taken := us.MultFac(h.balance, percentage/100.0)
	h.balance = us.Subtract(h.balance, taken)
	return h.Convert(taken), nil
}

func (h *Resource) outOfBounds(q us.Quantity) bool {
	return us.Less(q, h.min) || us.More(q, h.max)
}

// Balance returns the current balance.
func (h *Resource) Balance() us.Quantity {
	return h.Convert(h.balance)
}

// Min sets a new minimum balance. Returns true for success, false for incompatible unit
// or in case the min is more than the current balance.
func (h *Resource) Min(min us.Quantity) bool {
	if !us.AreCompatible(h.balance, min) || us.More(min, h.balance) {
		return false
	}
	h.min = min
	return true
}

// Max sets a new maximum balance. Returns true for success, false for incompatible unit
// or in case the max is less than the current balance.
func (h *Resource) Max(max us.Quantity) bool {
	if !us.AreCompatible(h.balance, max) || us.Less(max, h.balance) {
		return false
	}
	h.max = max
	return true
}

// Limits returns the min and max Measurements of the resource.
func (h *Resource) Limits() (min us.Quantity, max us.Quantity) {
	min, max = h.min, h.max
	return
}

// String returns a string value formatted according to the Context.
func (h Resource) String() string {
	return h.Context.String(h.balance)
}
