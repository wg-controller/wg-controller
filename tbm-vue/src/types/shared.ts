// Code generated by tygo. DO NOT EDIT.

//////////
// source: types.go

export interface Peer {
  uuid: string;
  hostname: string;
  enabled: boolean;
  privateKey: string; // Wireguard private key (stored encrypted with AES256)
  publicKey: string; // Wireguard public key
  preSharedKey: string; // Wireguard pre-shared key (stored encrypted with AES256)
  keepAliveSeconds: number /* int */; // Wireguard keep-alive interval in seconds
  localTunAddress: string; // The IP address of the client's tunnel interface
  remoteTunAddress: string; // The IP address of the server's tunnel interface (future use)
  remoteSubnets: string[]; // A list of CIDR subnets that the peer can provide access to
  allowedSubnets: string[]; // A list of CIDR subnets that the peer is allowed to access
  lastSeenUnixMillis: number /* int64 */;
  lastIPAddress: string;
  transmitBytes: number /* int64 */;
  receiveBytes: number /* int64 */;
  attributes: string[];
}
export interface PeerInit {
  uuid: string;
  privateKey: string;
  publicKey: string;
  preSharedKey: string;
  localTunAddress: string;
  remoteTunAddress: string;
}
export interface UserAccount {
  email: string;
  role: string; // "user", "admin"
  failedAttempts: number /* int */;
  lastActiveUnixMillis: number /* int64 */;
}
export interface UserAccountWithPass {
  email: string;
  role: string; // "user", "admin"
  password: string;
}
export interface LoginBody {
  email: string;
  password: string;
}
export interface APIKey {
  uuid: string;
  hash: string;
  expiresUnixMillis: number /* int64 */;
  role: string; // Future use
  name: string;
}
export interface ServerInfo {
  publicEndpoint: string;
  nameServers: string[];
}
