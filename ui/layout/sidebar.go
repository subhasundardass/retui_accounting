package layout

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
)

var sidebarTree = []components.TreeNode{
	{
		ID:    "dashboard",
		Label: "Dashboard",
		Children: []components.TreeNode{
			{
				ID:    "cmd-app",
				Label: "app",
				Children: []components.TreeNode{
					{ID: "home", Label: "Home"},
					{ID: "settings", Label: "Settings"},
					{ID: "about", Label: "About"},
				},
			},
		},
	},
	{
		ID:    "accounting",
		Label: "Accounting",
		Children: []components.TreeNode{
			{
				ID:    "accounting-ledger",
				Label: "Ledger",
				Children: []components.TreeNode{
					{ID: "ledger_group", Label: "Groups"},
					{ID: "ledger_list", Label: "Ledgers"},
				},
			},
			{
				ID:    "accounting-process",
				Label: "accounting",
				Children: []components.TreeNode{
					{ID: "journal_entry", Label: "Journal Entry"},
					// {ID: "receipt", Label: "Receipt"},
					// {ID: "payment", Label: "Payment"},
					// {ID: "sale", Label: "Sales"},
					// {ID: "purchase", Label: "Purchase"},
					// {ID: "contra", Label: "Contra"},
					// {ID: "credit_note", Label: "Credit Note"},
					// {ID: "debit_note", Label: "Debit Note"},
				},
			},
		},
	},
	{
		ID:    "report",
		Label: "Reports",
		Children: []components.TreeNode{
			{ID: "journal_report", Label: "Journals"},
			{ID: "cash_book", Label: "Cash Book"},
			{ID: "bank_book", Label: "Bank Book"},
			{ID: "trail_balance", Label: "Trail Balance"},
			{ID: "balance_sheet", Label: "Balance Sheet"},
			{ID: "party", Label: "Party",
				Children: []components.TreeNode{
					{ID: "sundry_debtors", Label: "Sundry Debtors"},
					{ID: "sundry_creditors", Label: "Sundry Debtors"},
				},
			},
		},
	},
	{
		ID:    "settings",
		Label: "Settings",
		Children: []components.TreeNode{
			{ID: "docs-readme", Label: "README.md"},
			{ID: "docs-api", Label: "API.md"},
		},
	},
	// {
	// 	ID:    "gomod",
	// 	Label: "go.mod",
	// },
	// {
	// 	ID:    "gosum",
	// 	Label: "go.sum",
	// },
}

func SidebarTree(ctx *context.AppContext, props retui.Props) retui.Element {
	// selected, setSelected := retui.UseState("")

	// Check if sidebar is focused
	isFocused := retui.IsFocused("sidebar")
	isLeafNode := make(map[string]bool)

	var findLeafNodes func([]components.TreeNode)
	findLeafNodes = func(nodes []components.TreeNode) {
		for _, node := range nodes {
			// If it has no children, it's a leaf node (the end of the branch)
			isLeafNode[node.ID] = len(node.Children) == 0

			if len(node.Children) > 0 {
				findLeafNodes(node.Children)
			}
		}
	}
	findLeafNodes(sidebarTree)

	return retui.Box(
		retui.Props{Direction: retui.Column, Padding: [4]int{0, 0, 0, 1}, Width: retui.Percent(15), Gap: 0},
		retui.NewStyle().Border(retui.Border{
			Top: true, Right: true, Bottom: true, Left: true,
			Chars: retui.BorderRounded, Color: retui.Gray(1),
		}),
		// retui.NewStyle().
		// 	Border(retui.BorderAll()).
		// Title(retui.Title{
		// 	Text:       "Bavigation",
		// 	Foreground: retui.Cyan,
		// 	Bold:       true,
		// 	Align:      retui.AlignStart,
		// }),

		// Sidebar tree panel
		retui.Box(
			retui.Props{Direction: retui.Column, Gap: 0},
			retui.NewStyle(),

			retui.Box(
				retui.Props{
					Justify: retui.JustifySpaceBetween,
					Padding: [4]int{0, 0, 1, 0}},
				retui.NewStyle(),
				retui.Text("Navigation", retui.NewStyle().Bold(true).Foreground(retui.Cyan)),
				retui.Text("[ F2 ]", retui.NewStyle().Foreground(retui.Cyan)),
			),

			components.Tree(
				"sidebar",
				sidebarTree,
				isFocused,
				func(id string) {

					// Check if the clicked ID is a leaf node
					if !isLeafNode[id] {
						// retui.Debug("Parent node clicked. Ignoring execution for ID: ", id)
						return
					}

					if ctx == nil {
						// retui.Debug(" ctx is nil!")
						return
					}

					if isFocused {
						retui.PushScreen(id)
						retui.SetFocus(id)
					}

				},
			),
		),
	)
}
