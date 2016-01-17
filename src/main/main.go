package main

import (
	"fmt"
	"unit"
)

func main() {
	m, err := unit.ParseSymbol("N/m2")
	m1 := unit.M(33, "N/m2")
	fmt.Println("m=", m, "err=", err, "m1=", m1)
	if m, ok := unit.M(6894.757, "N/m2").ConvertTo("lbf.in-2"); ok {
		fmt.Println(m)
	} else {
		fmt.Println("error")
	}
	fmt.Println("----")
	siFactor, err := unit.Define("foo", 7, "lbf/sq in")
	fmt.Println("siFactor", siFactor)
	r, _ := unit.M(3, "foo").ConvertTo("bar")
	fmt.Println(r)
	r = r.ToSI()
	fmt.Println(r)
	fmt.Println("----")
	t := unit.M(12345, "psi")
	fmt.Println(t.String())
	fmt.Println(t.Inspect())
	fmt.Println(t.Symbol())
	t1, t2 := t.Split()
	fmt.Println(t1, "*", t2)
	fmt.Println(t.In("bar"))
	value, ok := t.ConvertTo("bar")
	fmt.Println(value, ok)
	fmt.Println(t.ToSI())
	t.Normalize()
	fmt.Println(t)
	
}