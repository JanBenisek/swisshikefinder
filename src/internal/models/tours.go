package models

import (
	"database/sql"
	"errors"
)

type Tour struct {
	ID           string
	Record_type  string
	Name         string
	Abstract     string
	Logo_url     string
	Url_swmo     string
	Record_count int
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type TourModels struct {
	DB *sql.DB
}

func (m *TourModels) SearchTour(query string, limit int, offset int) ([]*Tour, error) {

	// initialise empty slice of pointers
	tours := []*Tour{}

	// using offset might not be the efficient, but who cares with this tiny dataset
	// another option is to have auto-increment id
	// also I need round-trip for the count, hacking it into one query for now
	stmt := `
		with 
		search_q as (
			select 
				t.ID, 
				t.record_type,
				t.name, 
				t.abstract,
				t.logo_url,
				t.url_swmo
			from gold.tours t
			where t.ID in (select ID from gold.tour_itinerary where name ilike $1)
		), 
		count_q as (
			select count(*) as cnt from search_q
		)
		select
			s.*,
			c.cnt
		from search_q s
		cross join count_q c
		limit $2
		offset $3
	`

	rows, err := m.DB.Query(stmt, query, limit, offset)
	if err != nil {
		return nil, err
	}
	// close before the method SearchTour returns
	// should be after we check for an error, otherwise get panic trying to close nil rows
	defer rows.Close()

	for rows.Next() {
		t := &Tour{} //pointer to Tour
		err := rows.Scan(&t.ID, &t.Record_type, &t.Name, &t.Abstract, &t.Logo_url, &t.Url_swmo, &t.Record_count)
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
