package dom

import (
	"encoding/xml"
	"errors"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/rickb777/expect"
)

var testDoc = `<?xml version="1.0" encoding="UTF-8"?>
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

func parseDoc() *Document {
	doc, err := Parse(strings.NewReader(testDoc))
	if err != nil {
		log.Panicf("Cannot parse test document. Error: %v", err)
	}
	return doc
}

func TestMoveChild(t *testing.T) {
	doc := parseDoc()
	root := doc.Root()
	node1 := root.Children()[0]
	sub := node1.Children()[0]
	sub.SetParent(doc.Root())
	// At this point, sub should be the 3rd of root's children.
	if root.Children()[3] != sub {
		t.Error("Failed to move sub from node1 to root")
	}
	// and trying to remove sub from node1 again should yield nil
	if node1.RemoveChild(sub) != nil {
		t.Error("sub is not a child of node1, but trying to remove it worked.")
	}
}

func TestElementRetrievalOrder(t *testing.T) {
	doc := parseDoc()
	res := doc.Root().All()
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

func TestAncestorOrder(t *testing.T) {
	doc := parseDoc()
	root := doc.Root()
	node1 := root.Children()[0]
	sub := node1.Children()[0]
	// Test the Parent() method while we are at it.
	if subParent := sub.Parent(); subParent != node1 {
		t.Errorf("sub should have %v as its parent, not %v", node1.Name, subParent.Name)
	}
	ancestors := sub.Ancestors()
	if len(ancestors) != 2 {
		t.Errorf("sub should have 2 ancestors, not %d", len(ancestors))
	}
	if ancestors[0] != node1 {
		t.Errorf("sub should have %v as its first ancestor, not %v",
			node1.Name, ancestors[0].Name)
	}
	if ancestors[1] != root {
		t.Errorf("sub should have %v as its second ancestor, not %v",
			node1.Name, ancestors[1].Name)
	}
}

func TestEncoding(t *testing.T) {
	doc := parseDoc()
	var sb strings.Builder
	enc := NewEncoder(&sb, "  ")
	err := doc.Encode(enc)
	expect.Error(err).ToBeNil(t)
	expect.String(sb.String()).ToBe(t, testDoc)
}

func TestElementString(t *testing.T) {
	refString := "<foo/>\n"
	refElement := Elem("foo", "")
	if res := refElement.String(); res != refString {
		t.Errorf("Expected stringification of reference to be '%s', not '%s'", refString, res)
	}
}

func TestParseElements(t *testing.T) {
	elems := "<foo/>\n<bar/>\n"
	elements, err := ParseElementString(elems)
	expect.Error(err).ToBeNil(t)
	if len(elements) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(elements))
	}
	names := []xml.Name{
		{Local: "foo"},
		{Local: "bar"},
	}

	for i, e := range names {
		if elements[i].Name != e {
			t.Errorf("Expected first element to be %v, it is %v", e, elements[i].Name)
		}
	}
}

func TestParseTooManyRootElements(t *testing.T) {
	elems := "<foo/>\n<bar/>\n"
	_, err := Parse(strings.NewReader(elems))
	if err == nil {
		t.Errorf("Did not get expected error parsing XML document %s", elems)
	}
	if !errors.Is(err, TooManyRootElements) {
		t.Errorf("Expected TooManyRootElements, got %v", err)
	}
}
