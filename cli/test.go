package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/v3lmx/goflowy/core"
	"github.com/v3lmx/goflowy/storage"
	_ "modernc.org/sqlite"
)

func Init() {
	fmt.Printf("Goflowy\n")

	fmt.Printf("Connecting to database...")
	db, err := sqlx.Connect("sqlite", "db.sqlite")
	if err != nil {
		panic("Could not connect to database")
	}

	storage := storage.NewSQLxStorage(db)

	origin, err := core.NewOriginNode(storage)
	if err != nil {
		panic("Could not create origin: " + err.Error())
	}

	_, err = core.NewNode(storage, "first node", origin.ID)
	if err != nil {
		panic("Could not create node: " + err.Error())
	}
	_, err = core.NewNode(storage, "second node", origin.ID)
	if err != nil {
		panic("Could not create node: " + err.Error())
	}
	_, err = core.NewNode(storage, "third node", origin.ID)
	if err != nil {
		panic("Could not create node: " + err.Error())
	}

	nodes, err := core.GetNodes(storage, origin.ID)
	if err != nil {
		panic("Could not list node: " + err.Error())
	}

	fmt.Printf("\nnodes: %v", nodes)
}

func main() {
	Init()
}
