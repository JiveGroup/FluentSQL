# FluentSQL Example

This example demonstrates how to use the [FluentSQL](https://github.com/jivegroup/fluentsql) library to perform CRUD (Create, Read, Update, Delete) operations with a PostgreSQL database.

## What This Example Demonstrates

This example shows:

1. How to connect to a PostgreSQL database using [sqlx](https://github.com/jmoiron/sqlx)
2. How to create a table (using raw SQL)
3. How to use FluentSQL to build and execute:
   - INSERT queries
   - SELECT queries (single record and multiple records)
   - UPDATE queries
   - DELETE queries
4. How to convert FluentSQL query builders to SQL strings and arguments
5. How to handle returned values from PostgreSQL (using the RETURNING clause)

## Prerequisites

To run this example, you need:

1. Go 1.24 or later
2. A PostgreSQL database server
3. The following environment setup:
   - PostgreSQL running on localhost:5432
   - A database named "test"
   - A user named "user" with password "secret" that has access to the "test" database

## How to Run

1. Make sure your PostgreSQL server is running and properly configured
2. Navigate to the example directory:
   ```bash
   cd example
   ```
3. Run the example:
   ```bash
   go run main.go
   ```

## Expected Output

When running the example successfully, you should see output similar to:

```
Starting application...
Connected to database successfully
Inserted user with ID: 1
Retrieved user: {ID:1 Username:johndoe Email:john.doe@example.com CreatedAt:2023-06-01 12:34:56 -0700 PDT UpdatedAt:2023-06-01 12:34:56 -0700 PDT}
Updated user successfully
Retrieved 1 users
User: {ID:1 Username:johndoe Email:john.doe.updated@example.com CreatedAt:2023-06-01 12:34:56 -0700 PDT UpdatedAt:2023-06-01 12:34:57 -0700 PDT}
Deleted user successfully
Retrieved 0 users after deletion
```

## Code Explanation

The example demonstrates a complete workflow:

1. **Database Connection**: Uses sqlx to connect to PostgreSQL
2. **Table Creation**: Creates a users table if it doesn't exist
3. **Insert Operation**: Adds a new user and retrieves the generated ID
4. **Select Operation**: Retrieves the user by ID to verify insertion
5. **Update Operation**: Modifies the user's email address
6. **Select All Operation**: Retrieves all users to verify the update
7. **Delete Operation**: Removes the user from the database
8. **Verification**: Retrieves all users again to verify deletion

Each database operation function demonstrates how to:
1. Build a query using FluentSQL's fluent interface
2. Convert the query to SQL and arguments
3. Execute the query using sqlx
4. Handle the results

## Customizing the Example

To use this example with your own database:

1. Modify the database connection parameters in `main.go`:
   ```
   const (
       dbUser = "your_username"
       dbPass = "your_password"
       dbName = "your_database"
   )
   ```

2. If needed, modify the connection string to include a different host or port:
   ```
   connStr := fmt.Sprintf("postgres://%s:%s@your_host:your_port/%s?sslmode=disable", dbUser, dbPass, dbName)
   ```

## Additional Resources

- [FluentSQL GitHub Repository](https://github.com/jivegroup/fluentsql)
- [sqlx Documentation](https://github.com/jmoiron/sqlx)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
