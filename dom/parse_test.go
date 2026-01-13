package dom

import (
	"testing"

	"github.com/rickb777/expect"
)

type tc struct {
	name       string
	creator    func() *Document
	sample     string
	nameSpaces map[string]string
}

var testCases = []tc{
	{
		name: "EmptyDoc",
		creator: func() *Document {
			return CreateDocument()
		},
		sample: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n",
	},
	{
		name: "OneEmptyNode",
		creator: func() *Document {
			doc := CreateDocument()
			doc.SetRoot(Elem("root", ""))
			return doc
		},
		sample: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<root/>\n",
	},
	{
		name: "MoreNodes",
		creator: func() *Document {
			doc := CreateDocument()
			doc.SetRoot(
				Elem("root", "").AddChildren(
					Elem("node1", "").AddChild(Elem("sub", "")),
					Elem("node2", "")))
			return doc
		},
		sample: `<?xml version="1.0" encoding="UTF-8"?>
<root>
  <node1>
    <sub/>
  </node1>
  <node2/>
</root>
`,
	},
	{
		name: "WithAttribs",
		creator: func() *Document {
			doc := CreateDocument()
			doc.SetRoot(
				Elem("root", "").AddChild(
					Elem("node1", "").Attr("attr1", "", "pouet")))
			return doc
		},
		sample: `<?xml version="1.0" encoding="UTF-8"?>
<root>
  <node1 attr1="pouet"/>
</root>
`,
	},
	{
		name: "WithContent",
		creator: func() *Document {
			doc := CreateDocument()
			root := Elem("root", "")
			node1 := ElemC("node1", "", "this is a text content")
			root.AddChild(node1)
			doc.SetRoot(root)
			return doc
		},
		sample: `<?xml version="1.0" encoding="UTF-8"?>
<root>
  <node1>this is a text content</node1>
</root>
`,
	},
	{
		name: "WithNamespaces",
		creator: func() *Document {
			doc := CreateDocument()
			ns := "http://schemas.xmlsoap.org/ws/2004/08/addressing"
			root := Elem("root", "")
			node1 := Elem("node1", ns)
			root.AddChild(node1)
			node1.Content = []byte("this is a text content")
			doc.SetRoot(root)
			return doc
		},
		sample: `<?xml version="1.0" encoding="UTF-8"?>
<root xmlns:ns0="http://schemas.xmlsoap.org/ws/2004/08/addressing">
  <ns0:node1>this is a text content</ns0:node1>
</root>
`,
	},
}

func TestParsing(t *testing.T) {
	for _, testCase := range testCases {
		manualdoc := testCase.creator()
		parsedoc, err := ParseString(testCase.sample)
		if err != nil {
			t.Errorf("Cannot parse testcase %s sample %s\n\nGot error %v",
				testCase.name, testCase.sample, err)
		}
		if sample := manualdoc.String(); sample != testCase.sample {
			t.Errorf("Manually created DOM for %s did not render.\nExpected: %s\n\nGot: %s\n",
				testCase.name, testCase.sample, sample)
		}
		if sample := parsedoc.String(); sample != testCase.sample {
			t.Errorf("Parsed DOM for %s did not render.\nExpected: %s\n\nGot: %s\n",
				testCase.name, testCase.sample, sample)
		}
		autoparse, err := Parse(parsedoc.Reader())
		if err != nil {
			t.Errorf("Parsing new document from a document.Reader() failed: %v", err)
		}
		s1 := autoparse.String()
		s2 := parsedoc.String()
		if s1 != s2 {
			t.Errorf("Expected copy of DOM to be the same, but there are differences:\nExpected:%s\n\nGot: %s\n", s1, s2)
		}

	}
}

func TestMalformedEarlyParse(t *testing.T) {
	_, err := ParseString(`<?xml version="1.0" encoding="UTF-8"?><root`)
	expect.Error(err).Not().ToBeNil(t)
}

func TestMalformedMiddleParse(t *testing.T) {
	_, err := ParseString(`<?xml version="1.0" encoding="UTF-8"?><root><chil`)
	expect.Error(err).Not().ToBeNil(t)
}

func TestMalformedEndParse(t *testing.T) {
	_, err := ParseString(`<?xml version="1.0" encoding="UTF-8"?><root></roo`)
	expect.Error(err).Not().ToBeNil(t)
}
