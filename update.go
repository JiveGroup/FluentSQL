package fluentsql

import (
	"fmt"
	"strings"
)

// Update clause
type Update struct {
	Table any    // Table indicates the target database table to be updated.
	Alias string // Alias represents an alternate name for the table that can be used in the query.
}

// String generates the SQL statement for an UPDATE clause.
//
// Returns:
//   - A string containing the formatted UPDATE statement.
func (u *Update) String() string {
	var sb strings.Builder // Used to efficiently build the SQL string.
	sb.WriteString(fmt.Sprintf("UPDATE %s", u.Table))

	if u.Alias != "" { // Add table alias to the statement if specified.
		sb.WriteString(" " + u.Alias)
	}

	return sb.String()
}

type UpdateItem struct {
	// Field name of column. Can be of type string or []string.
	Field any
	// Value data associated with the field. These could be of type string, int, ValueField, QueryBuilder, or []any.
	Value any
}

// String generates the SQL SET clause for an individual update item.
//
// Returns:
//   - A string representing the SET clause for the update field and value.
func (s *UpdateItem) String() string {
	// SET (field1, field2,...) = (int, string, ValueField...)
	// SET (field1, field2,...) = (SELECT * FROM table_name)
	if fieldStringSlice, ok := s.Field.([]string); ok { // Check if the field is of type []string.
		fieldStr := joinSlice(fieldStringSlice, ",") // Combine fields into a single comma-separated string.

		if valueQueryBuilder, ok := s.Value.(*QueryBuilder); ok { // Check if the value is a QueryBuilder.
			return fmt.Sprintf("(%s) = (%v)", fieldStr, valueQueryBuilder)
		}

		if fieldAnySlice, ok := s.Value.([]any); ok { // Check if the value is a slice of any type.
			var values []string
			for _, fieldAny := range fieldAnySlice {
				if valueField, ok := fieldAny.(IValueField); ok { // If the value is of ValueField type.
					values = append(values, valueField.Value())
				} else if valueString, ok := fieldAny.(string); ok { // If the value is of string type.
					values = append(values, "'"+valueString+"'")
				} else { // If the value is of numeric type.
					values = append(values, fmt.Sprintf("%v", fieldAny))
				}
			}

			valueStr := strings.Join(values, ", ") // Combine all values into a comma-separated string.

			return fmt.Sprintf("(%s) = (%v)", fieldStr, valueStr)
		}

		return ""
	}

	if valueQueryBuilder, ok := s.Value.(*QueryBuilder); ok { // Check if the value is a QueryBuilder.
		return fmt.Sprintf("%s = (%v)", s.Field, valueQueryBuilder)
	}

	if valueField, ok := s.Value.(IValueField); ok { // Check if the value is a ValueField.
		return fmt.Sprintf("%s = %s", s.Field, valueField.Value())
	}

	if valueString, ok := s.Value.(string); ok { // Check if the value is a string.
		return fmt.Sprintf("%s = '%s'", s.Field, valueString)
	}

	return fmt.Sprintf("%s = %v", s.Field, s.Value) // Default fallback for numeric types.
}

type UpdateSet struct {
	Items []UpdateItem // Items represents a collection of update assignments.
}

// Append adds a new update assignment to the Items slice.
//
// Parameters:
//   - field: The database field or column to be updated.
//   - value: The new value assigned to the field.
func (s *UpdateSet) Append(field, value any) {
	s.AppendItems(UpdateItem{
		Field: field,
		Value: value,
	})
}

// AppendItems appends multiple update assignments to the Items slice.
//
// Parameters:
//   - items: A variadic list of UpdateItem objects to be added.
func (s *UpdateSet) AppendItems(items ...UpdateItem) {
	s.Items = append(s.Items, items...)
}

// String generates the SQL SET clause combining all update assignments.
//
// Returns:
//   - A string representing the full SET clause of the update statement.
func (s *UpdateSet) String() string {
	var setColumns []string // A slice to hold the string representations of each update assignment.

	for _, item := range s.Items {
		setColumns = append(setColumns, item.String())
	}

	return fmt.Sprintf("SET %s", strings.Join(setColumns, ", ")) // Concatenates all assignments into a single string.
}
