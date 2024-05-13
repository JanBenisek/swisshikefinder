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
type TourModels struct {
	DB *sql.DB
}

func (m *TourModels) SearchTour(query string) ([]*Tour, error) {

	// initialise empty slice of pointers
	tours := []*Tour{}

	stmt := `select 
		t.ID, 
		t.record_type,
		t.name
	from gold.tours t
	where t.ID in (select ID from gold.tour_itinerary where name ilike $1)
	limit 10`

	rows, err := m.DB.Query(stmt, query)
	if err != nil {
		return nil, err
	}
	// close before the method SearchTour returns
	// should be after we check for an error, otherwise get panic trying to close nil rows
	defer rows.Close()

	for rows.Next() {
		t := &Tour{} //pointer to Tour
		err := rows.Scan(&t.ID, &t.Record_type, &t.Name)
		if err != nil {
			return nil, err
		}
		tours = append(tours, t)
	}
	// important to collect errors after the iterations
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tours, nil
}

func (m *TourModels) GetTour(n_rows int) (*Tour, error) {

	// initialise pointer to the new struct
	t := &Tour{}

	// row := m.DB.QueryRow("select ID, name, record_type from './data/gold_tours.parquet' limit ?", n_rows)
	row := m.DB.QueryRow("select ID, name, record_type from gold.tours limit ?", n_rows)
	// note that we are passing pointers
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
