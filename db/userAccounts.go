package db

import "github.com/lampy255/net-tbm/types"

func GetAccounts() ([]types.UserAccount, error) {
	// Query the database
	query := `SELECT
		email,
		role,
		failed_attempts,
		suspended
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
			&account.Suspended,
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
		suspended
		FROM user_accounts
		WHERE email = ?`
	row := DB.QueryRow(query, email)

	// Scan the row
	var account types.UserAccount
	err := row.Scan(
		&account.Email,
		&account.Role,
		&account.FailedAttempts,
		&account.Suspended,
	)
	if err != nil {
		return types.UserAccount{}, err
	}

	return account, nil
}

func GetAccountHashSalt(email string) (hash string, salt string, err error) {
	// Query the database
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

func InsertAccount(account types.UserAccount, passwordHash string, passwordSalt string) error {
	query := `INSERT INTO user_accounts (
		email,
		role,
		failed_attempts,
		suspended,
		password_hash,
		password_salt
	) VALUES (?, ?, ?, ?, ?, ?)`

	_, err := DB.Exec(query,
		account.Email,
		account.Role,
		account.FailedAttempts,
		account.Suspended,
		passwordHash,
		passwordSalt,
	)
	return err
}

func UpdateAccount(account types.UserAccount) error {
	query := `UPDATE user_accounts SET
		role = ?,
		failed_attempts = ?,
		suspended = ?
		WHERE email = ?`

	_, err := DB.Exec(query,
		account.Role,
		account.FailedAttempts,
		account.Suspended,
		account.Email,
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
