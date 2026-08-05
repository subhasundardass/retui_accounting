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
	value int,
	width int,
	focus bool,
	onChange func(id, value string),
) retui.Element {

	countries, setCountries := retui.UseState([]*ent.Country{})

	retui.UseEffect(func() func() {
		retui.Debug("CountryComponent UseEffect called")

		list, err := ctx.DB.Client.Country.Query().All(ctx.Context)
		if err != nil {
			retui.Debugf("Error: %v", err)
			return nil
		}

		retui.Debugf("Loaded %d countries", len(list))
		setCountries(list)

		return nil
	}, []any{})

	return renderCountryComponent(countries, value, width, focus, onChange)
}

func renderCountryComponent(
	countries []*ent.Country,
	value int,
	width int,
	focus bool,
	onChange func(id, value string),
) retui.Element {

	// options := make([]components.SelectOption, len(countries))
	options := retui.UseMemo(func() []components.SelectOption {
		opts := make([]components.SelectOption, len(countries))

		for i, country := range countries {
			opts[i] = components.SelectOption{
				Label: country.Name,
				Value: strconv.Itoa(country.ID),
			}
		}

		return opts
	}, []any{countries})

	for i, country := range countries {
		options[i] = components.SelectOption{
			Label: country.Name,
			Value: strconv.Itoa(country.ID),
		}
	}

	retui.Debugf("renderComponent: %d countries -> %d options", len(countries), len(options))

	return components.SelectDropdown().
		ID("country").
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

func FilterOptions(
	options []components.SelectOption,
	query string,
) []components.SelectOption {

	if query == "" {
		return options
	}

	query = strings.ToLower(query)

	out := make([]components.SelectOption, 0)

	for _, option := range options {
		if strings.Contains(strings.ToLower(option.Label), query) ||
			strings.Contains(strings.ToLower(option.Value), query) {
			out = append(out, option)
		}
	}

	return out
}
