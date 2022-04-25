package main

import (
	"github.com/bmedicke/bhdr/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	rootNode := tview.NewTreeNode(".")
	// rootNode.SetSelectable(false)

	haEntities := tview.NewTreeNode("home-assistant")

	// fill haEntities with nodes:
	entityNames := []string{"edison", "hue", "fan"} // TODO: read from json.
	for _, name := range entityNames {
		entity := tview.NewTreeNode(name)
		entity.SetReference("ref for " + name)
		haEntities.AddChild(entity)
	}

	rootNode.AddChild(haEntities)
	rootNode.AddChild(
		tview.NewTreeNode("localFunctions").
			AddChild(tview.NewTreeNode("func0")).
			AddChild(tview.NewTreeNode("func1").
				AddChild(tview.NewTreeNode("sub0"))),
	)

	// create the status view:
	status := tview.NewTextView()
	status.SetBorder(true)

	// create the switches view:
	switches := tview.NewTreeView()
	switches.SetBorder(true).SetTitle("switches")
	switches.SetRoot(rootNode)
	switches.SetTopLevel(1)

	// create the layout:
	layout := tview.NewFlex()
	layout.SetBorder(true).SetTitle("B H 🐙 D R")
	layout.AddItem(switches, 0, 1, false)
	layout.AddItem(status, 0, 1, false)

	// create the app:
	app := tview.NewApplication()
	app.SetRoot(layout, true)
	app.SetFocus(switches)

	// TODO: extract custom vim functionality.
	switches.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			sel := switches.GetCurrentNode()
			switch event.Rune() {
			case 'J', 'K': // remove default bindings.
				return nil
			case 'H': // collapse everything but the root node.
				// calling .CollapseAll() on the children of rootNode
				// does not work for some reason. do it manually:
				rootNode.Walk(func(node, parent *tview.TreeNode) bool {
					if parent != nil {
						node.Collapse()
					}
					return true // visit all nodes.
				})
			case 'L': // expand everything.
				rootNode.ExpandAll()
			case 'q': // quit the program.
				app.Stop()
			case '[': // TODO: jump to previous node on same level.
				break
			case ']': // TODO: jump to next node on same level.
				break
			case 'h': // move up OR collapse node.
				if sel.IsExpanded() && nil != sel.GetChildren() {
					sel.Collapse()
				} else if sel.GetLevel() > 1 {
					parent := util.GetParent(sel, rootNode)
					switches.SetCurrentNode(parent)
				}
			case 'l': // expand node.
				if !sel.IsExpanded() {
					sel.Expand()
				}
			case ';', '\'': // TODO: toggle entities, etc...
				status.SetText(string(event.Rune()) + " on " + sel.GetText())
			case 'i': // print information about current node.
				if parent := util.GetParent(sel, rootNode); parent != nil {
					t := "parent: " + parent.GetText() +
						"\ncurrent: " + sel.GetText()
					status.SetText(t)
				} else {
					t := "no parent found" +
						"\ncurrent: " + sel.GetText()
					status.SetText(t)
				}
			}
			return event
		},
	)

	// handle pressing Enter on a node:
	switches.SetSelectedFunc(func(node *tview.TreeNode) {})

	// handle focusing a node:
	switches.SetChangedFunc(func(node *tview.TreeNode) {})

	switches.SetCurrentNode(rootNode)
	app.Run()
}
