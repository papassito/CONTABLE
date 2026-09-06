package domain

// NodeIdentity modela la firma criptográfica asimétrica del nodo
type NodeIdentity struct {
	NodeID       string `json:"node_id"`
	PublicKeyPEM string `json:"public_key_pem"`
	KeyVersion   string `json:"key_version"`
}
