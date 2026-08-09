package views

import (
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/ledger"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/retui/window"
)

type LedgerFormComponent struct {
	controller *ledger.LedgerController
	win        *window.Window
	state      ledger.LedgerState
	editing    bool
	editID     int
}

func NewLedgerFormComponent(ctx *appctx.AppContext) *LedgerFormComponent {
	return &LedgerFormComponent{
		controller: ledger.NewController(ctx),
	}
}

func (c *LedgerFormComponent) bindKeys(form *retui.Form[ledger.LedgerState]) {
	key := retui.CurrentKey
	if key == (retui.Key{}) || key.Consumed {
		return
	}
	if retui.CapturedFocus() != "" {
		return
	}

	v := form.Values()

	switch retui.CurrentKey.Code {
	case retui.KeyEscape:
		retui.PopScreen()

	case retui.KeyDown, retui.KeyTab:
		// 26 fields + 2 buttons = 28 total focusable elements
		v.FocusIndex = (v.FocusIndex + 1) % 28
		form.SetValuesSilent(v)

	case retui.KeyUp, retui.KeyShiftTab:
		v.FocusIndex = (v.FocusIndex - 1 + 28) % 28
		form.SetValuesSilent(v)
	}
}

func (c *LedgerFormComponent) LedgerEditForm(ctx *appctx.AppContext) retui.Element {
	ledgerID := 0
	params := retui.CurrentScreenParams()

	if params != nil {
		if id, ok := params["ledgerID"].(int); ok {
			ledgerID = id
		}
	}

	state, err := c.controller.GetLedger(ledgerID)
	if err != nil {
		components.ShowError(err.Error())
	}

	c.editing = true
	c.editID = ledgerID

	c.state = ledger.LedgerState{
		Mode:        ledger.ModeUpdate,
		Code:        state.Code,
		Name:        state.Name,
		Alias:       state.Alias,
		GroupID:     state.GroupID,
		PartyType:   state.PartyType,
		Description: state.Description,
		IsActive:    state.IsActive,
		//--
		AddressLine1:        state.AddressLine1,
		AddressLine2:        state.AddressLine2,
		GSTRegistrationType: state.GSTRegistrationType,
		GSTIN:               state.GSTIN,
		PAN:                 state.PAN,
		City:                state.City,
		State:               state.State,
		Country:             state.Country,
		Pincode:             state.Pincode,
		Phone:               state.Phone,
		Mobile:              state.Mobile,
		Email:               state.Email,
		ContactPerson:       state.ContactPerson,
		BankName:            state.BankName,
		BankAccountNo:       state.BankAccountNo,
		BankIFSC:            state.BankIFSC,
		BankBranch:          state.BankBranch,
	}

	return c.buildForm()
}

func (c *LedgerFormComponent) buildForm() retui.Element {
	form := retui.UseForm(c.state)
	v := form.Values()

	c.bindKeys(form)

	// ==================== BASIC INFORMATION ====================
	// Code - Focus Index 0
	code := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Code", retui.NewStyle()),
		),
		components.TextInput().
			ID("code").
			Value(v.Code).
			Width(20).
			Prefix(" : ").
			Focused(v.FocusIndex == 0).
			OnChange(func(id, value string) {
				form.SetField("Code", value)
			}).
			Render(),
	)

	// Name - Focus Index 1
	name := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(2)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Name", retui.NewStyle()),
		),
		components.TextInput().
			ID("name").
			Value(v.Name).
			Width(40).
			Prefix(" : ").
			Focused(v.FocusIndex == 1).
			OnChange(func(id, value string) {
				form.SetField("Name", value)
			}).
			Render(),
	)

	// Alias - Focus Index 2
	alias := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Alias", retui.NewStyle()),
		),
		components.TextInput().
			ID("alias").
			Value(v.Alias).
			Width(40).
			Prefix(" : ").
			Focused(v.FocusIndex == 2).
			OnChange(func(id, value string) {
				form.SetField("Alias", value)
			}).
			Render(),
	)

	// Group - Focus Index 3
	group := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Group", retui.NewStyle()),
		),
		components.SelectDropdown().
			ID("group").
			Focused(v.FocusIndex == 3).
			Prefix(" : ").
			Width(50).
			OverlayAbsPos(80, 5).
			OnFilter(func(id, query string) []components.SelectOption {
				return c.controller.LedgerGroupFilterOptions(query)
			}).
			Value(v.GroupID).
			OnChange(func(id, value string) {
				form.SetField("GroupID", value)
			}).
			Render(),
	)

	// Party Type - Focus Index 4
	partyType := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Party Type", retui.NewStyle()),
		),
		components.SelectDropdown().
			ID("partyType").
			Focused(v.FocusIndex == 4).
			OverlayAbsPos(80, 5).
			Width(50).
			Prefix(" : ").
			Options([]components.SelectOption{
				{Value: "CUSTOMER", Label: "Customer"},
				{Value: "SUPPLIER", Label: "Supplier"},
				{Value: "BOTH", Label: "Both"},
				{Value: "INTERNAL", Label: "Internal"},
			}).
			Value(v.PartyType).
			OnChange(func(id, value string) {
				form.SetField("PartyType", value)
			}).
			Render(),
	)

	// Description - Focus Index 5
	description := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Description", retui.NewStyle()),
		),
		components.TextArea().
			ID("description").
			Prefix(" : ").
			Value(v.Description).
			Height(1).
			Focused(v.FocusIndex == 5).
			OnChange(func(id, value string) {
				form.SetField("Description", value)
			}).
			Render(),
	)

	// ==================== ADDRESS ====================
	// Address Line 1 - Focus Index 6
	addressLine1 := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Line 1", retui.NewStyle()),
		),
		components.TextInput().
			ID("addressLine1").
			Value(v.AddressLine1).
			Prefix(" : ").
			Focused(v.FocusIndex == 6).
			OnChange(func(id, value string) {
				form.SetField("AddressLine1", value)
			}).
			Render(),
	)

	// Address Line 2 - Focus Index 7
	addressLine2 := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Line 2", retui.NewStyle()),
		),
		components.TextInput().
			ID("addressLine2").
			Prefix(" : ").
			Value(v.AddressLine2).
			Focused(v.FocusIndex == 7).
			OnChange(func(id, value string) {
				form.SetField("AddressLine2", value)
			}).
			Render(),
	)

	// City - Focus Index 8
	city := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("City", retui.NewStyle()),
		),
		components.TextInput().
			ID("city").
			Value(v.City).
			Prefix(" : ").
			Focused(v.FocusIndex == 8).
			OnChange(func(id, value string) {
				form.SetField("City", value)
			}).
			Render(),
	)

	// State - Focus Index 9
	state := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("State", retui.NewStyle()),
		),
		components.TextInput().
			ID("state").
			Prefix(" : ").
			Value(v.State).
			Focused(v.FocusIndex == 9).
			OnChange(func(id, value string) {
				form.SetField("State", value)
			}).
			Render(),
	)

	// Country - Focus Index 10
	country := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Country", retui.NewStyle()),
		),
		components.TextInput().
			ID("country").
			Value(v.Country).
			Prefix(" : ").
			Focused(v.FocusIndex == 10).
			OnChange(func(id, value string) {
				form.SetField("Country", value)
			}).
			Render(),
	)

	// Pincode - Focus Index 11
	pincode := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Pincode", retui.NewStyle()),
		),
		components.TextInput().
			ID("pincode").
			Value(v.Pincode).
			Prefix(" : ").
			Focused(v.FocusIndex == 11).
			OnChange(func(id, value string) {
				form.SetField("Pincode", value)
			}).
			Render(),
	)

	// ==================== CONTACT ====================
	// Phone - Focus Index 12
	phone := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Phone", retui.NewStyle()),
		),
		components.TextInput().
			ID("phone").
			Value(v.Phone).
			Prefix(" : ").
			Focused(v.FocusIndex == 12).
			OnChange(func(id, value string) {
				form.SetField("Phone", value)
			}).
			Render(),
	)

	// Mobile - Focus Index 13
	mobile := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Mobile", retui.NewStyle()),
		),
		components.TextInput().
			ID("mobile").
			Prefix(" : ").
			Value(v.Mobile).
			Focused(v.FocusIndex == 13).
			OnChange(func(id, value string) {
				form.SetField("Mobile", value)
			}).
			Render(),
	)

	// Email - Focus Index 14
	email := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Email", retui.NewStyle()),
		),
		components.TextInput().
			ID("email").
			Value(v.Email).
			Prefix(" : ").
			Focused(v.FocusIndex == 14).
			OnChange(func(id, value string) {
				form.SetField("Email", value)
			}).
			Render(),
	)

	// Contact Person - Focus Index 15
	contactPerson := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Contact Person", retui.NewStyle()),
		),
		components.TextInput().
			ID("contactPerson").
			Value(v.ContactPerson).
			Prefix(" : ").
			Focused(v.FocusIndex == 15).
			OnChange(func(id, value string) {
				form.SetField("ContactPerson", value)
			}).
			Render(),
	)

	// ==================== TAX REGISTRATION ====================
	// GST Registration Type - Focus Index 16
	gstRegistrationType := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Registration Type", retui.NewStyle()),
		),
		components.SelectDropdown().
			ID("gstRegistrationType").
			Focused(v.FocusIndex == 16).
			OverlayAbsPos(80, 5).
			Prefix(" : ").
			Options([]components.SelectOption{
				{Value: "REGULAR", Label: "Regular"},
				{Value: "COMPOSITION", Label: "Composition"},
				{Value: "UNREGISTERED", Label: "Unregistered"},
				{Value: "CONSUMER", Label: "Consumer"},
				{Value: "SEZ", Label: "SEZ"},
				{Value: "OVERSEAS", Label: "Overseas"},
			}).
			Value(v.GSTRegistrationType).
			OnChange(func(id, value string) {
				form.SetField("GSTRegistrationType", value)
			}).
			Render(),
	)

	// GSTIN - Focus Index 17
	gstin := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("GSTIN", retui.NewStyle()),
		),
		components.TextInput().
			ID("gstin").
			Value(v.GSTIN).
			Prefix(" : ").
			Focused(v.FocusIndex == 17).
			OnChange(func(id, value string) {
				form.SetField("GSTIN", value)
			}).
			Render(),
	)

	// PAN - Focus Index 18
	pan := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("PAN", retui.NewStyle()),
		),
		components.TextInput().
			ID("pan").
			Value(v.PAN).
			Prefix(" : ").
			Focused(v.FocusIndex == 18).
			OnChange(func(id, value string) {
				form.SetField("PAN", value)
			}).
			Render(),
	)

	// ==================== BANK DETAILS ====================
	// Bank Name - Focus Index 19
	bankName := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Bank Name", retui.NewStyle()),
		),
		components.TextInput().
			ID("bankName").
			Value(v.BankName).
			Width(80).
			Prefix(" : ").
			Focused(v.FocusIndex == 19).
			OnChange(func(id, value string) {
				form.SetField("BankName", value)
			}).
			Render(),
	)

	// Bank Account No - Focus Index 20
	bankAccountNo := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Account No", retui.NewStyle()),
		),
		components.TextInput().
			ID("bankAccountNo").
			Value(v.BankAccountNo).
			Prefix(" : ").
			Focused(v.FocusIndex == 20).
			OnChange(func(id, value string) {
				form.SetField("BankAccountNo", value)
			}).
			Render(),
	)

	// Bank IFSC - Focus Index 21
	bankIFSC := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("IFSC", retui.NewStyle()),
		),
		components.TextInput().
			ID("bankIFSC").
			Value(v.BankIFSC).
			Prefix(" : ").
			Focused(v.FocusIndex == 21).
			OnChange(func(id, value string) {
				form.SetField("BankIFSC", value)
			}).
			Render(),
	)

	// Bank Branch - Focus Index 22
	bankBranch := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Branch", retui.NewStyle()),
		),
		components.TextInput().
			ID("bankBranch").
			Value(v.BankBranch).
			Prefix(" : ").
			Focused(v.FocusIndex == 22).
			OnChange(func(id, value string) {
				form.SetField("BankBranch", value)
			}).
			Render(),
	)

	// ==================== STATUS ====================
	// IsActive - Focus Index 23
	isActive := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		components.Checkbox().
			ID("isActive").
			Label("Active").
			Checked(v.IsActive).
			Focused(v.FocusIndex == 23).
			OnChange(func(id string, checked bool) {
				form.SetField("IsActive", checked)
			}).
			Render(),
	)

	// ==================== BUTTONS ====================
	// Save Button - Focus Index 24
	saveButton := retui.Box(
		retui.Props{},
		retui.NewStyle().Background(retui.Gray(2)),

		components.Button().
			Label("Save").
			Focused(v.FocusIndex == 24).
			Render(),
	)

	// Reset Button - Focus Index 25
	resetButton := retui.Box(
		retui.Props{},
		retui.NewStyle().
			Background(retui.Gray(2)).
			Foreground(retui.White),
		components.Button().
			Label("Reset").
			Focused(v.FocusIndex == 25).
			Render(),
	)

	// Button container
	buttonContainer := retui.Box(
		retui.Props{
			Gap:     1,
			Align:   retui.AlignCenter,
			Justify: retui.JustifyCenter,
		},
		retui.NewStyle(),

		saveButton,
		resetButton,
	)

	// ==================== MAIN PANEL ====================
	panel := components.Panel().
		Header(retui.Text("Ledger Create", retui.NewStyle().Bold(true))).
		Children(
			retui.Box(
				retui.Props{
					Direction: retui.Column,
					Width:     retui.Grow(1),
					Align:     retui.AlignStretch,
					Padding:   [4]int{0, 2, 0, 2},
				},
				retui.NewStyle(),

				// BASIC INFORMATION - Group 1
				retui.Box(
					retui.Props{Gap: 1},
					retui.NewStyle(),
					code,
					name,
					alias,
				),
				retui.Box(
					retui.Props{Gap: 1},
					retui.NewStyle(),
					group,
					partyType,
				),
				description,
			),
		).
		DividerWithText("Address & Contact").
		Children(
			retui.Box(
				retui.Props{
					Gap:     2,
					Width:   retui.Grow(3),
					Align:   retui.AlignStretch,
					Padding: [4]int{1, 2, 1, 2},
				},
				retui.NewStyle(),
				retui.Box(
					retui.Props{
						Direction: retui.Column,
						Width:     retui.Grow(1),
					},
					retui.NewStyle(),
					addressLine1,
					addressLine2,
					city,
					state,
					country,
					pincode,
				),
				retui.Box(
					retui.Props{
						Direction: retui.Column,
						Width:     retui.Grow(1),
					},
					retui.NewStyle(),
					phone,
					mobile,
					email,
					contactPerson,
				),
				retui.Box(
					retui.Props{
						Direction: retui.Column,
						Width:     retui.Grow(1),
					},
					retui.NewStyle(),
					gstRegistrationType,
					gstin,
					pan,
					bankName,
					bankAccountNo,
					bankIFSC,
					bankBranch,
					isActive,
				),
			),
		).
		DividerWithText("Bank Details").
		Children(
			// Buttons at the bottom
			buttonContainer,
		).
		Render()

	return retui.Box(
		retui.Props{
			Direction: retui.Column,
		},
		retui.NewStyle(),
		panel,
	)
}
