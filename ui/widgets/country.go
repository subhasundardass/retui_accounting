package widgets

import (
	"strconv"
	"strings"

	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
)

func CountryComponent(
	ctx *context.AppContext,
	id string,
	value int,
	width int,
	focus bool,
	onChange func(id, value string),
) retui.Element {

	countries, setCountries := retui.UseState([]*ent.Country{})

	retui.UseEffect(func() func() {
		list, err := ctx.DB.Client.Country.Query().
			Order(ent.Asc("name")).
			All(ctx.Context)
		if err != nil {
			retui.Debugf("CountryComponent: %v", err)
			return nil
		}
		setCountries(list)
		return nil
	}, []any{})

	options := retui.UseMemo(func() []components.SelectOption {
		opts := make([]components.SelectOption, len(countries))
		for i, c := range countries {
			opts[i] = components.SelectOption{
				Label: c.Name,
				Value: strconv.Itoa(c.ID),
			}
		}
		return opts
	}, []any{countries})

	return renderCountryComponent(id, options, value, width, focus, onChange)
}

func renderCountryComponent(
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

// FilterOptions is reusable across any component
func FilterOptions(
	options []components.SelectOption,
	query string,
) []components.SelectOption {
	if query == "" {
		return options
	}
	query = strings.ToLower(query)
	filtered := make([]components.SelectOption, 0)
	for _, option := range options {
		if strings.Contains(strings.ToLower(option.Label), query) ||
			strings.Contains(strings.ToLower(option.Value), query) {
			filtered = append(filtered, option)
		}
	}
	return filtered
}
