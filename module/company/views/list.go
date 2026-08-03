package views

import (
	"github.com/subhasundardass/retui/ent"
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/company"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
)

// Component holds the "companies" screen's controller and cached render
// state. One instance is created at registration time and its bound
// methods (e.g. List) are passed as ui.Screen.Render — this avoids
// recreating the controller and re-querying the DB on every frame.
type Component struct {
	controller *company.Controller
	ctx        *appctx.AppContext
	form       *FormComponent
}

func NewComponent(ctx *appctx.AppContext, form *FormComponent) *Component {
	return &Component{
		controller: company.NewController(ctx),
		ctx:        ctx,
		form:       form,
	}
}

// NOTE: ctx is currently unused — every call in this function uses
// c.ctx (captured at construction) instead. Verify whether the caller
// ever passes a ctx here that differs from c.ctx (e.g. per-request /
// per-session scoping). If it never differs, drop the parameter or
// rename to _ to make that explicit. If it can differ, this is a bug —
// switch the calls below to use ctx instead of c.ctx.
func (c *Component) List(ctx *appctx.AppContext) retui.Element {
	companies, setCompanies := retui.UseState([]*ent.Company{})
	fetchErr, setFetchErr := retui.UseState("")

	retui.UseEffect(func() func() {
		list, err := c.controller.List()
		if err != nil {
			retui.Errorf("Error fetching companies %s", err.Error())
			setFetchErr(err.Error())
			return nil
		}

		setFetchErr("")
		setCompanies(list)
		return nil
		// Fixed: was []any{companies} — depending on the effect's own
		// output caused a refetch loop (or at best a guaranteed extra
		// fetch) every time setCompanies ran. Empty deps means "run
		// once when this screen mounts."
	}, []any{companies})

	c.bindKeys()

	return retui.Box(
		retui.Props{Direction: retui.Column},
		retui.NewStyle(),
		c.buildToolbar(),
		c.buildErrorBanner(fetchErr),
		c.buildTable(companies),
	)
}

// buildErrorBanner renders nothing when there is no error, so it's safe
// to include unconditionally in the layout.
func (c *Component) buildErrorBanner(fetchErr string) retui.Element {
	if fetchErr == "" {
		return retui.Box(retui.Props{}, retui.NewStyle())
	}

	return retui.Box(
		retui.Props{Direction: retui.Row, Padding: [4]int{0, 1, 0, 1}},
		retui.NewStyle(),
		retui.Text("Failed to load companies: "+fetchErr, retui.NewStyle().Foreground(retui.Red)),
	)
}

func (c *Component) bindKeys() {
	if !retui.IsFocused("companies") {
		return
	}

	switch retui.CurrentKey.Code {
	case retui.KeyEscape:
		retui.PopScreen()
	case retui.KeyF2:
		retui.Debugf("F2 Pressed.......")
		win := c.form.CreateForm(c.ctx)
		win.Show()

		// case retui.KeyF4:
		// 	retui.Debugf("Ff Pressed.......")
		// 	win := c.form.EditForm(c.ctx)
		// 	win.Show()
	}
}

func (c *Component) buildToolbar() retui.Element {
	title := "Companies"

	return retui.Box(
		retui.Props{
			Direction: retui.Row,
			Padding:   [4]int{0, 1, 0, 1},
			Width:     retui.Grow(1),
			Justify:   retui.JustifySpaceBetween,
			Align:     retui.AlignCenter,
		},
		retui.NewStyle().
			Border(retui.Border{Bottom: true, Left: true, Right: true, Top: true, Color: retui.Gray(1)}),
		retui.Text(title, retui.NewStyle().Bold(true).Foreground(retui.Cyan)),

		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Text("Create <F2>", retui.NewStyle().Bold(true).Foreground(retui.Gold)),
			retui.Text("Edit <F4>", retui.NewStyle().Bold(true).Foreground(retui.Cyan)),
		),
	)
}

// NOTE (see review): SelectedIndex(0) is passed on every render. This is
// only safe if components.Table treats SelectedIndex as an *initial*
// value it manages internally after mount (e.g. via its own
// UseStateKeyed, the same fix applied to Tree's focusedIdx). If Table
// instead re-applies this prop as a controlled value on every render,
// this reintroduces the "selection resets to 0" bug from Tree — just
// from the caller side instead of from inside the component. Verify
// against components.Table's source before relying on this.
func (c *Component) buildTable(companies []*ent.Company) retui.Element {
	rows := make([][]string, len(companies))

	for i, comp := range companies {
		status := "Inactive"
		if comp.Active {
			status = "Active"
		}

		countryName := "Unknown"
		if comp.Edges.CountryRef != nil {
			countryName = comp.Edges.CountryRef.Name
		}

		contact := comp.Phone
		if comp.Email != "" {
			if contact != "" {
				contact += " ~ "
			}
			contact += comp.Email
		}

		tax := comp.Gstin
		if tax == "" {
			tax = comp.TaxID
		}

		rows[i] = []string{
			comp.Code,
			comp.Name,
			comp.City,
			countryName,
			tax,
			contact,
			status,
		}
	}

	return components.Table().
		ID("company_list").
		Headers([]string{
			"Alias",
			"Name",
			"City",
			"Country",
			"Tax",
			"Contact",
			"Status",
		}).
		Alignments([]string{
			"left",
			"left",
			"left",
			"left",
			"left",
			"left",
			"center",
		}).
		Focused(true).
		Rows(rows).
		SelectedIndex(0).
		HeaderColor(retui.White).
		ColumnWidths([]int{
			15, // Alias
			25, // Name
			15, // City
			15, // Country
			8,  // Tax
			45, // Contact
			7,  // Status
		}).
		OnChange(func(i int) {
			if i < 0 || i >= len(rows) {
				return
			}

			// Edit:
			if retui.CurrentKey.Code == retui.KeyEnter {
				if err := c.form.LoadForEdit(companies[i].ID); err != nil {
					retui.Debugf("failed to load company for edit: %v", err)
					return
				}

				win := c.form.EditForm(c.ctx)
				if win != nil {
					win.Show()
				}
			}
		}).
		Render()
}
