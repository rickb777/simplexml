package dom

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strings"
)

var TooManyRootElements = errors.New("no more than one root element is allowed")

func parseElement(decoder *xml.Decoder, tok xml.StartElement) (res *Element, err error) {
	res = CreateElement(tok.Name)
	for _, attr := range tok.Attr {
		res.AddAttr(attr)
	}

	for {
		newtok, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		switch rt := newtok.(type) {
		case xml.EndElement:
			return res, nil
		case xml.CharData:
			content := bytes.TrimSpace([]byte(rt.Copy()))
			if len(content) > 0 {
				res.Content = content
			}
		case xml.StartElement:
			child, err := parseElement(decoder, rt)
			if err != nil {
				return nil, err
			}
			res.AddChild(child)
		}
	}
}

// ParseElementString strictly parses the XML elements. If the input is malformed,
// an error is returned.
//
// This assumes our input is always UTF-8, no matter what the <?xml?> header says.
func ParseElementString(xml string) (elements []*Element, err error) {
	return ParseElements(strings.NewReader(xml))
}

// ParseElements strictly parses the XML elements. If the input is malformed,
// an error is returned.
//
// This assumes our input is always UTF-8, no matter what the <?xml?> header says.
func ParseElements(r io.Reader) (elements []*Element, err error) {
	decoder := xml.NewDecoder(r)
	decoder.Strict = true
	return ParseElementsWithDecoder(decoder)
}

// ParseElementsWithDecoder is like ParseElements but the decoder options can be specified.
func ParseElementsWithDecoder(decoder *xml.Decoder) (elements []*Element, err error) {
	elements = []*Element{}
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return elements, err
		}
		switch rt := tok.(type) {
		case xml.StartElement:
			element, err := parseElement(decoder, rt)
			if err != nil {
				return elements, err
			}
			elements = append(elements, element)
		}
	}
	return elements, nil
}

// ParseString strictly parses an XML document and returns a [Document] if input was well-formed.
// Otherwise, it returns an error.
func ParseString(xml string) (doc *Document, err error) {
	return Parse(strings.NewReader(xml))
}

// Parse strictly parses an XML document from a [io.Reader] and returns a [Document] if
// input was well-formed. Otherwise, it returns an error.
func Parse(r io.Reader) (doc *Document, err error) {
	decoder := xml.NewDecoder(r)
	decoder.Strict = true
	return ParseWithDecoder(decoder)
}

// ParseWithDecoder is like Parse but the decoder options can be specified.
func ParseWithDecoder(decoder *xml.Decoder) (doc *Document, err error) {
	elements, err := ParseElementsWithDecoder(decoder)
	if err != nil {
		return nil, err
	}
	if len(elements) > 1 {
		return nil, TooManyRootElements
	}
	doc = CreateDocument()
	if len(elements) == 1 {
		doc.SetRoot(elements[0])
	}
	return doc, nil
}
