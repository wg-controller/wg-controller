package types

type Env struct {
	PUBLIC_HOST    string   // Public host for web interface
	ADMIN_EMAIL    string   // Admin email
	ADMIN_PASS     string   // Admin password
	WG_PRIVATE_KEY string   // Private key for wireguard
	WG_PUBLIC_KEY  string   // Public key for wireguard
	DB_AES_KEY     string   // Base64 encoded 32 Byte AES key for encrypting private keys
	SERVER_CIDR    string   // CIDR Network for tunnel addresses (optional)
	NAME_SERVERS   []string // List of public DNS servers to use (optional)
	INTERFACE_NAME string   // Override kernel interface name (optional)
	WG_PORT        string   // Port for wireguard to listen on (optional)
	API_PORT       string   // Port for API to listen on (optional)
}

type Peer struct {
	UUID               string   `json:"uuid"`
	Hostname           string   `json:"hostname"`
	Enabled            bool     `json:"enabled"`
	PeerType           string   `json:"peerType"` // "wg-client", "tbm-client"
	UpdatedUnixMillis  int64    `json:"updatedUnixMillis"`
	PrivateKey         string   `json:"privateKey"`       // Wireguard private key (stored encrypted with AES256)
	PublicKey          string   `json:"publicKey"`        // Wireguard public key
	PreSharedKey       string   `json:"preSharedKey"`     // Wireguard pre-shared key (stored encrypted with AES256)
	KeepAliveMillis    int      `json:"keepAliveMillis"`  // Wireguard keep-alive interval in milliseconds
	LocalTunAddress    string   `json:"localTunAddress"`  // The IP address of the client's tunnel interface
	RemoteTunAddress   string   `json:"remoteTunAddress"` // The IP address of the server's tunnel interface (future use)
	RemoteSubnets      []string `json:"remoteSubnets"`    // A list of CIDR subnets that the peer can provide access to
	AllowedSubnets     []string `json:"allowedSubnets"`   // A list of CIDR subnets that the peer is allowed to access
	LastSeenUnixMillis int64    `json:"lastSeenUnixMillis"`
	LastIPAddress      string   `json:"lastIPAddress"`
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
}

type UserAccountWithPass struct {
	UserAccount
	Password string `json:"password"`
}

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type APIKey struct {
	UUID              string `json:"uuid"`
	Hash              string `json:"hash"`
	ExpiresUnixMillis int64  `json:"expiresUnixMillis"`
	Role              string `json:"role"` // Future use
	Name              string `json:"name"`
}
