package serie

import (
	"fmt"
	"strings"
)

func Object(v ...interface{}) Serie {
	s := New(ObjectValue{}, asObjectValue, compareObjectValue)
	if len(v) > 0 {
		s.Append(v...)
	}
	return s
}

type ObjectValue struct {
	Value map[string]interface{}
	Valid bool
}

func (a ObjectValue) Interface() interface{} {
	if a.Valid {
		return a.Value
	}
	return nil
}

func (a ObjectValue) String() string {
	return fmt.Sprint(a.Value)
}

func asObjectValue(i interface{}) ObjectValue {
	if av, ok := i.(ObjectValue); ok {
		return av
	}

	var a ObjectValue
	if i == nil {
		return a
	}

	if values, ok := i.(map[string]interface{}); ok {
		a.Value = values
		a.Valid = true
	}

	return a
}

func compareObjectValue(a, b ObjectValue) int {
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
