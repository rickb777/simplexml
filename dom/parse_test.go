package dom_test

import (
	"testing"

	"github.com/rickb777/expect"
	"github.com/rickb777/simplexml/dom"
)

type tc struct {
	name       string
	creator    func() *dom.Document
	sample     string
	nameSpaces map[string]string
}

var testCases = []tc{
	{
		name: "EmptyDoc",
		creator: func() *dom.Document {
			return dom.CreateDocument()
		},
		sample: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n",
	},
	{
		name: "OneEmptyNode",
		creator: func() *dom.Document {
			doc := dom.CreateDocument()
			doc.SetRoot(dom.Elem("root", ""))
			return doc
		},
		sample: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<root/>\n",
	},
	{
		name: "MoreNodes",
		creator: func() *dom.Document {
			doc := dom.CreateDocument()
			doc.SetRoot(
				dom.Elem("root", "").AddChildren(
					dom.Elem("node1", "").AddChild(dom.Elem("sub", "")),
					dom.Elem("node2", "")))
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
		creator: func() *dom.Document {
			doc := dom.CreateDocument()
			doc.SetRoot(
				dom.Elem("root", "").AddChild(
					dom.Elem("node1", "").AttrIfNonBlank("id", "", `"Fran & Freddie's Diner" <tasty@example.com>`)))
			return doc
		},
		sample: `<?xml version="1.0" encoding="UTF-8"?>
<root>
  <node1 id="&#34;Fran &amp; Freddie&#39;s Diner&#34; &lt;tasty@example.com&gt;"/>
</root>
`,
	},
	{
		name: "WithContent",
		creator: func() *dom.Document {
			doc := dom.CreateDocument()
			root := dom.Elem("root", "")
			node1 := dom.ElemC("node1", "", "this is a text content including < and >")
			root.AddChild(node1)
			doc.SetRoot(root)
			return doc
		},
		sample: `<?xml version="1.0" encoding="UTF-8"?>
<root>
  <node1>this is a text content including &lt; and &gt;</node1>
</root>
`,
	},
	{
		name: "WithNamespaces",
		creator: func() *dom.Document {
			doc := dom.CreateDocument()
			ns := "http://schemas.xmlsoap.org/ws/2004/08/addressing"
			root := dom.Elem("root", "")
			node1 := dom.Elem("node1", ns)
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
		parsedoc, err := dom.ParseString(testCase.sample)
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
		autoparse, err := dom.Parse(parsedoc.Reader())
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
	_, err := dom.ParseString(`<?xml version="1.0" encoding="UTF-8"?><root`)
	expect.Error(err).Not().ToBeNil(t)
}

func TestMalformedMiddleParse(t *testing.T) {
	_, err := dom.ParseString(`<?xml version="1.0" encoding="UTF-8"?><root><chil`)
	expect.Error(err).Not().ToBeNil(t)
}

func TestMalformedEndParse(t *testing.T) {
	_, err := dom.ParseString(`<?xml version="1.0" encoding="UTF-8"?><root></roo`)
	expect.Error(err).Not().ToBeNil(t)
}
