package widgets

import (
	"fmt"
	"strconv"

	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/internal/config"
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/company"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/window"
)

func CompanySwitcher(
	ctx *appctx.AppContext,
	controller *company.Controller,
	id string,
	onChange func(id, value string),
) retui.Element {

	if ctx == nil || controller == nil {
		return retui.Text(
			"Company: unavailable",
			retui.NewStyle().Foreground(retui.Red),
		)
	}

	companies, err := controller.List()
	if err != nil {
		retui.Debugf(
			"CompanySwitcher: failed to load companies: %v",
			err,
		)

		return retui.Text(
			"Company: error",
			retui.NewStyle().Foreground(retui.Red),
		)
	}

	selectedID := ctx.CompanyID()
	selected := findCompany(companies, selectedID)

	companyName := "Select Company"

	if selected != nil {
		companyName = selected.Name
	}

	// F1 opens the selector.
	if retui.CurrentKey.Code == retui.KeyF1 {
		openCompanyModal(
			ctx,
			companies,
			selectedID,
			id,
			onChange,
		)
	}

	return retui.Text(
		fmt.Sprintf("%s  [F1]", companyName),
		retui.NewStyle().Bold(true),
	)
}

func openCompanyModal(
	ctx *appctx.AppContext,
	companies []*ent.Company,
	selectedID int,
	id string,
	onChange func(id, value string),
) {
	if len(companies) == 0 {
		return
	}

	selectedIndex := 0

	for i, c := range companies {
		if c != nil && c.ID == selectedID {
			selectedIndex = i
			break
		}
	}

	win := window.NewWindow().
		SetTitle("Select Company").
		SetModal(true).
		SetSize(60, len(companies)+6).
		Center()

	win.SetRenderFn(func() retui.Element {
		return renderCompanyModal(
			companies,
			selectedIndex,
		)
	})

	win.OnKeyPress(func(key retui.Key) bool {

		switch key.Code {

		case retui.KeyEscape:
			win.Close()
			return true

		case retui.KeyUp:
			if selectedIndex > 0 {
				selectedIndex--
			}
			return true

		case retui.KeyDown:
			if selectedIndex < len(companies)-1 {
				selectedIndex++
			}
			return true

		case retui.KeyEnter:
			selected := companies[selectedIndex]

			if selected == nil {
				return true
			}

			// Update current application company.
			ctx.SetCompanyID(selected.ID)

			// Persist selection for next application start.
			if err := config.SaveLastCompany(selected.ID); err != nil {
				retui.Debugf(
					"CompanySwitcher: failed to save last company: %v",
					err,
				)
			}

			if onChange != nil {
				onChange(
					id,
					strconv.Itoa(selected.ID),
				)
			}

			win.Close()
			return true
		}

		return false
	})

	win.Show()
}

func renderCompanyModal(
	companies []*ent.Company,
	selectedIndex int,
) retui.Element {

	items := make([]retui.Element, 0, len(companies))

	for i, company := range companies {
		if company == nil {
			continue
		}

		prefix := "  "
		style := retui.NewStyle()

		if i == selectedIndex {
			prefix = "▶ "
			style = style.Bold(true)
		}

		items = append(
			items,
			retui.Text(
				prefix+company.Name,
				style,
			),
		)
	}

	companyList := retui.Box(
		retui.Props{
			Direction: retui.Column,
			Width:     retui.Grow(1),
			Height:    retui.Grow(1),
		},
		retui.NewStyle(),
		items...,
	)

	footer := retui.Box(
		retui.Props{
			Width:  retui.Grow(1),
			Height: retui.Fit(),
		},
		retui.NewStyle().Foreground(retui.Gray(5)),
		retui.Text(
			"↑ ↓ Navigate   Enter Select   Esc Close",
			retui.NewStyle(),
		),
	)

	return retui.Box(
		retui.Props{
			Direction: retui.Column,
			Width:     retui.Grow(1),
			Height:    retui.Grow(1),
			Padding:   [4]int{1, 2, 0, 2},
		},
		retui.NewStyle(),
		companyList,
		footer,
	)
}

func findCompany(
	companies []*ent.Company,
	id int,
) *ent.Company {

	for _, company := range companies {
		if company != nil && company.ID == id {
			return company
		}
	}

	return nil
}
