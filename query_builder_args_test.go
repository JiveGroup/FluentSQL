package fluentsql

import (
	"testing"
)

// TestFieldYearStringArgs tests the StringArgs function of FieldYear
func TestFieldYearStringArgs(t *testing.T) {
	// Save the original dialect to restore it after the test
	originalDialect := defaultDialect
	defer func() {
		defaultDialect = originalDialect
	}()

	// Test with PostgreSQL dialect
	pgDialect := new(PostgreSQLDialect)
	SetDialect(pgDialect)

	fieldYear := FieldYear("hire_date")
	var args []any
	sql, args := fieldYear.StringArgs(args)

	expectedSQL := "DATE_PART('year', $1)"
	if sql != expectedSQL {
		t.Fatalf("Expected SQL %s, got %s", expectedSQL, sql)
	}

	if len(args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(args))
	}

	if args[0] != "hire_date" {
		t.Fatalf("Expected argument to be 'hire_date', got %v", args[0])
	}

	// Test with MySQL dialect
	mysqlDialect := new(MySQLDialect)
	SetDialect(mysqlDialect)

	args = []any{}
	sql, args = fieldYear.StringArgs(args)

	expectedSQL = "YEAR(?)"
	if sql != expectedSQL {
		t.Fatalf("Expected SQL %s, got %s", expectedSQL, sql)
	}

	if len(args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(args))
	}

	if args[0] != "hire_date" {
		t.Fatalf("Expected argument to be 'hire_date', got %v", args[0])
	}

	// Test with SQLite dialect
	sqliteDialect := new(SQLiteDialect)
	SetDialect(sqliteDialect)

	args = []any{}
	sql, args = fieldYear.StringArgs(args)

	expectedSQL = "strftime('%Y', ?)"
	if sql != expectedSQL {
		t.Fatalf("Expected SQL %s, got %s", expectedSQL, sql)
	}

	if len(args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(args))
	}

	if args[0] != "hire_date" {
		t.Fatalf("Expected argument to be 'hire_date', got %v", args[0])
	}
}

// TestWhenCaseStringArgs tests the StringArgs function of WhenCase
func TestWhenCaseStringArgs(t *testing.T) {
	// Test with a simple condition
	whenCase := WhenCase{
		Conditions: 1,
		Value:      "One year",
	}

	var args []any
	sql, args := whenCase.StringArgs(args)

	expectedSQL := "WHEN 1 THEN '?'"
	if sql != expectedSQL {
		t.Fatalf("Expected SQL %s, got %s", expectedSQL, sql)
	}

	if len(args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(args))
	}

	if args[0] != "One year" {
		t.Fatalf("Expected argument to be 'One year', got %v", args[0])
	}

	// Test with multiple conditions
	conditions := []Condition{
		{
			Field: "salary",
			Opt:   Greater,
			Value: 5000,
		},
		{
			Field: "salary",
			Opt:   Lesser,
			Value: 10000,
		},
	}

	whenCase = WhenCase{
		Conditions: conditions,
		Value:      "High salary",
	}

	args = []any{}
	sql, args = whenCase.StringArgs(args)

	expectedSQL = "WHEN salary > $1 AND salary < $2 THEN '?'"
	if sql != expectedSQL {
		t.Fatalf("Expected SQL %s, got %s", expectedSQL, sql)
	}

	if len(args) != 3 {
		t.Fatalf("Expected 3 arguments, got %d", len(args))
	}

	if args[0] != 5000 {
		t.Fatalf("Expected second argument to be 5000, got %v", args[0])
	}

	if args[1] != 10000 {
		t.Fatalf("Expected third argument to be 10000, got %v", args[1])
	}

	if args[2] != "High salary" {
		t.Fatalf("Expected first argument to be 'High salary', got %v", args[2])
	}
}

// TestCaseStringArgs tests the StringArgs function of Case
func TestCaseStringArgs(t *testing.T) {
	// Create a Case with multiple When clauses
	caseTest := new(Case)
	caseTest.Exp = "(2000 - YEAR(hire_date))"
	caseTest.When(1, "1 year")
	caseTest.When(3, "3 years")
	caseTest.Name = "aniversary"

	var args []any
	sql, args := caseTest.StringArgs(args)

	expectedSQL := "CASE (2000 - YEAR(hire_date)) WHEN 1 THEN '?' WHEN 3 THEN '?' END aniversary"
	if sql != expectedSQL {
		t.Fatalf("Expected SQL %s, got %s", expectedSQL, sql)
	}

	if len(args) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(args))
	}

	if args[0] != "1 year" {
		t.Fatalf("Expected first argument to be '1 year', got %v", args[0])
	}

	if args[1] != "3 years" {
		t.Fatalf("Expected second argument to be '3 years', got %v", args[1])
	}

	// Test with conditions
	caseTest = new(Case)
	caseTest.Exp = ""

	var conditionsLow []Condition
	conditionsLow = append(conditionsLow, Condition{
		Field: "salary",
		Opt:   Lesser,
		Value: 3000,
	})

	var conditionsHigh []Condition
	conditionsHigh = append(conditionsHigh, Condition{
		Field: "salary",
		Opt:   Greater,
		Value: 5000,
	})

	caseTest.When(conditionsLow, "Low")
	caseTest.When(conditionsHigh, "High")
	caseTest.Name = "evaluation"

	args = []any{}
	sql, args = caseTest.StringArgs(args)

	expectedSQL = "CASE  WHEN salary < $1 THEN '?' WHEN salary > $3 THEN '?' END evaluation"
	if sql != expectedSQL {
		t.Fatalf("Expected SQL %s, got %s", expectedSQL, sql)
	}

	if len(args) != 4 {
		t.Fatalf("Expected 4 arguments, got %d", len(args))
	}

	if args[0] != 3000 {
		t.Fatalf("Expected first argument to be 3000, got %v", args[0])
	}

	if args[1] != "Low" {
		t.Fatalf("Expected second argument to be 'Low', got %v", args[1])
	}

	if args[2] != 5000 {
		t.Fatalf("Expected fourth argument to be 5000, got %v", args[2])
	}

	if args[3] != "High" {
		t.Fatalf("Expected third argument to be 'High', got %v", args[3])
	}
}
