package example

import (
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
)

var sidebarTree = []components.TreeNode{
	{
		ID:    "example",
		Label: "Example",
		Children: []components.TreeNode{
			{ID: "basic-inputs", Label: "Basic Inputs"},
			{ID: "list", Label: "List"},
			{ID: "windows", Label: "Windows"},
			{ID: "other", Label: "Other"},
			{ID: "state", Label: "State"},
		},
	},
}

// Sidebar renders the example navigation tree.
func Sidebar() retui.Element {
	isFocused := retui.IsFocused("sidebar")
	isLeafNode := make(map[string]bool)

	var findLeafNodes func([]components.TreeNode)
	findLeafNodes = func(nodes []components.TreeNode) {
		for _, node := range nodes {
			isLeafNode[node.ID] = len(node.Children) == 0
			if len(node.Children) > 0 {
				findLeafNodes(node.Children)
			}
		}
	}
	findLeafNodes(sidebarTree)

	return retui.Box(
		retui.Props{
			Direction: retui.Column,
			Padding:   [4]int{1, 0, 0, 1},
			Width:     retui.Grow(1),
			Height:    retui.Grow(1),
		},
		retui.NewStyle().Border(retui.Border{
			Top: true, Right: true, Bottom: true, Left: true,
			Chars: retui.BorderRounded,
			Color: retui.Red,
			Title: &retui.BorderTitle{
				Text:  "Navigation",
				Style: retui.NewStyle().Foreground(retui.Yellow).Bold(true),
				Align: retui.AlignCenter,
			},
		}),
		components.Tree(
			"sidebar",
			sidebarTree,
			isFocused,
			func(id string) {
				if !isLeafNode[id] || !isFocused {
					return
				}

				retui.PushScreen(id)
				retui.SetFocus("content")
			},
		),
	)
}
