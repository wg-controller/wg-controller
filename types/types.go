package types

type Peer struct {
	UUID               string   `json:"uuid"`
	Hostname           string   `json:"hostname"`
	Enabled            bool     `json:"enabled"`
	PrivateKey         string   `json:"privateKey"`       // Wireguard private key (stored encrypted with AES256)
	PublicKey          string   `json:"publicKey"`        // Wireguard public key
	PreSharedKey       string   `json:"preSharedKey"`     // Wireguard pre-shared key (stored encrypted with AES256)
	KeepAliveSeconds   int      `json:"keepAliveSeconds"` // Wireguard keep-alive interval in seconds
	LocalTunAddress    string   `json:"localTunAddress"`  // The IP address of the server's tunnel interface (future use)
	RemoteTunAddress   string   `json:"remoteTunAddress"` // The IP address of the peer's tunnel interface
	RemoteSubnets      []string `json:"remoteSubnets"`    // A list of CIDR subnets that the peer can provide access to
	AllowedSubnets     []string `json:"allowedSubnets"`   // A list of CIDR subnets that the peer is allowed to access
	LastSeenUnixMillis int64    `json:"lastSeenUnixMillis"`
	LastIPAddress      string   `json:"lastIPAddress"`
	TransmitBytes      int64    `json:"transmitBytes"`
	ReceiveBytes       int64    `json:"receiveBytes"`
	Attributes         []string `json:"attributes"`
}

type PeerInit struct {
	UUID             string `json:"uuid"`
	PrivateKey       string `json:"privateKey"`
	PublicKey        string `json:"publicKey"`
	PreSharedKey     string `json:"preSharedKey"`
	LocalTunAddress  string `json:"localTunAddress"`
	RemoteTunAddress string `json:"remoteTunAddress"`
	ServerCIDR       string `json:"serverCIDR"`
}

type UserAccount struct {
	Email                string `json:"email"`
	Role                 string `json:"role"` // "user", "admin"
	FailedAttempts       int    `json:"failedAttempts"`
	LastActiveUnixMillis int64  `json:"lastActiveUnixMillis"`
}

type UserAccountWithPass struct {
	Email          string `json:"email"`
	Role           string `json:"role"` // "user", "admin"
	FailedAttempts int    `json:"failedAttempts"`
	Password       string `json:"password"`
}

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type APIKey struct {
	UUID              string   `json:"uuid"`
	Name              string   `json:"name"`
	ExpiresUnixMillis int64    `json:"expiresUnixMillis"`
	Attributes        []string `json:"attributes"`
}

type APIKeyWithToken struct {
	UUID              string   `json:"uuid"`
	Name              string   `json:"name"`
	ExpiresUnixMillis int64    `json:"expiresUnixMillis"`
	Attributes        []string `json:"attributes"`
	Token             string   `json:"token"`
}

type APIKeyInit struct {
	UUID  string `json:"uuid"`
	Token string `json:"token"`
}

type ServerInfo struct {
	PublicKey        string   `json:"publicKey"`
	PublicEndpoint   string   `json:"publicEndpoint"`
	NameServers      []string `json:"nameServers"`
	Netmask          string   `json:"netmask"`
	ServerInternalIP string   `json:"serverInternalIP"`
}

type Password struct {
	Password string `json:"password"`
}
