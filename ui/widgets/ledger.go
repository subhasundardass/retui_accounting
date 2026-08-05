package widgets

import (
	"strconv"
	"strings"

	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
)

func LedgerComponent(
	ctx *context.AppContext,
	id string,
	value string,
	width int,
	focus bool,
	onChange func(id, value string),
) retui.Element {

	ledgers, setLedgers := retui.UseState([]*ent.Ledger{})

	retui.UseEffect(func() func() {

		list, err := ctx.DB.Client.Ledger.
			Query().
			Order(ent.Asc("name")).
			All(ctx.Context)

		if err != nil {
			retui.Debugf("LedgerComponent: %v", err)
			return nil
		}

		setLedgers(list)

		return nil
	}, []any{})

	return renderLedgerComponent(
		ledgers,
		value,
		width,
		focus,
		onChange,
	)
}

func renderLedgerComponent(
	ledgers []*ent.Ledger,
	value string,
	width int,
	focus bool,
	onChange func(id, value string),
) retui.Element {

	options := make([]components.SelectOption, len(ledgers))

	for i, ledger := range ledgers {
		options[i] = components.SelectOption{
			Label: ledger.Name,
			Value: strconv.Itoa(ledger.ID),
		}
	}

	return components.SelectDropdown().
		ID("ledger").
		Width(width).
		Options(options).
		Value(value).
		Focused(focus).
		OnFilter(func(id, query string) []components.SelectOption {

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
		}).
		OnChange(onChange).
		Render()
}
