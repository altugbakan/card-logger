package db

import (
	"database/sql"
	"fmt"
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

func GetAllSetCardsOfUser(abbr string) ([]UserCard, error) {
	query := `
	WITH CardPatternQuantities AS (
		SELECT
			Cards.number AS CardNumber,
			Cards.name AS CardName,
			Cards.id AS CardID,
			RP.pattern || ':' || COALESCE(UC.quantity, '0') AS PatternQuantity
		FROM
			Cards
			JOIN Sets ON Cards.set_abbr = Sets.abbr
			JOIN RarityPatterns RP ON Cards.set_abbr = RP.set_abbr AND Cards.rarity = RP.rarity
			LEFT JOIN UserCards UC ON Cards.id = UC.card_id AND RP.pattern = UC.pattern
		WHERE
			Sets.abbr = ?
		)
	SELECT 
		CardID,
		CardNumber,
		CardName,
		GROUP_CONCAT(PatternQuantity) AS Patterns
	FROM 
		CardPatternQuantities
	GROUP BY 
		CardNumber, CardName
	ORDER BY 
		CardNumber;`

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
