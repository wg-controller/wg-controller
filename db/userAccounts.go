package db

import (
	"github.com/lampy255/net-tbm/types"
)

func GetAccounts() ([]types.UserAccount, error) {
	// Query the database
	query := `SELECT
		email,
		role,
		failed_attempts,
		FROM user_accounts`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	// Loop through the rows
	var accounts []types.UserAccount
	for rows.Next() {
		var account types.UserAccount
		err = rows.Scan(
			&account.Email,
			&account.Role,
			&account.FailedAttempts,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func GetAccount(email string) (types.UserAccount, error) {
	// Query the database
	query := `SELECT
		email,
		role,
		failed_attempts,
		FROM user_accounts
		WHERE email = ?`
	row := DB.QueryRow(query, email)

	// Scan the row
	var account types.UserAccount
	err := row.Scan(
		&account.Email,
		&account.Role,
		&account.FailedAttempts,
	)
	if err != nil {
		return types.UserAccount{}, err
	}

	return account, nil
}

func InsertAccount(email string, role string, passwordHash string, passwordSalt string) error {
	query := `INSERT INTO user_accounts (
		email,
		role,
		failed_attempts,
		password_hash,
		password_salt
	) VALUES (?, ?, ?, ?, ?, ?)`

	_, err := DB.Exec(query,
		email,
		role,
		0,
		passwordHash,
		passwordSalt,
	)
	return err
}

// Delete an account from the database
func DeleteAccount(email string) error {
	query := `DELETE FROM user_accounts WHERE email = ?`
	_, err := DB.Exec(query, email)
	return err
}

// Delete all admin accounts from the database
func DeleteAdminAccounts() error {
	query := `DELETE FROM user_accounts WHERE role = 'admin'`
	_, err := DB.Exec(query)
	return err
}

func GetAccountPasswordHash(email string) (hash string, salt string, err error) {
	query := `SELECT
		password_hash,
		password_salt
		FROM user_accounts
		WHERE email = ?`
	row := DB.QueryRow(query, email)

	// Scan the row
	err = row.Scan(&hash, &salt)
	return
}

func UpdateAccountPasswordHash(email string, hash string, salt string) error {
	query := `UPDATE user_accounts SET
		password_hash = ?,
		password_salt = ?
		WHERE email = ?`

	_, err := DB.Exec(query, hash, salt, email)
	return err
}

func IncrementAccountFailedAttempts(email string) error {
	query := `UPDATE user_accounts SET
		failed_attempts = failed_attempts + 1
		WHERE email = ?`

	_, err := DB.Exec(query, email)
	return err
}

func ResetAccountFailedAttempts(email string) error {
	query := `UPDATE user_accounts SET
		failed_attempts = 0
		WHERE email = ?`

	_, err := DB.Exec(query, email)
	return err
}
