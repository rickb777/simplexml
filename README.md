# Simplexml Dom library for Go

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rickb777/simplexml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/simplexml)](https://goreportcard.com/report/github.com/rickb777/simplexml)
[![Build](https://github.com/rickb777/simplexml/actions/workflows/go.yml/badge.svg)](https://github.com/rickb777/simplexml/actions)
[![Coverage](https://coveralls.io/repos/github/rickb777/simplexml/badge.svg?branch=master)](https://coveralls.io/github/rickb777/simplexml?branch=v2)
[![Issues](https://img.shields.io/github/issues/rickb777/simplexml.svg)](https://github.com/rickb777/simplexml/issues)

This is a naive and simple Go library to build a XML DOM to be able to produce
XML content, and parse simple XML into an in-memory DOM.

It started as a fork of https://github.com/masterzen/simplexml, but has
since been massively refactored to make it work more closely with encoding/xml,
and to include a set of useful functions for doing simple searches against the
element tree.

## Contact

- Bugs: https://github.com/rickb777/simplexml/issues


### Building

You can build the library from source:

```sh
git clone https://github.com/rickb777/simplexml
cd simplexml
go build
```

## Usage

