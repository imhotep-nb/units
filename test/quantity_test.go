package t

import (
	"fmt"
	"os"
	"sort"
	"testing"
	"time"
	"unit"
)

func TestPanic(t *testing.T) {
	enablePanic := os.Getenv("GOUNITSPANIC") == "1"
	if enablePanic {
		fmt.Println("Panic if working with incompatible units")
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("TestPanic OK")
			}
		}()
		unit.Add(unit.Q(10, "kph"), unit.Q(20, "V"))
		t.Error("TestPanic didn't work as expected")
	}
}

func TestInvalid(t *testing.T) {
	defer func() {
		recover()
	}()
	m := unit.Q(0, "bla")
	t.Error(m.Inspect())
}

func TestIn(t *testing.T) {
	data := []struct {
		val  float64
		sym  string
		val1 string
		sym1 string
		fail bool
	}{
		{454.8, "kph", "245.5724", "kn", false},
		{454.8, "kph", "-1", "kn", true},
		{1500, "m", "0.9321", "mi", false},
		{0.9320568, "mi", "1500.0000", "m", false},
		{1, "m/s", "3.6000", "kph", false},
		{1, "m/s", "1", "m", true},
		{-1, "m/s", "-1.0000", "m/s", false},
		{34, "¤/m", "51.00", "$", true},
		{1000, "$", "1000.0000", "USD", false},
		{3.1, "us gal", "11.7348", "L", false},
		{7, "L/100km", "0.0700", "mm2", false},
		{3, "N", "3.0000", "kg.m/s2", false},
		{1, "psi", "0.0689", "bar", false},
		{6894.757, "Pa", "1.0000", "lbf.in-2", false},
	}
	for _, d := range data {
		m1 := unit.Q(d.val, d.sym)
		if m1.Invalid() {
			if !d.fail {
				t.Error("source unit not found:", d.sym)
			}
			continue
		}
		if m2, ok := m1.ConvertTo(d.sym1); ok {
			v, s := m2.Split()
			vs := fmt.Sprintf("%.4f", v)
			mismatch := vs != d.val1 || s != d.sym1
			if mismatch && !d.fail || !mismatch && d.fail {
				if d.fail {
					t.Error("expected to fail:", d.val, d.sym, "->", d.val1, d.sym1)
				} else {
					t.Error("expected:", d.val1, d.sym1, "; actual:", vs, s)
				}
			}
		} else {
			if !d.fail {
				t.Error("not expected to fail:", d.val, d.sym, "->", d.sym1)
			}
		}

	}
}

func TestString(t *testing.T) {
	data := []struct {
		input    unit.Quantity
		expected string
	}{
		{unit.Q(12.3456, "kn"), "12.3456 kn"},
		{unit.Q(0, "kn"), "0.0000 kn"},
		{unit.Q(-14.581699, "mph"), "-14.5817 mph"},
		{unit.Q(0.00001, "m"), "0.0000 m"},
	}
	for _, d := range data {
		s := d.input.String()
		if s != d.expected {
			t.Error("expected:", d.expected, "actual:", s)
		}
	}
	unit.DefaultFormat = "%.0f%s"
	if unit.Q(500.9999, "mph").String() != "501mph" {
		t.Error("setting default format failed")
	}
	unit.DefaultFormat = "%.4f %s"
	a := unit.Q(123.5, "NZD")
	if a.String() != "123.5000 NZD" {
		t.Error("currency formatting failed", a)
	}
}

func TestCalc1(t *testing.T) {
	q := unit.Q
	data := []struct {
		op       string
		x, y     unit.Quantity
		expected string
	}{
		{"+", q(10, "m"), q(8, "m"), "18.0000 m"},
		{"+", q(15, "km"), q(2, "mi"), "18218.6880 m"},
		{"-", q(5.301, "kg"), q(302, "g"), "4.9990 kg"},
		{"-", q(1.4, "mph"), q(3.0, "kn"), "-0.9175 m.s-1"},
		{"*", q(2, "kg"), q(15, "m"), "30.0000 m.kg"},
		{"/", q(9, "km"), q(2, "h"), "1.2500 m.s-1"},
		{"1/", q(100, "m/s"), unit.Quantity{}, "0.0100 m-1.s"},
		{"1/", q(8.0, "m"), unit.Quantity{}, "0.1250 m-1"},
	}
	for _, d := range data {
		var result unit.Quantity
		switch d.op {
		case "+":
			result = unit.Add(d.x, d.y)
		case "-":
			result = unit.Subtract(d.x, d.y)
		case "*":
			result = unit.Mult(d.x, d.y)
		case "/":
			result = unit.Div(d.x, d.y)
		case "1/":
			result = unit.Reciprocal(d.x)
		}
		if result.String() != d.expected {
			t.Error("expected:", d.expected, "actual:", result)
		}
	}
}

func TestCalc2(t *testing.T) {
	q := unit.Q
	data := []struct {
		op       string
		q        unit.Quantity
		f        float64
		expected string
	}{
		{"*", q(100, "m/s"), 1.2, "120.0000 m/s"},
		{"/", q(100, "g"), 4.0, "25.0000 g"},
		{"^", q(2.0, "m"), 3, "8.0000 m3"},
		{"^", q(8.4, "m"), -3, "0.0017 m-3"},
	}
	for _, d := range data {
		var result unit.Quantity
		switch d.op {
		case "*":
			result = unit.MultFac(d.q, d.f)
		case "/":
			result = unit.DivFac(d.q, d.f)
		case "^":
			result = unit.Power(d.q, int8(d.f))
		}
		if result.String() != d.expected {
			t.Error("expected:", d.expected, "actual:", result)
		}
	}
}

func TestCalc3(t *testing.T) {
	result := unit.Sum(unit.Q(5.1, "Pa"), unit.Q(0.3, "N.m-2"), unit.Q(0.11, "m-2.N"))
	expected := "5.5100 m-1.kg.s-2"
	if result.String() != expected {
		t.Error("expected:", expected, "actual:", result.String())
	}
	result = unit.Diff(unit.Q(100, "kph"), unit.Q(7, "mph"), unit.Q(1, "kn"))
	expected = "24.1341 m.s-1"
	if result.String() != expected {
		t.Error("expected:", expected, "actual:", result.String())
	}
}

func TestMixedUnits(t *testing.T) {
	p1 := unit.Q(7, "N.m-2")
	p2 := unit.Q(8, "Pa")
	if unit.AreCompatible(p1, p2) {
		p3 := unit.Add(p1, p2)
		const result = "15.0000 m-1.kg.s-2"
		if p3.String() != result {
			t.Error("expected:", result, "actual:", p3)
		}
	} else {
		t.Error("not same unit: ", p1.Symbol(), p2.Symbol())
	}
}

func TestPer(t *testing.T) {
	p1 := unit.Q(1, "km/h")
	p2 := unit.Q(2, "kph")
	p3 := unit.Q(3, "m/s")
	if !unit.AreCompatible(p1, p2) {
		t.Error("incompatible:", p1, "<>", p2)
	}
	if !unit.AreCompatible(p2, p3) {
		t.Error("incompatible:", p2, "<>", p3)
	}
	p4 := unit.Q(4, "kg.m/s2")
	p5 := unit.Q(5, "N")
	if !unit.AreCompatible(p4, p5) {
		t.Error("incompatible:", p4, "<>", p5)
	}
	p6 := unit.Q(6, "W")
	p7 := unit.Q(7, "J/s")
	if !unit.AreCompatible(p6, p7) {
		t.Error("same unit:", p6, p7)
	}
	p8 := unit.Subtract(unit.Q(8.8, "N.m/s"), unit.Q(8.8, "W"))
	if p8.String() != "0.0000 m2.kg.s-3" {
		t.Error()
	}
}

func TestEqual(t *testing.T) {
	p1 := unit.Q(999, "m")
	p2 := unit.Q(1, "km")
	if !unit.Equal(p1, p2, unit.Q(2, "m")) {
		t.Error("not equal: ", p1, p2)
	}
	if unit.Equal(p1, p2, unit.Q(1, "m")) {
		t.Error("false equality:", p1, p2)
	}
}

func TestNormalize(t *testing.T) {
	p1 := unit.Q(1.2, "mph")
	if p1.Value() != 1.2 || p1.Symbol() != "mph" {
		t.Error("unit initialization error", p1)
	}
	p1.Normalize()
	if fmt.Sprintf("%.4f", p1.Value()) != "0.5364" || p1.Symbol() != "m.s-1" {
		t.Error("unit initialization error", p1)
	}
}

func TestParse(t *testing.T) {
	p1 := unit.Q(12.4, "km.s-2")
	p2, err := unit.Parse("12.4 km/s2")
	if err != nil {
		t.Error(err)
	} else if !unit.Equal(p1, p2, unit.Q(0.01, "m.s-2")) {
		t.Error("not equal", p1, "<>", p2)
	}
	p3 := unit.Q(3894829.88, "sq in")
	p4, err := unit.Parse("  3,894,829.88sq in   ")
	if err != nil {
		t.Error(err)
	} else if !unit.Equal(p3, p4, unit.Q(0.001, "sq in")) {
		t.Error("not equal", p3, "<>", p4)
	}
}

func TestParse2(t *testing.T) {
	data := []struct {
		s    string
		fail bool
	}{
		{"38J", false},
		{"  -15.5  K  ", false},
		{"1,000 kW/sr", false},
		{"foo", true},
		{"/12309.8m", true},
		{"12,058,884.881 N/m2", false},
		{"5 chickens/m2", true},
		{"1.1 sq in", false},
		{"5.5.6 m", true},
	}
	for _, d := range data {
		_, err := unit.Parse(d.s)
		if err != nil && !d.fail {
			t.Error("failed but shouldn't: [", d.s, "]")
		} else if err == nil && d.fail {
			t.Error("should fail but didn't: [", d.s, "]")
		}
	}
}

func TestSort(t *testing.T) {
	arr := unit.Quantities{
		unit.Q(0.2, "M"),
		unit.Q(-3, "m"),
		unit.Q(-1.5, "m"),
		unit.Q(0.1, "cm"),
		unit.Q(0.1, "mm"),
		unit.Q(4, "ft"),
	}
	sort.Sort(arr)
	sa := fmt.Sprintf("%v", arr)
	if sa != "[-3.0000 m -1.5000 m 0.1000 mm 0.1000 cm 4.0000 ft 0.2000 M]" {
		t.Error("sort error", sa)
	}
}

func TestDuration(t *testing.T) {
	var t1 unit.Quantity = unit.Q(1.5, "d")
	var t2 time.Duration
	t2, err := unit.Duration(t1)
	if err != nil {
		t.Error(err)
	}
	if t2.Hours() != 36 {
		t.Error("expected:", 36, "actual:", t2.Hours())
	}
}

func TestPrefix(t *testing.T) {
	m1 := unit.Q(25*unit.Centi, "m")
	m2 := unit.Q(25, "cm")
	if !unit.Equal(m1, m2, unit.Q(1e-6, "m")) {
		t.Error("not equal:", m1, m2)
	}
	m3 := unit.Q(7*unit.Cubic(unit.Deci), "m3")
	m4 := unit.Q(7, "L")
	if !unit.AreCompatible(m3, m4) || !unit.Equal(m3, m4, unit.Q(1e-6, "m")) {
		t.Error("not equal:", m3, m4)
	}
}
