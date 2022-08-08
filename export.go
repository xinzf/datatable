package datatable

import "time"

// ExportOptions to add options for exporting (like showing hidden columns)
type ExportOptions struct {
	WithHiddenCols bool
	DefaultValue   map[ColumnType]any
}

type ExportOption func(*ExportOptions)

// ExportHidden to show a column when exporting (default false)
func ExportHidden(v bool) ExportOption {
	return func(opts *ExportOptions) {
		opts.WithHiddenCols = v
	}
}

func DefaultValue(tpe ColumnType, val any) ExportOption {
	return func(options *ExportOptions) {
		options.DefaultValue[tpe] = val
	}
}

// newExportOptions to build the ExportOptions in order to acces the parameters
func newExportOptions(opt ...ExportOption) ExportOptions {
	var opts ExportOptions
	opts.DefaultValue = map[ColumnType]any{
		Bool:        nil,
		String:      "",
		Int:         nil,
		Int32:       nil,
		Int64:       nil,
		Float32:     0.0,
		Float64:     0.0,
		Time:        time.Time{},
		Raw:         nil,
		Array:       []interface{}{},
		Object:      map[string]any{},
		ArrayObject: []interface{}{},
	}
	for _, o := range opt {
		o(&opts)
	}
	return opts

}

// ToMap to export the datatable to a json-like struct
//func (t *DataTable) ToMap(opt ...ExportOption) []map[string]interface{} {
//	if t == nil {
//		return nil
//	}
//
//	opts := newExportOptions(opt...)
//	if err := t.evaluateExpressions(); err != nil {
//		panic(err)
//	}
//
//	// visible columns
//	cols := make(map[string]int)
//	for i, col := range t.cols {
//		if opts.WithHiddenCols || col.IsVisible() {
//			cols[col.Name()] = i
//		}
//	}
//
//	rows := make([]map[string]interface{}, 0, t.nrows)
//	for i := 0; i < t.nrows; i++ {
//		r := make(map[string]interface{}, len(cols))
//		for name, pos := range cols {
//			r[name] = t.cols[pos].serie.Get(i)
//		}
//		rows = append(rows, r)
//	}
//	return rows
//}
func (t *DataTable) ToMap(opt ...ExportOption) []map[string]interface{} {
	if t == nil {
		return nil
	}

	opts := newExportOptions(opt...)
	if err := t.evaluateExpressions(); err != nil {
		panic(err)
	}

	type colsDesc struct {
		pos             int
		defaultValue    any
		hasDefaultValue bool
	}
	// visible columns
	//cols := make(map[string]int)
	cols := make(map[string]colsDesc)
	for i, col := range t.cols {
		if opts.WithHiddenCols || col.IsVisible() {
			defaultValue, hasDefaultValue := opts.DefaultValue[col.typ]
			cols[col.name] = colsDesc{
				pos:             i,
				defaultValue:    defaultValue,
				hasDefaultValue: hasDefaultValue,
			}
			//cols[col.Name()] = i
		}
	}

	rows := make([]map[string]interface{}, 0, t.nrows)
	for i := 0; i < t.nrows; i++ {
		r := make(map[string]interface{}, len(cols))
		for name, pos := range cols {
			val := t.cols[pos.pos].serie.Get(i)
			if val == nil && pos.hasDefaultValue {
				val = pos.defaultValue
			}
			r[name] = val
			//r[name] = t.cols[pos].serie.Get(i)
		}
		rows = append(rows, r)
	}
	return rows
}

// ToTable to export the datatable to a csv-like struct
func (t *DataTable) ToTable(opt ...ExportOption) [][]interface{} {
	if t == nil {
		return nil
	}

	opts := newExportOptions(opt...)
	if err := t.evaluateExpressions(); err != nil {
		panic(err)
	}

	rows := make([][]interface{}, 0, t.nrows+1)

	// visible columns
	var headers []interface{}
	var cols []int
	for i, col := range t.cols {
		if opts.WithHiddenCols || col.IsVisible() {
			cols = append(cols, i)
			headers = append(headers, col.Name())
		}
	}

	rows = append(rows, headers)
	for i := 0; i < t.nrows; i++ {
		r := make([]interface{}, 0, len(cols))
		for _, pos := range cols {
			r = append(r, t.cols[pos].serie.Get(i))
		}
		rows = append(rows, r)
	}
	return rows
}

// Schema describes a datatable
type Schema struct {
	Name    string          `json:"name"`
	Columns []SchemaColumn  `json:"cols"`
	Rows    [][]interface{} `json:"rows"`
}

type SchemaColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// ToSchema to export the datatable to a schema struct
func (t *DataTable) ToSchema(opt ...ExportOption) *Schema {
	if t == nil {
		return nil
	}

	opts := newExportOptions(opt...)
	if err := t.evaluateExpressions(); err != nil {
		panic(err)
	}

	schema := &Schema{
		Name: t.name,
		Rows: make([][]interface{}, 0, t.nrows),
	}

	// visible columns
	var cols []int
	for i, col := range t.cols {
		if opts.WithHiddenCols || col.IsVisible() {
			cols = append(cols, i)
			schema.Columns = append(schema.Columns, SchemaColumn{Type: col.UnderlyingType().Name(), Name: col.Name()})
		}
	}

	for i := 0; i < t.nrows; i++ {
		r := make([]interface{}, 0, len(cols))
		for _, pos := range cols {
			r = append(r, t.cols[pos].serie.Get(i))
		}
		schema.Rows = append(schema.Rows, r)
	}

	return schema
}
