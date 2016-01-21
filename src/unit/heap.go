package unit

// Heap is similar to an account, but can handle other values than money.
// For example use for inventory, limited resources. A heap has a min
// and max value and guarantees the balance is between these two at all times.
// Initially a Heap has a balance equal to the min value.
type Heap interface {
	Set(Measurement) bool
	Deposit(Measurement) bool
	Withdraw(Measurement) bool
	Balance() Measurement
	Min(Measurement) bool
	Max(Measurement) bool
	Limits() (min Measurement, max Measurement)
}

type heap struct {
	min, max, balance Measurement
}

func NewHeap(min Measurement, max Measurement) Heap {
	if SameUnit(min, max) && Less(min, max) {
		return &heap{min, max, min}
	}
	return nil
}

func (h *heap) Set(m Measurement) bool {
	if !SameUnit(h.balance, m) || h.outOfBounds(m) {
		return false
	}
	h.balance = m
	return true
}

func (h *heap) Deposit(m Measurement) bool {
	if !SameUnit(h.balance, m) {
		return false
	}
	n := Add(h.balance, m)
	if h.outOfBounds(n) {
		return false
	}
	h.balance = n
	return true
}

func (h *heap) Withdraw(m Measurement) bool {
	if !SameUnit(h.balance, m) {
		return false
	}
	n := Subtract(h.balance, m)
	if h.outOfBounds(n) {
		return false
	}
	h.balance = n
	return true
}

func (h *heap) outOfBounds(m Measurement) bool {
	return Less(m, h.min) || More(m, h.max)
}

func (h *heap) Balance() Measurement {
	return h.balance
}

func (h *heap) Min(min Measurement) bool {
	if More(min, h.max) || More(min, h.balance) {
		return false
	}
	h.min = min
	return true
}

func (h *heap) Max(max Measurement) bool {
	if Less(max, h.min) || Less(max, h.balance) {
		return false
	}
	h.max = max
	return true
}

func (h *heap) Limits() (min Measurement, max Measurement) {
	min, max = h.min, h.max
	return
}
