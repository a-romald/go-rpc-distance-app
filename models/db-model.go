package models

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"
)

// BaseModel is the type for database connection
type BaseModel struct {
	DB *sql.DB
}

type GeoLocation struct {
	Id        int
	Point1    GeoPoint
	Point2    GeoPoint
	Distance  float64
	IpAddress string
	CreatedAt time.Time
}

type GeoPoint struct {
	Lat float64
	Lng float64
}

type Pagination struct {
	Next        int
	Previous    int
	PerPage     int
	CurrentPage int
	TotalPage   int
	AllPages    []int
	Sort        string
}

func (m *BaseModel) GetLocation(gl GeoLocation) (GeoLocation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var loc GeoLocation

	query := `
		SELECT
			id 
		FROM 
			location
		WHERE 
			point1 = POINT(?, ?) AND point2 = POINT(?, ?) 
		LIMIT 1
		`

	row := m.DB.QueryRowContext(ctx, query, gl.Point1.Lat, gl.Point1.Lng, gl.Point2.Lat, gl.Point2.Lng)
	// Error if not exists: sql: no rows in result set
	err := row.Scan(
		&loc.Id,
	)
	if err != nil {
		return loc, err
	}
	return loc, nil
}

// InsertLocation inserts a new geoLocation loc, and returns its id
func (m *BaseModel) InsertLocation(loc GeoLocation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO location
			(point1, point2, distance, ip_address, created_at)
		VALUES (POINT(?, ?), POINT(?, ?), ?, ?, ?)		
	`

	result, err := m.DB.ExecContext(ctx, stmt,
		loc.Point1.Lat,
		loc.Point1.Lng,
		loc.Point2.Lat,
		loc.Point2.Lng,
		loc.Distance,
		loc.IpAddress,
		loc.CreatedAt,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *BaseModel) GetAllLocations(limit int, page int, sort_by string) ([]*GeoLocation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if page == 0 {
		page = 1
	}
	offset := limit * (page - 1)

	var locations []*GeoLocation

	// Sorting data
	var sort_order string
	var sort_field string
	if sort_by == "" {
		sort_by = os.Getenv("SORT_BY")
	}
	if strings.HasPrefix(sort_by, "-") {
		sort_order = "DESC"
		sort_field = strings.TrimPrefix(sort_by, "-")
	} else {
		sort_order = "ASC"
		sort_field = sort_by
	}

	query := `
		SELECT
			id, 
			ST_X(point1) AS latitude1, ST_Y(point1) AS longitude1, 
			ST_X(point2) AS latitude2, ST_Y(point2) AS longitude2, 
			distance, ip_address, created_at
		FROM 
			location
		ORDER BY
			%s %s
		LIMIT ? OFFSET ?
		`

	query = fmt.Sprintf(query, sort_field, sort_order)

	rows, err := m.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var loc GeoLocation
		err = rows.Scan(
			&loc.Id,
			&loc.Point1.Lat,
			&loc.Point1.Lng,
			&loc.Point2.Lat,
			&loc.Point2.Lng,
			&loc.Distance,
			&loc.IpAddress,
			&loc.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		locations = append(locations, &loc)
	}

	return locations, nil
}

func (m *BaseModel) Paginate(table string, limit int, page int, sort_by string) *Pagination {
	var (
		tmpl  = Pagination{}
		count int
	)
	if page == 0 {
		page = 1
	}
	// Count all record
	sqltable := fmt.Sprintf("SELECT count(id) FROM %s", table)

	m.DB.QueryRow(sqltable).Scan(&count)

	total := (count / limit)

	// Calculator Total Page
	remainder := (count % limit)
	if remainder == 0 {
		tmpl.TotalPage = total
	} else {
		tmpl.TotalPage = total + 1
	}

	// Set current/record per page
	tmpl.CurrentPage = page
	tmpl.PerPage = limit

	// Calculator the Next/Previous Page
	if page <= 0 {
		tmpl.Next = page + 1
	} else if page < tmpl.TotalPage {
		tmpl.Previous = page - 1
		tmpl.Next = page + 1
	} else if page == tmpl.TotalPage {
		tmpl.Previous = page - 1
		tmpl.Next = 0
	}

	// Slice from all pages
	totalNum := total + 1
	pagesSlice := []int{}
	for i := 1; i <= totalNum; i++ {
		pagesSlice = append(pagesSlice, i)
	}

	tmpl.AllPages = pagesSlice
	tmpl.Sort = sort_by

	return &tmpl
}
