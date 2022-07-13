package datatable

import (
	"fmt"
	"github.com/liushuochen/gotable"
	"github.com/liushuochen/gotable/cell"
	"io"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// PrintOptions to control the printer
type PrintOptions struct {
	ColumnName    bool
	ColumnLabel   bool
	ColumnType    bool
	RowNumber     bool
	MaxRows       int
	SelectColumns []string
}

type PrintOption func(opts *PrintOptions)

func PrintColumnLabel(v bool) PrintOption {
	return func(opts *PrintOptions) {
		opts.ColumnLabel = v
	}
}

func SelectColumns(name ...string) PrintOption {
	return func(opts *PrintOptions) {
		opts.SelectColumns = name
	}
}

func PrintColumnName(v bool) PrintOption {
	return func(opts *PrintOptions) {
		opts.ColumnName = v
	}
}

func PrintColumnType(v bool) PrintOption {
	return func(opts *PrintOptions) {
		opts.ColumnType = v
	}
}

func PrintRowNumber(v bool) PrintOption {
	return func(opts *PrintOptions) {
		opts.RowNumber = v
	}
}

func PrintMaxRows(v int) PrintOption {
	return func(opts *PrintOptions) {
		opts.MaxRows = v
	}
}

func (this *DataTable) Preview(opt ...PrintOption) {
	tCopy := this.Copy()
	options := PrintOptions{
		ColumnName:    true,
		ColumnLabel:   true,
		ColumnType:    true,
		RowNumber:     true,
		MaxRows:       100,
		SelectColumns: nil,
	}

	for _, o := range opt {
		o(&options)
	}

	//if writer == nil {
	//	writer = os.Stdout
	//}

	headers := make([]string, 0)
	{
		if options.RowNumber {
			headers = append(headers, "#")
		}
	}

	if options.ColumnName || options.ColumnLabel || options.ColumnType {
		for _, col := range tCopy.cols {
			if options.SelectColumns != nil {
				found := false
				for _, selectColumn := range options.SelectColumns {
					if col.name == selectColumn {
						found = true
						break
					}
				}
				if !found {
					col.hidden = true
				}
			}

			if !col.IsVisible() {
				continue
			}

			var h []string
			if options.ColumnName {
				h = append(h, col.Name())
			}
			if options.ColumnLabel {
				h = append(h, fmt.Sprintf("%s", col.Label()))
			}
			if options.ColumnType {
				h = append(h, fmt.Sprintf("%s", col.serie.Type().Name()))
			}
			headers = append(headers, strings.Join(h, "、"))
		}
	}

	tb, _ := gotable.Create(headers...)

	for _, header := range headers {
		tb.Align(header, cell.AlignCenter)
		tb.SetColumnColor(header, gotable.Underline, gotable.Write, gotable.NoneBackground)
	}

	numRows := options.MaxRows
	if tCopy.NumRows() < numRows {
		numRows = tCopy.NumRows()
	}

	if numRows < 1 {
		return
	}

	for i, record := range tCopy.Head(numRows).Records() {
		if options.RowNumber {
			record = append([]string{fmt.Sprintf("#%d", i)}, record...)
		}
		_ = tb.AddRow(record)
	}

	fmt.Printf("\nTable: %s, NumRows: %d, NumCols: %d\n", tCopy.name, this.NumRows(), this.NumCols())
	fmt.Println(tb)
}

// Print the tables with options
func (t *DataTable) Print(writer io.Writer, opt ...PrintOption) {
	options := PrintOptions{
		ColumnName:  true,
		ColumnLabel: true,
		ColumnType:  true,
		RowNumber:   true,
		MaxRows:     100,
	}

	for _, o := range opt {
		o(&options)
	}

	if writer == nil {
		writer = os.Stdout
	}

	tw := tablewriter.NewWriter(writer)
	tw.SetAutoWrapText(false)
	tw.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	tw.SetAlignment(tablewriter.ALIGN_LEFT)
	tw.SetCenterSeparator("")
	tw.SetColumnSeparator("")
	tw.SetRowSeparator("")
	tw.SetHeaderLine(false)
	tw.SetBorder(false)
	tw.SetTablePadding("\t")
	tw.SetNoWhiteSpace(true)

	if options.ColumnName || options.ColumnType {
		headers := make([]string, 0, len(t.cols))

		for _, col := range t.cols {
			if !col.IsVisible() {
				continue
			}
			var h []string
			if options.ColumnName {
				h = append(h, col.Name())
			}
			if options.ColumnType {
				h = append(h, fmt.Sprintf("<%s>", col.serie.Type().Name()))
			}
			headers = append(headers, strings.Join(h, " "))
		}
		tw.SetHeader(headers)
	}

	if options.MaxRows > 1 && options.MaxRows <= t.NumRows() {
		mr := options.MaxRows / 2
		tw.AppendBulk(t.Head(mr).Records())
		seps := make([]string, 0, len(t.cols))
		for _, col := range t.cols {
			if !col.IsVisible() {
				continue
			}
			seps = append(seps, "...")
		}
		tw.Append(seps)
		tw.AppendBulk(t.Tail(mr).Records())
	} else {
		tw.AppendBulk(t.Records())
	}

	tw.Render()
}
