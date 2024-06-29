package db

import (
	"strconv"
	"strings"

	"github.com/altugbakan/card-logger/utils"
)

type Set struct {
	Abbr       string
	Name       string
	TotalCards int
}

type UserCard struct {
	CardID   int
	Number   int
	Name     string
	Patterns []utils.Pattern
}

func GetSet(abbr string) (Set, error) {
	abbr = strings.ToUpper(abbr)
	query := `SELECT abbr, name, total_cards FROM sets WHERE abbr = ?`
	row := db.QueryRow(query, abbr)

	var set Set
	err := row.Scan(&set.Abbr, &set.Name, &set.TotalCards)

	return set, err
}

func GetAllSets() ([]Set, error) {
	query := `SELECT abbr, name, total_cards FROM sets`
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

func GetAllSetCardsOfUser(abbr string) ([]UserCard, error) {
	query := `
	WITH card_pattern_quantities AS (
		SELECT
			cards.number AS card_number,
			cards.name AS card_name,
			cards.id AS card_id,
			rp.pattern || ':' || COALESCE(uc.quantity, '0') AS pattern_quantity
		FROM
			cards
			JOIN sets ON cards.set_abbr = sets.abbr
			JOIN rarity_patterns rp ON cards.set_abbr = rp.set_abbr AND cards.rarity = rp.rarity
			LEFT JOIN user_cards uc ON cards.id = uc.card_id AND rp.pattern = uc.pattern
		WHERE
			sets.abbr = ?
		)
	SELECT 
		card_id,
		card_number,
		card_name,
		GROUP_CONCAT(pattern_quantity) AS patterns
	FROM 
		card_pattern_quantities
	GROUP BY 
		card_number, card_name
	ORDER BY 
		card_number;
	`

	rows, err := db.Query(query, abbr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := []UserCard{}
	for rows.Next() {
		var card UserCard
		var patterns string
		err := rows.Scan(&card.CardID, &card.Number, &card.Name, &patterns)
		if err != nil {
			return nil, err
		}

		patternsSplit := strings.Split(patterns, ",")
		for _, pattern := range patternsSplit {
			patternSplit := strings.Split(pattern, ":")
			quantity, err := strconv.Atoi(patternSplit[1])
			if err != nil {
				return nil, err
			}

			card.Patterns = append(card.Patterns, utils.Pattern{
				Name:     patternSplit[0],
				Quantity: quantity,
			})
		}

		cards = append(cards, card)
	}

	return cards, nil
}
