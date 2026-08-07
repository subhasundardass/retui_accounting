package widgets

import (
	"strconv"

	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/ent/state"
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
)

func StateComponent(
	ctx *context.AppContext,
	id string,
	countryID int,
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
			Order(ent.Asc("name")).
			All(ctx.Context)
		if err != nil {
			retui.Debugf("StateComponent: %v", err)
			return nil
		}
		setStates(list)
		return nil
	}, []any{countryID}) // re-fetch when countryID changes

	options := retui.UseMemo(func() []components.SelectOption {
		opts := make([]components.SelectOption, len(states))
		for i, s := range states {
			opts[i] = components.SelectOption{
				Label: s.Name,
				Value: strconv.Itoa(s.ID),
			}
		}
		return opts
	}, []any{states})

	return renderStateComponent(id, options, value, width, focus, onChange)
}

func renderStateComponent(
	id string,
	options []components.SelectOption,
	value int,
	width int,
	focus bool,
	onChange func(id, value string),
) retui.Element {

	return components.SelectDropdown().
		ID(id).
		Width(width).
		Options(options).
		Value(strconv.Itoa(value)).
		Focused(focus).
		OnFilter(func(filterID, query string) []components.SelectOption {
			return FilterOptions(options, query)
		}).
		OnChange(onChange).
		Render()
}
