// Package search provides matchers that can be applied to dom.Elements.
// For some basic usage examples, see match_test.go
package search

import (
	"bytes"
	"regexp"

	"github.com/rickb777/simplexml/dom"
	"github.com/rickb777/simplexml/ns"
)

// Match is the basic type of a search function. It takes a single element
// and returns a boolean indicating whether the element matched the func.
//
// Note that this is the same as dom.Match.
type Match = dom.Match

// And takes any number of Match functions, and returns another Match that will
// match if all of funcs match. Many operations on dom.Elements can be chained
// and this has the same effect.
func And(funcs ...Match) Match {
	return func(el *dom.Element) bool {
		for _, fn := range funcs {
			if !fn(el) {
				return false
			}
		}
		return true
	}
}

// Or takes any number of Match functions, and returns another Match that
// will match if any of funcs match. Many operations on dom.Elements have
// multiple predicates and this has the same effect.
func Or(funcs ...Match) Match {
	return func(el *dom.Element) bool {
		for _, fn := range funcs {
			if fn(el) {
				return true
			}
		}
		return false
	}
}

// Not takes a single Match, and returns another Match
// that matches if fn does not match.
func Not(fn Match) Match {
	return func(el *dom.Element) bool {
		return !fn(el)
	}
}

// NoParent returns a matcher that matches iff the element
// does not have a parent
func NoParent() Match {
	return func(el *dom.Element) bool {
		return el.Parent() == nil
	}
}

// Ancestor returns a matcher that matches iff the element has an
// ancestor that matches the passed matcher.
func Ancestor(fn Match) Match {
	return func(el *dom.Element) bool {
		return First(el.Ancestors(), fn) != nil
	}
}

// AncestorN returns a matcher that matches against the
// nth ancestor of the node being tested.
// If n == 0, then the node itself will be tested as a degenerate case.
// If there is no such ancestor the match fails.
func AncestorN(fn Match, distance uint) Match {
	return func(el *dom.Element) bool {
		if distance == 0 {
			return fn(el)
		}
		ancestors := el.Ancestors()
		if len(ancestors) < int(distance) {
			return false
		}
		return fn(ancestors[distance-1])
	}
}

// Parent returns a matcher that matches iff the element
// has a parent and that parent matches the passed fn.
func Parent(fn Match) Match {
	return func(el *dom.Element) bool {
		p := el.Parent()
		if p == nil {
			return false
		}
		return fn(p)
	}
}

// Child returns a matcher that matches iff the element has a child that
// matches the passed fn. Applying this has a similar effect to using
// dom.Elements With(fn) on an element's children.
func Child(fn Match) Match {
	return func(el *dom.Element) bool {
		for _, c := range el.Children() {
			if fn(c) {
				return true
			}
		}
		return false
	}
}

// Always returns a matcher that always matches
func Always() Match {
	return func(el *dom.Element) bool {
		return true
	}
}

// Never returns a matcher that never matches
func Never() Match {
	return Not(Always())
}

// All returns all the nodes that any fn matches
func All(els dom.Elements, fn ...Match) dom.Elements {
	return els.With(fn...)
}

// Partition divides the nodes into those that match and all the others.
func Partition(els dom.Elements, fn ...Match) (matching, others dom.Elements) {
	return els.Partition(fn...)
}

// First returns the first element that any fn matches.
func First(els dom.Elements, fn ...Match) *dom.Element {
	return els.First(fn...)
}

// Attr creates a Match against the names of the attributes of an element.
func Attr(matchName ns.MatchName) Match {
	return func(el *dom.Element) bool {
		for _, a := range el.Attributes {
			if matchName(a.Name) {
				return true
			}
		}
		return false
	}
}

// AttrV creates a Match against the names of the attributes of an element and a
// specified attribute value, which must also match.
func AttrV(matchName ns.MatchName, value string) Match {
	return func(el *dom.Element) bool {
		for _, a := range el.Attributes {
			if matchName(a.Name) && (value == a.Value) {
				return true
			}
		}
		return false
	}
}

// AttrRE creates a Match against the names of the attributes of an element.
func AttrRE(matchName ns.MatchName, value *regexp.Regexp) Match {
	return func(el *dom.Element) bool {
		for _, a := range el.Attributes {
			if matchName(a.Name) && value.MatchString(a.Value) {
				return true
			}
		}
		return false
	}
}

// ContentExists creates a Match against an element that has non-empty Content.
func ContentExists() Match {
	return func(el *dom.Element) bool {
		return len(el.Content) > 0
	}
}

// ContentIs creates a Match against an element that tests to see if
// it matches the supplied content.
func ContentIs(content []byte) Match {
	return func(el *dom.Element) bool {
		return bytes.Equal(el.Content, content)
	}
}

// ContentContains creates a Match against an element that tests to see if
// it contains the supplied content.
func ContentContains(subslice []byte) Match {
	return func(el *dom.Element) bool {
		return bytes.Contains(el.Content, subslice)
	}
}

// ContentRE creates a Match that applies a regular expression to the Content of an element.
func ContentRE(regex *regexp.Regexp) Match {
	return func(el *dom.Element) bool {
		return regex.Match(el.Content)
	}
}
