package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"

	"github.com/v3lmx/goflowy/core"
)

type Node struct {
	ID         string        `db:"id"`
	Sequence   int           `db:"sequence"`
	Type       string        `db:"node_type"`
	Contents   string        `db:"contents"`
	Children   string        `db:"children"`
	Parent     string        `db:"parent"`
	Deleted    int           `db:"deleted"`
	CreatedAt  int64         `db:"created_at"`
	ModifiedAt int64         `db:"modified_at"`
	DeletedAt  sql.NullInt64 `db:"deleted_at"`
}

func (n Node) Parse() (core.Node, error) {
	id, err := xid.FromString(n.ID)
	if err != nil {
		return core.Node{}, fmt.Errorf("Could not parse id for node with id='%s' : %v ;", n.ID, err)
	}

	nodeType, err := core.ToNodeType(n.Type)
	if err != nil {
		return core.Node{}, fmt.Errorf("Could not parse node_type for node with id='%s' : %v ;", n.ID, err)
	}

	children := strings.Split(n.Children, ",")
	childrenIDs := make([]xid.ID, len(children))

	for i, child := range children {
		id, err := xid.FromString(n.ID)
		if err != nil {
			return core.Node{}, fmt.Errorf("Could not parse child id = '%s' for node with id='%s' : %v ;", child, n.ID, err)
		}

		childrenIDs[i] = id
	}

	var parent xid.ID
	if nodeType != core.TypeOrigin {
		parent, err = xid.FromString(n.Parent)
		if err != nil {
			return core.Node{}, fmt.Errorf("Could not parse parent for node with id='%s' : %v ;", n.ID, err)
		}
	}

	var deleted bool
	if n.Deleted == 1 {
		deleted = true
	}

	createdAt := time.Unix(n.CreatedAt, 0)
	modifiedAt := time.Unix(n.ModifiedAt, 0)
	var deletedAt time.Time

	if deleted && n.DeletedAt.Valid {
		deletedAt = time.Unix(n.DeletedAt.Int64, 0)
	}

	node := core.Node{
		ID:       id,
		Sequence: n.Sequence,
		Type:     nodeType,
		Contents: n.Contents,
		Children: childrenIDs,
		Parent:   parent,
		Deleted:  deleted,
		Metadata: core.Metadata{
			CreatedAt:  createdAt,
			ModifiedAt: modifiedAt,
			DeletedAt:  deletedAt,
		},
	}

	return node, nil
}

func NewSQLxStorage(db *sqlx.DB) SQLxStorage {
	return SQLxStorage{
		DB: db,
	}
}

type SQLxInstance interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type SQLxStorage struct {
	DB SQLxInstance
}

func (s SQLxStorage) GetNodes(parentID xid.ID) ([]core.Node, error) {
	query := `select id, sequence, node_type, contents, children, parent, created_at, modified_at, deleted_at from nodes where nodes.parent = ? and nodes.deleted = false`

	facet := make([]Node, 0)
	err := s.DB.Select(&facet, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch nodes: %v", err)
	}

	nodes := make([]core.Node, len(facet))

	for i, row := range facet {
		node, err := row.Parse()
		if err != nil {
			return nodes, fmt.Errorf("Could not parse node : %v", err)
		}

		nodes[i] = node
	}

	return nodes, nil
}

func (s SQLxStorage) GetOrigins() ([]core.Node, error) {
	query := `select id, sequence, node_type, contents, children, created_at, modified_at, deleted_at from nodes where nodes.node_type = 'origin' and nodes.deleted = false`

	facet := make([]Node, 0)
	err := s.DB.Select(&facet, query)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch nodes: %v", err)
	}

	nodes := make([]core.Node, len(facet))

	for i, row := range facet {
		node, err := row.Parse()
		if err != nil {
			return nodes, fmt.Errorf("Could not parse node : %v", err)
		}

		nodes[i] = node
	}

	return nodes, nil
}

func (s SQLxStorage) AddNode(parentID xid.ID, node core.Node) error {
	query := `insert into nodes(id, sequence, node_type, contents, children, parent, created_at, modified_at)
	values (?, ?, ?, ?, ?, ?, ?, ?)`

	children := fmt.Sprintf("%v", node.Children)
	createdAt := node.Metadata.CreatedAt.Unix()
	modifiedAt := node.Metadata.ModifiedAt.Unix()

	_, err := s.DB.Exec(query, node.ID.String(), node.Sequence, node.Type, node.Contents, children, node.Parent, createdAt, modifiedAt)
	if err != nil {
		return fmt.Errorf("Could not add node: %v", err)
	}

	return nil
}

func (s SQLxStorage) DeleteNode(id xid.ID) error {
	query := `update nodes set deleted = true, deleted_at = ? where id = ?`

	now := time.Now().Unix()
	_, err := s.DB.Exec(query, now, id)
	if err != nil {
		return fmt.Errorf("Could not delete node: %v", err)
	}

	return nil
}

func (s SQLxStorage) Exists(id xid.ID) (bool, error) {
	query := `select 1 from nodes where id = ?`

	var exists int
	err := s.DB.Get(&exists, query, id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("Could not check for existence: %v", err)
	}
	return true, nil
}
