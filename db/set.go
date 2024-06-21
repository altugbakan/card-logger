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
