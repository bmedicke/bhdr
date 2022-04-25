package main

import (
	"github.com/rivo/tview"
)

func main() {
	entityNames := []string{"edison", "hue", "fan"}
	entities := tview.NewTreeNode("entities")
	entities.SetReference("json")
	// entities.SetSelectable(false)

	for _, name := range entityNames {
		entity := tview.NewTreeNode(name)
		entity.SetReference("json")
		entities.AddChild(entity)
	}

	status := tview.NewTextView()
	status.SetBorder(true).SetTitle("status")

	ha := tview.NewTreeView()
	ha.SetBorder(true).SetTitle("home")
	ha.SetRoot(entities)
	ha.SetCurrentNode(entities)

	ha.SetSelectedFunc(
		func(node *tview.TreeNode) {
		},
	)

	ha.SetChangedFunc(
		func(node *tview.TreeNode) {
			text := node.GetText()
			ref := node.GetReference().(string)
			status.SetTitle(text)
			status.SetText(text + " " + ref)
		},
	)

	root := tview.NewFlex()
	root.SetBorder(true).SetTitle("B H 🐙 D R")
	root.AddItem(ha, 0, 1, false)
	root.AddItem(status, 0, 1, false)

	app := tview.NewApplication()
	app.SetRoot(root, true)
	app.SetFocus(ha)
	app.Run()
}
