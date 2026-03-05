package ns

import (
	"encoding/xml"
	"regexp"
	"testing"

	"github.com/rickb777/expect"
)

var n = xml.Name{
	Space: "space",
	Local: "local",
}

func TestTrue(t *testing.T) {
	expect.Bool(Any()(n)).ToBeTrue(t)
	expect.Bool(Local("local")(n)).ToBeTrue(t)
	expect.Bool(Space("space")(n)).ToBeTrue(t)
	expect.Bool(Name("local", "space")(n)).ToBeTrue(t)
	expect.Bool(LocalRE(regexp.MustCompile(".oca."))(n)).ToBeTrue(t)
	expect.Bool(SpaceRE(regexp.MustCompile(".pac."))(n)).ToBeTrue(t)
	expect.Bool(NameRE(regexp.MustCompile(".oca."), regexp.MustCompile(".pac."))(n)).ToBeTrue(t)
}

func TestFalse(t *testing.T) {
	expect.Bool(NotName(Any())(n)).ToBeFalse(t)
	expect.Bool(Local("xocal")(n)).ToBeFalse(t)
	expect.Bool(Space("xpace")(n)).ToBeFalse(t)
	expect.Bool(Name("xocal", "space")(n)).ToBeFalse(t)
	expect.Bool(Name("local", "xpace")(n)).ToBeFalse(t)
	expect.Bool(LocalRE(regexp.MustCompile("xoca."))(n)).ToBeFalse(t)
	expect.Bool(SpaceRE(regexp.MustCompile("xpac."))(n)).ToBeFalse(t)
	expect.Bool(NameRE(regexp.MustCompile("xoca."), regexp.MustCompile(".pac."))(n)).ToBeFalse(t)
	expect.Bool(NameRE(regexp.MustCompile(".oca."), regexp.MustCompile("xpac."))(n)).ToBeFalse(t)
}
