package datatable_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xinzf/datatable"
)

// https://sql.sh/cours/jointures/inner-join
func sampleForJoin() (*datatable.DataTable, *datatable.DataTable) {
	customers := datatable.New("Customers")
	customers.AddColumn("id", datatable.Int)
	customers.AddColumn("prenom", datatable.String)
	customers.AddColumn("nom", datatable.String)
	customers.AddColumn("email", datatable.String)
	customers.AddColumn("ville", datatable.String)
	// customers.AddColumn("concat", datatable.String, datatable.Expr("CONCAT(`prenom`,`nom`)"))
	customers.AppendRow(1, "Aimée", "Marechal", "aime.marechal@example.com", "Paris")
	customers.AppendRow(2, "Esmée", "Lefort", "esmee.lefort@example.com", "Lyon")
	customers.AppendRow(3, "Marine", "Prevost", "m.prevost@example.com", "Lille")
	customers.AppendRow(4, "Luc", "Rolland", "lucrolland@example.com", "Marseille")

	orders := datatable.New("Orders")
	orders.AddColumn("user_id", datatable.Int, datatable.Values(1, 1, 2, 3, 5))
	orders.AddColumn("date_achat", datatable.Time, datatable.Values("2013-01-23", "2013-02-14", "2013-02-17", "2013-02-21", "2013-03-02"))
	orders.AddColumn("num_facture", datatable.String, datatable.Values("A00103", "A00104", "A00105", "A00106", "A00107"))
	orders.AddColumn("prix_total", datatable.Float64, datatable.Values(203.14, 124.00, 149.45, 235.35, 47.58))

	return customers, orders
}

func TestJoinOn(t *testing.T) {
	on := datatable.On("[customers].[id]")
	assert.NotNil(t, on)
	assert.Len(t, on, 1)
	assert.Equal(t, "customers", on[0].Table)
	assert.Equal(t, "id", on[0].Field)

	on = datatable.On("[id]")
	assert.NotNil(t, on)
	assert.Len(t, on, 1)
	assert.Equal(t, "*", on[0].Table)
	assert.Equal(t, "id", on[0].Field)

	on = datatable.On("id")
	assert.NotNil(t, on)
	assert.Len(t, on, 1)
	assert.Equal(t, "*", on[0].Table)
	assert.Equal(t, "id", on[0].Field)

	on = datatable.On("customers.[id]")
	assert.NotNil(t, on)
	assert.Len(t, on, 1)
	assert.Equal(t, "*", on[0].Table)
	assert.Equal(t, "customers.[id]", on[0].Field)
}

func TestInnerJoin(t *testing.T) {
	customers, orders := sampleForJoin()
	customers.AddColumn("concat", datatable.String, datatable.Expr("concat(`prenom`, `nom`)"))
	dt, err := customers.InnerJoin(orders, datatable.On("[Customers].[id]", "[Orders].[user_id]"))
	assert.NoError(t, err)
	assert.NotNil(t, dt)
	checkTable(t, dt,
		"id", "prenom", "nom", "email", "ville", "concat", "date_achat", "num_facture", "prix_total",
		1, "Aimée", "Marechal", "aime.marechal@example.com", "Paris", "AiméeMarechal", time.Date(2013, time.January, 23, 0, 0, 0, 0, time.UTC), "A00103", 203.14,
		1, "Aimée", "Marechal", "aime.marechal@example.com", "Paris", "AiméeMarechal", time.Date(2013, time.February, 14, 0, 0, 0, 0, time.UTC), "A00104", 124.00,
		2, "Esmée", "Lefort", "esmee.lefort@example.com", "Lyon", "EsméeLefort", time.Date(2013, time.February, 17, 0, 0, 0, 0, time.UTC), "A00105", 149.45,
		3, "Marine", "Prevost", "m.prevost@example.com", "Lille", "MarinePrevost", time.Date(2013, time.February, 21, 0, 0, 0, 0, time.UTC), "A00106", 235.35,
	)
}

func TestLeftJoin(t *testing.T) {
	customers, orders := sampleForJoin()

	dt, err := customers.LeftJoin(orders, datatable.On("[Customers].[id]", "[Orders].[user_id]"))
	assert.NoError(t, err)
	assert.NotNil(t, dt)

	checkTable(t, dt,
		"id", "prenom", "nom", "email", "ville", "date_achat", "num_facture", "prix_total",
		1, "Aimée", "Marechal", "aime.marechal@example.com", "Paris", time.Date(2013, time.January, 23, 0, 0, 0, 0, time.UTC), "A00103", 203.14,
		1, "Aimée", "Marechal", "aime.marechal@example.com", "Paris", time.Date(2013, time.February, 14, 0, 0, 0, 0, time.UTC), "A00104", 124.00,
		2, "Esmée", "Lefort", "esmee.lefort@example.com", "Lyon", time.Date(2013, time.February, 17, 0, 0, 0, 0, time.UTC), "A00105", 149.45,
		3, "Marine", "Prevost", "m.prevost@example.com", "Lille", time.Date(2013, time.February, 21, 0, 0, 0, 0, time.UTC), "A00106", 235.35,
		4, "Luc", "Rolland", "lucrolland@example.com", "Marseille", nil, nil, nil,
	)
}

func TestRightJoin(t *testing.T) {
	customers, orders := sampleForJoin()

	dt, err := customers.RightJoin(orders, datatable.On("[Customers].[id]", "[Orders].[user_id]"))
	assert.NoError(t, err)
	assert.NotNil(t, dt)

	checkTable(t, dt,
		"id", "prenom", "nom", "email", "ville", "date_achat", "num_facture", "prix_total",
		1, "Aimée", "Marechal", "aime.marechal@example.com", "Paris", time.Date(2013, time.January, 23, 0, 0, 0, 0, time.UTC), "A00103", 203.14,
		1, "Aimée", "Marechal", "aime.marechal@example.com", "Paris", time.Date(2013, time.February, 14, 0, 0, 0, 0, time.UTC), "A00104", 124.00,
		2, "Esmée", "Lefort", "esmee.lefort@example.com", "Lyon", time.Date(2013, time.February, 17, 0, 0, 0, 0, time.UTC), "A00105", 149.45,
		3, "Marine", "Prevost", "m.prevost@example.com", "Lille", time.Date(2013, time.February, 21, 0, 0, 0, 0, time.UTC), "A00106", 235.35,
		nil, nil, nil, nil, nil, time.Date(2013, time.March, 2, 0, 0, 0, 0, time.UTC), "A00107", 47.58,
	)
}

func TestOuterJoin(t *testing.T) {
	customers, orders := sampleForJoin()

	dt, err := customers.OuterJoin(orders, datatable.On("[Customers].[id]", "[Orders].[user_id]"))
	assert.NoError(t, err)
	assert.NotNil(t, dt)

	checkTable(t, dt,
		"id", "prenom", "nom", "email", "ville", "date_achat", "num_facture", "prix_total",
		1, "Aimée", "Marechal", "aime.marechal@example.com", "Paris", time.Date(2013, time.January, 23, 0, 0, 0, 0, time.UTC), "A00103", 203.14,
		1, "Aimée", "Marechal", "aime.marechal@example.com", "Paris", time.Date(2013, time.February, 14, 0, 0, 0, 0, time.UTC), "A00104", 124.00,
		2, "Esmée", "Lefort", "esmee.lefort@example.com", "Lyon", time.Date(2013, time.February, 17, 0, 0, 0, 0, time.UTC), "A00105", 149.45,
		3, "Marine", "Prevost", "m.prevost@example.com", "Lille", time.Date(2013, time.February, 21, 0, 0, 0, 0, time.UTC), "A00106", 235.35,
		4, "Luc", "Rolland", "lucrolland@example.com", "Marseille", nil, nil, nil,
		nil, nil, nil, nil, nil, time.Date(2013, time.March, 2, 0, 0, 0, 0, time.UTC), "A00107", 47.58,
	)
}

func TestInnerJoinWithExprOnHidden(t *testing.T) {
	customers, orders := sampleForJoin()
	customers.AddColumn("id2", datatable.Int, datatable.Expr("`id`+100"))
	orders.AddColumn("user_id2", datatable.Int, datatable.Expr("`user_id`+100"))
	customers.HideColumn("id")
	dt, err := customers.InnerJoin(orders, datatable.On("[Customers].[id2]", "[Orders].[user_id2]"))
	assert.NoError(t, err)
	assert.NotNil(t, dt)

	fmt.Println(dt)

	checkTable(t, dt,
		"prenom", "nom", "email", "ville", "id2", "user_id", "date_achat", "num_facture", "prix_total",
		"Aimée", "Marechal", "aime.marechal@example.com", "Paris", 101, 1, time.Date(2013, time.January, 23, 0, 0, 0, 0, time.UTC), "A00103", 203.14,
		"Aimée", "Marechal", "aime.marechal@example.com", "Paris", 101, 1, time.Date(2013, time.February, 14, 0, 0, 0, 0, time.UTC), "A00104", 124.00,
		"Esmée", "Lefort", "esmee.lefort@example.com", "Lyon", 102, 2, time.Date(2013, time.February, 17, 0, 0, 0, 0, time.UTC), "A00105", 149.45,
		"Marine", "Prevost", "m.prevost@example.com", "Lille", 103, 3, time.Date(2013, time.February, 21, 0, 0, 0, 0, time.UTC), "A00106", 235.35,
	)
}

func TestLeftJoinWithExpr(t *testing.T) {
	customers, orders := sampleForJoin()
	customers.AddColumn("id2", datatable.Int, datatable.Expr("`id`+100"))
	orders.AddColumn("user_id2", datatable.Int, datatable.Expr("`user_id`+100"))
	dt, err := customers.LeftJoin(orders, datatable.On("[Customers].[id2]", "[Orders].[user_id2]"))
	assert.NoError(t, err)
	assert.NotNil(t, dt)

	fmt.Println(dt)

	checkTable(t, dt,
		"id", "prenom", "nom", "email", "ville", "id2", "user_id", "date_achat", "num_facture", "prix_total",
		1, "Aimée", "Marechal", "aime.marechal@example.com", "Paris", 101, 1, time.Date(2013, time.January, 23, 0, 0, 0, 0, time.UTC), "A00103", 203.14,
		1, "Aimée", "Marechal", "aime.marechal@example.com", "Paris", 101, 1, time.Date(2013, time.February, 14, 0, 0, 0, 0, time.UTC), "A00104", 124.00,
		2, "Esmée", "Lefort", "esmee.lefort@example.com", "Lyon", 102, 2, time.Date(2013, time.February, 17, 0, 0, 0, 0, time.UTC), "A00105", 149.45,
		3, "Marine", "Prevost", "m.prevost@example.com", "Lille", 103, 3, time.Date(2013, time.February, 21, 0, 0, 0, 0, time.UTC), "A00106", 235.35,
		4, "Luc", "Rolland", "lucrolland@example.com", "Marseille", 104, nil, nil, nil, nil,
	)
}

func TestRightJoinWithExpr(t *testing.T) {
	customers, orders := sampleForJoin()
	customers.AddColumn("id2", datatable.Int, datatable.Expr("`id`+100"))
	orders.AddColumn("user_id2", datatable.Int, datatable.Expr("`user_id`+100"))
	dt, err := customers.RightJoin(orders, datatable.On("[Customers].[id2]", "[Orders].[user_id2]"))
	assert.NoError(t, err)
	assert.NotNil(t, dt)

	fmt.Println(dt)

	checkTable(t, dt,
		"id", "prenom", "nom", "email", "ville", "id2", "user_id", "date_achat", "num_facture", "prix_total",
		1, "Aimée", "Marechal", "aime.marechal@example.com", "Paris", 101, 1, time.Date(2013, time.January, 23, 0, 0, 0, 0, time.UTC), "A00103", 203.14,
		1, "Aimée", "Marechal", "aime.marechal@example.com", "Paris", 101, 1, time.Date(2013, time.February, 14, 0, 0, 0, 0, time.UTC), "A00104", 124.00,
		2, "Esmée", "Lefort", "esmee.lefort@example.com", "Lyon", 102, 2, time.Date(2013, time.February, 17, 0, 0, 0, 0, time.UTC), "A00105", 149.45,
		3, "Marine", "Prevost", "m.prevost@example.com", "Lille", 103, 3, time.Date(2013, time.February, 21, 0, 0, 0, 0, time.UTC), "A00106", 235.35,
		nil, nil, nil, nil, nil, nil, 5, time.Date(2013, time.March, 2, 0, 0, 0, 0, time.UTC), "A00107", 47.58,
	)
}

func TestJoinWithColumnName(t *testing.T) {
	customers, orders := sampleForJoin()
	assert.NoError(t, customers.RenameColumn("id", "ClientID"))

	dt, err := customers.InnerJoin(orders, datatable.On("[Customers].[ClientID]", "[Orders].[user_id]"))
	assert.NoError(t, err)
	assert.NotNil(t, dt)

	checkTable(t, dt,
		"ClientID", "prenom", "nom", "email", "ville", "date_achat", "num_facture", "prix_total",
		1, "Aimée", "Marechal", "aime.marechal@example.com", "Paris", time.Date(2013, time.January, 23, 0, 0, 0, 0, time.UTC), "A00103", 203.14,
		1, "Aimée", "Marechal", "aime.marechal@example.com", "Paris", time.Date(2013, time.February, 14, 0, 0, 0, 0, time.UTC), "A00104", 124.00,
		2, "Esmée", "Lefort", "esmee.lefort@example.com", "Lyon", time.Date(2013, time.February, 17, 0, 0, 0, 0, time.UTC), "A00105", 149.45,
		3, "Marine", "Prevost", "m.prevost@example.com", "Lille", time.Date(2013, time.February, 21, 0, 0, 0, 0, time.UTC), "A00106", 235.35,
	)
}
