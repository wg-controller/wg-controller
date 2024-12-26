package db

import "github.com/lampy255/net-tbm/types"

func GetPeers() ([]types.Peer, error) {
	// Query the database
	query := `SELECT
		uuid,
		hostname,
		enabled,
		peer_type,
		updated_millis,
		private_key,
		public_key,
		pre_shared_key,
		keep_alive_millis,
		local_tun_address,
		remote_tun_address,
		remote_subnets,
		allowed_subnets,
		last_seen_millis,
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
			&peer.UpdatedMillis,
			&peer.PrivateKey,
			&peer.PublicKey,
			&peer.PreSharedKey,
			&peer.KeepAliveMillis,
			&peer.LocalTunAddress,
			&peer.RemoteTunAddress,
			&peer.RemoteSubnets,
			&peer.AllowedSubnets,
			&peer.LastSeenMillis,
			&peer.LastIPAddress,
		)
		if err != nil {
			return nil, err
		}
		peers = append(peers, peer)
	}

	return peers, nil
}

func InsertPeer(peer types.Peer) error {
	// Insert the peer into the database
	query := `INSERT INTO peers (
		uuid,
		hostname,
		enabled,
		peer_type,
		updated_millis,
		private_key,
		public_key,
		pre_shared_key,
		keep_alive_millis,
		local_tun_address,
		remote_tun_address,
		remote_subnets,
		allowed_subnets,
		last_seen_millis,
		last_ip_address) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15)`

	_, err := DB.Exec(query,
		peer.UUID,
		peer.Hostname,
		peer.Enabled,
		peer.PeerType,
		peer.UpdatedMillis,
		peer.PrivateKey,
		peer.PublicKey,
		peer.PreSharedKey,
		peer.KeepAliveMillis,
		peer.LocalTunAddress,
		peer.RemoteTunAddress,
		peer.RemoteSubnets,
		peer.AllowedSubnets,
		peer.LastSeenMillis,
		peer.LastIPAddress)
	return err
}

func UpdatePeer(peer types.Peer) error {
	// Update the peer in the database
	query := `UPDATE peers SET
		hostname=@p2,
		enabled=@p3,
		peer_type=@p4,
		updated_millis=@p5,
		private_key=@p6,
		public_key=@p7,
		pre_shared_key=@p8,
		keep_alive_millis=@p9,
		local_tun_address=@p10,
		remote_tun_address=@p11,
		remote_subnets=@p12,
		allowed_subnets=@p13,
		last_seen_millis=@p14,
		last_ip_address=@p15
		WHERE uuid=@p1`

	_, err := DB.Exec(query,
		peer.UUID,
		peer.Hostname,
		peer.Enabled,
		peer.PeerType,
		peer.UpdatedMillis,
		peer.PrivateKey,
		peer.PublicKey,
		peer.PreSharedKey,
		peer.KeepAliveMillis,
		peer.LocalTunAddress,
		peer.RemoteTunAddress,
		peer.RemoteSubnets,
		peer.AllowedSubnets,
		peer.LastSeenMillis,
		peer.LastIPAddress)
	return err
}

func DeletePeer(uuid string) error {
	// Delete the peer from the database
	_, err := DB.Exec("DELETE FROM peers WHERE uuid = @p1", uuid)
	return err
}
