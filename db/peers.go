package db

import (
	"errors"
	"log"
	"strings"

	"github.com/lampy255/wg-controller/types"
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
		attributes) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14)`

	_, err = DB.Exec(query,
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
		strings.Join(peer.Attributes, ","))
	return err
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
		attributes=@p13
		WHERE uuid=@p14`

	_, err = DB.Exec(query,
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
		strings.Join(peer.Attributes, ","),
		peer.UUID)
	return err
}

func DeletePeer(uuid string) error {
	// Delete the peer from the database
	_, err := DB.Exec("DELETE FROM peers WHERE uuid = @p1", uuid)
	return err
}
