package core

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"
)

var OriginID = xid.NilID()

type NodeType string

const (
	TypeText   NodeType = "text"
	TypeOrigin NodeType = "origin"
)

func ToNodeType(s string) (NodeType, error) {
	switch s {
	case string(TypeOrigin):
		return TypeOrigin, nil
	case string(TypeText):
		return TypeText, nil
	default:
		return TypeOrigin, errors.New("Unknown node type")
	}
}

type Node struct {
	ID       xid.ID
	Sequence int
	Type     NodeType
	Contents string
	Children []xid.ID
	Parent   xid.ID
	Deleted  bool
	Metadata Metadata
}

type Metadata struct {
	CreatedAt  time.Time
	ModifiedAt time.Time
	DeletedAt  time.Time
}

func NewOriginNode(s Storage) (Node, error) {
	metadata := Metadata{
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	node := Node{
		ID:       xid.New(),
		Sequence: 0,
		Type:     TypeOrigin,
		Children: make([]xid.ID, 0),
		Parent:   OriginID,
		Metadata: metadata,
	}

	err := s.AddNode(OriginID, node)
	if err != nil {
		return Node{}, fmt.Errorf("Could not save node: %v", err)
	}

	return node, nil
}

func NewNode(s Storage, contents string, parent xid.ID) (Node, error) {
	exists, err := s.Exists(parent)
	if err != nil {
		return Node{}, fmt.Errorf("Could not check existence of the parent: %v ;", err)
	}
	if !exists {
		return Node{}, errors.New("The parent does not exist ;")
	}

	metadata := Metadata{
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	node := Node{
		ID:       xid.New(),
		Sequence: 0,
		Type:     TypeText,
		Contents: contents,
		Children: make([]xid.ID, 0),
		Parent:   parent,
		Metadata: metadata,
	}

	err = s.AddNode(parent, node)
	if err != nil {
		return Node{}, fmt.Errorf("Could not save node: %v ;", err)
	}

	return node, nil
}

func GetNodes(s Storage, parent xid.ID) ([]Node, error) {
	nodes := make([]Node, 0)
	exists, err := s.Exists(parent)
	if err != nil {
		return nodes, fmt.Errorf("Could not check existence of the parent: %v ;", err)
	}
	if !exists {
		return nodes, errors.New("The parent does not exist ;")
	}

	nodes, err = s.GetNodes(parent)
	if err != nil {
		return nodes, fmt.Errorf("Could not get nodes: %v ;", err)
	}

	return nodes, nil
}

func GetOrigins(s Storage) ([]Node, error) {
	nodes := make([]Node, 0)

	nodes, err := s.GetOrigins()
	if err != nil {
		return nodes, fmt.Errorf("Could not get origin nodes: %v ;", err)
	}

	return nodes, nil
}
