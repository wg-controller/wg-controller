package types

type Env struct {
	PUBLIC_HOST    string   // Public host for web interface
	ADMIN_EMAIL    string   // Admin email
	ADMIN_PASS     string   // Admin password
	PRIVATE_KEY    string   // Private key for wireguard
	PUBLIC_KEY     string   // Public key for wireguard
	SERVER_CIDR    string   // CIDR Network for tunnel addresses (optional)
	NAME_SERVERS   []string // List of public DNS servers to use (optional)
	INTERFACE_NAME string   // Override kernel interface name (optional)
	WG_PORT        string   // Port for wireguard to listen on (optional)
	UI_PORT        string   // Port for web interface to listen on (optional)
	API_PORT       string   // Port for API to listen on (optional)
}

type Peer struct {
	UUID             string   `json:"uuid"`
	Hostname         string   `json:"hostname"`
	Enabled          bool     `json:"enabled"`
	PeerType         string   `json:"peerType"` // "wg-client", "tbm-client"
	UpdatedMillis    int64    `json:"updatedMillis"`
	PrivateKey       string   `json:"privateKey"`
	PublicKey        string   `json:"publicKey"`
	PreSharedKey     string   `json:"preSharedKey"`
	KeepAliveMillis  int      `json:"keepAliveMillis"`
	LocalTunAddress  string   `json:"localTunAddress"`  // The IP address of the client's tunnel interface
	RemoteTunAddress string   `json:"remoteTunAddress"` // The IP address of the server's tunnel interface (future use)
	RemoteSubnets    []string `json:"remoteSubnets"`    // A list of CIDR subnets that the peer can provide access to
	AllowedSubnets   []string `json:"allowedSubnets"`   // A list of CIDR subnets that the peer is allowed to access
	LastSeenMillis   int64    `json:"lastSeenMillis"`
	LastIPAddress    string   `json:"lastIPAddress"`
}

type PeerExtended struct {
	Peer
	TransmitBytes int64 `json:"transmitBytes"`
	ReceiveBytes  int64 `json:"receiveBytes"`
}

type UserAccount struct {
	Email          string `json:"email"`
	Role           string `json:"role"` // "user", "admin"
	FailedAttempts int    `json:"failedAttempts"`
	Suspended      bool   `json:"suspended"`
}
