package main

import (
	"database/sql"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/aidenappl/RootedGeocoder/env"
	"github.com/aidenappl/RootedGeocoder/geocoder"
	"github.com/aidenappl/RootedGeocoder/structs"

	_ "github.com/lib/pq"
	"github.com/schollz/progressbar/v3"
)

func main() {

	db, err := sql.Open("postgres", "postgres://postgres:"+env.DBPassword+"@"+env.DBHost+":5432/rooted?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if the connection is established
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	} else {
		log.Println("Successfully connected to the database")
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.Select("COUNT(*)").
		From("website.organisations o").
		Join("website.organisation_locations ol ON ol.organisation_id = o.id").
		Where(sq.Eq{"ol.lat": nil}).
		Where(sq.Eq{"ol.lng": nil}).
		Where(sq.NotEq{"ol.address_line_1": ""}).
		Where(sq.NotEq{"ol.city": ""}).
		Where(sq.NotEq{"ol.state": ""}).
		Where(sq.NotEq{"ol.zip_code": ""}).
		ToSql()
	if err != nil {
		log.Fatal("Failed to build SQL query:", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal("Failed to execute query:", err)
	}
	defer rows.Close()
	var count int
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal("Failed to scan result:", err)
		}
		log.Printf("Total organisations: %d\n", count)
	} else {
		log.Println("No organisations found")
	}

	// Prompt user to confirm the query
	fmt.Println()
	fmt.Printf("[??] Are you sure you want to run this query with %d rows? (yes/no)\n", count)
	var response string
	fmt.Scanln(&response)
	if response != "yes" {
		log.Println("Query execution cancelled by user.")
		return
	}
	fmt.Println()

	log.Println("Working on all valid orgs...")

	bar := progressbar.Default(int64(count))

	// paginate through the results
	offset := 0
	limit := 100 // Adjust the limit as needed
	for {
		query, args, err = psql.
			Select(
				"o.id",
				"o.name",
				"ol.id as location_id",
				"ol.address_line_1",
				"ol.city",
				"ol.state",
				"ol.zip_code",
			).
			From("website.organisations o").
			Join("website.organisation_locations ol ON ol.organisation_id = o.id").
			Where(sq.Eq{"ol.lat": nil}).
			Where(sq.Eq{"ol.lng": nil}).
			Where(sq.NotEq{"ol.address_line_1": ""}).
			Where(sq.NotEq{"ol.city": ""}).
			Where(sq.NotEq{"ol.state": ""}).
			Where(sq.NotEq{"ol.zip_code": ""}).
			OrderBy("o.id").
			Limit(uint64(limit)).
			Offset(uint64(offset)).
			ToSql()
		if err != nil {
			log.Fatal("Failed to build SQL query:", err)
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			log.Fatal("Failed to execute query:", err)
		}

		var org structs.Organisation
		for rows.Next() {
			if err := rows.Scan(
				&org.ID,
				&org.Name,
				&org.LocationID,
				&org.AddressLine1,
				&org.City,
				&org.State,
				&org.ZipCode,
			); err != nil {
				log.Fatal("Failed to scan result:", err)
			}

			resp, err := geocoder.GetGeocodedAddress(org)
			if err != nil {
				log.Printf("Error geocoding organisation %d (%s): %v\n", org.ID, org.Name, err)
				continue // Skip to the next organisation if there's an error
			}
			if resp == nil {
				log.Printf("No geocoded address found for organisation %d (%s)\n", org.ID, org.Name)
				continue // Skip to the next organisation if no address is found
			}
			// Update the organisation with the geocoded address
			q, a, err := psql.Update("website.organisation_locations l").
				Set("lat", resp.Features[0].Geometry.Coordinates[1]).
				Set("lng", resp.Features[0].Geometry.Coordinates[0]).
				Where(sq.Eq{"l.id": org.LocationID}).
				ToSql()
			if err != nil {
				log.Printf("Failed to build update query for organisation %d (%s): %v\n", org.ID, org.Name, err)
				continue // Skip to the next organisation if there's an error
			}

			_, err = db.Exec(q, a...)
			if err != nil {
				log.Printf("Failed to update organisation %d (%s): %v\n", org.ID, org.Name, err)
				continue // Skip to the next organisation if there's an error
			}
			progressbar.Bprintf(bar, "Geocoded organisation %d (%s)\n", org.ID, org.Name)
			bar.Add(1)
		}

		if err := rows.Err(); err != nil {
			log.Fatal("Error iterating over rows:", err)
		}

		rows.Close()

		if offset+limit >= count {
			break
		}

		progressbar.Bprintf(bar, "[Update] Processed %d organisations so far...\n", offset+limit)
		offset += limit
	}

}
