package fluentsql

import (
	"fmt"
	"strings"
)

type Case struct {
	// Exp specifies the expression to be evaluated in the CASE statement.
	Exp string
	// WhenClauses is a list of WHEN clauses defined for the CASE statement.
	WhenClauses []WhenCase
	// Name is the alias for the CASE statement.
	Name string
}

// Case clause
/*
Simple CASE expression
   CASE (2000 - YEAR(hire_date))
       WHEN 1 THEN '1 year'
       WHEN 3 THEN '3 years'
       WHEN 5 THEN '5 years'
       WHEN 10 THEN '10 years'
       WHEN 15 THEN '15 years'
       WHEN 20 THEN '20 years'
       WHEN 25 THEN '25 years'
       WHEN 30 THEN '30 years'
   END anniversary

Search CASE expression
   CASE
       WHEN salary < 3000 THEN 'Low'
       WHEN salary >= 3000 AND salary <= 5000 THEN 'Average'
       WHEN salary > 5000 THEN 'High'
   END evaluation
*/

// FieldCase creates a new Case instance with the provided expression and name.
//
// Parameters:
//   - exp: The expression to be evaluated in the CASE clause.
//   - name: The alias for the CASE clause.
//
// Returns:
//   - *Case: A pointer to a new Case instance.
func FieldCase(exp, name string) *Case {
	return &Case{
		Exp:  exp,
		Name: name,
	}
}

// When appends a new WHEN clause to the Case instance.
//
// Parameters:
//   - conditions: The condition(s) to evaluate (can be a single value, string, or a slice of Condition).
//   - value: The value to return when the condition is met.
//
// Returns:
//   - *Case: A pointer to the Case instance, for method chaining.
func (c *Case) When(conditions any, value string) *Case {
	c.WhenClauses = append(c.WhenClauses, WhenCase{
		Conditions: conditions,
		Value:      value,
	})

	return c
}

type WhenCase struct {
	// Conditions represents the condition(s) evaluated in the WHEN clause. It can be a string, integer, or slice of Condition.
	Conditions any
	// Value represents the result to return when the conditions are met.
	Value string
}

// String generates the SQL representation of the WHEN clause.
//
// Returns:
//   - string: The SQL string of the WHEN clause.
func (c *WhenCase) String() string {
	if valueConditions, ok := c.Conditions.([]Condition); ok {
		var cons []string
		for _, condition := range valueConditions {
			cons = append(cons, condition.String())
		}

		return fmt.Sprintf("WHEN %s THEN '%s'", strings.Join(cons, " AND "), c.Value)
	}

	return fmt.Sprintf("WHEN %v THEN '%s'", c.Conditions, c.Value)
}

// String generates the SQL representation of the entire CASE statement.
//
// Returns:
//   - string: The SQL string of the CASE statement.
func (c *Case) String() string {
	var whenCases []string

	for _, whenClause := range c.WhenClauses {
		whenCases = append(whenCases, whenClause.String())
	}

	return fmt.Sprintf("CASE %s %s END %s", c.Exp, strings.Join(whenCases, " "), c.Name)
}
