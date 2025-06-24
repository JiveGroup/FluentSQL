package fluentsql

import (
	"testing"
)

// TestHavingBasic tests basic HAVING clause functionality
func TestHavingBasic(t *testing.T) {
	testCases := map[string]Condition{
		"HAVING salary > 1400": {
			Field: "salary",
			Opt:   Greater,
			Value: 1400,
		},
		"HAVING department_id = 5": {
			Field: "department_id",
			Opt:   Eq,
			Value: 5,
		},
		"HAVING COUNT(employee_id) > 10": {
			Field: "COUNT(employee_id)",
			Opt:   Greater,
			Value: 10,
		},
		"HAVING AVG(salary) >= 5000": {
			Field: "AVG(salary)",
			Opt:   GrEq,
			Value: 5000,
		},
	}

	for expected, condition := range testCases {
		havingTest := new(Having)
		havingTest.Append(condition)

		if havingTest.String() != expected {
			t.Fatalf(`Query %s != %s`, havingTest.String(), expected)
		}
	}
}

// TestHavingOr tests OR conditions in HAVING clause (covers lines 28-35 in having.go)
func TestHavingOr(t *testing.T) {
	var conditions []Condition

	// Test case 1: Simple OR condition
	conditions = append(conditions, Condition{
		Field: "COUNT(employee_id)",
		Opt:   Greater,
		Value: 10,
	}, Condition{
		Field: "AVG(salary)",
		Opt:   Greater,
		Value: 5000,
		AndOr: Or,
	})

	testCases := map[string][]Condition{
		"HAVING COUNT(employee_id) > 10 OR AVG(salary) > 5000": conditions,
	}

	for expected, condition := range testCases {
		havingTest := new(Having)
		havingTest.Append(condition...)

		if havingTest.String() != expected {
			t.Fatalf(`Query %s != %s`, havingTest.String(), expected)
		}
	}

	// Test case 2: Multiple OR conditions
	var multipleOrConditions []Condition
	multipleOrConditions = append(multipleOrConditions, Condition{
		Field: "SUM(quantity)",
		Opt:   Greater,
		Value: 100,
	}, Condition{
		Field: "COUNT(order_id)",
		Opt:   Greater,
		Value: 5,
		AndOr: Or,
	}, Condition{
		Field: "AVG(price)",
		Opt:   Greater,
		Value: 50,
		AndOr: Or,
	})

	multipleOrTestCases := map[string][]Condition{
		"HAVING SUM(quantity) > 100 OR COUNT(order_id) > 5 OR AVG(price) > 50": multipleOrConditions,
	}

	for expected, condition := range multipleOrTestCases {
		havingTest := new(Having)
		havingTest.Append(condition...)

		if havingTest.String() != expected {
			t.Fatalf(`Query %s != %s`, havingTest.String(), expected)
		}
	}
}

// TestHavingMixed tests mixed AND and OR conditions in HAVING clause
func TestHavingMixed(t *testing.T) {
	var conditions []Condition

	conditions = append(conditions, Condition{
		Field: "COUNT(employee_id)",
		Opt:   Greater,
		Value: 10,
	}, Condition{
		Field: "AVG(salary)",
		Opt:   Greater,
		Value: 5000,
		AndOr: Or,
	}, Condition{
		Field: "department_id",
		Opt:   Eq,
		Value: 5,
	})

	testCases := map[string][]Condition{
		"HAVING COUNT(employee_id) > 10 OR AVG(salary) > 5000 AND department_id = 5": conditions,
	}

	for expected, condition := range testCases {
		havingTest := new(Having)
		havingTest.Append(condition...)

		if havingTest.String() != expected {
			t.Fatalf(`Query %s != %s`, havingTest.String(), expected)
		}
	}
}

// TestHavingEmpty tests that an empty HAVING clause returns an empty string
func TestHavingEmpty(t *testing.T) {
	havingTest := new(Having)

	if havingTest.String() != "" {
		t.Fatalf(`Empty HAVING clause should return empty string, got %s`, havingTest.String())
	}
}
