// Package ledger_group
package ledger_group

import (
	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/retui/window"
)

type EditGroup struct {
	Name        string
	Code        string
	Description string
}

type GroupEditComponents struct {
	controller *Controller
}

func NewGroupEditComponents(controller *Controller) *GroupEditComponents {
	return &GroupEditComponents{controller: controller}
}

// WindowEditGroup creates a modal window for editing a ledger group
func (c *GroupEditComponents) WindowEditGroup(group ent.Ledger_Group) *window.Window {
	// Create state using the component pattern
	state := &EditGroupState{
		name:        group.Name,
		code:        group.Code,
		description: group.Description,
		focusIndex:  0,
		totalFields: 4,
	}

	// Create the render function
	renderFn := func() retui.Element {
		return c.renderEditGroup(state, group)
	}

	// Create the window
	w := window.NewWindow().
		SetTitle("Edit Ledger Group").
		SetModal(true).
		Center().
		SetSize(80, 40).
		SetRenderFn(renderFn)

	// Handle keyboard events
	w.OnKeyPress(func(key retui.Key) bool {
		switch key.Code {
		case retui.KeyEscape:
			w.Close()
			return true
		case retui.KeyDown:
			state.focusIndex = (state.focusIndex + 1) % state.totalFields
			return true
		case retui.KeyUp:
			state.focusIndex = (state.focusIndex - 1 + state.totalFields) % state.totalFields
			return true
		}
		return false
	})

	return w
}

// renderEditGroup renders the edit form
func (c *GroupEditComponents) renderEditGroup(state *EditGroupState, group ent.Ledger_Group) retui.Element {
	isFocused := func(idx int) bool { return state.focusIndex == idx }

	// Main content - form
	content := c.editGroupForm(
		state.name, func(value string) { state.name = value },
		state.code, func(value string) { state.code = value },
		state.description, func(value string) { state.description = value },
		isFocused,
	)

	// Footer with keyboard shortcuts
	footer := retui.Box(
		retui.Props{
			Direction: retui.Row,
			Padding:   [4]int{1, 1, 0, 1},
			Justify:   retui.JustifySpaceBetween,
			Width:     retui.Grow(1),
		},
		retui.NewStyle(),
		retui.Text("ESC: Close   ↑/↓: Navigate", retui.NewStyle()),
		components.Button().
			ID("submit").
			Label("Submit").
			Focused(isFocused(3)).
			OnClick(func(id string) {
				// println("Button", id, "clicked!")
				c.controller.HandleSave(state)
			}).Render(),
	)

	// Combine content and footer
	return retui.Box(
		retui.Props{
			Direction: retui.Column,
			Width:     retui.Grow(1),
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Width: retui.Grow(1),
			},
			retui.NewStyle(),
			content,
		),
		retui.Box(
			retui.Props{
				Width: retui.Grow(1),
			},
			retui.NewStyle(),
			footer,
		),
	)
}

// editGroupForm renders the form fields
func (c *GroupEditComponents) editGroupForm(
	name string, setName func(string),
	code string, setCode func(string),
	description string, setDescription func(string),
	isFocused func(int) bool,
) retui.Element {
	return retui.Box(
		retui.Props{
			Direction: retui.Column,
			Width:     retui.Grow(1),
			Gap:       0,
			Padding:   [4]int{0, 1, 0, 1},
		},
		retui.NewStyle(),

		// Group Code
		retui.Box(
			retui.Props{Gap: 1, Direction: retui.Row, Width: retui.Grow(1)},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(15)},
				retui.NewStyle(),
				retui.Text("Code:", retui.NewStyle().Bold(true)),
			),
			retui.Box(
				retui.Props{Width: retui.Grow(1)},
				retui.NewStyle(),
				components.TextInput().
					ID("groupCode").
					Width(retui.Fixed(20).Value).
					Focused(isFocused(0)).
					Value(code).
					Style(retui.NewStyle().Bold(true)).
					OnChange(func(id string, value string) {
						setCode(value)
					}).
					Render(),
			),
		),

		// Group Name
		retui.Box(
			retui.Props{Gap: 1, Direction: retui.Row, Width: retui.Grow(1)},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(15)},
				retui.NewStyle(),
				retui.Text("Name:", retui.NewStyle().Bold(true)),
			),
			retui.Box(
				retui.Props{Width: retui.Grow(1)},
				retui.NewStyle(),
				components.TextInput().
					ID("groupName").
					Width(retui.Fixed(50).Value).
					Focused(isFocused(1)).
					Value(name).
					Style(retui.NewStyle().Bold(true)).
					OnChange(func(id string, value string) {
						setName(value)
					}).
					Render(),
			),
		),

		// Description
		retui.Box(
			retui.Props{Gap: 1, Direction: retui.Row, Width: retui.Grow(1)},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(15)},
				retui.NewStyle(),
				retui.Text("Description:", retui.NewStyle().Bold(true)),
			),
			retui.Box(
				retui.Props{Width: retui.Grow(1)},
				retui.NewStyle(),
				components.TextInput().
					ID("groupDescription").
					Focused(isFocused(2)).
					Width(retui.Fixed(50).Value).
					Value(description).
					Style(retui.NewStyle().Bold(true)).
					OnChange(func(id string, value string) {
						setDescription(value)
					}).
					Render(),
			),
		),
	)
}

// SaveGroup saves a ledger group to the database
// func (c *Controller) SaveGroup(group ent.Ledger_Group) error {
// 	// Validation
// 	if group.Name == "" {
// 		// return retui.NewError("Group name is required")
// 	}
// 	if group.Code == "" {
// 		// return retui.NewError("Group code is required")
// 	}

// 	// TODO: Save to database
// 	// Example: return c.repo.Save(group)

// 	return nil
// }

// LoadGroup loads a ledger group from the database
// func (c *Controller) LoadGroup(id int) (*ent.Ledger_Group, error) {
// 	// TODO: Load from database
// 	// Example: return c.repo.FindByID(id)

// 	return &ent.Ledger_Group{
// 		ID:          id,
// 		Name:        "Sample Group",
// 		Code:        "SG001",
// 		Description: "Sample description",
// 	}, nil
// }

// DeleteGroup deletes a ledger group from the database
// func (c *Controller) DeleteGroup(id int) error {
// 	// TODO: Delete from database
// 	// Example: return c.repo.Delete(id)

// 	return nil
// }
