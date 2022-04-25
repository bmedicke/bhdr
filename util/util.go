package util

import "github.com/rivo/tview"

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
