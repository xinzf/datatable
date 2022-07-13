package csv_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/xinzf/datatable"

	"github.com/stretchr/testify/assert"
	"github.com/xinzf/datatable/import/csv"
)

func TestImport(t *testing.T) {
	dt, err := csv.Import("csv", "../../test/phone_data.csv",
		csv.HasHeader(true),
		csv.AcceptDate("02/01/06 15:04"),
		csv.AcceptDate("2006-01"),
	)
	assert.NoError(t, err)
	assert.NotNil(t, dt)

	dt.Print(os.Stdout, datatable.PrintMaxRows(24))

	dtc, err := dt.Aggregate(datatable.AggregateBy{Type: datatable.Count, Field: "index"})
	assert.NoError(t, err)
	fmt.Println(dtc)

	groups, err := dt.GroupBy(datatable.GroupBy{
		Name: "year",
		Type: datatable.Int,
		Keyer: func(row datatable.Row) (interface{}, bool) {
			if d, ok := row["date"]; ok {
				if tm, ok := d.(time.Time); ok {
					return tm.Year(), true
				}
			}
			return nil, false
		},
	})
	assert.NoError(t, err)
	out, err := groups.Aggregate(
		datatable.AggregateBy{Type: datatable.Sum, Field: "duration"},
		datatable.AggregateBy{Type: datatable.CountDistinct, Field: "network"},
	)
	assert.NoError(t, err)
	fmt.Println(out)

}
