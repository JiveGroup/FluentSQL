package fluentsql

import (
	"fmt"
	"strings"
)

type InsertRows struct {
	Rows []InsertRow
}

// Append adds a new row of values to the InsertRows.
//
// Parameters:
//   - values ...any: Variadic arguments representing the values of a new row.
func (r *InsertRows) Append(values ...any) {
	r.Rows = append(r.Rows, InsertRow{Values: values})
}

// String generates and returns the SQL VALUES clause representation of the InsertRows.
//
// Returns:
//   - string: The generated VALUES clause as a string.
func (r *InsertRows) String() string {
	var rowsStr []string

	// Generate string representation for each row.
	for _, row := range r.Rows {
		rowsStr = append(rowsStr, row.String())
	}

	// Return empty string if no rows were appended.
	if len(rowsStr) == 0 {
		return ""
	}

	return fmt.Sprintf("VALUES %s", strings.Join(rowsStr, ", "))
}

type InsertRow struct {
	Values []any
}

// String generates and returns the SQL representation of the row's values.
//
// Returns:
//   - string: The string representation of the row's values, formatted as a SQL tuple.
func (ir *InsertRow) String() string {
	var rowStr []string

	// Process each value in the row to generate its string representation.
	for _, col := range ir.Values {
		// Check if the value is of type ValueField.
		if colField, ok := col.(IValueField); ok {
			rowStr = append(rowStr, colField.String())
		} else if colString, ok := col.(string); ok { // Handle string values.
			rowStr = append(rowStr, "'"+colString+"'")
		} else { // Handle other types (int, float, etc.).
			rowStr = append(rowStr, fmt.Sprintf("%v", col))
		}
	}

	return fmt.Sprintf("(%s)", strings.Join(rowStr, ", "))
}
