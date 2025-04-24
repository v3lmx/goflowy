package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rivo/tview"
	_ "modernc.org/sqlite"

	"github.com/v3lmx/goflowy/core"
	"github.com/v3lmx/goflowy/storage"
)

func main() {
	db, err := sqlx.Connect("sqlite", "db.sqlite")
	if err != nil {
		panic("Could not connect to database")
	}

	storage := storage.NewSQLxStorage(db)

	origins, err := core.GetOrigins(storage)
	if err != nil {
		panic("Could not get origins: " + err.Error())
	}

	root := tview.NewTreeNode("root").
		SetColor(tcell.ColorRed)

	origin := origins[0]
	children, err := core.GetNodes(storage, origin.ID)
	if err != nil {
		panic("Could not get children" + err.Error())
	}

	for _, child := range children {
		treeNode := tview.NewTreeNode(child.Contents)
		treeNode.SetReference(child.ID)
		root.AddChild(treeNode)
	}

	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		id := node.GetReference()
		fmt.Println(id)
	})

	// frame := tview.NewFrame(tree).
	// 	SetBorder(true).
	// 	SetTitle(" Goflowy ")
	//
	// if err := tview.NewApplication().SetRoot(frame, true).SetFocus(tree).Run(); err != nil {
	if err := tview.NewApplication().SetRoot(tree, true).Run(); err != nil {
		panic(err)
	}
}
