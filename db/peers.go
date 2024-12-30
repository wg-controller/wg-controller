package db

import (
	"errors"
	"log"

	"github.com/lampy255/net-tbm/types"
)

func GetPeers() ([]types.Peer, error) {
	// Query the database
	query := `SELECT
		uuid,
		hostname,
		enabled,
		peer_type,
		updated_unixmillis,
		private_key,
		public_key,
		pre_shared_key,
		keep_alive_millis,
		local_tun_address,
		remote_tun_address,
		remote_subnets,
		allowed_subnets,
		last_seen_unixmillis,
		last_ip_address
		FROM peers`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	// Loop through the rows
	var peers []types.Peer
	for rows.Next() {
		var peer types.Peer
		err = rows.Scan(
			&peer.UUID,
			&peer.Hostname,
			&peer.Enabled,
			&peer.PeerType,
			&peer.UpdatedUnixMillis,
			&peer.PrivateKey,
			&peer.PublicKey,
			&peer.PreSharedKey,
			&peer.KeepAliveMillis,
			&peer.LocalTunAddress,
			&peer.RemoteTunAddress,
			&peer.RemoteSubnets,
			&peer.AllowedSubnets,
			&peer.LastSeenUnixMillis,
			&peer.LastIPAddress,
		)
		if err != nil {
			return nil, err
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
		peer_type,
		updated_unixmillis,
		private_key,
		public_key,
		pre_shared_key,
		keep_alive_millis,
		local_tun_address,
		remote_tun_address,
		remote_subnets,
		allowed_subnets,
		last_seen_unixmillis,
		last_ip_address
		FROM peers
		WHERE uuid = @p1`
	row := DB.QueryRow(query, uuid)

	// Scan the row
	var peer types.Peer
	err := row.Scan(
		&peer.UUID,
		&peer.Hostname,
		&peer.Enabled,
		&peer.PeerType,
		&peer.UpdatedUnixMillis,
		&peer.PrivateKey,
		&peer.PublicKey,
		&peer.PreSharedKey,
		&peer.KeepAliveMillis,
		&peer.LocalTunAddress,
		&peer.RemoteTunAddress,
		&peer.RemoteSubnets,
		&peer.AllowedSubnets,
		&peer.LastSeenUnixMillis,
		&peer.LastIPAddress,
	)
	if err != nil {
		return types.Peer{}, err
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
		peer_type,
		updated_unixmillis,
		private_key,
		public_key,
		pre_shared_key,
		keep_alive_millis,
		local_tun_address,
		remote_tun_address,
		remote_subnets,
		allowed_subnets,
		last_seen_unixmillis,
		last_ip_address) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15)`

	_, err = DB.Exec(query,
		peer.UUID,
		peer.Hostname,
		peer.Enabled,
		peer.PeerType,
		peer.UpdatedUnixMillis,
		peer.PrivateKey,
		peer.PublicKey,
		peer.PreSharedKey,
		peer.KeepAliveMillis,
		peer.LocalTunAddress,
		peer.RemoteTunAddress,
		peer.RemoteSubnets,
		peer.AllowedSubnets,
		peer.LastSeenUnixMillis,
		peer.LastIPAddress)
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
		hostname=@p2,
		enabled=@p3,
		peer_type=@p4,
		updated_unixmillis=@p5,
		private_key=@p6,
		public_key=@p7,
		pre_shared_key=@p8,
		keep_alive_millis=@p9,
		local_tun_address=@p10,
		remote_tun_address=@p11,
		remote_subnets=@p12,
		allowed_subnets=@p13,
		last_seen_unixmillis=@p14,
		last_ip_address=@p15
		WHERE uuid=@p1`

	_, err = DB.Exec(query,
		peer.UUID,
		peer.Hostname,
		peer.Enabled,
		peer.PeerType,
		peer.UpdatedUnixMillis,
		peer.PrivateKey,
		peer.PublicKey,
		peer.PreSharedKey,
		peer.KeepAliveMillis,
		peer.LocalTunAddress,
		peer.RemoteTunAddress,
		peer.RemoteSubnets,
		peer.AllowedSubnets,
		peer.LastSeenUnixMillis,
		peer.LastIPAddress)
	return err
}

func DeletePeer(uuid string) error {
	// Delete the peer from the database
	_, err := DB.Exec("DELETE FROM peers WHERE uuid = @p1", uuid)
	return err
}
