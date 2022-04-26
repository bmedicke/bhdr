package util

import (
	"errors"
	"fmt"
	"os"

	"github.com/rivo/tview"
)

// GetParent returns parent or nil (if it was not found).
// https://github.com/rivo/tview/issues/246#issuecomment-471173854
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

// IntuitiveViBindings provide H, L, h, l bindings that operate on a TreeView.
// H: recursively collapse all nodes but root.
// L: recursively expand all nodes.
// h: collapse OR move up (if already collapsed or not folder)
// l: expand node without switching.
func IntuitiveViBindings(rune int32, view *tview.TreeView) {
	selection := view.GetCurrentNode()

	switch rune {
	case 'H':
		// calling .CollapseAll() on the children of rootNode
		// does not work for some reason. do it manually:
		view.GetRoot().
			Walk(func(node, parent *tview.TreeNode) bool {
				// ignore root node:
				if parent != nil {
					node.Collapse()
				}
				return true // visit all nodes.
			})
	case 'L':
		view.GetRoot().ExpandAll()
	case 'h':
		if selection.IsExpanded() && nil != selection.GetChildren() {
			selection.Collapse()
		} else if selection.GetLevel() > 1 {
			parent := GetParent(selection, view.GetRoot())
			view.SetCurrentNode(parent)
		}
	case 'l':
		if !selection.IsExpanded() {
			selection.Expand()
		}
	}
}

// OverwriteFile replaces the content of a file.
// If the file does not exist, it will be created.
func OverwriteFile(file string, content string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	f.WriteString(content)
	defer f.Close()
	return nil
}

// CreateFileIfNotExist creates a file with a string as content.
// Returns an error if it already exists.
func CreateFileIfNotExist(file string, content string) error {
	_, err := os.Stat(file)
	// create file if it does not alreay exist:
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(file)
		if err != nil {
			return err
		}
		f.WriteString(content)
		defer f.Close()
	} else {
		return fmt.Errorf("file %v already present", file)
	}
	return nil
}
