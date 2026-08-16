package layout

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/ui/widgets"
)

func ShortcutPanel(ctx *context.AppContext, props retui.Props) retui.Element {

	return retui.Box(
		retui.Props{
			Width:  retui.Percent(15),
			Height: retui.Grow(1),
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Width: retui.Grow(1),
			},
			retui.NewStyle(),
			components.Panel().
				Header(retui.Text(" Shortcut", retui.NewStyle().Bold(true))).
				Children(
					retui.Box(
						retui.Props{
							Gap:       0,
							Direction: retui.Column,
							Width:     retui.Grow(1),
							Padding:   [4]int{0, 1, 0, 1},
						},
						retui.NewStyle(),

						// Company
						widgets.ShortCutItem("Select Company", "F1"),
						widgets.ShortCutItem("Change Data", "F2"),
						widgets.ShortCutItem("Contra Voucher", "F4"),
						// Data
						widgets.ShortCutItem("Change Data", "F5"),
					),
				).
				Divider().
				Children(
					retui.Box(
						retui.Props{
							Gap:       0,
							Direction: retui.Column,
							Width:     retui.Grow(1),
							Padding:   [4]int{0, 1, 0, 1},
						},
						retui.NewStyle(),

						// Vouchers
						widgets.ShortCutItem("Contra Voucher", "F3"),
						widgets.ShortCutItem("Payment Voucher", "F3"),
						widgets.ShortCutItem("Receipt Voucher", "F3"),
						widgets.ShortCutItem("Journal Voucher", "F3"),

						// Sales / Purchase
						widgets.ShortCutItem("Sale Voucher", "F4"),
						widgets.ShortCutItem("Purchase Voucher", "F4"),
					),
				).
				Divider().
				Children(
					retui.Box(
						retui.Props{
							Gap:       0,
							Direction: retui.Column,
							Width:     retui.Grow(1),
							Padding:   [4]int{0, 1, 0, 1},
						},
						retui.NewStyle(),

						// Features / Configuration
						widgets.ShortCutItem("Company Features", "F5"),
						widgets.ShortCutItem("Configuration", "F5"),
					),
				).
				Divider().
				Children(
					retui.Box(
						retui.Props{
							Gap:       0,
							Direction: retui.Column,
							Width:     retui.Grow(1),
							Padding:   [4]int{0, 1, 0, 1},
						},
						retui.NewStyle(),

						// Features / Configuration
						widgets.ShortCutItem("Company Features", "F5"),
					),
				).
				Render(),
		),
	)
}
