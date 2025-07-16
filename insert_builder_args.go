package fluentsql

import (
	"fmt"
	"strings"
)

// Sql retrieves the SQL INSERT statement and its arguments.
//
// Returns:
//   - string: The complete SQL INSERT statement.
//   - []any: A slice containing the arguments for the statement.
//   - error: An error value (always nil in this implementation).
func (ib *InsertBuilder) Sql() (string, []any, error) {
	var args []any

	return ib.StringArgs(args)
}

// StringArgs constructs the SQL INSERT statement along with its arguments.
//
// Parameters:
//   - args []any: A slice of arguments to be used in the statement.
//
// Returns:
//   - string: The constructed SQL INSERT statement.
//   - []any: A slice containing the arguments for the statement.
//   - error: An error value (always nil in this implementation).
func (ib *InsertBuilder) StringArgs(args []any) (string, []any, error) {
	var queryParts []string
	var sqlStr string

	// Generate SQL string and arguments for the INSERT clause.
	sqlStr, args = ib.insertStatement.StringArgs(args)
	queryParts = append(queryParts, sqlStr)

	// Generate SQL string and arguments for the VALUES clause.
	sqlStr, args = ib.rowStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	// Generate SQL string and arguments for the SUBQUERY clause.
	sqlStr, args = ib.queryStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	// Combine all parts into a complete SQL INSERT statement.
	sql := strings.Join(queryParts, " ")

	return sql, args, nil
}

// StringArgs generates the SQL INSERT statement for a table with specified columns.
//
// Parameters:
//   - args []any: A slice of arguments to be used in the statement.
//
// Returns:
//   - string: The SQL INSERT statement for the table and columns.
//   - []any: The updated slice of arguments.
func (i *Insert) StringArgs(args []any) (string, []any) {
	columnsStr := strings.Join(i.Columns, ", ")
	return fmt.Sprintf("INSERT INTO %s (%s)", i.Table, columnsStr), args
}

// StringArgs generates the VALUES clause for the INSERT statement, including all rows.
//
// Parameters:
//   - args []any: A slice of arguments to be used in the statement.
//
// Returns:
//   - string: The SQL VALUES clause.
//   - []any: The updated slice of arguments.
func (r *InsertRows) StringArgs(args []any) (string, []any) {
	var rowsStr []string
	var sqlStr string

	// Process each row in the VALUES clause.
	for _, row := range r.Rows {
		sqlStr, args = row.StringArgs(args)
		rowsStr = append(rowsStr, sqlStr)
	}

	// Return empty string if no rows are specified.
	if len(rowsStr) == 0 {
		return "", args
	}

	return fmt.Sprintf("VALUES %s", strings.Join(rowsStr, ", ")), args
}

// StringArgs generates the SQL string for a single row of values and updates the arguments slice.
//
// Parameters:
//   - args []any: A slice of arguments to be used in the statement.
//
// Returns:
//   - string: The string representation of the row's values.
//   - []any: The updated slice of arguments.
func (ir *InsertRow) StringArgs(args []any) (string, []any) {
	var rowStr []string

	// Process each value in the row.
	for _, col := range ir.Values {
		if colField, ok := col.(IValueField); ok {
			rowStr = append(rowStr, colField.Value())
		} else if colString, ok := col.(string); ok {
			args = append(args, colString)
			colStr := p(args)
			rowStr = append(rowStr, colStr)
		} else { // Value is of type int or float.
			args = append(args, col)
			colStr := p(args)
			rowStr = append(rowStr, colStr)
		}
	}

	return fmt.Sprintf("(%s)", strings.Join(rowStr, ", ")), args
}

// StringArgs generates the SQL string for the subquery in the INSERT statement.
//
// Parameters:
//   - args []any: A slice of arguments to be used in the statement.
//
// Returns:
//   - string: The SQL string for the subquery.
//   - []any: The updated slice of arguments.
func (q *InsertQuery) StringArgs(args []any) (string, []any) {
	if queryBuilder, ok := q.Query.(*QueryBuilder); ok {
		var sqlStr string

		// Generate SQL string and arguments for the subquery.
		sqlStr, args, _ = queryBuilder.StringArgs(args)
		return sqlStr, args
	}

	// Return empty string if no subquery is specified.
	return "", args
}
