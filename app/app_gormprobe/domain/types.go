package domain

import (
	"fmt"
	"strings"
	"strconv"
	"database/sql/driver"
)

type Point struct {
	X float64
	Y float64
}

// Implement the Valuer interface for saving to the database
func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("(%f,%f)", p.X, p.Y), nil
}

// Implement the Scanner interface for reading from the database
func (p *Point) Scan(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("could not scan type %T into Point", value)
	}

	s = strings.Trim(s, "()")
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return fmt.Errorf("invalid point format: %s", s)
	}

	x, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return fmt.Errorf("invalid X coordinate: %w", err)
	}
	y, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return fmt.Errorf("invalid Y coordinate: %w", err)
	}

	p.X = x
	p.Y = y
	return nil
}
