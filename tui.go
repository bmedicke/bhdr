package main

import (
	"github.com/bmedicke/bhdr/homeassistant"
	"github.com/bmedicke/bhdr/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TUI structure, all types from tview:
//
// app Application
//   │
// layout Flex
//   ├── switches TreeView
//   │    └── switchesRoot TreeNode
//   │					├── haEntities TreeNode
//   │					│     └── ...
//   │    			└── lorem TreeNode
//   │								└── ...
//   ├── status TextView
//   └── logs TextView

func spawnTUI(haConfig homeassistant.Config) {
	// fill haEntities with nodes:
	haEntities := tview.NewTreeNode("home-assistant")
	entityNames := []string{"edison", "hue", "fan"} // TODO: read from json.
	for _, name := range entityNames {
		entity := tview.NewTreeNode(name)
		entity.SetReference("ref for " + name)
		haEntities.AddChild(entity)
	}

	// create root tree node for the switches view:
	switchesRoot := tview.NewTreeNode(".")
	switchesRoot.SetSelectable(false)

	// attach subnodes:
	switchesRoot.AddChild(haEntities)
	switchesRoot.AddChild(
		tview.NewTreeNode("lorem").
			AddChild(tview.NewTreeNode("ipsum")).
			AddChild(tview.NewTreeNode("dolor").
				AddChild(tview.NewTreeNode("sit"))),
	)

	// create the logs view:
	logs := tview.NewTextView()
	logs.SetTitle("logs").SetBorder(true)

	// create the status view:
	status := tview.NewTextView()
	status.SetBorder(true)

	// create the switches view:
	switches := tview.NewTreeView()
	switches.SetBorder(true).SetTitle("switches")
	switches.SetBorderColor(tcell.ColorGreen)
	switches.SetRoot(switchesRoot)
	switches.SetTopLevel(1) // hide root node.

	// create the layout:
	layout := tview.NewFlex()
	layout.SetBorder(true).SetTitle("B H 🐙 D R")
	layout.AddItem(switches, 0, 1, false)
	layout.AddItem(status, 0, 1, false)
	layout.AddItem(logs, 0, 2, false)

	// create the app:
	app := tview.NewApplication()
	app.SetRoot(layout, true)
	app.SetFocus(switches)

	// switches keybindings:
	switches.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			selection := switches.GetCurrentNode()
			switch event.Rune() {
			case 'J', 'K': // disable tview's default bindings.
				return nil
			case 'H', 'L', 'h', 'l': // use custom vi bindings:
				util.IntuitiveViBindings(event.Rune(), switches)
				return nil // disable defaults.
			case 'i': // print information about current node.
				if parent := util.GetParent(selection, switchesRoot); parent != nil {
					t := "parent: " + parent.GetText() +
						"\ncurrent: " + selection.GetText()
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

	// logs keybindings:
	logs.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case 'd':
				logs.SetText("")
			case 'w':
				// TODO: write log to file.
			}
			return event
		},
	)

	// global keybindings:
	app.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case '[': // focus switches view.
				app.SetFocus(switches)
				switches.SetBorderColor(tcell.ColorGreen)
				logs.SetBorderColor(tcell.ColorWhite)
			case ']': // focus logs view.
				app.SetFocus(logs)
				switches.SetBorderColor(tcell.ColorWhite)
				logs.SetBorderColor(tcell.ColorGreen)
			case 'q': // quit the program.
				app.Stop()
			}
			return event
		},
	)

	// handle pressing Enter on a node:
	switches.SetSelectedFunc(func(node *tview.TreeNode) {})

	// called when focusing a node:
	switches.SetChangedFunc(
		func(node *tview.TreeNode) { status.SetTitle(node.GetText()) },
	)

	// preselect node:
	switches.SetCurrentNode(switchesRoot)

	// listen to Home Assistant events:
	haEvents := make(chan string)
	go homeassistant.GetEvents(haConfig, haEvents)

	// handle Home Assistant events:
	go func() {
		for {
			if logs != nil {
				message := <-haEvents // get message before GetText()!
				logs.SetText(logs.GetText(true) + "\n" + message)
				app.Draw() // required for external changes not based on key presses.
			}
		}
	}()

	app.Run()
}
