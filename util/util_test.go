package util

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestResetChord(t *testing.T) {
	chord := KeyChord{
		Active: true,
		Buffer: "buffer",
		Action: "action",
	}
	resetChord(&chord)

	if chord.Active != false {
		t.Fatal("Active not false")
	}
	if chord.Action != "" {
		t.Fatal("Action not empty")
	}
	if chord.Buffer != "" {
		t.Fatal("Buffer not empty")
	}
}

func TestHandleChordsInvalidNomen(t *testing.T) {
	keyrune := 'x'
	chord := KeyChord{}
	var chordmap map[string]interface{}
	testJSON := `{"c": {"c": "toggle:power"}}`
	json.Unmarshal([]byte(testJSON), &chordmap)

	err := HandleChords(keyrune, &chord, chordmap)
	expectedError := fmt.Errorf("invalid nomen [x]")
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error: '%v', got '%v'", expectedError, err)
	}

	expectedAction := ""
	if chord.Action != expectedAction {
		t.Errorf(
			"chord.Action should be '%v', got '%v'",
			expectedAction,
			chord.Action,
		)
	}

	expectedActive := false
	if chord.Active != expectedActive {
		t.Errorf(
			"chord.Active should be '%v', got '%v'",
			expectedActive,
			chord.Active,
		)
	}

	expectedBuffer := ""
	if chord.Buffer != expectedBuffer {
		t.Errorf(
			"chord.Buffer should be '%v', got '%v'",
			expectedBuffer,
			chord.Buffer,
		)
	}
}

func TestHandleChordsValidNomenVerb(t *testing.T) {
	keyrune := 'c'
	chord := KeyChord{}
	var chordmap map[string]interface{}
	testJSON := `{"c": {"c": "toggle:power"}}`
	json.Unmarshal([]byte(testJSON), &chordmap)

	// first call:
	err := HandleChords(keyrune, &chord, chordmap)
	if err != nil {
		t.Errorf("got unexpected error: '%v'", err)
	}

	expectedAction := ""
	if chord.Action != expectedAction {
		t.Errorf(
			"chord.Action should be '%v', got '%v'",
			expectedAction,
			chord.Action,
		)
	}

	expectedActive := true
	if chord.Active != expectedActive {
		t.Errorf(
			"chord.Active should be '%v', got '%v'",
			expectedActive,
			chord.Active,
		)
	}

	expectedBuffer := "c"
	if chord.Buffer != expectedBuffer {
		t.Errorf(
			"chord.Buffer should be '%v', got '%v'",
			expectedBuffer,
			chord.Buffer,
		)
	}

	// second call:
	err = HandleChords(keyrune, &chord, chordmap)
	if err != nil {
		t.Errorf("got unexpected error: '%v'", err)
	}

	expectedAction = "toggle:power"
	if chord.Action != expectedAction {
		t.Errorf(
			"chord.Action should be '%v', got '%v'",
			expectedAction,
			chord.Action,
		)
	}

	expectedActive = false
	if chord.Active != expectedActive {
		t.Errorf(
			"chord.Active should be '%v', got '%v'",
			expectedActive,
			chord.Active,
		)
	}

	expectedBuffer = ""
	if chord.Buffer != expectedBuffer {
		t.Errorf(
			"chord.Buffer should be '%v', got '%v'",
			expectedBuffer,
			chord.Buffer,
		)
	}
}

func TestHandleChordsValidNomenInvalidVerb(t *testing.T) {
	keyrune1 := 'c'
	keyrune2 := 'x'
	chord := KeyChord{}
	var chordmap map[string]interface{}
	testJSON := `{"c": {"c": "toggle:power"}}`
	json.Unmarshal([]byte(testJSON), &chordmap)

	// first call:
	HandleChords(keyrune1, &chord, chordmap)

	// second call:
	expectedError := fmt.Errorf("invalid verb [x]")
	err := HandleChords(keyrune2, &chord, chordmap)
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error: '%v', got '%v'", expectedError, err)
	}

	expectedAction := ""
	if chord.Action != expectedAction {
		t.Errorf(
			"chord.Action should be '%v', got '%v'",
			expectedAction,
			chord.Action,
		)
	}

	expectedActive := false
	if chord.Active != expectedActive {
		t.Errorf(
			"chord.Active should be '%v', got '%v'",
			expectedActive,
			chord.Active,
		)
	}

	expectedBuffer := ""
	if chord.Buffer != expectedBuffer {
		t.Errorf(
			"chord.Buffer should be '%v', got '%v'",
			expectedBuffer,
			chord.Buffer,
		)
	}
}

func TestHandleChordsValidNomenVerbPost(t *testing.T) {
	keyrune1 := 'c'
	keyrune2 := 'b'
	keyrune3 := '5'
	chord := KeyChord{}
	var chordmap map[string]interface{}
	testJSON := `{"c": {"b": "set:brightness:#"}}`
	json.Unmarshal([]byte(testJSON), &chordmap)

	// first call:
	HandleChords(keyrune1, &chord, chordmap)

	// second call:
	HandleChords(keyrune2, &chord, chordmap)

	// third call:
	err := HandleChords(keyrune3, &chord, chordmap)
	if err != nil {
		t.Errorf("got unexpected error: '%v'", err)
	}

	expectedAction := "set:brightness:#5"
	if chord.Action != expectedAction {
		t.Errorf(
			"chord.Action should be '%v', got '%v'",
			expectedAction,
			chord.Action,
		)
	}

	expectedActive := false
	if chord.Active != expectedActive {
		t.Errorf(
			"chord.Active should be '%v', got '%v'",
			expectedActive,
			chord.Active,
		)
	}

	expectedBuffer := ""
	if chord.Buffer != expectedBuffer {
		t.Errorf(
			"chord.Buffer should be '%v', got '%v'",
			expectedBuffer,
			chord.Buffer,
		)
	}
}

func TestHandleChordsValidNomenVerbInvalidPost(t *testing.T) {
	keyrune1 := 'c'
	keyrune2 := 'b'
	keyrune3 := 'x'
	chord := KeyChord{}
	var chordmap map[string]interface{}
	testJSON := `{"c": {"b": "set:brightness:#"}}`
	json.Unmarshal([]byte(testJSON), &chordmap)

	// first call:
	HandleChords(keyrune1, &chord, chordmap)

	// second call:
	HandleChords(keyrune2, &chord, chordmap)

	// third call:
	expectedError := fmt.Errorf("invalid value [x]")
	err := HandleChords(keyrune3, &chord, chordmap)
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error: '%v', got '%v'", expectedError, err)
	}

	expectedAction := ""
	if chord.Action != expectedAction {
		t.Errorf(
			"chord.Action should be '%v', got '%v'",
			expectedAction,
			chord.Action,
		)
	}

	expectedActive := false
	if chord.Active != expectedActive {
		t.Errorf(
			"chord.Active should be '%v', got '%v'",
			expectedActive,
			chord.Active,
		)
	}

	expectedBuffer := ""
	if chord.Buffer != expectedBuffer {
		t.Errorf(
			"chord.Buffer should be '%v', got '%v'",
			expectedBuffer,
			chord.Buffer,
		)
	}
}
