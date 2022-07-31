package serie

import (
	"fmt"
	"strings"
)

func Array(v ...interface{}) Serie {
	s := New(ArrayValue{}, asArrayValue, compareArrayValue)
	if len(v) > 0 {
		s.Append(v...)
	}
	return s
}

type ArrayValue struct {
	Value []interface{}
	Valid bool
}

func (a ArrayValue) Interface() interface{} {
	if a.Valid {
		return a.Value
	}
	return nil
}

func (a ArrayValue) String() string {
	return fmt.Sprint(a.Value)
}

func asArrayValue(i interface{}) ArrayValue {
	if av, ok := i.(ArrayValue); ok {
		return av
	}

	var a ArrayValue
	if i == nil {
		return a
	}

	if values, ok := i.(*[]interface{}); ok {
		a.Value = *values
		a.Valid = true
	}

	return a
}

func compareArrayValue(a, b ArrayValue) int {
	if !b.Valid {
		if !a.Valid {
			return Eq
		}
		return Gt
	}
	if !a.Valid {
		return Lt
	}

	return strings.Compare(a.String(), b.String())
}
