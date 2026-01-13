# SimpleXML Dom library for Go

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rickb777/simplexml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/simplexml)](https://goreportcard.com/report/github.com/rickb777/simplexml)
[![Build](https://github.com/rickb777/simplexml/actions/workflows/go.yml/badge.svg)](https://github.com/rickb777/simplexml/actions)
[![Coverage](https://coveralls.io/repos/github/rickb777/simplexml/badge.svg?branch=master)](https://coveralls.io/github/rickb777/simplexml?branch=master)
[![Issues](https://img.shields.io/github/issues/rickb777/simplexml.svg)](https://github.com/rickb777/simplexml/issues)

This is a deliberately-simple Go library 

 * to build XML DOM in memory, 
 * to produce XML content, and
 * to parse XML into an in-memory DOM.

## Origins

This started as a fork of [VictorLowther/simplexml](https://github.com/VictorLowther/simplexml) 
(which was originally from [masterzen/simplexml](https://github.com/masterzen/simplexml)), but has
since been massively refactored to make it work more closely with [encoding/xml](https://pkg.go.dev/encoding/xml),
and to include a set of functions for doing simple searches against the element tree.

### Building

To use, simply

```sh
go get github.com/rickb777/simplexml
```

You can build the library from source, e.g.

```sh
git clone https://github.com/rickb777/simplexml
cd simplexml
go build
```

It has a [Magefile](https://magefile.org).

