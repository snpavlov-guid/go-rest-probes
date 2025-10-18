package repo

import (
	"fmt"
	"database/sql"
	"github.com/snpavlov/app_aircraft/internal/model"
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


func executeRowsQueryAsync[T any](db *sql.DB, query string, args []interface{}, 
    scanFn func(*sql.Rows) (T, error)) <-chan model.ChannelListResult[T] {
	
	resultChan := make(chan model.ChannelListResult[T], 1)

	go func() {
        defer close(resultChan)

		rows, err := db.Query(query, args...)
		if err != nil {
			resultChan <- model.ChannelListResult[T]{Error: fmt.Errorf("ошибка выполнения запроса: %w", err)}
			return
		}
		defer rows.Close()

		var items []T
		for rows.Next() {
			item, err := scanFn(rows)
			if err != nil {
				resultChan <- model.ChannelListResult[T]{Error: fmt.Errorf("ошибка сканирования строки: %w", err)}
			}
			items = append(items, item)
		}
		resultChan <- model.ChannelListResult[T]{Items: &items}

		if err := rows.Err(); err != nil {
			resultChan <- model.ChannelListResult[T]{Error: fmt.Errorf("ошибка при обработке результатов: %w", err)}
		}

	}()
    
    return resultChan
}

func executeRowQueryAsync[T any](db *sql.DB, query string, args []interface{}, 
    scanFn func(*sql.Row) (T, error)) <-chan model.ChannelItemResult[T] {

	resultChan := make(chan model.ChannelItemResult[T], 1)

	go func() {
		var item T
		row := db.QueryRow(query, args...)
		item, err := scanFn(row)
		if err != nil {
			resultChan <- model.ChannelItemResult[T]{Error: fmt.Errorf("ошибка сканирования строки: %w", err)}
		}
		resultChan <- model.ChannelItemResult[T]{Item: &item }

	}()
    
    return resultChan
}

