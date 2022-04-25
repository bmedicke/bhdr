package main

import (
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
			AddChild(tview.NewTreeNode("func1")),
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

	// TODO: extract.
	switches.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			sel := switches.GetCurrentNode()
			switch event.Rune() {
			default:
				break
			case 'q':
				app.Stop()
			case '[':
				break // TODO: jump to previous node on same level.
			case ']':
				break // TODO: jump to next node on same level.
			case 'h':
				sel.Collapse()
			case 'l':
				sel.Expand()
			case ';', '\'':
				status.SetText(string(event.Rune()) + " on " + sel.GetText())
			case 'p':
				if parent := GetParent(sel, rootNode); parent != nil {
					status.SetText("parent: " + parent.GetText())
				} else {
					status.SetText("no parent found")
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

// GetParent returns parent or nil (if it was not found).
// https://github.com/rivo/tview/issues/246#issuecomment-471173854
// TODO: extract.
func GetParent(
	node, root *tview.TreeNode,
) *tview.TreeNode {
	var match *tview.TreeNode

	// walk the tree to find our node (and thus its parent):
	root.Walk(
		func(current, parent *tview.TreeNode) bool {
			// current node found:
			if current == node {
				if parent != nil {
					match = parent
				}
				// stop walk:
				return false
			}
			// current node not found, continue walk:
			return true
		},
	)
	return match
}
