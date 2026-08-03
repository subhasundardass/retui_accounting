// module/company/views/form.go
package views

import (
	"strconv"
	"strings"

	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/company"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/retui/window"
	"github.com/subhasundardass/retui/ui/widgets"
)

type FormComponent struct {
	controller *company.Controller
	win        *window.Window

	state    company.FormState
	setState func(company.FormState)

	editing bool
	editID  int

	// onSaved is called after a successful save so the caller (e.g. the
	// companies list) can refresh its data.
	onSaved func()
}

const totalFields = 18

func NewFormComponent(ctx *appctx.AppContext) *FormComponent {
	return &FormComponent{
		controller: company.NewController(ctx),
	}
}

// SetOnSaved registers a callback fired after a successful save.
func (c *FormComponent) SetOnSaved(fn func()) {
	c.onSaved = fn
}

// Reset clears all form fields back to zero values and clears edit mode.
// Call this before opening the window for a fresh "create" flow.
func (c *FormComponent) Reset() {
	c.state = company.FormState{}
	c.editing = false
	c.editID = 0
	if c.setState != nil {
		c.setState(c.state)
	}
}

// LoadForEdit fetches an existing company and populates the form state
// for editing. Call this before EditForm.
func (c *FormComponent) LoadForEdit(id int) error {
	comp, err := c.controller.Edit(id)
	if err != nil {
		return err
	}

	c.editing = true
	c.editID = id

	var countryID, stateID int
	if comp.Edges.CountryRef != nil {
		countryID = comp.Edges.CountryRef.ID
	}
	if comp.Edges.StateRef != nil {
		stateID = comp.Edges.StateRef.ID
	}

	c.state = company.FormState{
		Code:       comp.Code,
		Name:       comp.Name,
		LegalName:  comp.LegalName,
		Email:      comp.Email,
		Phone:      comp.Phone,
		Website:    comp.Website,
		Country:    countryID,
		State:      stateID,
		City:       comp.City,
		PostalCode: comp.PostalCode,
		Address:    comp.Address,
		TaxID:      comp.TaxID,
		GSTIN:      comp.Gstin,
		PAN:        comp.Pan,
		IsActive:   comp.Active,
	}
	// Do NOT call c.setState here — the form window hasn't rendered yet,
	// so setState (if non-nil) refers to a stale/unmounted hook context.
	// EditForm's SetRenderFn seeds UseState with c.state fresh, on open.
	return nil
}

func (c *FormComponent) CreateForm(ctx *appctx.AppContext) *window.Window {
	c.win = window.NewWindow().
		SetTitle("Create Company").
		SetModal(true).
		Center().
		SetSize(100, 40)

	c.win.SetRenderFn(func() retui.Element {
		c.state, c.setState = retui.UseState(company.FormState{})
		return c.buildWindow(ctx)
	})

	c.bindKeys()
	return c.win
}

func (c *FormComponent) EditForm(ctx *appctx.AppContext) *window.Window {

	c.win = window.NewWindow().
		SetTitle("Edit Company").
		SetModal(true).
		Center().
		SetSize(100, 40)

	c.win.SetRenderFn(func() retui.Element {
		// Seed UseState with the state already populated by LoadForEdit,
		// not a fresh zero-value — otherwise the loaded company data
		// would be wiped out on the first render.
		c.state, c.setState = retui.UseState(c.state)
		return c.buildWindow(ctx)
	})

	c.bindKeys()
	return c.win
}

func (c *FormComponent) bindKeys() {
	c.win.OnKeyPress(func(key retui.Key) bool {
		if retui.CapturedFocus() != "" {
			return false
		}

		switch key.Code {
		case retui.KeyDown:
			s := c.state
			s.FocusIndex = (s.FocusIndex + 1) % totalFields
			c.setState(s)
			return true

		case retui.KeyUp:
			s := c.state
			s.FocusIndex = (s.FocusIndex - 1 + totalFields) % totalFields
			c.setState(s)
			return true

		case retui.KeyEscape:
			c.win.Close()
			return true
		}
		return false
	})
}

func (c *FormComponent) save() {
	mode := company.ModeCreate
	id := 0
	if c.editing {
		mode = company.ModeUpdate
		id = c.editID
	}

	_, err := c.controller.Save(mode, id, c.state)
	if err != nil {
		retui.Debugf("Save failed: %v", err)
		components.ShowError("Save failed: " + err.Error())
		// TODO: surface this error in the UI (e.g. an Errors/status field on state)
		return
	}

	if c.onSaved != nil {
		c.onSaved()
	}
	c.win.Close()
}

func (c *FormComponent) buildWindow(ctx *appctx.AppContext) retui.Element {
	isFocused := func(index int) bool {
		return c.state.FocusIndex == index
	}

	row1 := retui.Box(
		retui.Props{
			Gap:     1,
			Padding: [4]int{0, 1, 0, 1},
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap:   1,
				Width: retui.Fixed(20),
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Width: retui.Fixed(5),
				},
				retui.NewStyle(),
				retui.Text("Code", retui.NewStyle()),
			),
			components.TextInput().
				ID("code").
				Focused(isFocused(0)).
				Value(c.state.Code).
				OnChange(func(id, value string) {
					var b strings.Builder
					for _, r := range strings.ToUpper(value) {
						if r >= 'A' && r <= 'Z' {
							b.WriteRune(r)
						}
					}
					s := c.state
					s.Code = b.String()
					c.setState(s)
				}).
				Render(),
		),

		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(10)},
				retui.NewStyle(),
				retui.Text("Name", retui.NewStyle()),
			),
			components.TextInput().
				ID("name").
				Width(retui.Fixed(40).Value).
				Focused(isFocused(1)).
				Value(c.state.Name).
				OnChange(func(id, value string) {
					s := c.state
					s.Name = value
					c.setState(s)
				}).
				Render(),
		),
		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(10)},
				retui.NewStyle(),
				retui.Text("Legal Name", retui.NewStyle()),
			),
			components.TextInput().
				ID("legal_name").
				Width(retui.Fixed(40).Value).
				Focused(isFocused(2)).
				Value(c.state.LegalName).
				OnChange(func(id, value string) {
					s := c.state
					s.LegalName = value
					c.setState(s)
				}).
				Render(),
		),
	)

	row2 := retui.Box(
		retui.Props{
			Gap:     1,
			Padding: [4]int{0, 1, 0, 1},
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(10)},
				retui.NewStyle(),
				retui.Text("Email", retui.NewStyle()),
			),
			components.TextInput().
				ID("email").
				Width(retui.Fixed(30).Value).
				Focused(isFocused(3)).
				Value(c.state.Email).
				OnChange(func(id, value string) {
					s := c.state
					s.Email = value
					c.setState(s)
				}).
				Render(),
		),

		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(10)},
				retui.NewStyle(),
				retui.Text("Phone", retui.NewStyle()),
			),
			components.TextInput().
				ID("phone").
				Width(retui.Fixed(30).Value).
				Focused(isFocused(4)).
				Value(c.state.Phone).
				OnChange(func(id, value string) {
					s := c.state
					s.Phone = value
					c.setState(s)
				}).
				Render(),
		),
		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(10)},
				retui.NewStyle(),
				retui.Text("Website", retui.NewStyle()),
			),
			components.TextInput().
				ID("website").
				Width(retui.Fixed(30).Value).
				Focused(isFocused(5)).
				Value(c.state.Website).
				OnChange(func(id, value string) {
					s := c.state
					s.Website = value
					c.setState(s)
				}).
				Render(),
		),
	)

	row3 := retui.Box(
		retui.Props{
			Gap:     1,
			Padding: [4]int{0, 1, 0, 1},
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(10)},
				retui.NewStyle(),
				retui.Text("Country", retui.NewStyle()),
			),
			widgets.CountryComponent(
				ctx,
				int(c.state.Country),
				30,
				isFocused(6),
				func(id, value string) {
					s := c.state
					i, err := strconv.Atoi(value)
					if err != nil {
						retui.Debugf("Invalid country ID: %v", err)
						return
					}
					s.Country = i
					c.setState(s)
				},
			),
		),
		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(10)},
				retui.NewStyle(),
				retui.Text("State", retui.NewStyle()),
			),
			widgets.StateComponent(
				c.state.Country,
				ctx,
				c.state.State,
				30,
				isFocused(7),
				func(id, value string) {
					s := c.state
					i, err := strconv.Atoi(value)
					if err != nil {
						retui.Debugf("Invalid State ID: %v", err)
						return
					}
					s.State = i
					c.setState(s)
				},
			),
		),
		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(10)},
				retui.NewStyle(),
				retui.Text("City", retui.NewStyle()),
			),
			components.TextInput().
				ID("city").
				Width(retui.Fixed(30).Value).
				Focused(isFocused(8)).
				Value(c.state.City).
				OnChange(func(id, value string) {
					s := c.state
					s.City = value
					c.setState(s)
				}).
				Render(),
		),
	)

	row4 := retui.Box(
		retui.Props{
			Gap:     1,
			Padding: [4]int{0, 1, 0, 1},
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(10)},
				retui.NewStyle(),
				retui.Text("Postal Code", retui.NewStyle()),
			),
			components.TextInput().
				ID("postal_code").
				Width(retui.Fixed(30).Value).
				Focused(isFocused(9)).
				Value(c.state.PostalCode).
				OnChange(func(id string, value string) {
					s := c.state
					s.PostalCode = value
					c.setState(s)
				}).
				Render(),
		),
		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(10)},
				retui.NewStyle(),
				retui.Text("Address", retui.NewStyle()),
			),
			components.TextInput().
				ID("address").
				Width(retui.Fixed(72).Value).
				Focused(isFocused(10)).
				Value(c.state.Address).
				OnChange(func(id, value string) {
					s := c.state
					s.Address = value
					c.setState(s)
				}).
				Render(),
		),
	)

	row10 := retui.Box(
		retui.Props{
			Gap:     1,
			Padding: [4]int{0, 1, 0, 1},
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(11)},
				retui.NewStyle(),
				retui.Text("Reg. No", retui.NewStyle()),
			),
			components.TextInput().
				ID("tax_id").
				Focused(isFocused(11)).
				Value(c.state.TaxID).
				OnChange(func(id, value string) {
					s := c.state
					s.TaxID = value
					c.setState(s)
				}).
				Render(),
		),
		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(10)},
				retui.NewStyle(),
				retui.Text("GSTIN", retui.NewStyle()),
			),
			components.TextInput().
				ID("gstin").
				Width(retui.Fixed(30).Value).
				Focused(isFocused(12)).
				Value(c.state.GSTIN).
				OnChange(func(id, value string) {
					s := c.state
					s.GSTIN = value
					c.setState(s)
				}).
				Render(),
		),
		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(10)},
				retui.NewStyle(),
				retui.Text("PAN", retui.NewStyle()),
			),
			components.TextInput().
				ID("pan").
				Width(retui.Fixed(30).Value).
				Focused(isFocused(13)).
				Value(c.state.PAN).
				OnChange(func(id, value string) {
					s := c.state
					s.PAN = value
					c.setState(s)
				}).
				Render(),
		),
	)

	return retui.Box(
		retui.Props{
			Direction: retui.Column,
			Padding:   [4]int{1, 1, 1, 1},
		},
		retui.NewStyle(),

		row1,
		retui.Box(
			retui.Props{Margin: [4]int{1, 0, 0, 0}},
			retui.NewStyle().Border(retui.Border{Top: true, Color: retui.Gray(2), Title: &retui.BorderTitle{
				Text:  "Contact",
				Align: retui.AlignStart,
				Style: retui.NewStyle().Foreground(retui.Teal),
			}}),
			row2,
		),
		row3,
		row4,

		retui.Box(
			retui.Props{Margin: [4]int{1, 0, 0, 0}},
			retui.NewStyle().Border(retui.Border{Top: true, Color: retui.Gray(2), Title: &retui.BorderTitle{
				Text:  "Tax Information",
				Align: retui.AlignStart,
				Style: retui.NewStyle().Foreground(retui.Teal),
			}}),
			row10,
		),

		retui.Box(
			retui.Props{Margin: [4]int{1, 0, 0, 0}},
			retui.NewStyle().Border(retui.Border{Top: true, Color: retui.Gray(2), Title: &retui.BorderTitle{
				Text:  "Actions",
				Align: retui.AlignStart,
				Style: retui.NewStyle().Foreground(retui.Teal),
			}}),

			retui.Box(
				retui.Props{
					Gap:     1,
					Justify: retui.JustifySpaceBetween,
					Width:   retui.Grow(1),
				},
				retui.NewStyle(),
				components.Checkbox().
					ID("IsActive").
					Label("Is Active").
					Checked(c.state.IsActive).
					Focused(isFocused(14)).
					OnChange(func(id string, checked bool) {
						s := c.state
						s.IsActive = checked
						c.setState(s)
					}).
					Render(),

				retui.Box(
					retui.Props{Gap: 1},
					retui.NewStyle(),
					components.Button().
						ID("save").
						Label("Save").
						Focused(isFocused(15)).
						Style(retui.NewStyle().Background(retui.Gray(2)).Foreground(retui.BrightWhite)).
						OnKeyPress(func(id string, key retui.Key) bool {
							if key.Code == retui.KeyEnter {
								c.save()
								return true
							}
							return false
						}).
						Render(),

					components.Button().
						ID("cancel").
						Label("Cancel").
						Focused(isFocused(16)).
						Style(retui.NewStyle().Background(retui.Gray(2)).Foreground(retui.BrightWhite)).
						OnKeyPress(func(id string, key retui.Key) bool {
							if key.Code == retui.KeyEnter {
								c.win.Close()
								return true
							}
							return false
						}).
						Render(),

					components.Button().
						ID("reset").
						Focused(isFocused(17)).
						Style(retui.NewStyle().Background(retui.Gray(2)).Foreground(retui.BrightWhite)).
						OnKeyPress(func(id string, key retui.Key) bool {
							if key.Code == retui.KeyEnter {
								c.Reset()
								return true
							}
							return false
						}).
						Render(),
				),
			),
		),
	)
}
