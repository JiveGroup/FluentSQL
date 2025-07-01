package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql
	qb "github.com/jivegroup/fluentsql"
	"github.com/jmoiron/sqlx"
)

// User represents a user in our system
type User struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

const (
	dbUser = "user"
	dbPass = "secret"
	dbName = "test"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting application...")

	// Connect to PostgreSQL
	// In a real application, you would get these values from environment variables or a config file
	connStr := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", dbUser, dbPass, dbName)
	db, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}(db)

	log.Println("Connected to database successfully")

	// Create the users table if it doesn't exist
	err = createUsersTable(db)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Perform CRUD operations
	user := User{
		Username:  "johndoe",
		Email:     "john.doe@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert a user
	userID, err := insertUser(db, user)
	if err != nil {
		log.Fatalf("Failed to insert user: %v", err)
	}
	log.Printf("Inserted user with ID: %d", userID)

	// Get the user by ID
	retrievedUser, err := getUserByID(db, userID)
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	log.Printf("Retrieved user: %+v", retrievedUser)

	// Update the user
	retrievedUser.Email = "john.doe.updated@example.com"
	retrievedUser.UpdatedAt = time.Now()
	err = updateUser(db, retrievedUser)
	if err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}
	log.Println("Updated user successfully")

	// Get all users
	users, err := getAllUsers(db)
	if err != nil {
		log.Fatalf("Failed to get all users: %v", err)
	}
	log.Printf("Retrieved %d users", len(users))
	for _, u := range users {
		log.Printf("User: %+v", u)
	}

	// Delete the user
	err = deleteUser(db, userID)
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}
	log.Println("Deleted user successfully")

	// Verify the user was deleted
	users, err = getAllUsers(db)
	if err != nil {
		log.Fatalf("Failed to get all users: %v", err)
	}
	log.Printf("Retrieved %d users after deletion", len(users))
}

// createUsersTable creates the users table if it doesn't exist
func createUsersTable(db *sqlx.DB) error {
	// Using raw SQL for table creation as fluentsql is primarily for DML operations
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		email VARCHAR(100) NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);`

	_, err := db.Exec(createTableSQL)
	return err
}

// insertUser inserts a new user into the database
func insertUser(db *sqlx.DB, user User) (int, error) {
	// Build the INSERT query using fluentsql
	query := qb.InsertInstance().
		Insert("users", "username", "email", "created_at", "updated_at").
		Row(user.Username, user.Email, user.CreatedAt, user.UpdatedAt)

	// Get the SQL and arguments
	sql, args, err := query.StringArgs([]any{})
	if err != nil {
		return 0, fmt.Errorf("failed to build insert query: %w", err)
	}

	// Add RETURNING clause for PostgreSQL to get the inserted ID
	sql = sql + " RETURNING id"

	// Execute the query
	var id int
	err = db.QueryRow(sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to execute insert query: %w", err)
	}

	return id, nil
}

// getUserByID retrieves a user by ID
func getUserByID(db *sqlx.DB, id int) (User, error) {
	// Build the SELECT query using fluentsql
	query := qb.QueryInstance().
		Select("id", "username", "email", "created_at", "updated_at").
		From("users").
		Where("id", qb.Eq, id)

	// Get the SQL and arguments
	sql, args, err := query.StringArgs([]any{})
	if err != nil {
		return User{}, fmt.Errorf("failed to build select query: %w", err)
	}

	// Execute the query
	var user User
	err = db.Get(&user, sql, args...)
	if err != nil {
		return User{}, fmt.Errorf("failed to execute select query: %w", err)
	}

	return user, nil
}

// getAllUsers retrieves all users from the database
func getAllUsers(db *sqlx.DB) ([]User, error) {
	// Build the SELECT query using fluentsql
	query := qb.QueryInstance().
		Select("id", "username", "email", "created_at", "updated_at").
		From("users").
		OrderBy("id", qb.Asc)

	// Get the SQL and arguments
	sql, args, err := query.StringArgs([]any{})
	if err != nil {
		return nil, fmt.Errorf("failed to build select all query: %w", err)
	}

	// Execute the query
	var users []User
	err = db.Select(&users, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select all query: %w", err)
	}

	return users, nil
}

// updateUser updates an existing user in the database
func updateUser(db *sqlx.DB, user User) error {
	// Build the UPDATE query using fluentsql
	query := qb.UpdateInstance().
		Update("users").
		Set("username", user.Username).
		Set("email", user.Email).
		Set("updated_at", user.UpdatedAt).
		Where("id", qb.Eq, user.ID)

	// Get the SQL and arguments
	sql, args, err := query.StringArgs()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	// Execute the query
	_, err = db.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	return nil
}

// deleteUser deletes a user from the database
func deleteUser(db *sqlx.DB, id int) error {
	// Build the DELETE query using fluentsql
	query := qb.DeleteInstance().
		Delete("users").
		Where("id", qb.Eq, id)

	// Get the SQL and arguments
	sql, args, err := query.StringArgs([]any{})
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	// Execute the query
	_, err = db.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	return nil
}
