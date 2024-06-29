package db

func IsPatternValidForRarity(setAbbr, rarity, pattern string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM rarity_patterns WHERE set_abbr = ? AND rarity = ? AND pattern = ?)`
	err := db.QueryRow(query, setAbbr, rarity, pattern).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func GetPatternsForRarity(setAbbr, rarity string) ([]string, error) {
	patterns := []string{}
	query := `SELECT pattern FROM rarity_patterns WHERE set_abbr = ? AND rarity = ?`
	rows, err := db.Query(query, setAbbr, rarity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pattern string
		if err := rows.Scan(&pattern); err != nil {
			return nil, err
		}
		patterns = append(patterns, pattern)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return patterns, nil
}
