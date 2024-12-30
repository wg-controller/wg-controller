package db

import "github.com/lampy255/net-tbm/types"

func GetApiKey(hash string) (expiresUnixMillis int64, err error) {
	// Query the database
	query := `SELECT
		expires_unix_millis
		FROM api_keys
		WHERE hash = ?`
	row := DB.QueryRow(query, hash)

	// Scan the row
	err = row.Scan(&expiresUnixMillis)
	if err != nil {
		return 0, err
	}

	return expiresUnixMillis, nil
}

func GetApiKeys() ([]types.APIKey, error) {
	// Query the database
	query := `SELECT
		uuid,
		expires_unix_millis,
		name
		FROM api_keys`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	// Loop through the rows
	var keys []types.APIKey
	for rows.Next() {
		var key types.APIKey
		err = rows.Scan(
			&key.UUID,
			&key.ExpiresUnixMillis,
			&key.Name,
		)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return keys, nil
}

func InsertApiKey(key types.APIKey) error {
	// Insert the session into the database
	query := `INSERT INTO api_keys
		(uuid, hash, expires_unix_millis, name)
		VALUES (?, ?, ?, ?)`
	_, err := DB.Exec(query, key.UUID, key.Hash, key.ExpiresUnixMillis, key.Name)
	return err
}

func DeleteApiKey(uuid string) error {
	// Delete the api key from the database
	query := `DELETE FROM api_keys WHERE uuid = ?`
	_, err := DB.Exec(query, uuid)
	return err
}
