package dom

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/url"
)

// Encoder holds the state needed to encode the DOM into a well-formed XML document.
type Encoder struct {
	*bufio.Writer
	depth           int
	indentation     string
	started         bool
	namespacesAdded int
	nsPrefixMap     map[string]string
	nsURLMap        map[string]string
}

// NewEncoder returns a new [Encoder] that will write to the [io.Writer].
//
// The encoded document will have all namespace declarations lifted to the
// root element of the document.
//
// Optional indentation may be specified.
func NewEncoder(writer io.Writer, indentation ...string) *Encoder {
	enc := &Encoder{Writer: bufio.NewWriter(writer)}
	enc.nsPrefixMap = make(map[string]string)
	enc.nsURLMap = make(map[string]string)
	if len(indentation) > 0 {
		enc.indentation = indentation[0]
	}
	return enc
}

func (e *Encoder) addNamespace(ns string, prefix string) {
	if e.started {
		log.Panic("Cannot add element namespaces after encoding starts!")
	}
	if ns == "" || ns == "xmlns" {
		return
	}
	if prefix != "" {
		if _, found := e.nsURLMap[ns]; found {
			delete(e.nsPrefixMap, e.nsURLMap[ns])
			delete(e.nsURLMap, prefix)
		}
		e.nsPrefixMap[prefix] = ns
		e.nsURLMap[ns] = prefix
		return
	}

	if _, found := e.nsURLMap[ns]; found {
		return
	}
	if _, err := url.Parse(ns); err != nil {
		log.Panic(err)
	}
	prefix = fmt.Sprintf("ns%v", e.namespacesAdded)
	e.namespacesAdded++
	e.nsPrefixMap[prefix] = ns
	e.nsURLMap[ns] = prefix
}

// prettyEnd relies on bufio.Writer error propagation.
func (e *Encoder) prettyEnd() {
	if len(e.indentation) > 0 {
		_, _ = e.WriteString("\n")
	}
}

// spaces relies on bufio.Writer error propagation.
func (e *Encoder) spaces() {
	if len(e.indentation) > 0 {
		for i := 0; i < e.depth; i++ {
			_, _ = e.WriteString(e.indentation)
		}
	}
}

//-------------------------------------------------------------------------------------------------

// Encode encodes an element using the passed-in [Encoder].
// If an error occurs during encoding, that error is returned.
func (node *Element) Encode(e *Encoder) (err error) {
	// This could use some refactoring. but it works Well Enough(tm)
	writeNamespaces := !e.started
	if writeNamespaces {
		node.addNamespaces(e)
		e.started = true
	}

	e.spaces()

	_, _ = fmt.Fprintf(e, "<%s", namespacedName(e, node.Name))
	for _, a := range node.Attributes {
		if a.Name.Space != "xmlns" {
			_, _ = fmt.Fprintf(e, " %s=\"", namespacedName(e, a.Name))
			if err = xml.EscapeText(e, []byte(a.Value)); err != nil {
				return err
			}
			_, _ = e.Write([]byte{'"'})
		}
	}

	if writeNamespaces {
		for prefix, uri := range e.nsPrefixMap {
			_, _ = fmt.Fprintf(e, " xmlns:%s=\"%s\"", prefix, uri)
		}
	}

	if len(node.children) == 0 && len(node.Content) == 0 {
		ctag := "/>"
		if len(e.indentation) > 0 {
			ctag = "/>\n"
		}
		_, _ = e.WriteString(ctag)
		return e.Flush()
	}

	_, _ = e.WriteString(">")

	if len(node.Content) > 0 {
		if err := xml.EscapeText(e, node.Content); err != nil {
			return err
		}
	}

	if len(node.children) > 0 {
		e.depth++
		e.prettyEnd()
		for _, c := range node.children {
			if err = c.Encode(e); err != nil {
				return err
			}
		}
		e.depth--
		e.spaces()
	}

	_, _ = fmt.Fprintf(e, "</%s>", namespacedName(e, node.Name))
	e.prettyEnd()
	return e.Flush()
}
