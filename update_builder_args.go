package fluentsql

import (
	"fmt"
	"strings"
)

// Sql generates the SQL query string and its corresponding arguments.
func (ub *UpdateBuilder) Sql() (string, []any, interface{}) {
	return ub.StringArgs()
}

// StringArgs constructs the SQL query string and collects the argument values.
// Returns the SQL query string, the list of arguments, and an error if any occurred.
func (ub *UpdateBuilder) StringArgs() (string, []any, error) {
	var queryParts []string // Holds different parts of the SQL query.
	var sql string          // The final SQL query string.
	var args []any          // A slice of arguments to be used in the query.

	// Add UPDATE statement.
	sql, args = ub.updateStatement.StringArgs(args)
	queryParts = append(queryParts, sql)

	// Add SET statement.
	sql, args = ub.setStatement.StringArgs(args)
	queryParts = append(queryParts, sql)

	// Add WHERE clause if present.
	sql, args = ub.whereStatement.StringArgs(args)
	if sql != "" {
		queryParts = append(queryParts, sql)
	}

	// Add ORDER BY clause if present.
	sql, args = ub.orderByStatement.StringArgs(args)
	if sql != "" {
		queryParts = append(queryParts, sql)
	}

	// Add LIMIT clause if present.
	sql, args = ub.limitStatement.StringArgs(args)
	if sql != "" {
		queryParts = append(queryParts, sql)
	}

	// Combine all query parts into a single string.
	sql = strings.Join(queryParts, " ")

	return sql, args, nil
}

// StringArgs generates the SQL fragment for the UPDATE statement and appends to provided arguments.
// Parameters:
// - args: A slice of arguments to be appended to.
//
// Returns:
// - A formatted SQL UPDATE string.
// - A slice of arguments.
func (u *Update) StringArgs(args []any) (string, []any) {
	var sb strings.Builder // Used for efficient string concatenation.
	sb.WriteString(fmt.Sprintf("UPDATE %s", u.Table))

	// Add table alias if present.
	if u.Alias != "" {
		sb.WriteString(" " + u.Alias)
	}

	return sb.String(), args
}

// StringArgs generates the SQL fragment for the SET clause and appends to provided arguments.
// Parameters:
// - args: A slice of arguments to be appended to.
//
// Returns:
// - A formatted SQL SET string.
// - A slice of arguments.
func (s *UpdateSet) StringArgs(args []any) (string, []any) {
	var setColumns []string // Holds the individual SET assignments.

	// Process each item in the SET clause.
	for _, item := range s.Items {
		var sql string

		sql, args = item.StringArgs(args)

		setColumns = append(setColumns, sql)
	}

	return fmt.Sprintf("SET %s", strings.Join(setColumns, ", ")), args
}

// StringArgs generates the SQL fragment for an individual assignment in the SET clause.
// Parameters:
// - args: A slice of arguments to be appended to.
//
// Returns:
// - A formatted SQL string for the assignment.
// - A slice of arguments.
func (s *UpdateItem) StringArgs(args []any) (string, []any) {
	// Check if Field is a slice of strings for multi-column updates.
	// SET (field1, field2,...) = (int, string, ValueField...)
	// SET (field1, field2,...) = (SELECT * FROM table_name)
	if fieldStringSlice, ok := s.Field.([]string); ok {
		fieldStr := joinSlice(fieldStringSlice, ",") // Join field names with commas.

		// If the value is a QueryBuilder, process the associated query.
		if valueQueryBuilder, ok := s.Value.(*QueryBuilder); ok {
			var _sql string
			_sql, args, _ = valueQueryBuilder.StringArgs(args)

			return fmt.Sprintf("(%s) = (%s)", fieldStr, _sql), args
		}

		// If the value is a slice, process each item in the slice.
		if fieldAnySlice, ok := s.Value.([]any); ok {
			var values []string
			for _, fieldAny := range fieldAnySlice {
				if valueField, ok := fieldAny.(IValueField); ok {
					values = append(values, valueField.Value())
				} else if valueString, ok := fieldAny.(string); ok {
					args = append(args, valueString)
					valueStr := p(args)

					values = append(values, valueStr)
				} else { // Value is an int or float.
					args = append(args, fieldAny)
					valueStr := p(args)

					values = append(values, valueStr)
				}
			}

			valueStr := strings.Join(values, ", ")

			return fmt.Sprintf("(%s) = (%v)", fieldStr, valueStr), args
		}

		return "", args
	}

	// If the value is a QueryBuilder, process the associated query.
	if valueQueryBuilder, ok := s.Value.(*QueryBuilder); ok {
		var _sql string
		_sql, args, _ = valueQueryBuilder.StringArgs(args)

		return fmt.Sprintf("%s = (%s)", s.Field, _sql), args
	}

	// If the value is a ValueField, format it as-is.
	if valueField, ok := s.Value.(IValueField); ok {
		return fmt.Sprintf("%s = %s", s.Field, valueField.Value()), args
	}

	// If the value is a string, add it to the arguments and format it.
	if valueString, ok := s.Value.(string); ok {
		args = append(args, valueString)
		valueStr := p(args)

		return fmt.Sprintf("%s = %s", s.Field, valueStr), args
	}

	// Default fallback for other types (e.g., int, float).
	args = append(args, s.Value)
	valueStr := p(args)

	return fmt.Sprintf("%s = %s", s.Field, valueStr), args
}
