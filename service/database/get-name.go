package database

import (
	"fmt"
)

func (db *appdbimpl) GetDatabaseTableContent(tableName string) ([]map[string]interface{}, error) {
	// Prepare the query to fetch all rows from the given table
	query := fmt.Sprintf("SELECT * FROM %s", tableName)

	// Execute the query
	rows, err := db.c.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Retrieve column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Prepare a slice to hold the map for each row
	var results []map[string]interface{}

	// Iterate over the rows
	for rows.Next() {
		// Prepare a slice to hold the values of each row
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row values into the values slice
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Create a map to hold column name and value pairs
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col] = values[i]
		}

		// Append the row map to the results slice
		results = append(results, rowMap)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetName() (string, error) {
	var name string
	err := db.c.QueryRow("SELECT name FROM example_table WHERE id=1").Scan(&name)
	return name, err
}
