package dom

import (
	"bytes"
	"io"
)

// A Document represents an entire XML document.  Documents hold the root Element.
type Document struct {
	root *Element
}

// CreateDocument creates a new XML document.
func CreateDocument() *Document {
	return &Document{}
}

// Root returns the root element of the document.
func (doc *Document) Root() (node *Element) {
	return doc.root
}

// SetRoot sets a new root element of the document.
func (doc *Document) SetRoot(node *Element) {
	node.parent = nil
	doc.root = node
}

// Encode encodes the entire [Document] using the [Encoder].
// The output is a well-formed XML document.
func (doc *Document) Encode(e *Encoder) error {
	_, _ = e.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	e.prettyEnd()
	if doc.root != nil {
		return doc.root.Encode(e)
	}
	return e.Flush()
}

// Bytes encodes a [Document] into a byte array. It can optionally be indented.
func (doc *Document) Bytes(indentation ...string) []byte {
	return doc.bytes(indentation...).Bytes()
}

// Reader returns a [io.Reader] that can be used wherever
// something wants to consume this document.
func (doc *Document) Reader() io.Reader {
	return doc.bytes()
}

func (doc *Document) bytes(indentation ...string) *bytes.Buffer {
	var b bytes.Buffer
	encoder := NewEncoder(&b, indentation...)
	// since we are encoding to a bytes.Buffer, assume Encode never fails.
	_ = doc.Encode(encoder)
	_ = encoder.Flush()
	return &b
}

// String converts to a string the result of [Document.Bytes] with 2-space indentation.
func (doc *Document) String() string {
	return string(doc.Bytes("  "))
}
