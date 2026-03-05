// Package ns provides matchers for XML names.
package ns

import (
	"encoding/xml"
	"regexp"
)

// MatchName is used for filtering elements and attributes by name.
type MatchName func(xml.Name) bool

// NotName returns a predicate that inverts another matcher.
func NotName(m MatchName) MatchName {
	return func(n xml.Name) bool {
		return !m(n)
	}
}

// Any returns a predicate that matches any name.
func Any() MatchName {
	return func(n xml.Name) bool { return true }
}

// Local returns a predicate that matches on the local name only.
func Local(local string) MatchName {
	return func(n xml.Name) bool {
		return local == n.Local
	}
}

// LocalRE returns a predicate that matches on the local name only.
func LocalRE(local *regexp.Regexp) MatchName {
	return func(n xml.Name) bool {
		return local.MatchString(n.Local)
	}
}

// Space returns a predicate that matches on the name space only.
func Space(space string) MatchName {
	return func(n xml.Name) bool {
		return space == n.Space
	}
}

// SpaceRE returns a predicate that matches on the name space only.
func SpaceRE(space *regexp.Regexp) MatchName {
	return func(n xml.Name) bool {
		return space.MatchString(n.Space)
	}
}

// Name returns a predicate that matches on both the local name and name space.
func Name(local, space string) MatchName {
	return func(n xml.Name) bool {
		return local == n.Local && space == n.Space
	}
}

// NameRE returns a predicate that matches on both the local name and name space.
func NameRE(local, space *regexp.Regexp) MatchName {
	return func(n xml.Name) bool {
		return local.MatchString(n.Local) && space.MatchString(n.Space)
	}
}
