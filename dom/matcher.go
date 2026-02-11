package dom

import (
	"encoding/xml"
	"regexp"
)

// NameMatcher is used for filtering elements and attributes by name.
type NameMatcher func(xml.Name) bool

// Not returns a predicate that inverts another matcher.
func Not(m NameMatcher) NameMatcher {
	return func(n xml.Name) bool {
		return !m(n)
	}
}

// Local returns a predicate that matches on the local name only.
func Local(local string) NameMatcher {
	return func(n xml.Name) bool {
		return local == n.Local
	}
}

// LocalRE returns a predicate that matches on the local name only.
func LocalRE(local *regexp.Regexp) NameMatcher {
	return func(n xml.Name) bool {
		return local.MatchString(n.Local)
	}
}

// Space returns a predicate that matches on the name space only.
func Space(space string) NameMatcher {
	return func(n xml.Name) bool {
		return space == n.Space
	}
}

// SpaceRE returns a predicate that matches on the name space only.
func SpaceRE(space *regexp.Regexp) NameMatcher {
	return func(n xml.Name) bool {
		return space.MatchString(n.Space)
	}
}

// Name returns a predicate that matches on both the local name and name space.
func Name(local, space string) NameMatcher {
	return func(n xml.Name) bool {
		return local == n.Local && space == n.Space
	}
}

// NameRE returns a predicate that matches on both the local name and name space.
func NameRE(local, space *regexp.Regexp) NameMatcher {
	return func(n xml.Name) bool {
		return local.MatchString(n.Local) && space.MatchString(n.Space)
	}
}
