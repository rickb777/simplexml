package search

import (
	"encoding/xml"
	"log"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/rickb777/simplexml/dom"
)

var testDoc string = `<?xml version="1.0" encoding="UTF-8"?>
<a:root idx="0" xmlns:a="http://schemas.xmlsoap.org/ws/2004/08/addressing">
 <node1 foo="bar" idx="1">
  <sub idx="4"/>
 </node1>
 <node2 order="0" idx="2">I am Node 2
  <node2 order="2" idx="5">I am Groot</node2>
 </node2>
 <node2 order="1" idx="3">I am a different Node 2</node2>
</a:root>
`

func parseDoc() *dom.Document {
	doc, err := dom.Parse(strings.NewReader(testDoc))
	if err != nil {
		log.Panicf("Cannot parse test document. Error: %v", err)
	}
	return doc
}

func TestTag(t *testing.T) {
	doc := parseDoc()
	res := First(Tag("sub", ""), doc.Root().All())
	if res == nil {
		t.Error("Could not find sub element!")
	}
	if res.Name.Local != "sub" || res.Name.Space != "" {
		t.Errorf("Looking for sub element gave me '%s' in namespace '%s'", res.Name.Local, res.Name.Space)
	}
}

func TestTagRE(t *testing.T) {
	doc := parseDoc()
	res := First(TagRE(regexp.MustCompile("^sub$"), nil), doc.Root().All())
	if res == nil {
		t.Error("Could not find sub element!")
	}
	if res.Name.Local != "sub" || res.Name.Space != "" {
		t.Errorf("Looking for sub element gave me '%s' in namespace '%s'", res.Name.Local, res.Name.Space)
	}
}

func TestTagOrderPreservation(t *testing.T) {
	doc := parseDoc()
	res := All(Tag("node2", ""), doc.Root().All())
	if len(res) != 3 {
		t.Errorf("Expected to find 2 elements, found %d", len(res))
	}
	for i, e := range res {
		if e.Name.Local != "node2" || e.Name.Space != "" {
			t.Errorf("Looking for node2 element gave me '%s' in namespace '%s'", e.Name.Local, e.Name.Space)
		}
		attr := e.Attributes[0]
		if attr.Name.Local != "order" {
			t.Error("Could not find expected order attribute on node2 element")
		}
		order, err := strconv.Atoi(attr.Value)
		if err != nil {
			t.Errorf("Could not extract order attribute value: %v", err)
		}
		if order != i {
			t.Errorf("Elements returned by All out of order! Expected %d, got %d", i, order)
		}
	}
}

func TestAttr(t *testing.T) {
	doc := parseDoc()
	res := All(Attr("idx", "", "*"), doc.Root().All())
	if len(res) != 6 {
		t.Errorf("Expected 6 elements, got %d", len(res))
	}
	res = All(Attr("idx", "", "0"), doc.Root().All())
	if len(res) != 1 {
		t.Errorf("Expected 1 elements, got %d", len(res))
	}
	if res[0] != doc.Root() {
		t.Errorf("Expected to get element %v, got %v",
			doc.Root().Name,
			res[0].Name)
	}
}

func TestAttrRE(t *testing.T) {
	doc := parseDoc()
	res := All(AttrRE(regexp.MustCompile("idx"), nil, regexp.MustCompile(".*")),
		doc.Root().All())
	if len(res) != 6 {
		t.Errorf("Expected 6 elements, got %d", len(res))
	}
	res = All(AttrRE(regexp.MustCompile("idx"), nil, regexp.MustCompile("^0$")),
		doc.Root().All())
	if len(res) != 1 {
		t.Errorf("Expected 1 elements, got %d", len(res))
	}
	if res[0] != doc.Root() {
		t.Errorf("Expected to get element %v, got %v",
			doc.Root().Name,
			res[0].Name)
	}
}

func TestAttrOrderPreservation(t *testing.T) {
	doc := parseDoc()
	res := All(Attr("idx", "", "*"), doc.Root().All())
	if len(res) != 6 {
		t.Errorf("Expected 6 elements, got %d", len(res))
	}

	for i, e := range res {
		var attr *xml.Attr
		for _, a := range e.Attributes {
			if a.Name.Local == "idx" {
				attr = &a
				break
			}
		}
		if attr == nil {
			t.Errorf("Could not find idx addr on element %s", e.Name.Local)
		}
		idx, err := strconv.Atoi(attr.Value)
		if err != nil {
			t.Errorf("Could not extract idx attribute value: %v", err)
		}
		if idx != i {
			t.Errorf("Elements returned by attr search are out of order.  Expected %d, got %d", i, idx)
		}
	}
}

func TestContentExists(t *testing.T) {
	doc := parseDoc()
	res := All(ContentExists(), doc.Root().All())
	if len(res) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(res))
	}
}

func TestContentRE(t *testing.T) {
	doc := parseDoc()
	res := All(ContentRE(regexp.MustCompile("^I am Groot$")), doc.Root().All())
	if len(res) != 1 {
		t.Errorf("Expected 1 element, got %d", len(res))
	}
	expected := "I am Groot"
	if string(res[0].Content) != expected {
		t.Errorf("Expected node content %s, got %s",
			expected, string(res[0].Content))
	}
}

func TestAndCombinator(t *testing.T) {
	doc := parseDoc()
	expected := "I am a different Node 2"
	res := All(And(
		Attr("*", "", "1"),
		Tag("node2", "")),
		doc.Root().All())
	if len(res) != 1 {
		t.Errorf("Expected 1 element, got %d", len(res))
	}
	if string(res[0].Content) != expected {
		t.Errorf("Expected node content not found!\nExpected: %s\n\nGot: %s",
			expected,
			string(res[0].Content))
	}
}

func TestOrCombinator(t *testing.T) {
	doc := parseDoc()
	res := All(Or(
		Attr("idx", "", "0"),
		Attr("foo", "", "bar")),
		doc.Root().All())
	if len(res) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(res))
	}
	if res[0].Name.Space != "http://schemas.xmlsoap.org/ws/2004/08/addressing" ||
		res[0].Name.Local != "root" {
		t.Errorf("Expected first element to be root, not %s", res[0].Name.Local)
	}
	if res[1].Name.Local != "node1" {
		t.Errorf("Expected second element to be node1, not %s", res[1].Name.Local)
	}
}

func TestNot(t *testing.T) {
	doc := parseDoc()
	match := Not(Tag("root", ""))
	if !match(doc.Root()) {
		t.Error("Not match testing failed!")
	}
}

func TestNoParent(t *testing.T) {
	doc := parseDoc()
	res := All(NoParent(),
		doc.Root().All())
	if len(res) != 1 {
		t.Errorf("NoParent matched %d elements, should only have matched 1", len(res))
	}
	if res[0] != doc.Root() {
		t.Errorf("NoParent matched element %v, instead of %v",
			res[0].Name,
			doc.Root().Name)
	}
}

func TestAncestor(t *testing.T) {
	doc := parseDoc()
	// The only node that will fail this test is the root node,
	// because it does not have ancestors.
	res := All(Ancestor(NoParent()),
		doc.Root().All())
	if len(res) != 5 {
		t.Errorf("TestAncestor matched %d elements instead of 5", len(res))
	}
}

func TestNestedAncestor(t *testing.T) {
	doc := parseDoc()
	res := All(Ancestor(Ancestor(NoParent())),
		doc.Root().All())
	if len(res) != 2 {
		t.Errorf("TestNestedAncestor matched %d elements instead of 2", len(res))
	}
}

func TestAncestorN(t *testing.T) {
	doc := parseDoc()
	answers := []int{1, 3, 2}
	for i, answer := range answers {
		res := All(AncestorN(NoParent(), uint(i)),
			doc.Root().All())
		if len(res) != answer {
			t.Errorf("TestAncestorN(%d) had %d matches instead of %d",
				i, len(res), answer)
		}
	}
}

func TestParent(t *testing.T) {
	doc := parseDoc()
	res := All(Parent(NoParent()),
		doc.Root().All())
	if len(res) != 3 {
		t.Errorf("TestParent matched %d elements instead of 3", len(res))
	}
}

func TestNestedParent(t *testing.T) {
	doc := parseDoc()
	res := All(Parent(Parent(NoParent())),
		doc.Root().All())
	if len(res) != 2 {
		t.Errorf("TestNestedParent matched %d elements instead of 2", len(res))
	}
}

func TestAlways(t *testing.T) {
	doc := parseDoc()
	if !Always()(doc.Root()) {
		t.Error("Always returned false")
	}
}

func TestNever(t *testing.T) {
	doc := parseDoc()
	if Never()(doc.Root()) {
		t.Error("Never returned true")
	}
}
