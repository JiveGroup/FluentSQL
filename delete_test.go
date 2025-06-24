package fluentsql

import "testing"

// TestDelete
func TestDelete(t *testing.T) {
	testCases := map[string]Delete{
		"DELETE FROM products": {
			Table: "products",
		},
		"DELETE FROM products p": {
			Table: "products",
			Alias: "p",
		},
	}

	for expected, table := range testCases {
		if table.String() != expected {
			t.Fatalf(`Query %s != %s`, table.String(), expected)
		}
	}
}

// TestDeleteBuilderStringArgs tests the StringArgs method of DeleteBuilder
func TestDeleteBuilderStringArgs(t *testing.T) {
	// Test with ORDER BY clause
	db := DeleteInstance().Delete("products")
	db.orderByStatement.Append("price", Desc)

	sql, args, err := db.StringArgs([]any{})
	if err != nil {
		t.Fatalf("Error generating SQL: %v", err)
	}
	expected := "DELETE FROM products ORDER BY price DESC"
	if sql != expected {
		t.Fatalf("Query %s != %s", sql, expected)
	}
	if len(args) != 0 {
		t.Fatalf("Expected 0 args, got %d", len(args))
	}

	// Test with LIMIT clause
	db = DeleteInstance().Delete("products")
	db.limitStatement.Limit = 10
	db.limitStatement.Offset = 5

	sql, args, err = db.StringArgs([]any{})
	if err != nil {
		t.Fatalf("Error generating SQL: %v", err)
	}
	expected = "DELETE FROM products LIMIT $1 OFFSET $2"
	if sql != expected {
		t.Fatalf("Query %s != %s", sql, expected)
	}
	if len(args) != 2 || args[0] != 10 || args[1] != 5 {
		t.Fatalf("Expected args [10, 5], got %v", args)
	}

	// Test with table alias
	db = DeleteInstance().Delete("products", "p")

	sql, args, err = db.StringArgs([]any{})
	if err != nil {
		t.Fatalf("Error generating SQL: %v", err)
	}
	expected = "DELETE FROM products p"
	if sql != expected {
		t.Fatalf("Query %s != %s", sql, expected)
	}
	if len(args) != 0 {
		t.Fatalf("Expected 0 args, got %d", len(args))
	}
}
