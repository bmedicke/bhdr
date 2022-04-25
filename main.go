package main

import (
	"fmt"

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

	// keybindings:
	switches.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			selection := switches.GetCurrentNode()
			switch event.Rune() {
			case 'J', 'K': // disable tview's default bindings.
				return nil
			case 'H', 'L', 'h', 'l': // use custom vi bindings:
				return util.IntuitiveViBindings(event.Rune(), switches)
			case 'q': // quit the program.
				app.Stop()
			case 'i': // print information about current node.
				if parent := util.GetParent(selection, rootNode); parent != nil {
					t := "parent: " + parent.GetText() +
						"\ncurrent: " + selection.GetText() +
						fmt.Sprintf("%T", event.Rune)
					status.SetText(t)
				} else {
					t := "no parent found" +
						"\ncurrent: " + selection.GetText()
					status.SetText(t)
				}
			case ';', '\'': // TODO: toggle entities, etc...
				status.SetText(
					string(event.Rune()) + " on " + selection.GetText(),
				)
			}
			return event
		},
	)

	// handle pressing Enter on a node:
	switches.SetSelectedFunc(func(node *tview.TreeNode) {})

	// handle focusing a node:
	switches.SetChangedFunc(func(node *tview.TreeNode) {})

	// preselect node and start app:
	switches.SetCurrentNode(rootNode)
	app.Run()
}
