package t

import (
	"fmt"
	"testing"
	"unit"
)

const (
	personHeight = "person height"
	landArea = "land area"
	money = "money"
	rainIntensity = "rain intensity"
)

func init() {
	unit.DefineContext(personHeight, "cm", "%.0[1]fcm")
	unit.DefineContext(landArea, "acre", "%0.[1]f acres")
	unit.DefineContext(money, "¤", "%[2]s%[1].2f")
	unit.DefineContext(rainIntensity, "mm/h", "%.1f %s")
	fmt.Println("init")
}

func TestContextDefinition(t *testing.T) {
	c := unit.Ctx(personHeight)
	if c == nil {
		t.Error("not found: person height")
	}
	c = unit.Ctx("foo")
	if c != nil {
		t.Error("nil expected, actual:", c.Name)
	}
	unit.DefineContext("letter weight", "g", "%f %s")
	c = unit.Ctx("letter weight")
	if c == nil {
		t.Error("not found: letter weight")
	}
	unit.DeleteContext(unit.Ctx("letter weight"))
	c = unit.Ctx("letter weight")
	if c != nil {
		t.Error("nil expected, found: ", c.Name)
	}
	c = unit.Ctx(rainIntensity)
	if c == nil || c.Name != rainIntensity || c.Symbol() != "mm/h" {
		t.Errorf("unexpected context: %v", c)
	}
}


func TestContextConversion(t *testing.T) {
	height := unit.Ctx(personHeight)
	m := height.M(1.75, "m")
	s := height.String(m)
	if s != "175cm" {
		t.Error("expected 175cm, actual:", s)
	}
	m = unit.Add(unit.M(5, "ft"), unit.M(11, "in"))
	s = height.String(m)
	if s != "180cm" {
		t.Error("expected 180cm, actual:", s)
	}
	
	rain := unit.Ctx(rainIntensity)
	m = rain.M(1, "in/d")  // inch/day
	s = rain.String(m)
	if s != "1.1 mm/h" {
		t.Error("expected 1.1 mm/h, actual:", s)
	}
	// todo .Format
}
