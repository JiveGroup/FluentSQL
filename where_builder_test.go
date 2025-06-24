package fluentsql

import "testing"

// TestWhereBuilderWhereGroup tests the WhereGroup function of WhereBuilder
func TestWhereBuilderWhereGroup(t *testing.T) {
	testCases := map[string]*WhereBuilder{
		"WHERE (first_name = 'John' AND last_name = 'Doe')": WhereInstance().
			WhereGroup(func(whereBuilder WhereBuilder) *WhereBuilder {
				whereBuilder.Where("first_name", Eq, "John").
					Where("last_name", Eq, "Doe")
				return &whereBuilder
			}),
		"WHERE active = true AND (department = 'IT' OR department = 'HR')": WhereInstance().
			Where("active", Eq, true).
			WhereGroup(func(whereBuilder WhereBuilder) *WhereBuilder {
				whereBuilder.Where("department", Eq, "IT").
					WhereOr("department", Eq, "HR")
				return &whereBuilder
			}),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestWhereBuilderWhereCondition tests the WhereCondition function of WhereBuilder
func TestWhereBuilderWhereCondition(t *testing.T) {
	condition1 := Condition{
		Field: "first_name",
		Opt:   Eq,
		Value: "John",
		AndOr: And,
	}
	condition2 := Condition{
		Field: "last_name",
		Opt:   Eq,
		Value: "Doe",
		AndOr: And,
	}
	condition3 := Condition{
		Field: "department",
		Opt:   Eq,
		Value: "IT",
		AndOr: Or,
	}

	testCases := map[string]*WhereBuilder{
		"WHERE first_name = 'John' AND last_name = 'Doe'": WhereInstance().
			WhereCondition(condition1, condition2),
		"WHERE first_name = 'John' AND last_name = 'Doe' OR department = 'IT'": WhereInstance().
			WhereCondition(condition1, condition2, condition3),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestWhereBuilderString tests the String function of WhereBuilder
func TestWhereBuilderString(t *testing.T) {
	testCases := map[string]*WhereBuilder{
		"WHERE id = 1": WhereInstance().
			Where("id", Eq, 1),
		"WHERE id = 1 AND name = 'John'": WhereInstance().
			Where("id", Eq, 1).
			Where("name", Eq, "John"),
		"WHERE id = 1 OR name = 'John'": WhereInstance().
			Where("id", Eq, 1).
			WhereOr("name", Eq, "John"),
	}

	for expected, query := range testCases {
		if query.String() != expected {
			t.Fatalf(`Query %s != %s`, query.String(), expected)
		}
	}
}

// TestWhereBuilderStringArgs tests the StringArgs function of WhereBuilder
func TestWhereBuilderStringArgs(t *testing.T) {
	testCases := map[string]*WhereBuilder{
		"WHERE id = $1": WhereInstance().
			Where("id", Eq, 1),
		"WHERE id = $1 AND name = $2": WhereInstance().
			Where("id", Eq, 1).
			Where("name", Eq, "John"),
		"WHERE id = $1 OR name = $2": WhereInstance().
			Where("id", Eq, 1).
			WhereOr("name", Eq, "John"),
	}

	for expected, query := range testCases {
		var args []any
		sql, args := query.StringArgs([]any{})

		if sql != expected {
			t.Fatalf(`Query %s != %s (%v)`, sql, expected, args)
		}
	}
}

// TestWhereBuilderConditions tests the Conditions function of WhereBuilder
func TestWhereBuilderConditions(t *testing.T) {
	wb := WhereInstance()
	wb.Where("id", Eq, 1)
	wb.Where("name", Eq, "John")

	conditions := wb.Conditions()
	if len(conditions) != 2 {
		t.Fatalf(`Expected 2 conditions, got %d`, len(conditions))
	}

	if conditions[0].Field != "id" || conditions[0].Opt != Eq || conditions[0].Value != 1 || conditions[0].AndOr != And {
		t.Fatalf(`First condition doesn't match expected values`)
	}

	if conditions[1].Field != "name" || conditions[1].Opt != Eq || conditions[1].Value != "John" || conditions[1].AndOr != And {
		t.Fatalf(`Second condition doesn't match expected values`)
	}
}
