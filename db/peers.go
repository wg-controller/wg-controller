package db

import (
	"errors"
	"log"
	"strings"

	"github.com/wg-controller/wg-controller/types"
)

func GetPeers() ([]types.Peer, error) {
	// Query the database
	query := `SELECT
		uuid,
		hostname,
		enabled,
		private_key,
		public_key,
		pre_shared_key,
		keep_alive_seconds,
		local_tun_address,
		remote_tun_address,
		remote_subnets,
		allowed_subnets,
		last_seen_unixmillis,
		last_ip_address,
		os,
		client_version,
		client_type,
		attributes
		FROM peers`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	// Loop through the rows
	var peers []types.Peer
	for rows.Next() {
		var peer types.Peer
		var remoteSubnets string
		var allowedSubnets string
		var attributes string
		err = rows.Scan(
			&peer.UUID,
			&peer.Hostname,
			&peer.Enabled,
			&peer.PrivateKey,
			&peer.PublicKey,
			&peer.PreSharedKey,
			&peer.KeepAliveSeconds,
			&peer.LocalTunAddress,
			&peer.RemoteTunAddress,
			&remoteSubnets,
			&allowedSubnets,
			&peer.LastSeenUnixMillis,
			&peer.LastIPAddress,
			&peer.OS,
			&peer.ClientVersion,
			&peer.ClientType,
			&attributes,
		)
		if err != nil {
			return nil, err
		}

		// Split arrays
		peer.RemoteSubnets = strings.Split(remoteSubnets, ",")
		if len(peer.RemoteSubnets) == 1 {
			if peer.RemoteSubnets[0] == "" {
				peer.RemoteSubnets = []string{}
			}
		}
		peer.AllowedSubnets = strings.Split(allowedSubnets, ",")
		if len(peer.AllowedSubnets) == 1 {
			if peer.AllowedSubnets[0] == "" {
				peer.AllowedSubnets = []string{}
			}
		}
		peer.Attributes = strings.Split(attributes, ",")
		if len(peer.Attributes) == 1 {
			if peer.Attributes[0] == "" {
				peer.Attributes = []string{}
			}
		}

		// Decrypt the private_key
		peer.PrivateKey, err = DecryptAES(peer.PrivateKey, AES_KEY)
		if err != nil {
			return nil, err
		}

		// Decrypt the pre_shared_key
		peer.PreSharedKey, err = DecryptAES(peer.PreSharedKey, AES_KEY)
		if err != nil {
			return nil, err
		}

		peers = append(peers, peer)
	}

	return peers, nil
}

func GetPeer(uuid string) (types.Peer, error) {
	// Query the database
	query := `SELECT
		uuid,
		hostname,
		enabled,
		private_key,
		public_key,
		pre_shared_key,
		keep_alive_seconds,
		local_tun_address,
		remote_tun_address,
		remote_subnets,
		allowed_subnets,
		last_seen_unixmillis,
		last_ip_address,
		os,
		client_version,
		client_type,
		attributes
		FROM peers
		WHERE uuid = @p1`

	row := DB.QueryRow(query, uuid)

	// Scan the row
	var peer types.Peer
	var remoteSubnets string
	var allowedSubnets string
	var attributes string
	err := row.Scan(
		&peer.UUID,
		&peer.Hostname,
		&peer.Enabled,
		&peer.PrivateKey,
		&peer.PublicKey,
		&peer.PreSharedKey,
		&peer.KeepAliveSeconds,
		&peer.LocalTunAddress,
		&peer.RemoteTunAddress,
		&remoteSubnets,
		&allowedSubnets,
		&peer.LastSeenUnixMillis,
		&peer.LastIPAddress,
		&peer.OS,
		&peer.ClientVersion,
		&peer.ClientType,
		&attributes,
	)
	if err != nil {
		return types.Peer{}, err
	}

	// Split arrays
	peer.RemoteSubnets = strings.Split(remoteSubnets, ",")
	if len(peer.RemoteSubnets) == 1 {
		if peer.RemoteSubnets[0] == "" {
			peer.RemoteSubnets = []string{}
		}
	}
	peer.AllowedSubnets = strings.Split(allowedSubnets, ",")
	if len(peer.AllowedSubnets) == 1 {
		if peer.AllowedSubnets[0] == "" {
			peer.AllowedSubnets = []string{}
		}
	}
	peer.Attributes = strings.Split(attributes, ",")
	if len(peer.Attributes) == 1 {
		if peer.Attributes[0] == "" {
			peer.Attributes = []string{}
		}
	}

	// Decrypt the private_key
	peer.PrivateKey, err = DecryptAES(peer.PrivateKey, AES_KEY)
	if err != nil {
		return types.Peer{}, err
	}

	// Decrypt the pre_shared_key
	peer.PreSharedKey, err = DecryptAES(peer.PreSharedKey, AES_KEY)
	if err != nil {
		return types.Peer{}, err
	}

	return peer, nil
}

func InsertPeer(peer types.Peer) (err error) {
	// Encrypt the private_key
	peer.PrivateKey, err = EncryptAES(peer.PrivateKey, AES_KEY)
	if err != nil {
		log.Println(err)
		return errors.New("encryption error")
	}

	// Encrypt the pre_shared_key
	peer.PreSharedKey, err = EncryptAES(peer.PreSharedKey, AES_KEY)
	if err != nil {
		log.Println(err)
		return errors.New("encryption error")
	}

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	// Insert the peer into the database
	query := `INSERT INTO peers (
		uuid,
		hostname,
		enabled,
		private_key,
		public_key,
		pre_shared_key,
		keep_alive_seconds,
		local_tun_address,
		remote_tun_address,
		remote_subnets,
		allowed_subnets,
		last_seen_unixmillis,
		last_ip_address,
		os,
		client_version,
		client_type,
		attributes) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15, @p16, @p17)`

	_, err = tx.Exec(query,
		peer.UUID,
		peer.Hostname,
		peer.Enabled,
		peer.PrivateKey,
		peer.PublicKey,
		peer.PreSharedKey,
		peer.KeepAliveSeconds,
		peer.LocalTunAddress,
		peer.RemoteTunAddress,
		strings.Join(peer.RemoteSubnets, ","),
		strings.Join(peer.AllowedSubnets, ","),
		peer.LastSeenUnixMillis,
		peer.LastIPAddress,
		peer.OS,
		peer.ClientVersion,
		peer.ClientType,
		strings.Join(peer.Attributes, ","))
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func UpdatePeer(peer types.Peer) (err error) {
	// Encrypt the private_key
	peer.PrivateKey, err = EncryptAES(peer.PrivateKey, AES_KEY)
	if err != nil {
		log.Println(err)
		return errors.New("encryption error")
	}

	// Encrypt the pre_shared_key
	peer.PreSharedKey, err = EncryptAES(peer.PreSharedKey, AES_KEY)
	if err != nil {
		log.Println(err)
		return errors.New("encryption error")
	}

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	// Update the peer in the database
	query := `UPDATE peers SET
		hostname=@p1,
		enabled=@p2,
		private_key=@p3,
		public_key=@p4,
		pre_shared_key=@p5,
		keep_alive_seconds=@p6,
		local_tun_address=@p7,
		remote_tun_address=@p8,
		remote_subnets=@p9,
		allowed_subnets=@p10,
		last_seen_unixmillis=@p11,
		last_ip_address=@p12,
		os=@p13,
		client_version=@p14,
		client_type=@p15,
		attributes=@p16
		WHERE uuid=@p17`

	_, err = tx.Exec(query,
		peer.Hostname,
		peer.Enabled,
		peer.PrivateKey,
		peer.PublicKey,
		peer.PreSharedKey,
		peer.KeepAliveSeconds,
		peer.LocalTunAddress,
		peer.RemoteTunAddress,
		strings.Join(peer.RemoteSubnets, ","),
		strings.Join(peer.AllowedSubnets, ","),
		peer.LastSeenUnixMillis,
		peer.LastIPAddress,
		peer.OS,
		peer.ClientVersion,
		peer.ClientType,
		strings.Join(peer.Attributes, ","),
		peer.UUID)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func DeletePeer(uuid string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM peers WHERE uuid = ?", uuid)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
