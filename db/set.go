package db

import (
	"database/sql"
	"fmt"
	"strings"
)

type Set struct {
	Abbr       string
	Name       string
	TotalCards int
}

func GetSet(abbr string) (Set, error) {
	abbr = strings.ToUpper(abbr)
	query := `SELECT abbr, name, total_cards FROM Sets WHERE abbr = ?`
	row := db.QueryRow(query, abbr)

	var set Set
	err := row.Scan(&set.Abbr, &set.Name, &set.TotalCards)
	if err != nil {
		if err == sql.ErrNoRows {
			return Set{}, fmt.Errorf("no set found with abbreviation %s", abbr)
		}
		return Set{}, err
	}

	return set, nil
}

func GetAllSets() ([]Set, error) {
	query := `SELECT abbr, name, total_cards FROM Sets`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sets := []Set{}
	for rows.Next() {
		var set Set
		err := rows.Scan(&set.Abbr, &set.Name, &set.TotalCards)
		if err != nil {
			return nil, err
		}

		sets = append(sets, set)
	}

	return sets, nil
}

func GetUserCardCountsBySet() (map[string]int, error) {
	query := `
    SELECT c.set_abbr, COUNT(DISTINCT uc.card_id) 
    FROM UserCards uc
    JOIN Cards c ON uc.card_id = c.id
    GROUP BY c.set_abbr
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	setCounts := make(map[string]int)

	for rows.Next() {
		var setAbbr string
		var count int

		err := rows.Scan(&setAbbr, &count)
		if err != nil {
			return nil, err
		}

		setCounts[setAbbr] = count
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return setCounts, nil
}
