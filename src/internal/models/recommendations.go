package models

import (
	"database/sql"
)

type Recommendation struct {
	ID          int
	Title       string
	Description string
	FieldErrors map[string]string
}

// Define a ToursModel type which wraps a sql.DB connection pool.
type RecModels struct {
	DB *sql.DB
}

func (m *RecModels) Insert(title string, description string) (int, error) {

	stmt := `
		insert into recommendations.tours (title, description, created_at) 
		values($1, $2, current_timestamp) RETURNING id;
	`

	var id int
	err := m.DB.QueryRow(stmt, title, description).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

func (m *RecModels) GetAll() ([]*Recommendation, error) {

	// initialise empty slice of pointers
	recommendations := []*Recommendation{}

	// probably want some paging, but at this point it is fine
	stmt := `select id, title, description from recommendations.tours`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// close before the method returns
	// should be after we check for an error, otherwise get panic trying to close nil rows
	defer rows.Close()

	for rows.Next() {
		r := &Recommendation{} //pointer to Recommendation
		err := rows.Scan(&r.ID, &r.Title, &r.Description)
		if err != nil {
			return nil, err
		}
		recommendations = append(recommendations, r)
	}
	// important to collect errors after the iterations
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return recommendations, nil
}
