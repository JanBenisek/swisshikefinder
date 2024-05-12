package models

import (
	"database/sql"
	"errors"
)

type Tour struct {
	ID          string
	Record_type string
	Name        string
	// Part_of_trip_id            string
	// Url_mysw                   string
	// Lat                        float64
	// Lon                        float64
	// Elevation                  int64
	// Logo_url                   string
	// Stage                      string
	// Requirements_technical     string
	// Requirements_endurance     string
	// Route_category             string
	// Url_swmo                   string
	// Season                     string
	// Distance                   int64
	// Duration                   int64
	// Duration_reverse_direction int64
	// Ascent                     int64
	// Descent                    int64
	// Barrier_free               bool
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type TourModel struct {
	DB *sql.DB
}

func (m *TourModel) GetTour(n_rows int) (*Tour, error) {

	t := &Tour{}

	// row := m.DB.QueryRow("select ID, name, record_type from './data/gold_tours.parquet' limit ?", n_rows)
	row := m.DB.QueryRow("select ID, name, record_type from gold.tours limit ?", n_rows)
	err := row.Scan(&t.ID, &t.Record_type, &t.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return t, nil
}
