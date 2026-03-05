package dom

import (
	"bytes"
	"encoding/xml"
	"io"
	"log"

	"github.com/rickb777/simplexml/ns"
)

// Element represents a node in an XML document.
// Elements are arranged in a tree that corresponds to
// the structure of the XML document.
type Element struct {
	Name     xml.Name
	children Elements
	parent   *Element
	// Unlike a full-fledged XML DOM, we only have a single Content field
	// instead of representing Text nodes separately.
	Content    []byte
	Attributes Attrs
}

// CreateElement creates a new element with the passed-in [xml.Name].
// The created [Element] has no parent, no children, no content, and no
// attributes.
func CreateElement(n xml.Name) *Element {
	return &Element{
		Name:       n,
		children:   Elements{},
		Attributes: Attrs{},
	}
}

// Elem creates a new [Element].  It is equivalent to creating a new
// [Element] with:
//
//	CreateElement(xml.Name{Local: name, Space: space})
func Elem(name, space string) *Element {
	return CreateElement(xml.Name{Space: space, Local: name})
}

// ElemC creates a new [Element] with content, e.g. "<b>hello</b>".
// It is equivalent to creating a new [Element] with:
//
//	e := Elem(name,space)
//	e.Content = []byte(content)
func ElemC(name, space, content string) *Element {
	res := Elem(name, space)
	res.Content = []byte(content)
	return res
}

// AddChild adds child to node.
// child will be reparented if needed.
// The altered node is returned.
func (node *Element) AddChild(child *Element) *Element {
	if child.parent != nil {
		child.parent.RemoveChild(child)
	}
	child.parent = node
	node.children = append(node.children, child)
	return node
}

// AddChildren adds children to the node.
// The children will be reparented as needed.
// The altered node is returned.
func (node *Element) AddChildren(children ...*Element) *Element {
	for _, c := range children {
		node.AddChild(c)
	}
	return node
}

// Replace performs an in-place replacement of node with other.
// The altered node is returned.
func (node *Element) Replace(other Element) *Element {
	node.Name = other.Name
	node.Content = other.Content
	node.Attributes = other.Attributes
	node.children = Elements{}
	node.AddChildren(other.children...)
	return node
}

// RemoveChild removes child from node, matching on its name and selecting the first
// match. The removed child will be returned if it was actually a child of node,
// otherwise nil will be returned. If multiple children would match the same name,
// use RemoveChildN instead.
func (node *Element) RemoveChild(child *Element) *Element {
	n := -1
	for i, v := range node.children {
		if v.Name == child.Name {
			n = i
			break
		}
	}

	if n == -1 {
		return nil
	}

	return node.RemoveChildN(n)
}

// RemoveChildN removes the nth child from node. For a group of N children,
// the indexes are 0 to N-1.
func (node *Element) RemoveChildN(n int) *Element {
	child := node.children[n]
	copy(node.children[n:], node.children[n+1:])
	node.children = node.children[0 : len(node.children)-1]
	child.parent = nil
	return child
}

// NumChildren gets the number of children.
func (node *Element) NumChildren() int {
	return len(node.children)
}

// ChildN returns nth child.
func (node *Element) ChildN(n int) *Element {
	return node.children[n]
}

// Child returns the first child with matching name. If there is
// no match, nil is returned. If multiple siblings may have the same
// name, use Children().With(...) instead.
func (node *Element) Child(pred ns.MatchName) *Element {
	for _, v := range node.children {
		if pred(v.Name) {
			return v
		}
	}
	return nil
}

// Children returns some or all the children of node.
func (node *Element) Children(pred ...ns.MatchName) Elements {
	res := make(Elements, 0, len(node.children))
	if len(pred) == 0 {
		return append(res, node.children...)
	}
	return node.children.WithName(pred...)
}

// Descendants returns all descendants of node as a flattened list in breadth-first order.
// The returned slice is a shallow copy.
func (node *Element) Descendants() Elements {
	res := make(Elements, 0, len(node.children))
	toProcess := node.Children()
	for len(toProcess) > 0 {
		nextToProcess := make(Elements, 0, len(node.children))
		for _, n := range toProcess {
			nextToProcess = append(nextToProcess, n.Children()...)
		}
		res = append(res, toProcess...)
		toProcess = nextToProcess
	}
	return res
}

// All returns node + Element.Descendants as a flattened list.
// The returned slice is a shallow copy.
func (node *Element) All() Elements {
	return append(Elements{node}, node.Descendants()...)
}

// Parent returns the parent of this node. If there is no parent, it returns nil.
func (node *Element) Parent() *Element {
	return node.parent
}

// Ancestors returns all the ancestors of this node with the most distant ancestor last.
func (node *Element) Ancestors() Elements {
	res := make(Elements, 0, 1)
	t := node.parent
	for t != nil {
		res = append(res, t)
		t = t.parent
	}
	return res
}

// SetParent makes parent the new parent of node, and returns node.
func (node *Element) SetParent(parent *Element) *Element {
	parent.AddChild(node)
	return node
}

//-------------------------------------------------------------------------------------------------

// AddAttr adds attr to node.
// Duplicates are ignored. If attr has the same name as a preexisting
// attribute, then it will replace the preexisting attribute.
// The altered node is returned.
func (node *Element) AddAttr(attr xml.Attr) *Element {
	for _, a := range node.Attributes {
		if a == attr {
			return node
		}
		if a.Name == attr.Name {
			a.Value = attr.Value
			return node
		}
	}
	node.Attributes = append(node.Attributes, attr)
	return node
}

// Attr creates a new [xml.Attr] and adds it to node.  It is equivalent to:
//
//	node.AddAttr(xml.Attr{
//	    Name: xml.Name{
//	        Space: space,
//	        Local: name,
//	    },
//	    Value: value,
//	})
//
// The altered node is returned.
func (node *Element) Attr(name, space, value string) *Element {
	return node.AddAttr(Attr(name, space, value))
}

// AttrIfNonBlank adds a new attribute to an element, but only if the value is not blank.
// This avoids the encoder emitting possibly many harmless but noisy blank attributes.
// The node is returned, whether altered or not.
func (node *Element) AttrIfNonBlank(name, space, value string) *Element {
	if value == "" {
		return node
	}
	return node.Attr(name, space, value)
}

func (node *Element) addNamespaces(encoder *Encoder) {
	// See if any of our attribs are in the xmlns namespace.
	// If they are, try to add them with their prefix
	for _, a := range node.Attributes {
		if a.Name.Space == "xmlns" {
			encoder.addNamespace(a.Value, a.Name.Local)
		}
	}

	encoder.addNamespace(node.Name.Space, "")
	for _, a := range node.Attributes {
		encoder.addNamespace(a.Name.Space, "")
	}
	for _, c := range node.children {
		c.addNamespaces(encoder)
	}
}

func namespacedName(e *Encoder, name xml.Name) string {
	if name.Space == "" {
		return name.Local
	}
	if name.Space == "xmlns" {
		return name.Space + ":" + name.Local
	}
	prefix, found := e.nsURLMap[name.Space]
	if !found {
		log.Panicf("No prefix found in %v for namespace %s", e.nsURLMap, name.Space)
	}
	return prefix + ":" + name.Local
}

//-------------------------------------------------------------------------------------------------

// Bytes returns the XML encoding of this part of the tree, with optional indentation.
func (node *Element) Bytes(indentation ...string) []byte {
	return node.bytes(indentation...).Bytes()
}

// Reader returns a [io.Reader] that can be used wherever
// something wants to consume this element tree.
func (node *Element) Reader() io.Reader {
	return node.bytes()
}

func (node *Element) bytes(indentation ...string) *bytes.Buffer {
	var b bytes.Buffer
	encoder := NewEncoder(&b, indentation...)
	// since we are encoding to a bytes.Buffer, assume Encode never fails.
	_ = node.Encode(encoder)
	_ = encoder.Flush()
	return &b
}

// String returns a pretty-printed XML encoding of this part of the tree.
func (node *Element) String() string {
	return string(node.Bytes("  "))
}
