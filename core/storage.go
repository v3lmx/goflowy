package core

import (
	"github.com/rs/xid"
)

type Storage interface {
	GetNodes(parentID xid.ID) ([]Node, error)
	GetOrigins() ([]Node, error)
	AddNode(parentID xid.ID, node Node) error
	DeleteNode(id xid.ID) error
	Exists(nodeID xid.ID) (bool, error)
}
