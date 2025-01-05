package db

import (
	"strings"

	"github.com/lampy255/net-tbm/types"
)

func GetApiKey(hash []byte) (expiresUnixMillis int64, attributes []string, err error) {
	// Query the database
	query := `SELECT
		expires_unixmillis,
		attributes
		FROM api_keys
		WHERE hash = ?`
	row := DB.QueryRow(query, hash)

	// Scan the row
	attributesString := ""
	err = row.Scan(&expiresUnixMillis, &attributesString)
	if err != nil {
		return 0, []string{}, err
	}

	// Split the attributes
	attributes = strings.Split(attributesString, ",")
	if len(attributes) == 1 {
		if attributes[0] == "" {
			attributes = []string{}
		}
	}

	return expiresUnixMillis, attributes, nil
}

func GetApiKeys() ([]types.APIKey, error) {
	// Query the database
	query := `SELECT
		uuid,
		name,
		expires_unixmillis,
		attributes
		FROM api_keys`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	// Loop through the rows
	var keys []types.APIKey
	for rows.Next() {
		var key types.APIKey
		var attributes string
		err = rows.Scan(
			&key.UUID,
			&key.Name,
			&key.ExpiresUnixMillis,
			&attributes,
		)
		if err != nil {
			return nil, err
		}

		// Split the attributes
		key.Attributes = strings.Split(attributes, ",")
		if len(key.Attributes) == 1 {
			if key.Attributes[0] == "" {
				key.Attributes = []string{}
			}
		}

		keys = append(keys, key)
	}

	return keys, nil
}

func InsertApiKey(key types.APIKey, hash []byte) error {
	// Insert the session into the database
	query := `INSERT INTO api_keys
		(uuid, name, expires_unixmillis, attributes, hash)
		VALUES (?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, key.UUID, key.Name, key.ExpiresUnixMillis, strings.Join(key.Attributes, ","), hash)
	return err
}

func UpdateApiKey(key types.APIKey) error {
	// Update the api key in the database
	query := `UPDATE api_keys
		SET name = ?, expires_unixmillis = ?, attributes = ?
		WHERE uuid = ?`
	_, err := DB.Exec(query, key.Name, key.ExpiresUnixMillis, strings.Join(key.Attributes, ","), key.UUID)
	return err
}

func DeleteApiKey(uuid string) error {
	// Delete the api key from the database
	query := `DELETE FROM api_keys WHERE uuid = ?`
	_, err := DB.Exec(query, uuid)
	return err
}
