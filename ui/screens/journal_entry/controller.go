package journal_entry

import (
	"time"

	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/repository"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/retui/window"
)

type State struct {
	IsDirty  bool
	Errors   map[string]string
	IsLoaded bool
}

type Controller struct {
	window     *window.Window
	ctx        *context.AppContext
	repo       *repository.JournalRepository
	repoLedger *repository.LedgerRepository
	state      State // embed state directly
}

func NewController(
	ctx *context.AppContext,
	repo *repository.JournalRepository,
	repoLedger *repository.LedgerRepository,
) *Controller {
	return &Controller{
		ctx:        ctx,
		repo:       repo,
		repoLedger: repoLedger,
		state: State{
			Errors:   make(map[string]string),
			IsDirty:  false,
			IsLoaded: false,
		},
	}
}

func (c *Controller) GetState() State {
	return c.state
}

// ledgerOptionsFor seeds the dropdown's initial Options (shown before any
// filter text is typed) with the default top-10, making sure the
// currently-selected value is present so the input box doesn't fall back
// to displaying the placeholder instead of the actual ledger name.
// LedgerFilterOptions is called from OnFilter on every keystroke while the
// dropdown is open. query is the text the user has typed so far.
func (c *Controller) LedgerFilterOptions(query string) []components.SelectOption {
	ledgers, err := c.repoLedger.Search(c.ctx.Context, query, 10)
	if err != nil {
		retui.Debug("LedgerFilterOptions error: " + err.Error())
		return nil
	}

	opts := make([]components.SelectOption, len(ledgers))
	for i, l := range ledgers {
		opts[i] = components.SelectOption{Label: l.Name, Value: l.Code}
	}
	return opts
}

// LedgerSeedOptions builds the initial Options list shown when the dropdown
// opens with no filter text typed yet. currentValue is the ledger already
// selected for this field (line.LedgerName) — it's included even if it
// falls outside the default top-10, so the input box displays the real
// name instead of falling back to the placeholder.
func (c *Controller) LedgerSeedOptions(currentValue string) []components.SelectOption {
	ledgers, err := c.repoLedger.Default(c.ctx.Context, 10)
	if err != nil {
		retui.Debug("LedgerSeedOptions error: " + err.Error())
		return nil
	}

	opts := make([]components.SelectOption, 0, len(ledgers)+1)
	found := false
	for _, l := range ledgers {
		opts = append(opts, components.SelectOption{Label: l.Name, Value: l.Code})
		if l.Code == currentValue {
			found = true
		}
	}

	// currentValue isn't among the default 10 — look it up directly and
	// prepend it, so the select shows the real name instead of falling
	// back to the placeholder.
	if !found && currentValue != "" {
		results, err := c.repoLedger.Search(c.ctx.Context, currentValue, 1)
		if err != nil {
			retui.Debug("LedgerSeedOptions Search error: " + err.Error())
		} else if len(results) > 0 {
			selected := results[0]
			opts = append([]components.SelectOption{
				{Label: selected.Name, Value: selected.Code},
			}, opts...)
		}
	}

	// retui.Debug(fmt.Sprintf("LedgerSeedOptions currentValue=%q -> %d options", currentValue, len(opts)))

	return opts
}

// ===Journal Entry
type JournalEntry struct {
	VoucherNo   string
	VoucherDate string
	Reference   string
	Narration   string
	Lines       []JournalLine
}

func (c *Controller) SaveJournal(input JournalEntry) (*ent.Journal, error) {
	voucherDate, err := time.Parse("02/01/2006", input.VoucherDate)
	if err != nil {
		return nil, err
	}

	journal := &ent.Journal{
		Date:          voucherDate,
		VoucherType:   "JV",
		VoucherNo:     input.VoucherNo,
		VoucherDate:   voucherDate,
		ReferenceNo:   &input.Reference,
		Narration:     &input.Narration,
		JournalStatus: "DRAFT",
	}

	var lines []repository.JournalLineInput

	for _, l := range input.Lines {
		lines = append(lines, repository.JournalLineInput{
			LedgerCode:  l.LedgerCode,
			Debit:       l.Debit,
			Credit:      l.Credit,
			Description: l.Remarks,
		})
	}

	jrnl, err := c.repo.CreateNew(c.ctx.Context, journal, lines)
	if err != nil {
		return nil, err
	}

	retui.Infof("Journal %s saved successfully.", jrnl.VoucherNo)

	return jrnl, nil
}
