package domain

import "time"

type NodeStatus string

const (
	NodeStatusActive     NodeStatus = "ACTIVE"
	NodeStatusInactive   NodeStatus = "INACTIVE"
	NodeStatusUnverified NodeStatus = "UNVERIFIED"
	NodeStatusSuspended  NodeStatus = "SUSPENDED"
)

// Node representa un nodo de procesamiento criptográfico validado dentro de la red FCOS.
type Node struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Status        NodeStatus    `json:"status"`
	Identity      *NodeIdentity `json:"identity,omitempty"`
	HardwareModel string        `json:"hardware_model"`
	OSVersion     string        `json:"os_version"`
	LastSeenAtUTC *time.Time    `json:"last_seen_at_utc"`
	CreatedAtUTC  time.Time     `json:"created_at_utc"`
	UpdatedAtUTC  time.Time     `json:"updated_at_utc"`
}
