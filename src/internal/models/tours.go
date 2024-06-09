package models

import (
	"database/sql"
	"errors"
)

type Tour struct {
	ID          string
	RecordType  string
	Name        string
	Abstract    string
	LogoURL     string
	URLswmo     string
	RecordCount int
}

type TourPicture struct {
	ID         string
	PictureURL string
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type TourModels struct {
	DB *sql.DB
}

// TODO:
// Review the Tour functions and struct
// could be simplified to basic and rich view (for function and structs)

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
				t.ID::varchar as ID, 
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
		err := rows.Scan(&t.ID, &t.RecordType, &t.Name, &t.Abstract, &t.LogoURL, &t.URLswmo, &t.RecordCount)
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

func (m *TourModels) RandomTourPics() ([]*TourPicture, error) {

	// initialise empty slice of pointers
	pics := []*TourPicture{}

	// stmt := `select url from gold.tour_images using sample $1`
	stmt := `select ID::varchar as ID, url from gold.tour_images using sample 3`

	// rows, err := m.DB.Query(stmt, sampleSize)
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// close before the method SearchTour returns
	// should be after we check for an error, otherwise get panic trying to close nil rows
	defer rows.Close()

	for rows.Next() {
		t := &TourPicture{} //pointer to TourPicture
		err := rows.Scan(&t.ID, &t.PictureURL)
		if err != nil {
			return nil, err
		}
		pics = append(pics, t)
	}
	// important to collect errors after the iterations
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pics, nil
}

func (m *TourModels) TourBasicInfo(id string) (*Tour, error) {

	// initialise pointer to the struct
	t := &Tour{}

	stmt := `
		select 
			t.ID::varchar as ID, 
			t.record_type,
			t.name, 
			t.abstract,
			t.logo_url,
			t.url_swmo,
			1 as cnt
		from gold.tours t
		where t.ID::varchar = $1
	`

	row := m.DB.QueryRow(stmt, id)

	err := row.Scan(&t.ID, &t.RecordType, &t.Name, &t.Abstract, &t.LogoURL, &t.URLswmo, &t.RecordCount)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return t, nil
}
