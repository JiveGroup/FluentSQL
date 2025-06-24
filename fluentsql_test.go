package fluentsql

import "testing"

func TestSetDialect(t *testing.T) {
	// Save the original dialect to restore it after the test
	originalDialect := defaultDialect
	defer func() {
		defaultDialect = originalDialect
	}()

	// Test setting MySQL dialect
	mysqlDialect := new(MySQLDialect)
	SetDialect(mysqlDialect)
	if defaultDialect != mysqlDialect {
		t.Fatalf("Expected dialect to be MySQL, got %v", defaultDialect.Name())
	}

	// Test setting SQLite dialect
	sqliteDialect := new(SQLiteDialect)
	SetDialect(sqliteDialect)
	if defaultDialect != sqliteDialect {
		t.Fatalf("Expected dialect to be SQLite, got %v", defaultDialect.Name())
	}

	// Test setting PostgreSQL dialect
	pgDialect := new(PostgreSQLDialect)
	SetDialect(pgDialect)
	if defaultDialect != pgDialect {
		t.Fatalf("Expected dialect to be PostgreSQL, got %v", defaultDialect.Name())
	}
}

func TestIsDialect(t *testing.T) {
	// Save the original dialect to restore it after the test
	originalDialect := defaultDialect
	defer func() {
		defaultDialect = originalDialect
	}()

	// Test with MySQL dialect
	mysqlDialect := new(MySQLDialect)
	SetDialect(mysqlDialect)
	if !IsDialect(MySQL) {
		t.Fatalf("Expected IsDialect(%s) to be true", MySQL)
	}
	if IsDialect(PostgreSQL) {
		t.Fatalf("Expected IsDialect(%s) to be false", PostgreSQL)
	}
	if IsDialect(SQLite) {
		t.Fatalf("Expected IsDialect(%s) to be false", SQLite)
	}

	// Test with PostgreSQL dialect
	pgDialect := new(PostgreSQLDialect)
	SetDialect(pgDialect)
	if !IsDialect(PostgreSQL) {
		t.Fatalf("Expected IsDialect(%s) to be true", PostgreSQL)
	}
	if IsDialect(MySQL) {
		t.Fatalf("Expected IsDialect(%s) to be false", MySQL)
	}
	if IsDialect(SQLite) {
		t.Fatalf("Expected IsDialect(%s) to be false", SQLite)
	}

	// Test with SQLite dialect
	sqliteDialect := new(SQLiteDialect)
	SetDialect(sqliteDialect)
	if !IsDialect(SQLite) {
		t.Fatalf("Expected IsDialect(%s) to be true", SQLite)
	}
	if IsDialect(MySQL) {
		t.Fatalf("Expected IsDialect(%s) to be false", MySQL)
	}
	if IsDialect(PostgreSQL) {
		t.Fatalf("Expected IsDialect(%s) to be false", PostgreSQL)
	}
}

func TestMySQLDialect_Name(t *testing.T) {
	dialect := new(MySQLDialect)
	if dialect.Name() != MySQL {
		t.Fatalf("Expected MySQL dialect name to be %s, got %s", MySQL, dialect.Name())
	}
}

func TestMySQLDialect_Placeholder(t *testing.T) {
	dialect := new(MySQLDialect)
	if dialect.Placeholder(1) != "?" {
		t.Fatalf("Expected MySQL placeholder to be ?, got %s", dialect.Placeholder(1))
	}
	if dialect.Placeholder(2) != "?" {
		t.Fatalf("Expected MySQL placeholder to be ?, got %s", dialect.Placeholder(2))
	}
}

func TestMySQLDialect_YearFunction(t *testing.T) {
	dialect := new(MySQLDialect)
	expected := "YEAR(hire_date)"
	result := dialect.YearFunction("hire_date")
	if result != expected {
		t.Fatalf("Expected MySQL year function to be %s, got %s", expected, result)
	}
}

func TestSQLiteDialect_Name(t *testing.T) {
	dialect := new(SQLiteDialect)
	if dialect.Name() != SQLite {
		t.Fatalf("Expected SQLite dialect name to be %s, got %s", SQLite, dialect.Name())
	}
}

func TestSQLiteDialect_Placeholder(t *testing.T) {
	dialect := new(SQLiteDialect)
	if dialect.Placeholder(1) != "?" {
		t.Fatalf("Expected SQLite placeholder to be ?, got %s", dialect.Placeholder(1))
	}
	if dialect.Placeholder(2) != "?" {
		t.Fatalf("Expected SQLite placeholder to be ?, got %s", dialect.Placeholder(2))
	}
}

func TestSQLiteDialect_YearFunction(t *testing.T) {
	dialect := new(SQLiteDialect)
	expected := "strftime('%Y', hire_date)"
	result := dialect.YearFunction("hire_date")
	if result != expected {
		t.Fatalf("Expected SQLite year function to be %s, got %s", expected, result)
	}
}
