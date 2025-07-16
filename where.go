package fluentsql

import (
	"fmt"
	"reflect"
	"strings"
)

// Where clause
type Where struct {
	// Conditions represent a slice of Condition structs that define the WHERE clause of a SQL query.
	Conditions []Condition
}

// Append adds one or more Condition instances to the Conditions slice of the Where struct.
//
// Parameters:
//   - conditions: One or more Condition instances to be added.
//
// Usage:
//
//	w.Append(condition1, condition2, ...)
//
// This function appends the given conditions to the existing Conditions slice.
func (w *Where) Append(conditions ...Condition) {
	w.Conditions = append(w.Conditions, conditions...)
}

// String generates and returns the SQL representation of the WHERE clause.
//
// Returns:
//   - string: A string containing the SQL WHERE clause. If no conditions are present, it returns an empty string.
//
// Usage:
//
//	query := w.String()
//
// Example:
//
//	  WHERE clause representation of conditions will be formatted as:
//		 WHERE condition1 AND condition2 OR condition3
func (w *Where) String() string {
	var conditions []string

	// Loop through each Condition in the Conditions slice.
	if len(w.Conditions) > 0 {
		for _, cond := range w.Conditions {
			// Convert the current condition to its string representation.
			var _condition = cond.String()

			// If the operator is OR, combine it with the previous condition.
			if cond.AndOr == Or && len(conditions) > 0 {
				_orCondition := fmt.Sprint(" OR ", _condition)

				// Modify the last condition by appending the OR condition.
				last := len(conditions) - 1
				conditions[last] += _orCondition
			} else {
				// Append the current condition to the slice.
				conditions = append(conditions, _condition)
			}
		}
	}

	// No WHERE condition
	if len(conditions) == 0 {
		return ""
	}

	// Join all conditions with "AND" and return the final WHERE clause.
	return fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND "))
}

// Condition type struct
type Condition struct {
	// Field represents the name of the column to compare. It can be of type `string` or `FieldNot`.
	Field any
	// Opt specifies the condition operator such as =, <>, >, <, >=, <=, LIKE, IN, NOT IN, BETWEEN, etc.
	Opt WhereOpt
	// Value holds the value to be compared against the field. Support ValueField for checking with table's column
	Value any
	// AndOr specifies the logical combination with the previous condition (AND, OR). Default is AND.
	AndOr WhereAndOr
	// Group contains sub-conditions enclosed in parentheses `()`.
	Group []Condition
}

// WhereOpt defines the operators used in SQL conditions.
type WhereOpt int

const (
	Eq         WhereOpt = iota // Equal to (=)
	NotEq                      // Not equal to (<>)
	Diff                       // Not equal to (!=)
	Greater                    // Greater than (>)
	Lesser                     // Less than (<)
	GrEq                       // Greater than or equal to (>=)
	LeEq                       // Less than or equal to (<=)
	Like                       // Pattern matching (LIKE)
	NotLike                    // Not pattern matching (NOT LIKE)
	In                         // Value in a list (IN)
	NotIn                      // Value not in a list (NOT IN)
	Between                    // Value in a range (BETWEEN)
	NotBetween                 // Value not in a range (NOT BETWEEN)
	Null                       // Null value (IS NULL)
	NotNull                    // Not null value (IS NOT NULL)
	Exists                     // Subquery results exist (EXISTS)
	NotExists                  // Subquery results do not exist (NOT EXISTS)
	EqAny                      // Equal to any value in a subquery (= ANY)
	NotEqAny                   // Not equal to any value in a subquery (<> ANY)
	DiffAny                    // Not equal to any value in a subquery (!= ANY)
	GreaterAny                 // Greater than any value in a subquery (> ANY)
	LesserAny                  // Less than any value in a subquery (< ANY)
	GrEqAny                    // Greater than or equal to any value in a subquery (>= ANY)
	LeEqAny                    // Less than or equal to any value in a subquery (<= ANY)
	EqAll                      // Equal to all values in a subquery (= ALL)
	NotEqAll                   // Not equal to all values in a subquery (<> ALL)
	DiffAll                    // Not equal to all values in a subquery (!= ALL)
	GreaterAll                 // Greater than all values in a subquery (> ALL)
	LesserAll                  // Less than all values in a subquery (< ALL)
	GrEqAll                    // Greater than or equal to all values in a subquery (>= ALL)
	LeEqAll                    // Less than or equal to all values in a subquery (<= ALL)
)

// opt determines and returns the SQL operator (e.g., =, >, LIKE) corresponding to the Opt field.
//
// Returns:
//   - string: The string representation of the SQL operator.
func (c *Condition) opt() string {
	var sign string

	// Map the WhereOpt value to its corresponding SQL operator.
	switch c.Opt {
	case Eq:
		sign = "="
	case NotEq:
		sign = "<>"
	case Diff:
		sign = "!="
	case Greater:
		sign = ">"
	case Lesser:
		sign = "<"
	case GrEq:
		sign = ">="
	case LeEq:
		sign = "<="
	case Like:
		sign = "LIKE"
	case NotLike:
		sign = "NOT LIKE"
	case In:
		sign = "IN"
	case NotIn:
		sign = "NOT IN"
	case Between:
		sign = "BETWEEN"
	case NotBetween:
		sign = "NOT BETWEEN"
	case Null:
		sign = "IS NULL"
	case NotNull:
		sign = "IS NOT NULL"
	case Exists:
		sign = "EXISTS"
	case NotExists:
		sign = "NOT EXISTS"
	case EqAny:
		sign = "= ANY"
	case NotEqAny:
		sign = "<> ANY"
	case DiffAny:
		sign = "!= ANY"
	case GreaterAny:
		sign = "> ANY"
	case LesserAny:
		sign = "< ANY"
	case GrEqAny:
		sign = ">= ANY"
	case LeEqAny:
		sign = "<= ANY"
	case EqAll:
		sign = "= ALL"
	case NotEqAll:
		sign = "<> ALL"
	case DiffAll:
		sign = "!= ALL"
	case GreaterAll:
		sign = "> ALL"
	case LesserAll:
		sign = "< ALL"
	case GrEqAll:
		sign = ">= ALL"
	case LeEqAll:
		sign = "<= ALL"
	}

	return sign
}

// String generates the SQL representation of the Condition.
//
// This function converts the condition into a SQL WHERE clause representation,
// supporting various SQL operators, subqueries, groups, and field-value conditions.
//
// Returns:
//   - string: A SQL string representation of the condition.
func (c *Condition) String() string {
	// Handle group conditions WhereGroup(groupCondition FnWhereBuilder)
	if len(c.Group) > 0 {
		var conditions []string

		// Iterate over the grouped conditions.
		for _, cond := range c.Group {
			// Generate the SQL string representation for each condition in the group.
			var _condition = cond.String()

			// If the logical operator is OR, append the condition with "OR".
			if cond.AndOr == Or && len(conditions) > 0 {
				_orCondition := fmt.Sprint(" OR ", _condition)

				// Update the last condition with the "OR" combination.
				last := len(conditions) - 1
				conditions[last] += _orCondition
			} else {
				// Append the condition to the list.
				conditions = append(conditions, _condition)
			}
		}

		// Return an empty string if no conditions are present in the group.
		if len(conditions) == 0 {
			return ""
		}

		// Join the group conditions with "AND" and enclose them in parentheses.
		return fmt.Sprintf("(%s)", strings.Join(conditions, " AND "))
	}

	// Handle IS NULL and IS NOT NULL conditions.
	// Example: WHERE Address IS NULL or WHERE Address IS NOT NULL
	if c.Opt == Null || c.Opt == NotNull {
		return fmt.Sprintf("%s %s", c.Field, c.opt())
	}

	// Handle IN and NOT IN conditions.
	// Example: WHERE Country IN ('Germany', 'France', 'UK')
	// Example: WHERE Age NOT IN (12, 31, 21)
	if c.Opt == In || c.Opt == NotIn {
		// Determine the type of the value being compared (e.g., slice or array).
		typ := reflect.TypeOf(c.Value)

		if typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
			valuesStr := ""
			// Convert slice of string values to SQL string format.
			if values, ok := c.Value.([]string); ok {
				valuesStr = "'" + strings.Join(values, "', '") + "'"
			}
			// Convert slice of integer values to comma-separated string format.
			if values, ok := c.Value.([]int); ok {
				valuesStr = joinSlice(values, ",")
			}

			// Generate the SQL representation.
			return fmt.Sprintf("%s %s (%s)", c.Field, c.opt(), valuesStr)
		}
	}

	// WHERE Price BETWEEN 10 AND 20;
	// WHERE ProductName BETWEEN 'Carnation Tigers' AND 'Mozzarella di Giovanni'
	// WHERE Price NOT BETWEEN 10 AND 20;
	// WHERE ProductName NOT BETWEEN 'Carnation Tigers' AND 'Mozzarella di Giovanni'
	// WHERE Price BETWEEN 10 AND 20
	if c.Opt == Between || c.Opt == NotBetween {
		return fmt.Sprintf("%s %s %v", c.Field, c.opt(), c.Value)
	}

	// WHERE salary = (SELECT DISTINCT salary FROM employees ORDER BY salary DESC LIMIT 1 , 1);
	// WHERE CustomerID IN (SELECT CustomerID FROM Orders);
	// WHERE CustomerID NOT IN (SELECT CustomerID FROM Orders);
	// WHERE EXISTS (SELECT ProductName FROM Products);
	// WHERE NOT EXISTS (SELECT ProductName FROM Products);
	// WHERE ProductID = ANY (SELECT ProductID FROM OrderDetails WHERE Quantity = 10);
	// WHERE ProductID > ALL (SELECT ProductID FROM OrderDetails WHERE Quantity = 10);
	if valueQueryBuilder, ok := c.Value.(*QueryBuilder); ok { // Column type is a complex query.
		return fmt.Sprintf("%s %s (%v)", c.Field, c.opt(), valueQueryBuilder)
	}

	// Handle string value conditions.
	// Example: WHERE Name = 'John'
	if valueString, ok := c.Value.(string); ok {
		return fmt.Sprintf("%s %s '%v'", c.Field, c.opt(), valueString)
	}

	// Default case: Handle simple field-value conditions.
	// Example: WHERE Age > 30
	return fmt.Sprintf("%s %s %v", c.Field, c.opt(), c.Value)
}

type WhereAndOr int

const (
	And WhereAndOr = iota // Logical AND operator for combining conditions
	Or                    // Logical OR operator for combining conditions
)

// ValueBetween for WhereOpt.Between or WhereOpt.NotBetween -
// A struct to define a range of values for SQL BETWEEN conditions.
type ValueBetween struct {
	// Low represents the lower bound of the range.
	Low any
	// High represents the upper bound of the range.
	High any
}

// String generates the SQL representation of the ValueBetween range.
//
// Returns:
//   - string: A string representing the range in the format "Low AND High"
//     If the bounds are strings, they are enclosed in single quotes.
//
// Examples:
//   - If Low = 1999 and High = 2000, it returns "1999 AND 2000"
//   - If Low = "1999-01-01" and High = "2000-12-31", it returns "'1999-01-01' AND '2000-12-31'"
func (v ValueBetween) String() string {
	if _, ok := v.Low.(string); ok {
		// If bounds are strings, include single quotes around them.
		// Eg: hire_date BETWEEN '1999-01-01' AND '2000-12-31'
		return fmt.Sprintf("'%v' AND '%v'", v.Low, v.High)
	}

	// Return the range as is for non-string bounds.
	// Eg: salary NOT BETWEEN 2500 AND 2900
	return fmt.Sprintf("%v AND %v", v.Low, v.High)
}

// ValueField represents a column/field in a SQL query as a string value.
//
// Methods:
//
//   - String(): Converts the ValueField to its string representation.
type ValueField string

// IValueField represents a field value interface that can be converted to a string value.
//
// The Value() method enables dynamic generation of SQL field values that need custom string formatting.
// This interface is commonly implemented by custom field types that need to control how they are
// represented in SQL queries.
//
// Example usage:
//
//	type CustomField struct {
//		name string
//	}
//
//	func (f CustomField) Value() string {
//		return fmt.Sprintf("CUSTOM(%s)", f.name)
//	}
//
//	// Can be used in SQL generation
//	field := CustomField{name: "test"}
//	sqlValue := field.Value() // Returns "CUSTOM(test)"
type IValueField interface {
	Value() string
}

// Value converts the ValueField to its string representation.
//
// Returns:
//   - string: The string representation of the ValueField.
//
// Examples
// ValueField to handle condition `WHERE c.column <WhereOpt> c.column_1`
//
//	So, When build condition Where("d.employee_id", qb.Eq, qb.ValueField("e.employee_id")) to keep SQL string as
//	    d.employee_id = e.employee_id
//	    instead of
//	    d.employee_id = 'e.employee_id'
func (v ValueField) Value() string {
	return string(v)
}

// FieldNot represents a SQL field prefixed with a NOT for negating conditions.
//
// Methods:
//
//   - String(): Returns the SQL string representation of the negated field.
type FieldNot string

// String generates the SQL representation of the FieldNot type.
//
// Returns:
//   - string: A string prefixed with "NOT" followed by the field name.
//
// Examples:
//   - FieldNot to handle condition `WHERE NOT salary > 5000` So, When build condition
//     Where(qb.FieldNot("salary"), qb.Greater, 5000) to keep SQL string as `NOT salary > 5000`
func (v FieldNot) String() string {
	return fmt.Sprintf("NOT %s", string(v))
}

// FieldEmpty represents an empty SQL field, often used in conditions like EXISTS or NOT EXISTS.
//
// Example:
//   - FieldEmpty to handle condition `WHERE NOT EXISTS (SELECT employee_id FROM dependents)`
type FieldEmpty string

// String returns the string representation of the FieldEmpty type.
//
// Returns:
//   - string: The string value of the FieldEmpty instance.
//
// Usage:
//   - This is typically used to generate SQL queries where a placeholder field is required for EXISTS or NOT EXISTS clauses.
func (v FieldEmpty) String() string {
	return string(v)
}

// FieldYear represents a SQL year extraction operation for a given field.
// It generates the SQL syntax for extracting the year portion from a date field based on the database type.
//
// Example usage:
//   - MySQL: YEAR(hire_date) Between 1990 AND 1993
//   - PostgreSQL: DATE_PART('year', hire_date) Between 1990 AND 1993
//   - SQLite: strftime('%Y', hire_date)
type FieldYear string

// String generates the SQL representation of a FieldYear based on the database type.
//
// Returns:
//   - string: The SQL string for extracting the year from a date field.
//
// Usage:
//   - To use as part of a WHERE clause or SELECT statement.
//   - The generated output varies depending on the database type (MySQL, PostgreSQL, SQLite).
func (v FieldYear) String() string {
	return DefaultDialect().YearFunction(string(v))
}
