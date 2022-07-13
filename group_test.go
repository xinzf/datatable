package datatable_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/xinzf/datatable"
)

import "testing"

func TestGroup(t *testing.T) {
	dt := datatable.New("group_test")
	_ = dt.AddColumn("UUID", datatable.String, datatable.Values("a", "b", "c", "d"))
	_ = dt.AddColumn("uid", datatable.Int, datatable.Values(1, 1, 2, 2), datatable.ColumnLabel("用户ID"))
	_ = dt.AddColumn("name", datatable.String, datatable.Values("向志", "向志", "刘志楠", "刘志楠"), datatable.ColumnLabel("姓名"))

	groups, err := dt.GroupBy(datatable.GroupBy{
		Name: "uid",
		Type: datatable.Float32,
		Keyer: func(row datatable.Row) (interface{}, bool) {
			return row.Get("uid"), true
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, groups)
	newDt, err := groups.Aggregate(datatable.AggregateBy{
		Type:  datatable.GroupAny,
		Field: "UUID",
	})
	assert.NoError(t, err)
	assert.NotNil(t, newDt)
	newDt.Preview()
}
