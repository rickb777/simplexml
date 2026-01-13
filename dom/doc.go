// Package dom implements a simple XML DOM that is a light wrapper on top of
// encoding/xml.  It is oriented towards processing XML used as an RPC
// encoding mechanism (XMLRPC, SOAP, etc.), and not for general XML document
// processing.  Specifically:
//
//  1. We ignore comments and document processing directives.  They are stripped
//     out as part of document processing.
//
//  2. We do not have separate Text fields.  Instead, each [Element] has a single
//     Content field that holds the contents of the text enclosed within a tag.
package dom
