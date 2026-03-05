package dom

import (
	"encoding/xml"
	"regexp"
	"strconv"
	"testing"

	"github.com/rickb777/expect"
	"github.com/rickb777/simplexml/ns"
)

func TestAttrsNames(t *testing.T) {
	a1 := Attr("a", "", "1")
	b1 := Attr("b", "", "1")
	xa3 := Attr("a", "x", "3")

	as := Attrs{a1, b1, xa3}

	expect.Slice(as.Names()).ToBe(t, a1.Name, b1.Name, xa3.Name)

	as = nil
	expect.Slice(as.Names()).ToBe(t)
}

func TestAttrsGet(t *testing.T) {
	a1 := Attr("a", "", "1")
	b1 := Attr("b", "", "1")
	a2 := Attr("a", "", "2")
	xa3 := Attr("a", "x", "3")
	xb4 := Attr("b", "x", "4")
	ya5 := Attr("a", "y", "5")

	as := Attrs{a1, b1, a2, xa3, xb4, ya5}

	expect.Value(as.Get(ns.Name("a", "x"))).ToBe(t, xa3)
	expect.Value(as.Get(ns.Local("a"))).ToBe(t, a1)
	expect.Value(as.Get(ns.Local("z"))).ToBe(t, xml.Attr{})

	as = nil
	expect.Value(as.Get(ns.Local("a"))).ToBe(t, xml.Attr{})
}

func TestAttrsWithName(t *testing.T) {
	a1 := Attr("a", "", "1")
	b1 := Attr("b", "", "1")
	a2 := Attr("a", "", "2")
	xa1 := Attr("a", "x", "1")
	xb1 := Attr("b", "x", "1")
	ya1 := Attr("a", "y", "1")

	as := Attrs{a1, b1, a2, xa1, xb1, ya1}

	expect.Slice(as.With(ns.Name("a", "x"))).ToBe(t, xa1)
	expect.Slice(as.With(ns.NameRE(regexp.MustCompile("."), regexp.MustCompile(".")))).ToBe(t, xa1, xb1, ya1)
	expect.Slice(as.With(ns.Local("a"))).ToBe(t, a1, a2, xa1, ya1)
	expect.Slice(as.With(ns.Local("b"))).ToBe(t, b1, xb1)
	expect.Slice(as.With(ns.Space("x"))).ToBe(t, xa1, xb1)
	expect.Slice(as.With(ns.Space("x"), ns.Space("y"))).ToBe(t, xa1, xb1, ya1)
	expect.Slice(as.With(ns.SpaceRE(regexp.MustCompile("[xy]")))).ToBe(t, xa1, xb1, ya1)

	as = nil
	expect.Slice(as.With(ns.Space("x"))).ToBe(t)
}

func TestAttrsValues(t *testing.T) {
	a1 := Attr("a", "", "1")
	b1 := Attr("b", "", "1")
	c2 := Attr("c", "", "2")
	d3 := Attr("d", "", "3")

	as := Attrs{a1, b1, c2, d3}

	expect.Slice(as.Values()).ToBe(t, "1", "1", "2", "3")

	as = nil
	expect.Slice(as.Values()).ToBeNil(t)
}

func TestAttrsToMap(t *testing.T) {
	a1 := Attr("a", "", "1")
	b1 := Attr("b", "", "1")
	c2 := Attr("c", "", "2")
	d3 := Attr("d", "", "3")

	as := Attrs{a1, b1, c2, d3}

	expect.Map(as.ToMap()).ToBe(t, map[xml.Name]string{
		{Local: "a"}: "1",
		{Local: "b"}: "1",
		{Local: "c"}: "2",
		{Local: "d"}: "3",
	})
}

func TestAttrsToSimpleMap(t *testing.T) {
	a1 := Attr("a", "", "1")
	b1 := Attr("b", "", "1")
	c2 := Attr("c", "", "2")
	d3 := Attr("d", "", "3")

	as := Attrs{a1, b1, c2, d3}

	expect.Map(as.ToSimpleMap()).ToBe(t, map[string]string{
		"a": "1",
		"b": "1",
		"c": "2",
		"d": "3",
	})
}

func TestCoerceAttrs(t *testing.T) {
	a1 := Attr("a", "", "true")
	b1 := Attr("b", "", "1")
	a2 := Attr("a", "", "2")
	xa3 := Attr("a", "x", "3")

	as := Attrs{a1, b1, a2, xa3}

	expect.Bool(Coerce(as, ns.Local("a"), strconv.ParseBool)).ToBeTrue(t)
	expect.Number(Coerce(as, ns.Name("a", "x"), strconv.Atoi)).ToBe(t, 3)
	expect.Number(Coerce(as, ns.Local("b"), strconv.Atoi)).ToBe(t, 1)
	expect.Number(Coerce(as, ns.Name("z", "z"), strconv.Atoi)).ToBe(t, 0)
}
