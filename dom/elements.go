package dom

import (
	"encoding/xml"

	"github.com/rickb777/simplexml/ns"
)

// Elements is a slice of Element pointers. Several filtering methods are provided:
// supplying multiple predicates will operate as a logical "or". Chaining the methods
// will operate as a logical "and".
type Elements []*Element

// Match is used for filtering elements.
// Note that this is the same as search.Match, which provides a selection of matchers.
type Match func(*Element) bool

// NameMatch adapts a name matcher as an element matcher based on the element's name.
func NameMatch(nm ns.MatchName) Match {
	return func(el *Element) bool {
		return nm(el.Name)
	}
}

// AttrNameMatch adapts a name matcher as an element matcher based on the name
// of any of the element's attributes.
func AttrNameMatch(nm ns.MatchName) Match {
	return func(el *Element) bool {
		for _, a := range el.Attributes {
			if nm(a.Name) {
				return true
			}
		}
		return false
	}
}

func NameMatchAny(nm ...ns.MatchName) []Match {
	res := make([]Match, len(nm))
	for i, fn := range nm {
		res[i] = NameMatch(fn)
	}
	return res
}

// First finds the first element that matches, or else nil.
// If multiple predicates are supplied, any matching element is returned.
// If els or pred is empty, nil is returned.
func (els Elements) First(pred ...Match) *Element {
	if len(els) == 0 {
		return nil
	}
	for _, e := range els {
		for _, p := range pred {
			if p(e) {
				return e
			}
		}
	}
	return nil
}

// With filters the elements and returns only those that match.
// If multiple predicates are supplied, any matching element is returned.
// If els or pred is empty, nil is returned.
func (els Elements) With(pred ...Match) Elements {
	if len(els) == 0 {
		return nil
	}
	res := make(Elements, 0, len(els))
	for _, e := range els {
	inner:
		for _, p := range pred {
			if p(e) {
				res = append(res, e)
				break inner
			}
		}
	}
	return res
}

// WithName filters the elements and returns only those that have matching name.
// If multiple predicates are supplied, any matching element is returned.
// If els or pred is empty, nil is returned.
func (els Elements) WithName(pred ...ns.MatchName) Elements {
	if len(els) == 0 || len(pred) == 0 {
		return nil
	}
	res := make(Elements, 0, len(els))
	for _, e := range els {
	inner:
		for _, p := range pred {
			if p(e.Name) {
				res = append(res, e)
				break inner
			}
		}
	}
	return res
}

// WithAttr filters the elements and returns only those that have any attribute
// with a matching name. If multiple predicates are supplied, any matching
// element is returned. If els or pred is empty, nil is returned.
func (els Elements) WithAttr(pred ...ns.MatchName) Elements {
	if len(els) == 0 || len(pred) == 0 {
		return nil
	}
	res := make(Elements, 0, len(els))
	for _, el := range els {
	inner:
		for _, a := range el.Attributes {
			for _, p := range pred {
				if p(a.Name) {
					res = append(res, el)
					break inner
				}
			}
		}
	}
	return res
}

// Partition divides the elements into those that match and those that don't match.
// If multiple predicates are supplied, any matching element is returned.
// If els or pred is empty, nil is returned.
func (els Elements) Partition(pred ...Match) (match, others Elements) {
	if len(els) == 0 || len(pred) == 0 {
		return nil, nil
	}
	match = make(Elements, 0, len(els))
	others = make(Elements, 0, len(els))
	for _, el := range els {
		matched := false
	inner:
		for _, p := range pred {
			if p(el) {
				match = append(match, el)
				matched = true
				break inner
			}
		}
		if !matched {
			others = append(others, el)
		}
	}
	return match, others
}

// Names gets the XML names from the elements in original order.
func (els Elements) Names() []xml.Name {
	if len(els) == 0 {
		return nil
	}
	res := make([]xml.Name, 0, len(els))
	for _, el := range els {
		res = append(res, el.Name)
	}
	return res
}
