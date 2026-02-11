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

// Get gets the attribute that has matching name (local and space).
// No wildcards are supported. If there is no match, a zero xml.Attr
// is returned.
func (as Attrs) Get(local, space string) xml.Attr {
	for i := range as {
		if local == as[i].Name.Local && space == as[i].Name.Space {
			return as[i]
		}
	}
	return xml.Attr{}
}

// WithName filters the attributes and returns only those that have
// matching local and space. Either of the local and space parameters
// can be "*", which is a wildcard.
func (as Attrs) WithName(local, space string) Attrs {
	if len(as) == 0 {
		return nil
	}
	res := make(Attrs, 0, len(as))
	for i := range as {
		if (local == "*" || local == as[i].Name.Local) &&
			(space == "*" || space == as[i].Name.Space) {
			res = append(res, as[i])
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
func Coerce[V any](as Attrs, local, space string, parse func(v string) (V, error)) V {
	v, err := parse(as.Get(local, space).Value)
	if err != nil {
		var zero V
		return zero
	}
	return v
}
