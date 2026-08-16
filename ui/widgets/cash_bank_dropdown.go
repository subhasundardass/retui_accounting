package widgets

import (
	"strconv"
	"strings"

	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/ent/ledger"
	"github.com/subhasundardass/retui/ent/ledger_group"
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
)

func CashBankDropdown(
	ctx *context.AppContext,
	id string,
	value string,
	width int,
	focus bool,
	onChange func(id, value string),
) retui.Element {
	ledgers, setLedgers := retui.UseState([]*ent.Ledger{})
	retui.UseEffect(func() func() {
		list, err := ctx.DB.Client.Ledger.Query().
			Where(ledger.HasGroupWith(ledger_group.CodeIn("CASH", "BANK"))).
			Order(ent.Asc("name")).
			All(ctx.Context)
		if err != nil {
			retui.Debugf("LedgerComponent: %v", err)
			return nil
		}
		setLedgers(list)
		return nil
	}, []any{})

	options := retui.UseMemo(func() []components.SelectOption {
		opts := make([]components.SelectOption, len(ledgers))
		for i, l := range ledgers {
			opts[i] = components.SelectOption{
				Label: l.Name,
				Value: strconv.Itoa(l.ID),
			}
		}
		return opts
	}, []any{ledgers})

	// Pass id and options — no ledgers, no duplicate UseMemo
	return renderCashBankComponent(id, options, value, width, focus, onChange)
}

func renderCashBankComponent(
	id string,
	options []components.SelectOption, // options, not ledgers
	value string,
	width int,
	focus bool,
	onChange func(id, value string),
) retui.Element {

	return components.SelectDropdown().
		ID(id).
		Width(width).
		Options(options).
		Value(value).
		Focused(focus).
		OnFilter(func(filterID, query string) []components.SelectOption {
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
