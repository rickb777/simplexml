package dom

import "encoding/xml"

// Attr creates a new [xml.Attr].  It is exactly equivalent to creating
// a new [xml.Attr] with:
//
//	xml.Attr{
//	    Name: xml.Name{
//	        Local: name,
//	        Space: space,
//	    },
//	    Value: value,
//	}
func Attr(name, space, value string) xml.Attr {
	return xml.Attr{
		Name:  xml.Name{Space: space, Local: name},
		Value: value,
	}
}

// Attrs is a slice of [xml.Attr].
type Attrs []xml.Attr

// Names gets the XML names from the attributes in original order.
func (as Attrs) Names() []xml.Name {
	if len(as) == 0 {
		return nil
	}
	res := make([]xml.Name, 0, len(as))
	for _, a := range as {
		res = append(res, a.Name)
	}
	return res
}

// Get gets the attribute that has matching name.
// If there is no match, a zero xml.Attr is returned.
func (as Attrs) Get(name NameMatcher) xml.Attr {
	for _, a := range as {
		if name(a.Name) {
			return a
		}
	}
	return xml.Attr{}
}

// With filters the attributes and returns only those that have
// matching name. If multiple predicates are supplied, any matching
// attribute is returned. If as is empty, nil is returned.
func (as Attrs) With(pred ...NameMatcher) Attrs {
	if len(as) == 0 {
		return nil
	}
	res := make(Attrs, 0, len(as))
	for _, a := range as {
	inner:
		for _, p := range pred {
			if p(a.Name) {
				res = append(res, a)
				break inner
			}
		}
	}
	return res
}

// Values extracts the value strings from all of the Attrs. A typical use-case
// is to get a single value from a list just of one xml.Attr.
// The result is a slice with the same length as 'as'; it will be nil if there
// are no attributes.
func (as Attrs) Values() []string {
	if len(as) == 0 {
		return nil
	}
	res := make([]string, 0, len(as))
	for i := range as {
		res = append(res, as[i].Value)
	}
	return res
}

// ToMap converts the name/value pairs to a map, keyed by XML name
// of each attribute. The result is a map that is not nil but may be empty.
func (as Attrs) ToMap() map[xml.Name]string {
	m := make(map[xml.Name]string)
	for _, a := range as {
		m[a.Name] = a.Value
	}
	return m
}

// ToSimpleMap converts the name/value pairs to a map, keyed only by
// the local name of each attribute. The namespaces of the attributes
// are discarded. The result is a map that is not nil but may be empty.
func (as Attrs) ToSimpleMap() map[string]string {
	m := make(map[string]string)
	for _, a := range as {
		m[a.Name.Local] = a.Value
	}
	return m
}

// Coerce will attempt to convert an attribute value to any type, given a parse function
// that accepts a string. The returned value will only contain the parsed result if the
// attribute existed and could be successfully parsed. Otherwise, it returns the zero
// value for type V. Parameters local and space must exactly match the required attribute.
//
// If error handling is required, use the parse function directly instead.
func Coerce[V any](as Attrs, name NameMatcher, parse func(v string) (V, error)) V {
	return CoerceString(as.Get(name).Value, parse)
}

// CoerceString will attempt to convert an attribute value to any type, given a parse function
// that accepts a string. The returned value will only contain the parsed result if the
// attribute could be successfully parsed. Otherwise, it returns the zero value
// for type V.
//
// If error handling is required, use the parse function directly instead.
func CoerceString[V any](s string, parse func(v string) (V, error)) V {
	v, err := parse(s)
	if err != nil {
		var zero V
		return zero
	}
	return v
}
