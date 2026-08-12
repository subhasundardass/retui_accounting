package companies

// import (
// 	"github.com/subhasundardass/retui/ent"
// 	"github.com/subhasundardass/retui/retui"
// 	"github.com/subhasundardass/retui/retui/components"
// 	createcompany "github.com/subhasundardass/retui/ui/windows/create_company"
// )

// type Components struct {
// 	controller *Controller
// }

// func NewComponents(controller *Controller) *Components {
// 	return &Components{controller: controller}
// }

// func (c *Components) bindKeys() {

// 	if retui.IsFocused("companies") {
// 		switch retui.CurrentKey.Code {
// 		case retui.KeyF2:
// 			retui.Debugf("F2 Pressed.......")
// 			win := createcompany.Window(c.controller.ctx)
// 			win.Show()
// 		}
// 	}
// }

// func (c *Components) RenderScreen() retui.Element {

// 	companies, setCompanies := retui.UseState([]*ent.Company{})

// 	retui.UseEffect(func() func() {
// 		list, err := c.controller.List()
// 		if err != nil {
// 			retui.Errorf("Error fetching companies %s", err.Error())
// 			return nil
// 		}

// 		setCompanies(list)
// 		return nil
// 	}, []any{})

// 	//--Key bindings
// 	c.bindKeys()

// 	return retui.Box(
// 		retui.Props{Direction: retui.Column, Gap: 0},
// 		retui.NewStyle(),

// 		c.buildHeader(),
// 		c.buildTable(companies),
// 	)
// }

// func (c *Components) buildHeader() retui.Element {

// 	title := "Companies"
// 	// if selected != nil {
// 	// 	title = fmt.Sprintf("Ledgers  %s", selected.Name)
// 	// }

// 	return retui.Box(
// 		retui.Props{
// 			Direction: retui.Row,
// 			Padding:   [4]int{0, 1, 0, 1},
// 			Width:     retui.Grow(1),
// 			Justify:   retui.JustifySpaceBetween,
// 			Align:     retui.AlignCenter,
// 		},
// 		retui.NewStyle().
// 			Border(retui.Border{Bottom: true, Left: true, Right: true, Top: true, Color: retui.Gray(1)}),
// 		retui.Text(title, retui.NewStyle().Bold(true).Foreground(retui.Cyan)),
// 		retui.Text("Create <F2>", retui.NewStyle().Bold(true).Foreground(retui.Gold)),
// 	)
// }

// func (c *Components) buildTable(companies []*ent.Company) retui.Element {

// 	rows := make([][]string, len(companies))

// 	for i, company := range companies {
// 		status := "Inactive"
// 		if company.Active {
// 			status = "Active"
// 		}

// 		contact := company.Phone
// 		if company.Email != "" {
// 			if contact != "" {
// 				contact += " / "
// 			}
// 			contact += company.Email
// 		}

// 		tax := company.Gstin
// 		if tax == "" {
// 			tax = company.TaxID
// 		}

// 		rows[i] = []string{
// 			company.Code,
// 			company.Name,
// 			company.City,
// 			// company.Country,
// 			tax,
// 			contact,
// 			status,
// 		}
// 	}

// 	return components.Table().
// 		ID("company_list").
// 		Headers([]string{
// 			"Alias",
// 			"Name",
// 			"City",
// 			"Country",
// 			"Tax",
// 			"Contact",
// 			"Status",
// 		}).
// 		Alignments([]string{
// 			"left",
// 			"left",
// 			"left",
// 			"left",
// 			"left",
// 			"left",
// 			"center",
// 		}).
// 		Focused(true).
// 		Rows(rows).
// 		SelectedIndex(0).
// 		HeaderColor(retui.Cyan).
// 		ColumnWidths([]int{
// 			15, // Alias
// 			25, // Name
// 			15, // City
// 			15, // Country
// 			20, // Tax
// 			30, // Contact
// 			10, // Status
// 		}).
// 		OnChange(func(i int) {
// 			if i < 0 || i >= len(companies) {
// 				return
// 			}

// 			// Example:
// 			// if retui.CurrentKey.Code == retui.KeyEnter {
// 			//     c.controller.ShowCompany(companies[i].ID)
// 			// }
// 		}).
// 		Render()
// }
