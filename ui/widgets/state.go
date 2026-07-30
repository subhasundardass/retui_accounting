package widgets

import (
	"strconv"
	"strings"

	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/ent/state"
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
)

func StateComponent(
	countryID int,
	ctx *context.AppContext,
	value int,
	width int,
	focus bool,
	onChange func(id, value string),
) retui.Element {

	states, setStates := retui.UseState([]*ent.State{})

	retui.UseEffect(func() func() {

		list, err := ctx.DB.Client.State.
			Query().
			Where(state.CountryIDEQ(countryID)).
			All(ctx.Context)

		if err != nil {
			retui.Debugf("Failed to load states: %v", err)
			return nil
		}

		setStates(list)

		return nil
	}, []any{countryID})

	return renderStateComponent(states, value, width, focus, onChange)
}

func renderStateComponent(
	states []*ent.State,
	value int,
	width int,
	focus bool,
	onChange func(id, value string),
) retui.Element {

	options := make([]components.SelectOption, len(states))

	for i, state := range states {
		options[i] = components.SelectOption{
			Label: state.Name,
			Value: strconv.Itoa(state.ID),
		}
	}

	retui.Debugf("renderComponent: %d states -> %d options", len(states), len(options))

	return components.SelectDropdown().
		ID("state").
		Width(width).
		OverlayAbsPos(40, 5).
		Options(options). // <-- explicitly push latest options every render
		OnFilter(func(id, query string) []components.SelectOption {

			retui.Debugf("OnFilter called, %d options available", len(options))

			if query == "" {
				return options
			}

			filtered := make([]components.SelectOption, 0)
			query = strings.ToLower(query)

			for _, option := range options {
				if strings.Contains(strings.ToLower(option.Label), query) ||
					strings.Contains(strings.ToLower(option.Value), query) {
					filtered = append(filtered, option)
				}
			}

			return filtered
		}).
		Value(strconv.Itoa(value)).
		Focused(focus).
		OnChange(onChange).
		Render()
}
