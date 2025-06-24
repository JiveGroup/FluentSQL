package fluentsql

import (
	"fmt"
	"reflect"
	"strings"
)

// Sql constructs the SQL query string and associated arguments.
// It uses the StringArgs method to combine all query parts (SELECT, FROM, WHERE, etc.).
//
// Returns:
// - string: The complete SQL query as a string.
// - []any: A slice of arguments to be used with the query (e.g., for prepared statements).
// - error: Any error encountered during the query construction.
func (qb *QueryBuilder) Sql() (string, []any, error) {
	var args []any // Slice to hold query arguments

	return qb.StringArgs(args)
}

// StringArgs constructs the SQL query string with placeholders and its associated arguments.
//
// Parameters:
// - args []any: An initial slice of arguments to be used in the query.
//
// Returns:
// - string: The complete SQL query as a string with placeholders for arguments.
// - []any: A slice containing all arguments for the query.
// - error: Any error encountered during query string construction.
func (qb *QueryBuilder) StringArgs(args []any) (string, []any, error) {
	var queryParts []string // Slice to hold the parts of the query
	var sqlStr string       // Variable to store the current query part

	sqlStr, args = qb.selectStatement.StringArgs(args)
	queryParts = append(queryParts, sqlStr)

	sqlStr, args = qb.fromStatement.StringArgs(args)
	queryParts = append(queryParts, sqlStr)

	sqlStr, args = qb.joinStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr, args = qb.whereStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr, args = qb.groupByStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr, args = qb.havingStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr, args = qb.orderByStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr, args = qb.limitStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr, args = qb.fetchStatement.StringArgs(args)
	if sqlStr != "" {
		queryParts = append(queryParts, sqlStr)
	}

	sqlStr = strings.Join(queryParts, " ") // Combine all query parts into a single string

	if qb.alias != "" {
		sqlStr = fmt.Sprintf("(%s) AS %s",
			sqlStr,
			qb.alias)
	}

	return sqlStr, args, nil
}

// StringArgs generates the SQL SELECT statement string and associated arguments.
//
// Parameters:
// - args []any: A slice of arguments for constructing the query.
//
// Returns:
// - string: The complete SQL SELECT statement as a string.
// - []any: A slice containing the arguments used in the query.
func (s *Select) StringArgs(args []any) (string, []any) {
	selectOf := "*" // Default to selecting all columns

	if len(s.Columns) > 0 {
		var columns []string

		// Iterate through each column to process its type and generate the corresponding SQL part
		for _, col := range s.Columns {
			var sqlPart string
			if _, ok := col.(*Case); ok { // Column is of type Case
				sqlPart, args = col.(*Case).StringArgs(args)
				columns = append(columns, sqlPart)
			} else if valueString, ok := col.(string); ok { // Column is a plain string
				columns = append(columns, valueString)
			} else if valueFieldYear, ok := col.(FieldYear); ok { // Column is of type FieldYear
				columns = append(columns, valueFieldYear.String())
			} else if valueQueryBuilder, ok := col.(*QueryBuilder); ok { // Column is a QueryBuilder
				var selectQuery string
				selectQuery, args, _ = valueQueryBuilder.StringArgs(args)

				// Wrap the query in parentheses if no alias is provided
				if col.(*QueryBuilder).alias == "" {
					selectQuery = fmt.Sprintf("(%s)", selectQuery)
				}

				columns = append(columns, selectQuery)
			}
		}

		// Combine all column representations into a comma-separated string
		selectOf = strings.Join(columns, ", ")
	}

	// Return the constructed SELECT statement and associated arguments
	return fmt.Sprintf("SELECT %s", selectOf), args
}

// StringArgs generates the SQL FROM clause string and associated arguments.
//
// Parameters:
// - args []any: A slice of arguments for constructing the query.
//
// Returns:
// - string: The SQL FROM clause string.
// - []any: A slice containing the arguments used in the clause.
func (f *From) StringArgs(args []any) (string, []any) {
	var sb strings.Builder // String builder for constructing the FROM clause

	// Process the table source based on its type
	if _, ok := f.Table.(string); ok { // Table is a plain string
		sb.WriteString(fmt.Sprintf("FROM %s", f.Table))
	} else if _, ok := f.Table.(*QueryBuilder); ok { // Table is a QueryBuilder
		var selectQuery string
		selectQuery, args, _ = f.Table.(*QueryBuilder).StringArgs(args)

		// Wrap the query in parentheses if no alias is provided
		if f.Table.(*QueryBuilder).alias == "" {
			sb.WriteString(fmt.Sprintf("FROM (%s)", selectQuery))
		} else {
			sb.WriteString(fmt.Sprintf("FROM %s", selectQuery))
		}
	}

	// Append the table alias if provided
	if f.Alias != "" {
		sb.WriteString(" " + f.Alias)
	}

	// Return the constructed FROM clause and associated arguments
	return sb.String(), args
}

// StringArgs generates the SQL JOIN clause string and associated arguments.
//
// Parameters:
// - args []any: A slice of arguments for constructing the query.
//
// Returns:
// - string: The SQL JOIN clause string. Returns an empty string if there are no JOIN items.
// - []any: A slice containing the arguments used in the clause.
func (j *Join) StringArgs(args []any) (string, []any) {
	// Return empty string if there are no join items
	if len(j.Items) == 0 {
		return "", args
	}

	var joinItems []string // Slice to hold each join statement

	// Process each join item to generate the full join statement
	for _, item := range j.Items {
		var cond string
		cond, args = item.Condition.StringArgs(args)

		// Construct the join string based on the type of JOIN
		joinStr := fmt.Sprintf("%s %s ON %s", item.opt(), item.Table, cond)

		// For CROSS JOIN, omit the ON clause
		if item.Join == CrossJoin {
			joinStr = fmt.Sprintf("%s %s", item.opt(), item.Table)
		}

		joinItems = append(joinItems, joinStr)
	}

	// Combine all join statements into a single string and return
	return strings.Join(joinItems, " "), args
}

// StringArgs generates the SQL WHERE clause string and associated arguments.
//
// Parameters:
// - args []any: A slice of arguments for constructing the WHERE clause.
//
// Returns:
// - string: The complete SQL WHERE clause string. Returns an empty string if no conditions are present.
// - []any: A slice containing the arguments used in the WHERE clause.
func (w *Where) StringArgs(args []any) (string, []any) {
	var conditions []string // Slice to hold individual condition strings.

	// Process each condition in the Where struct.
	if len(w.Conditions) > 0 {
		for _, cond := range w.Conditions {
			var _condition string
			_condition, args = cond.StringArgs(args)

			// Handle "OR" conditions.
			if cond.AndOr == Or && len(conditions) > 0 {
				_orCondition := fmt.Sprint(" OR ", _condition)

				last := len(conditions) - 1
				conditions[last] += _orCondition
			} else {
				conditions = append(conditions, _condition)
			}
		}
	}

	// Return an empty string if no conditions exist.
	if len(conditions) == 0 {
		return "", args
	}

	// Construct the WHERE clause by joining conditions with "AND".
	return fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND ")), args
}

// StringArgs generates the SQL condition string and associated arguments.
//
// Parameters:
// - args []any: A slice of arguments for building the condition.
//
// Returns:
// - string: The SQL condition as a string.
// - []any: A slice containing the arguments used in the condition.
func (c *Condition) StringArgs(args []any) (string, []any) {
	// Handle group conditions (nested conditions).
	if len(c.Group) > 0 {
		var conditions []string // Slice to store grouped condition strings.

		for _, cond := range c.Group {
			var _condition string
			_condition, args = cond.StringArgs(args)

			// Handle "OR" conditions within the group.
			if cond.AndOr == Or && len(conditions) > 0 {
				_orCondition := fmt.Sprint(" OR ", _condition)

				last := len(conditions) - 1
				conditions[last] += _orCondition
			} else {
				conditions = append(conditions, _condition)
			}
		}

		// Return an empty string if no conditions are in the group.
		if len(conditions) == 0 {
			return "", args
		}

		return fmt.Sprintf("(%s)", strings.Join(conditions, " AND ")), args
	}

	// Handle ValueField type, excluding it from arguments.
	if valueField, ok := c.Value.(ValueField); ok {
		return fmt.Sprintf("%s %s %s", c.Field, c.opt(), valueField), args
	}

	// Handle IS NULL and IS NOT NULL conditions.
	if c.Opt == Null || c.Opt == NotNull {
		return fmt.Sprintf("%s %s", c.Field, c.opt()), args
	}

	// Handle IN and NOT IN conditions.
	if c.Opt == In || c.Opt == NotIn {
		typ := reflect.TypeOf(c.Value)

		if typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
			var valuesStr []string // Slice to store stringified values.

			// Process string slices.
			if values, ok := c.Value.([]string); ok {
				for _, val := range values {
					args = append(args, val)
					valuesStr = append(valuesStr, p(args))
				}
			}

			// Process integer slices.
			if values, ok := c.Value.([]int); ok {
				for _, val := range values {
					args = append(args, val)
					valuesStr = append(valuesStr, p(args))
				}
			}

			return fmt.Sprintf("%s %s (%s)", c.Field, c.opt(), strings.Join(valuesStr, ", ")), args
		}
	}

	// Handle BETWEEN and NOT BETWEEN conditions.
	// WHERE Price BETWEEN 10 AND 20;
	// WHERE ProductName BETWEEN 'Carnation Tigers' AND 'Mozzarella di Giovanni'
	// WHERE Price NOT BETWEEN 10 AND 20;
	// WHERE ProductName NOT BETWEEN 'Carnation Tigers' AND 'Mozzarella di Giovanni'
	// WHERE Price BETWEEN 10 AND 20
	if c.Opt == Between || c.Opt == NotBetween {
		var betweenValue string
		betweenValue, args = c.Value.(ValueBetween).StringArgs(args)

		return fmt.Sprintf("%s %s %v", c.Field, c.opt(), betweenValue), args
	}

	// Handle subqueries and nested QueryBuilder objects.
	// WHERE salary = (SELECT DISTINCT salary FROM employees ORDER BY salary DESC LIMIT 1 , 1);
	// WHERE CustomerID IN (SELECT CustomerID FROM Orders);
	// WHERE CustomerID NOT IN (SELECT CustomerID FROM Orders);
	// WHERE EXISTS (SELECT ProductName FROM Products);
	// WHERE NOT EXISTS (SELECT ProductName FROM Products);
	// WHERE ProductID = ANY (SELECT ProductID FROM OrderDetails WHERE Quantity = 10);
	// WHERE ProductID > ALL (SELECT ProductID FROM OrderDetails WHERE Quantity = 10);
	if valueQueryBuilder, ok := c.Value.(*QueryBuilder); ok {
		var queryBuilderStr string
		queryBuilderStr, args, _ = valueQueryBuilder.StringArgs(args)

		return fmt.Sprintf("%s %s (%v)", c.Field, c.opt(), queryBuilderStr), args
	}

	// Handle string values directly.
	if valueString, ok := c.Value.(string); ok {
		args = append(args, valueString)

		return fmt.Sprintf("%s %s %s", c.Field, c.opt(), p(args)), args
	}

	// Handle all other value types.
	args = append(args, c.Value)

	return fmt.Sprintf("%s %s %s", c.Field, c.opt(), p(args)), args
}

// StringArgs generates the SQL representation for a ValueBetween range
// and appends the Low and High values to the arguments slice.
//
// Parameters:
// - args []any: The input slice to which the Low and High values will be appended.
//
// Returns:
// - string: The SQL representation of the range in the format "LOW_PLACEHOLDER AND HIGH_PLACEHOLDER".
// - []any: The updated slice of arguments, including Low and High values.
func (v ValueBetween) StringArgs(args []any) (string, []any) {
	// Append lower bound to arguments and get its placeholder.
	args = append(args, v.Low)
	pLow := p(args)

	// Append upper bound to arguments and get its placeholder.
	args = append(args, v.High)
	pHigh := p(args)

	// Return SQL representation and updated arguments.
	// hire_date BETWEEN '1999-01-01' AND '2000-12-31'
	// salary NOT BETWEEN 2500 AND 2900
	return fmt.Sprintf("%v AND %v", pLow, pHigh), args
}

// StringArgs generates the SQL representation for extracting a year value
// from a field and adds the field value to the arguments slice.
//
// Parameters:
// - args []any: The input slice to which the field value will be appended.
//
// Returns:
// - string: The SQL representation for the year extraction, customized for the database type.
// - []any: The updated slice of arguments, including the field value.
func (v FieldYear) StringArgs(args []any) (string, []any) {
	// Append the field value to arguments and get its placeholder.
	args = append(args, string(v))

	return DefaultDialect().YearFunction(p(args)), args
}

// StringArgs generates the SQL GROUP BY clause string.
//
// Parameters:
// - args []any: The input slice of arguments (unused in this case).
//
// Returns:
// - string: The SQL GROUP BY clause string. Returns an empty string if no items are present.
// - []any: The unchanged slice of arguments.
func (g *GroupBy) StringArgs(args []any) (string, []any) {
	// Return empty if there are no group by items.
	if len(g.Items) == 0 {
		return "", args
	}

	// Construct and return GROUP BY clause.
	return fmt.Sprintf("GROUP BY %s", strings.Join(g.Items, ", ")), args
}

// StringArgs generates the SQL HAVING clause string and appends the associated argument values.
//
// Parameters:
// - args []any: The input slice to which the HAVING condition values will be appended.
//
// Returns:
// - string: The SQL HAVING clause string, combining conditions with "AND". Returns an empty string if no conditions are present.
// - []any: The updated slice of arguments, including condition values.
func (w *Having) StringArgs(args []any) (string, []any) {
	var conditions []string

	// Process each HAVING condition.
	if len(w.Conditions) > 0 {
		for _, cond := range w.Conditions {
			// Generate SQL and update arguments for each condition.
			var _condition string
			_condition, args = cond.StringArgs(args)

			// Handle "OR" conditions.
			if cond.AndOr == Or && len(conditions) > 0 {
				_orCondition := fmt.Sprint(" OR ", _condition)

				last := len(conditions) - 1
				conditions[last] += _orCondition
			} else {
				conditions = append(conditions, _condition)
			}
		}
	}

	// Return empty string if no conditions exist.
	if len(conditions) == 0 {
		return "", args
	}

	// Construct HAVING clause by joining conditions with "AND".
	return fmt.Sprintf("HAVING %s", strings.Join(conditions, " AND ")), args
}

// StringArgs generates the SQL ORDER BY clause string.
//
// Parameters:
// - args []any: The input slice of arguments (unused in this case).
//
// Returns:
// - string: The SQL ORDER BY clause string. Returns an empty string if no items are present.
// - []any: The unchanged slice of arguments.
func (o *OrderBy) StringArgs(args []any) (string, []any) {
	// Return empty if there are no order by items.
	if len(o.Items) == 0 {
		return "", args
	}

	var orderItems []string
	// Process each ORDER BY item.
	for _, item := range o.Items {
		orderItems = append(orderItems, fmt.Sprintf("%s %s", item.Field, item.Dir()))
	}

	// Construct and return ORDER BY clause.
	return fmt.Sprintf("ORDER BY %s", strings.Join(orderItems, ", ")), args
}

// StringArgs generates the SQL LIMIT and OFFSET clause strings
// and appends the limit and offset values to the arguments slice.
//
// Parameters:
// - args []any: The input slice to which the limit and offset values will be appended.
//
// Returns:
// - string: The SQL LIMIT and OFFSET clause string. Returns an empty string if both values are zero.
// - []any: The updated slice of arguments, including limit and offset values.
func (l *Limit) StringArgs(args []any) (string, []any) {
	// Append limit and offset values, and generate placeholders.
	if l.Limit > 0 || l.Offset > 0 {
		args = append(args, l.Limit)
		pLimit := p(args)
		args = append(args, l.Offset)
		pOffset := p(args)

		// Construct and return LIMIT/OFFSET clause.
		return fmt.Sprintf("LIMIT %s OFFSET %s", pLimit, pOffset), args
	}

	// Return empty string if no limit or offset is set.
	return "", args
}

// StringArgs generates the SQL FETCH NEXT ROWS clause string
// and appends the fetch and offset values to the arguments slice.
//
// Parameters:
// - args []any: The input slice to which the fetch and offset values will be appended.
//
// Returns:
// - string: The SQL FETCH NEXT ROWS clause string. Returns an empty string if both values are zero.
// - []any: The updated slice of arguments, including fetch and offset values.
func (f *Fetch) StringArgs(args []any) (string, []any) {
	// Append fetch and offset values, and generate placeholders.
	if f.Fetch > 0 || f.Offset > 0 {
		args = append(args, f.Offset)
		pOffset := p(args)
		args = append(args, f.Fetch)
		pFetch := p(args)

		// Construct and return FETCH NEXT ROWS clause.
		return fmt.Sprintf("OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", pOffset, pFetch), args
	}

	// Return empty string if no fetch or offset is set.
	return "", args
}

// StringArgs generates the SQL WHEN clause string for a CASE statement
// and appends the value and condition arguments to the slice.
//
// Parameters:
// - args []any: The input slice to which the value and conditions will be appended.
//
// Returns:
// - string: The SQL WHEN clause string.
// - []any: The updated slice of arguments, including value and condition values.
func (c *WhenCase) StringArgs(args []any) (string, []any) {
	// Process conditions and construct the WHEN clause SQL.
	if valueConditions, ok := c.Conditions.([]Condition); ok {
		var cons []string
		for _, condition := range valueConditions {
			var sqlPart string
			sqlPart, args = condition.StringArgs(args)

			cons = append(cons, sqlPart)
		}

		// Append the value associated with the WHEN clause.
		args = append(args, c.Value)

		// Construct and return the WHEN clause.
		return fmt.Sprintf("WHEN %s THEN '?'", strings.Join(cons, " AND ")), args
	}

	// Append the value associated with the WHEN clause.
	args = append(args, c.Value)

	// Return the WHEN clause for a single condition.
	return fmt.Sprintf("WHEN %v THEN '?'", c.Conditions), args
}

// StringArgs generates the SQL CASE statement string
// and appends the associated arguments to the slice.
//
// Parameters:
// - args []any: The input slice to which the expression and WHEN clause arguments will be appended.
//
// Returns:
// - string: The SQL CASE statement string.
// - []any: The updated slice of arguments, including expression and WHEN clause values.
func (c *Case) StringArgs(args []any) (string, []any) {
	var whenCases []string

	// Process each WHEN clause in the CASE statement.
	for _, whenClause := range c.WhenClauses {
		var sqlPart string
		sqlPart, args = whenClause.StringArgs(args)

		whenCases = append(whenCases, sqlPart)
	}

	// Construct and return the CASE statement.
	return fmt.Sprintf("CASE %s %s END %s", c.Exp, strings.Join(whenCases, " "), c.Name), args
}
