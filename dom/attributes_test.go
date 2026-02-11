package dom

import (
	"encoding/xml"
	"strconv"
	"testing"

	"github.com/rickb777/expect"
)

func TestAttrsGet(t *testing.T) {
	a1 := Attr("a", "", "1")
	b1 := Attr("b", "", "1")
	a2 := Attr("a", "", "2")
	xa3 := Attr("a", "x", "3")
	xb4 := Attr("b", "x", "4")
	ya5 := Attr("a", "y", "5")

	as := Attrs{a1, b1, a2, xa3, xb4, ya5}

	expect.Value(as.Get("a", "x")).ToBe(t, xa3)
	expect.Value(as.Get("a", "")).ToBe(t, a1)
	expect.Value(as.Get("z", "")).ToBe(t, xml.Attr{})

	as = nil
	expect.Value(as.Get("a", "")).ToBe(t, xml.Attr{})
}

func TestAttrsWithName(t *testing.T) {
	a1 := Attr("a", "", "1")
	b1 := Attr("b", "", "1")
	a2 := Attr("a", "", "2")
	xa1 := Attr("a", "x", "1")
	xb1 := Attr("b", "x", "1")
	ya1 := Attr("a", "y", "1")

	as := Attrs{a1, b1, a2, xa1, xb1, ya1}

	expect.Slice(as.WithName("a", "x")).ToBe(t, xa1)
	expect.Slice(as.WithName("a", "*")).ToBe(t, a1, a2, xa1, ya1)
	expect.Slice(as.WithName("b", "*")).ToBe(t, b1, xb1)
	expect.Slice(as.WithName("*", "x")).ToBe(t, xa1, xb1)

	as = nil
	expect.Slice(as.WithName("*", "x")).ToBe(t)
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

	expect.Bool(Coerce(as, "a", "", strconv.ParseBool)).ToBeTrue(t)
	expect.Number(Coerce(as, "a", "x", strconv.Atoi)).ToBe(t, 3)
	expect.Number(Coerce(as, "b", "", strconv.Atoi)).ToBe(t, 1)
	expect.Number(Coerce(as, "z", "z", strconv.Atoi)).ToBe(t, 0)
}
