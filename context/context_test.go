package context

import (
	"bytes"
	"testing"
	. "github.com/zn8nz/units/quantity"
)

const (
	personHeight  = "person height"
	landArea      = "land area"
	money         = "money"
	rainIntensity = "rain intensity"
	pressureDrop  = "pressure drop"
)

func init() {
	DefineContext(personHeight, "cm", "%.0[1]fcm")
	DefineContext(landArea, "acre", "%0.[1]f acres")
	DefineContext(money, "¤", "%[2]s%.2[1]f")
	DefineContext(rainIntensity, "mm/h", "%.1f %s")
	DefineContext(pressureDrop, "hPa/km", "%.0f %s")
}

func TestContextDefinition(t *testing.T) {
	c := Ctx(personHeight)
	if c == nil {
		t.Error("not found: person height")
	}
	c = Ctx("foo")
	if c != nil {
		t.Error("nil expected, actual:", c.Name)
	}
	DefineContext("letter weight", "g", "%f %s")
	c = Ctx("letter weight")
	if c == nil {
		t.Error("not found: letter weight")
	}
	DeleteContext(Ctx("letter weight"))
	c = Ctx("letter weight")
	if c != nil {
		t.Error("nil expected, found: ", c.Name)
	}
	c = Ctx(rainIntensity)
	if c == nil || c.Name != rainIntensity || c.Symbol() != "mm/h" {
		t.Errorf("unexpected context: %v", c)
	}
	c = Ctx(money)
	s := c.String(Q(250.199, "$"))
	if s != "¤250.20" {
		t.Error("expected ¤250.20, actual:", s)
	}
}

func TestContextConversion(t *testing.T) {
	height := Ctx(personHeight)
	q := height.Q(1.75, "m")
	s := height.String(q)
	if s != "175cm" {
		t.Error("expected 175cm, actual:", s)
	}
	q = Add(Q(5, "ft"), Q(11, "in"))
	s = height.String(q)
	if s != "180cm" {
		t.Error("expected 180cm, actual:", s)
	}

	rain := Ctx(rainIntensity)
	q = rain.Q(1, "in/d") // inch/day
	s = rain.String(q)
	if s != "1.1 mm/h" {
		t.Error("(String) expected 1.1 mm/h, actual:", s)
	}
	var b bytes.Buffer
	rain.Format(&b, q)
	s = string(b.Bytes())
	if s != "1.1 mm/h" {
		t.Error("(Format) expected 1.1 mm/h, actual:", s)
	}
}

func TestUnregisteredContext(t *testing.T) {
	pressureChange, err := DefineContext("", "Pa/min", "%.0f %s")
	if err != nil {
		t.Error(err)
	}
	q := Q(3, "bar/h")
	s := pressureChange.String(q)
	if s != "5000 Pa/min" {
		t.Error("expected: 5000 Pa/min, actual:", s)
	}
	ctx := Ctx("")
	if ctx != nil {
		t.Error("should be nil:", ctx)
	}

	ctx = Ctx(pressureDrop)
	q = ctx.Q(11, "mbar/hm")
	s = ctx.String(q)
	if s != "110 hPa/km" {
		t.Error("expected 110 hPa/km, actual:", s)
	}
}
