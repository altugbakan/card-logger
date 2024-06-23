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

func GetCard(set string, number int) (Card, error) {
	query := `SELECT id, name, set_abbr, number, rarity FROM Cards WHERE set_abbr = ? AND number = ?`
	row := db.QueryRow(query, set, number)

	var card Card
	err := row.Scan(&card.ID, &card.Name, &card.Set, &card.Number, &card.Rarity)
	if err != nil {
		if err == sql.ErrNoRows {
			return Card{}, fmt.Errorf("no card found with set %s and card number %d", set, number)
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
	// if user already has the card id with the pattern, increase quantity
	query := `SELECT card_id FROM UserCards WHERE card_id = ? AND pattern = ?`
	row := db.QueryRow(query, cardID, pattern)

	var userCardID int
	err := row.Scan(&userCardID)
	if err != nil {
		// if user doesn't have the card id with the pattern, insert a new row
		query = `INSERT INTO UserCards (card_id, quantity, pattern) VALUES (?, 1, ?)`
		_, err = db.Exec(query, cardID, pattern)
		if err != nil {
			return err
		}
	} else {
		// if user already has the card id with the pattern, increase quantity
		query = `UPDATE UserCards SET quantity = quantity + 1 WHERE card_id = ? AND pattern = ?`
		_, err = db.Exec(query, userCardID, pattern)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveUserCard(cardID int, pattern string) error {
	// if user already has the card id with the pattern, decrease quantity
	query := `SELECT card_id FROM UserCards WHERE card_id = ? AND pattern = ?`
	row := db.QueryRow(query, cardID, pattern)

	var userCardID int
	err := row.Scan(&userCardID)
	if err != nil {
		return fmt.Errorf("no existing card found with card id %d with pattern %s", cardID, pattern)
	}

	// if user already has the card id with the pattern, decrease quantity
	query = `UPDATE UserCards SET quantity = quantity - 1 WHERE card_id = ? AND pattern = ?`
	_, err = db.Exec(query, userCardID, pattern)
	if err != nil {
		return err
	}

	// if quantity is 0, remove the row
	query = `DELETE FROM UserCards WHERE card_id = ? AND pattern = ? AND quantity = 0`
	_, err = db.Exec(query, userCardID, pattern)
	if err != nil {
		return err
	}

	return nil
}

func GetUserCardsForSet(set string) ([]Card, error) {
	query := `SELECT id, name, set_abbr, number, rarity FROM Cards WHERE set_abbr = ?`
	rows, err := db.Query(query, set)
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
