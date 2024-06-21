package db

import (
	"database/sql"
	"fmt"
)

type Card struct {
	ID     int
	Set    string
	Number int
	Rarity string
}

func GetCard(set string, number int) (Card, error) {
	query := `SELECT id, set_abbr, number, rarity FROM Cards WHERE set_abbr = ? AND number = ?`
	row := db.QueryRow(query, set, number)

	var card Card
	err := row.Scan(&card.ID, &card.Set, &card.Number, &card.Rarity)
	if err != nil {
		if err == sql.ErrNoRows {
			return Card{}, fmt.Errorf("no card found with set %s and card number %d", set, number)
		}
		return Card{}, err
	}

	return card, nil
}

func AddUserCard(cardID int, pattern string) error {
	// if user already has the card id with the pattern, increase quantity
	query := `SELECT id FROM UserCards WHERE card_id = ? AND pattern = ?`
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
		query = `UPDATE UserCards SET quantity = quantity + 1 WHERE id = ?`
		_, err = db.Exec(query, userCardID)
		if err != nil {
			return err
		}
	}

	return nil
}
