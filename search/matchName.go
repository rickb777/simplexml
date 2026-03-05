// Package search provides matchers that can be applied to dom.Elements.
// For some basic usage examples, see match_test.go
package search

import (
	"regexp"

	"github.com/rickb777/simplexml/dom"
	"github.com/rickb777/simplexml/ns"
)

// These functions provide adapters so that name matchers can be applied to elements.

// Local returns a predicate that matches on the local name only.
func Local(local string) Match {
	return dom.NameMatch(ns.Local(local))
}

// LocalRE returns a predicate that matches on the local name only.
func LocalRE(local *regexp.Regexp) Match {
	return dom.NameMatch(ns.LocalRE(local))
}

// Space returns a predicate that matches on the name space only.
func Space(space string) Match {
	return dom.NameMatch(ns.Space(space))
}

// SpaceRE returns a predicate that matches on the name space only.
func SpaceRE(space *regexp.Regexp) Match {
	return dom.NameMatch(ns.SpaceRE(space))
}

// Name returns a predicate that matches on both the local name and name space.
func Name(local, space string) Match {
	return dom.NameMatch(ns.Name(local, space))
}

// NameRE returns a predicate that matches on both the local name and name space.
func NameRE(local, space *regexp.Regexp) Match {
	return dom.NameMatch(ns.NameRE(local, space))
}
