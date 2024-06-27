package db

import (
	"database/sql"
	"fmt"
)

type Card struct {
	ID     int
	Name   string
	Set    string
	Number int
	Rarity string
}

type NoPatternCard struct {
	Number int
	Name   string
	Rarity string
}

type IncompleteCard struct {
	Number   int
	Name     string
	Rarity   string
	Patterns []string
}

func GetCard(abbr string, number int) (Card, error) {
	query := `SELECT id, name, set_abbr, number, rarity FROM Cards WHERE set_abbr = ? AND number = ?`
	row := db.QueryRow(query, abbr, number)

	var card Card
	err := row.Scan(&card.ID, &card.Name, &card.Set, &card.Number, &card.Rarity)
	if err != nil {
		if err == sql.ErrNoRows {
			return Card{}, fmt.Errorf("no card found with set %s and card number %d", abbr, number)
		}
		return Card{}, err
	}

	return card, nil
}

func GetAllUserPatternAmounts(cardID int) (map[string]int, error) {
	query := `SELECT card_id, pattern, quantity FROM UserCards WHERE card_id = ?`
	rows, err := db.Query(query, cardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	patterns := make(map[string]int)
	for rows.Next() {
		var pattern string
		var quantity int
		err := rows.Scan(&cardID, &pattern, &quantity)
		if err != nil {
			return nil, err
		}

		patterns[pattern] = quantity
	}

	return patterns, nil
}

func AddUserCard(cardID int, pattern string) error {
	query := `SELECT card_id FROM UserCards WHERE card_id = ? AND pattern = ?`
	row := db.QueryRow(query, cardID, pattern)

	var userCardID int
	err := row.Scan(&userCardID)
	if err != nil {
		query = `INSERT INTO UserCards (card_id, quantity, pattern) VALUES (?, 1, ?)`
		_, err = db.Exec(query, cardID, pattern)
		if err != nil {
			return err
		}
	} else {
		query = `UPDATE UserCards SET quantity = quantity + 1 WHERE card_id = ? AND pattern = ?`
		_, err = db.Exec(query, userCardID, pattern)
		if err != nil {
			return err
		}
	}

	hasChanges = true
	return nil
}

func RemoveUserCard(cardID int, pattern string) error {
	query := `SELECT card_id FROM UserCards WHERE card_id = ? AND pattern = ?`
	row := db.QueryRow(query, cardID, pattern)

	var userCardID int
	err := row.Scan(&userCardID)
	if err != nil {
		return fmt.Errorf("no existing card found with card id %d with pattern %s", cardID, pattern)
	}

	query = `UPDATE UserCards SET quantity = quantity - 1 WHERE card_id = ? AND pattern = ?`
	_, err = db.Exec(query, userCardID, pattern)
	if err != nil {
		return err
	}

	query = `DELETE FROM UserCards WHERE card_id = ? AND pattern = ? AND quantity = 0`
	_, err = db.Exec(query, userCardID, pattern)
	if err != nil {
		return err
	}

	hasChanges = true
	return nil
}

func GetUserCardsForSet(abbr string) ([]Card, error) {
	query := `SELECT id, name, set_abbr, number, rarity FROM Cards WHERE set_abbr = ?`
	rows, err := db.Query(query, abbr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := []Card{}
	for rows.Next() {
		var card Card
		err := rows.Scan(&card.ID, &card.Name, &card.Set, &card.Number, &card.Rarity)
		if err != nil {
			return nil, err
		}

		cards = append(cards, card)
	}

	return cards, nil
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

func GetUserCardsWithNoPatternsForSet(abbr string) ([]NoPatternCard, error) {
	query := `
	SELECT c.name, c.number, c.rarity
	FROM Cards c
	LEFT JOIN UserCards uc ON c.id = uc.card_id
	WHERE c.set_abbr = ? AND uc.pattern IS NULL
	ORDER BY c.number
	`

	rows, err := db.Query(query, abbr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := []NoPatternCard{}
	for rows.Next() {
		var card NoPatternCard
		err := rows.Scan(&card.Name, &card.Number, &card.Rarity)
		if err != nil {
			return nil, err
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func GetUserCardsWithIncompletePatternsForSet(abbr string) ([]IncompleteCard, error) {
	query := `
    SELECT c.number, c.name, c.rarity, rp.pattern
    FROM Cards c
    JOIN RarityPatterns rp ON c.set_abbr = rp.set_abbr AND c.rarity = rp.rarity
    LEFT JOIN UserCards uc ON uc.card_id = c.id AND uc.pattern = rp.pattern
    WHERE c.set_abbr = ? AND uc.card_id IS NULL
    ORDER BY c.number, rp.pattern;
    `
	rows, err := db.Query(query, abbr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []IncompleteCard
	var currentCard *IncompleteCard
	for rows.Next() {
		var number int
		var name, rarity, pattern string
		if err := rows.Scan(&number, &name, &rarity, &pattern); err != nil {
			return nil, err
		}

		if currentCard == nil || currentCard.Number != number {
			if currentCard != nil {
				results = append(results, *currentCard)
			}

			currentCard = &IncompleteCard{
				Number:   number,
				Name:     name,
				Rarity:   rarity,
				Patterns: []string{pattern},
			}
		} else {
			currentCard.Patterns = append(currentCard.Patterns, pattern)
		}
	}

	if currentCard != nil {
		results = append(results, *currentCard)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
