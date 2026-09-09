package domain

import "errors"

// NodeIdentity modela la firma criptográfica asimétrica del nodo
type NodeIdentity struct {
	NodeID       string `json:"node_id"`
	PublicKeyPEM string `json:"public_key_pem"`
	KeyVersion   string `json:"key_version"`
}

// Validate verifica que la identidad contenga llaves y metadatos consistentes.
func (ni *NodeIdentity) Validate() error {
	if ni.NodeID == "" {
		return errors.New("el ID del nodo es requerido")
	}
	if ni.PublicKeyPEM == "" {
		return errors.New("la llave pública PEM no puede estar vacía")
	}
	return nil
}
