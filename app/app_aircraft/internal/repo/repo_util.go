package repo

import (
	"fmt"
	"database/sql"
)

func executeRowsQuery[T any](db *sql.DB, query string, args []interface{}, 
    scanFn func(*sql.Rows) (T, error)) ([]T, error) {
	
    rows, err := db.Query(query, args...)
    if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

    var items []T
	for rows.Next() {
		//var item *T
		item, err := scanFn(rows)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
	}
    
    return items, nil
}

func executeRowQuery[T any](db *sql.DB, query string, args []interface{}, 
    scanFn func(*sql.Row) (T, error)) (*T, error) {
	
    var item T
    row := db.QueryRow(query, args...)
	item, err := scanFn(row)
    if err != nil {
        return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
    }
    
    return &item, nil
}
