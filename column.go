package datatable

import (
	"reflect"
	"strings"

	"github.com/datasweet/expr"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/xinzf/datatable/serie"
)

// ColumnType defines the valid column type in datatable
type ColumnType string

const (
	Bool   ColumnType = "bool"
	String ColumnType = "string"
	Int    ColumnType = "int"
	// Int8     ColumnType = "int8"
	// Int16    ColumnType = "int16"
	Int32 ColumnType = "int32"
	Int64 ColumnType = "int64"
	// Uint  ColumnType = "uint"
	// Uint8     ColumnType = "uint8"
	// Uint16    ColumnType = "uint16"
	// Uint32    ColumnType = "uint32"
	// Uint64    ColumnType = "uint64"
	Float32     ColumnType = "float32"
	Float64     ColumnType = "float64"
	Time        ColumnType = "time"
	Raw         ColumnType = "raw"
	Array       ColumnType = "array"
	Object      ColumnType = "object"
	ArrayObject ColumnType = "arrayObject"
)

// ColumnOptions describes options to be apply on a column
type ColumnOptions struct {
	Hidden      bool
	Expr        string
	Values      []interface{}
	TimeFormats []string
	Label       string
	Attrs       map[string]interface{}
}

// ColumnOption sets column options
type ColumnOption func(opts *ColumnOptions)

func ColumnLabel(label string) ColumnOption {
	return func(opts *ColumnOptions) {
		opts.Label = label
	}
}

func ColumnAttrs(attrs map[string]interface{}) ColumnOption {
	return func(opts *ColumnOptions) {
		opts.Attrs = attrs
	}
}

// ColumnHidden sets the visibility
func ColumnHidden(v bool) ColumnOption {
	return func(opts *ColumnOptions) {
		opts.Hidden = v
	}
}

// Expr sets the expr for the column
// <!> Incompatible with ColumnValues
func Expr(v string) ColumnOption {
	return func(opts *ColumnOptions) {
		opts.Expr = v
	}
}

// Values fills the column with the values
// <!> Incompatible with ColumnExpr
func Values(v ...interface{}) ColumnOption {
	return func(opts *ColumnOptions) {
		opts.Values = v
	}
}

// TimeFormats sets the valid time formats.
// <!> Only for Time Column
func TimeFormats(v ...string) ColumnOption {
	return func(opts *ColumnOptions) {
		opts.TimeFormats = append(opts.TimeFormats, v...)
	}
}

// ColumnSerier to create a serie from column options
type ColumnSerier func(ColumnOptions) serie.Serie

// ctypes is our column type registry
var ctypes map[ColumnType]ColumnSerier

func init() {
	ctypes = make(map[ColumnType]ColumnSerier)
	_ = RegisterColumnType(Bool, func(opts ColumnOptions) serie.Serie {
		return serie.BoolN(opts.Values...)
	})
	_ = RegisterColumnType(String, func(opts ColumnOptions) serie.Serie {
		return serie.StringN(opts.Values...)
	})
	_ = RegisterColumnType(Int, func(opts ColumnOptions) serie.Serie {
		return serie.IntN(opts.Values...)
	})
	_ = RegisterColumnType(Int32, func(opts ColumnOptions) serie.Serie {
		return serie.Int32N(opts.Values...)
	})
	_ = RegisterColumnType(Int64, func(opts ColumnOptions) serie.Serie {
		return serie.Int64N(opts.Values...)
	})
	_ = RegisterColumnType(Float32, func(opts ColumnOptions) serie.Serie {
		return serie.Float32N(opts.Values...)
	})
	_ = RegisterColumnType(Float64, func(opts ColumnOptions) serie.Serie {
		return serie.Float64N(opts.Values...)
	})
	_ = RegisterColumnType(Time, func(opts ColumnOptions) serie.Serie {
		sr := serie.TimeN(opts.TimeFormats...)
		if len(opts.Values) > 0 {
			sr.Append(opts.Values...)
		}
		return sr
	})
	_ = RegisterColumnType(Raw, func(opts ColumnOptions) serie.Serie {
		return serie.Raw(opts.Values...)
	})
	_ = RegisterColumnType(Object, func(options ColumnOptions) serie.Serie {
		return serie.Object(options.Values...)
	})
	_ = RegisterColumnType(Array, func(options ColumnOptions) serie.Serie {
		items := make([]*[]interface{}, 0)
		for _, value := range options.Values {
			arr := make([]interface{}, 0)
			bytes, err := jsoniter.Marshal(value)
			if err == nil {
				_ = jsoniter.Unmarshal(bytes, &arr)
			}

			items = append(items, &arr)
		}
		return serie.Array(items)
	})
}

// RegisterColumnType to extends the known type
func RegisterColumnType(name ColumnType, serier ColumnSerier) error {
	name = ColumnType(strings.TrimSpace(string(name)))
	if len(name) == 0 {
		return ErrEmptyName
	}
	if serier == nil {
		return ErrNilFactory
	}
	if _, ok := ctypes[name]; ok {
		err := errors.Errorf("type '%s' already exists", name)
		return errors.Wrap(err, ErrTypeAlreadyExists.Error())
	}
	ctypes[name] = serier
	return nil
}

// ColumnTypes to list all column type
func ColumnTypes() []ColumnType {
	ctyp := make([]ColumnType, 0, len(ctypes))
	for k := range ctypes {
		ctyp = append(ctyp, k)
	}
	return ctyp
}

// newColumnSerie to create a serie from a known type
func newColumnSerie(ctyp ColumnType, options ColumnOptions) (serie.Serie, error) {
	if s, ok := ctypes[ctyp]; ok {
		return s(options), nil
	}
	err := errors.Errorf("unknown column type '%s'", ctyp)
	return nil, errors.Wrap(err, ErrUnknownColumnType.Error())
}

// Column describes a column in our datatable
type Column interface {
	Name() string
	Type() ColumnType
	Label() string
	Attrs() map[string]interface{}
	SetLabel(label string) Column
	SetAttrs(attrs map[string]interface{}) Column
	Clone() Column
	UnderlyingType() reflect.Type
	IsVisible() bool
	IsComputed() bool
	Serie() serie.Serie
	//Clone(includeValues bool) Column
}

type column struct {
	name     string
	typ      ColumnType
	hidden   bool
	label    string
	attrs    map[string]interface{}
	formulae string
	expr     expr.Node
	serie    serie.Serie
}

func (c *column) Serie() serie.Serie {
	return c.serie
}

func (c *column) Name() string {
	return c.name
}

func (c *column) Label() string {
	if c.label == "" {
		return c.name
	}
	return c.label
}

func (c *column) Attrs() map[string]interface{} {
	if c.attrs == nil {
		return map[string]interface{}{}
	}
	return c.attrs
}

func (c *column) SetLabel(label string) Column {
	c.label = label
	return c
}

func (c *column) Clone() Column {
	return &column{
		name:     c.name,
		typ:      c.typ,
		hidden:   c.hidden,
		label:    c.label,
		attrs:    c.attrs,
		formulae: c.formulae,
		expr:     c.expr,
		serie:    c.serie.Copy(),
	}
}

func (c *column) SetAttrs(attrs map[string]interface{}) Column {
	c.attrs = attrs
	return c
}

func (c *column) Type() ColumnType {
	return c.typ
}

func (c *column) UnderlyingType() reflect.Type {
	return c.serie.Type()
}

func (c *column) IsVisible() bool {
	return !c.hidden
}

func (c *column) IsComputed() bool {
	return len(c.formulae) > 0
}

func (c *column) emptyCopy() *column {
	cpy := &column{
		name:     c.name,
		typ:      c.typ,
		label:    c.label,
		attrs:    c.attrs,
		hidden:   c.hidden,
		formulae: c.formulae,
		serie:    c.serie.EmptyCopy(),
	}
	if len(cpy.formulae) > 0 {
		if parsed, err := expr.Parse(cpy.formulae); err == nil {
			cpy.expr = parsed
		}
	}
	return cpy
}

func (c *column) copy() *column {
	cpy := &column{
		name:     c.name,
		typ:      c.typ,
		label:    c.label,
		attrs:    c.attrs,
		hidden:   c.hidden,
		formulae: c.formulae,
		serie:    c.serie.Copy(),
	}
	if len(cpy.formulae) > 0 {
		if parsed, err := expr.Parse(cpy.formulae); err == nil {
			cpy.expr = parsed
		}
	}
	return cpy
}
