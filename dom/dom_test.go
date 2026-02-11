package dom

import (
	"encoding/xml"
	"errors"
	"log"
	"regexp"
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

func TestChildren(t *testing.T) {
	var (
		node1 = xml.Name{Local: "node1"}
		node2 = xml.Name{Local: "node2"}
	)
	root := parseDoc().Root()
	expect.Number(root.NumChildren()).ToBe(t, 3)
	expect.Number(root.ChildN(1).NumChildren()).ToBe(t, 1)
	expect.Number(root.ChildN(2).NumChildren()).ToBe(t, 0)
	expect.Slice(root.Children(Local("node1")).Names()).ToBe(t, node1)
	expect.Slice(root.Children().With(Local("node2")).Names()).ToBe(t, node2, node2)
	expect.Slice(root.Children().With(Local("node1"), Local("node2")).Names()).ToBe(t, node1, node2, node2)
	expect.Slice(root.Children().With(Local("node0")).Names()).ToBe(t)
	expect.Slice(root.Children().WithAttr(Local("foo")).Names()).ToBe(t, node1)
	expect.Slice(root.Children().WithAttr(Local("order")).Names()).ToBe(t, node2, node2)
	expect.Slice(root.Children().With(LocalRE(regexp.MustCompile("node."))).Names()).ToBe(t, node1, node2, node2)
	expect.Slice(root.Children().WithContent().Names()).ToBe(t, node2, node2)
	expect.Slice(root.Children().WithContent(regexp.MustCompile(".* Node 2.*")).Names()).ToBe(t, node2, node2)
	expect.Slice(root.Descendants().With(Local("node2")).Names()).ToBe(t, node2, node2, node2)
	expect.Slice(root.Descendants().WithContent(regexp.MustCompile(".*Groot.*")).Names()).ToBe(t, node2)
}

func TestChild(t *testing.T) {
	root := parseDoc().Root()
	expect.Value(root.Child(Local("node1")).Name).ToBe(t, xml.Name{Local: "node1"})
	expect.Value(root.Child(Local("foobar"))).ToBeNil(t)
}

func TestMoveChild(t *testing.T) {
	doc := parseDoc()
	root := doc.Root()
	node1 := root.Children()[0]
	sub := node1.Children()[0]
	sub.SetParent(doc.Root())
	// At this point, sub should be the 3rd of root's children.
	expect.Value(root.Children()[3]).ToBe(t, sub)
	// and trying to remove sub from node1 again should yield nil
	expect.Value(node1.RemoveChild(sub)).ToBeNil(t)
}

func TestElementRetrievalOrder(t *testing.T) {
	doc := parseDoc()
	res := doc.Root().All()
	expect.Slice(res).ToHaveLength(t, 6)
	for i, e := range res {
		var attr *xml.Attr
		for _, a := range e.Attributes {
			if a.Name.Local == "idx" {
				attr = &a
				break
			}
		}
		expect.Value(attr).Not().ToBe(t, nil)
		idx, err := strconv.Atoi(attr.Value)
		expect.Error(err).Not().ToHaveOccurred(t)
		expect.Number(idx).ToBe(t, i)
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
	expect.Slice(ancestors).ToHaveLength(t, 2)
	expect.Value(ancestors[0]).ToBe(t, node1)
	expect.Value(ancestors[1]).ToBe(t, root)
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
	expect.Slice(elements).ToHaveLength(t, 2)
	names := []xml.Name{
		{Local: "foo"},
		{Local: "bar"},
	}

	for i, n := range names {
		expect.Value(n).ToBe(t, elements[i].Name)
	}
}

func TestParseTooManyRootElements(t *testing.T) {
	elems := "<foo/>\n<bar/>\n"
	_, err := Parse(strings.NewReader(elems))
	expect.Error(err).ToHaveOccurred(t)
	if !errors.Is(err, TooManyRootElements) {
		t.Errorf("Expected TooManyRootElements, got %v", err)
	}
}
