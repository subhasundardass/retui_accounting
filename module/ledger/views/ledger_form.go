package views

import (
	"strconv"

	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/util"
	"github.com/subhasundardass/retui/module/ledger"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/retui/window"
	"github.com/subhasundardass/retui/ui/widgets"
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

func (c *LedgerFormComponent) save(state ledger.LedgerState) {

	mode := ledger.ModeCreate
	id := 0
	if c.editing {
		mode = ledger.ModeUpdate
		id = c.editID
	}

	_, err := c.controller.LedgerSave(mode, id, state)
	if err != nil {
		retui.Debugf("Save failed: %v", err)
		components.ShowError("Save failed: " + err.Error())
		// TODO: surface this error in the UI (e.g. an Errors/status field on state)
		return
	}

	components.ShowSuccess("Ladger Saved ")
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
		StateID:             state.StateID,
		CountryID:           state.CountryID,
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

	return c.buildForm(ctx)
}

func (c *LedgerFormComponent) LedgerCreateForm(ctx *appctx.AppContext) retui.Element {

	c.editing = false
	c.editID = 0

	c.state = ledger.LedgerState{
		Mode: ledger.ModeUpdate,
	}

	return c.buildForm(ctx)
}

func (c *LedgerFormComponent) buildForm(ctx *appctx.AppContext) retui.Element {
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
				if err := form.SetField("Code", value); err != nil {
					// surface to the form's error state so the user sees it, rather than swallowing
					_ = form.SetField("Code", value)
				}
			}).
			Render(),
	)

	// Name - Focus Index 1
	name := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(2)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(10)},
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
				if err := form.SetField("Name", value); err != nil {
					// surface to the form's error state so the user sees it, rather than swallowing
					_ = form.SetField("Name", value)
				}
			}).
			Render(),
	)

	// Alias - Focus Index 2
	alias := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(10)},
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
				if err := form.SetField("Alias", value); err != nil {
					_ = form.SetField("Alias", value)
				}
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
		widgets.GroupComponent(
			ctx,
			"group",
			util.IntToString(form.Values().GroupID),
			30,
			v.FocusIndex == 3,
			" : ",
			func(id, value string) {
				i, err := strconv.Atoi(value)
				if err != nil {
					retui.Debugf("Invalid Group ID: %v", err)
					return
				}
				if err := form.SetField("GroupID", util.StringToInt(value, 0)); err != nil {
					_ = form.SetField("GroupID", i)
				}
			},
		),
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
				if err := form.SetField("PartyType", value); err != nil {
					_ = form.SetField("PartyType", value)
				}
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
				_ = form.SetField("Description", value)
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
			retui.Text("Address 1", retui.NewStyle()),
		),
		components.TextInput().
			ID("addressLine1").
			Value(v.AddressLine1).
			Prefix(" : ").
			Focused(v.FocusIndex == 6).
			OnChange(func(id, value string) {
				if err := form.SetField("AddressLine1", value); err != nil {
					_ = form.SetField("AddressLine1", value)
				}
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
			retui.Text("Address 2", retui.NewStyle()),
		),
		components.TextInput().
			ID("addressLine2").
			Prefix(" : ").
			Value(v.AddressLine2).
			Focused(v.FocusIndex == 7).
			OnChange(func(id, value string) {
				if err := form.SetField("AddressLine2", value); err != nil {
					_ = form.SetField("AddressLine2", value)
				}

			}).
			Render(),
	)

	// Country - Focus Index 8
	country := retui.Box(
		retui.Props{Gap: 1, Width: retui.Grow(1)},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(20)},
			retui.NewStyle(),
			retui.Text("Country", retui.NewStyle()),
		),

		widgets.CountryComponent(
			ctx,
			"country",
			form.Values().CountryID,
			30,
			v.FocusIndex == 8,
			" : ",
			func(id, value string) {
				i, err := strconv.Atoi(value)
				if err != nil {
					retui.Debugf("Invalid Country ID: %v", err)
					return
				}
				if err := form.SetField("CountryID", value); err != nil {

					_ = form.SetField("CountryID", i)
				}
			},
		),
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
		widgets.StateComponent(
			ctx,
			"state",
			form.Values().CountryID,
			form.Values().StateID,
			30,
			v.FocusIndex == 9,
			" : ",
			func(id, value string) {
				i, err := strconv.Atoi(value)
				if err != nil {
					retui.Debugf("Invalid State ID: %v", err)
					return
				}
				if err := form.SetField("StateID", value); err != nil {

					_ = form.SetField("StateID", i)
				}
			},
		),
	)

	// City - Focus Index 10
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
			Focused(v.FocusIndex == 10).
			OnChange(func(id, value string) {
				if err := form.SetField("City", value); err != nil {
					_ = form.SetField("City", value)
				}
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
				if err := form.SetField("Pincode", value); err != nil {
					_ = form.SetField("Pincode", value)
				}
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
				if err := form.SetField("Phone", value); err != nil {
					_ = form.SetField("Phone", value)
				}
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
				if err := form.SetField("Mobile", value); err != nil {
					_ = form.SetField("Mobile", value)
				}
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
				if err := form.SetField("Email", value); err != nil {
					_ = form.SetField("Email", value)
				}
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
				if err := form.SetField("ContactPerson", value); err != nil {
					_ = form.SetField("ContactPerson", value)
				}
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
				if err := form.SetField("GSTRegistrationType", value); err != nil {
					_ = form.SetField("GSTRegistrationType", value)
				}
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
				if err := form.SetField("GSTIN", value); err != nil {
					_ = form.SetField("GSTIN", value)
				}
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
			ID("save").
			Label("Save").
			Focused(v.FocusIndex == 24).
			OnKeyPress(func(id string, key retui.Key) bool {
				if key.Code == retui.KeyEnter {
					retui.Debugf("value: %v", v)
					c.save(v)
					form.Reset()
					return true
				}
				return false
			}).
			Render(),
	)

	// Reset Button - Focus Index 25
	resetButton := retui.Box(
		retui.Props{},
		retui.NewStyle().
			Background(retui.Gray(2)).
			Foreground(retui.White),
		components.Button().
			ID("reset").
			Label("Reset").
			Focused(v.FocusIndex == 25).
			OnKeyPress(func(id string, key retui.Key) bool {
				if key.Code == retui.KeyEnter {
					form.Reset()
					return true
				}
				return false
			}).
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
		FixedWidth(retui.CurrentScreenWidth - 70).
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
					country,
					state,
					city,
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
			Width:     retui.Grow(1),
		},
		retui.NewStyle(),
		panel,
	)
}
