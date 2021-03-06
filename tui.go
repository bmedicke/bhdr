package main

import (
	"encoding/json"
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
// frame Frame
//   │
// outerLayout Flex (FlexRow)
//   │
//   ├── innerLayout Flex (FlexColumn)
//   │     ├── switches TreeView
//   │     │    └── switchesRoot TreeNode
//   │     │					└── haEntities TreeNode
//   │     │					      └── ...
//   │     └── status TextView
//   ├── statusbar TextView
//   └── logs TextView

func spawnTUI(config map[string]interface{}, showLogs bool) {
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
	haEntities.SetReference(homeassistant.Data{})
	entitySlice := config["ha-entities"].([]interface{})

	// fill entities with nodes:
	for _, entityJSON := range entitySlice {
		entityMap := entityJSON.(map[string]interface{})
		entity := tview.NewTreeNode(entityMap["id"].(string))
		entity.SetReference(
			homeassistant.Data{
				EntityID: entityMap["entity-id"].(string),
				NickName: entityMap["id"].(string),
			},
		)
		haEntities.AddChild(entity)
	}

	// create root tree node for the switches view:
	switchesRoot := tview.NewTreeNode(".")
	switchesRoot.SetSelectable(false)

	// attach subnodes:
	switchesRoot.AddChild(haEntities)

	// create statusbar view:
	statusbar := tview.NewTextView()
	statusbar.SetBackgroundColor(tcell.ColorDarkOliveGreen)

	// create the status view:
	status := tview.NewTextView()
	status.SetBorder(true).SetTitle("status")

	// create the switches view:
	switches := tview.NewTreeView()
	switches.SetBorder(true).SetTitle("switches")
	switches.SetBorderColor(tcell.ColorGreen)
	switches.SetRoot(switchesRoot)
	switches.SetTopLevel(1) // hide root node.

	// create the layout:
	innerLayout := tview.NewFlex()
	innerLayout.AddItem(switches, 0, 1, false)
	innerLayout.AddItem(status, 0, 1, false)
	outerLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	outerLayout.AddItem(innerLayout, 0, 2, false)

	frame := tview.NewFrame(outerLayout)
	frame.SetBorders(0, 0, 0, 0, 0, 0)
	frame.AddText("B H 🐙 D R", true, tview.AlignCenter, tcell.ColorOlive)
	frame.SetBackgroundColor(tcell.Color236)

	var logs *tview.TextView
	if showLogs {
		// create the logs view:
		logs = tview.NewTextView()
		logs.SetTitle("logs").SetBorder(true)
		outerLayout.AddItem(logs, 0, 2, false)
	}
	outerLayout.AddItem(statusbar, 1, 0, false)

	// create the app:
	app := tview.NewApplication()
	app.SetRoot(frame, true)
	app.SetFocus(switches)

	// for keeping track of vi-like key chords:
	chord := util.KeyChord{Active: false, Buffer: "", Action: ""}
	chordmap := config["chordmap"].(map[string]interface{})

	// switches keybindings:
	switches.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			selection := switches.GetCurrentNode()
			key := event.Rune()

			if event.Key() == tcell.KeyEsc {
				util.ResetChord(&chord)
			}

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
				case 'c', 'd', 'o', 'y', 'p': // runes that start a chord:
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
				case 'R': // refetch all states from HA.
					haCommands <- homeassistant.Command{
						Type: "get_states",
					}
				case ';': // toggle entity.
					haCommands <- homeassistant.Command{
						EntityID: selection.GetReference().(homeassistant.Data).EntityID,
						Service:  "toggle",
						Type:     "call_service",
						Domain:   true,
					}
				}
			}
			statusbar.SetText(chord.Buffer)
			return event
		},
	)

	// logs keybindings:
	if showLogs {
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
	}

	// global keybindings:
	app.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			key := event.Rune()

			switch key {
			case '[': // focus switches view.
				if showLogs {
					app.SetFocus(switches)
					switches.SetBorderColor(tcell.ColorGreen)
					logs.SetBorderColor(tcell.ColorWhite)
				}
			case ']': // focus logs view.
				if showLogs {
					app.SetFocus(logs)
					switches.SetBorderColor(tcell.ColorWhite)
					logs.SetBorderColor(tcell.ColorGreen)
				}
			case 'q': // quit the program.
				app.Stop()
			}
			return event
		},
	)

	// preselect node:
	switches.SetCurrentNode(switchesRoot)

	// connect to Home Assistant:
	go homeassistant.Connect(haConfig, haEvents, haCommands)

	// handle Home Assistant events:
	go func() {
		// TODO clean up this entire function.
		nodeFormat := "%s == %s"

		for {
			m := homeassistant.Message{}
			message := <-haEvents
			json.Unmarshal([]byte(message), &m)

			// handle event type messages:
			if m.Event.Type == "state_changed" {
				for _, node := range haEntities.GetChildren() {
					r := node.GetReference().(homeassistant.Data)
					if r.EntityID == m.Event.Data.EntityID {
						node.SetText(
							fmt.Sprintf(
								nodeFormat,
								r.NickName,
								m.Event.Data.NewState.State,
							),
						)
					}
				}
			}

			// handle result type messages:
			if len(m.Result) > 0 {
				for _, node := range haEntities.GetChildren() {
					r := node.GetReference().(homeassistant.Data)
					for _, result := range m.Result {
						if r.EntityID == result.EntityID {
							node.SetText(
								fmt.Sprintf(
									nodeFormat,
									r.NickName,
									result.State,
								),
							)
						}
					}
				}
			}

			// update logs view:
			if logs != nil {
				current := logs.GetText(true)
				if len(current) == 0 {
					logs.SetText(message)
				} else {
					logs.SetText(current + ",\n" + message) // append message.
				}
			}

			app.Draw() // external changes require a manual redraw.
		}
	}()

	// fetch all states at startup.
	haCommands <- homeassistant.Command{Type: "get_states"}

	app.Run()
}
