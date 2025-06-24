package fluentsql

import (
	"testing"
)

// TestDeleteTable
func TestDeleteTable(t *testing.T) {
	testCases := map[string]*DeleteBuilder{
		"DELETE FROM customers WHERE contact_name = 'Alfred Schmidt' AND city = 'Frankfurt' AND customer_id = 1": DeleteInstance().
			Delete("customers").
			Where("contact_name", Eq, "Alfred Schmidt").
			Where("city", Eq, "Frankfurt").
			Where("customer_id", Eq, 1),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestDeleteTableArgs
func TestDeleteTableArgs(t *testing.T) {
	testCases := map[string]*DeleteBuilder{
		"DELETE FROM customers WHERE contact_name = $1 AND city = $2 AND customer_id = $3": DeleteInstance().
			Delete("customers").
			Where("contact_name", Eq, "Alfred Schmidt").
			Where("city", Eq, "Frankfurt").
			Where("customer_id", Eq, 1),
	}

	for expected, query := range testCases {
		var sql string
		var args []any

		sql, args, _ = query.Sql()

		if sql != expected {
			t.Fatalf(`Query %s != %s (%v)`, sql, expected, args)
		}
	}
}

// TestDeleteWhereOr tests the WhereOr function of DeleteBuilder
func TestDeleteWhereOr(t *testing.T) {
	testCases := map[string]*DeleteBuilder{
		"DELETE FROM customers WHERE contact_name = 'Alfred Schmidt' OR city = 'Frankfurt' OR customer_id = 1": DeleteInstance().
			Delete("customers").
			WhereOr("contact_name", Eq, "Alfred Schmidt").
			WhereOr("city", Eq, "Frankfurt").
			WhereOr("customer_id", Eq, 1),
		"DELETE FROM customers WHERE contact_name = 'Alfred Schmidt' AND city = 'Frankfurt' OR customer_id = 1": DeleteInstance().
			Delete("customers").
			Where("contact_name", Eq, "Alfred Schmidt").
			Where("city", Eq, "Frankfurt").
			WhereOr("customer_id", Eq, 1),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestDeleteWhereGroup tests the WhereGroup function of DeleteBuilder
func TestDeleteWhereGroup(t *testing.T) {
	testCases := map[string]*DeleteBuilder{
		"DELETE FROM employees WHERE (first_name = 'John' AND last_name = 'Doe')": DeleteInstance().
			Delete("employees").
			WhereGroup(func(whereBuilder WhereBuilder) *WhereBuilder {
				whereBuilder.Where("first_name", Eq, "John").
					Where("last_name", Eq, "Doe")
				return &whereBuilder
			}),
		"DELETE FROM employees WHERE contact_name = 'Alfred Schmidt' AND (city = 'Frankfurt' OR country = 'Germany')": DeleteInstance().
			Delete("employees").
			Where("contact_name", Eq, "Alfred Schmidt").
			WhereGroup(func(whereBuilder WhereBuilder) *WhereBuilder {
				whereBuilder.Where("city", Eq, "Frankfurt").
					WhereOr("country", Eq, "Germany")
				return &whereBuilder
			}),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestDeleteWhereCondition tests the WhereCondition function of DeleteBuilder
func TestDeleteWhereCondition(t *testing.T) {
	condition1 := Condition{
		Field: "contact_name",
		Opt:   Eq,
		Value: "Alfred Schmidt",
		AndOr: And,
	}
	condition2 := Condition{
		Field: "city",
		Opt:   Eq,
		Value: "Frankfurt",
		AndOr: And,
	}
	condition3 := Condition{
		Field: "customer_id",
		Opt:   Eq,
		Value: 1,
		AndOr: Or,
	}

	testCases := map[string]*DeleteBuilder{
		"DELETE FROM customers WHERE contact_name = 'Alfred Schmidt' AND city = 'Frankfurt'": DeleteInstance().
			Delete("customers").
			WhereCondition(condition1, condition2),
		"DELETE FROM customers WHERE contact_name = 'Alfred Schmidt' AND city = 'Frankfurt' OR customer_id = 1": DeleteInstance().
			Delete("customers").
			WhereCondition(condition1, condition2, condition3),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}
