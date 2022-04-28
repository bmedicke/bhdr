package main

import (
	"fmt"

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

func spawnTUI(config map[string]interface{}) {
	// channels for communicating with home-assistant:
	haEvents := make(chan string)
	haCommands := make(chan homeassistant.Command)

	// create HA config from global config:
	haConfig := homeassistant.Config{
		Scheme: config["scheme"].(string),
		Server: config["server"].(string),
		Token:  config["token"].(string),
	}

	// create node for home-assistant entities:
	haEntities := tview.NewTreeNode("home-assistant")
	entityMap := config["ha-entities"].(map[string]interface{})

	// fill haEntities with nodes:
	for name, entityID := range entityMap {
		entity := tview.NewTreeNode(name)
		entity.SetReference(entityID)
		haEntities.AddChild(entity)
	}

	// create root tree node for the switches view:
	switchesRoot := tview.NewTreeNode(".")
	switchesRoot.SetSelectable(false)

	// attach subnodes:
	switchesRoot.AddChild(haEntities)

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

	// for keeping track of vi-like key chords:
	chord := util.KeyChord{Active: false, Buffer: "", Action: ""}
	chordmap := config["chordmap"].(map[string]interface{})

	switches.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			selection := switches.GetCurrentNode()
			key := event.Rune()

			if chord.Active {
				if err := util.HandleChords(key, &chord, chordmap); err != nil {
					status.SetText(fmt.Sprint(err))
				}
				if chord.Action != "" {
					status.SetText(chord.Action)
				}
			} else {
				switch key {
				case 'J', 'K': // disable tview's default bindings.
					return nil
				case 'H', 'L', 'h', 'l': // use custom vi bindings:
					util.IntuitiveViBindings(key, switches)
					return nil // disable defaults.
				case 'x', 'c', 'd', 'o', 'y', 'p': // runes that start a chord:
					if err := util.HandleChords(key, &chord, chordmap); err != nil {
						status.SetText(fmt.Sprint(err))
					}
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
				case ';': // toggle entity.
					haCommands <- homeassistant.Command{
						EntityID: fmt.Sprint(selection.GetReference()),
						Service:  "toggle",
					}
				}
			}
			return event
		},
	)

	// logs keybindings:
	logs.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			key := event.Rune()

			switch key {
			case 'd':
				logs.SetText("")
			case 'w':
				util.OverwriteFile("bhdr_log.json", logs.GetText(true))
			}
			return event
		},
	)

	// global keybindings:
	app.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			key := event.Rune()

			switch key {
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

	// connect to Home Assistant:
	go homeassistant.Connect(haConfig, haEvents, haCommands)

	// handle Home Assistant events:
	go func() {
		for {
			if logs != nil {
				message := <-haEvents // get message before GetText()!
				current := logs.GetText(true)
				if len(current) == 0 {
					logs.SetText(message)
				} else {
					logs.SetText(current + ",\n" + message) // append message.
				}
				app.Draw() // required for external changes not based on key presses.
			}
		}
	}()

	app.Run()
}
